package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gs/lexer"
	"gs/operator"
)

func TestToAST_VariableDeclaration(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []Statement
	}{
		// {
		// 	name: "parse variable declaration statement (without initialization)",
		// 	in:   `var a`,
		// 	want: []Statement{
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "a",
		// 						Line:   1,
		// 						CharAt: 5,
		// 					},
		// 					Init:   nil,
		// 					Line:   1,
		// 					CharAt: 5,
		// 				},
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 	},
		// },
		// {
		// 	name: "parse variable declaration statement (without initialization)",
		// 	in: `var a
		// 		 var b`,
		// 	want: []Statement{
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "a",
		// 						Line:   1,
		// 						CharAt: 5,
		// 					},
		// 					Init:   nil,
		// 					Line:   1,
		// 					CharAt: 5,
		// 				},
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "b",
		// 						Line:   2,
		// 						CharAt: 5,
		// 					},
		// 					Init:   nil,
		// 					Line:   2,
		// 					CharAt: 5,
		// 				},
		// 			},
		// 			Line:   2,
		// 			CharAt: 1,
		// 		},
		// 	},
		// },
		// {
		// 	name: "parse number declaration statement",
		// 	in:   "var a=1",
		// 	want: []Statement{
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "a",
		// 						Line:   1,
		// 						CharAt: 5,
		// 					},
		// 					Init: LiteralExpression{
		// 						Type:   "number",
		// 						Value:  "1",
		// 						Line:   1,
		// 						CharAt: 7,
		// 					},
		// 					Line:   1,
		// 					CharAt: 5,
		// 				},
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 	},
		// },
		// {
		// 	name: "parse string declaration statement",
		// 	in:   `var a="xxx"`,
		// 	want: []Statement{
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "a",
		// 						Line:   1,
		// 						CharAt: 5,
		// 					},
		// 					Init: LiteralExpression{
		// 						Type:   "string",
		// 						Value:  `"xxx"`,
		// 						Line:   1,
		// 						CharAt: 7,
		// 					},
		// 					Line:   1,
		// 					CharAt: 5,
		// 				},
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 	},
		// },
		// {
		// 	name: "parse boolean declaration statement",
		// 	in:   "var a=false",
		// 	want: []Statement{
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "a",
		// 						Line:   1,
		// 						CharAt: 5,
		// 					},
		// 					Init: LiteralExpression{
		// 						Type:   "boolean",
		// 						Value:  "false",
		// 						Line:   1,
		// 						CharAt: 7,
		// 					},
		// 					Line:   1,
		// 					CharAt: 5,
		// 				},
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 	},
		// },
		// {
		// 	name: "parse variable expression",
		// 	in: `var a=false
		// 		 var b=a`,
		// 	want: []Statement{
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "a",
		// 						Line:   1,
		// 						CharAt: 5,
		// 					},
		// 					Init: LiteralExpression{
		// 						Type:   "boolean",
		// 						Value:  "false",
		// 						Line:   1,
		// 						CharAt: 7,
		// 					},
		// 					Line:   1,
		// 					CharAt: 5,
		// 				},
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "b",
		// 						Line:   2,
		// 						CharAt: 5,
		// 					},
		// 					Init: VariableExpression{
		// 						Name:   "a",
		// 						Line:   2,
		// 						CharAt: 7,
		// 					},
		// 					Line:   2,
		// 					CharAt: 5,
		// 				},
		// 			},
		// 			Line:   2,
		// 			CharAt: 1,
		// 		},
		// 	},
		// },
		// {
		// 	name: "parse declaration statement (same statement, multi variable)",
		// 	in:   "var a,b=1,'2'",
		// 	want: []Statement{
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "a",
		// 						Line:   1,
		// 						CharAt: 5,
		// 					},
		// 					Init: LiteralExpression{
		// 						Type:   "number",
		// 						Value:  "1",
		// 						Line:   1,
		// 						CharAt: 9,
		// 					},
		// 					Line:   1,
		// 					CharAt: 5,
		// 				},
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "b",
		// 						Line:   1,
		// 						CharAt: 7,
		// 					},
		// 					Init: LiteralExpression{
		// 						Type:   "string",
		// 						Value:  "'2'",
		// 						Line:   1,
		// 						CharAt: 11,
		// 					},
		// 					Line:   1,
		// 					CharAt: 7,
		// 				},
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 	},
		// },
		// {
		// 	name: "parse declaration statement (same statement, multi variable, b is not initialized)",
		// 	in: `var a,b=1
		// 	     var c=2`,
		// 	want: []Statement{
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "a",
		// 						Line:   1,
		// 						CharAt: 5,
		// 					},
		// 					Init: LiteralExpression{
		// 						Type:   "number",
		// 						Value:  "1",
		// 						Line:   1,
		// 						CharAt: 9,
		// 					},
		// 					Line:   1,
		// 					CharAt: 5,
		// 				},
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "b",
		// 						Line:   1,
		// 						CharAt: 7,
		// 					},
		// 					Init:   nil,
		// 					Line:   1,
		// 					CharAt: 7,
		// 				},
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "c",
		// 						Line:   2,
		// 						CharAt: 5,
		// 					},
		// 					Init: LiteralExpression{
		// 						Type:   "number",
		// 						Value:  "2",
		// 						Line:   2,
		// 						CharAt: 7,
		// 					},
		// 					Line:   2,
		// 					CharAt: 5,
		// 				},
		// 			},
		// 			Line:   2,
		// 			CharAt: 1,
		// 		},
		// 	},
		// },
		// {
		// 	name: "parse array declaration statement",
		// 	in: `var a=[
		// 			123,
		// 			"456",
		// 			1+1
		// 		]`,
		// 	want: []Statement{
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "a",
		// 						Line:   1,
		// 						CharAt: 5,
		// 					},
		// 					Init: ArrayExpression{
		// 						Elements: []Expression{
		// 							LiteralExpression{
		// 								Type:   "number",
		// 								Value:  "123",
		// 								Line:   2,
		// 								CharAt: 1,
		// 							},
		// 							LiteralExpression{
		// 								Type:   "string",
		// 								Value:  `"456`,
		// 								Line:   3,
		// 								CharAt: 1,
		// 							},
		// 							BinaryExpression{
		// 								Left: LiteralExpression{
		// 									Type:   "number",
		// 									Value:  "1",
		// 									Line:   4,
		// 									CharAt: 1,
		// 								},
		// 								Right: LiteralExpression{
		// 									Type:   "number",
		// 									Value:  "1",
		// 									Line:   4,
		// 									CharAt: 3,
		// 								},
		// 								Operator: "+",
		// 							},
		// 						},
		// 					},
		// 					Line:   1,
		// 					CharAt: 1,
		// 				},
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 	},
		// },
		// {
		// 	name: "parse function declaration statement",
		// 	in:   `var a=func(b){}`,
		// 	want: []Statement{
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "a",
		// 						Line:   1,
		// 						CharAt: 5,
		// 					},
		// 					Init: FunctionExpression{
		// 						Params: []Identifier{
		// 							Identifier{
		// 								Name:   "b",
		// 								Line:   1,
		// 								CharAt: 12,
		// 							},
		// 						},
		// 						Body:   []Statement{},
		// 						Line:   1,
		// 						CharAt: 7,
		// 					},
		// 				},
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 	},
		// },
		// {
		// 	name: "parse declaration statement and variable expression",
		// 	in: `var a
		// 		 a`,
		// 	want: []Statement{
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "a",
		// 						Line:   1,
		// 						CharAt: 5,
		// 					},
		// 					Init:   nil,
		// 					Line:   1,
		// 					CharAt: 5,
		// 				},
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 		ExpressionStatement{
		// 			Expression: VariableExpression{
		// 				Name:   "a",
		// 				Line:   2,
		// 				CharAt: 1,
		// 			},
		// 			Line:   2,
		// 			CharAt: 1,
		// 		},
		// 	},
		// },
		// {
		// 	name: "parse declaration statement",
		// 	in:   "a,b:=1,'2'",
		// 	want: []Statement{
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "a",
		// 						Line:   1,
		// 						CharAt: 1,
		// 					},
		// 					Init: LiteralExpression{
		// 						Type:   "number",
		// 						Value:  "1",
		// 						Line:   1,
		// 						CharAt: 6,
		// 					},
		// 					Line:   1,
		// 					CharAt: 1,
		// 				},
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "b",
		// 						Line:   1,
		// 						CharAt: 3,
		// 					},
		// 					Init: LiteralExpression{
		// 						Type:   "string",
		// 						Value:  "'2'",
		// 						Line:   1,
		// 						CharAt: 8,
		// 					},
		// 					Line:   1,
		// 					CharAt: 1,
		// 				},
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 	},
		// },
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ast, _ := ToAST(lexer.Lex(tt.in))
			require.Equal(t, tt.want, ast)
		})
	}
}

