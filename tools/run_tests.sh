#!/usr/bin/bash -e
# Copyright (C) 2020 - 2023 iDigitalFlame
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

build_tags=(
              bugs implant crypt stdrand nojson nosweep tiny small medium large nofrag regexp nopanic noservice ews noproxy nokeyset scripts
              bugs,implant
              bugs,implant,crypt
              bugs,implant,crypt,stdrand
              bugs,implant,crypt,stdrand,nojson
              bugs,implant,crypt,stdrand,nojson,nosweep
              bugs,implant,crypt,stdrand,nojson,nosweep,tiny
              bugs,implant,crypt,stdrand,nojson,nosweep,small
              bugs,implant,crypt,stdrand,nojson,nosweep,medium
              bugs,implant,crypt,stdrand,nojson,nosweep,nofrag
              bugs,implant,crypt,stdrand,nojson,nosweep,nofrag,regexp
              bugs,implant,crypt,stdrand,nojson,nosweep,nofrag,regexp,ews
              bugs,implant,crypt,stdrand,nojson,nosweep,nofrag,regexp,ews,noproxy
              implant,crypt,stdrand,nojson,nosweep,nofrag,regexp,ews,noproxy
              crypt,stdrand,nojson,nosweep,nofrag,regexp,ews,noproxy
              crypt,stdrand,nojson,nosweep,nofrag,regexp,noproxy
              stdrand,nojson,nosweep,nofrag,regexp,ews,noproxy
              nojson,nosweep,nofrag,regexp,ews,noproxy
              nosweep,nofrag,regexp,ews,noproxy
              nosweep,nofrag,ews,noproxy
              nosweep,nofrag,noproxy
              nosweep,nofrag,ews
              nofrag,regexp,ews
              nofrag,regexp
              nofrag,ews
              nofrag
              funcmap
           )

if ! go mod tidy; then
    printf "\x1b[1m\x1b[31m[!] Tidying modules failed!\x1b[0m\n"
    exit 1
fi

if ! go mod verify; then
    printf "\x1b[1m\x1b[31m[!] Verifying modules failed!\x1b[0m\n"
    exit 1
fi

run_vet() {
    printf "\x1b[1m\x1b[36m[+] Vetting GOARCH=$1 GOOS=$2..\x1b[0m\n"
    env GOARCH=$1 GOOS=$2 go vet ./c2 2>&1 | grep -vE '^# github.com'
    env GOARCH=$1 GOOS=$2 go vet ./cmd 2>&1 | grep -vE '^# github.com|: possible misuse of unsafe.Pointer$'
    env GOARCH=$1 GOOS=$2 go vet ./com 2>&1 | grep -vE '^# github.com'
    env GOARCH=$1 GOOS=$2 go vet ./data 2>&1 | grep -vE '^# github.com'
    env GOARCH=$1 GOOS=$2 go vet ./device 2>&1 | grep -vE '^# github.com|: possible misuse of unsafe.Pointer$|\.s:'
    env GOARCH=$1 GOOS=$2 go vet ./man 2>&1 | grep -vE '^# github.com'
    env GOARCH=$1 GOOS=$2 go vet ./util 2>&1 | grep -vE '^# github.com'
}
run_vet_all() {
    for entry in $(go tool dist list); do
        run_vet "$(echo $entry | cut -d '/' -f 2)" "$(echo $entry | cut -d '/' -f 1)"
    done
}
run_staticcheck() {
    printf "\x1b[1m\x1b[33m[+] Static Check GOARCH=$1 GOOS=$2..\x1b[0m\n"
    env GOARCH=$1 GOOS=$2 staticcheck -checks all -f text ./... | grep -vE '^-: # github.com/iDigitalFlame/xmt/examples$|^examples/|^unit_tests/|^c2/session.go:288|^c2/session.go:1155'
    for tags in ${build_tags[@]}; do
        printf "\x1b[1m\x1b[34m[+] Static Check GOARCH=$1 GOOS=$2 with tags \"${tags}\"..\x1b[0m\n"
        env GOARCH=$1 GOOS=$2 staticcheck -checks all -f text -tags $tags ./... | grep -vE '^-: # github.com/iDigitalFlame/xmt/examples$|^examples/|^unit_tests/|^c2/session.go:288|^c2/session.go:1155'
    done
}
run_staticcheck_all() {
    for entry in $(go tool dist list); do
        run_staticcheck "$(echo $entry | cut -d '/' -f 2)" "$(echo $entry | cut -d '/' -f 1)"
    done
}

run_vet_all
run_staticcheck_all
