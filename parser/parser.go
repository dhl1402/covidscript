package parser

import (
	"fmt"
	"gs/lexer"
	"gs/operator"
)

func ToAST(tokens []lexer.Token) ([]Statement, error) {
	ast := []Statement{}
	var i int
	for i = 0; i < len(tokens); i++ {
		t := tokens[i]
		switch t.Value {
		case "var":
			{
				s, processed, err := parseVariableDeclaration(tokens[i:])
				if err != nil {
					return nil, err
				}
				ast = append(ast, *s)
				i = i + processed - 1
			}
		}
	}
	return ast, nil
}

func parseVariableDeclaration(tokens []lexer.Token) (*VariableDeclaration, int, error) {
	if len(tokens) == 0 {
		return nil, 0, fmt.Errorf("TODO")
	}
	ids, i, err := parseSequentIdentifiers(tokens[1:]) // tokens[1:] -> skip 'var'
	if err != nil {
		return nil, 0, err
	}
	if len(ids) == 0 {
		return nil, 0, fmt.Errorf("TODO")
	}
	s := &VariableDeclaration{
		Line:   tokens[0].Line,
		CharAt: tokens[0].CharAt,
	}
	for _, id := range ids {
		s.Declarations = append(s.Declarations, VariableDeclarator{
			ID:     id,
			Init:   nil,
			Line:   id.Line,
			CharAt: id.CharAt,
		})
	}
	if i+1 >= len(tokens) || tokens[i+1].Value != "=" {
		return s, i + 1, nil
	}
	// start parsing variable initialization
	exps, processed, err := parseSequentExpressions(tokens[i+2:])
	if err != nil {
		return nil, 0, err
	}
	if exps == nil {
		return nil, 0, fmt.Errorf("TODO") // = without expressions
	}
	for ii, exp := range exps {
		s.Declarations[ii].Init = exp
	}
	return s, i + processed + 2, nil
}

func parseSequentIdentifiers(tokens []lexer.Token) ([]Identifier, int, error) {
	ids := []Identifier{}
	var i int
	for i = 0; i < len(tokens); i++ {
		t := tokens[i]
		if !t.IsIdentifier() {
			return ids, i, nil
		}
		if t.IsIdentifier() {
			ids = append(ids, Identifier{
				Name:   t.Value,
				Line:   t.Line,
				CharAt: t.CharAt,
			})
			i++ // increase before break or
			if i >= len(tokens) || tokens[i].Value != "," {
				break
			}
		}
	}
	return ids, i, nil
}

func parseSequentExpressions(tokens []lexer.Token) ([]Expression, int, error) {
	exps := []Expression{}
	var i int
	for i = 0; i < len(tokens); i++ {
		exp, processed, err := parseExpression(tokens[i:])
		if err != nil {
			return nil, 0, err
		}
		if exp == nil {
			return exps, i, nil
		}
		if exp != nil {
			exps = append(exps, exp)
			i = i + processed // do not need to - 1, skip ',' anyway
			if i >= len(tokens) || tokens[i].Value != "," {
				break
			}
		}
	}
	return exps, i, nil
}

// TODO: function, method, shorthand, compute
func parseObjectExpression(tokens []lexer.Token) (Expression, int, error) {
	if len(tokens) == 0 || tokens[0].Value != "{" {
		return nil, 0, fmt.Errorf("TODO")
	}
	obj := &ObjectExpression{
		Properties: []ObjectProperty{},
		Line:       tokens[0].Line,
		CharAt:     tokens[0].CharAt,
	}
	var prop *ObjectProperty
	var i int
	for i = 1; i < len(tokens); i++ { // skip first '{'
		t := tokens[i]
		pt := tokens[i-1]
		var nt *lexer.Token
		if i+1 < len(tokens) {
			nt = &tokens[i+1]
		}
		if t.Value == "}" {
			break
		}
		if t.Value == "," {
			if nt.Value == "}" {
				i++
				break
			}
			if prop != nil || !nt.IsIdentifier() { // TODO: or !compute prop key
				return nil, 0, fmt.Errorf("TODO")
			}
			continue
		}
		if prop == nil && pt.Value == "{" || pt.Value == "," {
			prop = &ObjectProperty{
				Line:   t.Line,
				CharAt: t.CharAt,
			}
			if t.IsIdentifier() {
				prop.KeyIdentifier = Identifier{
					Name:   t.Value,
					Line:   t.Line,
					CharAt: t.CharAt,
				}
			} else if t.Value == "[" {
				exp, processed, err := parseArrayExpression(tokens[i:])
				if err != nil {
					return nil, 0, err
				}
				aexp, _ := exp.(*ArrayExpression)
				if aexp == nil || len(aexp.Elements) != 1 {
					return nil, 0, fmt.Errorf("TODO")
				}
				prop.KeyExpression = aexp.Elements[0]
				prop.Computed = true
				i = i + processed - 1
			} else {
				return nil, 0, fmt.Errorf("TODO")
			}
		} else if prop != nil && t.Value == ":" {
			exp, processed, err := parseExpression(tokens[i+1:])
			if err != nil {
				return nil, 0, err
			}
			i = i + processed
			prop.Value = exp
			obj.Properties = append(obj.Properties, *prop)
			prop = nil
		} else {
			return nil, 0, fmt.Errorf("TODO")
		}
	}
	return obj, i + 1, nil
}

