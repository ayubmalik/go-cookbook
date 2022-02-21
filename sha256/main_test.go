package main

import (
	"fmt"
	"os"
	"testing"
)

func TestHash(t *testing.T) {

	t.Run("from file", func(t *testing.T) {
		f, err := os.Open("testdata/example.input")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		h, err := hash(f)
		if err != nil {
			t.Fatal(err)
		}

		got := fmt.Sprintf("%x", h)
		want := "0a1b2b0113d7a52992d4114128c285d864a970413abbcce8555f84a0c10a373f"
		if got != want {
			t.Errorf("got: %s, but want %s", got, want)
		}
	})

}
