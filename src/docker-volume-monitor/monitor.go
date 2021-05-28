package main

import (
	"flag"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func getVolumes(ctx context.Context, client *client.Client) (volumetypes.VolumesListOKBody, error) {
	args := filters.Args{}
	volumes, err := client.VolumeList(ctx, args)
	if err != nil {
		return volumetypes.VolumesListOKBody{}, err
	}
	return volumes, nil
}

func volumeInUse(ctx context.Context, client *client.Client, volume *types.Volume) bool {
	containers, err := client.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return false
	}
	inUse := false

	for _, container := range containers {
		for _, mount := range container.Mounts {
			if mount.Name == volume.Name {
				inUse = true
			}
		}
	}

	return inUse
}

func removeVolumes(ctx context.Context, client *client.Client, volumes []*types.Volume) error {
	for _, v := range volumes {
		if !volumeInUse(ctx, client, v) {
			log.Infof("%s Removing Volume %s \n", time.Now().Format(time.RFC3339), v.Name)
			if err := client.VolumeRemove(ctx, v.Name, true); err != nil {
				log.Errorf("Failed to remove Volume %s Reason %s", v.Name, err)
			} else {
				log.Infof("Removed Volume %s", v.Name)
			}
		} else {
			log.Infof("")
		}
	}
	return nil
}

func run() {
	// Setup parameters
	var interval int
	var pruneUnused bool

	flag.IntVar(&interval, "interval",
		10, "How often the monitor should check for unused volumes")
	flag.BoolVar(&pruneUnused, "prune-unused", true,
		"Whether the monitor should prune unused volumes")
	flag.Parse()

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	for {
		log.Infof("%s Checking for volumes \n",
			time.Now().Format(time.RFC3339))
		volumes, err := getVolumes(ctx, cli)
		if err != nil {
			panic(err)
		}

		log.Infof("%s Found %d volumes \n",
			time.Now().Format(time.RFC3339), len(volumes.Volumes))
		if pruneUnused {
			removeVolumes(ctx, cli, volumes.Volumes)
		}
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}
