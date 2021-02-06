package main

import (
	"fmt"
	"path"
)

func main() {
	dirpath := "/home/yu/go/src/new_file"
	file := "new_file"
	fmt.Println(path.Split(path.Clean(dirpath)))
	fmt.Println(path.Join(dirpath, file))
}
