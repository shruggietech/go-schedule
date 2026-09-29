// Command brand-check verifies the pinned official archive and its repository consumers.
package main

import (
	"fmt"
	"os"

	"github.com/shruggietech/go-schedule/internal/brandrelease"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "brand-check: working directory: %v\n", err)
		os.Exit(1)
	}
	files, consumers, err := brandrelease.CheckInstalled(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "brand-check: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("brand-check: OK - %d official files, %d consumers\n", files, consumers)
}
