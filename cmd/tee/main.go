package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	filenames := os.Args[1:]
	files := make([]*os.File, len(filenames))
	for i, fname := range filenames {
		f, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		files[i] = f
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Print(scanner.Text())
		for _, f := range files {
			fmt.Fprint(f, scanner.Text())
		}
	}
}