package main

import (
	"files"
	"flag"
	"fmt"
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

func PrintHelp() {
	fmt.Printf("usage: %s --path path_to_directory [options]\n", os.Args[0])
	fmt.Println("Options:")
	fmt.Println("  --help\tprint this help message")
	fmt.Println("  --report\toutput a tsv file listing")
	fmt.Printf("  --output-file\tname of the report to create, default: %s\n", outputFile)
	fmt.Println("  --verbose\toutput verbose messages")
	fmt.Println("  --workers\tnumber of threads to run")
}

func main() {
	flag.Parse()

	if help {
		PrintHelp()
		os.Exit(0)
	}

	//ensure root exists and is a directory
	if err := files.CheckDir(path); err != nil {
		fmt.Printf(err.Error())
		PrintHelp()
		os.Exit(1)
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
