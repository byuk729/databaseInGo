package lexer

func (token *Token) Equals(otherToken *Token) bool {
	return token.Type == otherToken.Type && token.Value == otherToken.Value
}
