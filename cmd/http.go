package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"

	"github.com/cenkalti/backoff/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewHttpCmd() *cobra.Command {
	var httpCmd = &cobra.Command{
		Use:   "http",
		Short: "Send events over to a http server",
		Long:  `Events are sent as single json object as a POST to http:\\{host}:{port}{path}`,
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

			host := "localhost"
			port := 3333

			// connect to remote tcp server to emit event
			bo := backoff.NewExponentialBackOff()
			url := url.URL{
				Scheme: "http",
				Host:   fmt.Sprintf("%s:%d", host, port),
				Path:   "/event",
			}

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

					log.Debugln(string(buf))
					err = backoff.Retry(func() error {
						resp, err := http.Post(url.String(), "application/json", bytes.NewReader(buf))
						if err != nil {
							return err
						}

						if resp.StatusCode < 200 && resp.StatusCode > 208 {
							return fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
						}

						_ = resp.Body.Close()
						return nil
					}, bo)
					if err != nil {
						log.Errorf("failed to write event to http: %s", err)
					}
				}
			}

			// Wait for generator goroutine completion
			wg.Wait()
		},
	}

	return httpCmd
}
