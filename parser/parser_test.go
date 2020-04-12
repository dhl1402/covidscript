package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gs/lexer"
	"gs/operator"
)

func Test(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want Expression
		// want []Statement
	}{
		// {
		// 	name: "parse call expression #9",
		// 	in:   "[a[1],[2],b]",
		// 	want: &CallExpression{
		// 		Callee: &VariableExpression{
		// 			Name:   "a",
		// 			Line:   1,
		// 			CharAt: 1,
		// 		},
		// 		Arguments: []Expression{
		// 			&CallExpression{
		// 				Callee: &VariableExpression{
		// 					Name:   "b",
		// 					Line:   1,
		// 					CharAt: 3,
		// 				},
		// 				Arguments: []Expression{
		// 					&CallExpression{
		// 						Callee: &VariableExpression{
		// 							Name:   "c",
		// 							Line:   1,
		// 							CharAt: 5,
		// 						},
		// 						Arguments: []Expression{},
		// 						Line:      1,
		// 						CharAt:    5,
		// 					},
		// 				},
		// 				Line:   1,
		// 				CharAt: 3,
		// 			},
		// 		},
		// 		Line:   1,
		// 		CharAt: 1,
		// 	},
		// },
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp, _, _ := parseExpression(lexer.Lex(tt.in))
			require.Equal(t, tt.want, exp)
			// ast, _ := ToAST(lexer.Lex(tt.in))
			// require.Equal(t, tt.want, ast)
		})
	}
}

