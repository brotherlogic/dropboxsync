package main

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/dropboxsync/proto"
)

func stripFile(f string) string {
	elems := strings.Split(f, "/")
	return elems[len(elems)-1]
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
	source, err := s.dropbox.listFiles(config.Key, config.Origin)
	dest, err2 := s.dropbox.listFiles(config.Key, config.Destination)
	s.listTime = time.Now().Sub(t) / 2

	if err != nil || err2 != nil {
		s.Log(fmt.Sprintf("Error listing files %v and %v", err, err2))
		return
	}

	diffs := diffFileList(dest, source)

	if len(diffs) > 0 {
		s.Log(fmt.Sprintf("Found %v diffs in %v", len(diffs), diffs))
		s.copies += int64(len(diffs))
	}
}
