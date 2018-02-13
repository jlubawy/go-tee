package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// main is the program entry point.
func main() {
	var (
		flagAppend bool
	)

	args := os.Args[1:]

	// Parse flags
	filesStart := 0
	for _, arg := range args {
		if len(arg) >= 2 {
			if arg[0] == '-' {
				if arg[1] == '-' {
					switch arg[2:] {
					case "help":
						usage()
					case "append":
						flagAppend = true
					default:
						unknownOption(arg[2:])
					}
				} else {
					switch arg[1] {
					case 'a':
						flagAppend = true
					default:
						unknownOption(arg[1:])
					}
				}

				filesStart += 1
			}
		}
	}

	// Open files
	files, err := openFiles(args[filesStart:], flagAppend)
	if err != nil {
		fatalf("Error opening file: %v\n", err)
	}
	defer closeFiles(files)

	// Scan for lines
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// Log to stdout
		fmt.Println(scanner.Text())

		// Then log to files
		for _, f := range files {
			fmt.Fprintln(f, scanner.Text())
		}
	}

	// Check for any errors
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}

// fatalf prints a message to stderr and then exits with error code 1.
func fatalf(format string, values ...interface{}) {
	fmt.Fprintf(os.Stderr, format, values...)
	os.Exit(1)
}

// openFiles opens all filenames provided (will append if flag is set). If an error occurs opening
// any file, the function will close all previously opened files and return the error.
func openFiles(filenames []string, flagAppend bool) (files []*os.File, err error) {
	files = make([]*os.File, 0, len(filenames))

	// On function return check for error, if there is any close
	// all previously opened files.
	defer func() {
		if err != nil {
			closeFiles(files)
		}
	}()

	flag := os.O_CREATE
	if flagAppend {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}

	for _, filename := range filenames {
		f, ferr := os.OpenFile(filename, flag, 0664)
		if ferr != nil {
			err = fmt.Errorf("opening filename '%s': %s", filename, err)
			return
		}

		files = append(files, f)
	}

	return
}

// closeFiles closes all files provided.
func closeFiles(files []*os.File) {
	for _, file := range files {
		file.Close()
	}
}

// usage prints the help message to stderr.
func usage() {
	fatalf(`Usage: go-tee [OPTION]... [FILE]...
Copy standard input to each FILE, and also to standard output.

  -a, --append   append to the given FILEs, do not overwrite
`)
}

// unknownOption prints a helpful message when an unknown option is encountered.
func unknownOption(arg string) {
	fatalf("go-tee: unknown option -- %s\nTry 'go-tee --help' for more information.\n", arg)
}
