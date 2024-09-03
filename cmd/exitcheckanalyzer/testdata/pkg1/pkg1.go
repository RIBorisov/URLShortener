package main

import (
	"fmt"
	"os"
)

func main() {
	os.Exit(15) // want "calling os.Exit in function main"
	Exit(15)
	callExit(1)
}

func Exit(i int) {
	fmt.Println(i)
}

func callExit(i int) {
	os.Exit(15)
}
