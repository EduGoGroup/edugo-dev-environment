package main

import (
	"fmt"
	"os"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run scripts/generate-password.go <password>")
		os.Exit(1)
	}
	password := os.Args[1]
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Hash:     %s\n", string(hash))
}
