package util

import (
	"math/rand"
	"strconv"
	"time"
)

func GetSixRandomCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Intn(900000) + 100000)
}
