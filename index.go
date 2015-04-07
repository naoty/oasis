package main

import (
	"bytes"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
	"path"

	"github.com/Sirupsen/logrus"
)

type Index struct {
	RootDir string
}

func NewIndex() *Index {
	currentUser, _ := user.Current()
	rootDir := path.Join(currentUser.HomeDir, ".oasis/index")
	return &Index{RootDir: rootDir}
}

func (index *Index) LookupPort(repositoryURL *url.URL, revision string) (string, error) {
	portFilePath := index.portFilePath(repositoryURL, revision)
	err := index.lookupPath(portFilePath)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadFile(portFilePath)
	if err == nil {
		return string(bytes.Trim(data, "\n")), nil
	} else {
		return "", err
	}
}

func (index *Index) UpdatePort(repositoryURL *url.URL, revision, port string) error {
	portFilePath := index.portFilePath(repositoryURL, revision)
	err := index.lookupPath(portFilePath)
	if err != nil {
		portFileDir, _ := path.Split(portFilePath)
		os.MkdirAll(portFileDir, 0755)
	}

	data := []byte(port)
	err = ioutil.WriteFile(portFilePath, data, 0644)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to write a port")
	}

	return err
}

func (index *Index) portFilePath(repositoryURL *url.URL, revision string) string {
	return path.Join(index.RootDir, repositoryURL.Host, repositoryURL.Path, revision)
}

func (index *Index) lookupPath(path string) error {
	_, err := os.Stat(path)
	return err
}
