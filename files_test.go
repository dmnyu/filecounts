package files

import (
	"path/filepath"
	"testing"
)

func TestFileCount(t *testing.T) {

	t.Run("test path exists and is a directory", func(t *testing.T) {
		path := "test-data/five-files"
		if err := CheckDir(path); err != nil {
			t.Error(err)
		}
	})

	t.Run("test count files in directory", func(t *testing.T) {
		path, err := filepath.Abs("test-data/five-files")
		if err != nil {
			t.Error(err)
		}

		got, err := getCount(path, 1, false)
		if err != nil {
			t.Error(err)
		}
		want := 5
		if want != got {
			t.Errorf("got: %d, wanted: %d", got, want)
		}
	})

	t.Run("test get Subdirectory slice", func(t *testing.T) {
		path := "test-data/multidirs"
		want := 2
		got, _, err := GetSubDirSlice(path)
		if err != nil {
			t.Error(err)
		}
		if want != len(got) {
			t.Errorf("got: %d, wanted: %d", len(got), want)
		}
	})

	t.Run("test empty directory handling", func(t *testing.T) {
		path := "test-data/empty-dir"
		_, _, err := GetSubDirSlice(path)
		if err == nil {
			t.Error(err)
		}

		path = "test-data/five-files"
		_, _, err = GetSubDirSlice(path)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("test getting dir counts", func(t *testing.T) {
		path := "test-data/multidirs"
		subdirSlice, pathCount, err := GetSubDirSlice(path)
		if err != nil {
			t.Error(err)
		}
		results := ProcessSubdirs(subdirSlice, 1, false)
		results = append(results, SubDirResult{path, pathCount, 0})
		got := len(results)
		want := 3
		if want != got {
			t.Errorf("got: %d, wanted: %d", got, want)
		}
	})

	t.Run("test getting file counts", func(t *testing.T) {
		path := "test-data/multidirs"
		subdirSlice, pathCount, err := GetSubDirSlice(path)
		if err != nil {
			t.Error(err)
		}
		results := ProcessSubdirs(subdirSlice, 1, false)
		results = append(results, SubDirResult{path, pathCount, 0})

		for _, result := range results {
			if result.Path == "test-data/multidirs" {
				if result.Count != 1 {
					t.Errorf("wanted %d, go %d", 1, result.Count)
				}
			} else {
				if result.Count != 5 {
					t.Errorf("wanted %d, go %d", 5, result.Count)
				}
			}
		}
	})

	t.Run("test get total count", func(t *testing.T) {
		path := "test-data/multidirs"
		subdirSlice, pathCount, err := GetSubDirSlice(path)
		if err != nil {
			t.Error(err)
		}
		results := ProcessSubdirs(subdirSlice, 1, false)
		results = append(results, SubDirResult{path, pathCount, 0})

		got := GetTotalPathCount(results)
		want := 11
		if want != got {
			t.Errorf("got: %d, wanted: %d", got, want)
		}
	})

	t.Run("test sort by count", func(t *testing.T) {
		path := "test-data/multidirs"
		subdirSlice, pathFileCount, err := GetSubDirSlice(path)
		if err != nil {
			t.Error(err)
		}
		results := ProcessSubdirs(subdirSlice, 1, false)
		results = append(results, SubDirResult{path, pathFileCount, 0})

		sortedSubdirMap := SortSubDirMapByCount(results)
		keys := sortKeys(sortedSubdirMap)

		if keys[0] != 5 {
			t.Errorf("wanted 5, got %d", keys[0])
		}

		if keys[1] != 1 {
			t.Errorf("wanted 1, got %d", keys[1])
		}

	})

}
