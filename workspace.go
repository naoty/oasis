package main

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
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
	workspace.buildImage()
	containerID, _ := workspace.runContainer()
	hostPort, _ := workspace.inspectHostPort(containerID)
	logrus.WithFields(logrus.Fields{"HostPort": hostPort}).Info("InspectHostPort")
	// workspace.updateIndex()
}

func (workspace *Workspace) clone() error {
	command := exec.Command("git", "clone", workspace.RepositoryURL.String(), workspace.Path)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()

	return err
}

func (workspace *Workspace) checkout(revision string) (string, error) {
	workspace.Revision = revision
	return workspace.exec("git", "checkout", revision)
}

func (workspace *Workspace) buildImage() (string, error) {
	return workspace.exec("docker", "build", "-t", workspace.imageName(), ".")
}

func (workspace *Workspace) runContainer() (string, error) {
	return workspace.exec("docker", "run", "-P", "-d", workspace.imageName())
}

func (workspace *Workspace) inspectHostPort(containerID string) (string, error) {
	result, err := workspace.exec("docker", "port", containerID)
	return workspace.parseHostPort(result), err
}

func (workspace *Workspace) exec(commands ...string) (string, error) {
	logrus.WithFields(logrus.Fields{
		"path": workspace.Path,
	}).Info(strings.Join(commands, " "))

	currentDir, _ := filepath.Abs(".")
	os.Chdir(workspace.Path)

	resultBuffer := bytes.NewBuffer(nil)

	command := exec.Command(commands[0], commands[1:]...)
	command.Stdout = resultBuffer
	command.Stderr = os.Stderr
	err := command.Run()

	os.Chdir(currentDir)

	result := resultBuffer.Bytes()
	return string(bytes.Trim(result, "\n")), err
}

func (workspace *Workspace) imageName() string {
	dir, projectName := path.Split(workspace.Path)
	username := path.Base(dir)
	return fmt.Sprintf("%s/%s:%s", username, projectName, workspace.Revision)
}

func (workspace *Workspace) parseHostPort(formattedPorts string) string {
	// formattedPorts is something like "3000/tcp -> 0.0.0.0:49155"
	ports := strings.Split(formattedPorts, " -> ")
	if len(ports) < 2 {
		return ""
	}

	elements := strings.Split(ports[1], ":")
	if len(elements) > 1 {
		return elements[1]
	} else {
		return ""
	}
}
