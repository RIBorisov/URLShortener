package pkg2

import (
	"fmt"
	"os"
)

func ExportMain() {
	main()
	os.Exit(1)
}

func main() {
	os.Exit(15)
	Exit(15)
	callExit(1)
}

func Exit(i int) {
	fmt.Println(i)
}

func callExit(i int) {
	os.Exit(15)
}
