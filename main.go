package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/LeeEirc/tclientlib"
	"github.com/gliderlabs/ssh"
	flags "github.com/jessevdk/go-flags"
)

type options struct {
	Addr    string `short:"a" long:"addr" description:"Address to listen on" default:"2222"`
	HostKey string `short:"k" long:"key"  description:"Path to the host key"`
}

func start(opts options) error {
	ssh.Handle(func(s ssh.Session) {
		_, _, isPty := s.Pty()
		if isPty {
			addr := net.JoinHostPort(s.User(), "23")
			fmt.Printf("Connecting to %s\n", addr)

			client, err := tclientlib.Dial("tcp", addr, &tclientlib.Config{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to connect to %s.\n", addr)
				s.Exit(1)
			} else {
				sigChan := make(chan struct{}, 1)

				go func() {
					_, _ = io.Copy(s, client)
					sigChan <- struct{}{}
				}()
				go func() {
					_, _ = io.Copy(client, s)
					sigChan <- struct{}{}
				}()

				<-sigChan
			}
		} else {
			fmt.Fprintf(os.Stderr, "No PTY requested.\n")
			s.Exit(1)
		}
	})

	var hostKey ssh.Option
	if _, err := os.Stat(opts.HostKey); err == nil {
		hostKey = ssh.HostKeyFile("id_rsa")
	} else {
		hostKey = func(srv *ssh.Server) error { return nil }
	}

	fmt.Printf("Starting ssh server on %s\n", opts.Addr)
	return ssh.ListenAndServe(opts.Addr, nil, hostKey)
}

func main() {
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		log.Fatal(err)
	}

	if err := start(opts); err != nil {
		log.Fatal(err)
	}
}
