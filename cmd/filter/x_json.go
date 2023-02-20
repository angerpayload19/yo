//go:build !nojson
// +build !nojson

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

package filter

import "encoding/json"

// MarshalJSON will attempt to convert the data in this Filter into the returned
// JSON byte array.
func (f Filter) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{"fallback": f.Fallback}
	if f.PID != 0 {
		m["pid"] = f.PID
	}
	if f.Session > Empty {
		m["session"] = f.Session
	}
	if f.Elevated > Empty {
		m["elevated"] = f.Elevated
	}
	if len(f.Exclude) > 0 {
		m["exclude"] = f.Elevated
	}
	if len(f.Include) > 0 {
		m["include"] = f.Include
	}
	return json.Marshal(m)
}

// UnmarshalJSON will attempt to parse the supplied JSON into this Filter.
func (f *Filter) UnmarshalJSON(b []byte) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	if len(m) == 0 {
		return nil
	}
	if v, ok := m["pid"]; ok {
		if err := json.Unmarshal(v, &f.PID); err != nil {
			return err
		}
	}
	if v, ok := m["session"]; ok {
		if err := json.Unmarshal(v, &f.Session); err != nil {
			return err
		}
	}
	if v, ok := m["elevated"]; ok {
		if err := json.Unmarshal(v, &f.Elevated); err != nil {
			return err
		}
	}
	if v, ok := m["exclude"]; ok {
		if err := json.Unmarshal(v, &f.Exclude); err != nil {
			return err
		}
	}
	if v, ok := m["include"]; ok {
		if err := json.Unmarshal(v, &f.Include); err != nil {
			return err
		}
	}
	if v, ok := m["fallback"]; ok {
		if err := json.Unmarshal(v, &f.Fallback); err != nil {
			return err
		}
	}
	return nil
}
func (b boolean) MarshalJSON() ([]byte, error) {
	switch b {
	case True:
		return []byte(`"true"`), nil
	case False:
		return []byte(`"false"`), nil
	default:
	}
	return []byte(`""`), nil
}
func (b *boolean) UnmarshalJSON(d []byte) error {
	if len(d) == 0 {
		*b = Empty
		return nil
	}
	if d[0] == '"' && len(d) >= 1 {
		switch d[1] {
		case '1', 'T', 't':
			*b = True
			return nil
		case '0', 'F', 'f':
			*b = False
			return nil
		}
		*b = Empty
		return nil
	}
	switch d[0] {
	case '1', 'T', 't':
		*b = True
		return nil
	case '0', 'F', 'f':
		*b = False
		return nil
	}
	*b = Empty
	return nil
}
