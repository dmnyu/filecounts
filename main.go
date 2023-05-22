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

type SubDirResult struct {
	Path   string
	Count  int
	Result int
}

var (
	path       string
	verbose    bool
	outputFile string
	report     bool
	workers    int
	help       bool
)

func init() {
	flag.StringVar(&path, "path", "", "")
	flag.BoolVar(&verbose, "verbose", false, "")
	flag.BoolVar(&report, "report", false, "")
	flag.StringVar(&outputFile, "output-file", "filecounts.tsv", "")
	flag.IntVar(&workers, "workers", 8, "")
	flag.BoolVar(&help, "help", false, "")
}

func main() {
	flag.Parse()

	if help {
		printHelp()
		os.Exit(0)
	}

	//ensure root exists and is a directory
	if err := checkDir(path); err != nil {
		panic(err)
	}

	//convert root to absolute
	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	//create a slice of subdirectories, get a count of files in path dir
	subdirSlice, filepathCount, err := getSubDirSlice(path)
	if err != nil {
		panic(err)
	}

	var results []SubDirResult

	//process the subdirectories
	if len(subdirSlice) > 0 {
		results = processSubdirs(subdirSlice)
	}

	//append the path to the results
	results = append(results, SubDirResult{path, filepathCount, 0})

	//sort the results by file-count
	sortedSubdirMap := sortSubDirMapByCount(results)

	//print the results
	printSortedMap(sortedSubdirMap, getTotalPathCount(results))

	//write the tsv report
	if report {
		if err := writeReport(sortedSubdirMap); err != nil {
			panic(err)
		}
	}

}

func printHelp() {
	fmt.Printf("usage: %s --path path_to_directory [options]\n", os.Args[0])
	fmt.Println("Options:")
	fmt.Println("  --help\tprint this help message")
	fmt.Println("  --report\toutput a tsv file listing")
	fmt.Printf("  --output-file\tname of the report to create, default: %s\n", outputFile)
	fmt.Println("  --verbose\toutput verbose messages")
}

func checkDir(p string) error {
	if p == "" {
		printHelp()
		os.Exit(1)
	}
	fi, err := os.Stat(p)
	if errors.Is(err, os.ErrNotExist) {
		return (err)
	} else if err != nil {
		return (err)
	} else if !fi.IsDir() {
		return (fmt.Errorf("%s is a not a directory\n", p))
	}
	return nil
}

func getSubDirSlice(p string) ([]string, int, error) {
	subdirSlice := []string{}
	pathFileCount := 0
	subdirs, err := os.ReadDir(p)
	if err != nil {
		return subdirSlice, pathFileCount, err
	}

	if len(subdirs) == 0 {
		return subdirSlice, pathFileCount, fmt.Errorf("path: %s is an empty directory", p)
	}

	for _, subdir := range subdirs {
		if subdir.IsDir() {
			subdirSlice = append(subdirSlice, filepath.Join(path, subdir.Name()))
		} else {
			pathFileCount = pathFileCount + 1
		}
	}
	return subdirSlice, pathFileCount, nil
}

func processSubdirs(subdirs []string) []SubDirResult {
	workers = updateNumWorkers(len(subdirs))
	subdirChunks := splitSubdirs(subdirs)

	resultChannel := make(chan []SubDirResult)
	for i, chunk := range subdirChunks {
		go processSubdir(chunk, resultChannel, i+1)
	}

	results := []SubDirResult{}
	for range subdirChunks {
		chunk := <-resultChannel
		results = append(results, chunk...)
	}

	return results
}

func sortSubDirMapByCount(subdirResults []SubDirResult) map[int][]string {

	subdirsSorted := make(map[int][]string)

	for _, subdirResult := range subdirResults {
		if subdirResult.Result == 0 {
			if contains(subdirResult.Count, &subdirsSorted) {
				subdirsSorted[subdirResult.Count] = append(subdirsSorted[subdirResult.Count], subdirResult.Path)
			} else {
				subdirsSorted[subdirResult.Count] = []string{subdirResult.Path}
			}
		}
	}
	return subdirsSorted
}

func printSortedMap(sortedMap map[int][]string, totalCount int) {
	fmt.Printf("total number of files in %s: %d\n", path, totalCount)
	keys := sortKeys(sortedMap)

	fmt.Println("\nnum files\tpath")
	fmt.Println("---------\t----")
	for i := range keys {
		key := keys[i]
		paths := sortedMap[key]
		for _, p := range paths {
			fmt.Printf("%d\t\t%s\n", key, p)
		}

	}
	fmt.Println()
}

func writeReport(sortedMap map[int][]string) error {
	reportData := ""
	reportData = reportData + "file count\tpath\n"
	keys := sortKeys(sortedMap)

	for i := range keys {
		key := keys[i]
		paths := sortedMap[key]
		for _, p := range paths {
			reportData = reportData + fmt.Sprintf("%d\t%s\n", key, p)
		}
	}

	if err := os.WriteFile(outputFile, []byte(reportData), 0755); err != nil {
		return err
	}
	return nil
}

func getTotalPathCount(results []SubDirResult) int {
	count := 0
	for _, r := range results {
		count = count + r.Count
	}
	return count
}

func sortKeys(sortedMap map[int][]string) []int {
	keys := []int{}
	for k := range sortedMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	sortedKeys := []int{}
	for i := len(keys) - 1; i > -1; i-- {
		key := keys[i]
		sortedKeys = append(sortedKeys, key)
	}

	return sortedKeys
}

func contains(i int, m *map[int][]string) bool {
	for k := range *m {
		if k == i {
			return true
		}
	}
	return false
}

func updateNumWorkers(numSubdirs int) int {
	if numSubdirs < workers {
		return numSubdirs
	} else {
		return workers
	}
}

func processSubdir(subdirChunk []string, resultChannel chan []SubDirResult, workerID int) {
	subDirResults := []SubDirResult{}
	for _, subdir := range subdirChunk {
		if verbose {
			fmt.Printf("worker %d counting files in: %s\n", workerID, subdir)
		}
		count, err := getCount(subdir)
		if err != nil {
			subDirResults = append(subDirResults, SubDirResult{subdir, count, 1})
		} else {
			subDirResults = append(subDirResults, SubDirResult{subdir, count, 0})
		}
	}

	resultChannel <- subDirResults
}

func splitSubdirs(subdirs []string) [][]string {
	var divided [][]string

	chunkSize := (len(subdirs) + workers - 1) / workers

	for i := 0; i < len(subdirs); i += chunkSize {
		end := i + chunkSize

		if end > len(subdirs) {
			end = len(subdirs)
		}

		divided = append(divided, subdirs[i:end])
	}
	return divided
}

func getCount(p string) (int, error) {
	count := 0
	if err := filepath.Walk(p, func(obj string, info fs.FileInfo, err error) error {
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
