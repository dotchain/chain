// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ast

// Normalize simplifies the AST for some typical cases (unary "+"
// operator for example) and flags some other cases as errors (missing
// term).
//
// Normalization results in new operators:
//
//  "[]" => array
//  "{}" => object
//  "()" => function call
//
// Full list of normalizations:
//    () => null
//    (()) => null,
//    (,,,) => list()
//    (,1,,) => list(1)
//    (=) => object()
//    (=,) => object()
//    (x=y,z=a) => object(list(x, z), list(y, a))
//    (x=y,z) => object(list(x, y), list(x, z))
//    (x=y,1) => object(list(x, y), list(1))
//    (x=,y=5) => error: incomplete object definition
//    (=x) => error: incomplete object definition
//    (fn)(x,,y,) => call(fn, x, y)
//    (fn)(x = 23, y) => call(fn, object(list(x, 23), list(y)))
func normalize(n *Node) *Node {
	if n == nil {
		return nil
	}

	result := *n

	if n.Op == "(" {
		result.Op = "[]"
		result.Children = nil
		switch {
		case isObject(n.Children[0]):
			result.Op = "{}"
			result.Children = normalizeObject(nil, n.Children[0])
		case n.Children[0] != nil && n.Children[0].Op == ",":
			result.Children = normalizeList(nil, n.Children[0])
		default:
			return normalize(n.Children[0])
		}

		return &result
	}

	if n.Op == "()" {
		c := []*Node{normalizeFunc(n.Children[0])}
		if len(n.Children) > 1 {
			c = normalizeList(c, n.Children[1])
		}
		result.Children = c
		return &result
	}

	if n.Op == "," || n.Op == "=" {
		panic(ParseError{n.Start, "unexpected " + n.Op})
	}

	return n
}

func normalizeFunc(n *Node) *Node {
	n = normalize(n)
	if n.Op == "[]" {
		panic(ParseError{n.Start, "cannot call a list"})
	}

	return n
}

func normalizeList(l []*Node, n *Node) []*Node {
	for _, child := range n.Children {
		if child != nil && child.Op == "," {
			l = normalizeList(l, child)
		} else if n := normalize(child); n != nil {
			l = append(l, n)
		}
	}
	return l
}

func collectArgs(l []*Node, n *Node) []*Node {
	if n == nil {
		return l
	}

	if n.Op != "," {
		l = append(l, n)
		return l
	}

	for _, child := range n.Children {
		l = collectArgs(l, child)
	}
	return l
}

func normalizeObject(l []*Node, n *Node) []*Node {
	for _, elt := range collectArgs(nil, n) {
		if elt == nil {
			continue
		}

		result := *elt
		result.Term = nil
		result.Op = "[]"

		if elt.Op != "=" {
			if elt = normalize(elt); elt != nil {
				result.Children = []*Node{elt}
				l = append(l, &result)
			}
			continue
		}

		c := []*Node{}
		for _, child := range elt.Children {
			child = normalize(child)
			if child != nil {
				c = append(c, child)
			}
		}
		if len(c) == 0 {
			continue
		}

		if len(c) == 1 {
			panic(ParseError{elt.Start, "incomplete object definition"})
		}

		result.Children = c
		l = append(l, &result)
	}

	return l
}

func isObject(n *Node) bool {
	if n != nil && n.Op == "=" {
		return true
	}

	if n == nil || n.Op != "," {
		return false
	}
	for _, child := range n.Children {
		if isObject(child) {
			return true
		}
	}
	return false
}
