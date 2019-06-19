// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ast

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

type scanner struct{}

func (sx scanner) scanTerm(s string, offset int) (*Term, int) {
	first, _ := utf8.DecodeRuneInString(s)
	if first == utf8.RuneError {
		panic(ParseError{offset, "invalid utf8 encoding"})
	}

	if first == '"' || first == '\'' {
		q, off := sx.scanQuote(s, offset)
		return &Term{Quote: q}, off
	}

	if q := sx.scanNumeric(s); q != "" {
		return &Term{Number: q}, offset + len(q)
	}

	if q := sx.scanID(s); q != "" {
		return &Term{Name: q}, offset + len(q)
	}

	return nil, offset
}

func (sx scanner) scanQuote(s string, offset int) (quoted string, off int) {
	var first rune
	skip := false
	result := []rune{}
	for idx, r := range s {
		switch {
		case idx == 0:
			first = r
		case !skip && r == '\\':
			skip = true
		case !skip && r == first:
			return string(result), idx + utf8.RuneLen(r)
		default:
			result = append(result, r)
			skip = false
		}
	}

	panic(ParseError{offset, "incomplete quote"})
}

func (sx scanner) scanNumeric(s string) string {
	last := 0
	for kk := range s {
		if _, err := strconv.ParseFloat(s[:kk], 32); kk > 0 && err != nil {
			return s[:last]
		}
		last = kk
	}
	if _, err := strconv.ParseFloat(s, 32); err != nil {
		return s[:last]
	}
	return s
}

func (sx scanner) scanID(s string) string {
	for idx, r := range s {
		if !unicode.IsLetter(r) && (idx == 0 || !unicode.IsDigit(r)) {
			return s[:idx]
		}
	}
	return s
}
