// Dup1 prints the text of each line that appears more than
// once in the standard input, preceded by its count.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Printf("%s", input)
}
