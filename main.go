package main

import (
	"densendither/cli"
	"fmt"
	"os"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
