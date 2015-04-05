package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

const (
	IndexRootDir = ".oasis/index"
)

type Index struct {
	Root string
}

func LoadIndex() *Index {
	usr, _ := user.Current()
	root := filepath.Join(usr.HomeDir, IndexRootDir)
	return &Index{Root: root}
}

func (i *Index) LookupPort(repositoryURLString, revision string) (string, error) {
	absRoot, _ := filepath.Abs(i.Root)
	repositoryPath := parseRepositoryPath(repositoryURLString)
	portFile := filepath.Join(absRoot, repositoryPath, revision)
	err := lookupFile(portFile)

	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadFile(portFile)

	if err == nil {
		return string(bytes.Trim(b, "\n")), nil
	} else {
		return "", err
	}
}

func parseRepositoryPath(repositoryURLString string) string {
	repositoryURL, err := url.Parse(repositoryURLString)
	if err != nil {
		log.Fatalf("URL parse error: %s", err)
	}

	return path.Join(repositoryURL.Host, repositoryURL.Path)
}

func lookupFile(path string) error {
	_, err := os.Stat(path)
	return err
}
