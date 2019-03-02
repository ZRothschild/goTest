/*
	规则比较
*/
package tj

import (
	"math/rand"
	"strconv"
	"time"
)

//是否是一对（宝）
type TjSlice []Tj

var expired []int

var self = make(map[string]TjSlice)

func (s TjSlice) Work(c TjSlice) int8 {
	cb, ci := s.ComCoup(c)
	if cb {
		return ci
	}
	kb, ki := s.ComKing(c)
	if kb {
		return ki
	}
	return s.ComNum(c)
}

//比较一对大小
func (s TjSlice) ComCoup(c TjSlice) (bool, int8) {
	//看双方是否为一对
	sCouple := IsCouple(s)
	cCouple := IsCouple(c)
	if sCouple && cCouple {
		if s[0].Level > c[0].Level {
			return true, 1
		} else {
			return true, 0
		}
	} else if sCouple == true {
		return true, 1
	} else if cCouple == true {
		return true, 0
	}
	return false, -1
}

//比较特殊大小
func (s TjSlice) ComKing(c TjSlice) (bool, int8) {
	sNum := IsKing(s)
	cNum := IsKing(c)
	if sNum > 0 && sNum > cNum {
		return true, 1
	} else if cNum > 0 && cNum > sNum {
		return true, 0
	} else if sNum > 0 {
		if (sNum == 1 || sNum == 2) && Com(s, c) == 0 {
			return true, 0
		}
		return true, 1
	}
	return false, -1
}

//杂项 点数对比
func (s TjSlice) ComNum(c TjSlice) int8 {
	var (
		sNum = (s[1].Number + s[0].Number) % 10
		cNum = (c[1].Number + c[0].Number) % 10
	)
	if sNum > cNum {
		return 1
	} else if sNum < cNum {
		return 0
	}
	return Com(s, c)
}

//是否是一对（宝）
func IsCouple(s TjSlice) bool {
	b := false
	if s[0].Level == s[1].Level {
		b = true
	}
	return b
}

//是否为天王：4 地王：3 天杠：2 地杠：1
func IsKing(s TjSlice) int8 {
	var level int8 = 0
	num := s[0].Number + s[1].Number
	if num == 21 && (s[0].Number == 12 || s[1].Number == 12) {
		level = 4
	} else if num == 11 && (s[0].Number == 2 || s[1].Number == 2) {
		level = 3
	} else if num == 20 && (s[0].Number == 12 || s[1].Number == 12) {
		level = 2
	} else if num == 10 && (s[0].Number == 2 || s[1].Number == 2) {
		level = 1
	}
	return level
}

//每一只牌对比
func Com(s TjSlice, c TjSlice) int8 {
	var (
		sMax = s[0].Flag
		sMin = s[1].Flag
		cMax = c[0].Flag
		cMin = c[1].Flag
	)
	if s[1].Flag > s[0].Flag {
		sMax = s[1].Flag
		sMin = s[0].Flag
	}
	if c[1].Flag > c[0].Flag {
		cMax = c[1].Flag
		cMin = c[0].Flag
	}
	if sMax > cMax {
		return 1
	} else if sMax < cMax {
		return 0
	} else if sMin > cMin {
		return 1
	} else if sMin < cMin {
		return 0
	}
	return 1
}

//每个人四张牌
func RandInt() map[string]TjSlice {
	rand.Seed(time.Now().UnixNano())
	selfName := rand.Intn(9999999)
	sliceInt := make(TjSlice, 0)
	for i := 0; i < 4; i++ {
		randInt := rands()
		sliceInt = append(sliceInt, Tjs[randInt])
		self[strconv.Itoa(selfName)] = sliceInt
	}
	return self
}

//32张牌随机选择一张，不可重复
func rands() int {
	rand.Seed(time.Now().UnixNano())
	randInt := rand.Intn(32)
	for _, v := range expired {
		if v == randInt {
			return rands()
		}
	}
	expired = append(expired, randInt)
	return randInt
}
