package cmd

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/errdefs"
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	RunE:    removeRunner,
	Aliases: []string{"delete", "rm"},
}

func removeRunner(cmd *cobra.Command, args []string) error {
	name := "helloweb"
	// create a container
	container, err := client.LoadContainer(rootCtx,
		name,
	)

	if err != nil {
		return err
	}

	t, err := container.Task(rootCtx, nil)

	if err != nil && !errdefs.IsNotFound(err) {
		return err
	}

	if t != nil {

		status, err := t.Status(rootCtx)
		if err != nil {
			return fmt.Errorf("Unable to get status for: %s, error: %s", name, err.Error())
		}

		log.Printf("Status: %s", status.Status)

		killTask(rootCtx, t)
	}

	log.Println("Delete container")
	if err = container.Delete(rootCtx, containerd.WithSnapshotCleanup); err != nil {
		return err
	}

	return nil
}

func killTask(ctx context.Context, task containerd.Task) error {

	killTimeout := 30 * time.Second

	wg := &sync.WaitGroup{}
	wg.Add(1)
	var err error

	go func() {
		defer wg.Done()
		if task != nil {
			wait, err := task.Wait(ctx)
			if err != nil {
				err = fmt.Errorf("error waiting on task: %s", err)
				return
			}
			if err := task.Kill(ctx, unix.SIGTERM, containerd.WithKillAll); err != nil {
				log.Printf("error killing container task: %s", err)
			}

			select {
			case <-wait:
				task.Delete(ctx)
				return
			case <-time.After(killTimeout):
				if err := task.Kill(ctx, unix.SIGKILL, containerd.WithKillAll); err != nil {
					log.Printf("error force killing container task: %s", err)
				}
				return
			}
		}
	}()
	wg.Wait()

	return err
}
