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
	app.Commands = []cli.Command{StartCommand}
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
	},
	Action: start,
}

func start(c *cli.Context) {
	host := c.String("host")
	backend := c.String("backend")

	if host == "" || backend == "" {
		cli.ShowCommandHelp(c, "start")
	}

	proxy := Proxy{Host: host, BackendHost: backend}
	log.Fatal(proxy.Start())
}
