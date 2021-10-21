package cmd

import (
	"encoding/json"
	"errors"
	"fake-event-generator/output"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewTcpCmd(timeoutDelay time.Duration) *cobra.Command {
	var (
		remoteAddr string
		remotePort uint16
	)

	tcpCmd := &cobra.Command{
		Use:   "tcp",
		Short: "Send events over a unencrypted tcp connection to a destination address",
		Long: `The program will connect to a destination host:port service that accept unencrypted tcp connection.
The events are sent in a new line delimited stream of json documents, each one representing a single event.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			wg := sync.WaitGroup{}
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt)
			done := false

			// Start event generation in a goroutine
			gen := createEventGenerator()
			go func() {
				wg.Add(1)
				err := runEventGenerator(gen)
				cobra.CheckErr(err)
				wg.Done()
			}()

			go func() {
				wg.Add(1)

				// connect to remote tcp server to emit event
				var tcp io.WriteCloser
				remoteHost := fmt.Sprintf("%s:%d", remoteAddr, remotePort)
				bo := backoff.NewExponentialBackOff()
				bo.InitialInterval = 100 * time.Millisecond
				bo.MaxInterval = 5 * time.Second
				bo.MaxElapsedTime = 15 * time.Second

				err = backoff.Retry(func() error {
					if done {
						return backoff.Permanent(errors.New("program is terminating"))
					}
					addr, err := net.ResolveTCPAddr("tcp", remoteHost)
					if err != nil {
						return err
					}

					tcp, err = output.NewTcpWriteCloser(addr, bo)
					return err
				}, bo)
				cobra.CheckErr(err)

				// Process and send event on the tcp output until termination signal received.
				for !done {
					select {
					case <-time.After(timeoutDelay):
						continue
					case ev := <-gen.Output():
						buf, err := json.Marshal(ev)
						cobra.CheckErr(err)

						log.Debugln(string(buf))
						n, err := tcp.Write(append(buf, '\n'))
						if err != nil {
							log.Errorf("failed to write event to tcp: %s", err)
						} else if n != len(buf)+1 {
							log.Errorf("partial event write, %d bytes written on %d", n, len(buf))
						}
					}
				}

				wg.Done()
			}()

			<-sigChan
			gen.Terminate()
			done = true

			// Wait for generator goroutine completion
			wg.Wait()
		},
	}

	tcpCmd.Flags().StringVarP(&remoteAddr, "host", "H", "localhost", "remote host address")
	tcpCmd.Flags().Uint16VarP(&remotePort, "port", "p", 3333, "remote port")

	return tcpCmd
}
