#!/usr/bin/python3
# Copyright (C) 2021-2022 iDigitalFlame
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.
#

from zlib import crc32
from hashlib import sha512
from secrets import token_bytes
from traceback import format_exc
from argparse import ArgumentParser
from base64 import b64decode, b64encode
from json import dumps, loads, JSONDecodeError
from sys import argv, exit, stderr, stdin, stdout


# System
S_HOST = 0xA0
S_SLEEP = 0xA1
S_JITTER = 0xA2

# Connectors
C_TCP = 0xC0
C_TLS = 0xC1
C_UDP = 0xC2
C_ICMP = 0xC3
C_PIPE = 0xC4
C_TLS_INSECURE = 0xC5

# Custom Connectors
C_IP = 0xB0
C_WC2 = 0xB1
C_TLS_EX = 0xB2
C_MTLS = 0xB3
C_TLS_CA = 0xB4
C_TLS_CERT = 0xB5

# Wrappers
W_HEX = 0xD0
W_ZLIB = 0xD1
W_GZIP = 0xD2
W_B64 = 0xD3
W_XOR = 0xD4
W_CBK = 0xD5
W_AES = 0xD6

# Transforms
T_B64 = 0xE0
T_DNS = 0xE1
T_B64S = 0xE2

WRAPPERS = ["hex", "zlib", "gzip", "b64", "xor", "aes", "cbk"]

NAMES = {
    S_HOST: "host",
    S_SLEEP: "sleep",
    S_JITTER: "jitter",
    C_TCP: "tcp",
    C_TLS: "tls",
    C_UDP: "udp",
    C_ICMP: "icmp",
    C_PIPE: "pipe",
    C_TLS_INSECURE: "tls-insecure",
    C_IP: "ip",
    C_WC2: "wc2",
    C_TLS_EX: "tls-ex",
    C_MTLS: "mtls",
    C_TLS_CA: "tls-ca",
    C_TLS_CERT: "tls-cert",
    W_HEX: "hex",
    W_ZLIB: "zlib",
    W_GZIP: "gzip",
    W_B64: "base64",
    W_XOR: "xor",
    W_CBK: "cbk",
    W_AES: "aes",
    T_B64: "base64t",
    T_DNS: "dns",
    T_B64S: "b64s",
}
UNITS = {
    "ns": 1,
    "us": 1000,
    "µs": 1000,
    "μs": 1000,
    "ms": 1000000,
    "s": 1000000000,
    "m": 60000000000,
    "h": 3600000000000,
}
NAMES_TO_ID = {
    "host": S_HOST,
    "sleep": S_SLEEP,
    "jitter": S_JITTER,
    "tcp": C_TCP,
    "tls": C_TLS,
    "udp": C_UDP,
    "icmp": C_ICMP,
    "pipe": C_PIPE,
    "tls-insecure": C_TLS_INSECURE,
    "ip": C_IP,
    "wc2": C_WC2,
    "tls-ex": C_TLS_EX,
    "mtls": C_MTLS,
    "tls-ca": C_TLS_CA,
    "tls-cert": C_TLS_CERT,
    "hex": W_HEX,
    "zlib": W_ZLIB,
    "gzip": W_GZIP,
    "base64": W_B64,
    "xor": W_XOR,
    "cbk": W_CBK,
    "aes": W_AES,
    "base64t": T_B64,
    "dns": T_DNS,
    "b64s": T_B64S,
}

BUILDER_HELP = """XMT cfg.Config Builder v1alpha

Basic Arguments:
  -h                            Show this help message and exit.
  --help

Input/Output Arguments:
  -f                <file>      Input file path. Use '-' for stdin.
  --in              <file>
  -o                <file>      Output file path. Stdout is used if
                                  empty.
  --out             <file>
  -j                            Output in JSON format. Omit for raw
                                  binary.
  --json

Operaton Arguments:
  -p                            List values contained in the file
                                  input. Fails if no input is found or
                                  invalid.
  --print

Build Arguments:
 System:
  --host            <hostname>  Hostname hint.
  --sleep           <secs|mod>  Sleep timeperiod. Defaults to seconds
                                  for integers, but can take modifiers
                                  such as 's', 'h', 'm'. (2m, 3s).
  --jitter          <jitter %>  Jitter as a percentage [0-100]. Values
                                  greater than 100 fail.

 Connection Hints (Max 1):
  --tcp                         Use the TCP Connection hint.
  --tls                         Use the TLS Connection hint.
  --udp                         Use the UDP Connection hint.
  --icmp                        Use the ICMP (Ping) Connection hint.
  --pipe                        Use the Windows Named Pipes Connection
                                  hint.
  -K                            Use the TLSNoVerify Connection hint.
  --tls-insecure
  --ip              <protocol>  Use the IP Connection hint with the
                                  specified protocol number [0-255].
  --wc2-url         <url>         Use the WC2 Connection hint with the
                                  URL expression or static string.
                                  This can be used with other WC2
                                  arguments without an error.
  --wc2-host        <host>      Use the WC2 Connection hint with the
                                  Host expression or static string.
                                  This can be used with other WC2
                                  arguments without an error.
  --wc2-agent       <agent>     Use the WC2 Connection hint with the
                                  User-Agent expression or static string.
                                  This can be used with other WC2
                                  arguments without an error.
  --wc2-header      <key>=<val> Use the WC2 Connection hint with the
                                  HTTP header expression or static string in
                                  a key=value formnat. This value will be
                                  parsed and will fail if 'key' is empty or
                                  no '=' is present in the string. This may be
                                  specified multiple times. This can be used
                                  with other WC2 arguments without an error.
  -H                <key>=<val>
  --mtls                        Use Mutual TLS Authentication (mTLS) with a TLS
                                  Connection hint. This just enables the flag for
                                  client auth and will fail if '--tls-pem' and
                                  '--tls-key' are empty or not specified.
  --tls-vers        <version>   Use the TLS version specified when using a TLS
                                  Connection hint. This will set the version
                                  required and can be used by itself. A value of
                                  zero means TLSv1. Can be used with other TLS
                                  options.
  --tls-ca          <file|pem>  Use the provided certificate to verify the server
                                  (for clients) or verify clients (for the server).
                                  Can be used on it's own and with '--mtls'.
                                  This argument can take a file path to a PEM
                                  formatted certificate or raw base64 encoded PEM
                                  data.
  --tls-pem         <file|pem>  Use the provided certificate for the generated
                                  TLS socket. This can be used for client or
                                  server listeners. Requires '--tls-key'.
                                  This argument can take a file path to a PEM
                                  formatted certificate or raw base64 encoded PEM
                                  data.
  --tls-key         <file|pem>  Use the provided certificate key for the generated
                                  TLS socket. This can be used for client or
                                  server listeners. Requires '--tls-pem'.
                                  This argument can take a file path to a PEM
                                  formatted certificate private key or raw
                                  base64 encoded PEM data.

 Wrappers (Multiple different types may be used):
  --hex                         Use the HEX Wrapper.
  --zlib                        Use the Zlib compression Wrapper.
  --gzip                        Use the Gzip compression Wrapper.
  --b64                         Use the Base64 encoding Wrapper.
  --xor             [key]       Encrypt with the XOR Wrapper using the provided
                                  key string. If omitted the key will be a
                                  randomally generated 64 byte array.
  --cbk             [A B C D]   Encrypt with the XOR Wrapper using the provided
                                  key string. If omitted the key will be a
                                  randomally generated 64 byte array.
  --aes             [key]       Encrypt with the AES Wrapper using the provided
                                  key string. If omitted the key will be a
                                  randomally generated 32 byte array. The AES IV
                                  may be supplied using the '--aes-iv' argument.
                                  If not specified a 16 byte IV will be generated.
  --aes-iv          [iv]        Encrypt with the AES Wrapper using the provided
                                  IV string. If omitted the IV will be a
                                  randomally generated 16 byte array. The AES key
                                  may be supplied using the '--aes-key' argument.
                                  If not specified a 32 byte key will be generated.

 Transforms (Max 1):
  --b64t            [shift]     Transform the data using a Base64 Transform. An
                                  option shift value [0-255] may be specified, but
                                  if omitted will not shift.
  --dns             [domain,*]  Use the DNS Packet Transform. optional DNS
                                  domain names may be specified (seperated by space)
                                  that will be used in the packets. This option may
                                  be used more than once to specify more domains.
  -D                [domain,*]
"""


