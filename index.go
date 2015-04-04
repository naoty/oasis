package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

const (
	RootDir = ".oasis"
)

type Index struct {
	Root string
}

func LoadIndex() *Index {
	usr, _ := user.Current()
	root := filepath.Join(usr.HomeDir, RootDir)
	return &Index{Root: root}
}

func (i *Index) LookupPort(repo, revision string) (int, error) {
	absRoot, _ := filepath.Abs(i.Root)
	portFile := filepath.Join(absRoot, repo, revision)
	err := lookupFile(portFile)

	if err != nil {
		return 0, err
	}

	b, err := ioutil.ReadFile(portFile)
	trimmed := bytes.Trim(b, "\n")
	port, err := strconv.Atoi(string(trimmed))

	if err == nil {
		return port, nil
	} else {
		return 0, err
	}
}

func lookupFile(path string) error {
	_, err := os.Stat(path)
	return err
}
