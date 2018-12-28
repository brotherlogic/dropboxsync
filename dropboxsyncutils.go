package main

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/dropboxsync/proto"
	pbt "github.com/brotherlogic/tracer/proto"
)

func stripFile(f string) string {
	elems := strings.Split(f, "/")
	return elems[len(elems)-1]
}

func diffFileListClever(master, new []string) []string {
	newOnes := []string{}

	mmap := make(map[string]string)
	for _, v := range master {
		mmap[stripFile(v)] = v
	}

	for _, m := range new {
		if _, ok := mmap[stripFile(m)]; !ok {
			newOnes = append(newOnes, m)
		}
	}

	return newOnes
}

func diffFileListPreStrip(master, new []string) []string {
	newOnes := []string{}

	mmap := make(map[string]string)
	for _, v := range master {
		mmap[v] = stripFile(v)
	}
	for _, v := range new {
		mmap[v] = stripFile(v)
	}

	for _, f := range new {
		found := false
		for _, m := range master {
			if mmap[m] == mmap[f] {
				found = true
			}
		}

		if !found {
			newOnes = append(newOnes, f)
		}
	}

	return newOnes
}

func diffFileList(master, new []string) []string {
	newOnes := []string{}

	for _, f := range new {
		found := false
		for _, m := range master {
			if stripFile(m) == stripFile(f) {
				found = true
			}
		}

		if !found {
			newOnes = append(newOnes, f)
		}
	}

	return newOnes
}

func (s *Server) runUpdate(ctx context.Context, config *pb.SyncConfig) {
	t := time.Now()
	s.LogTrace(ctx, "prelist-1", time.Now(), pbt.Milestone_START_EXTERNAL)
	source, err := s.dropbox.listFiles(config.Key, config.Origin)
	s.LogTrace(ctx, "postlist-1", time.Now(), pbt.Milestone_END_EXTERNAL)
	s.LogTrace(ctx, "prelist-2", time.Now(), pbt.Milestone_START_EXTERNAL)
	dest, err2 := s.dropbox.listFiles(config.Key, config.Destination)
	s.LogTrace(ctx, "postlist-2", time.Now(), pbt.Milestone_END_EXTERNAL)
	s.listTime = time.Now().Sub(t)

	if err != nil || err2 != nil {
		s.Log(fmt.Sprintf("Error listing files %v and %v", err, err2))
		return
	}
	diffs := diffFileListClever(dest, source)
	s.LogTrace(ctx, fmt.Sprintf("diffed-%v", len(dest)), time.Now(), pbt.Milestone_MARKER)

	for _, diff := range diffs {
		s.Log(fmt.Sprintf("Copying %v to %v", diff, config.Destination+"/"+stripFile(diff)))
		s.LogTrace(ctx, "precopy", time.Now(), pbt.Milestone_START_EXTERNAL)
		err = s.dropbox.copyFile(config.Key, diff, config.Destination+"/"+stripFile(diff))
		s.LogTrace(ctx, "postcopy", time.Now(), pbt.Milestone_END_EXTERNAL)
		if err != nil {
			s.Log(fmt.Sprintf("Error copying files: %v", err))
		} else {
			s.copies++
		}
	}
}
