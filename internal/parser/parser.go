package parser

import (
	"fmt"

	"github.com/dhl1402/covidscript/internal/core"
	"github.com/dhl1402/covidscript/internal/lexer"
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
		case t.Value == "break":
			ss = append(ss, core.BreakStatement{Line: t.Line, CharAt: t.CharAt})
		case t.Value == "continue":
			ss = append(ss, core.ContinueStatement{Line: t.Line, CharAt: t.CharAt})
		case t.Value == "{":
			s, processed, err := parseBlockStatement(tokens[i:])
			if err != nil {
				return nil, 0, err
			}
			ss = append(ss, *s)
			i = i + processed - 1
		case t.Value == "}":
			return ss, i + 1, nil
		case t.Value == "if":
			s, processed, err := parseIfStatement(tokens[i:])
			if err != nil {
				return nil, 0, err
			}
			ss = append(ss, *s)
			i = i + processed - 1
		case t.Value == "for":
			s, processed, err := parseForStatement(tokens[i:])
			if err != nil {
				return nil, 0, err
			}
			ss = append(ss, *s)
			i = i + processed - 1
		case t.Value == ";":
			continue
		default:
			s, processed, err := parseAssignmentStatement(tokens[i:])
			if err == nil {
				ss = append(ss, *s)
				i = i + processed - 1
			} else {
				s, processed, err := parseExpressionStatement(tokens[i:])
				if err != nil {
					return nil, 0, err
				}
				ss = append(ss, *s)
				i = i + processed - 1
			}
		}
	}
	return ss, i, nil
}

func parseBlockStatement(tokens []lexer.Token) (*core.BlockStatement, int, error) {
	if len(tokens) < 2 {
		return nil, 0, fmt.Errorf("Parsing error: cannot parse block statement")
	}
	if tokens[0].Value != "{" {
		return nil, 0, fmt.Errorf("Parsing error: unexpected token '%s', expected '{'. [%d,%d]", tokens[0].Value, tokens[0].Line, tokens[0].CharAt)
	}
	stmts, i, err := parseStatements(tokens[1:])
	if err != nil {
		return nil, 0, err
	}
	if tokens[i].Value != "}" {
		return nil, 0, fmt.Errorf("Parsing error: unexpected token '%s', expected '}'. [%d,%d]", tokens[i].Value, tokens[i].Line, tokens[i].CharAt)
	}
	return &core.BlockStatement{
		Statements: stmts,
		Line:       tokens[0].Line,
		CharAt:     tokens[0].CharAt,
	}, i + 1, nil
}

func parseAssignmentStatement(tokens []lexer.Token) (*core.AssignmentStatement, int, error) {
	if len(tokens) == 0 {
		return nil, 0, fmt.Errorf("Parsing error: cannot parse assignment statement")
	}
	exp, i, err := parseExpression(tokens[0:])
	if err != nil {
		return nil, 0, err
	}
	if i >= len(tokens) {
		lastToken := tokens[len(tokens)-1]
		return nil, 0, fmt.Errorf("Parsing error: unexpected end of statement. [%d,%d]", lastToken.Line, lastToken.CharAt)
	}
	if tokens[i].Value != "=" && tokens[i].Value != ":=" {
		return nil, 0, fmt.Errorf("Parsing error: unexpected token '%s', expected '=' or ':='. [%d,%d]", tokens[i].Value, tokens[i].Line, tokens[i].CharAt)
	}
	switch exp.(type) {
	case *core.VariableExpression:
	case *core.MemberAccessExpression:
	default:
		return nil, 0, fmt.Errorf("Parsing error: %s cannot be the left side of assignment statement. [%d,%d]", exp.GetType(), exp.GetLine(), exp.GetCharAt())
	}
	as := &core.AssignmentStatement{
		Left:   exp,
		Line:   tokens[0].Line,
		CharAt: tokens[0].CharAt,
	}
	if tokens[i].Value == ":=" {
		as.DeclarationShorthand = true
	}
	i++ // handle '=' -> +1
	rightExp, processed, err := parseExpression(tokens[i:])
	if err != nil {
		return nil, 0, err
	}
	as.Right = rightExp
	return as, i + processed, nil
}

