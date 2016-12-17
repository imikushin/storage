package driver

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	rancherClient "github.com/rancher/go-rancher/client"
	"github.com/rancher/storage/longhorn-driver/model"
	"github.com/rancher/storage/longhorn-driver/util"
	"os"
)

const (
	NAME = "name"
	SIZE = "size"
)

var Command = cli.Command{
	Name: "vol",
	Subcommands: []cli.Command{{
		Name:   "create",
		Action: Create,
		Flags:  []cli.Flag{cli.StringFlag{Name: NAME}, cli.IntFlag{Name: SIZE}},
	}, {
		Name:   "delete",
		Action: Delete,
	}, {
		Name:   "attach",
		Action: Attach,
	}, {
		Name:   "detach",
		Action: Detach,
	}},
}

func storageDaemon(c *cli.Context) *StorageDaemon {
	client, err := rancherClient.NewRancherClient(&rancherClient.ClientOpts{
		Url:       c.GlobalString("cattle-url"),
		AccessKey: c.GlobalString("cattle-access-key"),
		SecretKey: c.GlobalString("cattle-secret-key"),
	})
	if err != nil {
		logrus.Fatalf("Failed to establish connection to Rancher server")
	}

	md, err := util.GetMetadataConfig(c.GlobalString("metadata-url"))
	if err != nil {
		logrus.Fatalf("Unable to get metadata: %v", err)
	}

	sd, err := NewStorageDaemon(md.ContainerName, md.DriverName, md.Image, client)
	if err != nil {
		logrus.Fatalf("Error creating storage daemon: %v", err)
	}

	return sd
}

func volume(c *cli.Context) *model.Volume {
	v := &model.Volume{}
	if err := json.Unmarshal([]byte(c.Args().First()), v); err != nil {
		logrus.Fatalf("Error unmarshaling model.Volume: %s", err)
	}
	return v
}

func arg1(c *cli.Context) string {
	return c.Args().First()
}

func printVolumeErr(_ *model.Volume, err error) {
	printErr(err)
}

func printResultErr(result string, err error) {
	if _, err := fmt.Fprintln(os.Stdout, result); err != nil {
		logrus.Fatalf("Could not even print to stdout: %s", err)
	}
	printErr(err)
}

func printErr(err error) {
	if err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			logrus.Fatalf("Could not even print to stderr: %s", err)
		}
		os.Exit(1)
	}
}

func Create(c *cli.Context) {
	printVolumeErr(storageDaemon(c).Create(volume(c)))
}

func Delete(c *cli.Context) {
	printErr(storageDaemon(c).Delete(arg1(c), true))
}

func Attach(c *cli.Context) {
	printResultErr(storageDaemon(c).Attach(arg1(c)))
}

func Detach(c *cli.Context) {
	// TODO looks like it is not necessary
}
