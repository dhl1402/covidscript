package parser

import (
	"fmt"

	"github.com/dhl1402/covidscript/core"
	"github.com/dhl1402/covidscript/lexer"
)

func ToAST(tokens []lexer.Token) ([]core.Statement, error) {
	ast, _, err := parseStatements(tokens)
	return ast, err
}

func parseStatements(tokens []lexer.Token) ([]core.Statement, int, error) {
	ss := []core.Statement{}
	var i int
	for i = 0; i < len(tokens); i++ {
		t := tokens[i]
		switch {
		case t.Value == "var":
			s, processed, err := parseVariableDeclaration(tokens[i:])
			if err != nil {
				return nil, 0, err
			}
			ss = append(ss, *s)
			i = i + processed - 1
		case t.Value == "func":
			s, processed, err := parseFunctionDeclaration(tokens[i:])
			if err != nil {
				return nil, 0, err
			}
			ss = append(ss, *s)
			i = i + processed - 1
		case t.Value == "return":
			s, processed, err := parseReturnStatement(tokens[i:])
			if err != nil {
				return nil, 0, err
			}
			ss = append(ss, *s)
			i = i + processed - 1
		case t.Value == "}":
			return ss, i, nil
		default:
			e, processed, err := parseExpression(tokens[i:])
			if err != nil {
				return nil, 0, err
			}
			if i+processed < len(tokens) && (tokens[i+processed].Value == "=" || tokens[i+processed].Value == ":=") && e != nil {
				// parse AssignmentStatement
				i = i + processed + 1 // handle '=' -> +1
				switch e.(type) {
				case *core.VariableExpression:
				case *core.MemberAccessExpression:
				default:
					return nil, 0, fmt.Errorf("TODO") // left of assignment must be VariableExpression or MemberAccessExpression
				}
				as := core.AssignmentStatement{
					Left:   e,
					Line:   t.Line,
					CharAt: t.CharAt,
				}
				rightExp, processed, err := parseExpression(tokens[i:])
				if err != nil {
					return nil, 0, err
				}
				as.Right = rightExp
				ss = append(ss, as)
				i = i + processed - 1
			} else if e != nil {
				// parse ExpressionStatement
				ss = append(ss, core.ExpressionStatement{
					Expression: e,
					Line:       e.GetLine(),
					CharAt:     e.GetCharAt(),
				})
				i = i + processed - 1
			}
		}
	}
	return ss, i, nil
}

func parseFunctionDeclaration(tokens []lexer.Token) (*core.FunctionDeclaration, int, error) {
	if len(tokens) < 6 { // 6 is len of the most simple function
		return nil, 0, fmt.Errorf("TODO")
	}
	if !tokens[1].IsIdentifier() {
		return nil, 0, fmt.Errorf("TODO")
	}
	if tokens[2].Value != "(" {
		return nil, 0, fmt.Errorf("TODO")
	}
	i := 2
	f := &core.FunctionDeclaration{
		ID: core.Identifier{
			Name:   tokens[1].Value,
			Line:   tokens[1].Line,
			CharAt: tokens[1].CharAt,
		},
		Line:   tokens[0].Line,
		CharAt: tokens[0].CharAt,
	}
	params, statements, processed, err := parseFunctionParamAndBody(tokens[i:])
	if err != nil {
		return nil, 0, err
	}
	f.Params = params
	f.Body = statements
	return f, i + processed, nil
}

func parseReturnStatement(tokens []lexer.Token) (*core.ReturnStatement, int, error) {
	if len(tokens) == 0 {
		return nil, 0, fmt.Errorf("TODO")
	}
	r := &core.ReturnStatement{
		Line:   tokens[0].Line,
		CharAt: tokens[0].CharAt,
	}
	exp, i, err := parseExpression(tokens[1:]) // skip 'return'
	if err != nil {
		return nil, 0, err
	}
	r.Argument = exp
	return r, i + 1, nil
}