func parseExpressionStatement(tokens []lexer.Token) (*core.ExpressionStatement, int, error) {
	if len(tokens) == 0 {
		return nil, 0, fmt.Errorf("Parsing error: cannot parse expression statement")
	}
	exp, i, err := parseExpression(tokens[0:])
	if err != nil {
		return nil, 0, err
	}
	return &core.ExpressionStatement{
		Expression: exp,
		Line:       exp.GetLine(),
		CharAt:     exp.GetCharAt(),
	}, i, nil
}

func parseFunctionDeclaration(tokens []lexer.Token) (*core.FunctionDeclaration, int, error) {
	if len(tokens) < 6 { // 6 is len of the most simple function
		return nil, 0, fmt.Errorf("Parsing error: cannot parse function declaration")
	}
	if !tokens[1].IsIdentifier() {
		return nil, 0, fmt.Errorf("Parsing error: %s is not a valid variable name. [%d,%d]", tokens[1].Value, tokens[1].Line, tokens[1].CharAt)
	}
	if tokens[2].Value != "(" {
		return nil, 0, fmt.Errorf("Parsing error: unexpected token '%s', expected '('. [%d,%d]", tokens[2].Value, tokens[2].Line, tokens[2].CharAt)
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
		return nil, 0, fmt.Errorf("Parsing error: cannot parse return statement")
	}
	r := &core.ReturnStatement{
		Line:   tokens[0].Line,
		CharAt: tokens[0].CharAt,
	}
	exp, i, _ := parseExpression(tokens[1:]) // skip 'return'
	if exp != nil {
		r.Argument = exp
	}
	return r, i + 1, nil
}

func parseVariableDeclaration(tokens []lexer.Token) (*core.VariableDeclaration, int, error) {
	if len(tokens) == 0 {
		return nil, 0, fmt.Errorf("Parsing error: cannot parse variable declaration")
	}
	ids, i, err := parseSequentIdentifiers(tokens[1:]) // tokens[1:] -> skip 'var'
	if err != nil {
		return nil, 0, err
	}
	if len(ids) == 0 {
		return nil, 0, fmt.Errorf("Parsing error: cannot parse variable names. [%d,%d]", tokens[0].Line, tokens[0].CharAt)
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
		return nil, 0, fmt.Errorf("Parsing error: cannot parse variable initialization. [%d,%d]", tokens[i+1].Line, tokens[i+1].CharAt)
	}
	if len(exps) > len(s.Declarations) {
		return nil, 0, fmt.Errorf("Parsing error: too many expressions. [%d,%d]", exps[len(exps)-1].GetLine(), exps[len(exps)-1].GetCharAt())
	}
	for ii, exp := range exps {
		s.Declarations[ii].Init = exp
	}
	return s, i + processed + 2, nil
}

func parseIfStatement(tokens []lexer.Token) (*core.IfStatement, int, error) {
	if len(tokens) < 4 { // 4 is len of the most simple if
		return nil, 0, fmt.Errorf("Parsing error: cannot parse if statement")
	}
	ifstmt := &core.IfStatement{
		Line:   tokens[0].Line,
		CharAt: tokens[0].CharAt,
	}
	i := 1 // skip 'if'
	astmt, processed, err := parseAssignmentStatement(tokens[i:])
	if err == nil {
		ifstmt.Assignment = astmt
		i = i + processed
		if i >= len(tokens) {
			lastToken := tokens[len(tokens)-1]
			return nil, 0, fmt.Errorf("Parsing error: unexpected end of statement. [%d,%d]", lastToken.Line, lastToken.CharAt)
		}
		if tokens[i].Value != ";" {
			return nil, 0, fmt.Errorf("Parsing error: unexpected token '%s', expected ';'. [%d,%d]", tokens[i].Value, tokens[i].Line, tokens[i].CharAt)
		}
		i++ // skip ';'
	}
	estmt, processed, err := parseExpressionStatement(tokens[i:])
	if err != nil {
		return nil, 0, err
	}
	i = i + processed
	ifstmt.Test = estmt.Expression

	// parse if body block
	bstmt, processed, err := parseBlockStatement(tokens[i:])
	if err != nil {
		return nil, 0, err
	}
	i = i + processed
	ifstmt.Consequent = *bstmt

	// parse elif
	if i < len(tokens) && tokens[i].Value == "elif" {
		elifstmt, processed, err := parseIfStatement(tokens[i:])
		if err != nil {
			return nil, 0, err
		}
		ifstmt.Alternate = elifstmt
		i = i + processed
	}
	// parse else
	if i < len(tokens) && tokens[i].Value == "else" {
		bstmt, processed, err := parseBlockStatement(tokens[i+1:])
		if err != nil {
			return nil, 0, err
		}
		elseStmt := &core.IfStatement{
			Consequent: *bstmt,
			Line:       tokens[i].Line,
			CharAt:     tokens[i].CharAt,
		}
		ifstmt.Alternate = elseStmt
		i = i + processed + 1
	}
	return ifstmt, i, nil
}

