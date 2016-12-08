package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/rancher/storage/longhorn-driver/driver"
	"github.com/rancher/storage/longhorn-driver/storagepool"
)

func main() {
	app := cli.NewApp()
	app.Name = "docker-longhorn-driver"
	app.Version = "0.1.0"
	app.Author = "Rancher Labs"
	app.Usage = "Docker volume plugin for Rancher Longhorn"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "enable debug logging level",
		},
		cli.StringFlag{
			Name:   "cattle-url",
			Usage:  "The URL endpoint to communicate with cattle server",
			EnvVar: "CATTLE_URL",
		},
		cli.StringFlag{
			Name:   "cattle-access-key",
			Usage:  "The access key required to authenticate with cattle server",
			EnvVar: "CATTLE_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "cattle-secret-key",
			Usage:  "The secret key required to authenticate with cattle server",
			EnvVar: "CATTLE_SECRET_KEY",
		},
		cli.IntFlag{
			Name:  "healthcheck-interval",
			Value: 5000,
			Usage: "set the frequency of performing healthchecks",
		},
		cli.StringFlag{
			Name:  "metadata-url",
			Usage: "set the metadata url",
			Value: "http://rancher-metadata/2015-12-19",
		},
	}

	app.Commands = []cli.Command{storagepool.Command, driver.Command}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatalf("Error running longhorn driver: %v", err)
	}

}