def _dur_to_str(v):
    b = bytearray(32)
    n = len(b) - 1
    b[n] = ord("s")
    n, v = __fmt_frac(b, n, v)
    n = __fmt_int(b, n, v % 60)
    v /= 60
    if int(v) > 0:
        n -= 1
        b[n] = ord("m")
        n = __fmt_int(b, n, v % 60)
        v /= 60
        if int(v) > 0:
            n -= 1
            b[n] = ord("h")
            n = __fmt_int(b, n, v)
    return b[n:].decode("UTF-8")


def _str_to_dur(s):
    if not _nes(s):
        raise ValueError("str2dur: invalid duration")
    if s == "0":
        return 0
    d = 0
    while len(s) > 0:
        v = 0
        f = 0
        z = 1
        if not (s[0] == "." or (ord("0") <= ord(s[0]) and ord(s[0]) <= ord("9"))):
            raise ValueError("str2dur: invalid duration")
        p = len(s)
        v, s = __leading_int(s)
        r = p != len(s)
        y = False
        if len(s) > 0 and s[0] == ".":
            s = s[1:]
            p = len(s)
            f, z, s = __leading_fraction(s)
            y = p != len(s)
        if not r and not y:
            raise ValueError("str2dur: invalid duration")
        del r
        del y
        i = 0
        while i < len(s):
            c = ord(s[i])
            if c == ord(".") or (ord("0") <= c and c <= ord("9")):
                break
            i += 1
        if i == 0:
            # raise ValueError("str2dur: missing time unit")
            u = "s"
        else:
            u = s[:i]
            s = s[i:]
        del i
        if u not in UNITS:
            raise ValueError("str2dur: unknown unit")
        e = UNITS[u]
        del u
        if v > (((1 << 63) - 1) / e):
            raise ValueError("str2dur: invalid duration")
        v *= int(e)
        if f > 0:
            v += int(float(f) * float(float(e) / float(z)))
            if v < 0:
                raise ValueError("str2dur: invalid duration")
        del e
        d += v
        if d < 0:
            raise ValueError("str2dur: invalid duration")
        del v
        del f
        del z
    return d


def __leading_int(s):
    i = 0
    x = 0
    while i < len(s):
        c = ord(s[i])
        if c < ord("0") or c > ord("9"):
            break
        if x > (((1 << 63) - 1) / 10):
            raise OverflowError()
        x = int(x * 10) + int(c) - ord("0")
        if x < 0:
            raise OverflowError()
        i += 1
    return x, s[i:]


def __fmt_int(b, s, v):
    if int(v) == 0:
        s -= 1
        b[s] = ord("0")
        return s
    while int(v) > 0:
        s -= 1
        b[s] = int(v % 10) + ord("0")
        v /= 10
    return s


def __fmt_frac(b, s, v):
    p = False
    for _ in range(0, 9):
        d = v % 10
        p = p or d != 0
        if p:
            s -= 1
            b[s] = int(d) + ord("0")
        v /= 10
    if p:
        s -= 1
        b[s] = ord(".")
    del p
    return s, v


def __leading_fraction(s):
    i = 0
    x = 0
    v = 1
    o = False
    while i < len(s):
        c = ord(s[i])
        if c < ord("0") or c > ord("9"):
            break
        if o:
            continue
        if x > (((1 << 63) - 1) / 10):
            o = True
            continue
        y = int(x * 10) + int(c) - ord("0")
        if y < 0:
            o = True
            continue
        x = y
        v *= 10
        i += 1
        del y
    del o
    return x, v, s[1:]


def _nes(s, min=0, max=-1):
    if max > min:
        return isinstance(s, str) and len(s) < max and len(s) > min
    return isinstance(s, str) and len(s) > min


