package util

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const MIN = 0.000001

func Round(value float64, scale int) float64 {
	v := fmt.Sprintf("%.*f", scale, value)
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		Log.Error(err)
		return 0
	}
	return f
}

func DecimalFormat(v float64) string {
	if v == 0 {
		return "0"
	}
	return fmt.Sprintf("%.2f", v)
}

func DecimalEqual(f1, f2 float64) bool {
	return math.Dim(f1, f2) < MIN
}

func SplitInt(src, sep string) ([]int, error) {
	src = strings.TrimSuffix(src, sep)
	strs := strings.Split(src, sep)
	var nums []int
	for _, v := range strs {
		num, err := strconv.Atoi(v)
		if err != nil {
			Log.Error(err)
			return nil, err
		}
		nums = append(nums, num)
	}
	return nums, nil
}
