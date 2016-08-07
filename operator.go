package abnf

import "strings"

// Terminal Values
// match -> []rune of match value
// different or EOF -> nil

// V match single character
func V(c rune) Rule {
	return func(s *scanner) []rune {
		if ch := s.next(); ch != nil && ch[0] == c {
			s.pp++
			return ch
		}
		return nil
	}
}

// VL match multiple characters
func VL(cl ...rune) Rule {
	return func(s *scanner) []rune {
		ch := s.next()
		// end of data
		if ch == nil {
			return nil
		}
		// verify
		for _, c := range cl {
			if ch[0] == c {
				s.pp++
				return ch
			}
		}
		return nil
	}
}

// VI match single string (case insensitive)
func VI(s string) Rule {
	s = strings.ToLower(s)
	r := make([]Rule, len(s))
	for i, ch1 := range s {
		ch2 := ch1
		if ch2 >= 0x61 && ch2 <= 0x7a {
			ch2 -= 0x20
		}
		r[i] = VL(ch1, ch2)
	}
	return C(r...)
}

// VIL match multiple strings (case insensitive)
func VIL(sl ...string) Rule {
	r := make([]Rule, len(sl))
	for i, s := range sl {
		r[i] = VI(s)
	}
	return A(r...)
}

// VS match single string (case sensitive)
func VS(s string) Rule {
	r := make([]Rule, len(s))
	for i, ch := range s {
		r[i] = V(ch)
	}
	return C(r...)
}

// VSL match multiple strings (case sensitive)
func VSL(sl ...string) Rule {
	r := make([]Rule, len(sl))
	for i, s := range sl {
		r[i] = VS(s)
	}
	return A(r...)
}

// C is Concatenation(Rule1 Rule2)
// match -> []rune of match value
// different or EOF -> nil
func C(r ...Rule) Rule {
	return func(s *scanner) []rune {
		s.mark()
		for _, f := range r {
			if f(s) == nil {
				return s.rollback()
			}
		}
		return s.commit()
	}
}

// A is Alternatives(Rule1 / Rule2)
// Incremental Alternatives: Rule1 =/ Rule2
// different or EOF -> nil
func A(r ...Rule) Rule {
	return func(s *scanner) []rune {
		for _, f := range r {
			if b := f(s); b != nil {
				return b
			}
		}
		return nil
	}
}

// VR is Value Range Alternatives(%c##-##)
// match -> []rune of match value
// different or EOF -> nil
func VR(h rune, t rune) Rule {
	return func(s *scanner) []rune {
		if ch := s.next(); ch != nil && ch[0] >= h && ch[0] <= t {
			s.pp++
			return ch
		}
		return nil
	}
}

// RN is Specific Repetition(nRule)
// match -> []rune of match value
// different or EOF -> nil
func RN(n int, r Rule) Rule {
	return func(s *scanner) []rune {
		return repet(s, n, n, r)
	}
}

// RV is Variable Repetition(*Rule)
// match -> []rune of match value
// different or EOF -> nil
// If max=-1 then max=infinity.
func RV(min int, max int, r Rule) Rule {
	return func(s *scanner) []rune {
		return repet(s, min, max, r)
	}
}

// R0 is default repeat count of Variable Repetition
// From zero to infinity
func R0(r Rule) Rule {
	return func(s *scanner) []rune {
		return repet(s, 0, -1, r)
	}
}

// R1 is more than one repeat count of Variable Repetition
// From one to infinity
func R1(r Rule) Rule {
	return func(s *scanner) []rune {
		return repet(s, 1, -1, r)
	}
}

// subroutin of Repeat
func repet(s *scanner, min int, max int, r Rule) []rune {
	s.mark()
	i := 0
	for ; max < 0 || i < max; i++ {
		if r(s) == nil {
			break
		}
	}
	// check minimum count
	if i < min {
		return s.rollback()
	}
	return s.commit()
}

// O is Optional Sequence([RULE])
// match -> []rune of match value
// different or EOF -> empty []rune
func O(r Rule) Rule {
	return func(s *scanner) []rune {
		b := r(s)
		if b == nil {
			b = make([]rune, 0)
		}
		return b
	}
}

/*
   Additional operator
*/

// N is Not match
// match -> nil
// different or EOF -> empty []rune
func N(r Rule) Rule {
	return func(s *scanner) []rune {
		var b []rune
		s.mark()
		if r(s) != nil {
			b = nil
		}
		s.rollback()
		return b
	}
}

// K add new tree node with key k
func K(r Rule, k int) Rule {
	return func(s *scanner) []rune {
		p := s.wp
		s.wp = &Tree{k, nil, nil, nil}

		v := r(s)
		if v != nil {
			s.wp.V = v
			p.add(s.wp)
		}
		s.wp = p
		return v
	}
}