func parseVariableDeclaration(tokens []lexer.Token) (*core.VariableDeclaration, int, error) {
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
	s := &core.VariableDeclaration{
		Line:   tokens[0].Line,
		CharAt: tokens[0].CharAt,
	}
	for _, id := range ids {
		s.Declarations = append(s.Declarations, core.VariableDeclarator{
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
	if len(exps) > len(s.Declarations) {
		return nil, 0, fmt.Errorf("TODO") // = too many expressions
	}
	for ii, exp := range exps {
		s.Declarations[ii].Init = exp
	}
	return s, i + processed + 2, nil
}

// primitive, object, array, function, function call, binary expression, object property, array element, identifier
func parseExpression(tokens []lexer.Token) (core.Expression, int, error) {
	if len(tokens) == 0 {
		return nil, 0, fmt.Errorf("TODO")
	}
	var tmpExp core.Expression
	var bexp *core.BinaryExpression
	var lastBexp *core.BinaryExpression
	var maexp *core.MemberAccessExpression
	bexpsAfterGroup := []*core.BinaryExpression{}
	openParen := 0
	var i int
	for i = 0; i < len(tokens); i++ {
		t := tokens[i]
		var nt *lexer.Token
		if i+1 < len(tokens) {
			nt = &tokens[i+1]
		}
		if exp, processed, err := parseTempExpression(tokens[i:]); exp != nil {
			aexp, _ := exp.(*core.ArrayExpression)
			if tmpExp != nil && (aexp == nil || len(aexp.Elements) != 1) {
				return tmpExp, i, nil
			}
			i = i + processed - 1
			if tmpExp != nil && aexp != nil && len(aexp.Elements) == 1 {
				if lastBexp != nil {
					maexp = &core.MemberAccessExpression{
						Object:             lastBexp.Right,
						PropertyExpression: aexp.Elements[0],
						Compute:            true,
						Line:               lastBexp.Right.GetLine(),
						CharAt:             lastBexp.Right.GetCharAt(),
					}
					lastBexp.Right = maexp
					lastBexp.Group = false
				} else if maexp == nil {
					maexp = &core.MemberAccessExpression{
						Object:             tmpExp,
						PropertyExpression: aexp.Elements[0],
						Compute:            true,
						Line:               tmpExp.GetLine(),
						CharAt:             tmpExp.GetCharAt(),
					}
				} else if maexp.PropertyExpression != nil || maexp.PropertyIdentifier.Name != "" {
					maexp = &core.MemberAccessExpression{
						Object:             maexp,
						PropertyExpression: aexp.Elements[0],
						Compute:            true,
						Line:               maexp.Line,
						CharAt:             maexp.CharAt,
					}
				} else {
					return nil, 0, fmt.Errorf("TODO")
				}
				tmpExp = maexp
			} else {
				tmpExp = exp
			}
			if maexp != nil && maexp.PropertyExpression == nil && maexp.PropertyIdentifier.Name == "" {
				vexp, _ := tmpExp.(*core.VariableExpression)
				if vexp == nil {
					return nil, 0, fmt.Errorf("TODO")
				}
				maexp.PropertyIdentifier = core.Identifier{
					Name:   vexp.Name,
					Line:   vexp.Line,
					CharAt: vexp.CharAt,
				}
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
				if lbexp, ok := bexp.Left.(*core.BinaryExpression); ok && lbexp.Operator.Compare(bexp.Operator) > 0 && !lbexp.Group && (lastBexpAfterGroup == nil || bexp != lastBexpAfterGroup) {
					bexp.Left = lbexp.Right
					bexp.CharAt = lbexp.Right.GetCharAt()
					lbexp.Right = bexp
					bexp = lbexp
				}
			}
		} else if t.Value == "." {
			if lastBexp != nil {
				maexp = &core.MemberAccessExpression{
					Object: lastBexp.Right,
					Line:   lastBexp.Right.GetLine(),
					CharAt: lastBexp.Right.GetCharAt(),
				}
				lastBexp.Right = maexp
				lastBexp.Group = false
			} else if maexp == nil {
				if tmpExp == nil {
					return nil, 0, fmt.Errorf("TODO")
				}
				maexp = &core.MemberAccessExpression{
					Object: tmpExp,
					Line:   tmpExp.GetLine(),
					CharAt: tmpExp.GetCharAt(),
				}
			} else if maexp.PropertyExpression != nil || maexp.PropertyIdentifier.Name != "" {
				maexp = &core.MemberAccessExpression{
					Object: maexp,
					Line:   maexp.Line,
					CharAt: maexp.CharAt,
				}
			} else {
				return nil, 0, fmt.Errorf("TODO")
			}
			tmpExp = nil
		} else if t.IsOperatorSymbol() {
			op := core.Operator{
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
				bexp = &core.BinaryExpression{
					Left:     tmpExp,
					Line:     t.Line, // tmpExp.getCharAt()
					CharAt:   tmpExp.GetCharAt(),
					Operator: op,
				}
			} else if bexp.Right != nil {
				bexp = &core.BinaryExpression{
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
			openParen++
			if tmpExp != nil {
				args, processed, err := parseSequentExpressions(tokens[i+1:])
				if err != nil {
					return nil, 0, err
				}
				i = i + processed
				if i+1 >= len(tokens) || tokens[i+1].Value != ")" {
					return nil, 0, fmt.Errorf("TODO")
				}
				if lastBexp != nil {
					tmpExp = &core.CallExpression{
						Callee:    lastBexp.Right,
						Arguments: args,
						Line:      lastBexp.Right.GetLine(),
						CharAt:    lastBexp.Right.GetCharAt(),
					}
					lastBexp.Right = tmpExp
					lastBexp.Group = false
				} else {
					tmpExp = &core.CallExpression{
						Callee:    tmpExp,
						Arguments: args,
						Line:      tmpExp.GetLine(),
						CharAt:    tmpExp.GetCharAt(),
					}
				}
			}
		} else if t.Value == ")" {
			if openParen == 0 && len(bexpsAfterGroup) == 0 {
				// it means this ')' doesn't belong to current expreesion
				break
			}
			if openParen > 0 && openParen >= len(bexpsAfterGroup) {
				// if openParen is increased by some exp like (a)
				openParen--
			} else if bexp != nil {
				if len(bexpsAfterGroup) > 0 {
					// Swap expression order to match expression group
					lastBexpAfterGroup, i := getLastBexpAfterGroup(bexpsAfterGroup, true)
					if bexpBeforeGroup, ok := lastBexpAfterGroup.Left.(*core.BinaryExpression); ok {
						lastBexpAfterGroup.Left = bexpBeforeGroup.Right
						lastBexpAfterGroup.CharAt = bexpBeforeGroup.Right.GetCharAt()
						bexpBeforeGroup.Right = bexp
						groupedBexp := bexp
						bexp = bexpBeforeGroup

						// loop from groupedBexp to lastBexpAfterGroup and update exp CharAt
						groupedBexp.CharAt = lastBexpAfterGroup.CharAt
						for tmpBexp, _ := groupedBexp.Left.(*core.BinaryExpression); tmpBexp != nil; {
							tmpBexp.CharAt = lastBexpAfterGroup.CharAt
							tmpBexp, _ = tmpBexp.Left.(*core.BinaryExpression)
						}
					}
					bexpsAfterGroup = append(bexpsAfterGroup[i+1:], bexpsAfterGroup[:i]...)
				}
			}
		} else if (bexp == nil || bexp.Right == nil) && tmpExp == nil {
			return nil, 0, err
		} else {
			break
		}
	}
	if bexp != nil && bexp.Right != nil {
		return bexp, i, nil
	}
	return tmpExp, i, nil
}

func getLastBexpAfterGroup(bexps []*core.BinaryExpression, rightShouldNotNil bool) (*core.BinaryExpression, int) {
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

// Check first token, if it is the start of an expression then parse it and return
func parseTempExpression(tokens []lexer.Token) (core.Expression, int, error) {
	if len(tokens) == 0 {
		return nil, 0, fmt.Errorf("TODO")
	}
	t := tokens[0]
	if ptype, ok := t.ParsePrimitiveType(); ok {
		v := t.Value
		if ptype == core.LiteralTypeString {
			v = t.Value[1 : len(t.Value)-1]
		}
		return &core.LiteralExpression{
			Type:   ptype,
			Value:  v,
			Line:   t.Line,
			CharAt: t.CharAt,
		}, 1, nil
	}
	if t.IsIdentifier() {
		return &core.VariableExpression{
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
	if t.Value == "func" {
		return parseFunctionExpression(tokens)
	}
	return nil, 0, fmt.Errorf("TODO")
}

func parseSequentIdentifiers(tokens []lexer.Token) ([]core.Identifier, int, error) {
	ids := []core.Identifier{}
	var i int
	for i = 0; i < len(tokens); i++ {
		t := tokens[i]
		if !t.IsIdentifier() {
			return ids, i, nil
		}
		if t.IsIdentifier() {
			ids = append(ids, core.Identifier{
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

func parseSequentExpressions(tokens []lexer.Token) ([]core.Expression, int, error) {
	exps := []core.Expression{}
	var i int
	for i = 0; i < len(tokens); i++ {
		exp, processed, _ := parseExpression(tokens[i:])
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

// TODO: method, shorthand
func parseObjectExpression(tokens []lexer.Token) (core.Expression, int, error) {
	if len(tokens) == 0 || tokens[0].Value != "{" {
		return nil, 0, fmt.Errorf("TODO")
	}
	obj := &core.ObjectExpression{
		Properties: []*core.ObjectProperty{},
		Line:       tokens[0].Line,
		CharAt:     tokens[0].CharAt,
	}
	var prop *core.ObjectProperty
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
			prop = &core.ObjectProperty{
				Line:   t.Line,
				CharAt: t.CharAt,
			}
			if t.IsIdentifier() {
				prop.KeyIdentifier = core.Identifier{
					Name:   t.Value,
					Line:   t.Line,
					CharAt: t.CharAt,
				}
			} else if t.Value == "[" {
				exp, processed, err := parseArrayExpression(tokens[i:])
				if err != nil {
					return nil, 0, err
				}
				aexp, _ := exp.(*core.ArrayExpression)
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
			obj.Properties = append(obj.Properties, prop)
			prop = nil
		} else {
			return nil, 0, fmt.Errorf("TODO")
		}
	}
	return obj, i + 1, nil
}

func parseArrayExpression(tokens []lexer.Token) (core.Expression, int, error) {
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
	return &core.ArrayExpression{
		Elements: exps,
		Line:     tokens[0].Line,
		CharAt:   tokens[0].CharAt,
	}, processed + 2, nil
}

func parseFunctionExpression(tokens []lexer.Token) (*core.FunctionExpression, int, error) {
	if len(tokens) < 5 { // 6 is len of the most simple function
		return nil, 0, fmt.Errorf("TODO")
	}
	if tokens[1].Value != "(" {
		return nil, 0, fmt.Errorf("TODO")
	}
	i := 1
	f := &core.FunctionExpression{
		Line:   tokens[0].Line,
		CharAt: tokens[0].CharAt,
	}
	params, statements, processed, err := parseFunctionParamAndBody(tokens[i:])
	if err != nil {
		return nil, 0, err
	}
	f.Params = params
	f.Body = statements
	return f, i + processed, nil
}

func parseFunctionParamAndBody(tokens []lexer.Token) ([]core.Identifier, []core.Statement, int, error) {
	if len(tokens) < 4 { // (){} -> min len = 4
		return nil, nil, 0, fmt.Errorf("TODO")
	}
	if tokens[0].Value != "(" {
		return nil, nil, 0, fmt.Errorf("TODO")
	}
	i := 0
	params, processed, err := parseSequentIdentifiers(tokens[i+1:])
	i = i + processed + 1
	if err != nil {
		return nil, nil, 0, err
	}
	if tokens[i].Value != ")" {
		return nil, nil, 0, fmt.Errorf("TODO")
	}
	i++
	if tokens[i].Value != "{" {
		return nil, nil, 0, fmt.Errorf("TODO")
	}
	statements, processed, err := parseStatements(tokens[i+1:])
	i = i + processed + 1
	if err != nil {
		return nil, nil, 0, err
	}
	if tokens[i].Value != "}" {
		return nil, nil, 0, fmt.Errorf("TODO")
	}
	return params, statements, i + 1, nil
}
