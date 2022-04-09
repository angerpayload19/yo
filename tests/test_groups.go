package main

import (
	"os"
	"strconv"
	"time"

	"github.com/PurpleSec/logx"
	"github.com/iDigitalFlame/xmt/c2"
	"github.com/iDigitalFlame/xmt/c2/cfg"
	// Uncomment to enable profiling.
	// "net/http"
	// _ "net/http/pprof"
)

func testGroups() {
	// This example shows off the multi-Profile group!
	//
	// The client process will reconnect to any listener that first
	// accepts it in a 'RoundRobin' fashion.
	//
	// Try it by running both server and client then using iptables
	// to filter out the listener ports (8085, 8086, 8087) randomally.

	// Create initial config
	c := cfg.Pack(
		cfg.Host("127.0.0.1:8085"),
		cfg.ConnectTCP,
		cfg.Sleep(time.Second*5),
		cfg.Jitter(0),
	)

	// Add a Group!
	c.AddGroup(
		cfg.Host("127.0.0.1:8086"),
		cfg.ConnectTCP,
		cfg.Sleep(time.Second*10),
		cfg.Jitter(50),
	)

	// Add a Group!
	c.AddGroup(
		cfg.Host("127.0.0.1:8087"),
		cfg.ConnectTCP,
		cfg.Sleep(time.Second*5),
		cfg.Jitter(0),
	)

	// Use the last valid group, unless an error happens.
	c.Add(cfg.SelectorLastValid)

	if len(os.Args) == 1 {
		server(c)
		return
	}
	client(c)
}

func server(c cfg.Config) {
	// Uncomment to enable profiling.
	// go http.ListenAndServe("localhost:9090", nil)

	s := c2.NewServer(logx.Console(logx.Trace))

	// Loop over the config Groups.
	for i := 0; i < c.Groups(); i++ {
		// Get group at section and Build it.
		// We're ignoring the error since we KNOW what the result is.
		var (
			v    = c.Group(i)
			p, _ = v.Build()
		)
		// Create a Listener with the name "tcp<offset>".
		if _, err := s.Listen("", "tcp"+strconv.Itoa(i+1), p); err != nil {
			panic(err)
		}
	}
	// Block
	s.Wait()
}
func client(c cfg.Config) {
	// Uncomment to enable profiling.
	// go http.ListenAndServe("localhost:9091", nil)

	var (
		p, _   = c.Build()
		s, err = c2.Connect(logx.Console(logx.Trace), p)
	)
	if err != nil {
		panic(err)
	}

	// Block
	s.Wait()
}
