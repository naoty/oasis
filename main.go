package main

import (
	"fmt"
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
	app.Commands = []cli.Command{StartCommand, InspectCommand}
	app.Run(os.Args)
}

var StartCommand = cli.Command{
	Name:  "start",
	Usage: "start a HTTP proxy",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Usage: "host for the proxy to listen to",
		},
		cli.StringFlag{
			Name:  "backend",
			Usage: "backend host which the proxy redirects to",
		},
		cli.StringFlag{
			Name:  "repo",
			Usage: "repository URL",
		},
	},
	Action: start,
}

func start(c *cli.Context) {
	host := c.String("host")
	backend := c.String("backend")
	repo := c.String("repo")

	if host == "" || backend == "" || repo == "" {
		cli.ShowCommandHelp(c, "start")
	}

	proxy := Proxy{Host: host, BackendHost: backend, Repository: repo}
	log.Fatal(proxy.Start())
}

var InspectCommand = cli.Command{
	Name:  "inspect",
	Usage: "inspect for debug",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "repo",
		},
		cli.StringFlag{
			Name: "revision",
		},
	},
	Action: inspect,
}

func inspect(c *cli.Context) {
	repo := c.String("repo")
	revision := c.String("revision")

	if repo == "" || revision == "" {
		cli.ShowCommandHelp(c, "inspect")
	}

	index := LoadIndex()
	port, err := index.LookupPort(repo, revision)

	if err == nil {
		fmt.Printf("Port found: %s\n", port)
	} else {
		fmt.Printf("Error: %s\n", err)
	}
}
