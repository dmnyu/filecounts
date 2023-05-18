package main

import (
	"path/filepath"
	"testing"
)

func TestFileCount(t *testing.T) {

	t.Run("test path exists and is a directory", func(t *testing.T) {
		path := "test-data/five-files"
		if err := checkDir(path); err != nil {
			t.Error(err)
		}
	})

	t.Run("test files in test-data", func(t *testing.T) {
		path, err := filepath.Abs("test-data/five-files")
		if err != nil {
			t.Error(err)
		}

		got, err := getCount(path)
		if err != nil {
			t.Error(err)
		}
		want := 5
		if want != got {
			t.Errorf("got: %d, wanted: %d", got, want)
		}
	})

	t.Run("test get Subdirectory slice", func(t *testing.T) {
		path = "test-data/multidirs"
		want := 2
		got, _, err := getSubDirSlice(path)
		if err != nil {
			t.Error(err)
		}
		if want != len(got) {
			t.Errorf("got: %d, wanted: %d", len(got), want)
		}
	})

	t.Run("test getting dir counts", func(t *testing.T) {
		path = "test-data/multidirs"
		subdirSlice, pathCount, err := getSubDirSlice(path)
		if err != nil {
			t.Error(err)
		}
		results := processSubdirs(subdirSlice)
		results = append(results, SubDirResult{path, pathCount, 0})
		got := len(results)
		want := 3
		if want != got {
			t.Errorf("got: %d, wanted: %d", got, want)
		}
	})

	t.Run("test getting file counts", func(t *testing.T) {
		path = "test-data/multidirs"
		subdirSlice, pathCount, err := getSubDirSlice(path)
		if err != nil {
			t.Error(err)
		}
		results := processSubdirs(subdirSlice)
		results = append(results, SubDirResult{path, pathCount, 0})

		for _, result := range results {
			want := 5
			got := result.Count
			if want != got {
				t.Errorf("got: %d, wanted: %d", got, want)
			}
		}
	})
	/*

		t.Run("test get total count", func(t *testing.T) {
			path = "test-data/multidirs"
			subdirMap, err := getSubDirMap()
			if err != nil {
				t.Error(err)
			}

			got := getTotalPathCount(subdirMap)
			want := 15
			if want != got {
				t.Errorf("got: %d, wanted: %d", got, want)
			}
		})

		t.Run("test sort by count", func(t *testing.T) {
			path = "test-data/multidirs"
			subdirMap, err := getSubDirMap()
			if err != nil {
				t.Error(err)
			}
			sortedMap := sortSubDirMapByCount(subdirMap)
			want := 1
			got := len(sortedMap)
			if want != got {
				t.Errorf("wanted: %d, got %d\n", want, got)
			}
		})

	*/

}
