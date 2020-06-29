package device

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"github.com/iDigitalFlame/xmt/data"
)

const (
	// Windows represents the Windows family of Operating Systems.
	Windows deviceOS = 0x0
	// Linux represents the Linux family of Operating Systems
	Linux deviceOS = 0x1
	// Unix represents the Unix family of Operating Systems
	Unix deviceOS = 0x2
	// Mac represents the MacOS/BSD family of Operating Systems
	Mac deviceOS = 0x3

	// Arch64 represents the 64-bit chipset family.
	Arch64 deviceArch = 0x0
	// Arch86 represents the 32-bit chipset family.
	Arch86 deviceArch = 0x1
	// ArchARM represents the ARM chipset family.
	ArchARM deviceArch = 0x2
	// ArchPowerPC represents the PowerPC chipset family.
	ArchPowerPC deviceArch = 0x3
	// ArchMips represents the MIPS chipset family.
	ArchMips deviceArch = 0x4
	// ArchUnknown represents an unknown chipset family.
	ArchUnknown deviceArch = 0x5

	// IDSize is the amount of bytes used to store the Host ID and
	// SessionID values.  The ID is the (HostID + SessionID).
	IDSize = 32
	// SmallIDSize is the amount of bytes used for printing the Host ID
	// value using the ID function.
	SmallIDSize = MachineIDSize
	// MachineIDSize is the amount of bytes that is used as the Host
	// specific ID value that does not change when on the same host.
	MachineIDSize = 28
)

const (
	xmtID              = "xmtFramework"
	xmtIDPrime  uint32 = 16777619
	xmtIDOffset uint32 = 2166136261
)

var (
	// ErrNoWindows is an error that is returned when a non-Windows device attempts a Windows specific function.
	ErrNoWindows = errors.New("not supported on non-Windows devices")

	envRegexp = regexp.MustCompile(`(%([\w\d()-_]+)%|\$([[\w\d-_]+))`)
)

// ID is an alias for a byte array that represents a 48 byte client identification number. This is used for
// tracking and detection purposes.
type ID []byte
type deviceOS uint8
type deviceArch uint8

func getID() ID {
	var (
		i      = ID(make([]byte, IDSize))
		s, err = machineid.ProtectedID(xmtID)
	)
	if err == nil {
		copy(i, s)
	} else {
		rand.Read(i)
	}
	rand.Read(i[MachineIDSize:])
	return i
}
func getArch() deviceArch {
	switch runtime.GOARCH {
	case "386":
		return Arch86
	case "amd64", "amd64p32":
		return Arch64
	case "ppc", "ppc64", "ppc64le":
		return ArchPowerPC
	case "arm", "armbe", "arm64", "arm64be":
		return ArchARM
	case "mips", "mipsle", "mips64", "mips64le", "mips64p32", "mips64p32le":
		return ArchMips
	}
	return ArchUnknown
}

// Full returns the full string representation of this ID instance. Full is an alias of the 'FullString' function.
func (i ID) Full() string {
	return i.FullString()
}

// Hash returns the 32bit hash sum of this ID value. The hash mechanism used is similar to the hash/fnv mechanism.
func (i ID) Hash() uint32 {
	h := xmtIDOffset
	for x := range i {
		h *= xmtIDPrime
		h ^= uint32(i[x])
	}
	return h
}

// String returns a representation of this ID instance.
func (i ID) String() string {
	if len(i) < SmallIDSize {
		return strings.ToUpper(hex.EncodeToString(i))
	}
	return strings.ToUpper(hex.EncodeToString(i[SmallIDSize:]))
}

// Expand attempts to determine environment variables from the current session and translate them from
// the supplied string.
func Expand(s string) string {
	v := envRegexp.FindAllStringSubmatch(s, -1)
	if len(v) == 0 {
		return s
	}
	for _, i := range v {
		if len(i) != 4 {
			continue
		}
		n := i[2]
		if len(i[3]) > 0 {
			n = i[3]
		}
		if d, ok := Environment[strings.ToLower(n)]; ok {
			s = strings.ReplaceAll(s, i[0], d)
		}
	}
	return s
}

// FullString returns the full string representation of this ID instance.
func (i ID) FullString() string {
	return strings.ToUpper(hex.EncodeToString(i))
}
func getEnv() map[string]string {
	m := make(map[string]string)
	for _, v := range os.Environ() {
		if i := strings.IndexRune(v, '='); i > 0 {
			m[strings.ToLower(v[:i])] = v[i+1:]
		}
	}
	t := os.TempDir()
	m["tmp"], m["temp"], m["tmpdir"], m["tempdir"] = t, t, t, t
	return m
}

// LoadSession will attempt to load the Session UUID from the specified file. This function will return an
// error if the file cannot be read or not found.
func LoadSession(s string) error {
	r, err := os.Open(s)
	if err != nil {
		return err
	}
	n, err := data.ReadFully(r, UUID)
	if r.Close(); err != nil {
		return err
	}
	if n != IDSize {
		return io.EOF
	}
	return nil
}

// SaveSession will attempt to save the Session UUID to the specified file. This function will return an
// error if the file cannot be written to or created.
func SaveSession(s string) error {
	w, err := os.Create(s)
	if err != nil {
		return err
	}
	_, err = w.Write(UUID)
	if w.Close(); err != nil {
		return err
	}
	return nil
}
func (d deviceOS) String() string {
	switch d {
	case Windows:
		return "Windows"
	case Linux:
		return "Linux"
	case Unix:
		return "Unix/BSD"
	case Mac:
		return "MacOS"
	}
	return "Unknown"
}
func (d deviceArch) String() string {
	switch d {
	case Arch86:
		return "32bit"
	case Arch64:
		return "64bit"
	case ArchARM:
		return "ARM"
	case ArchMips:
		return "MIPS"
	case ArchPowerPC:
		return "PowerPC"
	}
	return "Unknown"
}

// IDFromString attempts to convert the hex string supplied into an ID value.
func IDFromString(s string) (ID, error) {
	i, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	if len(i) == 0 {
		return nil, fmt.Errorf("invalid ID value %q", s)
	}
	switch len(i) {
	case 4, MachineIDSize, IDSize:
		break
	default:
		return nil, fmt.Errorf("ID size %d is invalid", len(i))
	}
	return ID(i), nil
}

// MarshalStream transform this struct into a binary format and writes to the supplied data.Writer.
func (i ID) MarshalStream(w data.Writer) error {
	_, err := w.Write(i)
	return err
}

// UnmarshalStream transforms this struct from a binary format that is read from the supplied data.Reader.
func (i *ID) UnmarshalStream(r data.Reader) error {
	if *i == nil {
		*i = append(*i, make([]byte, IDSize)...)
	}
	n, err := data.ReadFully(r, *i)
	if err != nil {
		return err
	}
	if n != IDSize {
		return io.EOF
	}
	return nil
}