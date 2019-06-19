// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package ast has the parser and render for the chain language
package ast

// Parse parses the input string and returns the AST for it.
//
// If there is a parse error, the returned error will be of type ParseError.
func Parse(s string) (n *Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(ParseError); !ok {
				panic(r)
			}
			err = r.(ParseError)
		}
	}()

	result, p := (&parser{s: s}).parse()
	if p.s != "" {
		panic(ParseError{p.offset, "unexpected character"})
	}
	return normalize(result), nil
}

// Term represents a number or name or quoted string
type Term struct {
	Number string
	Name   string
	Quote  string
}

// Node represents a node in the tree
//
// It can be a terminal node (in which case Term will be non-nil) or
// an op node in which Children may be present.
//
// If both Term and Op are nil, it is a "missing node" (or a
// placeholder node) -- Start, End are still expected to be good.
type Node struct {
	Start, End int
	*Term
	Op       string
	Children []*Node
}

// NewNode creates a new node.
//
// The start and end can be -1 to indicate unknown. They are
// normalized from the provided children.  Either Term or Op must be
// specified and there can be any number of children
func NewNode(start, end int, term *Term, op string, children ...*Node) *Node {
	n := &Node{Start: start, End: end, Term: term, Op: op, Children: children}
	var min, max *Node
	for _, node := range children {
		if node != nil {
			if min == nil || min.Start > node.Start {
				min = node
			}
			if max == nil || max.End < node.End {
				max = node
			}
		}
	}

	if start < 0 && min != nil {
		n.Start = min.Start
	}
	if end < 0 && max != nil {
		n.End = max.End
	}
	return n
}

// operators and their priorities
var priority = map[string]int{
	",": 1,
	"=": 2,
	"|": 3,
	"&": 4,

	"<":  5,
	">":  5,
	"<=": 5,
	">=": 5,

	"==": 6,
	"!=": 6,

	"+": 7,
	"-": 7,

	"*": 8,
	"/": 8,

	"(": 9,
	")": 10,

	".": 11,
}