class Cfg(object):
    @staticmethod
    def host(v):
        if not _nes(v):
            raise ValueError("host: invalid name object")
        f = v.encode("UTF-8")
        n = len(f)
        if n > 0xFFFF:
            n = 0xFFFF
        s = Setting(3 + n)
        s[0] = S_HOST
        s[1] = (n >> 8) & 0xFF
        s[2] = n & 0xFF
        for x in range(0, n):
            s[x + 3] = f[x]
        del n
        del f
        return s

    @staticmethod
    def sleep(t):
        if _nes(t):
            t = _str_to_dur(t)
        if not isinstance(t, int) or t <= 0:
            raise ValueError("sleep: invalid duration")
        s = Setting(9)
        s[0] = S_SLEEP
        s[1] = (t >> 56) & 0xFF
        s[2] = (t >> 48) & 0xFF
        s[3] = (t >> 40) & 0xFF
        s[4] = (t >> 32) & 0xFF
        s[5] = (t >> 24) & 0xFF
        s[6] = (t >> 16) & 0xFF
        s[7] = (t >> 8) & 0xFF
        s[8] = t & 0xFF
        return s

    @staticmethod
    def single(v):
        if not isinstance(v, int):
            raise ValueError("single: invalid bit")
        s = Setting(1)
        s[0] = v
        return s

    @staticmethod
    def jitter(p):
        if _nes(p):
            if "%" in p:
                p = p.replace("%", "")
            try:
                p = int(p)
            except ValueError:
                raise ValueError("jitter: invalid percentage")
        if not isinstance(p, int) or p < 0 or p > 100:
            raise ValueError("jitter: invalid percentage")
        s = Setting(2)
        s[0] = S_JITTER
        s[1] = p & 0xFF
        return s

    @staticmethod
    def connect_ip(p):
        if not isinstance(p, int) or p <= 0 or p > 0xFF:
            raise ValueError("ip: invalid protocol")
        s = Setting(2)
        s[0] = C_IP
        s[1] = p & 0xFF
        return s

    @staticmethod
    def connect_tls_ex(v):
        if not isinstance(v, int) or v <= 0 or v > 0xFF:
            raise ValueError("tls-ex: invalid version")
        s = Setting(2)
        s[0] = C_TLS_EX
        s[1] = v & 0xFF
        return s

    @staticmethod
    def wrap_xor(key=None):
        if isinstance(key, list) and len(key) > 0:
            key = key[0]
        if key is None:
            key = token_bytes(64)
        elif _nes(key):
            key = key.encode("UTF-8")
        elif not isinstance(key, bytes) and not isinstance(key, bytearray):
            raise ValueError("xor: invalid KEY value")
        n = len(key)
        if n > 0xFFFF:
            n = 0xFFFF
        s = Setting(3 + n)
        s[0] = W_XOR
        s[1] = (n >> 8) & 0xFF
        s[2] = n & 0xFF
        for x in range(0, n):
            s[x + 3] = key[x]
        del n
        return s

    @staticmethod
    def connect_tls_ca(v, ca):
        if isinstance(ca, bytes) or isinstance(ca, bytearray):
            f = ca
        elif _nes(ca):
            f = ca.encode("UTF-8")
        else:
            raise ValueError("tls-ca: invalid CA")
        if not isinstance(v, int) or v <= 0 or v > 0xFF:
            raise ValueError("tls-ca: invalid version")
        n = len(f)
        if n > 0xFFFF:
            n = 0xFFFF
        s = Setting(4 + n)
        s[0] = C_TLS_CA
        s[1] = v & 0xFF
        s[2] = (n >> 8) & 0xFF
        s[3] = n & 0xFF
        for x in range(0, n):
            s[x + 4] = f[x]
        del f
        del n
        return s

    @staticmethod
    def transform_b64_shift(v):
        if not isinstance(v, int) or v <= 0 or v > 0xFF:
            raise ValueError("base64s: invalid shift")
        s = Setting(2)
        s[0] = T_B64S
        s[1] = v & 0xFF
        return s

    @staticmethod
    def wrap_aes(key=None, iv=None):
        if isinstance(iv, list) and len(iv) > 0:
            iv = iv[0]
        if isinstance(key, list) and len(key) > 0:
            key = key[0]
        if key is None:
            key = token_bytes(32)
        elif _nes(key):
            key = key.encode("UTF-8")
        elif not isinstance(key, bytes) and not isinstance(key, bytearray):
            raise ValueError("aes: invalid KEY value")
        if iv is None:
            iv = token_bytes(16)
        elif _nes(iv):
            iv = iv.encode("UTF-8")
        elif not isinstance(iv, bytes) and not isinstance(iv, bytearray):
            raise ValueError("aes: invalid IV value")
        if len(key) > 32:
            raise ValueError("aes: invalid KEY size")
        if len(iv) != 16:
            raise ValueError("aes: invalid IV size")
        s = Setting(3 + len(key) + len(iv))
        s[0] = W_AES
        s[1] = len(key) & 0xFF
        s[2] = len(iv) & 0xFF
        for x in range(0, len(key)):
            s[x + 3] = key[x]
        for x in range(0, len(iv)):
            s[x + len(key) + 3] = iv[x]
        return s

    @staticmethod
    def connect_mtls(v, ca, pem, key):
        if isinstance(ca, bytes) or isinstance(ca, bytearray):
            f = ca
        elif _nes(ca):
            f = ca.encode("UTF-8")
        else:
            raise ValueError("mtls: invalid CA")
        if isinstance(pem, bytes) or isinstance(pem, bytearray):
            if len(pem) == 0:
                raise ValueError("mtls: invalid PEM")
            p = pem
        elif _nes(pem):
            p = pem.encode("UTF-8")
        else:
            raise ValueError("mtls: invalid PEM")
        if isinstance(key, bytes) or isinstance(key, bytearray):
            if len(key) == 0:
                raise ValueError("mtls: invalid KEY")
            k = key
        elif _nes(key):
            k = key.encode("UTF-8")
        else:
            raise ValueError("mtls: invalid KEY")
        if not isinstance(v, int) or v <= 0 or v > 0xFF:
            raise ValueError("mtls invalid version")
        if len(p) == 0 or len(k) == 0:
            raise ValueError("mtls: invalid PEM or KEY version")
        o = len(f)
        if o > 0xFFFF:
            o = 0xFFFF
        n = len(p)
        if n > 0xFFFF:
            n = 0xFFFF
        m = len(k)
        if m > 0xFFFF:
            m = 0xFFFF
        s = Setting(8 + o + n + m)
        s[0] = C_MTLS
        s[1] = v & 0xFF
        s[2] = (o >> 8) & 0xFF
        s[3] = o & 0xFF
        s[2] = (n >> 8) & 0xFF
        s[3] = n & 0xFF
        s[4] = (m >> 8) & 0xFF
        s[5] = m & 0xFF
        for x in range(0, n):
            s[x + 8] = f[x]
        for x in range(0, n):
            s[x + o + 8] = p[x]
        for x in range(0, m):
            s[x + o + n + 8] = k[x]
        del f
        del p
        del k
        del o
        del n
        del m
        return s

    @staticmethod
    def connect_tls_certs(v, pem, key):
        if isinstance(pem, bytes) or isinstance(pem, bytearray):
            if len(pem) == 0:
                raise ValueError("tls-cert: invalid PEM")
            p = pem
        elif _nes(pem):
            p = pem.encode("UTF-8")
        else:
            raise ValueError("tls-cert: invalid PEM")
        if isinstance(key, bytes) or isinstance(key, bytearray):
            if len(key) == 0:
                raise ValueError("tls-cert: invalid KEY")
            k = key
        elif _nes(key):
            k = key.encode("UTF-8")
        else:
            raise ValueError("tls-cert: invalid KEY")
        if not isinstance(v, int) or v <= 0 or v > 0xFF:
            raise ValueError("tls-cert: invalid version")
        if len(p) == 0 or len(k) == 0:
            raise ValueError("tls-cert: invalid PEM or KEY version")
        n = len(p)
        if n > 0xFFFF:
            n = 0xFFFF
        m = len(k)
        if m > 0xFFFF:
            m = 0xFFFF
        s = Setting(6 + n + m)
        s[0] = C_TLS_CA
        s[1] = v & 0xFF
        s[2] = (n >> 8) & 0xFF
        s[3] = n & 0xFF
        s[4] = (m >> 8) & 0xFF
        s[5] = m & 0xFF
        for x in range(0, n):
            s[x + 6] = p[x]
        for x in range(0, m):
            s[x + n + 6] = k[x]
        del p
        del k
        del n
        del m
        return s

    @staticmethod
    def connect_wc2(u, h, a, head=None):
        if _nes(u):
            c = u.encode("UTF-8")
        else:
            c = bytearray()
        if _nes(h):
            v = h.encode("UTF-8")
        else:
            v = bytearray()
        if _nes(a):
            b = a.encode("UTF-8")
        else:
            b = bytearray()
        j = len(c)
        if j > 0xFFFF:
            j = 0xFFFF
        k = len(v)
        if k > 0xFFFF:
            k = 0xFFFF
        n = len(b)
        if n > 0xFFFF:
            n = 0xFFFF
        s = Setting(8 + j + k + n)
        s[0] = C_WC2
        s[1] = (j >> 8) & 0xFF
        s[2] = j & 0xFF
        s[3] = (k >> 8) & 0xFF
        s[4] = k & 0xFF
        s[5] = (n >> 8) & 0xFF
        s[6] = n & 0xFF
        for x in range(0, j):
            s[x + 8] = c[x]
        for x in range(0, k):
            s[x + j + 8] = v[x]
        for x in range(0, n):
            s[x + j + k + 8] = b[x]
        del j
        del k
        del n
        del c
        del v
        del b
        if not isinstance(head, dict):
            s[7] = 0
            return s
        i = 0
        s[7] = len(head) & 0xFF
        for k, v in head.items():
            if i >= 0xFF:
                break
            if not _nes(k):
                raise ValueError("wc2: invalid header")
            if _nes(v):
                z = v.encode("UTF-8")
            else:
                z = bytearray()
            o = k.encode("UTF-8")
            f = len(o)
            if f > 0xFF:
                f = 0xFF
            g = len(z)
            if g > 0xFF:
                g = 0xFF
            s.append(f & 0xFF)
            s.append(g & 0xFF)
            s.extend(o)
            s.extend(z)
            i += 1
            del o
            del z
            del f
            del g
        return s

    @staticmethod
    def wrap_cbk(a=None, b=None, c=None, d=None, size=128, key=None):
        if (
            not isinstance(a, int)
            and not isinstance(b, int)
            and not isinstance(c, int)
            and not isinstance(d, int)
        ):
            if _nes(key):
                v = key.encode("UTF-8")
            elif isinstance(key, bytes) or isinstance(key, bytearray):
                v = key
            else:
                v = token_bytes(64)
            if len(v) == 0:
                v - token_bytes(64)
            h = sha512()
            for _ in range(0, 256):
                h.update(v)
            del v
            n = crc32(h.digest()).to_bytes(4, byteorder="big", signed=False)
            del h
            a = n[0]
            b = n[1]
            c = n[2]
            d = n[3]
            del n
        if (
            not isinstance(a, int)
            or not isinstance(b, int)
            or not isinstance(c, int)
            or not isinstance(d, int)
            or a < 0
            or a > 0xFF
            or b < 0
            or b > 0xFF
            or c < 0
            or c > 0xFF
            or d < 0
            or d > 0xFF
        ):
            raise ValueError("cbk: invalid ABCD keys")
        if not isinstance(size, int) or size not in [16, 32, 64, 128]:
            raise ValueError("cbk: invalid size")
        s = Setting(6)
        s[0] = W_CBK
        s[1] = size & 0xFF
        s[2] = a
        s[3] = b
        s[4] = c
        s[5] = d
        return s


