package main

import (
	"fmt"
	"strings"
)

func main() {
	str1 := []string{"maxuefei", "nihao", "happy"}
	res := strings.Join(str1, " ")
	fmt.Printf("%s\n", res)
}
