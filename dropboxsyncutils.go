package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"golang.org/x/net/context"

	pb "github.com/brotherlogic/dropboxsync/proto"
)

func stripFile(f string) string {
	elems := strings.Split(f, "/")
	return strings.ToLower(elems[len(elems)-1])
}

func (s *Server) diffFileListClever(master, new []string) []string {
	newOnes := []string{}

	mmap := make(map[string]string)
	for _, v := range master {
		mmap[stripFile(v)] = v
	}

	for _, m := range new {
		if _, ok := mmap[stripFile(m)]; !ok {
			s.Log(fmt.Sprintf("Adding %v because %v", m, mmap))
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
	source, err := s.dropbox.listFiles(config.Key, config.GetOrigin())
	dest, err2 := s.dropbox.listFiles(config.Key, config.GetDestination())
	s.listTime = time.Now().Sub(t)

	if err != nil || err2 != nil {
		s.Log(fmt.Sprintf("Error listing files %v and %v", err, err2))
		return
	}
	diffs := s.diffFileListClever(dest, source)

	for _, diff := range diffs {
		s.Log(fmt.Sprintf("Copying %v to %v", diff, config.GetDestination()+"/"+stripFile(diff)))
		err = s.dropbox.copyFile(config.Key, diff, config.GetDestination()+"/"+stripFile(diff))
		s.CtxLog(ctx, fmt.Sprintf("Copy error: %v", err))
		if err != nil {
			str1, ok := err.(files.CopyAPIError)
			str := "bad conversion"
			if ok {
				str = fmt.Sprintf("%+v", str1.EndpointError)
			}
			s.Log(fmt.Sprintf("Error copying files (%v), %v -> %v: %v, %v, %v",
				config.GetKey(),
				diff,
				config.GetDestination()+"/"+stripFile(diff),
				err,
				str,
				str1))
		} else {
			s.copies++
		}
	}
}
