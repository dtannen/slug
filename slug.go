// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package slug transforms strings into a normalized form that is safe for use
// in URLs. It is also useful for sanitizing user-submitted data because it
// removes or replaces everything except for Latin alphanumeric characters and
// removes diacritical marks.
package slug

import (
	"code.google.com/p/go.text/unicode/norm"
	"unicode"
)

var (
	// Non alphanumeric unicode characters will be replaced with this byte.
	Replacement = '_'

	alphanum = &unicode.RangeTable{
		R16: []unicode.Range16{
			{0x0030, 0x0039, 1}, // 0-9
			{0x0041, 0x005A, 1}, // A-Z
			{0x0061, 0x007A, 1}, // a-z
		},
	}
	nop = []*unicode.RangeTable{
		unicode.Mark,
		unicode.Sk, // Symbol - modifier
		unicode.Lm, // Letter - modifier
		unicode.Cc, // Other - control
		unicode.Cf, // Other - format
	}
)

// Slug replaces each run of characters which are not ASCII letters or numbers
// with the Replacement character, except for leading or trailing runs. Letters
// will be stripped of diacritical marks and lowercased. Letter or number
// codepoints that do not have combining marks or a lower-cased variant will be
// passed through unaltered.
func Clean(s string) string {
	buf := make([]rune, 0, len(s))
	replacement := false

	for _, r := range norm.NFKD.String(s) {
		switch {
		case unicode.In(r, alphanum):
			buf = append(buf, unicode.ToLower(r))
			replacement = true
		case unicode.IsOneOf(nop, r):
			// skip
		case replacement:
			buf = append(buf, Replacement)
			replacement = false
		}
	}

	// Strip trailing Replacement byte.
	if i := len(buf) - 1; i >= 0 && buf[i] == Replacement {
		buf = buf[:i]
	}

	return string(buf)
}
