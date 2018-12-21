package main

import (
	"fmt"
	"testing"

	pb "github.com/brotherlogic/dropboxsync/proto"
	"golang.org/x/net/context"
)

type dbTest struct {
	failCopy bool
	failList bool
	paths    map[string][]string
}

func (d *dbTest) copyFile(key string, origin, dest string) error {
	if d.failCopy {
		return fmt.Errorf("Built to fail")
	}
	return nil
}

func (d *dbTest) listFiles(key string, path string) ([]string, error) {
	if d.failList {
		return []string{}, fmt.Errorf("Failed to list")
	}
	if _, ok := d.paths[path]; !ok {
		return []string{}, nil
	}
	return d.paths[path], nil
}

func TestDiffFileList(t *testing.T) {
	files := diffFileList([]string{"file1", "file2", "file3"},
		[]string{"file1", "file5"})

	if len(files) != 1 && files[0] != "file5" {
		t.Errorf("Diff file list has failed: %v", files)
	}
}

func InitTest() *Server {
	s := Init()
	s.SkipLog = true
	return s
}

func TestSuccessList(t *testing.T) {
	s := InitTest()
	s.dropbox = &dbTest{paths: map[string][]string{"one": []string{"one", "two"}, "two": []string{"two"}}}

	s.runUpdate(context.Background(), &pb.SyncConfig{Origin: "one", Destination: "two"})
	if s.copies != 1 {
		t.Errorf("Not enough copying")
	}
}

func TestSuccessFail(t *testing.T) {
	s := InitTest()
	s.dropbox = &dbTest{failList: true, paths: map[string][]string{"one": []string{"one", "two"}, "two": []string{"two"}}}

	s.runUpdate(context.Background(), &pb.SyncConfig{Origin: "one", Destination: "two"})
	if s.copies != 0 {
		t.Errorf("Too much copying")
	}
}
