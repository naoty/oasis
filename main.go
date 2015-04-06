package main

import (
	"net/url"
	"os"
	"regexp"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "oasis"
	app.Usage = "a HTTP proxy building docker containers for each commits"
	app.Version = "0.0.1"
	app.Author = "Naoto Kaneko"
	app.Email = "naoty.k@gmail.com"
	app.Commands = []cli.Command{StartCommand}
	app.Run(os.Args)
}

var StartCommand = cli.Command{
	Name:  "start",
	Usage: "start a HTTP proxy",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "proxy, p",
			Usage: "Proxy URL",
		},
		cli.StringFlag{
			Name:  "container-host, c",
			Usage: "Container Host URL",
		},
		cli.StringFlag{
			Name:  "repository, r",
			Usage: "Repository URL",
		},
	},
	Action: start,
}

func start(context *cli.Context) {
	proxyURLString := context.String("proxy")
	containerHostURLString := context.String("container-host")
	repositoryURLString := context.String("repository")

	if proxyURLString == "" || containerHostURLString == "" || repositoryURLString == "" {
		cli.ShowCommandHelp(context, "start")
		os.Exit(1)
	}

	proxyURLString = normalizeURLString(proxyURLString)
	containerHostURLString = normalizeURLString(containerHostURLString)
	repositoryURLString = normalizeRepositoryURLString(repositoryURLString)

	logFields := logrus.Fields{
		"proxy":          proxyURLString,
		"container-host": containerHostURLString,
		"repository":     repositoryURLString,
	}

	logrus.WithFields(logFields).Info("Start")

	proxyURL, err := url.Parse(proxyURLString)
	containerHostURL, err := url.Parse(containerHostURLString)
	repositoryURL, err := url.Parse(repositoryURLString)

	if err != nil {
		logrus.WithFields(logFields).Fatal("Failed to parse URL")
	}

	proxy := NewProxy(proxyURL, containerHostURL, repositoryURL)
	err = proxy.Start()

	logrus.WithFields(logrus.Fields{"error": err}).Fatal("Stop a proxy")
}

var urlPattern = regexp.MustCompile("^[^:]+://")

func normalizeURLString(urlString string) string {
	if urlPattern.MatchString(urlString) {
		return urlString
	} else {
		return "http://" + urlString
	}
}

func normalizeRepositoryURLString(urlString string) string {
	if urlPattern.MatchString(urlString) {
		return urlString
	} else {
		return "https://" + urlString
	}
}
