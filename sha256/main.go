package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	input := "hello world\n"

	h, err := hash(strings.NewReader(input))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("hash = %x\n", h)
}

func hash(r io.Reader) ([]byte, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
