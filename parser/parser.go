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
	var maexp *MemberExpression
	var i int
	for i = 0; i < len(tokens); i++ { // skip first '{'
		t := tokens[i]
		var nt *lexer.Token
		if i+1 < len(tokens) {
			nt = &tokens[i+1]
		}
		if ptype, ok := t.ParsePrimitiveType(); ok || t.IsIdentifier() {
			var tmpExp Expression
			if ok {
				tmpExp = LiteralExpression{
					Type:   ptype,
					Value:  t.Value,
					Line:   t.Line,
					CharAt: t.CharAt,
				}
			} else if t.IsIdentifier() {
				tmpExp = VariableExpression{
					Name:   t.Value,
					Line:   t.Line,
					CharAt: t.CharAt,
				}
				if maexp != nil {
					maexp.PropertyIdentifier = Identifier{
						Name:   t.Value,
						Line:   t.Line,
						CharAt: t.CharAt,
					}
					if nt == nil || !nt.IsOperatorSymbol() {
						return *maexp, i, nil
					}
					bexp = &BinaryExpression{
						Left:    *maexp,
						Nesting: 0,
						Line:    t.Line,
						CharAt:  t.CharAt,
					}
					continue
				}
			}
			if bexp == nil {
				if nt == nil || !nt.IsOperatorSymbol() {
					return tmpExp, 1, nil
				}
				bexp = &BinaryExpression{
					Left:    tmpExp,
					Nesting: 0,
					Line:    t.Line,
					CharAt:  t.CharAt,
				}
			} else {
				exp, processed, err := parseExpression(tokens[i:])
				if err != nil {
					return nil, 0, err
				}
				if rightBexp, ok := exp.(BinaryExpression); ok {
					bexp.Nesting = rightBexp.Nesting + 1
					mergeBinaryExpression(bexp, rightBexp)
				} else {
					bexp.Right = exp
				}
				return *bexp, i + processed, nil
			}
		} else if t.IsOperatorSymbol() {
			if bexp == nil {
				return nil, 0, fmt.Errorf("TODO")
			}
			if t.Value == "." {
				maexp = &MemberExpression{
					Object: bexp.Left,
					Line:   bexp.Line,
					CharAt: bexp.CharAt,
				}
				bexp = nil
			} else {
				bexp.Operator = operator.Operator{
					Symbol: t.Value,
					Line:   t.Line,
					CharAt: t.CharAt,
				}
			}
		} else if t.Value == "(" {
			exp, processed, err := parseExpression(tokens[i+1:]) // i + 1 -> skip '('
			if err != nil {
				return nil, 0, err
			}
			openParenIndex := i
			i = i + processed + 1 // tokens[i].Value should be '(' now
			if i >= len(tokens) || tokens[i].Value != ")" {
				return nil, 0, fmt.Errorf("TODO")
			}
			switch e := exp.(type) {
			case LiteralExpression, BinaryExpression:
				if bexp == nil {
					expEnded := i+1 >= len(tokens) || !tokens[i+1].IsOperatorSymbol()
					bexp = &BinaryExpression{
						Left:    e,
						Nesting: 0,
						Line:    t.Line,
						CharAt:  t.CharAt,
					}

					bexpInParen, ok := e.(BinaryExpression)
					if ok {
						bexpInParen.Group = true
						bexp.Left = bexpInParen
						bexp.Nesting = bexpInParen.Nesting + 1
					}

					if expEnded {
						if ok {
							return bexpInParen, i, nil
						}
						return e, i, nil
					}
				} else {
					exp, processed, err := parseExpression(tokens[openParenIndex:])
					if err != nil {
						return nil, 0, err
					}
					if rightBexp, ok := exp.(BinaryExpression); ok {
						bexp.Nesting = rightBexp.Nesting + 1
						mergeBinaryExpression(bexp, rightBexp)
					} else {
						bexp.Right = exp
					}
					return *bexp, openParenIndex + processed, nil
				}
			case VariableExpression:
			case FunctionExpression:
			case ArrayExpression:
			case ObjectExpression:
			default:
				fmt.Printf("TODO: I don't know about type %T!\n", e)
			}
		} else if t.Value == "{" {
			return parseObjectExpression(tokens)
		}
	}
	return nil, 0, nil
}

func mergeBinaryExpression(bexp *BinaryExpression, rightBexp BinaryExpression) {
	if bexp.Operator.Compare(rightBexp.Operator) <= 0 && !rightBexp.Group {
		reorderBinaryExpression(bexp, rightBexp)
	} else {
		bexp.Right = rightBexp
	}
	root := bexp
	tmp := []*BinaryExpression{}
	for root.Nesting > 1 {
		if bexp1, ok := root.Left.(BinaryExpression); ok {
			if bexp2, ok := bexp1.Right.(BinaryExpression); ok {
				if bexp1.Operator.Compare(bexp2.Operator) <= 0 && !bexp2.Group {
					reorderBinaryExpression(&bexp1, bexp2)
				}
				tmp = append(tmp, &bexp1)
				root = &bexp1
			} else {
				break
			}
		} else {
			break
		}
	}
	if len(tmp) > 0 {
		for i := len(tmp) - 2; i >= 0; i-- {
			tmp[i].Left = *tmp[i+1]
		}
		bexp.Left = *tmp[0]
	}
}

func reorderBinaryExpression(bexp1 *BinaryExpression, bexp2 BinaryExpression) {
	bexp1.Left = BinaryExpression{
		Left:     bexp1.Left,
		Right:    bexp2.Left,
		Operator: bexp1.Operator,
		Line:     bexp1.Line,
		CharAt:   bexp1.CharAt,
		Nesting:  bexp2.Nesting,
	}
	bexp1.Operator = bexp2.Operator
	bexp1.Right = bexp2.Right
}
