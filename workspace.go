package main

import (
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

type Workspace struct {
	Path          string
	Revision      string
	RepositoryURL *url.URL
	DockerClient  *docker.Client
	Index         *Index
}

func NewWorkspace(repositoryURL, containerHostURL *url.URL, index *Index) *Workspace {
	currentUser, _ := user.Current()
	localPath := path.Join(currentUser.HomeDir, ".oasis/repositories", repositoryURL.Host, repositoryURL.Path)
	dockerClient := NewDockerClient(containerHostURL)
	return &Workspace{
		Path:          localPath,
		Revision:      "",
		RepositoryURL: repositoryURL,
		DockerClient:  dockerClient,
		Index:         index,
	}
}

func NewDockerClient(containerHostURL *url.URL) *docker.Client {
	dockerCertPath := os.Getenv("DOCKER_CERT_PATH")
	ca := path.Join(dockerCertPath, "ca.pem")
	cert := path.Join(dockerCertPath, "cert.pem")
	key := path.Join(dockerCertPath, "key.pem")
	client, err := docker.NewTLSClient(containerHostURL.String(), cert, key, ca)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"path": dockerCertPath,
			"ca":   ca,
			"cert": cert,
			"key":  key,
		}).Fatal("Failed to initialize a docker client")
	}

	return client
}

func (workspace *Workspace) Setup(revision string) {
	workspace.clone()
	workspace.checkout(revision)
	workspace.buildContainer()
	// workspace.runContainer()
	// workspace.updateIndex()
}

func (workspace *Workspace) clone() error {
	command := exec.Command("git", "clone", workspace.RepositoryURL.String(), workspace.Path)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()

	return err
}

func (workspace *Workspace) checkout(revision string) error {
	logrus.WithFields(logrus.Fields{
		"command": strings.Join([]string{"git", "checkout", revision}, " "),
		"path":    workspace.Path,
	}).Info("Command exec")

	workspace.Revision = revision

	currentDir, _ := filepath.Abs(".")
	os.Chdir(workspace.Path)

	command := exec.Command("git", "checkout", revision)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()

	os.Chdir(currentDir)

	return err
}

func (workspace *Workspace) buildContainer() error {
	options := docker.BuildImageOptions{
		OutputStream: os.Stdout,
		ContextDir:   workspace.Path,
	}
	err := workspace.DockerClient.BuildImage(options)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to build an image")
	}

	return err
}
