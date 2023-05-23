package main

import (
	"files"
	"flag"
	"os"
	"path/filepath"
)

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
		files.PrintHelp()
		os.Exit(0)
	}

	//ensure root exists and is a directory
	if err := files.CheckDir(path); err != nil {
		panic(err)
	}

	//convert root to absolute
	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	//create a slice of subdirectories, get a count of files in path dir
	subdirSlice, filepathCount, err := files.GetSubDirSlice(path)
	if err != nil {
		panic(err)
	}

	var results []files.SubDirResult

	//process the subdirectories
	if len(subdirSlice) > 0 {
		results = files.ProcessSubdirs(subdirSlice, workers, verbose)
	}

	//append the path to the results
	results = append(results, files.SubDirResult{path, filepathCount, 0})

	//sort the results by file-count
	sortedSubdirMap := files.SortSubDirMapByCount(results)

	//print the results
	files.PrintSortedMap(sortedSubdirMap, files.GetTotalPathCount(results), path)

	//write the tsv report
	if report {
		if err := files.WriteReport(sortedSubdirMap, outputFile); err != nil {
			panic(err)
		}
	}

}
