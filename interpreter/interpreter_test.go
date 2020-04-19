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
				 var c = #f`,
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
							Value:  "#f",
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
				   1+#f,
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
									Value:  "1#f",
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
			name: "execute function declaration #3",
			in: `var a=1
				 func b() {
					a=2
				 }
				 b()`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				gec := &core.ExecutionContext{
					Variables: map[string]core.Expression{},
				}
				gec.Variables["a"] = &core.LiteralExpression{
					Type:   "number",
					Value:  "2",
					Line:   3,
					CharAt: 3,
				}
				gec.Variables["b"] = &core.FunctionExpression{
					Params: []core.Identifier{},
					Body: []core.Statement{
						core.AssignmentStatement{
							Left: &core.VariableExpression{
								Name:   "a",
								Line:   3,
								CharAt: 1,
							},
							Right: &core.LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   3,
								CharAt: 3,
							},
							Line:   3,
							CharAt: 1,
						},
					},
					EC: &core.ExecutionContext{
						Outer:     gec,
						Variables: map[string]core.Expression{},
					},
					Line:   5,
					CharAt: 1,
				}
				return gec
			},
			err: nil,
		},
		{
			name: "execute function declaration #4",
			in: `func a(){
				break
			}
			a()
			`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: nil,
			err:    core.BreakError{Message: "break is not in a loop. [2,1]"},
		},
		{
			name: "execute function declaration #5",
			in: `func a(){
				continue
			}
			a()
			`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: nil,
			err:    core.ContinueError{Message: "continue is not in a loop. [2,1]"},
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
		{
			name: "execute if statement #1",
			in: `
				var a
				if #t {
					a=1
				}`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				gec := &core.ExecutionContext{
					Variables: map[string]core.Expression{},
				}
				gec.Variables["a"] = &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   4,
					CharAt: 3,
				}
				return gec
			},
			err: nil,
		},
		{
			name: "execute if statement #2",
			in: `
				var a=1
				if a==1 {
					a=2
				}`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				gec := &core.ExecutionContext{
					Variables: map[string]core.Expression{},
				}
				gec.Variables["a"] = &core.LiteralExpression{
					Type:   "number",
					Value:  "2",
					Line:   4,
					CharAt: 3,
				}
				return gec
			},
			err: nil,
		},
		{
			name: "execute if statement #4",
			in: `
				var a=1
				if !a {
					a=2
				}`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				gec := &core.ExecutionContext{
					Variables: map[string]core.Expression{},
				}
				gec.Variables["a"] = &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   3,
					CharAt: 5,
				}
				return gec
			},
			err: nil,
		},
		{
			name: "execute if statement #5",
			in: `
				var a=1
				if !a {
					a=2
				} elif a==1 {
					a=3	
				}`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				gec := &core.ExecutionContext{
					Variables: map[string]core.Expression{},
				}
				gec.Variables["a"] = &core.LiteralExpression{
					Type:   "number",
					Value:  "3",
					Line:   6,
					CharAt: 3,
				}
				return gec
			},
			err: nil,
		},
		{
			name: "execute if statement #6",
			in: `
				var a=1
				if !a {
					a=2
				} elif a==1 && #f {
					a=3	
				} elif a {
					a=4	
				}`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				gec := &core.ExecutionContext{
					Variables: map[string]core.Expression{},
				}
				gec.Variables["a"] = &core.LiteralExpression{
					Type:   "number",
					Value:  "4",
					Line:   8,
					CharAt: 3,
				}
				return gec
			},
			err: nil,
		},
		{
			name: "execute if statement #8",
			in: `
				var a=1
				if !a {
					a=2
				} elif a==1 && #f {
					a=3	
				} elif #f {
					a=4	
				} else {
					a=5	
				}`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				gec := &core.ExecutionContext{
					Variables: map[string]core.Expression{},
				}
				gec.Variables["a"] = &core.LiteralExpression{
					Type:   "number",
					Value:  "5",
					Line:   10,
					CharAt: 3,
				}
				return gec
			},
			err: nil,
		},
		{
			name: "execute if statement #9",
			in: `
				var a=1
				if b:=2; a {
					a=b
				}`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				gec := &core.ExecutionContext{
					Variables: map[string]core.Expression{},
				}
				gec.Variables["a"] = &core.LiteralExpression{
					Type:   "number",
					Value:  "2",
					Line:   4,
					CharAt: 3,
				}
				return gec
			},
			err: nil,
		},
		{
			name: "execute if statement #10",
			in: `
				var a=1
				var b
				if a:=2; a!=1 {
					b=3
				}`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   2,
							CharAt: 7,
						},
						"b": &core.LiteralExpression{
							Type:   "number",
							Value:  "3",
							Line:   5,
							CharAt: 3,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute if statement #11",
			in: `
				var b
				if a:=2; a!=2 {
					b=3
				}elif a==2{
					b=4
				}`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"b": &core.LiteralExpression{
							Type:   "number",
							Value:  "4",
							Line:   6,
							CharAt: 3,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute if statement #12",
			in: `
				var c
				if a:=2; a!=2 {
					c=3
				}elif b:=1; (a+(b))>=0{
					c=4
				}`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"c": &core.LiteralExpression{
							Type:   "number",
							Value:  "4",
							Line:   6,
							CharAt: 3,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute if statement #13",
			in: `
				if #t {
					break
				}`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: nil,
			err:    core.BreakError{Message: "break is not in a loop. [3,1]"},
		},
		{
			name: "execute if statement #14",
			in: `
				if #t {
					continue
				}`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: nil,
			err:    core.ContinueError{Message: "continue is not in a loop. [3,1]"},
		},
		{
			name: "execute for statement #1",
			in: `
				a:=0
				for i:=0;i<3;i=i+1 {
					a=a+1
				}
				`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.LiteralExpression{
							Type:   "number",
							Value:  "3",
							Line:   4,
							CharAt: 3,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute for statement #2",
			in: `
				a:=0
				i:=0
				for i<3;i=i+1 {
					a=a+1
				}
				`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.LiteralExpression{
							Type:   "number",
							Value:  "3",
							Line:   5,
							CharAt: 3,
						},
						"i": &core.LiteralExpression{
							Type:   "number",
							Value:  "3",
							Line:   4,
							CharAt: 5,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute for statement #3",
			in: `
				a:=0
				for i:=0;i<3 {
					a=a+1
					i=i+1
				}
				`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.LiteralExpression{
							Type:   "number",
							Value:  "3",
							Line:   4,
							CharAt: 3,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute for statement #4",
			in: `
				a:=0
				for a<3 {
					a=a+1
				}
				`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.LiteralExpression{
							Type:   "number",
							Value:  "3",
							Line:   3,
							CharAt: 5,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute for statement #5",
			in: `
				a:=0
				for a<3 {
					a=a+1
					if a==2 {
						return
					}
				}
				`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   5,
							CharAt: 4,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute for statement #6",
			in: `
				a:=0
				for {
					a=a+1
					if a==10 {
						return
					}
				}
				`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.LiteralExpression{
							Type:   "number",
							Value:  "10",
							Line:   5,
							CharAt: 4,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute for statement #7",
			in: `
				a:=0
				for {
					a=a+1
					break
				}
				`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   4,
							CharAt: 3,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute for statement #8",
			in: `
				a:=0
				for {
					a=a+1
					if a == 2 {
						break
					}
				}
				`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   5,
							CharAt: 4,
						},
					},
				}
			},
			err: nil,
		},
		{
			name: "execute for statement #9",
			in: `
				a:=0
				for i:=0;i<3;i=i+1 {
					if i<2 {
						continue
					}
					a=a+1
				}
				`,
			inEC: &core.ExecutionContext{
				Variables: map[string]core.Expression{},
			},
			wantEC: func() *core.ExecutionContext {
				return &core.ExecutionContext{
					Variables: map[string]core.Expression{
						"a": &core.LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   7,
							CharAt: 3,
						},
					},
				}
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
	}{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := lexer.Lex(tt.in)
			require.Equal(t, err, nil)
			stmts, _ := parser.ToAST(tokens)
			err = Execute(tt.inEC, stmts)
			require.Equal(t, tt.err, err)
			if err == nil {
				require.Equal(t, tt.wantEC(), tt.inEC)
			}
		})
	}
}
