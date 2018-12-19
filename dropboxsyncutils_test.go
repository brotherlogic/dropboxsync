package main

import (
	"testing"
)

func TestDiffFileList(t *testing.T) {
	files := diffFileList([]string{"file1", "file2", "file3"},
		[]string{"file1", "file5"})

	if len(files) != 1 && files[0] != "file5" {
		t.Errorf("Diff file list has failed: %v", files)
	}
}
