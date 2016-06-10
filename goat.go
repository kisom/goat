// goat is a basic netcat clone.
package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

func dieIf(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] %s\n", err)
		os.Exit(1)
	}
}

func usage(w io.Writer) {
	fmt.Fprintf(w, `goat is a basic netcat-like program.
	
Usage: goat [-4] [-6] [-k] [-l] [address] port

Flags:
	-4	IPv4 only
	-6	IPv6 only
	-k	Keep the listener open
	-l	Listen on the specified host/port.
`)
}

func listener(l net.Listener) {
	conn, err := l.Accept()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] %s\n", err)
		conn.Close()
	}

	_, err = io.Copy(os.Stdout, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] %s\n", err)
		conn.Close()
	}
}

func init() {
	flag.Usage = func() { usage(os.Stdout) }
}

func main() {
	var addr, network, port string
	var keep, listen bool
	var ipv4, ipv6 bool

	flag.BoolVar(&ipv4, "4", false, "IPv4 only")
	flag.BoolVar(&ipv6, "6", false, "IPv6 only")
	flag.BoolVar(&keep, "k", false, "listen")
	flag.BoolVar(&listen, "l", false, "listen")
	flag.Parse()

	switch {
	case ipv4 && ipv6:
		fmt.Fprintf(os.Stderr, "[!] -4 and -6 are mutually exclusive.\n")
		os.Exit(1)
	case ipv4:
		network = "tcp4"
	case ipv6:
		network = "tcp6"
	default:
		network = "tcp"
	}

	args := flag.Args()
	if len(args) == 0 || len(args) > 2 {
		usage(os.Stdout)
		return
	}

	if !listen && len(args) == 1 {
		fmt.Fprintf(os.Stderr, "[!] need a host and port to connect to.\n")
		usage(os.Stderr)
		os.Exit(1)
	}

	if len(args) == 1 {
		port = args[0]
	} else {
		addr = args[0]
		port = args[1]
	}
	addr = net.JoinHostPort(addr, port)

	if listen {
		l, err := net.Listen("tcp", addr)
		dieIf(err)
		defer l.Close()

		for {
			listener(l)
			if !keep {
				break
			}
		}
	} else {
		conn, err := net.Dial(network, addr)
		dieIf(err)
		defer conn.Close()

		_, err = io.Copy(conn, os.Stdin)
		dieIf(err)
	}
}
