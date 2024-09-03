package pkg2

import (
	"fmt"
	"os"
)

// ExportMain is an exported function that can be accessed from other packages.
func ExportMain() {
	main()
	os.Exit(1)
}

func main() {
	os.Exit(15)
	Exit(15)
	callExit(1)
}

// Exit is not an exported function and is not accessible from other packages.
func Exit(i int) {
	fmt.Println(i)
}

func callExit(i int) {
	os.Exit(15)
}
