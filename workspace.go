package main

import (
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Workspace struct {
	Path          string
	RepositoryURL *url.URL
}

func NewWorkspace(repositoryURL *url.URL) *Workspace {
	currentUser, _ := user.Current()
	localPath := path.Join(currentUser.HomeDir, ".oasis/repositories", repositoryURL.Host, repositoryURL.Path)
	return &Workspace{Path: localPath, RepositoryURL: repositoryURL}
}

func (workspace *Workspace) Setup(revision string) {
	workspace.clone()
	workspace.checkout(revision)
}

func (workspace *Workspace) clone() error {
	command := exec.Command("git", "clone", workspace.RepositoryURL.String(), workspace.Path)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()

	return err
}

func (workspace *Workspace) checkout(revision string) error {
	log.WithFields(log.Fields{
		"command": strings.Join(commands, " "),
		"path":    workspace.Path,
	}).Info("Command exec")

	currentDir, _ := filepath.Abs(".")
	os.Chdir(workspace.Path)

	command := exec.Command("git", "checkout", revision)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()

	os.Chdir(currentDir)

	return err
}
