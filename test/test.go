package main

import (
	//"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type s1 struct {
	a int
}

func (s *s1) geta() int {
	return s.a
}

type A interface {
	geta() int
}

type s2 struct {
	s1
}

func (s *s2) geta() int {
	return s.s1.a * 2
}

type s3 struct {
	c int
	s A
}

type t1 struct {
	m map[string]int
}

func test1(data []int) {
	data = make([]int, 6, 6)
}

func test2(data []int) {
	data[1] = 6
}

func test3(data *[]int) *[]int {
	d := make([]int, 6, 6)
	data = &d
	return data
}

type ma struct {
	a int
	m map[int]string
}

func clear(s *ma) {
	s.m = make(map[int]string)
}
func swat(a, b int) {
	temp := a
	a = b
	b = temp
}

type zxcv struct {
	a *ma
}

func testsli() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	go func() {
		for _, v := range arr {
			fmt.Println(v)
			<-time.After(100 * time.Millisecond)
		}
	}()
}
func main() {
	ss := &s3{
		c: 1,
		s: &s2{
			s1: s1{
				a: 1,
			},
		},
	}
	fmt.Println(ss.s.geta())

	mm := make(map[string]int)
	mm["holly"] = 1
	mm["dml"] = 2
	t := t1{
		m: mm,
	}
	fmt.Println(t)
	ii, err := strconv.Atoi("1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ii)
	fmt.Println(int(1.0))
	fmt.Println(int(math.Floor(0.91)))
	now, err := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
	if err != nil {
		fmt.Println(err)
	}
	d := time.Time{}
	fmt.Println((&now).Unix())
	fmt.Println(d)
	res := "1^2^3"
	fmt.Println(strings.Trim(res, "^"))
	fmt.Println(strings.Split(res, "^")[1])
	fmt.Println(strings.Replace(res, "^", "", -1))
	//fmt.Println([]byte("1") - []byte("0"))
	poker := "94@12@73@@@103@133@22@@@"

	ttt := strings.Split(strings.TrimSuffix(poker, "@"), "@")

	for i, r := range ttt {
		fmt.Printf("ttt[%d]=%s\n", i, r)
	}
	fmt.Println(string(49))
	//fmt.Println(byte("49"))
	data := []int{1, 2, 3, 4, 5, 6}
	test1(data)
	fmt.Println(data)
	test2(data)
	fmt.Println(data)
	fmt.Println(test3(&data))
	fmt.Println(data)

	sli := make([]int, 6, 6)
	fmt.Println(sli[4])

	fmt.Println(strconv.Quote("a"))
	var bb byte
	bb = 'a'
	fmt.Println(bb + 1)
	//sli2 := []byte{2, 2, 2}
	//rrr := "hello"
	//sli = append(sli, sli2)
	fmt.Println("a" + "byy"[:1])
	fmt.Println('^')

	am := ma{
		a: 1,
		m: map[int]string{1: "dml", 2: "lzl"},
	}
	fmt.Println(am)
	clear(&am)
	fmt.Println(am)

	var x float64
	x = 1.2
	if x == 0 {
		fmt.Println(x)
	}
	a := 3
	b := 7
	swat(a, b)
	fmt.Println(a)
	fmt.Println(b)

	go func(a *int) {
		b := a
		fmt.Printf("1: [%d]\n", *b)
		time.Sleep(2 * time.Second)
		fmt.Printf("2: [%d]\n", *b)
	}(&a)

	go func(a *int) {
		time.Sleep(1 * time.Second)
		*a += 1
	}(&a)
	time.Sleep(3 * time.Second)

	sssss := make([]int, 6)
	fmt.Println(len(sssss))

	var bcd byte = 1
	fmt.Println(float64(bcd))

	var ffff float64 = 6.25388

	wessss := fmt.Sprintf("%.*f", 2, ffff)
	wesf, _ := strconv.ParseFloat(wessss, 64)
	fmt.Println(wesf)

	sss1 := []int{1, 2, 3, 4, 5, 6}
	sss2 := make([]int, 0, 6)
	copy(sss2, sss1)
	fmt.Println(sss2)
	//strconv.FormatFloat
	var fffff float64 = 0

	fmt.Println(fmt.Sprintf("%.2f", fffff))
	fmt.Println(strconv.FormatFloat(fffff, 'f', -1, 64))
	testsli()
	<-time.After(2 * time.Second)
	sad := zxcv{}
	fmt.Println(sad.a)
	for index := range sss1 {
		fmt.Println(index)
	}
	map1 := map[string]int{
		"lzl": 1,
		"ss":  2,
	}
	for k := range map1 {
		fmt.Println(k)
	}
	string1 := "1@2@3@@"
	fmt.Println(strings.TrimRight(string1, "@"))
	fmt.Println(strings.TrimSuffix(string1, "@"))

	string1 += "hello"
}
