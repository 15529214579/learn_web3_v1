// package main

// import(
// 	"fmt"
// 	"crypto/sha256"
// )

// func main() {
// 	data := "helloworld"
// 	for i:=0;i<1000000;i++{
// 		hash :=sha256.Sum256([]byte(data + string(i)))
// 		fmt.Printf("%x\n",hash[:])
// 	}

// }