func TestParseExpression(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want Expression
	}{
		{
			name: "parse number expression",
			in:   "1",
			want: &LiteralExpression{
				Type:   "number",
				Value:  "1",
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse number expression",
			in:   "1 1",
			want: &LiteralExpression{
				Type:   "number",
				Value:  "1",
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse string expression",
			in:   `"1"`,
			want: &LiteralExpression{
				Type:   "string",
				Value:  `"1"`,
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse boolean expression",
			in:   "false",
			want: &LiteralExpression{
				Type:   "boolean",
				Value:  `false`,
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression",
			in:   "1+1",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &LiteralExpression{
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
		{
			name: "parse literal expression (1)",
			in:   "(1)",
			want: &LiteralExpression{
				Type:   "number",
				Value:  "1",
				Line:   1,
				CharAt: 2,
			},
		},
		{
			name: "parse literal expression (((1)))",
			in:   "(((1)))",
			want: &LiteralExpression{
				Type:   "number",
				Value:  "1",
				Line:   1,
				CharAt: 4,
			},
		},
		{
			name: "parse binary expression (1)+2",
			in:   "(1)+2",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 2,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "2",
					Line:   1,
					CharAt: 5,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 4,
				},
				Line:   1,
				CharAt: 2,
			},
		},
		{
			name: "parse binary expression (1+2)+3",
			in:   "1+2+3",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "2",
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
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "3",
					Line:   1,
					CharAt: 5,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 4,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression 1+((2*3)/4)",
			in:   "1+2*3/4",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &BinaryExpression{
					Left: &BinaryExpression{
						Left: &LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 3,
						},
						Right: &LiteralExpression{
							Type:   "number",
							Value:  "3",
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
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "4",
						Line:   1,
						CharAt: 7,
					},
					Operator: operator.Operator{
						Symbol: "/",
						Line:   1,
						CharAt: 6,
					},
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
			name: "parse binary expression ((1*2)*3)-4",
			in:   "1*2*3-4",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &BinaryExpression{
						Left: &LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 1,
						},
						Right: &LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 3,
						},
						Operator: operator.Operator{
							Symbol: "*",
							Line:   1,
							CharAt: 2,
						},
						Line:   1,
						CharAt: 1,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "3",
						Line:   1,
						CharAt: 5,
					},
					Operator: operator.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 4,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "4",
					Line:   1,
					CharAt: 7,
				},
				Operator: operator.Operator{
					Symbol: "-",
					Line:   1,
					CharAt: 6,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression ((((1+2)+3)+4)+5",
			in:   "1+2+3+4+5",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &BinaryExpression{
						Left: &BinaryExpression{
							Left: &LiteralExpression{
								Type:   "number",
								Value:  "1",
								Line:   1,
								CharAt: 1,
							},
							Right: &LiteralExpression{
								Type:   "number",
								Value:  "2",
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
						Right: &LiteralExpression{
							Type:   "number",
							Value:  "3",
							Line:   1,
							CharAt: 5,
						},
						Operator: operator.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 4,
						},
						Line:   1,
						CharAt: 1,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "4",
						Line:   1,
						CharAt: 7,
					},
					Operator: operator.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 6,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "5",
					Line:   1,
					CharAt: 9,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 8,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression (2*3)/4",
			in:   "2*3/4",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 1,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "3",
						Line:   1,
						CharAt: 3,
					},
					Operator: operator.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "4",
					Line:   1,
					CharAt: 5,
				},
				Operator: operator.Operator{
					Symbol: "/",
					Line:   1,
					CharAt: 4,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression with parantheses (1+2)*3",
			in:   "(1+2)*3",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 2,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 4,
					},
					Operator: operator.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 3,
					},
					Group:  true,
					Line:   1,
					CharAt: 2,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "3",
					Line:   1,
					CharAt: 7,
				},
				Operator: operator.Operator{
					Symbol: "*",
					Line:   1,
					CharAt: 6,
				},
				Line:   1,
				CharAt: 2,
			},
		},
		{
			name: "parse binary expression with parantheses ((1+2)+3)*4",
			in:   "(1+2+3)*4",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &BinaryExpression{
						Left: &LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 2,
						},
						Right: &LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 4,
						},
						Operator: operator.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 3,
						},
						Line:   1,
						CharAt: 2,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "3",
						Line:   1,
						CharAt: 6,
					},
					Operator: operator.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 5,
					},
					Group:  true,
					Line:   1,
					CharAt: 2,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "4",
					Line:   1,
					CharAt: 9,
				},
				Operator: operator.Operator{
					Symbol: "*",
					Line:   1,
					CharAt: 8,
				},
				Line:   1,
				CharAt: 2,
			},
		},
		{
			name: "parse binary expression with parantheses 1+(2+3)",
			in:   "1+(2+3)",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 4,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "3",
						Line:   1,
						CharAt: 6,
					},
					Operator: operator.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 5,
					},
					Group:  true,
					Line:   1,
					CharAt: 4,
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
			name: "parse binary expression with parantheses 1+(2*3)",
			in:   "1+(2*3)",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 4,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "3",
						Line:   1,
						CharAt: 6,
					},
					Operator: operator.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 5,
					},
					Group:  true,
					Line:   1,
					CharAt: 4,
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
			name: "parse binary expression with parantheses (1+(2+3))+4",
			in:   "1+(2+3)+4",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &BinaryExpression{
						Left: &LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 4,
						},
						Right: &LiteralExpression{
							Type:   "number",
							Value:  "3",
							Line:   1,
							CharAt: 6,
						},
						Operator: operator.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 5,
						},
						Group:  true,
						Line:   1,
						CharAt: 4,
					},
					Operator: operator.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "4",
					Line:   1,
					CharAt: 9,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 8,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression with parantheses 1+((2+3))*4)",
			in:   "1+(2+3)*4",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &BinaryExpression{
					Left: &BinaryExpression{
						Left: &LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 4,
						},
						Right: &LiteralExpression{
							Type:   "number",
							Value:  "3",
							Line:   1,
							CharAt: 6,
						},
						Operator: operator.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 5,
						},
						Group:  true,
						Line:   1,
						CharAt: 4,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "4",
						Line:   1,
						CharAt: 9,
					},
					Operator: operator.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 8,
					},
					Line:   1,
					CharAt: 4,
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
			name: "parse binary expression with parantheses 1+(((2+3)+4)*5)",
			in:   "1+(2+3+4)*5",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &BinaryExpression{
					Left: &BinaryExpression{
						Left: &BinaryExpression{
							Left: &LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 4,
							},
							Right: &LiteralExpression{
								Type:   "number",
								Value:  "3",
								Line:   1,
								CharAt: 6,
							},
							Operator: operator.Operator{
								Symbol: "+",
								Line:   1,
								CharAt: 5,
							},
							Line:   1,
							CharAt: 4,
						},
						Right: &LiteralExpression{
							Type:   "number",
							Value:  "4",
							Line:   1,
							CharAt: 8,
						},
						Operator: operator.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 7,
						},
						Group:  true,
						Line:   1,
						CharAt: 4,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "5",
						Line:   1,
						CharAt: 11,
					},
					Operator: operator.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 10,
					},
					Line:   1,
					CharAt: 4,
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
			name: "parse binary expression with parantheses (1+((2+3)*4))-5",
			in:   "1+((2+3)*4)-5",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &BinaryExpression{
						Left: &BinaryExpression{
							Left: &LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 5,
							},
							Right: &LiteralExpression{
								Type:   "number",
								Value:  "3",
								Line:   1,
								CharAt: 7,
							},
							Operator: operator.Operator{
								Symbol: "+",
								Line:   1,
								CharAt: 6,
							},
							Group:  true,
							Line:   1,
							CharAt: 5,
						},
						Right: &LiteralExpression{
							Type:   "number",
							Value:  "4",
							Line:   1,
							CharAt: 10,
						},
						Operator: operator.Operator{
							Symbol: "*",
							Line:   1,
							CharAt: 9,
						},
						Group:  true,
						Line:   1,
						CharAt: 5,
					},
					Operator: operator.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "5",
					Line:   1,
					CharAt: 13,
				},
				Operator: operator.Operator{
					Symbol: "-",
					Line:   1,
					CharAt: 12,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression with parantheses 1+((2+3)+4)-5",
			in:   "1+((2+3)+4)-5",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &BinaryExpression{
						Left: &BinaryExpression{
							Left: &LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 5,
							},
							Right: &LiteralExpression{
								Type:   "number",
								Value:  "3",
								Line:   1,
								CharAt: 7,
							},
							Operator: operator.Operator{
								Symbol: "+",
								Line:   1,
								CharAt: 6,
							},
							Group:  true,
							Line:   1,
							CharAt: 5,
						},
						Right: &LiteralExpression{
							Type:   "number",
							Value:  "4",
							Line:   1,
							CharAt: 10,
						},
						Operator: operator.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 9,
						},
						Group:  true,
						Line:   1,
						CharAt: 5,
					},
					Operator: operator.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "5",
					Line:   1,
					CharAt: 13,
				},
				Operator: operator.Operator{
					Symbol: "-",
					Line:   1,
					CharAt: 12,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression with parantheses (1*((2+(3/4))+5))+6",
			in:   "1*(2+(3/4)+5)+6",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &BinaryExpression{
						Left: &BinaryExpression{
							Left: &LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 4,
							},
							Right: &BinaryExpression{
								Left: &LiteralExpression{
									Type:   "number",
									Value:  "3",
									Line:   1,
									CharAt: 7,
								},
								Right: &LiteralExpression{
									Type:   "number",
									Value:  "4",
									Line:   1,
									CharAt: 9,
								},
								Operator: operator.Operator{
									Symbol: "/",
									Line:   1,
									CharAt: 8,
								},
								Group:  true,
								Line:   1,
								CharAt: 7,
							},
							Operator: operator.Operator{
								Symbol: "+",
								Line:   1,
								CharAt: 5,
							},
							Line:   1,
							CharAt: 4,
						},
						Right: &LiteralExpression{
							Type:   "number",
							Value:  "5",
							Line:   1,
							CharAt: 12,
						},
						Operator: operator.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 11,
						},
						Group:  true,
						Line:   1,
						CharAt: 4,
					},
					Operator: operator.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "6",
					Line:   1,
					CharAt: 15,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 14,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression with parantheses 1+(((2+(3/4))+5)*6)",
			in:   "1+(2+(3/4)+5)*6",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &BinaryExpression{
					Left: &BinaryExpression{
						Left: &BinaryExpression{
							Left: &LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 4,
							},
							Right: &BinaryExpression{
								Left: &LiteralExpression{
									Type:   "number",
									Value:  "3",
									Line:   1,
									CharAt: 7,
								},
								Right: &LiteralExpression{
									Type:   "number",
									Value:  "4",
									Line:   1,
									CharAt: 9,
								},
								Operator: operator.Operator{
									Symbol: "/",
									Line:   1,
									CharAt: 8,
								},
								Group:  true,
								Line:   1,
								CharAt: 7,
							},
							Operator: operator.Operator{
								Symbol: "+",
								Line:   1,
								CharAt: 5,
							},
							Line:   1,
							CharAt: 4,
						},
						Right: &LiteralExpression{
							Type:   "number",
							Value:  "5",
							Line:   1,
							CharAt: 12,
						},
						Operator: operator.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 11,
						},
						Group:  true,
						Line:   1,
						CharAt: 4,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "6",
						Line:   1,
						CharAt: 15,
					},
					Operator: operator.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 14,
					},
					Line:   1,
					CharAt: 4,
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
			name: "parse variable expression",
			in:   "a",
			want: &VariableExpression{
				Name:   "a",
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse variable expression ((a))",
			in:   "((a))",
			want: &VariableExpression{
				Name:   "a",
				Line:   1,
				CharAt: 3,
			},
		},
		{
			name: "parse binary expression, operand is variable expression a+b",
			in:   "a+b",
			want: &BinaryExpression{
				Left: &VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				Right: &VariableExpression{
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
			name: "parse binary expression, operand is variable expression (a)+b",
			in:   "(a)+b",
			want: &BinaryExpression{
				Left: &VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 2,
				},
				Right: &VariableExpression{
					Name:   "b",
					Line:   1,
					CharAt: 5,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 4,
				},
				Line:   1,
				CharAt: 2,
			},
		},
		{
			name: "parse binary expression, operand is variable expression a+1",
			in:   "a+1",
			want: &BinaryExpression{
				Left: &VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				Right: &LiteralExpression{
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
			name: "parse binary expression, operand is variable expression (1+abc)+1",
			in:   "1+abc+1",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &VariableExpression{
						Name:   "abc",
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
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 7,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 6,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse member access expression a.b",
			in:   "a.b",
			want: &MemberAccessExpression{
				Object: &VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				Property: &VariableExpression{
					Name:   "b",
					Line:   1,
					CharAt: 3,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse member access expression a.b.c",
			in:   "a.b.c",
			want: &MemberAccessExpression{
				Object: &MemberAccessExpression{
					Object: &VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					Property: &VariableExpression{
						Name:   "b",
						Line:   1,
						CharAt: 3,
					},
					Line:   1,
					CharAt: 1,
				},
				Property: &VariableExpression{
					Name:   "c",
					Line:   1,
					CharAt: 5,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse member access expression ((a.b)).c",
			in:   "((a.b)).c",
			want: &MemberAccessExpression{
				Object: &MemberAccessExpression{
					Object: &VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 3,
					},
					Property: &VariableExpression{
						Name:   "b",
						Line:   1,
						CharAt: 5,
					},
					Line:   1,
					CharAt: 3,
				},
				Property: &VariableExpression{
					Name:   "c",
					Line:   1,
					CharAt: 9,
				},
				Line:   1,
				CharAt: 3,
			},
		},
		{
			name: "parse binary expression, operand is member access expression (a.b)*1",
			in:   "a.b*1",
			want: &BinaryExpression{
				Left: &MemberAccessExpression{
					Object: &VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					Property: &VariableExpression{
						Name:   "b",
						Line:   1,
						CharAt: 3,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &LiteralExpression{
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
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression, operand is member access expression ((a).b).c*1",
			in:   "((a).b).c*1",
			want: &BinaryExpression{
				Left: &MemberAccessExpression{
					Object: &MemberAccessExpression{
						Object: &VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 3,
						},
						Property: &VariableExpression{
							Name:   "b",
							Line:   1,
							CharAt: 6,
						},
						Line:   1,
						CharAt: 3,
					},
					Property: &VariableExpression{
						Name:   "c",
						Line:   1,
						CharAt: 9,
					},
					Line:   1,
					CharAt: 3,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 11,
				},
				Operator: operator.Operator{
					Symbol: "*",
					Line:   1,
					CharAt: 10,
				},
				Line:   1,
				CharAt: 3,
			},
		},
		{
			name: "parse binary expression, operand is member access expression ((a).b).c+1*2",
			in:   "((a).b).c+1*2",
			want: &BinaryExpression{
				Left: &MemberAccessExpression{
					Object: &MemberAccessExpression{
						Object: &VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 3,
						},
						Property: &VariableExpression{
							Name:   "b",
							Line:   1,
							CharAt: 6,
						},
						Line:   1,
						CharAt: 3,
					},
					Property: &VariableExpression{
						Name:   "c",
						Line:   1,
						CharAt: 9,
					},
					Line:   1,
					CharAt: 3,
				},
				Right: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 11,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 13,
					},
					Operator: operator.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 12,
					},
					Line:   1,
					CharAt: 11,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 10,
				},
				Line:   1,
				CharAt: 3,
			},
		},
		{
			name: "parse binary expression, operand is member access expression ((a).b).c+(1+2)",
			in:   "((a).b).c+(1+2)",
			want: &BinaryExpression{
				Left: &MemberAccessExpression{
					Object: &MemberAccessExpression{
						Object: &VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 3,
						},
						Property: &VariableExpression{
							Name:   "b",
							Line:   1,
							CharAt: 6,
						},
						Line:   1,
						CharAt: 3,
					},
					Property: &VariableExpression{
						Name:   "c",
						Line:   1,
						CharAt: 9,
					},
					Line:   1,
					CharAt: 3,
				},
				Right: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 12,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 14,
					},
					Operator: operator.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 13,
					},
					Group:  true,
					Line:   1,
					CharAt: 12,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 10,
				},
				Line:   1,
				CharAt: 3,
			},
		},
		{
			name: "parse binary expression, operand is member access expression 1+((a.b)*2)",
			in:   "1+a.b*2",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &BinaryExpression{
					Left: &MemberAccessExpression{
						Object: &VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 3,
						},
						Property: &VariableExpression{
							Name:   "b",
							Line:   1,
							CharAt: 5,
						},
						Line:   1,
						CharAt: 3,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 7,
					},
					Operator: operator.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 6,
					},
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
			name: "parse binary expression 1+((a).b).c+2",
			in:   "1+((a).b).c+2",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &MemberAccessExpression{
						Object: &MemberAccessExpression{
							Object: &VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Property: &VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 8,
							},
							Line:   1,
							CharAt: 5,
						},
						Property: &VariableExpression{
							Name:   "c",
							Line:   1,
							CharAt: 11,
						},
						Line:   1,
						CharAt: 5,
					},
					Operator: operator.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "2",
					Line:   1,
					CharAt: 13,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 12,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression 1+((a).b).c*2",
			in:   "1+((a).b).c*2",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &BinaryExpression{
					Left: &MemberAccessExpression{
						Object: &MemberAccessExpression{
							Object: &VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Property: &VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 8,
							},
							Line:   1,
							CharAt: 5,
						},
						Property: &VariableExpression{
							Name:   "c",
							Line:   1,
							CharAt: 11,
						},
						Line:   1,
						CharAt: 5,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 13,
					},
					Operator: operator.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 12,
					},
					Line:   1,
					CharAt: 5,
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
			name: "parse binary expression 1+2+((a).b).c*2",
			in:   "1+2+((a).b).c*2",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "2",
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
				Right: &BinaryExpression{
					Left: &MemberAccessExpression{
						Object: &MemberAccessExpression{
							Object: &VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 7,
							},
							Property: &VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 10,
							},
							Line:   1,
							CharAt: 7,
						},
						Property: &VariableExpression{
							Name:   "c",
							Line:   1,
							CharAt: 13,
						},
						Line:   1,
						CharAt: 7,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 15,
					},
					Operator: operator.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 14,
					},
					Line:   1,
					CharAt: 7,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 4,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression 1+2*(((a).b).c)*3",
			in:   "1+2*(((a).b).c)*3",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &BinaryExpression{
					Left: &BinaryExpression{
						Left: &LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 3,
						},
						Right: &MemberAccessExpression{
							Object: &MemberAccessExpression{
								Object: &VariableExpression{
									Name:   "a",
									Line:   1,
									CharAt: 8,
								},
								Property: &VariableExpression{
									Name:   "b",
									Line:   1,
									CharAt: 11,
								},
								Line:   1,
								CharAt: 8,
							},
							Property: &VariableExpression{
								Name:   "c",
								Line:   1,
								CharAt: 14,
							},
							Line:   1,
							CharAt: 8,
						},
						Operator: operator.Operator{
							Symbol: "*",
							Line:   1,
							CharAt: 4,
						},
						Line:   1,
						CharAt: 3,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "3",
						Line:   1,
						CharAt: 17,
					},
					Operator: operator.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 16,
					},
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
			name: "parse binary expression, operand is member access expression 1+a.b",
			in:   "1+a.b",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &MemberAccessExpression{
					Object: &VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 3,
					},
					Property: &VariableExpression{
						Name:   "b",
						Line:   1,
						CharAt: 5,
					},
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
			name: "parse binary expression 1+((a).b).c",
			in:   "1+((a).b).c",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &MemberAccessExpression{
					Object: &MemberAccessExpression{
						Object: &VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 5,
						},
						Property: &VariableExpression{
							Name:   "b",
							Line:   1,
							CharAt: 8,
						},
						Line:   1,
						CharAt: 5,
					},
					Property: &VariableExpression{
						Name:   "c",
						Line:   1,
						CharAt: 11,
					},
					Line:   1,
					CharAt: 5,
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
			name: "parse binary expression 1+(((a).b).c)",
			in:   "1+(((a).b).c)",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &MemberAccessExpression{
					Object: &MemberAccessExpression{
						Object: &VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 6,
						},
						Property: &VariableExpression{
							Name:   "b",
							Line:   1,
							CharAt: 9,
						},
						Line:   1,
						CharAt: 6,
					},
					Property: &VariableExpression{
						Name:   "c",
						Line:   1,
						CharAt: 12,
					},
					Line:   1,
					CharAt: 6,
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
			name: "parse binary expression 1+2*(((a).b).c)",
			in:   "1+2*(((a).b).c)",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 3,
					},
					Right: &MemberAccessExpression{
						Object: &MemberAccessExpression{
							Object: &VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 8,
							},
							Property: &VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 11,
							},
							Line:   1,
							CharAt: 8,
						},
						Property: &VariableExpression{
							Name:   "c",
							Line:   1,
							CharAt: 14,
						},
						Line:   1,
						CharAt: 8,
					},
					Operator: operator.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 4,
					},
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
			name: "parse binary expression 1+2*(((a).b).c)",
			in:   "1+2+(((a).b).c)",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "2",
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
				Right: &MemberAccessExpression{
					Object: &MemberAccessExpression{
						Object: &VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 8,
						},
						Property: &VariableExpression{
							Name:   "b",
							Line:   1,
							CharAt: 11,
						},
						Line:   1,
						CharAt: 8,
					},
					Property: &VariableExpression{
						Name:   "c",
						Line:   1,
						CharAt: 14,
					},
					Line:   1,
					CharAt: 8,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 4,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression (1+((a).b).c)*2",
			in:   "(1+((a).b).c)*2",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 2,
					},
					Right: &MemberAccessExpression{
						Object: &MemberAccessExpression{
							Object: &VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 6,
							},
							Property: &VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 9,
							},
							Line:   1,
							CharAt: 6,
						},
						Property: &VariableExpression{
							Name:   "c",
							Line:   1,
							CharAt: 12,
						},
						Line:   1,
						CharAt: 6,
					},
					Operator: operator.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 3,
					},
					Group:  true,
					Line:   1,
					CharAt: 2,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "2",
					Line:   1,
					CharAt: 15,
				},
				Operator: operator.Operator{
					Symbol: "*",
					Line:   1,
					CharAt: 14,
				},
				Line:   1,
				CharAt: 2,
			},
		},
		{
			name: "parse binary expression 1+(2+((a).b).c)",
			in:   "1+(2+((a).b).c)",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 4,
					},
					Right: &MemberAccessExpression{
						Object: &MemberAccessExpression{
							Object: &VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 8,
							},
							Property: &VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 11,
							},
							Line:   1,
							CharAt: 8,
						},
						Property: &VariableExpression{
							Name:   "c",
							Line:   1,
							CharAt: 14,
						},
						Line:   1,
						CharAt: 8,
					},
					Operator: operator.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 5,
					},
					Group:  true,
					Line:   1,
					CharAt: 4,
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
			name: "parse binary expression 1+(2+(((a).b).c))",
			in:   "1+(2+(((a).b).c))+3",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &BinaryExpression{
						Left: &LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 4,
						},
						Right: &MemberAccessExpression{
							Object: &MemberAccessExpression{
								Object: &VariableExpression{
									Name:   "a",
									Line:   1,
									CharAt: 9,
								},
								Property: &VariableExpression{
									Name:   "b",
									Line:   1,
									CharAt: 12,
								},
								Line:   1,
								CharAt: 9,
							},
							Property: &VariableExpression{
								Name:   "c",
								Line:   1,
								CharAt: 15,
							},
							Line:   1,
							CharAt: 9,
						},
						Operator: operator.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 5,
						},
						Group:  true,
						Line:   1,
						CharAt: 4,
					},
					Operator: operator.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "3",
					Line:   1,
					CharAt: 19,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 18,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression 1+((((a).b).c)+2)+3",
			in:   "1+((((a).b).c)+2)+3",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &BinaryExpression{
						Left: &MemberAccessExpression{
							Object: &MemberAccessExpression{
								Object: &VariableExpression{
									Name:   "a",
									Line:   1,
									CharAt: 7,
								},
								Property: &VariableExpression{
									Name:   "b",
									Line:   1,
									CharAt: 10,
								},
								Line:   1,
								CharAt: 7,
							},
							Property: &VariableExpression{
								Name:   "c",
								Line:   1,
								CharAt: 13,
							},
							Line:   1,
							CharAt: 7,
						},
						Right: &LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 16,
						},
						Operator: operator.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 15,
						},
						Group:  true,
						Line:   1,
						CharAt: 7,
					},
					Operator: operator.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "3",
					Line:   1,
					CharAt: 19,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 18,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression a.b+c.d",
			in:   "a.b+c.d",
			want: &BinaryExpression{
				Left: &MemberAccessExpression{
					Object: &VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					Property: &VariableExpression{
						Name:   "b",
						Line:   1,
						CharAt: 3,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &MemberAccessExpression{
					Object: &VariableExpression{
						Name:   "c",
						Line:   1,
						CharAt: 5,
					},
					Property: &VariableExpression{
						Name:   "d",
						Line:   1,
						CharAt: 7,
					},
					Line:   1,
					CharAt: 5,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 4,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression ((a).b)+c.d",
			in:   "((a).b)+c.d",
			want: &BinaryExpression{
				Left: &MemberAccessExpression{
					Object: &VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 3,
					},
					Property: &VariableExpression{
						Name:   "b",
						Line:   1,
						CharAt: 6,
					},
					Line:   1,
					CharAt: 3,
				},
				Right: &MemberAccessExpression{
					Object: &VariableExpression{
						Name:   "c",
						Line:   1,
						CharAt: 9,
					},
					Property: &VariableExpression{
						Name:   "d",
						Line:   1,
						CharAt: 11,
					},
					Line:   1,
					CharAt: 9,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 8,
				},
				Line:   1,
				CharAt: 3,
			},
		},
		{
			name: "parse binary expression 1+a.b*2+c.d+3",
			in:   "1+a.b*2+c.d+3",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &BinaryExpression{
						Left: &LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 1,
						},
						Right: &BinaryExpression{
							Left: &MemberAccessExpression{
								Object: &VariableExpression{
									Name:   "a",
									Line:   1,
									CharAt: 3,
								},
								Property: &VariableExpression{
									Name:   "b",
									Line:   1,
									CharAt: 5,
								},
								Line:   1,
								CharAt: 3,
							},
							Right: &LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 7,
							},
							Operator: operator.Operator{
								Symbol: "*",
								Line:   1,
								CharAt: 6,
							},
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
					Right: &MemberAccessExpression{
						Object: &VariableExpression{
							Name:   "c",
							Line:   1,
							CharAt: 9,
						},
						Property: &VariableExpression{
							Name:   "d",
							Line:   1,
							CharAt: 11,
						},
						Line:   1,
						CharAt: 9,
					},
					Operator: operator.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 8,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "3",
					Line:   1,
					CharAt: 13,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 12,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression 1+((((a).b()).c)+2)+3",
			in:   "1+((((a).b(4+(5+6)))[(func(){})()])+2)+3",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &BinaryExpression{
						Left: &MemberAccessExpression{
							Object: &CallExpression{
								Callee: &MemberAccessExpression{
									Object: &VariableExpression{
										Name:   "a",
										Line:   1,
										CharAt: 7,
									},
									Property: &VariableExpression{
										Name:   "b",
										Line:   1,
										CharAt: 10,
									},
									Line:   1,
									CharAt: 7,
								},
								Arguments: []Expression{
									&BinaryExpression{
										Left: &LiteralExpression{
											Type:   "number",
											Value:  "4",
											Line:   1,
											CharAt: 12,
										},
										Right: &BinaryExpression{
											Left: &LiteralExpression{
												Type:   "number",
												Value:  "5",
												Line:   1,
												CharAt: 15,
											},
											Right: &LiteralExpression{
												Type:   "number",
												Value:  "6",
												Line:   1,
												CharAt: 17,
											},
											Operator: operator.Operator{
												Symbol: "+",
												Line:   1,
												CharAt: 16,
											},
											Group:  true,
											Line:   1,
											CharAt: 15,
										},
										Operator: operator.Operator{
											Symbol: "+",
											Line:   1,
											CharAt: 13,
										},
										Line:   1,
										CharAt: 12,
									},
								},
								Line:   1,
								CharAt: 7,
							},
							Property: &CallExpression{
								Callee: &FunctionExpression{
									Params: []Identifier{},
									Body:   []Statement{},
									Line:   1,
									CharAt: 23,
								},
								Arguments: []Expression{},
								Line:      1,
								CharAt:    23,
							},
							Line:   1,
							CharAt: 7,
						},
						Right: &LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 37,
						},
						Operator: operator.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 36,
						},
						Group:  true,
						Line:   1,
						CharAt: 7,
					},
					Operator: operator.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "3",
					Line:   1,
					CharAt: 40,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 39,
				},
				Line:   1,
				CharAt: 1,
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp, _, _ := parseExpression(lexer.Lex(tt.in))
			require.Equal(t, tt.want, exp)
		})
	}
}

