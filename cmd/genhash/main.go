package main

import (
	"fmt"
	"mindblog/internal/auth"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run ./cmd/genhash <password>")
		os.Exit(1)
	}
	hash, err := auth.HashPassword(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println(hash)
}
