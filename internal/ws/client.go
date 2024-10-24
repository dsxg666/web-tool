package ws

import (
	"github.com/dsxg666/web-tool/global"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 写入消息的超时时间
	writeWait = 10 * time.Second

	// 等待客户端 pong 响应的超时时间
	pongWait = 60 * time.Second

	// 发送 ping 消息的时间间隔
	pingPeriod = (pongWait * 9) / 10

	// 最大消息大小
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	UserId string
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c // 从 Hub 中注销客户端
		c.Conn.Close()        // 关闭 WebSocket 连接
	}()
	// 设置最大消息读取大小
	c.Conn.SetReadLimit(maxMessageSize)
	// 设置消息读取的超时时间
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	// 设置用户响应处理
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	lastActive := time.Now()

	// 启动一个 goroutine 来监控超时
	go func() {
		someTimeout := 3 * time.Minute
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if time.Since(lastActive) > someTimeout {
					global.Logger.Errorf("Client inactive for too long, closing connection")
					c.Hub.Unregister <- c
					c.Conn.Close()
					return
				}
			}
		}
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				global.Logger.Errorf("error: %v", err)
			}
			break
		}

		lastActive = time.Now()

		c.Hub.Broadcast <- message
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod) // 定期发送 ping 消息
	defer func() {
		ticker.Stop()  // 停止定时器
		c.Conn.Close() // 关闭连接
	}()
	for {
		select {
		case message, ok := <-c.Send: // 读取发送队列中的消息
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{}) // 连接关闭
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage) // 准备发送消息
			if err != nil {
				global.Logger.Errorf("error: %v", err)
				return
			}
			w.Write(message) // 写入消息

			n := len(c.Send) // 批量发送队列中的其他消息
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				global.Logger.Errorf("error: %v", err)
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				global.Logger.Errorf("error: %v", err)
				return
			}
		}
	}
}

func ServeWs(hub *Hub, c *gin.Context) {
	id := c.Query("id")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		global.Logger.Errorf("websocket upgrade error: %v", err)
		return
	}
	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256), UserId: id}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
