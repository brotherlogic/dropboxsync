package main

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type dbProd struct {
	log func(s string)
}

func (d *dbProd) copyFile(key string, origin, dest string) error {
	config := dropbox.Config{
		Token: key,
	}

	arg := files.NewRelocationArg(origin, dest)
	dbx := files.New(config)
	_, err := dbx.CopyV2(arg)
	return err
}

func (d *dbProd) listFiles(key string, path string) ([]string, error) {
	config := dropbox.Config{
		Token: key,
	}

	arg := files.NewListFolderArg(path)

	dbx := files.New(config)
	resp, err := dbx.ListFolder(arg)

	if err != nil {
		return []string{}, err
	}

	fs := []string{}
	for _, entry := range resp.Entries {
		if conv, ok := entry.(*files.FileMetadata); ok {
			fs = append(fs, conv.PathLower)
		}
	}
	for resp.HasMore {
		resp, err := dbx.ListFolderContinue(&files.ListFolderContinueArg{Cursor: resp.Cursor})
		if err != nil {
			return []string{}, err
		}
		for _, entry := range resp.Entries {
			if conv, ok := entry.(*files.FileMetadata); ok {
				fs = append(fs, conv.PathLower)
			}
		}
	}

	return fs, nil
}
