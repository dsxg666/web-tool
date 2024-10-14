package api

import (
	"math"
	"strconv"
	"time"

	"github.com/dsxg666/web-tool/internal/model"
	"github.com/dsxg666/web-tool/pkg/result"
	"github.com/gin-gonic/gin"
)

type DailyData struct {
	Num  int    `json:"num"`
	Rate string `json:"rate"`
	Res  string `json:"res"`
}

type ChartData struct {
	Date string `json:"date"`
	Num  int    `json:"num"`
}

type Statistic struct{}

func NewStatistic() Statistic {
	return Statistic{}
}

func (Statistic) GetChartMessageData(c *gin.Context) {
	today := time.Now().Truncate(24 * time.Hour)
	dates := make([]string, 7)
	for i := 6; i >= 0; i-- {
		date := today.AddDate(0, 0, -i).Format("2006-01-02")
		dates[6-i] = date
	}

	m := &model.Messages{}
	gm := &model.GroupMessages{}
	messageData := m.GetPastWeekData(dates)
	groupMessageData := gm.GetPastWeekData(dates)

	var totalMessageChartDatas []*ChartData

	for i := 0; i < 7; i++ {
		totalMessageChartDatas = append(totalMessageChartDatas, &ChartData{Date: dates[i], Num: messageData[i] + groupMessageData[i]})
	}

	c.JSON(200, result.SuccessWithData(totalMessageChartDatas))
}

func (Statistic) GetChartTodolistData(c *gin.Context) {
	today := time.Now().Truncate(24 * time.Hour)
	dates := make([]string, 7)
	for i := 6; i >= 0; i-- {
		date := today.AddDate(0, 0, -i).Format("2006-01-02")
		dates[6-i] = date
	}

	t := &model.TodoList{}
	datas := t.GetPastWeekData(dates)

	var todolistChartDatas []*ChartData

	for i := 0; i < 7; i++ {
		todolistChartDatas = append(todolistChartDatas, &ChartData{Date: dates[i], Num: datas[i]})
	}

	c.JSON(200, result.SuccessWithData(todolistChartDatas))
}

func (Statistic) GetChartDauData(c *gin.Context) {
	today := time.Now().Truncate(24 * time.Hour)
	dates := make([]string, 7)
	for i := 6; i >= 0; i-- {
		date := today.AddDate(0, 0, -i).Format("2006-01-02")
		dates[6-i] = date
	}

	dau := &model.Dau{}
	datas := dau.GetPastWeekData(dates)

	var dauChartDatas []*ChartData

	for i := 0; i < 7; i++ {
		dauChartDatas = append(dauChartDatas, &ChartData{Date: dates[i], Num: datas[i]})
	}

	c.JSON(200, result.SuccessWithData(dauChartDatas))
}

func (Statistic) GetDailyMessageData(c *gin.Context) {
	message := &model.Messages{}
	groupMessage := &model.GroupMessages{}
	todayCount := message.GetTodayMessageNum() + groupMessage.GetTodayGroupMessageNum()
	yesterdayCount := message.GetYesterdayMessageNum() + groupMessage.GetYesterdayGroupMessageNum()

	var dauRate string
	var dauRes string

	if yesterdayCount == 0 {
		if todayCount > 0 {
			dauRate = "100"
			dauRes = "1"
		} else {
			dauRate = "0"
			dauRes = "0"
		}
	} else {
		if todayCount > yesterdayCount {
			// 计算增长率
			growthRate := float64(todayCount-yesterdayCount) / float64(yesterdayCount) * 100
			dauRate = strconv.Itoa(int(math.Round(growthRate)))
			dauRes = "1"
		} else if todayCount < yesterdayCount {
			// 计算降低率
			declineRate := float64(yesterdayCount-todayCount) / float64(yesterdayCount) * 100
			dauRate = strconv.Itoa(int(math.Round(declineRate)))
			dauRes = "-1"
		} else {
			dauRate = "0" // 如果今天和昨天相同
			dauRes = "0"
		}
	}

	c.JSON(200, result.SuccessWithData(&DailyData{
		Num:  todayCount,
		Rate: dauRate + "%",
		Res:  dauRes,
	}))
}

func (Statistic) GetDailyTodolistData(c *gin.Context) {
	todolist := &model.TodoList{}
	todayCount := todolist.GetTodayTodolistNum()
	yesterdayCount := todolist.GetYesterdayTodolistNum()

	var dauRate string
	var dauRes string

	if yesterdayCount == 0 {
		if todayCount > 0 {
			dauRate = "100"
			dauRes = "1"
		} else {
			dauRate = "0"
			dauRes = "0"
		}
	} else {
		if todayCount > yesterdayCount {
			// 计算增长率
			growthRate := float64(todayCount-yesterdayCount) / float64(yesterdayCount) * 100
			dauRate = strconv.Itoa(int(math.Round(growthRate)))
			dauRes = "1"
		} else if todayCount < yesterdayCount {
			// 计算降低率
			declineRate := float64(yesterdayCount-todayCount) / float64(yesterdayCount) * 100
			dauRate = strconv.Itoa(int(math.Round(declineRate)))
			dauRes = "-1"
		} else {
			dauRate = "0" // 如果今天和昨天相同
			dauRes = "0"
		}
	}

	c.JSON(200, result.SuccessWithData(&DailyData{
		Num:  todayCount,
		Rate: dauRate + "%",
		Res:  dauRes,
	}))
}

func (Statistic) GetDailyDauData(c *gin.Context) {
	dau := &model.Dau{}
	todayCount := dau.GetTodayDauNum()
	yesterdayCount := dau.GetYesterdayDauNum()

	var dauRate string
	var dauRes string

	if yesterdayCount == 0 {
		if todayCount > 0 {
			dauRate = "100"
			dauRes = "1"
		} else {
			dauRate = "0"
			dauRes = "0"
		}
	} else {
		if todayCount > yesterdayCount {
			// 计算增长率
			growthRate := float64(todayCount-yesterdayCount) / float64(yesterdayCount) * 100
			dauRate = strconv.Itoa(int(math.Round(growthRate)))
			dauRes = "1"
		} else if todayCount < yesterdayCount {
			// 计算降低率
			declineRate := float64(yesterdayCount-todayCount) / float64(yesterdayCount) * 100
			dauRate = strconv.Itoa(int(math.Round(declineRate)))
			dauRes = "-1"
		} else {
			dauRate = "0" // 如果今天和昨天相同
			dauRes = "0"
		}
	}

	c.JSON(200, result.SuccessWithData(&DailyData{
		Num:  todayCount,
		Rate: dauRate + "%",
		Res:  dauRes,
	}))
}