func parseForStatement(tokens []lexer.Token) (*core.ForStatement, int, error) {
	if len(tokens) < 3 { // 3 is len of the most simple for
		return nil, 0, fmt.Errorf("Parsing error: cannot parse for statement")
	}
	forstmt := &core.ForStatement{
		Line:   tokens[0].Line,
		CharAt: tokens[0].CharAt,
	}
	i := 1 // skip 'for'
	if tokens[i].Value == "{" {
		bstmt, processed, err := parseBlockStatement(tokens[i:])
		if err != nil {
			return nil, 0, err
		}
		forstmt.Body = *bstmt
		return forstmt, i + processed, nil
	}
	astmt, processed, err := parseAssignmentStatement(tokens[i:])
	if err == nil {
		forstmt.Assignment = astmt
		i = i + processed
	}
	if i < len(tokens) && tokens[i].Value == ";" {
		i++ // skip ';'
	}
	if i < len(tokens) && tokens[i].Value == "{" {
		return nil, 0, fmt.Errorf("Parsing error: unexpected token '%s'. [%d,%d]", tokens[i].Value, tokens[i].Line, tokens[i].CharAt)
	}
	estmt, processed, err := parseExpressionStatement(tokens[i:])
	if err == nil {
		forstmt.Test = estmt.Expression
		i = i + processed
	}
	if i < len(tokens) && tokens[i].Value == ";" {
		i++ // skip ';'
	}
	ustmt, processed, err := parseAssignmentStatement(tokens[i:])
	if err == nil {
		forstmt.Update = ustmt
		i = i + processed
	}
	bstmt, processed, err := parseBlockStatement(tokens[i:])
	if err != nil {
		return nil, 0, err
	}
	i = i + processed
	forstmt.Body = *bstmt
	return forstmt, i, nil
}

