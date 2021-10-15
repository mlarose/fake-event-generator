package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:3333")
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	terminated := false
	go func() {
		wg.Add(1)
		defer wg.Done()
		for !terminated {
			conn, err := l.Accept()
			if err != nil && err != net.ErrClosed {
				fmt.Println(err)
				return
			}

			// handle connection
			go func() {
				err := handleConnection(conn)
				if err != nil && err != io.EOF {
					fmt.Println("err: ", err)
				}
				conn.Close()
			}()
		}
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, os.Interrupt)

	<-termChan
	terminated = true
	l.Close()
	wg.Wait()
}

func handleConnection(conn net.Conn) error {
	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return err
		}

		fmt.Printf(data)
	}
}
