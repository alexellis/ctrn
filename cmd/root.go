package cmd

import (
	"context"
	"log"
	"os"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/spf13/cobra"
)

var (
	client  *containerd.Client
	rootCtx context.Context
)

var rootCmd = &cobra.Command{Use: "ctrn"}

const defaultSnapshotter = "overlayfs"

func init() {
	sock := "/run/containerd/containerd.sock"
	if val, ok := os.LookupEnv("CONTAINERD"); ok && len(val) > 0 {
		sock = val
	}

	c, err := containerd.New(sock)
	if err != nil {
		log.Fatalf("fail to connect to containerd: %v", err)
	}

	client = c

	rootCtx = namespaces.WithNamespace(context.Background(), "default")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(removeCmd)

	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(specCmd)
	rootCmd.AddCommand(startCmd)

	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
}

// Execute is the entry point of the application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