func parseExpression(tokens []lexer.Token) (core.Expression, int, error) {
	if len(tokens) == 0 {
		return nil, 0, fmt.Errorf("Parsing error: cannot parse expression")
	}
	var tmpExp core.Expression
	var bexp *core.BinaryExpression
	var lastBexp *core.BinaryExpression
	var maexp *core.MemberAccessExpression
	bexpsAfterGroup := []*core.BinaryExpression{}
	openParen := 0
	unaryQueue := [][]lexer.Token{{}} // by paren level
	var i int
	for i = 0; i < len(tokens); i++ {
		t := tokens[i]
		parenLevel := openParen + len(bexpsAfterGroup)
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
					tt := tokens[i-processed+1]
					return nil, 0, fmt.Errorf("Parsing error: unexpected expression. [%d,%d]", tt.Line, tt.CharAt)
				}
				tmpExp = maexp
			} else if len(unaryQueue[parenLevel]) > 0 {
				tmpExp = wrapInUnaryExpression(exp, unaryQueue, parenLevel)
			} else {
				tmpExp = exp
			}
			if maexp != nil && maexp.PropertyExpression == nil && maexp.PropertyIdentifier.Name == "" {
				vexp, _ := tmpExp.(*core.VariableExpression)
				if vexp == nil {
					return nil, 0, fmt.Errorf("Parsing error: unexpected token '%s'. [%d,%d]", t.Value, t.Line, t.CharAt)
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
				tmpExp = bexp
			}
		} else if t.Value == "!" {
			unaryQueue[parenLevel] = append(unaryQueue[parenLevel], t)
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
					return nil, 0, fmt.Errorf("Parsing error: unexpected token '%s'. [%d,%d]", t.Value, t.Line, t.CharAt)
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
				return nil, 0, fmt.Errorf("Parsing error: unexpected token '%s'. [%d,%d]", t.Value, t.Line, t.CharAt)
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
					return nil, 0, fmt.Errorf("Parsing error: unexpected token '%s'. [%d,%d]", t.Value, t.Line, t.CharAt)
				}
				if tmpExp == maexp {
					maexp = nil
				}
				bexp = &core.BinaryExpression{
					Left:     tmpExp,
					Line:     tmpExp.GetLine(),
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
				return nil, 0, fmt.Errorf("Parsing error: unexpected token '%s'. [%d,%d]", t.Value, t.Line, t.CharAt)
			}
			if bexp != nil && openParen > 0 {
				bexpsAfterGroup = append(bexpsAfterGroup, bexp)
				openParen--
			}
			tmpExp = nil
		} else if t.Value == "(" {
			openParen++
			if openParen+len(bexpsAfterGroup) >= len(unaryQueue) {
				unaryQueue = append(unaryQueue, []lexer.Token{}) // increase level
			}
			if tmpExp != nil {
				args, processed, err := parseSequentExpressions(tokens[i+1:])
				if err != nil {
					return nil, 0, err
				}
				i = i + processed
				if i+1 >= len(tokens) {
					lastToken := tokens[len(tokens)-1]
					return nil, 0, fmt.Errorf("Parsing error: unexpected end of expression. [%d,%d]", lastToken.Line, lastToken.CharAt)
				}
				if tokens[i+1].Value != ")" {
					return nil, 0, fmt.Errorf("Parsing error: unexpected token '%s', expected ')'. [%d,%d]", tokens[i+1].Value, tokens[i+1].Line, tokens[i+1].CharAt)
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
					tmpExp = lastBexp
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
			parenLevel--
			if openParen == 0 && len(bexpsAfterGroup) == 0 {
				// it means this ')' doesn't belong to current expreesion
				break
			}
			if openParen > 0 && openParen >= len(bexpsAfterGroup) {
				// if openParen is increased by some exp like (a) or ((a+b)+c) -> do not need to swap
				if len(unaryQueue[parenLevel]) > 0 {
					if openParen == len(bexpsAfterGroup) && bexp != nil {
						tmpExp = wrapInUnaryExpression(bexp, unaryQueue, parenLevel).(*core.UnaryExpression)
						bexp = nil
					} else if bexp == nil {
						tmpExp = wrapInUnaryExpression(tmpExp, unaryQueue, parenLevel)
					} else {
						bexp.Right = wrapInUnaryExpression(bexp.Right, unaryQueue, parenLevel)
					}
				}
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
						if len(unaryQueue[parenLevel]) > 0 {
							bexp.Right = wrapInUnaryExpression(bexp.Right, unaryQueue, parenLevel)
						}

						// loop from groupedBexp to lastBexpAfterGroup and update exp CharAt
						groupedBexp.CharAt = lastBexpAfterGroup.CharAt
						for tmpBexp, _ := groupedBexp.Left.(*core.BinaryExpression); tmpBexp != nil; {
							tmpBexp.CharAt = lastBexpAfterGroup.CharAt
							tmpBexp, _ = tmpBexp.Left.(*core.BinaryExpression)
						}
					} else if len(unaryQueue[parenLevel]) > 0 {
						tmpExp, _ = wrapInUnaryExpression(bexp, unaryQueue, parenLevel).(*core.UnaryExpression)
						bexp = nil
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

func wrapInUnaryExpression(exp core.Expression, unaryQueue [][]lexer.Token, parenLevel int) core.Expression {
	for j := len(unaryQueue[parenLevel]) - 1; j >= 0; j-- {
		exp = &core.UnaryExpression{
			Expression: exp,
			Line:       unaryQueue[parenLevel][j].Line,
			CharAt:     unaryQueue[parenLevel][j].CharAt,
		}
	}
	unaryQueue[parenLevel] = []lexer.Token{}
	return exp
}

// Check first token, if it is the start of an expression then parse it and return
func parseTempExpression(tokens []lexer.Token) (core.Expression, int, error) {
	if len(tokens) == 0 {
		return nil, 0, fmt.Errorf("Parsing error: cannot parse expression")
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
	return nil, 0, fmt.Errorf("Parsing error: cannot parse expression")
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

func parseObjectExpression(tokens []lexer.Token) (core.Expression, int, error) {
	if len(tokens) == 0 || tokens[0].Value != "{" {
		return nil, 0, fmt.Errorf("Parsing error: cannot parse object")
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
			if prop != nil || (!nt.IsIdentifier() && nt.Value != "[") {
				return nil, 0, fmt.Errorf("Parsing error: unexpected token '%s'. [%d,%d]", t.Value, t.Line, t.CharAt)
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
					return nil, 0, fmt.Errorf("Parsing error: cannot parse object property key. [%d,%d]", t.Line, t.CharAt)
				}
				prop.KeyExpression = aexp.Elements[0]
				prop.Computed = true
				i = i + processed - 1
			} else {
				return nil, 0, fmt.Errorf("Parsing error: unexpected token '%s'. [%d,%d]", t.Value, t.Line, t.CharAt)
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
			return nil, 0, fmt.Errorf("Parsing error: unexpected token '%s'. [%d,%d]", t.Value, t.Line, t.CharAt)
		}
	}
	return obj, i + 1, nil
}

func parseArrayExpression(tokens []lexer.Token) (core.Expression, int, error) {
	if len(tokens) == 0 || tokens[0].Value != "[" {
		return nil, 0, fmt.Errorf("Parsing error: cannot parse array expession")
	}
	exps, processed, err := parseSequentExpressions(tokens[1:]) // skip '['
	if err != nil {
		return nil, 0, err
	}
	if processed+1 >= len(tokens) {
		lastToken := tokens[len(tokens)-1]
		return nil, 0, fmt.Errorf("Parsing error: missing token ']'. [%d,%d]", lastToken.Line, lastToken.CharAt)
	}
	if tokens[processed+1].Value != "]" {
		return nil, 0, fmt.Errorf("Parsing error: unexpected token '%s', expected ']. [%d,%d]", tokens[processed+1].Value, tokens[processed+1].Line, tokens[processed+1].CharAt)
	}
	return &core.ArrayExpression{
		Elements: exps,
		Line:     tokens[0].Line,
		CharAt:   tokens[0].CharAt,
	}, processed + 2, nil
}

func parseFunctionExpression(tokens []lexer.Token) (*core.FunctionExpression, int, error) {
	if len(tokens) < 5 { // 6 is len of the most simple function
		return nil, 0, fmt.Errorf("Parsing error: cannot parse function expression")
	}
	if tokens[1].Value != "(" {
		return nil, 0, fmt.Errorf("Parsing error: unexpected token '%s', expected '('. [%d,%d]", tokens[1].Value, tokens[1].Line, tokens[1].CharAt)
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
		return nil, nil, 0, fmt.Errorf("Parsing error: cannot parse function")
	}
	if tokens[0].Value != "(" {
		return nil, nil, 0, fmt.Errorf("Parsing error: unexpected token '%s', expected '('. [%d,%d]", tokens[0].Value, tokens[0].Line, tokens[0].CharAt)
	}
	i := 0
	params, processed, err := parseSequentIdentifiers(tokens[i+1:])
	i = i + processed + 1
	if err != nil {
		return nil, nil, 0, err
	}
	if tokens[i].Value != ")" {
		return nil, nil, 0, fmt.Errorf("Parsing error: unexpected token '%s', expected ')'. [%d,%d]", tokens[i].Value, tokens[i].Line, tokens[i].CharAt)
	}
	i++
	if tokens[i].Value != "{" {
		return nil, nil, 0, fmt.Errorf("Parsing error: unexpected token '%s', expected '{'. [%d,%d]", tokens[i].Value, tokens[i].Line, tokens[i].CharAt)
	}
	statements, processed, err := parseStatements(tokens[i+1:])
	i = i + processed
	if err != nil {
		return nil, nil, 0, err
	}
	if tokens[i].Value != "}" {
		return nil, nil, 0, fmt.Errorf("Parsing error: unexpected token '%s', expected '}'. [%d,%d]", tokens[i].Value, tokens[i].Line, tokens[i].CharAt)
	}
	return params, statements, i + 1, nil
}