func TestToAST_ObjectDeclaration(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []Statement
	}{
		// {
		// 	name: "parse object declaration statement",
		// 	in: `var a={
		// 			b: 1,
		// 			c: "abc",
		// 			d: false
		// 		}`,
		// 	want: []Statement{
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "a",
		// 						Line:   1,
		// 						CharAt: 5,
		// 					},
		// 					Init: ObjectExpression{
		// 						Properties: []ObjectProperty{
		// 							ObjectProperty{
		// 								KeyIdentifier: Identifier{
		// 									Name:   "b",
		// 									Line:   2,
		// 									CharAt: 1,
		// 								},
		// 								Value: LiteralExpression{
		// 									Type:   "number",
		// 									Value:  "1",
		// 									Line:   2,
		// 									CharAt: 4,
		// 								},
		// 								Computed: false,
		// 							},
		// 							ObjectProperty{
		// 								KeyIdentifier: Identifier{
		// 									Name:   "c",
		// 									Line:   3,
		// 									CharAt: 1,
		// 								},
		// 								Value: LiteralExpression{
		// 									Type:   "string",
		// 									Value:  `"abc"`,
		// 									Line:   3,
		// 									CharAt: 4,
		// 								},
		// 								Computed: false,
		// 							},
		// 							ObjectProperty{
		// 								KeyIdentifier: Identifier{
		// 									Name:   "d",
		// 									Line:   4,
		// 									CharAt: 1,
		// 								},
		// 								Value: LiteralExpression{
		// 									Type:   "boolean",
		// 									Value:  "false",
		// 									Line:   4,
		// 									CharAt: 4,
		// 								},
		// 								Computed: false,
		// 							},
		// 						},
		// 						Line:   1,
		// 						CharAt: 7,
		// 					},
		// 					Line:   1,
		// 					CharAt: 5,
		// 				},
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 	},
		// },
		// {
		// 	name: "parse nesting object declaration statement",
		// 	in: `var a={
		// 			b: {
		// 				c: 1
		// 			}
		// 		}`,
		// 	want: []Statement{
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "a",
		// 						Line:   1,
		// 						CharAt: 5,
		// 					},
		// 					Init: ObjectExpression{
		// 						Properties: []ObjectProperty{
		// 							ObjectProperty{
		// 								KeyIdentifier: Identifier{
		// 									Name:   "b",
		// 									Line:   2,
		// 									CharAt: 1,
		// 								},
		// 								Value: ObjectExpression{
		// 									Properties: []ObjectProperty{
		// 										ObjectProperty{
		// 											KeyIdentifier: Identifier{
		// 												Name:   "c",
		// 												Line:   3,
		// 												CharAt: 1,
		// 											},
		// 											Value: LiteralExpression{
		// 												Type:   "number",
		// 												Value:  "1",
		// 												Line:   3,
		// 												CharAt: 4,
		// 											},
		// 											Computed: false,
		// 										},
		// 									},
		// 									Line:   2,
		// 									CharAt: 4,
		// 								},
		// 								Computed: false,
		// 							},
		// 						},
		// 						Line:   1,
		// 						CharAt: 7,
		// 					},
		// 					Line:   1,
		// 					CharAt: 5,
		// 				},
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 	},
		// },
		// {
		// 	name: "parse object declaration statement with array as property value",
		// 	in: `var a={
		// 			b: [1,2]
		// 		}`,
		// 	want: []Statement{
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "a",
		// 						Line:   1,
		// 						CharAt: 5,
		// 					},
		// 					Init: ObjectExpression{
		// 						Properties: []ObjectProperty{
		// 							ObjectProperty{
		// 								KeyIdentifier: Identifier{
		// 									Name:   "b",
		// 									Line:   2,
		// 									CharAt: 1,
		// 								},
		// 								Value: ArrayExpression{
		// 									Elements: []Expression{
		// 										LiteralExpression{
		// 											Type:   "number",
		// 											Value:  "1",
		// 											Line:   2,
		// 											CharAt: 5,
		// 										},
		// 										LiteralExpression{
		// 											Type:   "number",
		// 											Value:  "2",
		// 											Line:   2,
		// 											CharAt: 7,
		// 										},
		// 									},
		// 								},
		// 								Computed: false,
		// 							},
		// 						},
		// 						Line:   1,
		// 						CharAt: 9,
		// 					},
		// 					Line:   1,
		// 					CharAt: 1,
		// 				},
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 	},
		// },
		// {
		// 	name: "parse object declaration statement with function as property value",
		// 	in: `var a={
		// 			b: func(c){}
		// 		}`,
		// 	want: []Statement{
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "a",
		// 						Line:   1,
		// 						CharAt: 5,
		// 					},
		// 					Init: ObjectExpression{
		// 						Properties: []ObjectProperty{
		// 							ObjectProperty{
		// 								KeyIdentifier: Identifier{
		// 									Name:   "b",
		// 									Line:   2,
		// 									CharAt: 1,
		// 								},
		// 								Value: FunctionExpression{
		// 									Params: []Identifier{
		// 										Identifier{
		// 											Name:   "c",
		// 											Line:   2,
		// 											CharAt: 9,
		// 										},
		// 									},
		// 									Body:   []Statement{},
		// 									Line:   2,
		// 									CharAt: 4,
		// 								},
		// 								Computed: false,
		// 							},
		// 						},
		// 						Line:   1,
		// 						CharAt: 9,
		// 					},
		// 					Line:   1,
		// 					CharAt: 1,
		// 				},
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 	},
		// },
		// {
		// 	name: "parse object declaration statement with compute key",
		// 	in: `var a={
		// 			[1+1]: 1
		// 		}`,
		// 	want: []Statement{
		// 		VariableDeclaration{
		// 			Declarations: []VariableDeclarator{
		// 				VariableDeclarator{
		// 					ID: Identifier{
		// 						Name:   "a",
		// 						Line:   1,
		// 						CharAt: 5,
		// 					},
		// 					Init: ObjectExpression{
		// 						Properties: []ObjectProperty{
		// 							ObjectProperty{
		// 								KeyExpression: BinaryExpression{
		// 									Left: LiteralExpression{
		// 										Type:   "number",
		// 										Value:  "1",
		// 										Line:   2,
		// 										CharAt: 2,
		// 									},
		// 									Right: LiteralExpression{
		// 										Type:   "number",
		// 										Value:  "1",
		// 										Line:   2,
		// 										CharAt: 4,
		// 									},
		// 									Operator: "+",
		// 								},
		// 								Value: LiteralExpression{
		// 									Type:   "number",
		// 									Value:  "1",
		// 									Line:   2,
		// 									CharAt: 4,
		// 								},
		// 								Computed: true,
		// 							},
		// 						},
		// 						Line:   1,
		// 						CharAt: 9,
		// 					},
		// 					Line:   1,
		// 					CharAt: 1,
		// 				},
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 	},
		// },
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ast, _ := ToAST(lexer.Lex(tt.in))
			require.Equal(t, tt.want, ast)
		})
	}
}

