package abnf

import (
	"bytes"
	"io"
)

type Rule func(s *Scanner) []rune

// reading Scanner
type Scanner struct {
	// main read Scanner
	b []rune
	// Scanner expand function
	r io.RuneReader

	// marker stack
	rp []int
	wp *Tree

	// pre-reading pointer
	pp int
}

// Create new reading Scanner
func NewScanner(r io.RuneReader, t *Tree, blen, slen int) *Scanner {
	s := Scanner{}
	s.b = make([]rune, 0, blen)
	s.r = r
	s.rp = make([]int, 0, slen)
	s.wp = t
	s.pp = 0

	return &s
}

func StringScanner(b string) (s *Scanner, t *Tree) {
	t = &Tree{-1, []rune(b), nil, nil}
	s = NewScanner(bytes.NewReader([]byte(b)), t, len(b), len(b))
	return
}

// Pars string value by Rule
func ParseString(b string, f Rule) *Tree {
	t = &Tree{-1, []rune(b), nil, nil}
	s = NewScanner(bytes.NewReader([]byte(b)), t, len(b), len(b))
	if f(s) != nil {
		return t
	}
	return nil
}

// set marker
func (s *Scanner) mark() {
	// push marker to read stack
    s.rp = append(s.rp, s.pp)
}

// commit
func (s *Scanner) commit() []rune {
	// delete marker from read stack, and get buffer value
	ret := s.b[s.rp[len(s.rp) - 1]: s.pp]
    s.rp = s.rp[0: len(s.rp) - 1]
    
	return ret
}

// rollback
func (s *Scanner) rollback() []rune {
	// pop marker from red stack
	s.pp = s.rp[len(s.rp) - 1]
    s.rp = s.rp[0: len(s.rp) - 1]
    
	return nil
}

// get one charactor
func (s *Scanner) next() []rune {
	// end of current Scanner
	if len(s.b) <= s.pp {
		// read additional data
		ch, _, e := s.r.ReadRune()
		if e != nil {
			// println("next()=nil")
			return nil
		}
		// append to Scanner
		s.b = append(s.b, ch)
	}
	return s.b[s.pp: s.pp + 1]
}




// create function for add new tree node
func K(f Rule, k int) Rule {
	return func(s *Scanner) []rune {
        p := s.wp
        s.wp = &Tree{k, nil, nil, nil}
        
        v := f(s)
		if v != nil {
            s.wp.V = v
            p.add(s.wp)
        }
        s.wp = p
    	return v
	}
}

