package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "oasis"
	app.Usage = "a HTTP proxy building docker containers for each commits"
	app.Version = "0.0.1"
	app.Author = "Naoto Kaneko"
	app.Email = "naoty.k@gmail.com"
	app.Commands = []cli.Command{StartCommand, DebugCommand}
	app.Run(os.Args)
}

var StartCommand = cli.Command{
	Name:  "start",
	Usage: "start a HTTP proxy",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "proxy",
			Usage: "Proxy host",
		},
		cli.StringFlag{
			Name:  "backend",
			Usage: "Backend host",
		},
		cli.StringFlag{
			Name:  "repository",
			Usage: "Repository URL",
		},
	},
	Action: start,
}

func start(c *cli.Context) {
	proxyHost := c.String("proxy")
	backendHost := c.String("backend")
	repositoryURLString := c.String("repository")

	if proxyHost == "" || backendHost == "" || repositoryURLString == "" {
		cli.ShowCommandHelp(c, "start")
		os.Exit(1)
	}

	proxy := NewProxy(proxyHost, backendHost, repositoryURLString)
	log.Fatal(proxy.Start())
}

var DebugCommand = cli.Command{
	Name:  "debug",
	Usage: "debug",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "repo",
		},
		cli.StringFlag{
			Name: "revision",
		},
	},
	Action: debug,
}

func debug(c *cli.Context) {
	repo := c.String("repository")
	revision := c.String("revision")

	if repo == "" || revision == "" {
		cli.ShowCommandHelp(c, "debug")
	}
}
