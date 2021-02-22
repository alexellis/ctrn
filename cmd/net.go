package cmd

import (
	"context"

	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
)

var _ oci.SpecOpts = WithVMNetwork

// WithVMNetwork modifies a container to use its host network settings (host netns, host utsns, host
// /etc/resolv.conf and host /etc/hosts). It's intended to configure Firecracker-containerd containers
// to have access to the network (if any) their VM was configured with.
func WithVMNetwork(ctx context.Context, cli oci.Client, ctr *containers.Container, spec *oci.Spec) error {
	for _, opt := range []oci.SpecOpts{
		oci.WithHostNamespace(specs.NetworkNamespace),
		oci.WithHostNamespace(specs.UTSNamespace),
		oci.WithHostResolvconf,
		oci.WithHostHostsFile,
	} {
		err := opt(ctx, cli, ctr, spec)
		if err != nil {
			return err
		}
	}

	return nil
}
