package traindown

import (
	"fmt"
)

// TokenType is just the token enum for use in a Parser
type TokenType int

// Token contains a TokenType and the raw input
type Token struct {
	Type  TokenType
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("[%v] %v", t.Type, t.Value)
}
