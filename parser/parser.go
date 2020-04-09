package parser

import (
	"fmt"
	"gs/lexer"
	"gs/operator"
)

func ToAST(tokens []lexer.Token) ([]Statement, error) {
	ast := []Statement{}
	for len(tokens) > 0 {
		t := tokens[0]
		switch t.Value {
		case "var":
			{
				s, pt, err := parseVariableDeclaration(tokens)
				if err != nil {
					return nil, err
				}
				ast = append(ast, *s)
				tokens = tokens[pt:]
			}
		default:
			tokens = tokens[1:]
		}
	}
	return ast, nil
}

func parseVariableDeclaration(tokens []lexer.Token) (*VariableDeclaration, int, error) {
	s := &VariableDeclaration{
		Line:   tokens[0].Line,
		CharAt: tokens[0].CharAt,
	}
	var i int
	for i = 1; i < len(tokens); i++ { // i = 1 -> skip 'var'
		t := tokens[i]
		if t.IsIdentifier() {
			if len(s.Declarations) != 0 && tokens[i-1].Value != "," {
				i--
				break
			}
			s.Declarations = append(s.Declarations, VariableDeclarator{
				ID: Identifier{
					Name:   t.Value,
					Line:   t.Line,
					CharAt: t.CharAt,
				},
				Init:   nil,
				Line:   t.Line,
				CharAt: t.CharAt,
			})
		} else if t.Value == "=" {
			exps, processed, err := parseVariableInitialization(tokens[i+1:])
			if err != nil {
				return nil, 0, err
			}
			if exps == nil {
				return nil, 0, fmt.Errorf("TODO")
			}
			for i, exp := range exps {
				s.Declarations[i].Init = exp
			}
			i = i + processed - 1
			break
		} else if t.Value == "," {
			if !tokens[i-1].IsIdentifier() {
				return nil, 0, fmt.Errorf("TODO")
			}
			if i+1 >= len(tokens) || !tokens[i+1].IsIdentifier() {
				return nil, 0, fmt.Errorf("TODO")
			}
		} else {
			i--
			break
		}
	}

	return s, i, nil
}

func parseVariableInitialization(tokens []lexer.Token) ([]Expression, int, error) {
	exps := []Expression{}
	var i int
	for i = 0; i < len(tokens); i++ {
		t := tokens[i]
		if ptype, ok := t.ParsePrimitiveType(); ok {
			if len(exps) != 0 && tokens[i-1].Value != "," {
				i--
				break
			}
			exps = append(exps, LiteralExpression{
				Type:   ptype,
				Value:  t.Value,
				Line:   t.Line,
				CharAt: t.CharAt,
			})
		} else if t.IsIdentifier() {
			if len(exps) != 0 && tokens[i-1].Value != "," {
				i--
				break
			}
			exps = append(exps, VariableExpression{
				Name:   t.Value,
				Line:   t.Line,
				CharAt: t.CharAt,
			})
		} else if t.Value == "{" {
			obj, processed, err := parseObjectExpression(tokens[i:])
			if err != nil {
				return nil, 0, err
			}
			i = i + processed - 1
			exps = append(exps, obj)
		} else if t.Value == "," {
			// object, array, function,... is also ok
			if i == 0 || (!tokens[i-1].IsIdentifier() && !tokens[i-1].IsPrimitiveValue()) {
				return nil, 0, fmt.Errorf("TODO")
			}
			// object, array, function,... is also ok
			if i+1 >= len(tokens) || (!tokens[i+1].IsIdentifier() && !tokens[i+1].IsPrimitiveValue()) {
				return nil, 0, fmt.Errorf("TODO")
			}
		} else {
			i--
			break
		}
	}
	return exps, i, nil
}

func parseObjectExpression(tokens []lexer.Token) (Expression, int, error) {
	obj := ObjectExpression{
		Line:   tokens[0].Line,
		CharAt: tokens[0].CharAt,
	}
	var prop *ObjectProperty
	var i int
	for i = 1; i < len(tokens); i++ { // skip first '{'
		t := tokens[i]
		pt := tokens[i-1]
		if t.Value == ":" || t.Value == "," {
			continue
		}
		if prop == nil && pt.Value == "{" || pt.Value == "," {
			if t.IsIdentifier() {
				prop = &ObjectProperty{
					KeyIdentifier: Identifier{
						Name:   t.Value,
						Line:   t.Line,
						CharAt: t.CharAt,
					},
				}
			} else if t.Value == "[" {

			} else {
				return nil, 0, fmt.Errorf("TODO")
			}
		} else if prop != nil && pt.Value == ":" {
			exp, processed, err := parseExpression(tokens[i:])
			if err != nil {
				return nil, 0, err
			}
			i = i + processed - 1
			prop.Value = exp
			obj.Properties = append(obj.Properties, *prop)
			prop = nil
		} else {
			i--
			break
		}
	}
	return obj, i, nil
}

