package cmd

import (
	"encoding/json"
	"fmt"
	"homework-event-generator/event"
	"homework-event-generator/event/auth"
	"homework-event-generator/output"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var rootCmd = &cobra.Command{
	Use:   "eventgen",
	Short: "Event generator produces fake events for testing and homeworks",
	Long: `
Event generator mimicks security appliances and generate structured json log events for use 
in synthetic testing or homework scenario`,
	Run: func(cmd *cobra.Command, args []string) {

		jitter := event.NewJitterTicker(time.Millisecond*100, time.Second*4)

		gen := event.NewGenerator()
		err := gen.SetRandomSeed(424242421)
		if err != nil {
			panic(err)
		}

		err = gen.RegisterPatternFactory(auth.NewForeignLoginFactory(jitter), 0.1, 1)
		if err != nil {
			panic(err)
		}

		err = gen.RegisterPatternFactory(auth.NewAccountLockedFactory(jitter), 0.2, 1)
		if err != nil {
			panic(err)
		}

		wg := sync.WaitGroup{}

		go func() {
			wg.Add(1)
			err := gen.Run(
				event.WrapTimeTicker(time.NewTicker(10*time.Millisecond)),
				event.WrapTimeTicker(time.NewTicker(100*time.Millisecond)),
			)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}()

		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, os.Interrupt)

		addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:3333")
		if err != nil {
			panic(err)
		}

		tcp, err := output.NewTcpWriteCloser(addr)
		if err != nil {
			panic(err)
		}

		done := false
		for !done {
			select {
			case <-sigChan:
				gen.Terminate()
				done = true
			case ev := <-gen.Output():
				buf, err := json.Marshal(ev)
				if err != nil {
					panic(err)

				}

				log.Debugln(string(buf))
				n, err := tcp.Write(append(buf, '\n'))
				if err != nil {
					log.Errorf("failed to write event to tcp: %s", err)
				} else if n != len(buf)+1 {
					log.Errorf("partial event write, %d bytes written on %d", n, len(buf))
				}

			}
		}
		wg.Wait()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
