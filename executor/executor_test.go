package executor

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gs/core"
	"gs/lexer"
	"gs/parser"
)

func TestExecute(t *testing.T) {
	cases := []struct {
		name   string
		in     string
		inEC   *core.ExecutionContext
		wantEC func() *core.ExecutionContext
		err    error
	}{
		{
			name: "execute variable declaration #1",
			in:   "var a = 1",
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 9,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute variable declaration #2",
			in: `var a = 1
				 var a = 2`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   2,
							CharAt: 9,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute variable declaration #3",
			in: `var a = 1
				 var b = 2`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 9,
						},
						"b": &core.LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   2,
							CharAt: 9,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute variable declaration #4",
			in: `var a,b = 1,"2"
				 var c = false`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 11,
						},
						"b": &core.LiteralExpression{
							Type:   "string",
							Value:  "2",
							Line:   1,
							CharAt: 13,
						},
						"c": &core.LiteralExpression{
							Type:   "boolean",
							Value:  "false",
							Line:   2,
							CharAt: 9,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute variable declaration #5",
			in:   `var a,b = 1`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 11,
						},
						"b": &core.LiteralExpression{
							Type:   "undefined",
							Line:   1,
							CharAt: 7,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute variable declaration #5",
			in: `var a = {
				   b: {
				     c: 1,
				   },
				 }
				 var d = a.b.c`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.ObjectExpression{
							Properties: []core.ObjectProperty{
								{
									KeyIdentifier: core.Identifier{
										Name:   "b",
										Line:   2,
										CharAt: 1,
									},
									Value: &core.ObjectExpression{
										Properties: []core.ObjectProperty{
											{
												KeyIdentifier: core.Identifier{
													Name:   "c",
													Line:   3,
													CharAt: 1,
												},
												Value: &core.LiteralExpression{
													Type:   "number",
													Value:  "1",
													Line:   3,
													CharAt: 4,
												},
												Line:   3,
												CharAt: 1,
											},
										},
										Line:   2,
										CharAt: 4,
									},
									Line:   2,
									CharAt: 1,
								},
							},
							Line:   6,
							CharAt: 9,
						},
						"d": &core.LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   3,
							CharAt: 4,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute variable declaration #6",
			in: `var a = [
				   "1"+"2",
				   1+2,
				   1+false,
				 ]`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.ArrayExpression{
							Elements: []core.Expression{
								&core.LiteralExpression{
									Type:   "string",
									Value:  "12",
									Line:   2,
									CharAt: 1,
								},
								&core.LiteralExpression{
									Type:   "number",
									Value:  "3",
									Line:   3,
									CharAt: 1,
								},
								&core.LiteralExpression{
									Type:   "string",
									Value:  "1false",
									Line:   4,
									CharAt: 1,
								},
							},
							Line:   1,
							CharAt: 9,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute variable declaration #7",
			in:   `var a = func(){}`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				gec := &core.ExecutionContext{
					Variables: map[string]core.Expression{},
				}
				gec.Variables["a"] = &core.FunctionExpression{
					Params: []core.Identifier{},
					Body:   []core.Statement{},
					EC: &core.ExecutionContext{
						Outer:     gec,
						Variables: map[string]core.Expression{},
					},
					Line:   1,
					CharAt: 9,
				}
				return gec
			},
			err: nil,
		},
		{
			name: "execute variable declaration #8",
			in: `
			var a = func(b,c){
			  return b+c
			}
			var d = a(1,2)`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				gec := &core.ExecutionContext{
					Variables: map[string]core.Expression{},
				}
				gec.Variables["a"] = &core.FunctionExpression{
					Params: []core.Identifier{
						{Name: "b", Line: 2, CharAt: 14},
						{Name: "c", Line: 2, CharAt: 16},
					},
					Body: []core.Statement{
						core.ReturnStatement{
							Argument: &core.BinaryExpression{
								Left: &core.VariableExpression{
									Name:   "b",
									Line:   3,
									CharAt: 8,
								},
								Right: &core.VariableExpression{
									Name:   "c",
									Line:   3,
									CharAt: 10,
								},
								Operator: core.Operator{
									Symbol: "+",
									Line:   3,
									CharAt: 9,
								},
								Line:   3,
								CharAt: 8,
							},
							Line:   3,
							CharAt: 1,
						},
					},
					EC: &core.ExecutionContext{
						Outer: gec,
						Variables: map[string]core.Expression{
							"b": &core.LiteralExpression{
								Type:   "number",
								Value:  "1",
								Line:   3,
								CharAt: 8,
							},
							"c": &core.LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   3,
								CharAt: 10,
							},
						},
					},
					Line:   5,
					CharAt: 9,
				}
				gec.Variables["d"] = &core.LiteralExpression{
					Type:   "number",
					Value:  "3",
					Line:   3,
					CharAt: 8,
				}
				return gec
			},
			err: nil,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			stmts, _ := parser.ToAST(lexer.Lex(tt.in))
			require.Equal(t, tt.err, Execute(tt.inEC, stmts))
			require.Equal(t, tt.wantEC(), tt.inEC)
		})
	}
}

func TestTMP(t *testing.T) {
	cases := []struct {
		name   string
		in     string
		inEC   *core.ExecutionContext
		wantEC func() *core.ExecutionContext
		err    error
	}{
		{
			name: "execute variable declaration #8",
			in: `
			var a = func(b,c){
			  return b+c
			}
			var d = a(1,2)`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				gec := &core.ExecutionContext{
					Variables: map[string]core.Expression{},
				}
				gec.Variables["a"] = &core.FunctionExpression{
					Params: []core.Identifier{
						{Name: "b", Line: 2, CharAt: 14},
						{Name: "c", Line: 2, CharAt: 16},
					},
					Body: []core.Statement{
						core.ReturnStatement{
							Argument: &core.BinaryExpression{
								Left: &core.VariableExpression{
									Name:   "b",
									Line:   3,
									CharAt: 8,
								},
								Right: &core.VariableExpression{
									Name:   "c",
									Line:   3,
									CharAt: 10,
								},
								Operator: core.Operator{
									Symbol: "+",
									Line:   3,
									CharAt: 9,
								},
								Line:   3,
								CharAt: 8,
							},
							Line:   3,
							CharAt: 1,
						},
					},
					EC: &core.ExecutionContext{
						Outer: gec,
						Variables: map[string]core.Expression{
							"b": &core.LiteralExpression{
								Type:   "number",
								Value:  "1",
								Line:   3,
								CharAt: 8,
							},
							"c": &core.LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   3,
								CharAt: 10,
							},
						},
					},
					Line:   5,
					CharAt: 9,
				}
				gec.Variables["d"] = &core.LiteralExpression{
					Type:   "number",
					Value:  "3",
					Line:   3,
					CharAt: 8,
				}
				return gec
			},
			err: nil,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			stmts, _ := parser.ToAST(lexer.Lex(tt.in))
			require.Equal(t, tt.err, Execute(tt.inEC, stmts))
			require.Equal(t, tt.wantEC(), tt.inEC)
		})
	}
}
