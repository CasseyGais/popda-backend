package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	password := "Popda2026"
	hash := sha256.Sum256([]byte(password))
	hashStr := hex.EncodeToString(hash[:])

	fmt.Println("Password:", password)
	fmt.Println("SHA256 Hash:", hashStr)
	fmt.Println("Length:", len(hashStr), "bytes")
	fmt.Println("\nCopy hash ini ke SQL:")
	fmt.Println(hashStr)
}