func parseArrayExpression(tokens []lexer.Token) (Expression, int, error) {
	if len(tokens) == 0 || tokens[0].Value != "[" {
		return nil, 0, fmt.Errorf("TODO")
	}
	exps, processed, err := parseSequentExpressions(tokens[1:]) // skip '['
	if err != nil {
		return nil, 0, err
	}
	if processed+1 >= len(tokens) || tokens[processed+1].Value != "]" {
		return nil, 0, fmt.Errorf("TODO")
	}
	return &ArrayExpression{
		Elements: exps,
		Line:     tokens[0].Line,
		CharAt:   tokens[0].CharAt,
	}, processed + 2, nil
}

// primitive, object, array, function, function call, binary expression, object property, array element, identifier
// identifier: call expression
func parseExpression(tokens []lexer.Token) (Expression, int, error) {
	if len(tokens) == 0 {
		return nil, 0, fmt.Errorf("TODO")
	}
	var tmpExp Expression
	var bexp *BinaryExpression
	var lastBexp *BinaryExpression
	var maexp *MemberAccessExpression
	bexpsAfterGroup := []*BinaryExpression{}
	openParen := 0
	var i int
	for i = 0; i < len(tokens); i++ {
		t := tokens[i]
		var nt *lexer.Token
		if i+1 < len(tokens) {
			nt = &tokens[i+1]
		}
		if exp, processed, _ := parseTempExpression(tokens[i:]); exp != nil {
			aexp, _ := exp.(*ArrayExpression)
			if tmpExp != nil && (aexp == nil || len(aexp.Elements) != 1) {
				return tmpExp, i, nil
			}
			i = i + processed - 1
			if tmpExp != nil && aexp != nil && len(aexp.Elements) == 1 {
				if lastBexp != nil {
					maexp = &MemberAccessExpression{
						Object:   lastBexp.Right,
						Property: aexp.Elements[0],
						Line:     t.Line, // TODO: lastBexp.Right.GetLine()
						CharAt:   lastBexp.Right.GetCharAt(),
					}
					lastBexp.Right = maexp
					lastBexp.Group = false
				} else if maexp == nil {
					if tmpExp == nil {
						return nil, 0, fmt.Errorf("TODO")
					}
					maexp = &MemberAccessExpression{
						Object:   tmpExp,
						Property: aexp.Elements[0],
						Line:     t.Line, // TODO: tmpExp.GetLine()
						CharAt:   tmpExp.GetCharAt(),
					}
				} else if maexp.Property != nil {
					maexp = &MemberAccessExpression{
						Object:   maexp,
						Property: aexp.Elements[0],
						Line:     maexp.Line,
						CharAt:   maexp.CharAt,
					}
				} else {
					return nil, 0, fmt.Errorf("TODO")
				}
				tmpExp = maexp
			} else {
				tmpExp = exp
			}
			if maexp != nil && maexp.Property == nil {
				maexp.Property = tmpExp
				tmpExp = maexp
			}
			if bexp != nil && nt != nil && nt.Value == ")" && openParen <= len(bexpsAfterGroup) {
				bexp.Group = true
			}
			if bexp != nil && bexp.Right == nil {
				bexp.Right = tmpExp
				lastBexp = bexp
				lastBexpAfterGroup, _ := getLastBexpAfterGroup(bexpsAfterGroup, false)
				// Check prececdence and swap expression order if needed
				if lbexp, ok := bexp.Left.(*BinaryExpression); ok && lbexp.Operator.Compare(bexp.Operator) > 0 && !lbexp.Group && (lastBexpAfterGroup == nil || bexp != lastBexpAfterGroup) {
					bexp.Left = lbexp.Right
					bexp.CharAt = lbexp.Right.GetCharAt()
					lbexp.Right = bexp
					bexp = lbexp
				}
			}
		} else if t.Value == "." {
			if lastBexp != nil {
				maexp = &MemberAccessExpression{
					Object: lastBexp.Right,
					Line:   t.Line, // TODO: lastBexp.Right.GetLine()
					CharAt: lastBexp.Right.GetCharAt(),
				}
				lastBexp.Right = maexp
				lastBexp.Group = false
			} else if maexp == nil {
				if tmpExp == nil {
					return nil, 0, fmt.Errorf("TODO")
				}
				maexp = &MemberAccessExpression{
					Object: tmpExp,
					Line:   t.Line, // TODO: tmpExp.GetLine()
					CharAt: tmpExp.GetCharAt(),
				}
			} else if maexp.Property != nil {
				maexp = &MemberAccessExpression{
					Object: maexp,
					Line:   maexp.Line,
					CharAt: maexp.CharAt,
				}
			} else {
				return nil, 0, fmt.Errorf("TODO")
			}
			tmpExp = nil
		} else if t.IsOperatorSymbol() {
			op := operator.Operator{
				Symbol: t.Value,
				Line:   t.Line,
				CharAt: t.CharAt,
			}
			if bexp == nil {
				if tmpExp == nil {
					return nil, 0, fmt.Errorf("TODO")
				}
				if tmpExp == maexp {
					maexp = nil
				}
				bexp = &BinaryExpression{
					Left:     tmpExp,
					Line:     t.Line, // tmpExp.getCharAt()
					CharAt:   tmpExp.GetCharAt(),
					Operator: op,
				}
			} else if bexp.Right != nil {
				bexp = &BinaryExpression{
					Left:     bexp,
					Line:     bexp.Line,
					CharAt:   bexp.CharAt,
					Operator: op,
				}
			} else {
				return nil, 0, fmt.Errorf("TODO")
			}
			if bexp != nil && openParen > 0 {
				bexpsAfterGroup = append(bexpsAfterGroup, bexp)
				openParen--
			}
			tmpExp = nil
		} else if t.Value == "(" {
			// TODO: handle function call
			openParen++
		} else if t.Value == ")" {
			// TODO: handle function call
			if openParen >= len(bexpsAfterGroup) {
				// if openParen is increased by some exp like (a)
				// TODO: check if not function call
				openParen--
			} else if bexp != nil {
				if len(bexpsAfterGroup) > 0 {
					// Swap expression order to match expression group
					lastBexpAfterGroup, i := getLastBexpAfterGroup(bexpsAfterGroup, true)
					if bexpBeforeGroup, ok := lastBexpAfterGroup.Left.(*BinaryExpression); ok {
						lastBexpAfterGroup.Left = bexpBeforeGroup.Right
						lastBexpAfterGroup.CharAt = bexpBeforeGroup.Right.GetCharAt()
						bexpBeforeGroup.Right = bexp
						groupedBexp := bexp
						bexp = bexpBeforeGroup

						// loop from groupedBexp to lastBexpAfterGroup and update exp CharAt
						groupedBexp.CharAt = lastBexpAfterGroup.CharAt
						for tmpBexp, _ := groupedBexp.Left.(*BinaryExpression); tmpBexp != nil; {
							tmpBexp.CharAt = lastBexpAfterGroup.CharAt
							tmpBexp, _ = tmpBexp.Left.(*BinaryExpression)
						}
					}
					bexpsAfterGroup = append(bexpsAfterGroup[i+1:], bexpsAfterGroup[:i]...)
				}
			}
		} else {
			break
		}
	}
	if bexp != nil {
		return bexp, i, nil
	}
	if maexp != nil {
		return maexp, i, nil
	}
	return tmpExp, i, nil
}

// Check first token, if it is the start of an expression then parse it and return
func parseTempExpression(tokens []lexer.Token) (Expression, int, error) {
	if len(tokens) == 0 {
		return nil, 0, fmt.Errorf("TODO")
	}
	t := tokens[0]
	if ptype, ok := t.ParsePrimitiveType(); ok {
		return &LiteralExpression{
			Type:   ptype,
			Value:  t.Value,
			Line:   t.Line,
			CharAt: t.CharAt,
		}, 1, nil
	}
	if t.IsIdentifier() {
		return &VariableExpression{
			Name:   t.Value,
			Line:   t.Line,
			CharAt: t.CharAt,
		}, 1, nil
	}
	if t.Value == "{" {
		return parseObjectExpression(tokens)
	}
	if t.Value == "[" {
		return parseArrayExpression(tokens)
	}
	return nil, 0, fmt.Errorf("TODO")
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