func TestParseExpression_Object(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want Expression
	}{
		{
			name: "parse object expression {}",
			in:   "{}",
			want: &ObjectExpression{
				Properties: []ObjectProperty{},
				Line:       1,
				CharAt:     1,
			},
		},
		{
			name: "parse object expression {a:1,b:2}",
			in:   "{a:1,b:2}",
			want: &ObjectExpression{
				Properties: []ObjectProperty{
					ObjectProperty{
						KeyIdentifier: Identifier{
							Name:   "a",
							Line:   1,
							CharAt: 2,
						},
						Value: &LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 4,
						},
						Line:   1,
						CharAt: 2,
					},
					ObjectProperty{
						KeyIdentifier: Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 6,
						},
						Value: &LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 8,
						},
						Line:   1,
						CharAt: 6,
					},
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse object expression {[a+b]:1,c:2}",
			in:   "{[a+b]:1,c:2}",
			want: &ObjectExpression{
				Properties: []ObjectProperty{
					ObjectProperty{
						KeyExpression: &BinaryExpression{
							Left: &VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 3,
							},
							Right: &VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 5,
							},
							Operator: operator.Operator{
								Symbol: "+",
								Line:   1,
								CharAt: 4,
							},
							Line:   1,
							CharAt: 3,
						},
						Value: &LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 8,
						},
						Computed: true,
						Line:     1,
						CharAt:   2,
					},
					ObjectProperty{
						KeyIdentifier: Identifier{
							Name:   "c",
							Line:   1,
							CharAt: 10,
						},
						Value: &LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 12,
						},
						Line:   1,
						CharAt: 10,
					},
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse object expression {a:{c:1},b:2}",
			in:   "{a:{c:1},b:2}",
			want: &ObjectExpression{
				Properties: []ObjectProperty{
					ObjectProperty{
						KeyIdentifier: Identifier{
							Name:   "a",
							Line:   1,
							CharAt: 2,
						},
						Value: &ObjectExpression{
							Properties: []ObjectProperty{
								ObjectProperty{
									KeyIdentifier: Identifier{
										Name:   "c",
										Line:   1,
										CharAt: 5,
									},
									Value: &LiteralExpression{
										Type:   "number",
										Value:  "1",
										Line:   1,
										CharAt: 7,
									},
									Line:   1,
									CharAt: 5,
								},
							},
							Line:   1,
							CharAt: 4,
						},
						Line:   1,
						CharAt: 2,
					},
					ObjectProperty{
						KeyIdentifier: Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 10,
						},
						Value: &LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 12,
						},
						Line:   1,
						CharAt: 10,
					},
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse object expression {a:{c:1},b:{d:2}}",
			in:   "{a:{c:1},b:{d:2}}",
			want: &ObjectExpression{
				Properties: []ObjectProperty{
					ObjectProperty{
						KeyIdentifier: Identifier{
							Name:   "a",
							Line:   1,
							CharAt: 2,
						},
						Value: &ObjectExpression{
							Properties: []ObjectProperty{
								ObjectProperty{
									KeyIdentifier: Identifier{
										Name:   "c",
										Line:   1,
										CharAt: 5,
									},
									Value: &LiteralExpression{
										Type:   "number",
										Value:  "1",
										Line:   1,
										CharAt: 7,
									},
									Line:   1,
									CharAt: 5,
								},
							},
							Line:   1,
							CharAt: 4,
						},
						Line:   1,
						CharAt: 2,
					},
					ObjectProperty{
						KeyIdentifier: Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 10,
						},
						Value: &ObjectExpression{
							Properties: []ObjectProperty{
								ObjectProperty{
									KeyIdentifier: Identifier{
										Name:   "d",
										Line:   1,
										CharAt: 13,
									},
									Value: &LiteralExpression{
										Type:   "number",
										Value:  "2",
										Line:   1,
										CharAt: 15,
									},
									Line:   1,
									CharAt: 13,
								},
							},
							Line:   1,
							CharAt: 12,
						},
						Line:   1,
						CharAt: 10,
					},
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse object expression {a:{c:1},b:{d:2}}",
			in: `{
						   a:{c:1+a.b*2+c.d+3},
						   b:{d:2}
						 }`,
			want: &ObjectExpression{
				Properties: []ObjectProperty{
					ObjectProperty{
						KeyIdentifier: Identifier{
							Name:   "a",
							Line:   2,
							CharAt: 1,
						},
						Value: &ObjectExpression{
							Properties: []ObjectProperty{
								ObjectProperty{
									KeyIdentifier: Identifier{
										Name:   "c",
										Line:   2,
										CharAt: 4,
									},
									Value: &BinaryExpression{
										Left: &BinaryExpression{
											Left: &BinaryExpression{
												Left: &LiteralExpression{
													Type:   "number",
													Value:  "1",
													Line:   2,
													CharAt: 6,
												},
												Right: &BinaryExpression{
													Left: &MemberAccessExpression{
														Object: &VariableExpression{
															Name:   "a",
															Line:   2,
															CharAt: 8,
														},
														Property: &VariableExpression{
															Name:   "b",
															Line:   2,
															CharAt: 10,
														},
														Line:   2,
														CharAt: 8,
													},
													Right: &LiteralExpression{
														Type:   "number",
														Value:  "2",
														Line:   2,
														CharAt: 12,
													},
													Operator: operator.Operator{
														Symbol: "*",
														Line:   2,
														CharAt: 11,
													},
													Line:   2,
													CharAt: 8,
												},
												Operator: operator.Operator{
													Symbol: "+",
													Line:   2,
													CharAt: 7,
												},
												Line:   2,
												CharAt: 6,
											},
											Right: &MemberAccessExpression{
												Object: &VariableExpression{
													Name:   "c",
													Line:   2,
													CharAt: 14,
												},
												Property: &VariableExpression{
													Name:   "d",
													Line:   2,
													CharAt: 16,
												},
												Line:   2,
												CharAt: 14,
											},
											Operator: operator.Operator{
												Symbol: "+",
												Line:   2,
												CharAt: 13,
											},
											Line:   2,
											CharAt: 6,
										},
										Right: &LiteralExpression{
											Type:   "number",
											Value:  "3",
											Line:   2,
											CharAt: 18,
										},
										Operator: operator.Operator{
											Symbol: "+",
											Line:   2,
											CharAt: 17,
										},
										Line:   2,
										CharAt: 6,
									},
									Line:   2,
									CharAt: 4,
								},
							},
							Line:   2,
							CharAt: 3,
						},
						Line:   2,
						CharAt: 1,
					},
					ObjectProperty{
						KeyIdentifier: Identifier{
							Name:   "b",
							Line:   3,
							CharAt: 1,
						},
						Value: &ObjectExpression{
							Properties: []ObjectProperty{
								ObjectProperty{
									KeyIdentifier: Identifier{
										Name:   "d",
										Line:   3,
										CharAt: 4,
									},
									Value: &LiteralExpression{
										Type:   "number",
										Value:  "2",
										Line:   3,
										CharAt: 6,
									},
									Line:   3,
									CharAt: 4,
								},
							},
							Line:   3,
							CharAt: 3,
						},
						Line:   3,
						CharAt: 1,
					},
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse object expression {a:{c:1},b:[2,3]}",
			in:   "{a:{c:1},b:[2,3]}",
			want: &ObjectExpression{
				Properties: []ObjectProperty{
					ObjectProperty{
						KeyIdentifier: Identifier{
							Name:   "a",
							Line:   1,
							CharAt: 2,
						},
						Value: &ObjectExpression{
							Properties: []ObjectProperty{
								ObjectProperty{
									KeyIdentifier: Identifier{
										Name:   "c",
										Line:   1,
										CharAt: 5,
									},
									Value: &LiteralExpression{
										Type:   "number",
										Value:  "1",
										Line:   1,
										CharAt: 7,
									},
									Line:   1,
									CharAt: 5,
								},
							},
							Line:   1,
							CharAt: 4,
						},
						Line:   1,
						CharAt: 2,
					},
					ObjectProperty{
						KeyIdentifier: Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 10,
						},
						Value: &ArrayExpression{
							Elements: []Expression{
								&LiteralExpression{
									Type:   "number",
									Value:  "2",
									Line:   1,
									CharAt: 13,
								},
								&LiteralExpression{
									Type:   "number",
									Value:  "3",
									Line:   1,
									CharAt: 15,
								},
							},
							Line:   1,
							CharAt: 12,
						},
						Line:   1,
						CharAt: 10,
					},
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression 1+({a:1,b:2})",
			in:   "1+({a:1,b:2})",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &ObjectExpression{
					Properties: []ObjectProperty{
						ObjectProperty{
							KeyIdentifier: Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Value: &LiteralExpression{
								Type:   "number",
								Value:  "1",
								Line:   1,
								CharAt: 7,
							},
							Line:   1,
							CharAt: 5,
						},
						ObjectProperty{
							KeyIdentifier: Identifier{
								Name:   "b",
								Line:   1,
								CharAt: 9,
							},
							Value: &LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 11,
							},
							Line:   1,
							CharAt: 9,
						},
					},
					Line:   1,
					CharAt: 4,
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
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp, _, _ := parseExpression(lexer.Lex(tt.in))
			require.Equal(t, tt.want, exp)
		})
	}
}

