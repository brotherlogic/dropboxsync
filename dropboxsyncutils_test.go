package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

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
	files := diffFileList([]string{"file1", "file2", "file3", "somefolder/file4"},
		[]string{"file1", "file4", "file5"})

	if len(files) != 1 && files[0] != "file5" {
		t.Errorf("Diff file list has failed: %v", files)
	}
}

func TestDiffFileListPreStrip(t *testing.T) {
	files := diffFileListPreStrip([]string{"file1", "file2", "file3", "somefolder/file4"},
		[]string{"file1", "file4", "file5"})

	if len(files) != 1 && files[0] != "file5" {
		t.Errorf("Diff file list has failed: %v", files)
	}
}

func TestDiffFileListClever(t *testing.T) {
	files := diffFileListClever([]string{"file1", "file2", "file3", "somefolder/file4"},
		[]string{"file1", "file4", "file5"})

	if len(files) != 1 && files[0] != "file5" {
		t.Errorf("Diff file list has failed: %v", files)
	}
}

func benchmarkDiff(strlen int, b *testing.B) {
	// Create two 1000 piece arrays
	files := []string{}

	for i := 0; i < strlen; i++ {
		str := fmt.Sprintf("stringington-%v", i)
		files = append(files, str)
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	shuf1 := make([]string, len(files))
	perm := r.Perm(len(shuf1))
	for i, randIndex := range perm {
		shuf1[i] = files[randIndex]
	}

	shuf2 := make([]string, len(files))
	perm = r.Perm(len(shuf2))
	for i, randIndex := range perm {
		shuf2[i] = files[randIndex]
	}

	for n := 0; n < b.N; n++ {
		diffFileList(shuf2, shuf2)
	}
}

func benchmarkDiffPreStrip(strlen int, b *testing.B) {
	// Create two 1000 piece arrays
	files := []string{}

	for i := 0; i < strlen; i++ {
		str := fmt.Sprintf("stringington-%v", i)
		files = append(files, str)
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	shuf1 := make([]string, len(files))
	perm := r.Perm(len(shuf1))
	for i, randIndex := range perm {
		shuf1[i] = files[randIndex]
	}

	shuf2 := make([]string, len(files))
	perm = r.Perm(len(shuf2))
	for i, randIndex := range perm {
		shuf2[i] = files[randIndex]
	}

	for n := 0; n < b.N; n++ {
		diffFileListPreStrip(shuf2, shuf2)
	}
}

func benchmarkDiffClever(strlen int, b *testing.B) {
	// Create two 1000 piece arrays
	files := []string{}

	for i := 0; i < strlen; i++ {
		str := fmt.Sprintf("stringington-%v", i)
		files = append(files, str)
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	shuf1 := make([]string, len(files))
	perm := r.Perm(len(shuf1))
	for i, randIndex := range perm {
		shuf1[i] = files[randIndex]
	}

	shuf2 := make([]string, len(files))
	perm = r.Perm(len(shuf2))
	for i, randIndex := range perm {
		shuf2[i] = files[randIndex]
	}

	for n := 0; n < b.N; n++ {
		diffFileListClever(shuf2, shuf2)
	}
}

func BenchmarkDiff1(b *testing.B)    { benchmarkDiff(1, b) }
func BenchmarkDiff10(b *testing.B)   { benchmarkDiff(10, b) }
func BenchmarkDiff100(b *testing.B)  { benchmarkDiff(100, b) }
func BenchmarkDiff1000(b *testing.B) { benchmarkDiff(1000, b) }

func BenchmarkDiffPS1(b *testing.B)    { benchmarkDiffPreStrip(1, b) }
func BenchmarkDiffPS10(b *testing.B)   { benchmarkDiffPreStrip(10, b) }
func BenchmarkDiffPS100(b *testing.B)  { benchmarkDiffPreStrip(100, b) }
func BenchmarkDiffPS1000(b *testing.B) { benchmarkDiffPreStrip(1000, b) }

func BenchmarkDiffC1(b *testing.B)    { benchmarkDiffClever(1, b) }
func BenchmarkDiffC10(b *testing.B)   { benchmarkDiffClever(10, b) }
func BenchmarkDiffC100(b *testing.B)  { benchmarkDiffClever(100, b) }
func BenchmarkDiffC1000(b *testing.B) { benchmarkDiffClever(1000, b) }

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

func TestCopyFail(t *testing.T) {
	s := InitTest()
	s.dropbox = &dbTest{failCopy: true, paths: map[string][]string{"one": []string{"one", "two"}, "two": []string{"two"}}}

	s.runUpdate(context.Background(), &pb.SyncConfig{Origin: "one", Destination: "two"})
	if s.copies != 0 {
		t.Errorf("Too much copying")
	}
}

func TestListFail(t *testing.T) {
	s := InitTest()
	s.dropbox = &dbTest{failList: true, paths: map[string][]string{"one": []string{"one", "two"}, "two": []string{"two"}}}

	s.runUpdate(context.Background(), &pb.SyncConfig{Origin: "one", Destination: "two"})
	if s.copies != 0 {
		t.Errorf("Too much copying")
	}
}
