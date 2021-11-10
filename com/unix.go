package com

import (
	"crypto/tls"
	"net"
	"time"
)

type unixConnector tcpConnector

// NewUNIX creates a new simple UNIX socket based connector with the supplied timeout.
func NewUNIX(t time.Duration) Connector {
	return unixConnector(tcpConnector{Dialer: net.Dialer{Timeout: t, KeepAlive: t, DualStack: true}})
}
func (u unixConnector) Connect(s string) (net.Conn, error) {
	return newStreamConn("unix", s, tcpConnector(u))
}

// NewSecureUNIX creates a new simple TLS wrapped UNIX socket based connector with the supplied timeout.
func NewSecureUNIX(t time.Duration, c *tls.Config) Connector {
	if t < 0 {
		t = DefaultTimeout
	}
	return unixConnector(tcpConnector{tls: c, Dialer: net.Dialer{Timeout: t, KeepAlive: t, DualStack: true}})
}
func (u unixConnector) Listen(s string) (net.Listener, error) {
	c, err := newStreamListener("unix", s, tcpConnector(u))
	if err != nil {
		return nil, err
	}
	return &tcpListener{timeout: u.Timeout, Listener: c}, nil
}
