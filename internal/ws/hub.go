package ws

import "github.com/dsxg666/web-tool/global"

type Hub struct {
	Clients    map[*Client]bool // 活跃客户端集合
	Broadcast  chan []byte      // 广播消息通道
	Register   chan *Client     // 注册新客户端通道
	Unregister chan *Client     // 注销客户端通道
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			global.OnlineUser = append(global.OnlineUser, client.UserId)
			//global.Logger.Infof("用户进入: %v", client.UserId)
			//global.Logger.Infof("在线人数: %v", len(h.Clients))
			//global.Logger.Infof("global: %v", global.OnlineUser)
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			global.OnlineUser = removeElement(global.OnlineUser, client.UserId)
			//global.Logger.Infof("用户退出: %v", client.UserId)
			//global.Logger.Infof("在线人数: %v", len(h.Clients))
			//global.Logger.Infof("global: %v", global.OnlineUser)
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

func removeElement(slice []string, element string) []string {
	for i, v := range slice {
		if v == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
