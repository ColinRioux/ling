package parser

import (
	"duck/ling/js/lexer/token"
	"duck/ling/js/parser/ast/node"
	"duck/ling/js/syntax/literal"
	"duck/ling/js/syntax/unicode"
	"fmt"
	"strconv"
)

var Pos = 0
var CurrentToken *token.Token = nil
var CurrentChar rune
var Text string

func Peek() rune {
	peekPos := Pos + 1
	if peekPos > len(Text)-1 {
		return 0
	}
	return rune(Text[peekPos])
}

func SkipWhitespace() {
	for CurrentChar != 0 && unicode.IsWhitespace(CurrentChar) {
		Advance()
	}
}

func GetInt() int {
	var res string = ""
	for CurrentChar != 0 && literal.IsDecimalDigit(CurrentChar) {
		res += string(CurrentChar)
		Advance()
	}
	resI, _ := strconv.Atoi(res)
	return resI
}

func GetNextToken() *token.Token {
	for CurrentChar != 0 {
		if unicode.IsWhitespace(CurrentChar) {
			SkipWhitespace()
			continue
		}

		if literal.IsAlpha(CurrentChar) || CurrentChar == '_' {
			return ParseIdentifier()
		}

		if CurrentChar == '=' {
			Advance()
			return token.NewToken(token.ASSIGN, "=")
		}

		if CurrentChar == ';' {
			Advance()
			return token.NewToken(token.SEMICOLON, ";")
		}

		if literal.IsDecimalDigit(CurrentChar) || CurrentChar == '.' {
			return token.NewToken(token.NUMBER, ParseNumber())
		}

		if CurrentChar == '+' {
			Advance()
			return token.NewToken(token.ADD, "+")
		}

		if CurrentChar == '-' {
			Advance()
			return token.NewToken(token.SUB, "-")
		}

		if CurrentChar == '*' {
			Advance()
			return token.NewToken(token.MUL, "*")
		}

		if CurrentChar == '/' {
			Advance()
			return token.NewToken(token.DIV, "/")
		}

		if CurrentChar == '(' {
			Advance()
			return token.NewToken(token.LPAREN, "(")
		}

		if CurrentChar == ')' {
			Advance()
			return token.NewToken(token.RPAREN, ")")
		}
	}
	return token.NewToken(token.EOF, 0)
}

func eat(tokenType token.Type) {
	if tokenType == CurrentToken.Type {
		CurrentToken = GetNextToken()
	}
}

func Advance() {
	Pos++
	if Pos > len(Text)-1 {
		CurrentChar = 0
	} else {
		CurrentChar = rune(Text[Pos])
	}
}

func Factor() *node.IASTNode {
	tok := CurrentToken
	var nde node.IASTNode
	if tok.Type == token.ADD {
		eat(token.ADD)
		nde = node.NewUnaryOperatorNode(nil, nil, tok, Factor())
		return &nde
	} else if tok.Type == token.SUB {
		eat(token.SUB)
		nde = node.NewUnaryOperatorNode(nil, nil, tok, Factor())
		return &nde
	} else if tok.Type == token.LPAREN {
		eat(token.LPAREN)
		nde := Expr()
		eat(token.RPAREN)
		return nde
	} else if tok.Type == token.NUMBER {
		eat(token.NUMBER)
		nde = node.NewNumberNode(tok)
		return &nde
	}
	nde = node.NewASTNode(nil, nil, tok)
	return &nde
}

func Term() *node.IASTNode {
	var nde *node.IASTNode = Factor()
	for CurrentToken.Type == token.MUL || CurrentToken.Type == token.DIV {
		tok := CurrentToken
		if CurrentToken.Type == token.MUL {
			eat(token.MUL)
		} else if CurrentToken.Type == token.DIV {
			eat(token.DIV)
		}
		var tmp node.IASTNode = node.NewBinaryOperatorNode(nde, Term(), tok)
		nde = &tmp
	}
	return nde
}

func Expr() *node.IASTNode {
	var nde *node.IASTNode = Term()
	for CurrentToken.Type == token.ADD || CurrentToken.Type == token.SUB {
		tok := CurrentToken
		if tok.Type == token.ADD {
			eat(token.ADD)
		} else if tok.Type == token.SUB {
			eat(token.SUB)
		}
		var tmp node.IASTNode = node.NewBinaryOperatorNode(nde, Term(), tok)
		nde = &tmp
	}

	return nde
}

func Visit(root *node.IASTNode) {
	if root == nil {
		return
	}
	Visit((*root).GetLeft())
	Visit((*root).GetRight())
	fmt.Println((*root).GetToken())
}
