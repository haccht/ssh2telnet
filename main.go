package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/gliderlabs/ssh"
	"github.com/ziutek/telnet"

	flags "github.com/jessevdk/go-flags"
)

type options struct {
	Addr    string `short:"a" long:"addr" description:"Address to listen on" default:":2222"`
	HostKey string `short:"k" long:"key"  description:"Path to the host key"`
}

func start(opts options) error {
	server := &ssh.Server{Addr: opts.Addr}

	if _, err := os.Stat(opts.HostKey); err == nil {
		hostKeyFile := ssh.HostKeyFile(opts.HostKey)
		server.SetOption(hostKeyFile)
	}

	server.Handle(func(s ssh.Session) {
		_, _, isPty := s.Pty()
		if isPty {
			addr := net.JoinHostPort(s.User(), "23")
			fmt.Printf("Connecting to %s\n", addr)

			conn, err := telnet.Dial("tcp", addr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to connect to %s.\n", addr)
				s.Exit(1)
			} else {
				sigChan := make(chan struct{}, 1)

				go func() {
					_, _ = io.Copy(s, conn)
					sigChan <- struct{}{}
				}()
				go func() {
					_, _ = io.Copy(conn, s)
					sigChan <- struct{}{}
				}()

				<-sigChan
			}
		} else {
			fmt.Fprintf(os.Stderr, "No PTY requested.\n")
			s.Exit(1)
		}
	})

	fmt.Printf("Starting ssh server on %s\n", opts.Addr)
	return server.ListenAndServe()
}

func main() {
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		if fe, ok := err.(*flags.Error); ok && fe.Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatal(err)
	}

	if err := start(opts); err != nil {
		log.Fatal(err)
	}
}
