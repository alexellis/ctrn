package cmd

import (
	"context"
	"log"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/spf13/cobra"
)

var (
	client  *containerd.Client
	rootCtx context.Context
)

var rootCmd = &cobra.Command{Use: "ctrn"}

func init() {
	c, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		log.Fatalf("fail to connect to containerd: %v", err)
	}

	client = c

	rootCtx = namespaces.WithNamespace(context.Background(), "default")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(removeCmd)

	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(specCmd)
	rootCmd.AddCommand(netCmd)
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
}

// Execute is the entry point of the application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
