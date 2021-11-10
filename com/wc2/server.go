package wc2

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/iDigitalFlame/xmt/com"
	"github.com/iDigitalFlame/xmt/device/devtools"
	"github.com/iDigitalFlame/xmt/util/xerr"
)

// Server is a C2 profile that mimics a standard web server and client setup.
// This struct inherits the http.Server struct and can be used to serve real files and pages.
//
// Use the Target Rules a URL mapping that can be used by clients to access the C2 functions.
type Server struct {
	Target  Target
	tls     *tls.Config
	dialer  *net.Dialer
	handler *http.ServeMux
	ctx     context.Context
	cancel  context.CancelFunc

	Client *http.Client
	rules  []Rule
}
type directory string

// Close terminates this Web instance and signals all current listeners and clients to disconnect. This
// will close all connections related to this struct.
//
// This function does not block.
func (s *Server) Close() error {
	s.cancel()
	return nil
}

// TargetAsRule will take the server 'Target' and create a matching ruleset using the 'Rule' function.
func (s *Server) TargetAsRule() {
	if s.Target.empty() {
		return
	}
	s.Rule(s.Target.Rule())
}

// Rule adds the specified rules to the Web instance to assist in determing real and C2 traffic.
func (s *Server) Rule(r ...Rule) {
	if len(r) == 0 {
		return
	}
	if s.rules == nil {
		s.rules = make([]Rule, 0, len(r))
	}
	s.rules = append(s.rules, r...)
}

// New creates a new Web C2 server instance. This can be passed to the Listen function of a controller to
// serve a Web Server that also acts as a C2 instance.
//
// This struct supports all the default Golang http.Server functions and can be used to serve real web pages.
// Rules must be defined using the 'Rule' function to allow the server to differentiate between C2 and real traffic.
func New(t time.Duration) *Server {
	return NewTLS(t, nil)
}

// Serve attempts to serve the specified filesystem path 'f' at the URL mapped path 'p'. This function will
// determine if the path represents a file or directory and will call 'ServeFile' or 'ServeDirectory' depending
// on the path result. This function will return an error if the filesystem path does not exist or is invalid.
func (s *Server) Serve(p, f string) error {
	i, err := os.Stat(f)
	if err != nil {
		return err
	}
	if i.IsDir() {
		s.handler.Handle(p, http.FileServer(http.Dir(f)))
		return nil
	}
	s.handler.Handle(p, http.FileServer(directory(f)))
	return nil
}

// ServeFile attempts to serve the specified filesystem path 'f' at the URL mapped path 'p'. This function is used
// to serve files and will return an error if the filesystem path does not exist or the path destination is not
// a file.
func (s *Server) ServeFile(p, f string) error {
	i, err := os.Stat(f)
	if err != nil {
		return err
	}
	if !i.IsDir() {
		s.handler.Handle(p, http.FileServer(directory(f)))
		return nil
	}
	return xerr.New(`path "` + f + `" is not a file`)
}

// Handle registers the handler for the given pattern. If a handler already exists for pattern, Handle panics.
func (s *Server) Handle(p string, h http.Handler) {
	s.handler.Handle(p, h)
}

// ServeDirectory attempts to serve the specified filesystem path 'f' at the URL mapped path 'p'. This function is used
// to serve directories and will return an error if the filesystem path does not exist or the path destination is not
// a directory.
func (s *Server) ServeDirectory(p, f string) error {
	i, err := os.Stat(f)
	if err != nil {
		return err
	}
	if i.IsDir() {
		s.handler.Handle(p, http.FileServer(http.Dir(f)))
		return nil
	}
	return xerr.New(`path "` + f + `" is not a directory`)
}

// NewTLS creates a new TLS wrapped Web C2 server instance. This can be passed to the Listen function of a Controller
// to serve a Web Server that also acts as a C2 instance. This struct supports all the default Golang http.Server
// functions and can be used to serve real web pages. Rules must be defined using the 'Rule' function to allow the
// server to differentiate between C2 and real traffic.
func NewTLS(t time.Duration, c *tls.Config) *Server {
	w := &Server{
		tls:     c,
		dialer:  &net.Dialer{Timeout: t, KeepAlive: t, DualStack: true},
		handler: new(http.ServeMux),
	}
	w.ctx, w.cancel = context.WithCancel(context.Background())
	w.Client = &http.Client{
		Timeout: t,
		Transport: &http.Transport{
			Proxy:                 devtools.Proxy,
			DialContext:           w.dialer.DialContext,
			MaxIdleConns:          256,
			IdleConnTimeout:       w.dialer.Timeout,
			ForceAttemptHTTP2:     false,
			TLSHandshakeTimeout:   w.dialer.Timeout,
			ExpectContinueTimeout: w.dialer.Timeout,
			ResponseHeaderTimeout: w.dialer.Timeout,
			ReadBufferSize:        1,
			WriteBufferSize:       1,
			DisableKeepAlives:     true,
		},
	}
	return w
}

// Connect creates a C2 client connector that uses the same properties of the Web struct parent.
func (s *Server) Connect(a string) (net.Conn, error) {
	return connect(&s.Target, s.Client, a)
}
func (d directory) Open(_ string) (http.File, error) {
	return os.OpenFile(string(d), os.O_RDONLY, 0)
}

// Listen returns a new C2 listener for this Web instance. This function creates a separate server, but still
// shares the handler for the base Web instance that it's created from.
func (s *Server) Listen(a string) (net.Listener, error) {
	if s.tls != nil && (len(s.tls.Certificates) == 0 || s.tls.GetCertificate == nil) {
		return nil, com.ErrInvalidTLSConfig
	}
	x, err := com.ListenConfig.Listen(s.ctx, "tcp", a)
	if err != nil {
		return nil, err
	}
	l := &listener{
		p:     s,
		ch:    make(chan complete, 1),
		new:   make(chan *conn, 256),
		ctx:   s.ctx,
		rules: make([]Rule, len(s.rules)),
		Server: &http.Server{
			Addr:              a,
			TLSConfig:         s.tls,
			ReadTimeout:       s.dialer.Timeout,
			IdleTimeout:       0,
			WriteTimeout:      s.dialer.Timeout,
			ReadHeaderTimeout: s.dialer.Timeout,
		},
	}
	l.Handler, l.BaseContext = l, l.context
	if copy(l.rules, s.rules); s.tls != nil {
		go l.listen(tls.NewListener(x, s.tls))
	} else {
		go l.listen(x)
	}
	return l, nil
}

// HandleFunc registers the handler function for the given pattern.
func (s *Server) HandleFunc(p string, h func(http.ResponseWriter, *http.Request)) {
	s.handler.HandleFunc(p, h)
}
