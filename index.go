package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/user"
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

func (i *Index) LookupPort(repo, revision string) (string, error) {
	absRoot, _ := filepath.Abs(i.Root)
	portFile := filepath.Join(absRoot, repo, revision)
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

func lookupFile(path string) error {
	_, err := os.Stat(path)
	return err
}
