package scanner

func isAlpha(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

func (s *Scanner) isAtEnd() bool {
	return s.currentPos >= len(s.src)
}

func (s *Scanner) next() rune {
	s.currentPos++
	return rune(s.src[s.currentPos-1])
}

func (s *Scanner) current() rune {
	return rune(s.src[s.currentPos])
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return EOF
	}

	return s.current()
}

func (s *Scanner) peekNext() rune {
	if s.currentPos+1 >= len(s.src) {
		return EOF
	}

	return rune(s.src[s.currentPos+1])
}

func (s *Scanner) match(c rune) bool {
	if !s.isAtEnd() && s.current() == c {
		s.currentPos++
		return true
	}

	return false
}
