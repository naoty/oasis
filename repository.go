package main

import (
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
)

type Repository struct {
	Path      string
	RemoteURL *url.URL
	Revision  string
}

const RepositoryRootDir = ".oasis/repositories"

func NewRepository(remoteURLString, branch string) *Repository {
	remoteURL, err := url.Parse(normalizeRepositoryURLString(remoteURLString))
	if err != nil {
		log.Fatalf("URL parse error: %s", err)
	}

	localPath := parseLocalPath(remoteURL.String())
	if !isExist(localPath) {
		clone(remoteURL.String(), localPath, branch)
	}
	return &Repository{Path: localPath, RemoteURL: remoteURL, Revision: branch}
}

func (repository *Repository) Checkout(revision string) {
	repository.Revision = revision
	err := repository.Exec("git", "checkout", revision)
	if err != nil {
		log.Fatalf("git checkout %s: %s", revision, err)
	}
}

func (repository *Repository) Exec(commands ...string) (err error) {
	currentDir, _ := filepath.Abs(".")
	workspaceDir, _ := filepath.Abs(repository.Path)
	os.Chdir(workspaceDir)

	command := exec.Command(commands[0], commands[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Run()

	os.Chdir(currentDir)

	return err
}

func clone(remoteURLString, localPath, branch string) {
	command := exec.Command("git", "clone", "--branch", branch, remoteURLString, localPath)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		log.Fatalf("git clone %s %s: %s", remoteURLString, localPath, err)
	}
}

var urlPattern = regexp.MustCompile("^[^:]+://")

func normalizeRepositoryURLString(urlString string) string {
	if urlPattern.MatchString(urlString) {
		return urlString
	} else {
		return "https://" + urlString
	}
}

func parseLocalPath(remoteURLString string) string {
	remoteURL, err := url.Parse(remoteURLString)
	if err != nil {
		log.Fatalf("URL Parse error: %s", err)
	}

	currentUser, _ := user.Current()
	return path.Join(currentUser.HomeDir, RepositoryRootDir, remoteURL.Host, remoteURL.Path)
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
