package main

import (
	"fmt"
	"os"

	qrcode "github.com/skip2/go-qrcode"
)

func main() {
	err := qrcode.Print(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
