package main

import (
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"github.com/docker/docker/api/types/filters"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/api/types"
	"time"
	"fmt"
)

func getVolumes(ctx context.Context, client *client.Client) (volumetypes.VolumeListOKBody, error) {
	args := filters.Args{}
	volumes, err := client.VolumeList(ctx, args)
	if err != nil {
		return volumetypes.VolumeListOKBody{}, err
	}
	return volumes, nil
}

func removeVolumes(ctx context.Context, client *client.Client, volumes []*types.Volume) error {
	for _, v := range volumes {
		fmt.Printf("%s Removing Volume %s \n", time.Now().Format(time.RFC3339), v.Name)
		if err := client.VolumeRemove(ctx, v.Name, true); err != nil {
			return err
		}
		fmt.Printf("Removed Volume %s \n", v.Name)
	}
	return nil
}

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
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
			time.Now().Format(time.RFC3339),
			len(volumes.Volumes))

		if err := removeVolumes(ctx, cli, volumes.Volumes); err != nil {
			panic(err)
		}
		time.Sleep(5 * time.Minute)
	}
}