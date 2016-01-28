package abnf

/*
   Core ABNF rules
*/

// ALPHA = %x41-5A / %x61-7A
func ALPHA() Rule {
	return A(VR(0x41, 0x5a), VR(0x61, 0x7a))
}

// DIGIT = %x30-39
func DIGIT() Rule {
	return VR(0x30, 0x39)
}

// HEXDIG = DIGIT / %x41-46
func HEXDIG() Rule {
	return A(DIGIT(), VR(0x41, 0x46))
}

// DQUOTE = %x22
func DQUOTE() Rule {
	return V(0x22)
}

// SP = %x20
func SP() Rule {
	return V(0x20)
}

// HTAB = %x09
func HTAB() Rule {
	return V(0x09)
}

// WSP = SP / HTAB
func WSP() Rule {
	return A(SP(), HTAB())
}

// LWSP = *(WSP / CRLF WSP)
func LWSP() Rule {
	return R0(A(WSP(), C(WSP(), CRLF())))
}

// VCHAR = %x21-7E
func VCHAR() Rule {
	return VR(0x21, 0x7e)
}

// CHAR = %x01-7F
func CHAR() Rule {
	return VR(0x01, 0x7f)
}

// OCTET = %x00-FF
func OCTET() Rule {
	return VR(0x00, 0xff)
}

// CTL = %x00-1F / %x7F
func CTL() Rule {
	return A(VR(0x00, 0x1f), V(0x7f))
}

// CR = %x0D
func CR() Rule {
	return V(0x0d)
}

// LF = %x0A
func LF() Rule {
	return V(0x0a)
}

// CRLF = CR LF
func CRLF() Rule {
	return C(CR(), LF())
}

// BIT = "0" / "1"
func BIT() Rule {
	return VL(0x30, 0x31)
}

// LHEX = DIGIT / %x61-66
func LHEX() Rule {
	return A(DIGIT(), VR(0x61, 0x66))
}

// ALPHANUM = ALPHA / DIGIT
func ALPHANUM() Rule {
	return A(DIGIT(), ALPHA())
}


/*
   Additional rules
*/

// EOF = [EOF]
func EOF() Rule {
	return func(s *Scanner) []rune {
		if s.next() == nil {
			return make([]rune, 0)
		}
		return nil
	}
}
