package game

import (
	"go/cmkj_server_go/util"

	"strconv"
	"strings"
)

//游戏结果
const (
	X  = 0
	Z  = 1
	H  = 2
	XD = 3
	ZD = 4
)

//@Description  结果转为数组
//@Param
func converResult(result string) []byte {
	tmp := strings.TrimSuffix(result, "^")
	res := strings.Split(tmp, "^")
	re := make([]byte, 0, len(res))
	for _, r := range res {
		i, _ := strconv.Atoi(r)
		re = append(re, byte(i))
	}
	return re
}

//@Description  牌面结果分析
//@Param
func getNewResult(poker string) []byte {
	res := strings.Split(poker, "@")
	if len(res) < 6 {
		util.Log.Errorf("poker num = %d\n", len(res))
		return nil
	}
	point := anPoker(res)
	re := make([]byte, 26, 26)
	//闲 和 庄
	if point[1] > point[2] {
		re[X] = 1
	} else if point[1] < point[2] {
		re[Z] = 1
	} else {
		re[H] = 1
	}
	//对子
	if getPokerPoint(res[0]) == getPokerPoint(res[1]) {
		re[XD] = 1
	}
	if getPokerPoint(res[3]) == getPokerPoint(res[4]) {
		re[ZD] = 1
	}
	//total := point[1] + point[2]
	if re[Z] == 1 && point[2] == 6 {
		re[25] = 1
	}
	return re
}

//@Description  分析庄闲牌面，计算点数
//@Param
func anPoker(poker []string) []int {
	point := make([]int, 7, 7)
	p := make([]int, len(poker), len(poker))
	for i, po := range poker {
		pi := getPokerPoint(po)
		if pi > 9 {
			pi = 0
		}
		p[i] = pi
	}
	point[1] = int((p[0] + p[1] + p[2]) % 10)
	point[2] = int((p[3] + p[4] + p[5]) % 10)
	point[3] = int((p[0] + p[1]) % 10)
	point[4] = int((p[3] + p[4]) % 10)
	return point
}

//@Description  牌面转换为点数
//@Param
func getPokerPoint(poker string) int {
	p, err := strconv.Atoi(poker)
	if err != nil {
		util.Log.Error(err)
	}
	return p / 10
}

//@Description  生成路单
//@Param
func converWay(res []byte) string {
	if len(res) < 4 {
		return ""
	}
	var c byte
	if res[Z] == 1 {
		c = 'a'
	} else if res[X] == 1 {
		c = 'e'
	} else if res[H] == 1 {
		c = 'i'
	}
	if res[XD] == 1 {
		c++
	}
	if res[ZD] == 1 {
		c += 2
	}
	return string(c)
}

//@Description  结果统计
//@Param
func waysCount(counts []int, result []byte) string {
	if len(result) < 7 {
		return ""
	}
	var sb []byte
	for i := 0; i < len(counts); i++ {
		counts[i] += int(result[i])
		sb = strconv.AppendInt(sb, int64(counts[i]), 10)
		sb = append(sb, '^')
	}
	return string(sb)
}

func waysCountOld(counts []int, oldResult []byte) {
	if len(oldResult) < 7 {
		return
	}
	for i := 0; i < len(counts); i++ {
		counts[i] -= int(oldResult[i])
	}
}

func waysCountWay(counts []int, ways string) string {
	for i := 0; i < len(counts); i++ {
		counts[i] = 0
	}
	wayss := []byte(ways)
	for _, a := range wayss {
		n := a - 97
		r := n / 4
		if r == 0 {
			counts[1]++
		} else if r == 1 {
			counts[0]++
		} else {
			counts[2]++
		}
		r = n % 4
		if r == 1 || r == 3 {
			counts[3]++
		}
		if r >= 2 {
			counts[4]++
		}
	}
	var sb []byte
	for _, c := range counts {
		sb = strconv.AppendInt(sb, int64(c), 10)
		sb = append(sb, '^')
	}
	return string(sb)
}
