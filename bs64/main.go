package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	r := strings.NewReader("YQ==")
	NewDecoder(r, os.Stdout).Decode()
	fmt.Println()
}

