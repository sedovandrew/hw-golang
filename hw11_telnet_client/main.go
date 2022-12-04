package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"
)

const (
	minPort = 0
	maxPort = 65535

	defaultTimeout = time.Duration(10 * time.Second)

	errorReturnCode = 1
)

func logf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "..."+format, a...)
}

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", defaultTimeout, "Connection timeout")
	flag.Parse()

	if len(flag.Args()) != 2 {
		logf("Need 'host' and 'port' arguments\n")
		os.Exit(errorReturnCode)
	}
	host, port := flag.Arg(0), flag.Arg(1)

	// Validate IP address
	if net.ParseIP(host) == nil {
		_, err := net.LookupHost(host)
		if err != nil {
			logf("Error resolve: %v\n", host)
			os.Exit(errorReturnCode)
		}
	}

	// Validate port
	portInt, err := strconv.Atoi(port)
	if err != nil {
		logf("Port is not digit: %s\n", port)
		os.Exit(errorReturnCode)
	}
	if portInt < minPort || portInt > maxPort {
		logf("Wrong port range (%d-%d): %s\n", minPort, maxPort, port)
		os.Exit(errorReturnCode)
	}

	// Connect to server
	hostPort := net.JoinHostPort(host, port)
	telnetClient := NewTelnetClient(hostPort, timeout, os.Stdin, os.Stdout)
	err = telnetClient.Connect()
	if err != nil {
		logf("Error connecting: %v\n", err)
		os.Exit(errorReturnCode)
	}
	logf("Connected to %s\n", hostPort)

	// Start Receiver
	receiverErrCh := make(chan error, 1)
	go func() {
		defer close(receiverErrCh)
		receiverErrCh <- telnetClient.Receive()
	}()

	// Start Sender
	senderErrCh := make(chan error, 1)
	go func() {
		defer close(senderErrCh)
		senderErrCh <- telnetClient.Send()
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	select {
	case <-signalCh:
		logf("Interrupt\n")
	case err = <-receiverErrCh:
		if err != nil {
			logf("Receiver error: %v\n", err)
		}
	case err = <-senderErrCh:
		if err != nil {
			logf("Sender error: %v\n", err)
		}
	}
	close(signalCh)

	err = telnetClient.Close()
	if err != nil {
		logf("Error closing connection: %v\n", err)
		os.Exit(errorReturnCode)
	}
	logf("Disconnected\n")
}
