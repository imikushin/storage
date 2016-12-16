package driver

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	rancherClient "github.com/rancher/go-rancher/client"

	"encoding/json"
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
	}, {
		Name:   "mountdest",
		Action: Mountdest,
	}, {
		Name:   "unmount",
		Action: Unmount,
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

func args(c *cli.Context) []string {
	a1 := c.Args().First()
	if a1 == "" {
		return []string{}
	}
	return append([]string{a1}, c.Args().Tail()...)
}

func logVolumeError(_ *model.Volume, err error) {
	logError(err)
}

func logError(err error) {
	if err != nil {
		println(err)
		os.Exit(1)
	}
}

func Create(c *cli.Context) {
	logVolumeError(storageDaemon(c).Create(volume(c)))
}

func Delete(c *cli.Context) {
	logError(storageDaemon(c).Delete(arg1(c), true))
}

func Attach(c *cli.Context) {
	logVolumeError(storageDaemon(c).Mount(arg1(c))) // TODO extract Attach from Mount
}

func Detach(c *cli.Context) {
	logError(storageDaemon(c).Unmount(arg1(c))) // TODO extract Detach from Unmount
}

func Mountdest(c *cli.Context) {
	logVolumeError(storageDaemon(c).Mount(args(c))) // TODO MountDest
}

func Unmount(c *cli.Context) {
	logError(storageDaemon(c).Unmount(arg1(c))) // TODO UnmountDest
}
