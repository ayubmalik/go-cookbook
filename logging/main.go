package main

import (
	"io"
	"log"
	"os"
)

func main() {

	log.SetFlags(log.Ldate)
	log.Println("Using log.Ldate")

	// the time in the local time zone: 01:23:23
	log.SetFlags(log.Ltime)
	log.Println("Using log.Ltime")

	// microsecond resolution: 01:23:23.123123
	log.SetFlags(log.Lmicroseconds)
	log.Println("Using log.Lmicroseconds")

	// full file name and line number: /a/b/c/d.go:23
	log.SetFlags(log.Llongfile)
	log.Println("Using log.Llongfile")

	// final file name element and line number: d.go:23
	log.SetFlags(log.Lshortfile)
	log.Println("Using log.Lshortfile")

	// if Ldate or Ltime is set, use UTC rather than the local time zone
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	log.Println("Using log.LUTC")

	// initial values for the standard logger
	log.SetFlags(log.LstdFlags)
	log.Println("Using log.LstdFlags")

	// set prefix
	log.SetPrefix("INFO ")
	log.Println("Using a prefix at beginning of line")

	// set prefix
	log.SetPrefix("INFO ")
	log.SetFlags(log.Lmsgprefix | log.LstdFlags)
	log.Println("Using a prefix before actual message")

	// log to file only. If it exists will append
	f, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("hello, world (file)")
	log.Println("hello, world (file)")
	log.Println("hello, world (file)")

	// log to a file and standard out
	w := io.MultiWriter(os.Stdout, f)
	log.SetOutput(w)
	log.Println("hello, world (both)")
	log.Println("hello, world (both)")
	log.Println("hello, world (both)")
}
