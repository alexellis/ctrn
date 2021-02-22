package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/containerd/containerd/errdefs"
	"github.com/spf13/cobra"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/oci"
)

var createCmd = &cobra.Command{
	Use:  "create",
	RunE: createRunner,
}

func createRunner(cmd *cobra.Command, args []string) error {
	ref := "docker.io/functions/figlet:latest"
	name := "helloweb"

	snapshotter := defaultSnapshotter
	if val, ok := os.LookupEnv("SNAPSHOTTER"); ok && len(val) > 0 {
		snapshotter = val
	}

	image, err := prepareImage(rootCtx, client, ref, snapshotter)

	if err != nil {
		return err
	}

	runtime := ""
	if val, ok := os.LookupEnv("RUNTIME"); ok && len(val) > 0 {
		runtime = val
	}

	fmt.Printf("Runtime: %s\tSnapshotter: %s\n", runtime, snapshotter)

	// create a container
	container, err := client.NewContainer(
		rootCtx,
		name,
		containerd.WithImage(image),
		containerd.WithSnapshotter(snapshotter),
		containerd.WithRuntime(runtime, nil),
		containerd.WithNewSnapshot(name+"-snapshot", image),
		containerd.WithNewSpec(
			oci.WithImageConfig(image),
			oci.WithCapabilities([]string{"CAP_NET_RAW"}),
			WithVMNetwork,
		),
	)

	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}

	fmt.Println(container)
	return nil
}

func prepareImage(ctx context.Context, client *containerd.Client, imageName, snapshotter string) (containerd.Image, error) {

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
