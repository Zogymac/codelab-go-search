package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	recursiveFlag = flag.Bool("r", true, "recursive search: for directories")
	nFlag         = flag.Bool("n", false, "line number on which the entry was found")
)

type ScanResult struct {
	file       string
	lineNumber []int
	line       []string
}

func scanFile(fpath, pattern string) (ScanResult, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return ScanResult{}, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	result := make([]string, 0)
	index := make([]int, 0)
	i := 0 //0
	for scanner.Scan() {
		i++
		line := scanner.Text()
		contains := strings.Contains(line, pattern)
		if contains {
			index = append(index, i)
			result = append(result, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return ScanResult{}, err
	}
	//return ScanResult{}, nil
	scanEnd := ScanResult{fpath, index, result}
	return scanEnd, nil
}

func exit(format string, val ...interface{}) {
	if len(val) == 0 {
		fmt.Println(format)
	} else {
		fmt.Printf(format, val)
		fmt.Println()
	}
	os.Exit(1)
}

func processFile(fpath string, pattern string) {
	res, err := scanFile(fpath, pattern)
	if err != nil {
		exit("Error scanning %s: %s", fpath, err.Error())
	}
	/*
		for _, line := range res {
			fmt.Println(line)
		}*/
	for i, _ := range res.line {
		if *nFlag {
			fmt.Print(fpath, ":", strconv.Itoa(res.lineNumber[i]), ":", res.line[i], "\n")
		} else {
			fmt.Print(fpath, ":", res.line[i], "\n")
		}
	}
}
func walkFunc(path string, info os.FileInfo, err error, pattern string) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	fmt.Println(info.Name())
	scanFile(path, pattern)

	return nil
}
func processDirectory(dir string, pattern string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//res, nil := scanFile(dir, pattern)
		if !info.IsDir() {
			//fmt.Println(path)
			processFile(path, pattern)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	if flag.NArg() < 2 {
		exit("usage: go-search <path> <pattern> to search")
	}

	path := flag.Arg(0)
	pattern := flag.Arg(1)

	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	recursive := *recursiveFlag
	if info.IsDir() && !recursive {
		exit("%s: is a directory", info.Name())
	}

	if info.IsDir() && recursive {
		processDirectory(path, pattern)
	} else {
		processFile(path, pattern)
	}
}
