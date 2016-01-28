package abnf

import (
	"bytes"
	"io"
	// "unicode/utf8"
)

// length of mark pointer stack
const (
	STACK_LENGTH  = 65536
	BUFFER_LENGTH = 65536
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
func NewScanner(r io.RuneReader, t *Tree) *Scanner {
	s := Scanner{}
	s.b = make([]rune, 0, BUFFER_LENGTH)
	s.r = r
	s.rp = make([]int, 1, STACK_LENGTH)
	s.rp[0] = 0
	s.wp = t
	s.pp = 0

	return &s
}

func StringScanner(b string, k int) (s *Scanner, t *Tree) {
	t = &Tree{k, []rune(b), nil, nil}
	s = NewScanner(bytes.NewReader([]byte(b)), t)
	return
}

// add mark pointer
func (s *Scanner) mark() {
	// push marker to read stack
    s.rp = append(s.rp, s.pp)
}

// commit
func (s *Scanner) commit() []rune {
	// delete marker from read stack, and get buffer value
	// println("commit()=" + string(s.b[s.rp[s.st]: s.pp]))
	ret := s.b[s.rp[len(s.rp) - 1]: s.pp]
    s.rp = s.rp[0: len(s.rp) - 1]
    
	return ret
}

// rollback
func (s *Scanner) rollback() []rune {
	// pop marker from red stack
	// println("rollback()=" + string(s.b[s.rp[s.st]: s.pp]))
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
	// println("next()=" + string(s.b[s.pp: s.pp + 1]))
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

