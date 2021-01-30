package main

import (
	"fmt"
)

func main() {
	fmt.Println(getAxis(66, 12))
}

// 通过x和y生成excel坐标文本
func getAxis(x, y int) string {
	return fmt.Sprintf("%s%d", getX(x), y+1)
}

// 生成横坐标字母
func getX(x int) string {
	abc := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if x < 26 {
		return string(abc[x])
	}
	mod := (x + 1) % 26
	return getX(int((x+1)/26)-1) + string(abc[mod-1])
}
