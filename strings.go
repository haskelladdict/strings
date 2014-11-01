package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"unicode"
)

var (
	minStrLen  int
	incLetters bool
	incNumbers bool
	incSpace   bool
	incPunct   bool
	offset     bool
)

func init() {
	flag.IntVar(&minStrLen, "l", 5, "minimum length of unicode strings to be printed")
	flag.BoolVar(&incLetters, "L", true, "include letters in printable strings")
	flag.BoolVar(&incNumbers, "N", true, "include numbers in printable strings")
	flag.BoolVar(&incSpace, "S", true, "include space characters in printable strings")
	flag.BoolVar(&incPunct, "P", true, "include punctuation characters in printable strings")
	flag.BoolVar(&offset, "o", false, "write the offset of each string from the "+
		"beginning of the file")
}

// parser takes a reader and parser it for consecutive stretches of valid
// unicode characters. The provided function isValid determines if a rune is
// valid or not. Parser writes the detected stretches to the provided io.Writer.
func parser(reader *bufio.Reader, w *bufio.Writer, isValid func(rune) bool) {
	var r rune
	var err error
	var str []rune
	var off, n int
	for {
		r, n, err = reader.ReadRune()
		if err != nil {
			return
		}
		off += n

		if !isValid(r) {
			if len(str) >= minStrLen {
				if offset {
					fmt.Fprintf(w, "%d: %s\n", off, string(str))
				} else {
					fmt.Fprintf(w, "%s\n", string(str))
				}
			}
			str = []rune{}
			continue
		}
		str = append(str, r)
	}

}

// createValidator generates a validator function for encountered input runs
// based on the provided command line options
func createValidator() func(rune) bool {
	return func(r rune) bool {
		if incLetters && unicode.IsLetter(r) {
			return true
		} else if incNumbers && unicode.IsNumber(r) {
			return true
		} else if incSpace && unicode.IsSpace(r) {
			return true
		} else if incPunct && unicode.IsPunct(r) {
			return true
		}
		return false
	}
}

// main entry point
func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		usage()
		os.Exit(1)
	}

	for _, f := range flag.Args() {
		file, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		defer file.Close()

		r := bufio.NewReader(file)
		w := bufio.NewWriter(os.Stdout)
		parser(r, w, createValidator())
		w.Flush()
		file.Close()
	}
}

// usage prints a short usage string to stdout
func usage() {
	fmt.Println("strings v0.1 (C) 2014 Markus Dittrich")
	fmt.Println()
	fmt.Println("usage: strings [options] filename")
	fmt.Println("options:")
	flag.PrintDefaults()
}
