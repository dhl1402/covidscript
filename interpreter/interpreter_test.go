package interpreter

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dhl1402/covidscript/core"
	"github.com/dhl1402/covidscript/lexer"
	"github.com/dhl1402/covidscript/parser"
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
							Type:   core.LiteralTypeNumber,
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
							Type:   core.LiteralTypeNumber,
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
							Type:   core.LiteralTypeNumber,
							Value:  "1",
							Line:   1,
							CharAt: 9,
						},
						"b": &core.LiteralExpression{
							Type:   core.LiteralTypeNumber,
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
							Type:   core.LiteralTypeNumber,
							Value:  "1",
							Line:   1,
							CharAt: 11,
						},
						"b": &core.LiteralExpression{
							Type:   core.LiteralTypeString,
							Value:  "2",
							Line:   1,
							CharAt: 13,
						},
						"c": &core.LiteralExpression{
							Type:   core.LiteralTypeBoolean,
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
							Type:   core.LiteralTypeNumber,
							Value:  "1",
							Line:   1,
							CharAt: 11,
						},
						"b": &core.LiteralExpression{
							Type:   core.LiteralTypeUndefined,
							Line:   1,
							CharAt: 7,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute variable declaration #6",
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
							Properties: []*core.ObjectProperty{
								{
									KeyIdentifier: core.Identifier{
										Name:   "b",
										Line:   2,
										CharAt: 1,
									},
									Value: &core.ObjectExpression{
										Properties: []*core.ObjectProperty{
											{
												KeyIdentifier: core.Identifier{
													Name:   "c",
													Line:   3,
													CharAt: 1,
												},
												Value: &core.LiteralExpression{
													Type:   core.LiteralTypeNumber,
													Value:  "1",
													Line:   6,
													CharAt: 9,
												},
												Line:   3,
												CharAt: 1,
											},
										},
										Line:   6,
										CharAt: 9,
									},
									Line:   2,
									CharAt: 1,
								},
							},
							Line:   6,
							CharAt: 9,
						},
						"d": &core.LiteralExpression{
							Type:   core.LiteralTypeNumber,
							Value:  "1",
							Line:   6,
							CharAt: 9,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute variable declaration #7",
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
									Type:   core.LiteralTypeString,
									Value:  "12",
									Line:   2,
									CharAt: 1,
								},
								&core.LiteralExpression{
									Type:   core.LiteralTypeNumber,
									Value:  "3",
									Line:   3,
									CharAt: 1,
								},
								&core.LiteralExpression{
									Type:   core.LiteralTypeString,
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
			name: "execute variable declaration #8",
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
			name: "execute variable declaration #9",
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
								Type:   core.LiteralTypeNumber,
								Value:  "1",
								Line:   3,
								CharAt: 8,
							},
							"c": &core.LiteralExpression{
								Type:   core.LiteralTypeNumber,
								Value:  "2",
								Line:   3,
								CharAt: 10,
							},
							"_args0_": &core.LiteralExpression{
								Type:   core.LiteralTypeNumber,
								Value:  "1",
								Line:   3,
								CharAt: 8,
							},
							"_args1_": &core.LiteralExpression{
								Type:   core.LiteralTypeNumber,
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
					Type:   core.LiteralTypeNumber,
					Value:  "3",
					Line:   3,
					CharAt: 8,
				}
				return gec
			},
			err: nil,
		},
		{
			name: "execute variable declaration #10",
			in:   `a:=1`,
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
							CharAt: 4,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute function declaration #1",
			in:   `func a(){}`,
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
					CharAt: 1,
				}
				return gec
			},
			err: nil,
		},
		{
			name: "execute function declaration #2",
			in: `
			func a(b,c){
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
						{Name: "b", Line: 2, CharAt: 8},
						{Name: "c", Line: 2, CharAt: 10},
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
								Type:   core.LiteralTypeNumber,
								Value:  "1",
								Line:   3,
								CharAt: 8,
							},
							"c": &core.LiteralExpression{
								Type:   core.LiteralTypeNumber,
								Value:  "2",
								Line:   3,
								CharAt: 10,
							},
							"_args0_": &core.LiteralExpression{
								Type:   core.LiteralTypeNumber,
								Value:  "1",
								Line:   3,
								CharAt: 8,
							},
							"_args1_": &core.LiteralExpression{
								Type:   core.LiteralTypeNumber,
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
					Type:   core.LiteralTypeNumber,
					Value:  "3",
					Line:   3,
					CharAt: 8,
				}
				return gec
			},
			err: nil,
		},
		{
			name: "execute assignment statement #1",
			in: `
				var a
				a=1
			`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.LiteralExpression{
							Type:   core.LiteralTypeNumber,
							Value:  "1",
							Line:   3,
							CharAt: 3,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute assignment statement #2",
			in: `
				var a,b
				a=1
				b=a+1
			`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.LiteralExpression{
							Type:   core.LiteralTypeNumber,
							Value:  "1",
							Line:   4,
							CharAt: 3,
						},
						"b": &core.LiteralExpression{
							Type:   core.LiteralTypeNumber,
							Value:  "2",
							Line:   4,
							CharAt: 3,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute assignment statement #3",
			in: `
				var a=[1,2]
				a[0]="xxx"
			`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.ArrayExpression{
							Elements: []core.Expression{
								&core.LiteralExpression{
									Type:   core.LiteralTypeString,
									Value:  "xxx",
									Line:   3,
									CharAt: 6,
								},
								&core.LiteralExpression{
									Type:   core.LiteralTypeNumber,
									Value:  "2",
									Line:   2,
									CharAt: 10,
								},
							},
							Line:   3,
							CharAt: 1,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute assignment statement #4",
			in: `
				var a=[1,[2,3]]
				a[1][1]="xxx"
			`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.ArrayExpression{
							Elements: []core.Expression{
								&core.LiteralExpression{
									Type:   core.LiteralTypeNumber,
									Value:  "1",
									Line:   2,
									CharAt: 8,
								},
								&core.ArrayExpression{
									Elements: []core.Expression{
										&core.LiteralExpression{
											Type:   core.LiteralTypeNumber,
											Value:  "2",
											Line:   2,
											CharAt: 11,
										},
										&core.LiteralExpression{
											Type:   core.LiteralTypeString,
											Value:  "xxx",
											Line:   3,
											CharAt: 9,
										},
									},
									Line:   3,
									CharAt: 1,
								},
							},
							Line:   3,
							CharAt: 1,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute assignment statement #5",
			in: `
				var a=[1,[2,3]]
				var b=[0,1]
				a[b[1]][1]="xxx"
			`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.ArrayExpression{
							Elements: []core.Expression{
								&core.LiteralExpression{
									Type:   core.LiteralTypeNumber,
									Value:  "1",
									Line:   2,
									CharAt: 8,
								},
								&core.ArrayExpression{
									Elements: []core.Expression{
										&core.LiteralExpression{
											Type:   core.LiteralTypeNumber,
											Value:  "2",
											Line:   2,
											CharAt: 11,
										},
										&core.LiteralExpression{
											Type:   core.LiteralTypeString,
											Value:  "xxx",
											Line:   4,
											CharAt: 12,
										},
									},
									Line:   4,
									CharAt: 1,
								},
							},
							Line:   4,
							CharAt: 1,
						},
						"b": &core.ArrayExpression{
							Elements: []core.Expression{
								&core.LiteralExpression{
									Type:   core.LiteralTypeNumber,
									Value:  "0",
									Line:   3,
									CharAt: 8,
								},
								&core.LiteralExpression{
									Type:   core.LiteralTypeNumber,
									Value:  "1",
									Line:   4,
									CharAt: 3,
								},
							},
							Line:   4,
							CharAt: 3,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute assignment statement #7",
			in: `
				var a={b:1}
				a.b=2
			`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.ObjectExpression{
							Properties: []*core.ObjectProperty{
								{
									KeyIdentifier: core.Identifier{
										Name:   "b",
										Line:   2,
										CharAt: 8,
									},
									Value: &core.LiteralExpression{
										Type:   core.LiteralTypeNumber,
										Value:  "2",
										Line:   3,
										CharAt: 5,
									},
									Line:   2,
									CharAt: 8,
								},
							},
							Line:   3,
							CharAt: 1,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute assignment statement #8",
			in: `
				var a={b:1}
				a.c=2
			`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.ObjectExpression{
							Properties: []*core.ObjectProperty{
								{
									KeyIdentifier: core.Identifier{
										Name:   "b",
										Line:   2,
										CharAt: 8,
									},
									Value: &core.LiteralExpression{
										Type:   core.LiteralTypeNumber,
										Value:  "1",
										Line:   2,
										CharAt: 10,
									},
									Line:   2,
									CharAt: 8,
								},
								{
									KeyIdentifier: core.Identifier{
										Name:   "c",
										Line:   3,
										CharAt: 3,
									},
									Value: &core.LiteralExpression{
										Type:   core.LiteralTypeNumber,
										Value:  "2",
										Line:   3,
										CharAt: 5,
									},
									Line:   3,
									CharAt: 5,
								},
							},
							Line:   3,
							CharAt: 1,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute assignment statement #8",
			in: `
				var a={b:1}
				a["c"]=2
			`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.ObjectExpression{
							Properties: []*core.ObjectProperty{
								{
									KeyIdentifier: core.Identifier{
										Name:   "b",
										Line:   2,
										CharAt: 8,
									},
									Value: &core.LiteralExpression{
										Type:   core.LiteralTypeNumber,
										Value:  "1",
										Line:   2,
										CharAt: 10,
									},
									Line:   2,
									CharAt: 8,
								},
								{
									KeyExpression: &core.LiteralExpression{
										Type:   core.LiteralTypeString,
										Value:  "c",
										Line:   3,
										CharAt: 3,
									},
									Value: &core.LiteralExpression{
										Type:   core.LiteralTypeNumber,
										Value:  "2",
										Line:   3,
										CharAt: 8,
									},
									Computed: true,
									Line:     3,
									CharAt:   8,
								},
							},
							Line:   3,
							CharAt: 1,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute assignment statement #9",
			in: `
				var a={b:1,c:[2,3]}
				a.c[1]="xxx"
			`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.ObjectExpression{
							Properties: []*core.ObjectProperty{
								{
									KeyIdentifier: core.Identifier{
										Name:   "b",
										Line:   2,
										CharAt: 8,
									},
									Value: &core.LiteralExpression{
										Type:   core.LiteralTypeNumber,
										Value:  "1",
										Line:   2,
										CharAt: 10,
									},
									Line:   2,
									CharAt: 8,
								},
								{
									KeyIdentifier: core.Identifier{
										Name:   "c",
										Line:   2,
										CharAt: 12,
									},
									Value: &core.ArrayExpression{
										Elements: []core.Expression{
											&core.LiteralExpression{
												Type:   core.LiteralTypeNumber,
												Value:  "2",
												Line:   2,
												CharAt: 15,
											},
											&core.LiteralExpression{
												Type:   core.LiteralTypeString,
												Value:  "xxx",
												Line:   3,
												CharAt: 8,
											},
										},
										Line:   3,
										CharAt: 1,
									},
									Line:   2,
									CharAt: 12,
								},
							},
							Line:   3,
							CharAt: 1,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute currying function #1",
			in: `
				func a(b){
					return func(c){
						return b+c
					}
				}
				var d=a(1)(2)
			`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				gec := &core.ExecutionContext{
					Variables: map[string]core.Expression{},
				}
				aec := &core.ExecutionContext{
					Outer: gec,
					Variables: map[string]core.Expression{
						"b": &core.LiteralExpression{
							Type:   core.LiteralTypeNumber,
							Value:  "1",
							Line:   4,
							CharAt: 8,
						},
						"_args0_": &core.LiteralExpression{
							Type:   core.LiteralTypeNumber,
							Value:  "1",
							Line:   4,
							CharAt: 8,
						},
					},
				}
				gec.Variables["a"] = &core.FunctionExpression{
					Params: []core.Identifier{
						{
							Name:   "b",
							Line:   2,
							CharAt: 8,
						},
					},
					Body: []core.Statement{
						core.ReturnStatement{
							Argument: &core.FunctionExpression{
								Params: []core.Identifier{
									{
										Name:   "c",
										Line:   3,
										CharAt: 13,
									},
								},
								Body: []core.Statement{
									core.ReturnStatement{
										Argument: &core.BinaryExpression{
											Left: &core.VariableExpression{
												Name:   "b",
												Line:   4,
												CharAt: 8,
											},
											Right: &core.VariableExpression{
												Name:   "c",
												Line:   4,
												CharAt: 10,
											},
											Operator: core.Operator{
												Symbol: "+",
												Line:   4,
												CharAt: 9,
											},
											Line:   4,
											CharAt: 8,
										},
										Line:   4,
										CharAt: 1,
									},
								},
								EC: &core.ExecutionContext{
									Outer: aec,
									Variables: map[string]core.Expression{
										"c": &core.LiteralExpression{
											Type:   core.LiteralTypeNumber,
											Value:  "2",
											Line:   4,
											CharAt: 10,
										},
										"_args0_": &core.LiteralExpression{
											Type:   core.LiteralTypeNumber,
											Value:  "2",
											Line:   4,
											CharAt: 10,
										},
									},
								},
								Line:   3,
								CharAt: 8,
							},
							Line:   3,
							CharAt: 1,
						},
					},
					EC:     aec,
					Line:   7,
					CharAt: 7,
				}
				gec.Variables["d"] = &core.LiteralExpression{
					Type:   core.LiteralTypeNumber,
					Value:  "3",
					Line:   4,
					CharAt: 8,
				}
				return gec
			},
			err: nil,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := lexer.Lex(tt.in)
			require.Equal(t, err, nil)
			stmts, _ := parser.ToAST(tokens)
			require.Equal(t, tt.err, Execute(tt.inEC, stmts))
			if tt.err == nil {
				require.Equal(t, tt.wantEC(), tt.inEC)
			}
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
			name: "execute function declaration #2",
			in: `
			func a(b,c){
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
						{Name: "b", Line: 2, CharAt: 8},
						{Name: "c", Line: 2, CharAt: 10},
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
								Type:   core.LiteralTypeNumber,
								Value:  "1",
								Line:   3,
								CharAt: 8,
							},
							"c": &core.LiteralExpression{
								Type:   core.LiteralTypeNumber,
								Value:  "2",
								Line:   3,
								CharAt: 10,
							},
							"_args0_": &core.LiteralExpression{
								Type:   core.LiteralTypeNumber,
								Value:  "1",
								Line:   3,
								CharAt: 8,
							},
							"_args1_": &core.LiteralExpression{
								Type:   core.LiteralTypeNumber,
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
					Type:   core.LiteralTypeNumber,
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
			tokens, err := lexer.Lex(tt.in)
			require.Equal(t, err, nil)
			stmts, _ := parser.ToAST(tokens)
			require.Equal(t, tt.err, Execute(tt.inEC, stmts))
			require.Equal(t, tt.wantEC(), tt.inEC)
		})
	}
}
