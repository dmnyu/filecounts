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

var (
	path       string
	verbose    bool
	outputFile string
	report     bool
)

func init() {
	flag.StringVar(&path, "path", "", "")
	flag.BoolVar(&verbose, "verbose", false, "")
	flag.BoolVar(&report, "report", false, "")
	flag.StringVar(&outputFile, "output-file", "filecounts.tsv", "")
}

func main() {
	flag.Parse()

	//ensure root exists and is a directory
	if err := checkDir(path); err != nil {
		panic(err)
	}

	//convert root to absolute
	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	if verbose {
		fmt.Printf("Counting files in subdirectories in %s\n", path)
	}

	subdirMap, err := getSubDirMap()
	if err != nil {
		panic(err)
	}
	if verbose {
		fmt.Printf("%v\n", subdirMap)
	}
	sortedSubdirMap := sortSubDirMapByCount(subdirMap)

	if verbose {
		fmt.Printf("%v\n", sortedSubdirMap)
	}

	if report {
		if err := writeReport(sortedSubdirMap); err != nil {
			panic(err)
		}
	}
	printSortedMap(sortedSubdirMap)

}

func writeReport(sortedMap map[int][]string) error {
	reportData := ""
	reportData = reportData + "file count\tpath\n"
	keys := []int{}
	for k, _ := range sortedMap {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	for i := len(keys) - 1; i > -1; i-- {
		key := keys[i]
		for _, p := range sortedMap[key] {
			reportData = reportData + fmt.Sprintf("%d\t%s\n", key, p)
		}
	}

	if err := os.WriteFile(outputFile, []byte(reportData), 0755); err != nil {
		return err
	}
	return nil
}

func printSortedMap(sortedMap map[int][]string) {
	keys := []int{}
	for k, _ := range sortedMap {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	fmt.Println("num files\tpath")
	fmt.Println("---------\t----")
	for i := len(keys) - 1; i > -1; i-- {
		key := keys[i]
		for _, p := range sortedMap[key] {
			fmt.Printf("%d\t\t%s\n", key, p)
		}
	}
}

func sortSubDirMapByCount(subdirMap map[string]int) map[int][]string {
	subdirsSorted := make(map[int][]string)
	for p, c := range subdirMap {
		if contains(c, subdirsSorted) {
			subdirsSorted[c] = append(subdirsSorted[c], p)
		} else {
			subdirsSorted[c] = []string{p}
		}
	}
	return subdirsSorted
}

func getTotalPathCount(subdirMap map[string]int) int {
	count := 0
	for _, v := range subdirMap {
		count = count + v
	}
	return count
}

func getSubDirMap() (map[string]int, error) {
	subdirMap := make(map[string]int)
	subdirs, err := os.ReadDir(path)
	if err != nil {
		return subdirMap, err
	}

	pathFileCount := 0
	for _, subdir := range subdirs {
		if subdir.IsDir() {
			subdirPath := filepath.Join(path, subdir.Name())
			subdirMap[subdirPath], err = getCount(subdirPath)
			if err != nil {
				return subdirMap, err
			}
		} else {
			pathFileCount = pathFileCount + 1
		}
	}
	subdirMap[path] = pathFileCount
	return subdirMap, nil
}

func checkDir(path string) error {
	fi, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return (err)
	} else if err != nil {
		return (err)
	} else if !fi.IsDir() {
		return (fmt.Errorf("%s is a not a directory\n", path))
	}
	return nil
}

func contains(i int, m map[int][]string) bool {
	for k, _ := range m {
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