// primitive, object, array, function, function call, binary expression, object property, array element, identifier
// identifier: variable expression, binary expression, call expression, member access exression
func parseExpression(tokens []lexer.Token) (Expression, int, error) {
	var bexp *BinaryExpression
	bexpsAfterGroup := []*BinaryExpression{}
	openParen := 0
	var i int
	for i = 0; i < len(tokens); i++ {
		t := tokens[i]
		var nt *lexer.Token
		// var pt *lexer.Token
		// if i > 0 {
		// 	pt = &tokens[i-1]
		// }
		if i+1 < len(tokens) {
			nt = &tokens[i+1]
		}
		if ptype, ok := t.ParsePrimitiveType(); ok {
			var tmpExp Expression
			tmpExp = &LiteralExpression{
				Type:   ptype,
				Value:  t.Value,
				Line:   t.Line,
				CharAt: t.CharAt,
			}
			if bexp == nil {
				if nt == nil || !nt.IsOperatorSymbol() {
					return tmpExp, 1, nil
				}
				bexp = &BinaryExpression{
					Left:   tmpExp,
					Line:   t.Line,
					CharAt: t.CharAt,
				}
			} else {
				bexp.Right = tmpExp
				if nt != nil && nt.Value == ")" {
					bexp.Group = true
				}
				lastBexpAfterGroup, _ := getLastBexpAfterGroup(bexpsAfterGroup, false)
				if lbexp, ok := bexp.Left.(*BinaryExpression); ok && lbexp.Operator.Compare(bexp.Operator) > 0 && !lbexp.Group && (lastBexpAfterGroup == nil || lbexp != lastBexpAfterGroup.Left) {
					bexp.Left = lbexp.Right
					bexp.CharAt = lbexp.Right.GetCharAt()
					lbexp.Right = bexp
					bexp = &BinaryExpression{
						Left:   lbexp,
						Line:   lbexp.Line,
						CharAt: lbexp.CharAt,
					}
				} else {
					bexp = &BinaryExpression{
						Left:   bexp,
						Line:   bexp.Line,
						CharAt: bexp.CharAt,
					}
				}
			}
			if openParen > 0 {
				bexpsAfterGroup = append(bexpsAfterGroup, bexp)
				openParen--
			}
		} else if t.IsOperatorSymbol() {
			if bexp == nil {
				return nil, 0, fmt.Errorf("TODO")
			}
			bexp.Operator = operator.Operator{
				Symbol: t.Value,
				Line:   t.Line,
				CharAt: t.CharAt,
			}
		} else if t.Value == "(" {
			// TODO: handle function call
			openParen++
		} else if t.Value == ")" {
			// TODO: handle function call
			if bexp != nil {
				if lbexp, ok := bexp.Left.(*BinaryExpression); ok {
					if len(bexpsAfterGroup) > 0 {
						lastBexpAfterGroup, i := getLastBexpAfterGroup(bexpsAfterGroup, true)
						if bexpBeforeGroup, ok := lastBexpAfterGroup.Left.(*BinaryExpression); ok {
							lastBexpAfterGroup.Left = bexpBeforeGroup.Right
							lastBexpAfterGroup.CharAt = bexpBeforeGroup.Right.GetCharAt()
							bexpBeforeGroup.Right = lbexp
							bexp.Left = bexpBeforeGroup

							// loop from lbexp to lastBexpAfterGroup and update exp CharAt
							lbexp.CharAt = lastBexpAfterGroup.CharAt
							for tmpBexp, _ := lbexp.Left.(*BinaryExpression); tmpBexp != nil; {
								tmpBexp.CharAt = lastBexpAfterGroup.CharAt
								tmpBexp, _ = tmpBexp.Left.(*BinaryExpression)
							}
						}
						bexpsAfterGroup = append(bexpsAfterGroup[i+1:], bexpsAfterGroup[:i]...)
					}
				}
			}
		} else if t.Value == "{" {
			return parseObjectExpression(tokens)
		}
	}
	return bexp.Left, i, nil
}

func getLastBexpAfterGroup(bexps []*BinaryExpression, rightShouldNotNil bool) (*BinaryExpression, int) {
	for i := len(bexps) - 1; i >= 0; i-- {
		bexp := bexps[i]
		if !rightShouldNotNil {
			return bexp, i
		}
		if bexp.Right != nil {
			return bexp, i
		}
	}
	return nil, -1
}
