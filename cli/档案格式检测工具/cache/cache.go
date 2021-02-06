package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("HEllo world")
}

func isDigit(s string) bool {
	for _, a := range s {
		i := int(a)
		println(i)
		if i < 48 || i > 57 {
			return false
		}
	}
	return true
}

func verifyID(id string) bool {
	a := [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	b := [11]string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}
	sum := 0
	for i, v := range id[:17] {
		n, err := strconv.Atoi(string(v))
		if err != nil {
			return false
		}
		sum += a[i] * n
	}
	if strings.ToUpper(string(id[17])) != b[sum%11] {
		return false
	}
	return true
}
