package abnf

import "strings"

/*
   Terminal Values

   Rules resolve into a string of terminal values, sometimes called
   characters.  In ABNF, a character is merely a non-negative integer.
   In certain contexts, a specific mapping (encoding) of values into a
   character set (such as ASCII) will be specified.
   
   Terminals are specified by one or more numeric characters, with the
   base interpretation of those characters indicated explicitly.  The
   following bases are currently defined:
         b           =  binary
         d           =  decimal
         x           =  hexadecimal
   Hence:
         CR          =  %d13
         CR          =  %x0D
   respectively specify the decimal and hexadecimal representation of
   [US-ASCII] for carriage return.

   A concatenated string of such values is specified compactly, using a
   period (".") to indicate a separation of characters within that
   value.  Hence:
         CRLF        =  %d13.10
   ABNF permits the specification of literal text strings directly,
   enclosed in quotation marks.  Hence:
         command     =  "command string"
   Literal text strings are interpreted as a concatenated set of
   printable characters.

   NOTE:
      ABNF strings are case insensitive and the character set for these
      strings is US-ASCII.

   Hence:
         rulename = "abc"
   and:
         rulename = "aBc"
   will match "abc", "Abc", "aBc", "abC", "ABc", "aBC", "AbC", and
   "ABC".

      To specify a rule that is case sensitive, specify the characters
      individually.

   For example:
         rulename    =  %d97 %d98 %d99
   or
         rulename    =  %d97.98.99
   will match only the string that comprises only the lowercase
   characters, abc.
   
   different or EOF -> nil
*/
func V(c rune) Rule {
	return func(s *Scanner) []rune {
		if ch := s.next(); ch != nil && ch[0] == c {
            s.pp++
            return ch
        }
		return nil
	}
}

func VS(str string) Rule {
    str = strings.ToLower(str)
    fs := make([]Rule, len(str))
	for i, ch1 := range str {
        ch2 := ch1
        if ch2 >= 0x61 && ch2 <= 0x7a {
            ch2 -= 0x20
        }
		fs[i] = VL(ch1, ch2)
	}
    
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

/*
   Concatenation:  Rule1 Rule2

   A rule can define a simple, ordered string of values (i.e., a
   concatenation of contiguous characters) by listing a sequence of rule
   names.  For example:
         foo         =  %x61           ; a
         bar         =  %x62           ; b
         mumble      =  foo bar foo
   So that the rule <mumble> matches the lowercase string "aba".

   Linear white space: Concatenation is at the core of the ABNF parsing
   model.  A string of contiguous characters (values) is parsed
   according to the rules defined in ABNF.  For Internet specifications,
   there is some history of permitting linear white space (space and
   horizontal tab) to be freely and implicitly interspersed around major
   constructs, such as delimiting special characters or atomic strings.

   NOTE:
      This specification for ABNF does not provide for implicit
      specification of linear white space.

   Any grammar that wishes to permit linear white space around
   delimiters or string segments must specify it explicitly.  It is
   often useful to provide for such white space in "core" rules that are
   then used variously among higher-level rules.  The "core" rules might
   be formed into a lexical analyzer or simply be part of the main
   ruleset.
*/
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

/*
   Alternatives:  Rule1 / Rule2
   
   Elements separated by a forward slash ("/") are alternatives.
   Therefore,
         foo / bar
   will accept <foo> or <bar>.

   NOTE:
      A quoted string containing alphabetic characters is a special form
      for specifying alternative characters and is interpreted as a non-
      terminal representing the set of combinatorial strings with the
      contained characters, in the specified order but with any mixture
      of upper- and lowercase.
*/
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

/*
   Incremental Alternatives: Rule1 =/ Rule2

   It is sometimes convenient to specify a list of alternatives in
   fragments.  That is, an initial rule may match one or more
   alternatives, with later rule definitions adding to the set of
   alternatives.  This is particularly useful for otherwise independent
   specifications that derive from the same parent ruleset, such as
   often occurs with parameter lists.  ABNF permits this incremental
   definition through the construct:
         oldrule     =/ additional-alternatives

   So that the ruleset
         ruleset     =  alt1 / alt2
         ruleset     =/ alt3
         ruleset     =/ alt4 / alt5
   is the same as specifying
         ruleset     =  alt1 / alt2 / alt3 / alt4 / alt5

   different or EOF -> nil
*/
func IA(fs ...Rule) Rule {
	return A(fs...)
}

/*
   Value Range Alternatives:  %c##-##

   A range of alternative numeric values can be specified compactly,
   using a dash ("-") to indicate the range of alternative values.
   Hence:
         DIGIT       =  %x30-39
   is equivalent to:
         DIGIT       =  "0" / "1" / "2" / "3" / "4" / "5" / "6" /
                        "7" / "8" / "9"
   Concatenated numeric values and numeric value ranges cannot be
   specified in the same string.  A numeric value may use the dotted
   notation for concatenation or it may use the dash notation to specify
   one value range.  Hence, to specify one printable character between
   end-of-line sequences, the specification could be:
         char-line = %x0D.0A %x20-7E %x0D.0A
*/
func VR(head rune, tail rune) Rule {
	return func(s *Scanner) []rune {
		if ch := s.next(); ch != nil && ch[0] >= head && ch[0] <= tail {
			s.pp++
			return ch
		}
		return nil
	}
}

/*
   Specific Repetition:  nRule

   A rule of the form:
         <n>element
   is equivalent to
         <n>*<n>element
   That is, exactly <n> occurrences of <element>.  Thus, 2DIGIT is a
   2-digit number, and 3ALPHA is a string of three alphabetic
   characters.
*/
func RS(n int, fs Rule) Rule {
	return func(s *Scanner) []rune {
		return repet(s, n, n, fs)
	}
}

/*
   Variable Repetition:  *Rule

   The operator "*" preceding an element indicates repetition.  The full
   form is:
         <a>*<b>element
   where <a> and <b> are optional decimal values, indicating at least
   <a> and at most <b> occurrences of the element.

   Default values are 0 and infinity so that *<element> allows any
   number, including zero; 1*<element> requires at least one;
   3*3<element> allows exactly 3; and 1*2<element> allows one or two.
   
   If max=-1 then max=infinity.
*/
func RV(min int, max int, fs Rule) Rule {
	return func(s *Scanner) []rune {
		return repet(s, min, max, fs)
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

/*
   Optional Sequence:  [RULE]

   Square brackets enclose an optional element sequence:
         [foo bar]
   is equivalent to
         *1(foo bar).
*/
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
/*
   Alternatives of Terminal Values
*/
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

/*
   Alternatives of Terminal String Values
*/
func VSL(slist ...string) Rule {
    fs := make([]Rule, len(slist))
	for i, s := range slist {
        fs[i] = VS(s)
	}
    
	return func(s *Scanner) []rune {
		for _, f := range fs {
            if r := f(s); r != nil {
       			return r
			}
		}
		return nil
	}
}

/*
   Not match
*/
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

/*
   Default repeat count of Variable Repetition
   From zero to infinity
*/
func R0(fs Rule) Rule {
	return func(s *Scanner) []rune {
		return repet(s, 0, -1, fs)
	}
}

/*
   More than one repeat count of Variable Repetition
   From one to infinity
*/
func R1(fs Rule) Rule {
	return func(s *Scanner) []rune {
		return repet(s, 1, -1, fs)
	}
}