func TestParseExpression(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want Expression
	}{
		// {
		// 	name: "parse number expression",
		// 	in:   "1",
		// 	want: LiteralExpression{
		// 		Type:   "number",
		// 		Value:  "1",
		// 		Line:   1,
		// 		CharAt: 1,
		// 	},
		// },
		// {
		// 	name: "parse string expression",
		// 	in:   `"1"`,
		// 	want: LiteralExpression{
		// 		Type:   "string",
		// 		Value:  `"1"`,
		// 		Line:   1,
		// 		CharAt: 1,
		// 	},
		// },
		// {
		// 	name: "parse boolean expression",
		// 	in:   "false",
		// 	want: LiteralExpression{
		// 		Type:   "boolean",
		// 		Value:  `false`,
		// 		Line:   1,
		// 		CharAt: 1,
		// 	},
		// },
		// {
		// 	name: "parse binary expression",
		// 	in:   "1+1",
		// 	want: BinaryExpression{
		// 		Left: LiteralExpression{
		// 			Type:   "number",
		// 			Value:  "1",
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 		Right: LiteralExpression{
		// 			Type:   "number",
		// 			Value:  "1",
		// 			Line:   1,
		// 			CharAt: 3,
		// 		},
		// 		Operator: operator.Operator{
		// 			Symbol: "+",
		// 			Line:   1,
		// 			CharAt: 2,
		// 		},
		// 		Line:   1,
		// 		CharAt: 1,
		// 	},
		// },
		// {
		// 	name: "parse object expression",
		// 	in: `{
		// 			a: 1,
		// 			b: 2,
		// 		 }`,
		// 	want: ObjectExpression{
		// 		Properties: []ObjectProperty{
		// 			ObjectProperty{
		// 				KeyIdentifier: Identifier{
		// 					Name:   "a",
		// 					Line:   2,
		// 					CharAt: 1,
		// 				},
		// 				Value: LiteralExpression{
		// 					Type:   "number",
		// 					Value:  "1",
		// 					Line:   2,
		// 					CharAt: 4,
		// 				},
		// 				Computed: false,
		// 			},
		// 			ObjectProperty{
		// 				KeyIdentifier: Identifier{
		// 					Name:   "b",
		// 					Line:   3,
		// 					CharAt: 1,
		// 				},
		// 				Value: LiteralExpression{
		// 					Type:   "number",
		// 					Value:  "2",
		// 					Line:   3,
		// 					CharAt: 4,
		// 				},
		// 				Computed: false,
		// 			},
		// 		},
		// 	},
		// },
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp, _, _ := parseExpression(lexer.Lex(tt.in))
			require.Equal(t, tt.want, exp)
		})
	}
}

