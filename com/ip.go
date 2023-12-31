// Copyright (C) 2020 - 2023 iDigitalFlame
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package com

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/iDigitalFlame/xmt/util"
	"github.com/iDigitalFlame/xmt/util/bugtrack"
)

type ipStream struct {
	udpStream
}
type ipListener struct {
	_ [0]func()
	net.Listener
	proto byte
}
type ipConnector struct {
	net.Dialer
	proto byte
}
type ipPacketConn struct {
	_ [0]func()
	net.PacketConn
}

// NewIP creates a new simple IP based connector with the supplied timeout and
// protocol number.
func NewIP(t time.Duration, p byte) Connector {
	return &ipConnector{proto: p, Dialer: net.Dialer{Timeout: t, KeepAlive: t}}
}
func (i *ipStream) Read(b []byte) (int, error) {
	n, err := i.udpStream.Read(b)
	if n > 20 {
		if bugtrack.Enabled {
			bugtrack.Track("com.(*ipStream).Read(): Cutting off IP header n=%d, after n=%d", n, n-20)
		}
		copy(b, b[20:])
		n -= 20
	}
	if err == nil && n < len(b)-20 {
		err = io.EOF
	}
	return n, err
}
func (i *ipConnector) Connect(x context.Context, s string) (net.Conn, error) {
	c, err := i.DialContext(x, NameIP+":"+util.Uitoa(uint64(i.proto)), s)
	if err != nil {
		return nil, err
	}
	return &ipStream{udpStream{Conn: c}}, nil
}
func (i *ipConnector) Listen(x context.Context, s string) (net.Listener, error) {
	c, err := ListenConfig.ListenPacket(x, NameIP+":"+util.Uitoa(uint64(i.proto)), s)
	if err != nil {
		return nil, err
	}
	l := &udpListener{
		del:  make(chan udpAddr, 16),
		new:  make(chan *udpConn, 16),
		cons: make(map[udpAddr]*udpConn),
		sock: &udpCompat{&ipPacketConn{PacketConn: c}},
	}
	l.ctx, l.cancel = context.WithCancel(x)
	go l.purge()
	go l.listen()
	return &ipListener{proto: i.proto, Listener: l}, nil
}
