package filter

import (
	"fmt"
	"testing"
)

func PrintList(msg string, r, s []uint64) {
	fmt.Print(msg)
	for _, ele := range r {
		fmt.Print(ele, ",")
	}
	fmt.Println()
	for _, ele := range s {
		fmt.Print(ele, ",")
	}
	fmt.Println()
}

// 基本功能测试
func TestBasic(t *testing.T) {
	tool := Init("127.0.0.1:6379", 100000, 0.1)

	key := "uid1"
	var l = []uint64{uint64(201805261420)}
	r, s, err := tool.Filter(key, l)
	if err == nil {
		PrintList("Before Update():", r, s)
	}

	n, err := tool.Update(key, l)
	if err == nil {
		fmt.Printf("Update %d success\n", n)
	}

	r, s, err = tool.Filter(key, l)
	if err == nil {
		PrintList("After Update():", r, s)
	}

	var ll = []uint64{uint64(2019052)}
	r, s, err = tool.Filter(key, ll)
	if err == nil {
		PrintList("Before Update():", r, s)
	}
}
