package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/ziutek/telnet"

	flags "github.com/jessevdk/go-flags"
)

type options struct {
	Addr           string `short:"a" long:"addr" description:"Address to listen on" default:":2222"`
	HostKey        string `short:"k" long:"key" description:"Path to the host key"`
	LoginPrompt    string `long:"login-prompt" description:"Login prompt" default:"\"login: \""`
	PasswordPrompt string `long:"password-prompt" description:"Password prompt" default:"\"Password: \""`
}

func start(opts options) error {
	server := &ssh.Server{Addr: opts.Addr}

	if _, err := os.Stat(opts.HostKey); err == nil {
		hostKeyFile := ssh.HostKeyFile(opts.HostKey)
		server.SetOption(hostKeyFile)
	}

	var username, password, hostname string
	passwordAuth := ssh.PasswordAuth(func(ctx ssh.Context, s string) bool {
		t := strings.SplitN(ctx.User(), "@", 2)
		if len(t) != 2 {
			return false
		}

		username, hostname = t[0], t[1]
		password = s
		return true
	})
	server.SetOption(passwordAuth)

	server.Handle(func(s ssh.Session) {
		_, _, isPty := s.Pty()
		if isPty {
			addr := net.JoinHostPort(hostname, "23")
			fmt.Printf("Connecting to %s\n", addr)

			conn, err := telnet.Dial("tcp", addr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to connect to %s.\n", addr)
				s.Exit(1)
			} else {
				_, err = conn.ReadUntil(opts.LoginPrompt)
				conn.Write([]byte(fmt.Sprintf("%s\n", username)))
				_, err = conn.ReadUntil(opts.PasswordPrompt)
				conn.Write([]byte(fmt.Sprintf("%s\n", password)))

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
