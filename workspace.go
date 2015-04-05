package main

import (
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
)

const (
	ReposRootDir = ".oasis/repos"
)

type Workspace struct {
	Path string
}

func NewWorkspace(repositoryURLString, revision string) *Workspace {
	var workspace *Workspace

	local := parseLocalPath(repositoryURLString)
	if isExist(local) {
		workspace = &Workspace{Path: local}
	} else {
		workspace = gitClone(repositoryURLString)
	}
	workspace.gitCheckout(revision)
	workspace.gitPull()

	// TODO: docker build

	return workspace
}

func (w *Workspace) Run() {
	// TODO: docker run
}

func (w *Workspace) InspectPort() string {
	// TODO: inspect the port of the container
	return ""
}

func (w *Workspace) Exec(commands ...string) error {
	return w.ExecFunc(func() error {
		command := exec.Command(commands[0], commands[1:]...)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		return command.Run()
	})
}

func (w *Workspace) ExecFunc(f func() error) error {
	currentDir, _ := filepath.Abs(".")

	workspaceDir, _ := filepath.Abs(w.Path)
	os.Chdir(workspaceDir)

	err := f()

	os.Chdir(currentDir)

	return err
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func gitClone(repositoryURLString string) *Workspace {
	local := parseLocalPath(repositoryURLString)

	command := exec.Command("git", "clone", repositoryURLString, local)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()

	if err != nil {
		log.Fatalf("git clone %s %s", repositoryURLString, local, err)
	}

	return &Workspace{Path: local}
}

func (w *Workspace) gitCheckout(revision string) error {
	return w.Exec("git", "checkout", revision)
}

func (w *Workspace) gitPull() error {
	return w.Exec("git", "pull")
}

func parseLocalPath(rawURLString string) string {
	rawURL, err := url.Parse(rawURLString)

	if err != nil {
		log.Fatalf("While parsing repository URL: %s", err)
	}

	usr, _ := user.Current()
	return path.Join(usr.HomeDir, ReposRootDir, rawURL.Host, rawURL.Path)
}
