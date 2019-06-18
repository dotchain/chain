// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ast_test

import (
	"strings"
	"testing"

	"github.com/dotchain/chain/lang/ast"
	"github.com/tvastar/test"
)

func TestNormalize(t *testing.T) {
	test.File(t.Error, "norm.txt", "normed.json", func(input string) interface{} {
		result := []interface{}{}
		for _, expr := range strings.Split(input, "\n\n") {
			parsed, err := ast.Parse(expr)
			result = append(result, struct {
				Expr       string
				Normalized *ast.Node
				Error      error
			}{expr, parsed, err})
		}
		return result
	})

}
