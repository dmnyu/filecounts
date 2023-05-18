package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

var root string
var indexByIntMap = map[int][]string{}
var verbose bool

func init() {
	flag.StringVar(&root, "root", "", "")
	flag.BoolVar(&verbose, "verbose", false, "")
}

func main() {
	flag.Parse()

	//ensure root exists and is a directory
	fi, err := os.Stat(root)
	if errors.Is(err, os.ErrNotExist) {
		panic(err)
	} else if err != nil {
		panic(err)
	} else if !fi.IsDir() {
		panic(fmt.Errorf("%s is a not a directory\n", root))
	}

	//convert root to absolute
	root, err = filepath.Abs(root)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Counting files in subdirectories in %s\n", root)
	subdirs, err := os.ReadDir(root)
	if err != nil {
		panic(err)
	}

	subdirMap := make(map[string]int)
	rootCount := 0

	subdirCount := 0
	for _, subdir := range subdirs {
		if subdir.IsDir() {
			subdirCount = subdirCount + 1
			subdirMap[filepath.Join(root, subdir.Name())] = 0
		} else {
			rootCount = rootCount + 1
		}
	}

	currentCount := 1
	for k, _ := range subdirMap {
		fmt.Printf("  * counting files in %s (%d/%d)\n", k, currentCount, subdirCount)
		c, err := getCount(k)
		if err != nil {
			panic(err)
		} else {
			subdirMap[k] = c
		}
		currentCount = currentCount + 1
	}

	subdirMap[root] = rootCount

	for k, v := range subdirMap {
		if contains(v) {
			indexByIntMap[v] = append(indexByIntMap[v], k)
		} else {
			indexByIntMap[v] = []string{k}
		}
	}

	keys := []int{}
	for k, _ := range indexByIntMap {
		keys = append(keys, k)
	}

	fmt.Println("\n\t--- counts ---")
	sort.Ints(keys)
	for i := len(keys) - 1; i > -1; i-- {
		for _, subdir := range indexByIntMap[keys[i]] {
			fmt.Printf("%d\t%s\n", i, subdir)
		}
	}
}

func contains(i int) bool {
	for k, _ := range indexByIntMap {
		if k == i {
			return true
		}
	}
	return false
}

func getCount(path string) (int, error) {
	count := 0
	if err := filepath.Walk(path, func(obj string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			count = count + 1
			if verbose {
				fmt.Println("    * found ", obj)
			}
		}
		return nil
	}); err != nil {
		return 0, err
	}
	return count, nil
}
