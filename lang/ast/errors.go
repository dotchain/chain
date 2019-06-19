// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ast

import "strconv"

// ParseError provides the offset of any errors.
type ParseError struct {
	Offset  int
	Message string
}

// Error returns a string version of the error
func (p ParseError) Error() string {
	return strconv.Itoa(p.Offset) + ": " + p.Message
}
