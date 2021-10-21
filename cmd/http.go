package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewHttpCmd(timeoutDelay time.Duration) *cobra.Command {
	var (
		remoteAddr string
		remotePort uint16
	)

	var httpCmd = &cobra.Command{
		Use:   "http",
		Short: "Send events over to a http server",
		Long:  `Events are sent as single json object as a POST to http:\\{host}:{port}{path}`,
		Run: func(cmd *cobra.Command, args []string) {
			wg := sync.WaitGroup{}
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt)

			// create event generator
			gen := createEventGenerator()

			// connect to remote tcp server to emit event
			bo := backoff.NewExponentialBackOff()
			bo.MaxElapsedTime = 10 * time.Second
			u := url.URL{
				Scheme: "http",
				Host:   fmt.Sprintf("%s:%d", remoteAddr, remotePort),
				Path:   "/event",
			}

			// Process and send event on http destination until termination signal received.
			done := false
			go func() {
				wg.Add(1)
				for !done {
					select {
					case <-time.After(timeoutDelay):
						continue

					case ev := <-gen.Output():
						buf, err := json.Marshal(ev)
						cobra.CheckErr(err)

						err = backoff.Retry(func() error {
							if done {
								return backoff.Permanent(errors.New("program is terminating"))
							}

							resp, err := http.Post(u.String(), "application/json", bytes.NewReader(buf))
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
				wg.Done()
			}()

			// Start event generation in a goroutine
			go func() {
				wg.Add(1)
				err := runEventGenerator(gen)
				cobra.CheckErr(err)
				wg.Done()
			}()

			<-sigChan
			gen.Terminate()
			done = true

			// Wait for generator goroutine completion
			wg.Wait()
		},
	}

	httpCmd.Flags().StringVarP(&remoteAddr, "host", "H", "localhost", "remote host address")
	httpCmd.Flags().Uint16VarP(&remotePort, "port", "p", 3333, "remote port")

	return httpCmd
}
