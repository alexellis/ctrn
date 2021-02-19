package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/google/uuid"

	"github.com/containerd/containerd/cio"
	gocni "github.com/containerd/go-cni"
)

var startCmd = &cobra.Command{
	Use:  "start",
	RunE: startRunner,
}

func startRunner(cmd *cobra.Command, args []string) error {
	name := "helloweb"

	container, err := client.LoadContainer(rootCtx, name)
	if err != nil {
		return fmt.Errorf("failed to load container %s: %v", name, err)
	}

	fmt.Printf("Container %v\n", container)

	task, err := container.NewTask(rootCtx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return fmt.Errorf("Fail to create task: %v", err)
	}
	fmt.Printf("Task %v\n", task)
	id := uuid.New().String()
	netns := getNetns(task.Pid())

	cni, err := gocni.New(
		gocni.WithPluginConfDir("./net.d/"),
		gocni.WithPluginDir([]string{"/opt/cni/bin/"}),
	)

	if err != nil {
		return err
	}

	// Load the cni configuration
	if err := cni.Load(gocni.WithLoNetwork, gocni.WithDefaultConf); err != nil {
		return fmt.Errorf("failed to load cni configuration: %v", err)
	}

	labels := map[string]string{
		// "OPENFAAS": "yes",
	}

	result, err := cni.Setup(rootCtx, id, netns, gocni.WithLabels(labels))
	if err != nil {
		return fmt.Errorf("failed to setup network for namespace %q: %v", id, err)
	}

	for name, config := range result.Interfaces {
		fmt.Printf("Config of interface %s: %v\n", name, config)
	}

	if err := task.Start(rootCtx); err != nil {
		return fmt.Errorf("failed to start task: %v", err)
	}
	return nil
}

func getNetns(pid uint32) string {
	return fmt.Sprintf("/proc/%d/ns/net", pid)
}
