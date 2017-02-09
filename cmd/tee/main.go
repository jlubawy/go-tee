package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var filenames []string
var files []*os.File

func main() {
	filenames = os.Args[1:]
	files = make([]*os.File, len(filenames))

	for i, fname := range filenames {
		f, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		defer f.Close()

		files[i] = f
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		for _, f := range files {
			fmt.Fprintln(f, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}