class Config(bytearray):
    def __init__(self):
        self._c = False
        self._t = False

    def json(self):
        i = 0
        n = 0
        e = list()
        while n >= 0 and n < len(self):
            n = self.next(i)
            if self[i] not in NAMES:
                raise ValueError(f"json: invalid setting id {self[i]}")
            o = None
            if self[i] == T_B64:
                pass
            elif self[i] >= W_HEX and self[i] <= W_B64:
                pass
            elif self[i] >= C_TCP and self[i] <= C_TLS_INSECURE:
                pass
            elif self[i] == S_HOST:
                o = self[
                    i + 3 : (int(self[i + 2]) | int(self[i + 1]) << 8) + i + 3
                ].decode("UTF-8")
            elif self[i] == S_SLEEP:
                o = _dur_to_str(
                    (
                        int(self[i + 8])
                        | int(self[i + 7]) << 8
                        | int(self[i + 6]) << 16
                        | int(self[i + 5]) << 24
                        | int(self[i + 4]) << 32
                        | int(self[i + 3]) << 40
                        | int(self[i + 2]) << 48
                        | int(self[i + 1]) << 56
                    )
                )
            elif self[i] == S_JITTER:
                o = int(self[i + 1])
            elif self[i] == C_IP or self[i] == C_TLS_EX or self[i] == T_B64S:
                o = int(self[i + 1])
            elif self[i] == C_WC2:
                z = i + 8
                v = (int(self[i + 2]) | int(self[i + 1]) << 8) + i + 8
                o = dict()
                if v > z:
                    o["url"] = self[z:v].decode("UTF-8")
                z = v
                v = (int(self[i + 4]) | int(self[i + 3]) << 8) + v
                if v > z:
                    o["host"] = self[z:v].decode("UTF-8")
                z = v
                v = (int(self[i + 6]) | int(self[i + 5]) << 8) + v
                if v > z:
                    o["agent"] = self[z:v].decode("UTF-8")
                if self[i + 7] > 0:
                    o["headers"] = dict()
                    j = 0
                    while v < n and z < n and j < n:
                        j = int(self[v]) + v + 2
                        z = v + 2
                        v = int(self[v + 1]) + j
                        if z == j:
                            raise ValueError("wc2: invalid header")
                        o["headers"][self[z:j].decode("UTF-8")] = self[j:v].decode(
                            "UTF-8"
                        )
                    del j
                del z
                del v
            elif self[i] == C_MTLS:
                a = (int(self[i + 3]) | int(self[i + 2]) << 8) + i + 8
                p = (int(self[i + 5]) | int(self[i + 4]) << 8) + a
                k = (int(self[i + 7]) | int(self[i + 6]) << 8) + p
                o = {"version": int(self[i + 1])}
                o["ca"] = b64encode(self[i + 8 : a]).decode("UTF-8")
                o["pem"] = b64encode(self[a:p]).decode("UTF-8")
                o["key"] = b64encode(self[p:k]).decode("UTF-8")
                del a
                del p
                del k
            elif self[i] == C_TLS_CA:
                a = (int(self[i + 3]) | int(self[i + 2]) << 8) + i + 4
                o = {"version": int(self[i + 1])}
                o["ca"] = b64encode(self[i + 4 : a]).decode("UTF-8")
                del a
            elif self[i] == C_TLS_CERT:
                p = (int(self[i + 3]) | int(self[i + 2]) << 8) + i + 6
                k = (int(self[i + 5]) | int(self[i + 4]) << 8) + p
                o = {"version": int(self[i + 1])}
                o["pem"] = b64encode(self[i + 6 : p]).decode("UTF-8")
                o["key"] = b64encode(self[p:k]).decode("UTF-8")
                del p
                del k
            elif self[i] == W_XOR:
                o = b64encode(
                    self[i + 3 : (int(self[i + 2]) | int(self[i + 1]) << 8) + i + 3]
                ).decode("UTF-8")
            elif self[i] == W_CBK:
                o = {
                    "size": int(self[i + 1]),
                    "A": int(self[i + 2]),
                    "B": int(self[i + 3]),
                    "C": int(self[i + 4]),
                    "D": int(self[i + 5]),
                }
            elif self[i] == W_AES:
                v = int(self[i + 1]) + i + 3
                z = int(self[i + 2]) + v
                if v == z or i + 3 == v:
                    raise ValueError("aes: invalid KEY/IV values")
                o = {
                    "key": b64encode(self[i + 3 : v]).decode("UTF-8"),
                    "iv": b64encode(self[v:z]).decode("UTF-8"),
                }
                del v
                del z
            elif self[i] == T_DNS:
                o = None
            r = {"type": NAMES[self[i]]}
            if o is not None:
                r["args"] = o
            del o
            e.append(r)
            i = n
        return e

    def add(self, s):
        if not isinstance(s, Setting):
            raise ValueError("add: cannot add a non-Settings object")
        if len(s) == 0 or s[0] == 0:
            raise ValueError("add: invalid Settings object")
        if (s[0] >= C_TCP and s[0] <= C_TLS_INSECURE) or (
            s[0] >= C_IP and s[0] <= C_TLS_CERT
        ):
            if self._c:
                raise ValueError("add: attempted to add multiple Connection hints")
            self._c = True
        if s[0] >= T_B64 and s[0] <= T_B64S:
            if self._t:
                raise ValueError("add: attempted to add multiple Transforms")
            self._t = True
        if s._single():
            return self.append(s[0])
        for i in s:
            self.append(i)

    def read(self, b):
        if not isinstance(b, bytes) and not isinstance(b, bytearray):
            raise ValueError("read: invalid raw type")
        self.extend(b)

    def next(self, i):
        if i > len(self):
            return -1
        if self[i] == T_B64:
            return i + 1
        if self[i] >= W_HEX and self[i] <= W_B64:
            return i + 1
        if self[i] >= C_TCP and self[i] <= C_TLS_INSECURE:
            return i + 1
        if (
            self[i] == C_IP
            or self[i] == T_B64S
            or self[i] == S_JITTER
            or self[i] == C_TLS_EX
        ):
            return i + 2
        if self[i] == W_CBK:
            return i + 6
        if self[i] == S_SLEEP:
            return i + 9
        if self[i] == C_WC2:
            n = (
                i
                + 8
                + (int(self[i + 2]) | int(self[i + 1]) << 8)
                + (int(self[i + 4]) | int(self[i + 3]) << 8)
                + (int(self[i + 6]) | int(self[i + 5]) << 8)
            )
            if self[i + 7] == 0:
                return n
            for _ in range(self[i + 7], 0, -1):
                n += int(self[n]) + int(self[n + 1]) + 2
            return n
        if self[i] == W_XOR or self[i] == S_HOST:
            return i + 3 + int(self[i + 2]) | int(self[i + 1]) << 8
        if self[i] == W_AES:
            return i + 3 + int(self[i + 1]) + int(self[i + 2])
        if self[i] == C_MTLS:
            return (
                i
                + 8
                + (int(self[i + 3]) | int(self[i + 2]) << 8)
                + (int(self[i + 5]) | int(self[i + 4]) << 8)
                + (int(self[i + 7]) | int(self[i + 6]) << 8)
            )
        if self[i] == C_TLS_CA:
            return i + 4 + int(self[i + 3]) | int(self[i + 2]) << 8
        if self[i] == C_TLS_CERT:
            return (
                i
                + 6
                + (int(self[i + 3]) | int(self[i + 2]) << 8)
                + (int(self[i + 5]) | int(self[i + 4]) << 8)
            )
        return -1

    def parse(self, j):
        v = loads(j)
        if not isinstance(v, list):
            raise ValueError("parse: invalid JSON value")
        if len(v) == 0:
            return
        for x in v:
            if not isinstance(x, dict) or len(x) == 0:
                raise ValueError("parse: invalid JSON value")
            if "type" not in x or x["type"].lower() not in NAMES_TO_ID:
                raise ValueError("parse: invalid JSON value")
            m = NAMES_TO_ID[x["type"].lower()]
            if m == T_B64:
                self.add(Cfg.single(m))
                continue
            if m >= W_HEX and m <= W_B64:
                self.add(Cfg.single(m))
                continue
            if m >= C_TCP and m <= C_TLS_INSECURE:
                self.add(Cfg.single(m))
                continue
            if "args" not in x:
                raise ValueError("parse: invalid JSON payload")
            p = x["args"]
            if m == S_HOST:
                if not _nes(p):
                    raise ValueError("host: invalid JSON value")
                self.add(Cfg.host(p))
            elif m == S_SLEEP:
                if not _nes(p):
                    raise ValueError("sleep: invalid JSON value")
                self.add(Cfg.sleep(p))
            elif m == S_JITTER:
                if not isinstance(p, int) and p > 0:
                    raise ValueError("jitter: invalid JSON value")
                self.add(Cfg.jitter(p))
            elif m == C_IP:
                if not isinstance(p, int) and p > 0:
                    raise ValueError("ip: invalid JSON value")
                self.add(Cfg.connect_ip(p))
            elif m == C_WC2:
                if not isinstance(p, dict):
                    raise ValueError("wc2: invalid JSON value")
                u = p.get("url")
                h = p.get("host")
                a = p.get("agent")
                j = p.get("headers")
                if j is not None and not isinstance(j, dict):
                    raise ValueError("wc2: invalid JSON header value")
                self.add(Cfg.connect_wc2(u, h, a, j))
                del u
                del h
                del a
                del j
            elif m == C_TLS_EX:
                if not isinstance(p, int) and p > 0:
                    raise ValueError("tls-ex: invalid JSON value")
                self.add(Cfg.connect_tls_ex(p))
            elif m == C_MTLS:
                if not isinstance(p, dict):
                    raise ValueError("mtls: invalid JSON value")
                a = p.get("ca")
                y = p.get("pem")
                k = p.get("key")
                n = p.get("version", 0)
                if not _nes(y) or not _nes(k):
                    raise ValueError("mtls: invalid JSON PEM/KEY values")
                if n is not None and not isinstance(n, int):
                    raise ValueError("mtls: invalid JSON version value")
                self.add(
                    Cfg.connect_mtls(
                        n,
                        b64decode(a, validate=True),
                        b64decode(y, validate=True),
                        b64decode(k, validate=True),
                    )
                )
                del a
                del y
                del k
                del n
            elif m == C_TLS_CA:
                if not isinstance(p, dict):
                    raise ValueError("tls-ca: invalid JSON value")
                a = p.get("ca")
                n = p.get("version", 0)
                if n is not None and not isinstance(n, int):
                    raise ValueError("tls-ca: invalid JSON version value")
                self.add(Cfg.connect_tls_ca(n, b64decode(a, validate=True)))
                del a
                del n
            elif m == C_TLS_CERT:
                if not isinstance(p, dict):
                    raise ValueError("tls-cert: invalid JSON value")
                y = p.get("pem")
                k = p.get("key")
                n = p.get("version", 0)
                if not _nes(y) or not _nes(k):
                    raise ValueError("tls-cert: invalid JSON PEM/KEY values")
                if n is not None and not isinstance(n, int):
                    raise ValueError("tls-cert: invalid JSON version value")
                self.add(
                    Cfg.connect_tls_certs(
                        n,
                        b64decode(y, validate=True),
                        b64decode(k, validate=True),
                    )
                )
                del y
                del k
                del n
            elif m == W_XOR:
                if not _nes(p):
                    raise ValueError("xor: invalid JSON value")
                self.add(Cfg.wrap_xor(b64decode(p, validate=True)))
            elif m == W_AES:
                if not isinstance(p, dict):
                    raise ValueError("aes: invalid JSON value")
                y = p.get("iv")
                k = p.get("key")
                if not _nes(y) or not _nes(k):
                    raise ValueError("aes: invalid JSON KEY/IV values")
                self.add(
                    Cfg.wrap_aes(
                        b64decode(k, validate=True),
                        b64decode(y, validate=True),
                    )
                )
                del y
                del k
            elif m == W_CBK:
                if not isinstance(p, dict):
                    raise ValueError("aes: invalid JSON value")
                A = p.get("A")
                B = p.get("B")
                C = p.get("C")
                D = p.get("D")
                z = p.get("size", 128)
                if not isinstance(A, int):
                    raise ValueError("cbk: invalid JSON A value")
                if not isinstance(B, int):
                    raise ValueError("cbk: invalid JSON B value")
                if not isinstance(C, int):
                    raise ValueError("cbk: invalid JSON C value")
                if not isinstance(D, int):
                    raise ValueError("cbk: invalid JSON D value")
                self.add(Cfg.wrap_cbk(a=A, b=B, c=C, d=D, size=z))
                del z
                del A
                del B
                del C
                del D
            elif m == T_DNS:
                pass
            elif m == T_B64S:
                if not isinstance(p, int) and p > 0:
                    raise ValueError("b64s: invalid JSON value")
                self.add(Cfg.transform_b64_shift(p))
            del p
            del m
        del v


