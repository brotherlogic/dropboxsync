package main

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

func copyFile(key string, origin, dest string) {
	config := dropbox.Config{
		Token: key,
	}

	arg := files.NewRelocationArg(origin, dest)
	dbx := files.New(config)
	dbx.CopyV2(arg)
}

func listFiles(key string, path string) []string {
	config := dropbox.Config{
		Token: key,
	}

	arg := files.NewListFolderArg(path)

	dbx := files.New(config)
	resp, _ := dbx.ListFolder(arg)

	fs := []string{}
	for _, entry := range resp.Entries {
		if conv, ok := entry.(*files.FileMetadata); ok {
			fs = append(fs, conv.PathLower)
		}
	}

	return fs
}
