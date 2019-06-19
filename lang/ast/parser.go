// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ast

import (
	"strings"
	"unicode"
)

// parser implements a two-stack "shunting-yard" parse algorithm
type parser struct {
	s      string
	offset int

	ops   []*Node
	nodes []*Node
}

// parse returns an AST and any remaining unparseable string
func (p *parser) parse() (*Node, *parser) {
	var lastWasTerm bool
	p = p.trimSpaceLeft()
main:
	for p.s != "" {
		n, next := p.scan()
		switch {
		case n == nil:
			p = next
			break main
		case n.Op == ")":
			break main
		case n.Op == "(":
			n, next = next.parseGroup(n)
			if lastWasTerm {
				next = next.reduceNodes(priority[n.Op])
				l := len(next.nodes) - 1
				next.nodes[l] = NewNode(-1, -1, nil, "()", next.nodes[l], n)
				break
			}
			fallthrough
		case n.Term != nil:
			p.failIf(lastWasTerm, p.offset, "unexpected term")
			next.nodes = append(next.nodes, n)
			lastWasTerm = true
		case n.Op != "":
			if !lastWasTerm {
				fixable := strings.Contains("-,=", n.Op)
				p.failIf(!fixable, p.offset, "unexpected operator")
				next.nodes = append(next.nodes, nil)
			}
			next = next.reduceNodes(priority[n.Op])
			lastWasTerm = false
			next.ops = append(next.ops, n)
		}

		p = next.trimSpaceLeft()
	}

	if len(p.nodes) == 0 && len(p.ops) == 0 {
		return nil, p
	}

	if !lastWasTerm {
		p.nodes = append(p.nodes, nil)
	}

	p = p.reduceNodes(-1)
	return p.nodes[0], p
}

func (p *parser) parseGroup(begin *Node) (*Node, *parser) {
	group, next := (&parser{s: p.s, offset: p.offset}).parse()

	p.failIf(next.s == "", next.offset, "incomplete (")
	n, next2 := next.scan()

	p.failIf(n.Op != ")", next.offset, "unexpected char")

	*next = *p
	next.offset, next.s = next2.offset, next2.s
	return NewNode(begin.Start, next.offset, nil, begin.Op, group), next
}

func (p *parser) reduceNodes(pri int) *parser {
	result := *p
	p = &result

	t := len(p.nodes) - 1
	for l := len(p.ops) - 1; l >= 0 && priority[p.ops[l].Op] >= pri; l-- {
		opn := p.ops[l]
		p.nodes[t-1] = NewNode(opn.Start, opn.End, nil, opn.Op, p.nodes[t-1], p.nodes[t])
		p.ops, p.nodes = p.ops[:l], p.nodes[:t]
		t--
	}
	return p
}

// scan looks for an op or a term and also returns the updated parser state.
func (p *parser) scan() (*Node, *parser) {
	var op string

	if len(p.s) > 1 && priority[p.s[:2]] > 0 {
		op = p.s[:2]
	} else if priority[p.s[:1]] > 0 {
		op = p.s[:1]
	}

	if op != "" {
		n := &Node{Start: p.offset, End: p.offset + len(op), Op: op}
		return n, p.advance(len(op))
	}

	term, offset := (scanner{}).scanTerm(p.s, p.offset)
	if term != nil {
		n := &Node{Start: p.offset, End: offset, Term: term}
		return n, p.advance(offset - p.offset)
	}

	return nil, p
}

func (p *parser) trimSpaceLeft() *parser {
	idx := strings.IndexFunc(p.s, func(r rune) bool {
		return !unicode.IsSpace(r)
	})
	if idx < 0 {
		idx = len(p.s)
	}
	return p.advance(idx)
}

func (p *parser) advance(count int) *parser {
	next := *p
	next.s, next.offset = p.s[count:], p.offset+count
	return &next
}

func (p *parser) failIf(cond bool, offset int, message string) {
	if cond {
		panic(ParseError{offset, message})
	}
}
