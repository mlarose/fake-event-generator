package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"sync"
)

func NewStdoutCmd() *cobra.Command {
	var stdoutCmd = &cobra.Command{
		Use:   "stdout",
		Short: "Write events to stdout",
		Long:  `Events are sent line separated json objects`,
		Run: func(cmd *cobra.Command, args []string) {
			wg := sync.WaitGroup{}
			sigChan := make(chan os.Signal)
			signal.Notify(sigChan, os.Interrupt)

			// Start event generation in a goroutine
			gen := createEventGenerator()
			go func() {
				wg.Add(1)
				err := runEventGenerator(gen)
				cobra.CheckErr(err)

				wg.Done()
			}()

			// Process and send event on http destination until termination signal received.
			done := false
			for !done {
				select {
				case <-sigChan:
					gen.Terminate()
					done = true
				case ev := <-gen.Output():
					buf, err := json.Marshal(ev)
					cobra.CheckErr(err)

					fmt.Println(string(buf))
				}
			}

			// Wait for generator goroutine completion
			wg.Wait()
		},
	}

	return stdoutCmd
}
