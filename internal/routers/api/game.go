package api

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

type Game struct{}

func NewGame() Game {
	return Game{}
}

var (
	colors      = [4]string{"♣", "♦", "♠", "♥"}
	numbers     = [13]string{"3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A", "2"}
	playingCard []string
)

func init() {
	// 初始化牌堆
	for _, color := range colors {
		for _, number := range numbers {
			playingCard = append(playingCard, color+number)
		}
	}
	playingCard = append(playingCard, "BigKing", "SmallKing")
}

// shuffle 洗牌功能
func shuffle(cards []string) []string {
	// 创建一个新的随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]string, len(cards))
	perm := r.Perm(len(cards)) // 随机排列索引
	for i, v := range perm {
		shuffled[i] = cards[v]
	}
	return shuffled
}

// deal 发牌功能
func deal(cards []string) (player1, player2, player3, landlord []string) {
	if len(cards) < 54 {
		panic("牌堆不完整")
	}

	// 玩家各17张，留3张作为地主牌
	player1 = cards[:17]
	player2 = cards[17:34]
	player3 = cards[34:51]
	landlord = cards[51:]
	return
}

func (Game) LD(c *gin.Context) {

}