class Setting(bytearray):
    def __str__(self):
        if len(self) == 0 or self[0] == 0:
            return "invalid"
        if self[0] not in NAMES:
            return "invalid"
        return NAMES[self[0]] + self.decode("UTF-8", "replace")

    def _single(self):
        if len(self) == 0 or self[0] == 0:
            return True
        if self[0] == T_B64:
            return True
        if self[0] >= W_HEX and self[0] <= W_B64:
            return True
        if self[0] >= C_TCP and self[0] <= C_TLS_INSECURE:
            return True
        return False


class Builder(ArgumentParser):
    def __init__(self):
        ArgumentParser.__init__(self)
        self.add_argument(
            "-f", "--in", type=str, dest="input", default=None, metavar="file"
        )
        self.add_argument(
            "-o", "--out", type=str, dest="output", default=None, metavar="file"
        )
        self.add_argument("-j", "--json", dest="json", action="store_true")
        self.add_argument("-p", "--print", dest="print", action="store_true")

        # System Args
        self.add_argument(
            "--host", type=str, dest="host", default=None, metavar="hostname"
        )
        self.add_argument(
            "--sleep", type=str, dest="sleep", default=None, metavar="seconds"
        )
        self.add_argument(
            "--jitter", type=int, dest="jitter", default=None, metavar="jitter"
        )

        cons = self.add_mutually_exclusive_group(required=False)

        # Connector Args
        cons.add_argument("--tcp", dest="tcp", action="store_true")
        cons.add_argument("--tls", dest="tls", action="store_true")
        cons.add_argument("--udp", dest="udp", action="store_true")
        cons.add_argument("--icmp", dest="icmp", action="store_true")
        cons.add_argument("--pipe", dest="pipe", action="store_true")
        cons.add_argument(
            "-K", "--tls-insecure", dest="tls_insecure", action="store_true"
        )

        # Custom Connector Args
        cons.add_argument("--ip", type=int, dest="ip", default=None, metavar="protocol")
        self.add_argument(
            "--wc2-url", type=str, dest="wc2_url", default=None, metavar="url"
        )
        self.add_argument(
            "--wc2-host", type=str, dest="wc2_host", default=None, metavar="host"
        )
        self.add_argument(
            "--wc2-agent", type=str, dest="wc2_agent", default=None, metavar="agent"
        )
        self.add_argument(
            "-H",
            "--wc2_header",
            nargs="*",
            type=str,
            dest="wc2_headers",
            action="append",
            default=None,
            metavar="key=value",
        )
        self.add_argument("--mtls", dest="mtls", action="store_true")
        self.add_argument(
            "--tls-ver", type=int, dest="tls_ver", default=None, metavar="version"
        )
        self.add_argument(
            "--tls-ca", type=str, dest="tls_ca", default=None, metavar="pem"
        )
        self.add_argument(
            "--tls-pem", type=str, dest="tls_pem", default=None, metavar="pem"
        )
        self.add_argument(
            "--tls-key", type=str, dest="tls_key", default=None, metavar="pem"
        )

        # Wrapper Args
        self.add_argument("--hex", dest="hex", action="store_true")
        self.add_argument("--zlib", dest="zlib", action="store_true")
        self.add_argument("--gzip", dest="gzip", action="store_true")
        self.add_argument("--b64", dest="b64", action="store_true")

        # Custom Wrapper Args
        self.add_argument(
            "--xor",
            nargs="?",
            type=str,
            dest="xor",
            action="append",
            default=None,
            metavar="key",
        )
        self.add_argument(
            "--cbk",
            nargs="?",
            type=str,
            dest="cbk",
            action="append",
            default=None,
            metavar="key",
        )
        self.add_argument(
            "--aes",
            nargs="?",
            type=str,
            dest="aes",
            action="append",
            default=None,
            metavar="key",
        )
        self.add_argument(
            "--aes-iv",
            nargs="?",
            type=str,
            dest="aes_iv",
            action="append",
            default=None,
            metavar="iv",
        )

        # Transform Args
        self.add_argument(
            "--b64t",
            nargs="?",
            type=int,
            dest="b64t",
            action="append",
            default=None,
            metavar="shift",
        )
        self.add_argument(
            "-D",
            "--dns",
            nargs="*",
            type=str,
            dest="dns",
            action="append",
            default=None,
            metavar="domain",
        )

    def run(self):
        if len(argv) <= 1:
            self.print_help()
        a = self.parse_args()
        if a.input:
            c = Builder._parse_input(a.input)
        else:
            c = Config()
        Builder._run_with(a, c)
        f = stdout
        if a.output and a.output != "-":
            if not a.json and not a.print:
                f = open(a.output, "wb")
            else:
                f = open(a.output, "w")
        if a.print or a.json:
            print(
                dumps(c.json(), sort_keys=False, indent=(4 if a.print else None)),
                file=f,
            )
        else:
            if f == stdout and not f.isatty():
                f.buffer.write(c)
            elif f.mode == "wb":
                f.write(c)
            else:
                f.write(b64encode(c).decode("UTF-8"))
        f.close()
        del a
        del f

    @staticmethod
    def _parse_pos():
        a = False
        w = list()
        d = dict()
        for i in range(0, len(argv)):
            if len(argv[i]) < 3:
                continue
            if argv[i][0] != "-":
                continue
            if not a and argv[i].lower() == "--aes-iv":
                v = "aes"
                a = True
            else:
                v = argv[i].lower()[2:]
            if v not in WRAPPERS:
                continue
            if v in d:
                raise ValueError(f'duplicate argument "--{v}" found')
            w.append(v)
            d[v] = len(w) - 1
        e = [None] * len(w)
        del w
        del a
        return d, e

    @staticmethod
    def _parse_input(v):
        if v.strip() == "-" and not stdin.isatty():
            if hasattr(stdin, "buffer"):
                b = stdin.buffer.read()
            else:
                b = stdin.read()
            stdin.close()
        else:
            with open(v, "rb") as f:
                b = f.read()
        c = Config()
        if len(b) == 0:
            raise ValueError("input: empty input data")
        try:
            c.parse(b.decode("UTF-8"))
            return c
        except (ValueError, JSONDecodeError):
            pass
        try:
            # This will fail if we're not base64.
            b.decode("UTF-8")
            c.read(b64decode(b, validate=True))
            return c
        except (ValueError, UnicodeDecodeError):
            pass
        c.read(b)
        del b
        return c

    @staticmethod
    def _run_with(a, c):
        p, w = Builder._parse_pos()
        if a.host:
            c.add(Cfg.host(a.host))
        if a.sleep:
            c.add(Cfg.sleep(a.sleep))
        if isinstance(a.jitter, int):
            c.add(Cfg.jitter(a.sleep))
        if a.tcp:
            c.add(Cfg.single(C_TCP))
        if a.tls:
            c.add(Cfg.single(C_TLS))
        if a.udp:
            c.add(Cfg.single(C_UDP))
        if a.icmp:
            c.add(Cfg.single(C_ICMP))
        if a.pipe:
            c.add(Cfg.single(C_PIPE))
        if a.tls_insecure:
            c.add(Cfg.single(C_TLS_INSECURE))
        if isinstance(a.ip, int):
            c.add(Cfg.connect_ip(a.ip))
        if a.wc2_url or a.wc2_host or a.wc2_agent or a.wc2_headers:
            c.add(
                Cfg.connect_wc2(
                    a.wc2_url,
                    a.wc2_host,
                    a.wc2_agent,
                    Builder._parse_headers(a.wc2_headers),
                )
            )
        if a.tls_ca or a.tls_pem or a.tls_key:
            c.add(Builder._parse_tls(a.tls_ca, a.tls_pem, a.tls_key, a.mtls, a.tls_ver))
        elif a.mtls:
            raise ValueError("mtls: missing CA, PEM and KEY values")
        elif isinstance(a.tls_ver, int):
            c.add(Cfg.connect_tls_ex(a.tls_ver))
        if a.hex:
            w[p["hex"]] = Cfg.single(W_HEX)
        if a.zlib:
            w[p["zlib"]] = Cfg.single(W_ZLIB)
        if a.gzip:
            w[p["gzip"]] = Cfg.single(W_GZIP)
        if a.b64:
            w[p["b64"]] = Cfg.single(W_B64)
        if a.xor:
            w[p["xor"]] = Cfg.wrap_xor(a.xor[0])
        if a.cbk:
            w[p["cbk"]] = Cfg.wrap_cbk(key=a.cbk[0])
        if a.aes or a.aes_iv:
            w[p["aes"]] = Cfg.wrap_aes(a.aes, a.aes_iv)
        for i in w:
            c.add(i)
        del w

    @staticmethod
    def _parse_headers(v):
        if not isinstance(v, list) or len(v) == 0:
            return None
        d = dict()
        for e in v:
            Builder._parse_header(d, e, False)
        if len(d) == 0:
            return None
        return d

    @staticmethod
    def _parse_header(d, e, r):
        if isinstance(e, str):
            if len(e) == 0 or "=" not in e:
                raise ValueError("wc2: invalid header")
            p = e.find("=")
            if p == 0 or p == len(e) - 1:
                raise ValueError("wc2: empty header")
            d[e[:p].strip()] = e[p + 1 :].strip()
            return
        if isinstance(e, list) and len(e) > 0:
            if r:
                raise ValueError("wc2: too many nested lists")
            for v in e:
                Builder._parse_header(d, v, True)
            return
        raise ValueError("wc2: Invalid header")

    def print_help(self, file=None):
        print(BUILDER_HELP, file=file)
        exit(2)

    @staticmethod
    def _parse_tls(ca, pem, key, mtls, ver):
        a = None
        p = None
        k = None
        if isinstance(ca, str) and len(ca) > 0:
            try:
                a = b64decode(ca, validate=True)
            except ValueError:
                with open(ca, "rb") as f:
                    a = f.read()
        if isinstance(pem, str) and len(pem) > 0:
            try:
                a = b64decode(pem, validate=True)
            except ValueError:
                with open(pem, "rb") as f:
                    p = f.read()
        if isinstance(key, str) and len(key) > 0:
            try:
                a = b64decode(key, validate=True)
            except ValueError:
                with open(key, "rb") as f:
                    k = f.read()
        if mtls and (p is None or k is None or a is None):
            raise ValueError("mtls: CA, PEM and KEY must be provided")
        if (p is not None and k is None) or (k is not None and p is None):
            raise ValueError("tls-cert: PEM and KEY must be provided")
        if not isinstance(ver, int):
            ver = 0
        if a is None and p is None and k is None:
            return Cfg.connect_tls_ex(ver)
        if a is not None and p is None and k is None:
            return Cfg.connect_tls_ca(ver, a)
        if a is None:
            return Cfg.connect_tls_certs(ver, p, k)
        return Cfg.connect_mtls(ver, a, p, k)


if __name__ == "__main__":
    b = Builder()
    try:
        b.run()
    except Exception as err:
        print(f"Error: {err}\n{format_exc(limit=3)}", file=stderr)
        exit(1)