package abnf

import (
	"bytes"
	"io"
)

// Rule is rule function of ABNF
type Rule func(s *scanner) []rune

// reading scanner
type scanner struct {
	// main read scanner
	b []rune
	// scanner expand function
	r io.RuneReader

	// marker stack
	rp []int
	wp *Tree

	// pre-reading pointer
	pp int
}

// ParseString parse string value by Rule
func ParseString(b string, f Rule) *Tree {
	return ParseReader(bytes.NewReader([]byte(b)), len(b), len(b), f)
}

// ParseReader parse reader value by Rule
func ParseReader(r io.RuneReader, blen, slen int, f Rule) *Tree {
	s := scanner{}
	s.b = make([]rune, 0, blen)
	s.r = r
	s.rp = make([]int, 0, slen)
	s.wp = &Tree{-1, nil, nil, nil}
	s.pp = 0

	if v := f(&s); v != nil {
		s.wp.V = v
		return s.wp
	}
	return nil
}

// set marker
func (s *scanner) mark() {
	// push marker to read stack
	s.rp = append(s.rp, s.pp)
}

// commit
func (s *scanner) commit() []rune {
	// delete marker from read stack, and get buffer value
	ret := s.b[s.rp[len(s.rp)-1]:s.pp]
	s.rp = s.rp[0 : len(s.rp)-1]

	return ret
}

// rollback
func (s *scanner) rollback() []rune {
	// pop marker from red stack
	s.pp = s.rp[len(s.rp)-1]
	s.rp = s.rp[0 : len(s.rp)-1]

	return nil
}

// get one charactor
func (s *scanner) next() []rune {
	// end of current scanner
	if len(s.b) <= s.pp {
		// read additional data
		ch, _, e := s.r.ReadRune()
		if e != nil {
			// println("next()=nil")
			return nil
		}
		// append to scanner
		s.b = append(s.b, ch)
	}
	return s.b[s.pp : s.pp+1]
}
