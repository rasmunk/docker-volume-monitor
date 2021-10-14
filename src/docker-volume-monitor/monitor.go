package main

import (
	"flag"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
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

func removeVolumes(ctx context.Context, client *client.Client, volumes []*types.Volume, debug bool) error {
	for _, v := range volumes {
		if !volumeInUse(ctx, client, v) {
			if debug {
				log.Infof("%s Removing Volume %s", time.Now().Format(time.RFC3339), v.Name)
			}
			if err := client.VolumeRemove(ctx, v.Name, true); err != nil {
				log.Errorf("%s Failed to remove Volume %s Reason %s", time.Now().Format(time.RFC3339), v.Name, err)
			} else {
				if debug {
					log.Infof("%s Removed Volume %s", time.Now().Format(time.RFC3339), v.Name)
				}
			}
		}
	}
	return nil
}

func run() {
	// Setup parameters
	var interval int
	var pruneUnused bool
	var debug bool

	flag.IntVar(&interval, "interval",
		10, "How often the monitor should check for unused volumes")
	flag.BoolVar(&pruneUnused, "prune-unused", true,
		"Whether the monitor should prune unused volumes")
	flag.BoolVar(&debug, "debug", false,
		"Set the debug flag to run the monitor in debug mode")
	flag.Parse()

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	for {
		log.Infof("%s Checking for volumes",
			time.Now().Format(time.RFC3339))
		volumes, err := getVolumes(ctx, cli)
		if err != nil {
			panic(err)
		}
		
		if debug {
			log.Infof("%s Found %d volumes",
			time.Now().Format(time.RFC3339), len(volumes.Volumes))
		}
		if pruneUnused {
			removeVolumes(ctx, cli, volumes.Volumes, debug)
		}
		log.Infof("%s Finished checking for volumes",
			time.Now().Format(time.RFC3339))
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}
