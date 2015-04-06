package main

import (
	"bytes"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
	"path"
)

type Index struct {
	RootDir       string
	RepositoryURL *url.URL
}

func NewIndex(repositoryURL *url.URL) *Index {
	currentUser, _ := user.Current()
	rootDir := path.Join(currentUser.HomeDir, ".oasis/index")
	return &Index{RootDir: rootDir, RepositoryURL: repositoryURL}
}

func (index *Index) LookupPort(revision string) (string, error) {
	portFilePath := path.Join(index.RootDir, index.RepositoryURL.Host, index.RepositoryURL.Path, revision)
	err := lookupPath(portFilePath)
	if err != nil {
		return "", err
	}

	binary, err := ioutil.ReadFile(portFilePath)
	if err == nil {
		return string(bytes.Trim(binary, "\n")), nil
	} else {
		return "", err
	}
}

func lookupPath(path string) error {
	_, err := os.Stat(path)
	return err
}
