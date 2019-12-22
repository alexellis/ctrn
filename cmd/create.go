package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/containerd/containerd/errdefs"

	"github.com/spf13/cobra"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/oci"
)

var createCmd = &cobra.Command{
	Use: "create",
	Run: createRunner,
}

func createRunner(cmd *cobra.Command, args []string) {
	ref := "docker.io/functions/figlet:latest"

	image, err := prepareImage(rootCtx, client, ref)

	if err != nil {
		log.Fatal(err)
	}

	// create a container
	container, errC := client.NewContainer(
		rootCtx,
		"helloweb",
		containerd.WithImage(image),
		containerd.WithNewSnapshot("hello-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image),
			oci.WithCapabilities([]string{"CAP_NET_RAW"})),
	)

	if errC != nil {
		log.Fatalf("Fail to create container: %v\n", errC)
	}
	fmt.Println(container)
}

const defaultSnapshotter = "overlayfs"

func prepareImage(ctx context.Context, client *containerd.Client, imageName string) (containerd.Image, error) {
	snapshotter := defaultSnapshotter

	var empty containerd.Image
	image, err := client.GetImage(ctx, imageName)
	if err != nil {
		if !errdefs.IsNotFound(err) {
			return empty, err
		}

		img, err := client.Pull(ctx, imageName, containerd.WithPullUnpack)
		if err != nil {
			return empty, fmt.Errorf("cannot pull: %s", err)
		}
		image = img
	}

	unpacked, err := image.IsUnpacked(ctx, snapshotter)
	if err != nil {
		return empty, fmt.Errorf("cannot check if unpacked: %s", err)
	}

	if !unpacked {
		if err := image.Unpack(ctx, snapshotter); err != nil {
			return empty, fmt.Errorf("cannot unpack: %s", err)
		}
	}

	return image, nil
}