func TestParseExpression_Array(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want Expression
	}{
		{
			name: "parse object expression []",
			in:   "[]",
			want: &ArrayExpression{
				Elements: []Expression{},
				Line:     1,
				CharAt:   1,
			},
		},
		{
			name: "parse object expression [1,a]",
			in:   "[1,a]",
			want: &ArrayExpression{
				Elements: []Expression{
					&LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 2,
					},
					&VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 4,
					},
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse object expression [1,[a,b]]",
			in:   "[1,[a,b]]",
			want: &ArrayExpression{
				Elements: []Expression{
					&LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 2,
					},
					&ArrayExpression{
						Elements: []Expression{
							&VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							&VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 7,
							},
						},
						Line:   1,
						CharAt: 4,
					},
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse object expression [1,[a,b]]",
			in:   "[1,{a:1,b:2}]",
			want: &ArrayExpression{
				Elements: []Expression{
					&LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 2,
					},
					&ObjectExpression{
						Properties: []ObjectProperty{
							ObjectProperty{
								KeyIdentifier: Identifier{
									Name:   "a",
									Line:   1,
									CharAt: 5,
								},
								Value: &LiteralExpression{
									Type:   "number",
									Value:  "1",
									Line:   1,
									CharAt: 7,
								},
								Line:   1,
								CharAt: 5,
							},
							ObjectProperty{
								KeyIdentifier: Identifier{
									Name:   "b",
									Line:   1,
									CharAt: 9,
								},
								Value: &LiteralExpression{
									Type:   "number",
									Value:  "2",
									Line:   1,
									CharAt: 11,
								},
								Line:   1,
								CharAt: 9,
							},
						},
						Line:   1,
						CharAt: 4,
					},
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse object expression [1,[a,1+a.b*2+c.d+3]]",
			in:   "[1,[a,1+a.b*2+c.d+3]]",
			want: &ArrayExpression{
				Elements: []Expression{
					&LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 2,
					},
					&ArrayExpression{
						Elements: []Expression{
							&VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							&BinaryExpression{
								Left: &BinaryExpression{
									Left: &BinaryExpression{
										Left: &LiteralExpression{
											Type:   "number",
											Value:  "1",
											Line:   1,
											CharAt: 7,
										},
										Right: &BinaryExpression{
											Left: &MemberAccessExpression{
												Object: &VariableExpression{
													Name:   "a",
													Line:   1,
													CharAt: 9,
												},
												Property: &VariableExpression{
													Name:   "b",
													Line:   1,
													CharAt: 11,
												},
												Line:   1,
												CharAt: 9,
											},
											Right: &LiteralExpression{
												Type:   "number",
												Value:  "2",
												Line:   1,
												CharAt: 13,
											},
											Operator: operator.Operator{
												Symbol: "*",
												Line:   1,
												CharAt: 12,
											},
											Line:   1,
											CharAt: 9,
										},
										Operator: operator.Operator{
											Symbol: "+",
											Line:   1,
											CharAt: 8,
										},
										Line:   1,
										CharAt: 7,
									},
									Right: &MemberAccessExpression{
										Object: &VariableExpression{
											Name:   "c",
											Line:   1,
											CharAt: 15,
										},
										Property: &VariableExpression{
											Name:   "d",
											Line:   1,
											CharAt: 17,
										},
										Line:   1,
										CharAt: 15,
									},
									Operator: operator.Operator{
										Symbol: "+",
										Line:   1,
										CharAt: 14,
									},
									Line:   1,
									CharAt: 7,
								},
								Right: &LiteralExpression{
									Type:   "number",
									Value:  "3",
									Line:   1,
									CharAt: 19,
								},
								Operator: operator.Operator{
									Symbol: "+",
									Line:   1,
									CharAt: 18,
								},
								Line:   1,
								CharAt: 7,
							},
						},
						Line:   1,
						CharAt: 4,
					},
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression a+([b])",
			in:   "a+([b])",
			want: &BinaryExpression{
				Left: &VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				Right: &ArrayExpression{
					Elements: []Expression{
						&VariableExpression{
							Name:   "b",
							Line:   1,
							CharAt: 5,
						},
					},
					Line:   1,
					CharAt: 4,
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
			name: "parse member access expression a[b]",
			in:   "a[b]",
			want: &MemberAccessExpression{
				Object: &VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				Property: &VariableExpression{
					Name:   "b",
					Line:   1,
					CharAt: 3,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse member access expression a[b][c]",
			in:   "a[b][c]",
			want: &MemberAccessExpression{
				Object: &MemberAccessExpression{
					Object: &VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					Property: &VariableExpression{
						Name:   "b",
						Line:   1,
						CharAt: 3,
					},
					Line:   1,
					CharAt: 1,
				},
				Property: &VariableExpression{
					Name:   "c",
					Line:   1,
					CharAt: 6,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse member access expression ((a[b]))[c]",
			in:   "((a[b]))[c]",
			want: &MemberAccessExpression{
				Object: &MemberAccessExpression{
					Object: &VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 3,
					},
					Property: &VariableExpression{
						Name:   "b",
						Line:   1,
						CharAt: 5,
					},
					Line:   1,
					CharAt: 3,
				},
				Property: &VariableExpression{
					Name:   "c",
					Line:   1,
					CharAt: 10,
				},
				Line:   1,
				CharAt: 3,
			},
		},
		{
			name: "parse member access expression a.b[c]",
			in:   "a.b[c]",
			want: &MemberAccessExpression{
				Object: &MemberAccessExpression{
					Object: &VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					Property: &VariableExpression{
						Name:   "b",
						Line:   1,
						CharAt: 3,
					},
					Line:   1,
					CharAt: 1,
				},
				Property: &VariableExpression{
					Name:   "c",
					Line:   1,
					CharAt: 5,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse member access expression (a.b)[c]",
			in:   "(a.b)[c]",
			want: &MemberAccessExpression{
				Object: &MemberAccessExpression{
					Object: &VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 2,
					},
					Property: &VariableExpression{
						Name:   "b",
						Line:   1,
						CharAt: 4,
					},
					Line:   1,
					CharAt: 2,
				},
				Property: &VariableExpression{
					Name:   "c",
					Line:   1,
					CharAt: 7,
				},
				Line:   1,
				CharAt: 2,
			},
		},
		{
			name: "parse member access expression a[b[c]]",
			in:   "a[b[c]]",
			want: &MemberAccessExpression{
				Object: &VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				Property: &MemberAccessExpression{
					Object: &VariableExpression{
						Name:   "b",
						Line:   1,
						CharAt: 3,
					},
					Property: &VariableExpression{
						Name:   "c",
						Line:   1,
						CharAt: 5,
					},
					Line:   1,
					CharAt: 3,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression a[b[c]]+1",
			in:   "a[b[c]]+1",
			want: &BinaryExpression{
				Left: &MemberAccessExpression{
					Object: &VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					Property: &MemberAccessExpression{
						Object: &VariableExpression{
							Name:   "b",
							Line:   1,
							CharAt: 3,
						},
						Property: &VariableExpression{
							Name:   "c",
							Line:   1,
							CharAt: 5,
						},
						Line:   1,
						CharAt: 3,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 9,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 8,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression 1+a[b[c]]",
			in:   "1+a[b[c]]",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &MemberAccessExpression{
					Object: &VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 3,
					},
					Property: &MemberAccessExpression{
						Object: &VariableExpression{
							Name:   "b",
							Line:   1,
							CharAt: 5,
						},
						Property: &VariableExpression{
							Name:   "c",
							Line:   1,
							CharAt: 7,
						},
						Line:   1,
						CharAt: 5,
					},
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
			name: "parse binary expression 1+a[b[c]]+2",
			in:   "1+a[b[c]]+2",
			want: &BinaryExpression{
				Left: &BinaryExpression{
					Left: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &MemberAccessExpression{
						Object: &VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 3,
						},
						Property: &MemberAccessExpression{
							Object: &VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 5,
							},
							Property: &VariableExpression{
								Name:   "c",
								Line:   1,
								CharAt: 7,
							},
							Line:   1,
							CharAt: 5,
						},
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
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "2",
					Line:   1,
					CharAt: 11,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 10,
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression 1+a[b[c]]*2",
			in:   "1+a[b[c]]*2",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &BinaryExpression{
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 11,
					},
					Left: &MemberAccessExpression{
						Object: &VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 3,
						},
						Property: &MemberAccessExpression{
							Object: &VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 5,
							},
							Property: &VariableExpression{
								Name:   "c",
								Line:   1,
								CharAt: 7,
							},
							Line:   1,
							CharAt: 5,
						},
						Line:   1,
						CharAt: 3,
					},
					Operator: operator.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 10,
					},
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
			name: "parse binary expression 1+((a[b[c]])+2)",
			in:   "1+((a[b[c]])+2)",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &BinaryExpression{
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 14,
					},
					Left: &MemberAccessExpression{
						Object: &VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 5,
						},
						Property: &MemberAccessExpression{
							Object: &VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 7,
							},
							Property: &VariableExpression{
								Name:   "c",
								Line:   1,
								CharAt: 9,
							},
							Line:   1,
							CharAt: 7,
						},
						Line:   1,
						CharAt: 5,
					},
					Operator: operator.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 13,
					},
					Group:  true,
					Line:   1,
					CharAt: 5,
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
			name: "parse binary expression 1+((a[b[c]])+2)*3",
			in:   "1+((a[b[c]])+2)*3",
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &BinaryExpression{
					Left: &BinaryExpression{
						Right: &LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 14,
						},
						Left: &MemberAccessExpression{
							Object: &VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Property: &MemberAccessExpression{
								Object: &VariableExpression{
									Name:   "b",
									Line:   1,
									CharAt: 7,
								},
								Property: &VariableExpression{
									Name:   "c",
									Line:   1,
									CharAt: 9,
								},
								Line:   1,
								CharAt: 7,
							},
							Line:   1,
							CharAt: 5,
						},
						Operator: operator.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 13,
						},
						Group:  true,
						Line:   1,
						CharAt: 5,
					},
					Right: &LiteralExpression{
						Type:   "number",
						Value:  "3",
						Line:   1,
						CharAt: 17,
					},
					Operator: operator.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 16,
					},
					Line:   1,
					CharAt: 5,
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
			name: "parse binary expression 1+(func (){})",
			in:   `1+(func (){})`,
			want: &BinaryExpression{
				Left: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &FunctionExpression{
					Params: []Identifier{},
					Body:   []Statement{},
					Line:   1,
					CharAt: 4,
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
			name: "parse binary expression (func (){})+1",
			in:   `(func (){})+1`,
			want: &BinaryExpression{
				Right: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 13,
				},
				Left: &FunctionExpression{
					Params: []Identifier{},
					Body:   []Statement{},
					Line:   1,
					CharAt: 2,
				},
				Operator: operator.Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 12,
				},
				Line:   1,
				CharAt: 2,
			},
		},
		{
			name: "parse member access expression (func (){}).abc",
			in:   `(func (){}).abc`,
			want: &MemberAccessExpression{
				Object: &FunctionExpression{
					Params: []Identifier{},
					Body:   []Statement{},
					Line:   1,
					CharAt: 2,
				},
				Property: &VariableExpression{
					Name:   "abc",
					Line:   1,
					CharAt: 13,
				},
				Line:   1,
				CharAt: 2,
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp, _, _ := parseExpression(lexer.Lex(tt.in))
			require.Equal(t, tt.want, exp)
		})
	}
}

func TestToAST_VariableDeclaration(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []Statement
	}{
		{
			name: "parse variable declaration statement (without initialization)",
			in:   `var a`,
			want: []Statement{
				VariableDeclaration{
					Declarations: []VariableDeclarator{
						VariableDeclarator{
							ID: Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init:   nil,
							Line:   1,
							CharAt: 5,
						},
					},
					Line:   1,
					CharAt: 1,
				},
			},
		},
		{
			name: "parse variable declaration statement (without initialization)",
			in: `var a
		 		 var b`,
			want: []Statement{
				VariableDeclaration{
					Declarations: []VariableDeclarator{
						VariableDeclarator{
							ID: Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init:   nil,
							Line:   1,
							CharAt: 5,
						},
					},
					Line:   1,
					CharAt: 1,
				},
				VariableDeclaration{
					Declarations: []VariableDeclarator{
						VariableDeclarator{
							ID: Identifier{
								Name:   "b",
								Line:   2,
								CharAt: 5,
							},
							Init:   nil,
							Line:   2,
							CharAt: 5,
						},
					},
					Line:   2,
					CharAt: 1,
				},
			},
		},
		{
			name: "parse number declaration statement",
			in:   "var a=1",
			want: []Statement{
				VariableDeclaration{
					Declarations: []VariableDeclarator{
						VariableDeclarator{
							ID: Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init: &LiteralExpression{
								Type:   "number",
								Value:  "1",
								Line:   1,
								CharAt: 7,
							},
							Line:   1,
							CharAt: 5,
						},
					},
					Line:   1,
					CharAt: 1,
				},
			},
		},
		{
			name: "parse string declaration statement",
			in:   `var a="xxx"`,
			want: []Statement{
				VariableDeclaration{
					Declarations: []VariableDeclarator{
						VariableDeclarator{
							ID: Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init: &LiteralExpression{
								Type:   "string",
								Value:  `"xxx"`,
								Line:   1,
								CharAt: 7,
							},
							Line:   1,
							CharAt: 5,
						},
					},
					Line:   1,
					CharAt: 1,
				},
			},
		},
		{
			name: "parse boolean declaration statement",
			in:   "var a=false",
			want: []Statement{
				VariableDeclaration{
					Declarations: []VariableDeclarator{
						VariableDeclarator{
							ID: Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init: &LiteralExpression{
								Type:   "boolean",
								Value:  "false",
								Line:   1,
								CharAt: 7,
							},
							Line:   1,
							CharAt: 5,
						},
					},
					Line:   1,
					CharAt: 1,
				},
			},
		},
		{
			name: "parse array declaration statement",
			in: `var a=[
					123,
					"456",
					1+1
				 ]`,
			want: []Statement{
				VariableDeclaration{
					Declarations: []VariableDeclarator{
						VariableDeclarator{
							ID: Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init: &ArrayExpression{
								Elements: []Expression{
									&LiteralExpression{
										Type:   "number",
										Value:  "123",
										Line:   2,
										CharAt: 1,
									},
									&LiteralExpression{
										Type:   "string",
										Value:  `"456"`,
										Line:   3,
										CharAt: 1,
									},
									&BinaryExpression{
										Left: &LiteralExpression{
											Type:   "number",
											Value:  "1",
											Line:   4,
											CharAt: 1,
										},
										Right: &LiteralExpression{
											Type:   "number",
											Value:  "1",
											Line:   4,
											CharAt: 3,
										},
										Operator: operator.Operator{
											Symbol: "+",
											Line:   4,
											CharAt: 2,
										},
										Line:   4,
										CharAt: 1,
									},
								},
								Line:   1,
								CharAt: 7,
							},
							Line:   1,
							CharAt: 5,
						},
					},
					Line:   1,
					CharAt: 1,
				},
			},
		},
		{
			name: "parse object declaration statement",
			in:   `var c={a:1,b:2}`,
			want: []Statement{
				VariableDeclaration{
					Declarations: []VariableDeclarator{
						VariableDeclarator{
							ID: Identifier{
								Name:   "c",
								Line:   1,
								CharAt: 5,
							},
							Init: &ObjectExpression{
								Properties: []ObjectProperty{
									ObjectProperty{
										KeyIdentifier: Identifier{
											Name:   "a",
											Line:   1,
											CharAt: 8,
										},
										Value: &LiteralExpression{
											Type:   "number",
											Value:  "1",
											Line:   1,
											CharAt: 10,
										},
										Line:   1,
										CharAt: 8,
									},
									ObjectProperty{
										KeyIdentifier: Identifier{
											Name:   "b",
											Line:   1,
											CharAt: 12,
										},
										Value: &LiteralExpression{
											Type:   "number",
											Value:  "2",
											Line:   1,
											CharAt: 14,
										},
										Line:   1,
										CharAt: 12,
									},
								},
								Line:   1,
								CharAt: 7,
							},
							Line:   1,
							CharAt: 5,
						},
					},
					Line:   1,
					CharAt: 1,
				},
			},
		},
		{
			name: "parse variable expression",
			in: `var a=false
		         var b=a`,
			want: []Statement{
				VariableDeclaration{
					Declarations: []VariableDeclarator{
						VariableDeclarator{
							ID: Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init: &LiteralExpression{
								Type:   "boolean",
								Value:  "false",
								Line:   1,
								CharAt: 7,
							},
							Line:   1,
							CharAt: 5,
						},
					},
					Line:   1,
					CharAt: 1,
				},
				VariableDeclaration{
					Declarations: []VariableDeclarator{
						VariableDeclarator{
							ID: Identifier{
								Name:   "b",
								Line:   2,
								CharAt: 5,
							},
							Init: &VariableExpression{
								Name:   "a",
								Line:   2,
								CharAt: 7,
							},
							Line:   2,
							CharAt: 5,
						},
					},
					Line:   2,
					CharAt: 1,
				},
			},
		},
		{
			name: "parse declaration statement (same statement, multi variable)",
			in:   "var a,b=1,'2'",
			want: []Statement{
				VariableDeclaration{
					Declarations: []VariableDeclarator{
						VariableDeclarator{
							ID: Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init: &LiteralExpression{
								Type:   "number",
								Value:  "1",
								Line:   1,
								CharAt: 9,
							},
							Line:   1,
							CharAt: 5,
						},
						VariableDeclarator{
							ID: Identifier{
								Name:   "b",
								Line:   1,
								CharAt: 7,
							},
							Init: &LiteralExpression{
								Type:   "string",
								Value:  "'2'",
								Line:   1,
								CharAt: 11,
							},
							Line:   1,
							CharAt: 7,
						},
					},
					Line:   1,
					CharAt: 1,
				},
			},
		},
		{
			name: "parse declaration statement (same statement, multi variable, b is not initialized)",
			in: `var a,b=1
		         var c=2`,
			want: []Statement{
				VariableDeclaration{
					Declarations: []VariableDeclarator{
						VariableDeclarator{
							ID: Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init: &LiteralExpression{
								Type:   "number",
								Value:  "1",
								Line:   1,
								CharAt: 9,
							},
							Line:   1,
							CharAt: 5,
						},
						VariableDeclarator{
							ID: Identifier{
								Name:   "b",
								Line:   1,
								CharAt: 7,
							},
							Init:   nil,
							Line:   1,
							CharAt: 7,
						},
					},
					Line:   1,
					CharAt: 1,
				},
				VariableDeclaration{
					Declarations: []VariableDeclarator{
						VariableDeclarator{
							ID: Identifier{
								Name:   "c",
								Line:   2,
								CharAt: 5,
							},
							Init: &LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   2,
								CharAt: 7,
							},
							Line:   2,
							CharAt: 5,
						},
					},
					Line:   2,
					CharAt: 1,
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ast, _ := ToAST(lexer.Lex(tt.in))
			require.Equal(t, tt.want, ast)
		})
	}
}

func TestToAST_FunctionDeclaration(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []Statement
	}{
		{
			name: "parse function #1",
			in:   `func a(){}`,
			want: []Statement{
				FunctionDeclaration{
					ID: Identifier{
						Name:   "a",
						Line:   1,
						CharAt: 6,
					},
					Params: []Identifier{},
					Body:   []Statement{},
					Line:   1,
					CharAt: 1,
				},
			},
		},
		{
			name: "parse function #2",
			in:   `func a(b,c){}`,
			want: []Statement{
				FunctionDeclaration{
					ID: Identifier{
						Name:   "a",
						Line:   1,
						CharAt: 6,
					},
					Params: []Identifier{
						Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 8,
						},
						Identifier{
							Name:   "c",
							Line:   1,
							CharAt: 10,
						},
					},
					Body:   []Statement{},
					Line:   1,
					CharAt: 1,
				},
			},
		},
		{
			name: "parse function #3",
			in: `func a(b,c){
					var a,b=1
					var c=2
				 }`,
			want: []Statement{
				FunctionDeclaration{
					ID: Identifier{
						Name:   "a",
						Line:   1,
						CharAt: 6,
					},
					Params: []Identifier{
						Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 8,
						},
						Identifier{
							Name:   "c",
							Line:   1,
							CharAt: 10,
						},
					},
					Body: []Statement{
						VariableDeclaration{
							Declarations: []VariableDeclarator{
								VariableDeclarator{
									ID: Identifier{
										Name:   "a",
										Line:   2,
										CharAt: 5,
									},
									Init: &LiteralExpression{
										Type:   "number",
										Value:  "1",
										Line:   2,
										CharAt: 9,
									},
									Line:   2,
									CharAt: 5,
								},
								VariableDeclarator{
									ID: Identifier{
										Name:   "b",
										Line:   2,
										CharAt: 7,
									},
									Init:   nil,
									Line:   2,
									CharAt: 7,
								},
							},
							Line:   2,
							CharAt: 1,
						},
						VariableDeclaration{
							Declarations: []VariableDeclarator{
								VariableDeclarator{
									ID: Identifier{
										Name:   "c",
										Line:   3,
										CharAt: 5,
									},
									Init: &LiteralExpression{
										Type:   "number",
										Value:  "2",
										Line:   3,
										CharAt: 7,
									},
									Line:   3,
									CharAt: 5,
								},
							},
							Line:   3,
							CharAt: 1,
						},
					},
					Line:   1,
					CharAt: 1,
				},
			},
		},
		{
			name: "parse function #4",
			in: `func a(b,c){
					var a,b=1
					var c=2
					return a+b
				 }`,
			want: []Statement{
				FunctionDeclaration{
					ID: Identifier{
						Name:   "a",
						Line:   1,
						CharAt: 6,
					},
					Params: []Identifier{
						Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 8,
						},
						Identifier{
							Name:   "c",
							Line:   1,
							CharAt: 10,
						},
					},
					Body: []Statement{
						VariableDeclaration{
							Declarations: []VariableDeclarator{
								VariableDeclarator{
									ID: Identifier{
										Name:   "a",
										Line:   2,
										CharAt: 5,
									},
									Init: &LiteralExpression{
										Type:   "number",
										Value:  "1",
										Line:   2,
										CharAt: 9,
									},
									Line:   2,
									CharAt: 5,
								},
								VariableDeclarator{
									ID: Identifier{
										Name:   "b",
										Line:   2,
										CharAt: 7,
									},
									Init:   nil,
									Line:   2,
									CharAt: 7,
								},
							},
							Line:   2,
							CharAt: 1,
						},
						VariableDeclaration{
							Declarations: []VariableDeclarator{
								VariableDeclarator{
									ID: Identifier{
										Name:   "c",
										Line:   3,
										CharAt: 5,
									},
									Init: &LiteralExpression{
										Type:   "number",
										Value:  "2",
										Line:   3,
										CharAt: 7,
									},
									Line:   3,
									CharAt: 5,
								},
							},
							Line:   3,
							CharAt: 1,
						},
						ReturnStatement{
							Argument: &BinaryExpression{
								Left: &VariableExpression{
									Name:   "a",
									Line:   4,
									CharAt: 8,
								},
								Right: &VariableExpression{
									Name:   "b",
									Line:   4,
									CharAt: 10,
								},
								Operator: operator.Operator{
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
					Line:   1,
					CharAt: 1,
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ast, _ := ToAST(lexer.Lex(tt.in))
			require.Equal(t, tt.want, ast)
		})
	}
}

func TestParseExpression_Function(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want Expression
	}{
		{
			name: "parse function expression #1",
			in:   `func (b,c){}`,
			want: &FunctionExpression{
				Params: []Identifier{
					Identifier{
						Name:   "b",
						Line:   1,
						CharAt: 7,
					},
					Identifier{
						Name:   "c",
						Line:   1,
						CharAt: 9,
					},
				},
				Body:   []Statement{},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse function expression #2",
			in: `func (b,c){
					var a,b=1
					var c=2
					return a+b
				 }`,
			want: &FunctionExpression{
				Params: []Identifier{
					Identifier{
						Name:   "b",
						Line:   1,
						CharAt: 7,
					},
					Identifier{
						Name:   "c",
						Line:   1,
						CharAt: 9,
					},
				},
				Body: []Statement{
					VariableDeclaration{
						Declarations: []VariableDeclarator{
							VariableDeclarator{
								ID: Identifier{
									Name:   "a",
									Line:   2,
									CharAt: 5,
								},
								Init: &LiteralExpression{
									Type:   "number",
									Value:  "1",
									Line:   2,
									CharAt: 9,
								},
								Line:   2,
								CharAt: 5,
							},
							VariableDeclarator{
								ID: Identifier{
									Name:   "b",
									Line:   2,
									CharAt: 7,
								},
								Init:   nil,
								Line:   2,
								CharAt: 7,
							},
						},
						Line:   2,
						CharAt: 1,
					},
					VariableDeclaration{
						Declarations: []VariableDeclarator{
							VariableDeclarator{
								ID: Identifier{
									Name:   "c",
									Line:   3,
									CharAt: 5,
								},
								Init: &LiteralExpression{
									Type:   "number",
									Value:  "2",
									Line:   3,
									CharAt: 7,
								},
								Line:   3,
								CharAt: 5,
							},
						},
						Line:   3,
						CharAt: 1,
					},
					ReturnStatement{
						Argument: &BinaryExpression{
							Left: &VariableExpression{
								Name:   "a",
								Line:   4,
								CharAt: 8,
							},
							Right: &VariableExpression{
								Name:   "b",
								Line:   4,
								CharAt: 10,
							},
							Operator: operator.Operator{
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
				Line:   1,
				CharAt: 1,
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp, _, _ := parseExpression(lexer.Lex(tt.in))
			require.Equal(t, tt.want, exp)
		})
	}
}

func TestParseExpression_CallExpression(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want Expression
	}{
		{
			name: "parse call expression #1",
			in:   "a()",
			want: &CallExpression{
				Callee: &VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				Arguments: []Expression{},
				Line:      1,
				CharAt:    1,
			},
		},
		{
			name: "parse call expression #2",
			in:   "((a))(1)",
			want: &CallExpression{
				Callee: &VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 3,
				},
				Arguments: []Expression{
					&LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 7,
					},
				},
				Line:   1,
				CharAt: 3,
			},
		},
		{
			name: "parse call expression #3",
			in:   "((a))(1,(2+3))",
			want: &CallExpression{
				Callee: &VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 3,
				},
				Arguments: []Expression{
					&LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 7,
					},
					&BinaryExpression{
						Left: &LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 10,
						},
						Right: &LiteralExpression{
							Type:   "number",
							Value:  "3",
							Line:   1,
							CharAt: 12,
						},
						Operator: operator.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 11,
						},
						Group:  true,
						Line:   1,
						CharAt: 10,
					},
				},
				Line:   1,
				CharAt: 3,
			},
		},
		{
			name: "parse call expression #4",
			in:   "((a))(1,(2+3))(4)",
			want: &CallExpression{
				Callee: &CallExpression{
					Callee: &VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 3,
					},
					Arguments: []Expression{
						&LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 7,
						},
						&BinaryExpression{
							Left: &LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 10,
							},
							Right: &LiteralExpression{
								Type:   "number",
								Value:  "3",
								Line:   1,
								CharAt: 12,
							},
							Operator: operator.Operator{
								Symbol: "+",
								Line:   1,
								CharAt: 11,
							},
							Group:  true,
							Line:   1,
							CharAt: 10,
						},
					},
					Line:   1,
					CharAt: 3,
				},
				Arguments: []Expression{
					&LiteralExpression{
						Type:   "number",
						Value:  "4",
						Line:   1,
						CharAt: 16,
					},
				},
				Line:   1,
				CharAt: 3,
			},
		},
		{
			name: "parse call expression #5",
			in:   "((a(1,(2+3))))(4)",
			want: &CallExpression{
				Callee: &CallExpression{
					Callee: &VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 3,
					},
					Arguments: []Expression{
						&LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 5,
						},
						&BinaryExpression{
							Left: &LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 8,
							},
							Right: &LiteralExpression{
								Type:   "number",
								Value:  "3",
								Line:   1,
								CharAt: 10,
							},
							Operator: operator.Operator{
								Symbol: "+",
								Line:   1,
								CharAt: 9,
							},
							Group:  true,
							Line:   1,
							CharAt: 8,
						},
					},
					Line:   1,
					CharAt: 3,
				},
				Arguments: []Expression{
					&LiteralExpression{
						Type:   "number",
						Value:  "4",
						Line:   1,
						CharAt: 16,
					},
				},
				Line:   1,
				CharAt: 3,
			},
		},
		{
			name: "parse call expression #6",
			in:   "a[1](2)",
			want: &CallExpression{
				Callee: &MemberAccessExpression{
					Object: &VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					Property: &LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 3,
					},
					Line:   1,
					CharAt: 1,
				},
				Arguments: []Expression{
					&LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 6,
					},
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse call expression #7",
			in:   "[a,b(a),c]",
			want: &ArrayExpression{
				Elements: []Expression{
					&VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 2,
					},
					&CallExpression{
						Callee: &VariableExpression{
							Name:   "b",
							Line:   1,
							CharAt: 4,
						},
						Arguments: []Expression{
							&VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 6,
							},
						},
						Line:   1,
						CharAt: 4,
					},
					&VariableExpression{
						Name:   "c",
						Line:   1,
						CharAt: 9,
					},
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse call expression #8",
			in:   "{a:b(a)}",
			want: &ObjectExpression{
				Properties: []ObjectProperty{
					ObjectProperty{
						KeyIdentifier: Identifier{
							Name:   "a",
							Line:   1,
							CharAt: 2,
						},
						Value: &CallExpression{
							Callee: &VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 4,
							},
							Arguments: []Expression{
								&VariableExpression{
									Name:   "a",
									Line:   1,
									CharAt: 6,
								},
							},
							Line:   1,
							CharAt: 4,
						},
						Line:   1,
						CharAt: 2,
					},
				},
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse call expression #9",
			in:   "a(b(c()))",
			want: &CallExpression{
				Callee: &VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				Arguments: []Expression{
					&CallExpression{
						Callee: &VariableExpression{
							Name:   "b",
							Line:   1,
							CharAt: 3,
						},
						Arguments: []Expression{
							&CallExpression{
								Callee: &VariableExpression{
									Name:   "c",
									Line:   1,
									CharAt: 5,
								},
								Arguments: []Expression{},
								Line:      1,
								CharAt:    5,
							},
						},
						Line:   1,
						CharAt: 3,
					},
				},
				Line:   1,
				CharAt: 1,
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp, _, _ := parseExpression(lexer.Lex(tt.in))
			require.Equal(t, tt.want, exp)
		})
	}
}
