package abnf

import "strings"

// Terminal Values
// match -> []rune of match value
// different or EOF -> nil

// V:   Match single character
func V(c rune) Rule {
	return func(s *Scanner) []rune {
		if ch := s.next(); ch != nil && ch[0] == c {
            s.pp++
            return ch
        }
		return nil
	}
}

// VL:  Alternatives of multiple characters
func VL(clist ...rune) Rule {
	return func(s *Scanner) []rune {
		ch := s.next()
		// end of data
		if ch == nil {
			return nil
		}
		// verify
		for _, c := range clist {
			if ch[0] == c {
				s.pp++
				return ch
			}
		}
		return nil
	}
}

// VI:  Match single string (case insensitive)
func VI(str string) Rule {
    str = strings.ToLower(str)
    fs := make([]Rule, len(str))
	for i, ch1 := range str {
        ch2 := ch1
        if ch2 >= 0x61 && ch2 <= 0x7a {
            ch2 -= 0x20
        }
		fs[i] = VL(ch1, ch2)
	}
    return C(fs...)
}

// VIL: Alternatives of multiple strings (case insensitive)
func VIL(slist ...string) Rule {
    fs := make([]Rule, len(slist))
	for i, s := range slist {
        fs[i] = VI(s)
	}
	return A(fs...)
}

// VS:  Match single string (case sensitive)
func VS(str string) Rule {
    fs := make([]Rule, len(str))
	for i, ch := range str {
		fs[i] = V(ch)
	}
    return C(fs...)
}

// VSL: Alternatives of multiple strings (case sensitive)
func VSL(slist ...string) Rule {
    fs := make([]Rule, len(slist))
	for i, s := range slist {
        fs[i] = VS(s)
	}
	return A(fs...)
}

// Concatenation:  Rule1 Rule2
// match -> []rune of match value
// different or EOF -> nil
func C(fs ...Rule) Rule {
	return func(s *Scanner) []rune {
		s.mark()
		for _, f := range fs {
			if f(s) == nil {
				return s.rollback()
			}
		}
		return s.commit()
	}
}

// Alternatives:  Rule1 / Rule2
// Incremental Alternatives: Rule1 =/ Rule2
// different or EOF -> nil
func A(fs ...Rule) Rule {
	return func(s *Scanner) []rune {
		for _, f := range fs {
			if b := f(s); b != nil {
				return b
			}
		}
		return nil
	}
}

// Value Range Alternatives:  %c##-##
// match -> []rune of match value
// different or EOF -> nil
func VR(head rune, tail rune) Rule {
	return func(s *Scanner) []rune {
		if ch := s.next(); ch != nil && ch[0] >= head && ch[0] <= tail {
			s.pp++
			return ch
		}
		return nil
	}
}

// Specific Repetition:  nRule
// match -> []rune of match value
// different or EOF -> nil
func RN(n int, fs Rule) Rule {
	return func(s *Scanner) []rune {
		return repet(s, n, n, fs)
	}
}

// Variable Repetition:  *Rule
// match -> []rune of match value
// different or EOF -> nil
// If max=-1 then max=infinity.
func RV(min int, max int, fs Rule) Rule {
	return func(s *Scanner) []rune {
		return repet(s, min, max, fs)
	}
}

// Default repeat count of Variable Repetition
// From zero to infinity
func R0(fs Rule) Rule {
	return func(s *Scanner) []rune {
		return repet(s, 0, -1, fs)
	}
}

// More than one repeat count of Variable Repetition
// From one to infinity
func R1(fs Rule) Rule {
	return func(s *Scanner) []rune {
		return repet(s, 1, -1, fs)
	}
}

// subroutin of Repeat
func repet(s *Scanner, min int, max int, f Rule) []rune {
	s.mark()
	i := 0
	for ; max < 0 || i < max; i++ {
		if f(s) == nil {
			break
		}
	}
	// check minimum count
	if i < min {
		return s.rollback()
	}
	return s.commit()
}

// Optional Sequence:  [RULE]
// match -> []rune of match value
// different or EOF -> empty []rune
func O(fs Rule) Rule {
	return func(s *Scanner) []rune {
		b := fs(s)
		if b == nil {
			b = make([]rune, 0)
		}
		return b
	}
}



/*
   Additional operator
*/

// Not match
// match -> nil
// different or EOF -> empty []rune
func N(fs Rule) Rule {
	return func(s *Scanner) []rune {
        r := make([]rune, 0)
		s.mark()
		if fs(s) != nil {
            r = nil
        }
		s.rollback()
		return r
	}
}

