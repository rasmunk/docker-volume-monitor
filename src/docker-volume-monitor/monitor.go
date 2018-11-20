package main


import (
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"github.com/docker/docker/api/types/filters"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/api/types"
	"time"
	"fmt"
	"flag"
)

func getVolumes(ctx context.Context, client *client.Client) (volumetypes.VolumesListOKBody, error) {
	args := filters.Args{}
	volumes, err := client.VolumeList(ctx, args)
	if err != nil {
		return volumetypes.VolumesListOKBody{}, err
	}
	return volumes, nil
}

func removeVolumes(ctx context.Context, client *client.Client, volumes []*types.Volume) error {
	for _, v := range volumes {
		fmt.Printf("%s Removing Volume %s \n",
			time.Now().Format(time.RFC3339), v.Name)
		if err := client.VolumeRemove(ctx, v.Name, true); err != nil {
			fmt.Printf("Failed to remove Volume %s Reason %s \n",
				v.Name, err)
		} else {
			fmt.Printf("Removed Volume %s \n", v.Name)
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
		fmt.Printf("%s Checking for volumes \n",
			time.Now().Format(time.RFC3339))
		volumes, err := getVolumes(ctx, cli)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s Found %d volumes \n",
			time.Now().Format(time.RFC3339), len(volumes.Volumes))

		if pruneUnused {
			removeVolumes(ctx, cli, volumes.Volumes)
		}
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}