func TestParseExpression_Precedence(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want Expression
	}{
		// {
		// 	name: "parse binary expression (1+1)-2",
		// 	in:   "1+1-2",
		// 	want: BinaryExpression{
		// 		Left: BinaryExpression{
		// 			Left: LiteralExpression{
		// 				Type:   "number",
		// 				Value:  "1",
		// 				Line:   1,
		// 				CharAt: 1,
		// 			},
		// 			Right: LiteralExpression{
		// 				Type:   "number",
		// 				Value:  "1",
		// 				Line:   1,
		// 				CharAt: 3,
		// 			},
		// 			Operator: operator.Operator{
		// 				Symbol: "+",
		// 				Line:   1,
		// 				CharAt: 2,
		// 			},
		// 			Nesting: 0,
		// 			Line:    1,
		// 			CharAt:  1,
		// 		},
		// 		Right: LiteralExpression{
		// 			Type:   "number",
		// 			Value:  "2",
		// 			Line:   1,
		// 			CharAt: 5,
		// 		},
		// 		Operator: operator.Operator{
		// 			Symbol: "-",
		// 			Line:   1,
		// 			CharAt: 4,
		// 		},
		// 		Nesting: 1,
		// 		Line:    1,
		// 		CharAt:  1,
		// 	},
		// },
		// {
		// 	name: "parse binary expression ((1*2)*3)-4",
		// 	in:   "1*2*3-4",
		// 	want: BinaryExpression{
		// 		Left: BinaryExpression{
		// 			Left: BinaryExpression{
		// 				Left: LiteralExpression{
		// 					Type:   "number",
		// 					Value:  "1",
		// 					Line:   1,
		// 					CharAt: 1,
		// 				},
		// 				Right: LiteralExpression{
		// 					Type:   "number",
		// 					Value:  "2",
		// 					Line:   1,
		// 					CharAt: 3,
		// 				},
		// 				Operator: operator.Operator{
		// 					Symbol: "*",
		// 					Line:   1,
		// 					CharAt: 2,
		// 				},
		// 				Nesting: 0,
		// 				Line:    1,
		// 				CharAt:  1,
		// 			},
		// 			Right: LiteralExpression{
		// 				Type:   "number",
		// 				Value:  "3",
		// 				Line:   1,
		// 				CharAt: 5,
		// 			},
		// 			Operator: operator.Operator{
		// 				Symbol: "*",
		// 				Line:   1,
		// 				CharAt: 4,
		// 			},
		// 			Nesting: 1,
		// 			Line:    1,
		// 			CharAt:  1,
		// 		},
		// 		Right: LiteralExpression{
		// 			Type:   "number",
		// 			Value:  "4",
		// 			Line:   1,
		// 			CharAt: 7,
		// 		},
		// 		Operator: operator.Operator{
		// 			Symbol: "-",
		// 			Line:   1,
		// 			CharAt: 6,
		// 		},
		// 		Nesting: 2,
		// 		Line:    1,
		// 		CharAt:  1,
		// 	},
		// },
		// {
		// 	name: "parse binary expression ((((1+2)+3)+4)+5",
		// 	in:   "1+2+3+4+5",
		// 	want: BinaryExpression{
		// 		Left: BinaryExpression{
		// 			Left: BinaryExpression{
		// 				Left: BinaryExpression{
		// 					Left: LiteralExpression{
		// 						Type:   "number",
		// 						Value:  "1",
		// 						Line:   1,
		// 						CharAt: 1,
		// 					},
		// 					Right: LiteralExpression{
		// 						Type:   "number",
		// 						Value:  "2",
		// 						Line:   1,
		// 						CharAt: 3,
		// 					},
		// 					Operator: operator.Operator{
		// 						Symbol: "+",
		// 						Line:   1,
		// 						CharAt: 2,
		// 					},
		// 					Nesting: 0,
		// 					Line:    1,
		// 					CharAt:  1,
		// 				},
		// 				Right: LiteralExpression{
		// 					Type:   "number",
		// 					Value:  "3",
		// 					Line:   1,
		// 					CharAt: 5,
		// 				},
		// 				Operator: operator.Operator{
		// 					Symbol: "+",
		// 					Line:   1,
		// 					CharAt: 4,
		// 				},
		// 				Nesting: 1,
		// 				Line:    1,
		// 				CharAt:  1,
		// 			},
		// 			Right: LiteralExpression{
		// 				Type:   "number",
		// 				Value:  "4",
		// 				Line:   1,
		// 				CharAt: 7,
		// 			},
		// 			Operator: operator.Operator{
		// 				Symbol: "+",
		// 				Line:   1,
		// 				CharAt: 6,
		// 			},
		// 			Nesting: 2,
		// 			Line:    1,
		// 			CharAt:  1,
		// 		},
		// 		Right: LiteralExpression{
		// 			Type:   "number",
		// 			Value:  "5",
		// 			Line:   1,
		// 			CharAt: 9,
		// 		},
		// 		Operator: operator.Operator{
		// 			Symbol: "+",
		// 			Line:   1,
		// 			CharAt: 8,
		// 		},
		// 		Nesting: 3,
		// 		Line:    1,
		// 		CharAt:  1,
		// 	},
		// },
		// {
		// 	name: "parse binary expression 1+(2*3)",
		// 	in:   "1+2*3",
		// 	want: BinaryExpression{
		// 		Left: LiteralExpression{
		// 			Type:   "number",
		// 			Value:  "1",
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 		Right: BinaryExpression{
		// 			Left: LiteralExpression{
		// 				Type:   "number",
		// 				Value:  "2",
		// 				Line:   1,
		// 				CharAt: 3,
		// 			},
		// 			Right: LiteralExpression{
		// 				Type:   "number",
		// 				Value:  "3",
		// 				Line:   1,
		// 				CharAt: 5,
		// 			},
		// 			Operator: operator.Operator{
		// 				Symbol: "*",
		// 				Line:   1,
		// 				CharAt: 4,
		// 			},
		// 			Nesting: 0,
		// 			Line:    1,
		// 			CharAt:  3,
		// 		},
		// 		Operator: operator.Operator{
		// 			Symbol: "+",
		// 			Line:   1,
		// 			CharAt: 2,
		// 		},
		// 		Nesting: 1,
		// 		Line:    1,
		// 		CharAt:  1,
		// 	},
		// },
		// {
		// 	name: "parse binary expression (2*3)/4",
		// 	in:   "2*3/4",
		// 	want: BinaryExpression{
		// 		Left: BinaryExpression{
		// 			Left: LiteralExpression{
		// 				Type:   "number",
		// 				Value:  "2",
		// 				Line:   1,
		// 				CharAt: 1,
		// 			},
		// 			Right: LiteralExpression{
		// 				Type:   "number",
		// 				Value:  "3",
		// 				Line:   1,
		// 				CharAt: 3,
		// 			},
		// 			Operator: operator.Operator{
		// 				Symbol: "*",
		// 				Line:   1,
		// 				CharAt: 2,
		// 			},
		// 			Nesting: 0,
		// 			Line:    1,
		// 			CharAt:  1,
		// 		},
		// 		Right: LiteralExpression{
		// 			Type:   "number",
		// 			Value:  "4",
		// 			Line:   1,
		// 			CharAt: 5,
		// 		},
		// 		Operator: operator.Operator{
		// 			Symbol: "/",
		// 			Line:   1,
		// 			CharAt: 4,
		// 		},
		// 		Nesting: 1,
		// 		Line:    1,
		// 		CharAt:  1,
		// 	},
		// },
		// {
		// 	name: "parse binary expression with parantheses (1+2)*3",
		// 	in:   "(1+2)*3",
		// 	want: BinaryExpression{
		// 		Left: BinaryExpression{
		// 			Left: LiteralExpression{
		// 				Type:   "number",
		// 				Value:  "1",
		// 				Line:   1,
		// 				CharAt: 2,
		// 			},
		// 			Right: LiteralExpression{
		// 				Type:   "number",
		// 				Value:  "2",
		// 				Line:   1,
		// 				CharAt: 4,
		// 			},
		// 			Operator: operator.Operator{
		// 				Symbol: "+",
		// 				Line:   1,
		// 				CharAt: 3,
		// 			},
		// 			Group:   true,
		// 			Nesting: 0,
		// 			Line:    1,
		// 			CharAt:  2,
		// 		},
		// 		Right: LiteralExpression{
		// 			Type:   "number",
		// 			Value:  "3",
		// 			Line:   1,
		// 			CharAt: 7,
		// 		},
		// 		Operator: operator.Operator{
		// 			Symbol: "*",
		// 			Line:   1,
		// 			CharAt: 6,
		// 		},
		// 		Nesting: 1,
		// 		Line:    1,
		// 		CharAt:  1,
		// 	},
		// },
		// {
		// 	name: "parse binary expression with parantheses 1+(2+3)",
		// 	in:   "1+(2+3)",
		// 	want: BinaryExpression{
		// 		Left: LiteralExpression{
		// 			Type:   "number",
		// 			Value:  "1",
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 		Right: BinaryExpression{
		// 			Left: LiteralExpression{
		// 				Type:   "number",
		// 				Value:  "2",
		// 				Line:   1,
		// 				CharAt: 4,
		// 			},
		// 			Right: LiteralExpression{
		// 				Type:   "number",
		// 				Value:  "3",
		// 				Line:   1,
		// 				CharAt: 6,
		// 			},
		// 			Operator: operator.Operator{
		// 				Symbol: "+",
		// 				Line:   1,
		// 				CharAt: 5,
		// 			},
		// 			Group:   true,
		// 			Nesting: 0,
		// 			Line:    1,
		// 			CharAt:  4,
		// 		},
		// 		Operator: operator.Operator{
		// 			Symbol: "+",
		// 			Line:   1,
		// 			CharAt: 2,
		// 		},
		// 		Nesting: 1,
		// 		Line:    1,
		// 		CharAt:  1,
		// 	},
		// },
		// {
		// 	name: "parse binary expression with parantheses (1+(2+3))+4",
		// 	in:   "1+(2+3)+4",
		// 	want: BinaryExpression{
		// 		Left: BinaryExpression{
		// 			Left: LiteralExpression{
		// 				Type:   "number",
		// 				Value:  "1",
		// 				Line:   1,
		// 				CharAt: 1,
		// 			},
		// 			Right: BinaryExpression{
		// 				Left: LiteralExpression{
		// 					Type:   "number",
		// 					Value:  "2",
		// 					Line:   1,
		// 					CharAt: 4,
		// 				},
		// 				Right: LiteralExpression{
		// 					Type:   "number",
		// 					Value:  "3",
		// 					Line:   1,
		// 					CharAt: 6,
		// 				},
		// 				Operator: operator.Operator{
		// 					Symbol: "+",
		// 					Line:   1,
		// 					CharAt: 5,
		// 				},
		// 				Group:   true,
		// 				Nesting: 0,
		// 				Line:    1,
		// 				CharAt:  4,
		// 			},
		// 			Operator: operator.Operator{
		// 				Symbol: "+",
		// 				Line:   1,
		// 				CharAt: 2,
		// 			},
		// 			Nesting: 1,
		// 			Line:    1,
		// 			CharAt:  1,
		// 		},
		// 		Right: LiteralExpression{
		// 			Type:   "number",
		// 			Value:  "4",
		// 			Line:   1,
		// 			CharAt: 9,
		// 		},
		// 		Operator: operator.Operator{
		// 			Symbol: "+",
		// 			Line:   1,
		// 			CharAt: 8,
		// 		},
		// 		Nesting: 2,
		// 		Line:    1,
		// 		CharAt:  1,
		// 	},
		// },
		// {
		// 	name: "parse binary expression with parantheses (1*((2+(3/4))+5))+6",
		// 	in:   "1*(2+(3/4)+5)+6",
		// 	want: BinaryExpression{
		// 		Left: BinaryExpression{
		// 			Left: LiteralExpression{
		// 				Type:   "number",
		// 				Value:  "1",
		// 				Line:   1,
		// 				CharAt: 1,
		// 			},
		// 			Right: BinaryExpression{
		// 				Left: BinaryExpression{
		// 					Left: LiteralExpression{
		// 						Type:   "number",
		// 						Value:  "2",
		// 						Line:   1,
		// 						CharAt: 4,
		// 					},
		// 					Right: BinaryExpression{
		// 						Left: LiteralExpression{
		// 							Type:   "number",
		// 							Value:  "3",
		// 							Line:   1,
		// 							CharAt: 7,
		// 						},
		// 						Right: LiteralExpression{
		// 							Type:   "number",
		// 							Value:  "4",
		// 							Line:   1,
		// 							CharAt: 9,
		// 						},
		// 						Operator: operator.Operator{
		// 							Symbol: "/",
		// 							Line:   1,
		// 							CharAt: 8,
		// 						},
		// 						Group:   true,
		// 						Nesting: 0,
		// 						Line:    1,
		// 						CharAt:  7,
		// 					},
		// 					Operator: operator.Operator{
		// 						Symbol: "+",
		// 						Line:   1,
		// 						CharAt: 5,
		// 					},
		// 					Nesting: 1,
		// 					Line:    1,
		// 					CharAt:  4,
		// 				},
		// 				Right: LiteralExpression{
		// 					Type:   "number",
		// 					Value:  "5",
		// 					Line:   1,
		// 					CharAt: 12,
		// 				},
		// 				Operator: operator.Operator{
		// 					Symbol: "+",
		// 					Line:   1,
		// 					CharAt: 11,
		// 				},
		// 				Group:   true,
		// 				Nesting: 2,
		// 				Line:    1,
		// 				CharAt:  4,
		// 			},
		// 			Operator: operator.Operator{
		// 				Symbol: "*",
		// 				Line:   1,
		// 				CharAt: 2,
		// 			},
		// 			Nesting: 3,
		// 			Line:    1,
		// 			CharAt:  1,
		// 		},
		// 		Right: LiteralExpression{
		// 			Type:   "number",
		// 			Value:  "6",
		// 			Line:   1,
		// 			CharAt: 15,
		// 		},
		// 		Operator: operator.Operator{
		// 			Symbol: "+",
		// 			Line:   1,
		// 			CharAt: 14,
		// 		},
		// 		Nesting: 4,
		// 		Line:    1,
		// 		CharAt:  1,
		// 	},
		// },
		{
			name: "parse variable expression",
			in:   "a",
			want: VariableExpression{
				Name:   "a",
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression, operand is variable expression a+b",
			in:   "a+b",
			want: BinaryExpression{
				Left: VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				Right: VariableExpression{
					Name:   "b",
					Line:   1,
					CharAt: 3,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 2,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression, operand is variable expression a+1",
			in:   "a+1",
			want: BinaryExpression{
				Left: VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				Right: LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 3,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 2,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse member access expression a.b",
			in:   "a.b",
			want: MemberExpression{
				Object: VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				PropertyIdentifier: Identifier{
					Name:   "b",
					Line:   1,
					CharAt: 3,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression, operand is member access expression (a.b)*1",
			in:   "a.b*1",
			want: BinaryExpression{
				Left: MemberExpression{
					Object: VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					PropertyIdentifier: Identifier{
						Name:   "b",
						Line:   1,
						CharAt: 3,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 5,
				},
				Operator: operator.Operator{
					Symbol: "*",
					Line:   1,
					CharAt: 4,
				},
				Line:   1,
				CharAt: 3,
			},
		},
		{
			name: "parse member access expression (a.b).c",
			in:   "a.b.c",
			want: MemberExpression{
				Object: MemberExpression{
					Object: VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					PropertyIdentifier: Identifier{
						Name:   "b",
						Line:   1,
						CharAt: 3,
					},
					Line:   1,
					CharAt: 1,
				},
				PropertyIdentifier: Identifier{
					Name:   "c",
					Line:   1,
					CharAt: 5,
				},
				Line:   1,
				CharAt: 3,
			},
		},
		// {
		// 	name: "parse binary expression, operand is variable expression (1+abc)+1",
		// 	in:   "1+abc+1",
		// 	want: BinaryExpression{
		// 		Left: BinaryExpression{
		// 			Left: LiteralExpression{
		// 				Type:   "number",
		// 				Value:  "1",
		// 				Line:   1,
		// 				CharAt: 1,
		// 			},
		// 			Right: VariableExpression{
		// 				Name:   "abc",
		// 				Line:   1,
		// 				CharAt: 3,
		// 			},
		// 			Operator: operator.Operator{
		// 				Symbol: "+",
		// 				Line:   1,
		// 				CharAt: 2,
		// 			},
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 		Right: LiteralExpression{
		// 			Type:   "number",
		// 			Value:  "1",
		// 			Line:   1,
		// 			CharAt: 7,
		// 		},
		// 		Operator: operator.Operator{
		// 			Symbol: "+",
		// 			Line:   1,
		// 			CharAt: 6,
		// 		},
		// 		Nesting: 1,
		// 		Line:    1,
		// 		CharAt:  1,
		// 	},
		// },
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp, _, _ := parseExpression(lexer.Lex(tt.in))
			require.Equal(t, tt.want, exp)
		})
	}
}
