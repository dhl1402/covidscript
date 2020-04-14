package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gs/core"
	"gs/lexer"
)

func Test_TMP(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want core.Expression
		// want []core.Statement
	}{
		{
			name: "parse member access expression a.b.c",
			in:   "a.b.c",
			want: &core.MemberAccessExpression{
				Object: &core.MemberAccessExpression{
					Object: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					PropertyIdentifier: core.Identifier{
						Name:   "b",
						Line:   1,
						CharAt: 3,
					},
					Line:   1,
					CharAt: 1,
				},
				PropertyIdentifier: core.Identifier{
					Name:   "c",
					Line:   1,
					CharAt: 5,
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
			// ast, _ := ToAST(lexer.Lex(tt.in))
			// require.Equal(t, tt.want, ast)
		})
	}
}

func TestParseExpression(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want core.Expression
	}{
		{
			name: "parse number expression",
			in:   "1",
			want: &core.LiteralExpression{
				Type:   "number",
				Value:  "1",
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse number expression",
			in:   "1 1",
			want: &core.LiteralExpression{
				Type:   "number",
				Value:  "1",
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse string expression",
			in:   `"1"`,
			want: &core.LiteralExpression{
				Type:   "string",
				Value:  "1",
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse boolean expression",
			in:   "false",
			want: &core.LiteralExpression{
				Type:   "boolean",
				Value:  `false`,
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse binary expression",
			in:   "1+1",
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 3,
				},
				Operator: core.Operator{
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
		want core.Expression
	}{
		{
			name: "parse literal expression (1)",
			in:   "(1)",
			want: &core.LiteralExpression{
				Type:   "number",
				Value:  "1",
				Line:   1,
				CharAt: 2,
			},
		},
		{
			name: "parse literal expression (((1)))",
			in:   "(((1)))",
			want: &core.LiteralExpression{
				Type:   "number",
				Value:  "1",
				Line:   1,
				CharAt: 4,
			},
		},
		{
			name: "parse binary expression (1)+2",
			in:   "(1)+2",
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 2,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "2",
					Line:   1,
					CharAt: 5,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 3,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "3",
					Line:   1,
					CharAt: 5,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.BinaryExpression{
					Left: &core.BinaryExpression{
						Left: &core.LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 3,
						},
						Right: &core.LiteralExpression{
							Type:   "number",
							Value:  "3",
							Line:   1,
							CharAt: 5,
						},
						Operator: core.Operator{
							Symbol: "*",
							Line:   1,
							CharAt: 4,
						},
						Line:   1,
						CharAt: 3,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "4",
						Line:   1,
						CharAt: 7,
					},
					Operator: core.Operator{
						Symbol: "/",
						Line:   1,
						CharAt: 6,
					},
					Line:   1,
					CharAt: 3,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.BinaryExpression{
						Left: &core.LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 1,
						},
						Right: &core.LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 3,
						},
						Operator: core.Operator{
							Symbol: "*",
							Line:   1,
							CharAt: 2,
						},
						Line:   1,
						CharAt: 1,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "3",
						Line:   1,
						CharAt: 5,
					},
					Operator: core.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 4,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "4",
					Line:   1,
					CharAt: 7,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.BinaryExpression{
						Left: &core.BinaryExpression{
							Left: &core.LiteralExpression{
								Type:   "number",
								Value:  "1",
								Line:   1,
								CharAt: 1,
							},
							Right: &core.LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 3,
							},
							Operator: core.Operator{
								Symbol: "+",
								Line:   1,
								CharAt: 2,
							},
							Line:   1,
							CharAt: 1,
						},
						Right: &core.LiteralExpression{
							Type:   "number",
							Value:  "3",
							Line:   1,
							CharAt: 5,
						},
						Operator: core.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 4,
						},
						Line:   1,
						CharAt: 1,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "4",
						Line:   1,
						CharAt: 7,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 6,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "5",
					Line:   1,
					CharAt: 9,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 1,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "3",
						Line:   1,
						CharAt: 3,
					},
					Operator: core.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "4",
					Line:   1,
					CharAt: 5,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 2,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 4,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 3,
					},
					Group:  true,
					Line:   1,
					CharAt: 2,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "3",
					Line:   1,
					CharAt: 7,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.BinaryExpression{
						Left: &core.LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 2,
						},
						Right: &core.LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 4,
						},
						Operator: core.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 3,
						},
						Line:   1,
						CharAt: 2,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "3",
						Line:   1,
						CharAt: 6,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 5,
					},
					Group:  true,
					Line:   1,
					CharAt: 2,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "4",
					Line:   1,
					CharAt: 9,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 4,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "3",
						Line:   1,
						CharAt: 6,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 5,
					},
					Group:  true,
					Line:   1,
					CharAt: 4,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 4,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "3",
						Line:   1,
						CharAt: 6,
					},
					Operator: core.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 5,
					},
					Group:  true,
					Line:   1,
					CharAt: 4,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &core.BinaryExpression{
						Left: &core.LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 4,
						},
						Right: &core.LiteralExpression{
							Type:   "number",
							Value:  "3",
							Line:   1,
							CharAt: 6,
						},
						Operator: core.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 5,
						},
						Group:  true,
						Line:   1,
						CharAt: 4,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "4",
					Line:   1,
					CharAt: 9,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.BinaryExpression{
					Left: &core.BinaryExpression{
						Left: &core.LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 4,
						},
						Right: &core.LiteralExpression{
							Type:   "number",
							Value:  "3",
							Line:   1,
							CharAt: 6,
						},
						Operator: core.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 5,
						},
						Group:  true,
						Line:   1,
						CharAt: 4,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "4",
						Line:   1,
						CharAt: 9,
					},
					Operator: core.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 8,
					},
					Line:   1,
					CharAt: 4,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.BinaryExpression{
					Left: &core.BinaryExpression{
						Left: &core.BinaryExpression{
							Left: &core.LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 4,
							},
							Right: &core.LiteralExpression{
								Type:   "number",
								Value:  "3",
								Line:   1,
								CharAt: 6,
							},
							Operator: core.Operator{
								Symbol: "+",
								Line:   1,
								CharAt: 5,
							},
							Line:   1,
							CharAt: 4,
						},
						Right: &core.LiteralExpression{
							Type:   "number",
							Value:  "4",
							Line:   1,
							CharAt: 8,
						},
						Operator: core.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 7,
						},
						Group:  true,
						Line:   1,
						CharAt: 4,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "5",
						Line:   1,
						CharAt: 11,
					},
					Operator: core.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 10,
					},
					Line:   1,
					CharAt: 4,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &core.BinaryExpression{
						Left: &core.BinaryExpression{
							Left: &core.LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 5,
							},
							Right: &core.LiteralExpression{
								Type:   "number",
								Value:  "3",
								Line:   1,
								CharAt: 7,
							},
							Operator: core.Operator{
								Symbol: "+",
								Line:   1,
								CharAt: 6,
							},
							Group:  true,
							Line:   1,
							CharAt: 5,
						},
						Right: &core.LiteralExpression{
							Type:   "number",
							Value:  "4",
							Line:   1,
							CharAt: 10,
						},
						Operator: core.Operator{
							Symbol: "*",
							Line:   1,
							CharAt: 9,
						},
						Group:  true,
						Line:   1,
						CharAt: 5,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "5",
					Line:   1,
					CharAt: 13,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &core.BinaryExpression{
						Left: &core.BinaryExpression{
							Left: &core.LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 5,
							},
							Right: &core.LiteralExpression{
								Type:   "number",
								Value:  "3",
								Line:   1,
								CharAt: 7,
							},
							Operator: core.Operator{
								Symbol: "+",
								Line:   1,
								CharAt: 6,
							},
							Group:  true,
							Line:   1,
							CharAt: 5,
						},
						Right: &core.LiteralExpression{
							Type:   "number",
							Value:  "4",
							Line:   1,
							CharAt: 10,
						},
						Operator: core.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 9,
						},
						Group:  true,
						Line:   1,
						CharAt: 5,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "5",
					Line:   1,
					CharAt: 13,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &core.BinaryExpression{
						Left: &core.BinaryExpression{
							Left: &core.LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 4,
							},
							Right: &core.BinaryExpression{
								Left: &core.LiteralExpression{
									Type:   "number",
									Value:  "3",
									Line:   1,
									CharAt: 7,
								},
								Right: &core.LiteralExpression{
									Type:   "number",
									Value:  "4",
									Line:   1,
									CharAt: 9,
								},
								Operator: core.Operator{
									Symbol: "/",
									Line:   1,
									CharAt: 8,
								},
								Group:  true,
								Line:   1,
								CharAt: 7,
							},
							Operator: core.Operator{
								Symbol: "+",
								Line:   1,
								CharAt: 5,
							},
							Line:   1,
							CharAt: 4,
						},
						Right: &core.LiteralExpression{
							Type:   "number",
							Value:  "5",
							Line:   1,
							CharAt: 12,
						},
						Operator: core.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 11,
						},
						Group:  true,
						Line:   1,
						CharAt: 4,
					},
					Operator: core.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "6",
					Line:   1,
					CharAt: 15,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.BinaryExpression{
					Left: &core.BinaryExpression{
						Left: &core.BinaryExpression{
							Left: &core.LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 4,
							},
							Right: &core.BinaryExpression{
								Left: &core.LiteralExpression{
									Type:   "number",
									Value:  "3",
									Line:   1,
									CharAt: 7,
								},
								Right: &core.LiteralExpression{
									Type:   "number",
									Value:  "4",
									Line:   1,
									CharAt: 9,
								},
								Operator: core.Operator{
									Symbol: "/",
									Line:   1,
									CharAt: 8,
								},
								Group:  true,
								Line:   1,
								CharAt: 7,
							},
							Operator: core.Operator{
								Symbol: "+",
								Line:   1,
								CharAt: 5,
							},
							Line:   1,
							CharAt: 4,
						},
						Right: &core.LiteralExpression{
							Type:   "number",
							Value:  "5",
							Line:   1,
							CharAt: 12,
						},
						Operator: core.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 11,
						},
						Group:  true,
						Line:   1,
						CharAt: 4,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "6",
						Line:   1,
						CharAt: 15,
					},
					Operator: core.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 14,
					},
					Line:   1,
					CharAt: 4,
				},
				Operator: core.Operator{
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
			want: &core.VariableExpression{
				Name:   "a",
				Line:   1,
				CharAt: 1,
			},
		},
		{
			name: "parse variable expression ((a))",
			in:   "((a))",
			want: &core.VariableExpression{
				Name:   "a",
				Line:   1,
				CharAt: 3,
			},
		},
		{
			name: "parse binary expression, operand is variable expression a+b",
			in:   "a+b",
			want: &core.BinaryExpression{
				Left: &core.VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.VariableExpression{
					Name:   "b",
					Line:   1,
					CharAt: 3,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 2,
				},
				Right: &core.VariableExpression{
					Name:   "b",
					Line:   1,
					CharAt: 5,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 3,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &core.VariableExpression{
						Name:   "abc",
						Line:   1,
						CharAt: 3,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 7,
				},
				Operator: core.Operator{
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
			want: &core.MemberAccessExpression{
				Object: &core.VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				PropertyIdentifier: core.Identifier{
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
			want: &core.MemberAccessExpression{
				Object: &core.MemberAccessExpression{
					Object: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					PropertyIdentifier: core.Identifier{
						Name:   "b",
						Line:   1,
						CharAt: 3,
					},
					Line:   1,
					CharAt: 1,
				},
				PropertyIdentifier: core.Identifier{
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
			want: &core.MemberAccessExpression{
				Object: &core.MemberAccessExpression{
					Object: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 3,
					},
					PropertyIdentifier: core.Identifier{
						Name:   "b",
						Line:   1,
						CharAt: 5,
					},
					Line:   1,
					CharAt: 3,
				},
				PropertyIdentifier: core.Identifier{
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
			want: &core.BinaryExpression{
				Left: &core.MemberAccessExpression{
					Object: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					PropertyIdentifier: core.Identifier{
						Name:   "b",
						Line:   1,
						CharAt: 3,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 5,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.MemberAccessExpression{
					Object: &core.MemberAccessExpression{
						Object: &core.VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 3,
						},
						PropertyIdentifier: core.Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 6,
						},
						Line:   1,
						CharAt: 3,
					},
					PropertyIdentifier: core.Identifier{
						Name:   "c",
						Line:   1,
						CharAt: 9,
					},
					Line:   1,
					CharAt: 3,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 11,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.MemberAccessExpression{
					Object: &core.MemberAccessExpression{
						Object: &core.VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 3,
						},
						PropertyIdentifier: core.Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 6,
						},
						Line:   1,
						CharAt: 3,
					},
					PropertyIdentifier: core.Identifier{
						Name:   "c",
						Line:   1,
						CharAt: 9,
					},
					Line:   1,
					CharAt: 3,
				},
				Right: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 11,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 13,
					},
					Operator: core.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 12,
					},
					Line:   1,
					CharAt: 11,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.MemberAccessExpression{
					Object: &core.MemberAccessExpression{
						Object: &core.VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 3,
						},
						PropertyIdentifier: core.Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 6,
						},
						Line:   1,
						CharAt: 3,
					},
					PropertyIdentifier: core.Identifier{
						Name:   "c",
						Line:   1,
						CharAt: 9,
					},
					Line:   1,
					CharAt: 3,
				},
				Right: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 12,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 14,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 13,
					},
					Group:  true,
					Line:   1,
					CharAt: 12,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.BinaryExpression{
					Left: &core.MemberAccessExpression{
						Object: &core.VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 3,
						},
						PropertyIdentifier: core.Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 5,
						},
						Line:   1,
						CharAt: 3,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 7,
					},
					Operator: core.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 6,
					},
					Line:   1,
					CharAt: 3,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &core.MemberAccessExpression{
						Object: &core.MemberAccessExpression{
							Object: &core.VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							PropertyIdentifier: core.Identifier{
								Name:   "b",
								Line:   1,
								CharAt: 8,
							},
							Line:   1,
							CharAt: 5,
						},
						PropertyIdentifier: core.Identifier{
							Name:   "c",
							Line:   1,
							CharAt: 11,
						},
						Line:   1,
						CharAt: 5,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "2",
					Line:   1,
					CharAt: 13,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.BinaryExpression{
					Left: &core.MemberAccessExpression{
						Object: &core.MemberAccessExpression{
							Object: &core.VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							PropertyIdentifier: core.Identifier{
								Name:   "b",
								Line:   1,
								CharAt: 8,
							},
							Line:   1,
							CharAt: 5,
						},
						PropertyIdentifier: core.Identifier{
							Name:   "c",
							Line:   1,
							CharAt: 11,
						},
						Line:   1,
						CharAt: 5,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 13,
					},
					Operator: core.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 12,
					},
					Line:   1,
					CharAt: 5,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 3,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.BinaryExpression{
					Left: &core.MemberAccessExpression{
						Object: &core.MemberAccessExpression{
							Object: &core.VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 7,
							},
							PropertyIdentifier: core.Identifier{
								Name:   "b",
								Line:   1,
								CharAt: 10,
							},
							Line:   1,
							CharAt: 7,
						},
						PropertyIdentifier: core.Identifier{
							Name:   "c",
							Line:   1,
							CharAt: 13,
						},
						Line:   1,
						CharAt: 7,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 15,
					},
					Operator: core.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 14,
					},
					Line:   1,
					CharAt: 7,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.BinaryExpression{
					Left: &core.BinaryExpression{
						Left: &core.LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 3,
						},
						Right: &core.MemberAccessExpression{
							Object: &core.MemberAccessExpression{
								Object: &core.VariableExpression{
									Name:   "a",
									Line:   1,
									CharAt: 8,
								},
								PropertyIdentifier: core.Identifier{
									Name:   "b",
									Line:   1,
									CharAt: 11,
								},
								Line:   1,
								CharAt: 8,
							},
							PropertyIdentifier: core.Identifier{
								Name:   "c",
								Line:   1,
								CharAt: 14,
							},
							Line:   1,
							CharAt: 8,
						},
						Operator: core.Operator{
							Symbol: "*",
							Line:   1,
							CharAt: 4,
						},
						Line:   1,
						CharAt: 3,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "3",
						Line:   1,
						CharAt: 17,
					},
					Operator: core.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 16,
					},
					Line:   1,
					CharAt: 3,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.MemberAccessExpression{
					Object: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 3,
					},
					PropertyIdentifier: core.Identifier{
						Name:   "b",
						Line:   1,
						CharAt: 5,
					},
					Line:   1,
					CharAt: 3,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.MemberAccessExpression{
					Object: &core.MemberAccessExpression{
						Object: &core.VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 5,
						},
						PropertyIdentifier: core.Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 8,
						},
						Line:   1,
						CharAt: 5,
					},
					PropertyIdentifier: core.Identifier{
						Name:   "c",
						Line:   1,
						CharAt: 11,
					},
					Line:   1,
					CharAt: 5,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.MemberAccessExpression{
					Object: &core.MemberAccessExpression{
						Object: &core.VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 6,
						},
						PropertyIdentifier: core.Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 9,
						},
						Line:   1,
						CharAt: 6,
					},
					PropertyIdentifier: core.Identifier{
						Name:   "c",
						Line:   1,
						CharAt: 12,
					},
					Line:   1,
					CharAt: 6,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 3,
					},
					Right: &core.MemberAccessExpression{
						Object: &core.MemberAccessExpression{
							Object: &core.VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 8,
							},
							PropertyIdentifier: core.Identifier{
								Name:   "b",
								Line:   1,
								CharAt: 11,
							},
							Line:   1,
							CharAt: 8,
						},
						PropertyIdentifier: core.Identifier{
							Name:   "c",
							Line:   1,
							CharAt: 14,
						},
						Line:   1,
						CharAt: 8,
					},
					Operator: core.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 4,
					},
					Line:   1,
					CharAt: 3,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 3,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.MemberAccessExpression{
					Object: &core.MemberAccessExpression{
						Object: &core.VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 8,
						},
						PropertyIdentifier: core.Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 11,
						},
						Line:   1,
						CharAt: 8,
					},
					PropertyIdentifier: core.Identifier{
						Name:   "c",
						Line:   1,
						CharAt: 14,
					},
					Line:   1,
					CharAt: 8,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 2,
					},
					Right: &core.MemberAccessExpression{
						Object: &core.MemberAccessExpression{
							Object: &core.VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 6,
							},
							PropertyIdentifier: core.Identifier{
								Name:   "b",
								Line:   1,
								CharAt: 9,
							},
							Line:   1,
							CharAt: 6,
						},
						PropertyIdentifier: core.Identifier{
							Name:   "c",
							Line:   1,
							CharAt: 12,
						},
						Line:   1,
						CharAt: 6,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 3,
					},
					Group:  true,
					Line:   1,
					CharAt: 2,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "2",
					Line:   1,
					CharAt: 15,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 4,
					},
					Right: &core.MemberAccessExpression{
						Object: &core.MemberAccessExpression{
							Object: &core.VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 8,
							},
							PropertyIdentifier: core.Identifier{
								Name:   "b",
								Line:   1,
								CharAt: 11,
							},
							Line:   1,
							CharAt: 8,
						},
						PropertyIdentifier: core.Identifier{
							Name:   "c",
							Line:   1,
							CharAt: 14,
						},
						Line:   1,
						CharAt: 8,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 5,
					},
					Group:  true,
					Line:   1,
					CharAt: 4,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &core.BinaryExpression{
						Left: &core.LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 4,
						},
						Right: &core.MemberAccessExpression{
							Object: &core.MemberAccessExpression{
								Object: &core.VariableExpression{
									Name:   "a",
									Line:   1,
									CharAt: 9,
								},
								PropertyIdentifier: core.Identifier{
									Name:   "b",
									Line:   1,
									CharAt: 12,
								},
								Line:   1,
								CharAt: 9,
							},
							PropertyIdentifier: core.Identifier{
								Name:   "c",
								Line:   1,
								CharAt: 15,
							},
							Line:   1,
							CharAt: 9,
						},
						Operator: core.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 5,
						},
						Group:  true,
						Line:   1,
						CharAt: 4,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "3",
					Line:   1,
					CharAt: 19,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &core.BinaryExpression{
						Left: &core.MemberAccessExpression{
							Object: &core.MemberAccessExpression{
								Object: &core.VariableExpression{
									Name:   "a",
									Line:   1,
									CharAt: 7,
								},
								PropertyIdentifier: core.Identifier{
									Name:   "b",
									Line:   1,
									CharAt: 10,
								},
								Line:   1,
								CharAt: 7,
							},
							PropertyIdentifier: core.Identifier{
								Name:   "c",
								Line:   1,
								CharAt: 13,
							},
							Line:   1,
							CharAt: 7,
						},
						Right: &core.LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 16,
						},
						Operator: core.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 15,
						},
						Group:  true,
						Line:   1,
						CharAt: 7,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "3",
					Line:   1,
					CharAt: 19,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.MemberAccessExpression{
					Object: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					PropertyIdentifier: core.Identifier{
						Name:   "b",
						Line:   1,
						CharAt: 3,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.MemberAccessExpression{
					Object: &core.VariableExpression{
						Name:   "c",
						Line:   1,
						CharAt: 5,
					},
					PropertyIdentifier: core.Identifier{
						Name:   "d",
						Line:   1,
						CharAt: 7,
					},
					Line:   1,
					CharAt: 5,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.MemberAccessExpression{
					Object: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 3,
					},
					PropertyIdentifier: core.Identifier{
						Name:   "b",
						Line:   1,
						CharAt: 6,
					},
					Line:   1,
					CharAt: 3,
				},
				Right: &core.MemberAccessExpression{
					Object: &core.VariableExpression{
						Name:   "c",
						Line:   1,
						CharAt: 9,
					},
					PropertyIdentifier: core.Identifier{
						Name:   "d",
						Line:   1,
						CharAt: 11,
					},
					Line:   1,
					CharAt: 9,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.BinaryExpression{
						Left: &core.LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 1,
						},
						Right: &core.BinaryExpression{
							Left: &core.MemberAccessExpression{
								Object: &core.VariableExpression{
									Name:   "a",
									Line:   1,
									CharAt: 3,
								},
								PropertyIdentifier: core.Identifier{
									Name:   "b",
									Line:   1,
									CharAt: 5,
								},
								Line:   1,
								CharAt: 3,
							},
							Right: &core.LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 7,
							},
							Operator: core.Operator{
								Symbol: "*",
								Line:   1,
								CharAt: 6,
							},
							Line:   1,
							CharAt: 3,
						},
						Operator: core.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 2,
						},
						Line:   1,
						CharAt: 1,
					},
					Right: &core.MemberAccessExpression{
						Object: &core.VariableExpression{
							Name:   "c",
							Line:   1,
							CharAt: 9,
						},
						PropertyIdentifier: core.Identifier{
							Name:   "d",
							Line:   1,
							CharAt: 11,
						},
						Line:   1,
						CharAt: 9,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 8,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "3",
					Line:   1,
					CharAt: 13,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &core.BinaryExpression{
						Left: &core.MemberAccessExpression{
							Object: &core.CallExpression{
								Callee: &core.MemberAccessExpression{
									Object: &core.VariableExpression{
										Name:   "a",
										Line:   1,
										CharAt: 7,
									},
									PropertyIdentifier: core.Identifier{
										Name:   "b",
										Line:   1,
										CharAt: 10,
									},
									Line:   1,
									CharAt: 7,
								},
								Arguments: []core.Expression{
									&core.BinaryExpression{
										Left: &core.LiteralExpression{
											Type:   "number",
											Value:  "4",
											Line:   1,
											CharAt: 12,
										},
										Right: &core.BinaryExpression{
											Left: &core.LiteralExpression{
												Type:   "number",
												Value:  "5",
												Line:   1,
												CharAt: 15,
											},
											Right: &core.LiteralExpression{
												Type:   "number",
												Value:  "6",
												Line:   1,
												CharAt: 17,
											},
											Operator: core.Operator{
												Symbol: "+",
												Line:   1,
												CharAt: 16,
											},
											Group:  true,
											Line:   1,
											CharAt: 15,
										},
										Operator: core.Operator{
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
							PropertyExpression: &core.CallExpression{
								Callee: &core.FunctionExpression{
									Params: []core.Identifier{},
									Body:   []core.Statement{},
									Line:   1,
									CharAt: 23,
								},
								Arguments: []core.Expression{},
								Line:      1,
								CharAt:    23,
							},
							Compute: true,
							Line:    1,
							CharAt:  7,
						},
						Right: &core.LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 37,
						},
						Operator: core.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 36,
						},
						Group:  true,
						Line:   1,
						CharAt: 7,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "3",
					Line:   1,
					CharAt: 40,
				},
				Operator: core.Operator{
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
		want core.Expression
	}{
		{
			name: "parse object expression {}",
			in:   "{}",
			want: &core.ObjectExpression{
				Properties: []core.ObjectProperty{},
				Line:       1,
				CharAt:     1,
			},
		},
		{
			name: "parse object expression {a:1,b:2}",
			in:   "{a:1,b:2}",
			want: &core.ObjectExpression{
				Properties: []core.ObjectProperty{
					core.ObjectProperty{
						KeyIdentifier: core.Identifier{
							Name:   "a",
							Line:   1,
							CharAt: 2,
						},
						Value: &core.LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 4,
						},
						Line:   1,
						CharAt: 2,
					},
					core.ObjectProperty{
						KeyIdentifier: core.Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 6,
						},
						Value: &core.LiteralExpression{
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
			want: &core.ObjectExpression{
				Properties: []core.ObjectProperty{
					core.ObjectProperty{
						KeyExpression: &core.BinaryExpression{
							Left: &core.VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 3,
							},
							Right: &core.VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 5,
							},
							Operator: core.Operator{
								Symbol: "+",
								Line:   1,
								CharAt: 4,
							},
							Line:   1,
							CharAt: 3,
						},
						Value: &core.LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 8,
						},
						Computed: true,
						Line:     1,
						CharAt:   2,
					},
					core.ObjectProperty{
						KeyIdentifier: core.Identifier{
							Name:   "c",
							Line:   1,
							CharAt: 10,
						},
						Value: &core.LiteralExpression{
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
			want: &core.ObjectExpression{
				Properties: []core.ObjectProperty{
					core.ObjectProperty{
						KeyIdentifier: core.Identifier{
							Name:   "a",
							Line:   1,
							CharAt: 2,
						},
						Value: &core.ObjectExpression{
							Properties: []core.ObjectProperty{
								core.ObjectProperty{
									KeyIdentifier: core.Identifier{
										Name:   "c",
										Line:   1,
										CharAt: 5,
									},
									Value: &core.LiteralExpression{
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
					core.ObjectProperty{
						KeyIdentifier: core.Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 10,
						},
						Value: &core.LiteralExpression{
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
			want: &core.ObjectExpression{
				Properties: []core.ObjectProperty{
					core.ObjectProperty{
						KeyIdentifier: core.Identifier{
							Name:   "a",
							Line:   1,
							CharAt: 2,
						},
						Value: &core.ObjectExpression{
							Properties: []core.ObjectProperty{
								core.ObjectProperty{
									KeyIdentifier: core.Identifier{
										Name:   "c",
										Line:   1,
										CharAt: 5,
									},
									Value: &core.LiteralExpression{
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
					core.ObjectProperty{
						KeyIdentifier: core.Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 10,
						},
						Value: &core.ObjectExpression{
							Properties: []core.ObjectProperty{
								core.ObjectProperty{
									KeyIdentifier: core.Identifier{
										Name:   "d",
										Line:   1,
										CharAt: 13,
									},
									Value: &core.LiteralExpression{
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
			want: &core.ObjectExpression{
				Properties: []core.ObjectProperty{
					core.ObjectProperty{
						KeyIdentifier: core.Identifier{
							Name:   "a",
							Line:   2,
							CharAt: 1,
						},
						Value: &core.ObjectExpression{
							Properties: []core.ObjectProperty{
								core.ObjectProperty{
									KeyIdentifier: core.Identifier{
										Name:   "c",
										Line:   2,
										CharAt: 4,
									},
									Value: &core.BinaryExpression{
										Left: &core.BinaryExpression{
											Left: &core.BinaryExpression{
												Left: &core.LiteralExpression{
													Type:   "number",
													Value:  "1",
													Line:   2,
													CharAt: 6,
												},
												Right: &core.BinaryExpression{
													Left: &core.MemberAccessExpression{
														Object: &core.VariableExpression{
															Name:   "a",
															Line:   2,
															CharAt: 8,
														},
														PropertyIdentifier: core.Identifier{
															Name:   "b",
															Line:   2,
															CharAt: 10,
														},
														Line:   2,
														CharAt: 8,
													},
													Right: &core.LiteralExpression{
														Type:   "number",
														Value:  "2",
														Line:   2,
														CharAt: 12,
													},
													Operator: core.Operator{
														Symbol: "*",
														Line:   2,
														CharAt: 11,
													},
													Line:   2,
													CharAt: 8,
												},
												Operator: core.Operator{
													Symbol: "+",
													Line:   2,
													CharAt: 7,
												},
												Line:   2,
												CharAt: 6,
											},
											Right: &core.MemberAccessExpression{
												Object: &core.VariableExpression{
													Name:   "c",
													Line:   2,
													CharAt: 14,
												},
												PropertyIdentifier: core.Identifier{
													Name:   "d",
													Line:   2,
													CharAt: 16,
												},
												Line:   2,
												CharAt: 14,
											},
											Operator: core.Operator{
												Symbol: "+",
												Line:   2,
												CharAt: 13,
											},
											Line:   2,
											CharAt: 6,
										},
										Right: &core.LiteralExpression{
											Type:   "number",
											Value:  "3",
											Line:   2,
											CharAt: 18,
										},
										Operator: core.Operator{
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
					core.ObjectProperty{
						KeyIdentifier: core.Identifier{
							Name:   "b",
							Line:   3,
							CharAt: 1,
						},
						Value: &core.ObjectExpression{
							Properties: []core.ObjectProperty{
								core.ObjectProperty{
									KeyIdentifier: core.Identifier{
										Name:   "d",
										Line:   3,
										CharAt: 4,
									},
									Value: &core.LiteralExpression{
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
			want: &core.ObjectExpression{
				Properties: []core.ObjectProperty{
					core.ObjectProperty{
						KeyIdentifier: core.Identifier{
							Name:   "a",
							Line:   1,
							CharAt: 2,
						},
						Value: &core.ObjectExpression{
							Properties: []core.ObjectProperty{
								core.ObjectProperty{
									KeyIdentifier: core.Identifier{
										Name:   "c",
										Line:   1,
										CharAt: 5,
									},
									Value: &core.LiteralExpression{
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
					core.ObjectProperty{
						KeyIdentifier: core.Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 10,
						},
						Value: &core.ArrayExpression{
							Elements: []core.Expression{
								&core.LiteralExpression{
									Type:   "number",
									Value:  "2",
									Line:   1,
									CharAt: 13,
								},
								&core.LiteralExpression{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.ObjectExpression{
					Properties: []core.ObjectProperty{
						core.ObjectProperty{
							KeyIdentifier: core.Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Value: &core.LiteralExpression{
								Type:   "number",
								Value:  "1",
								Line:   1,
								CharAt: 7,
							},
							Line:   1,
							CharAt: 5,
						},
						core.ObjectProperty{
							KeyIdentifier: core.Identifier{
								Name:   "b",
								Line:   1,
								CharAt: 9,
							},
							Value: &core.LiteralExpression{
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
				Operator: core.Operator{
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

func TestParseExpression_ArrayExpression(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want core.Expression
	}{
		{
			name: "parse array expression []",
			in:   "[]",
			want: &core.ArrayExpression{
				Elements: []core.Expression{},
				Line:     1,
				CharAt:   1,
			},
		},
		{
			name: "parse array expression [1,a]",
			in:   "[1,a]",
			want: &core.ArrayExpression{
				Elements: []core.Expression{
					&core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 2,
					},
					&core.VariableExpression{
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
			name: "parse array expression [1,[a,b]]",
			in:   "[1,[a,b]]",
			want: &core.ArrayExpression{
				Elements: []core.Expression{
					&core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 2,
					},
					&core.ArrayExpression{
						Elements: []core.Expression{
							&core.VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							&core.VariableExpression{
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
			name: "parse array expression [1,[a,b]]",
			in:   "[1,{a:1,b:2}]",
			want: &core.ArrayExpression{
				Elements: []core.Expression{
					&core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 2,
					},
					&core.ObjectExpression{
						Properties: []core.ObjectProperty{
							core.ObjectProperty{
								KeyIdentifier: core.Identifier{
									Name:   "a",
									Line:   1,
									CharAt: 5,
								},
								Value: &core.LiteralExpression{
									Type:   "number",
									Value:  "1",
									Line:   1,
									CharAt: 7,
								},
								Line:   1,
								CharAt: 5,
							},
							core.ObjectProperty{
								KeyIdentifier: core.Identifier{
									Name:   "b",
									Line:   1,
									CharAt: 9,
								},
								Value: &core.LiteralExpression{
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
			name: "parse array expression [1,[a,1+a.b*2+c.d+3]]",
			in:   "[1,[a,1+a.b*2+c.d+3]]",
			want: &core.ArrayExpression{
				Elements: []core.Expression{
					&core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 2,
					},
					&core.ArrayExpression{
						Elements: []core.Expression{
							&core.VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							&core.BinaryExpression{
								Left: &core.BinaryExpression{
									Left: &core.BinaryExpression{
										Left: &core.LiteralExpression{
											Type:   "number",
											Value:  "1",
											Line:   1,
											CharAt: 7,
										},
										Right: &core.BinaryExpression{
											Left: &core.MemberAccessExpression{
												Object: &core.VariableExpression{
													Name:   "a",
													Line:   1,
													CharAt: 9,
												},
												PropertyIdentifier: core.Identifier{
													Name:   "b",
													Line:   1,
													CharAt: 11,
												},
												Line:   1,
												CharAt: 9,
											},
											Right: &core.LiteralExpression{
												Type:   "number",
												Value:  "2",
												Line:   1,
												CharAt: 13,
											},
											Operator: core.Operator{
												Symbol: "*",
												Line:   1,
												CharAt: 12,
											},
											Line:   1,
											CharAt: 9,
										},
										Operator: core.Operator{
											Symbol: "+",
											Line:   1,
											CharAt: 8,
										},
										Line:   1,
										CharAt: 7,
									},
									Right: &core.MemberAccessExpression{
										Object: &core.VariableExpression{
											Name:   "c",
											Line:   1,
											CharAt: 15,
										},
										PropertyIdentifier: core.Identifier{
											Name:   "d",
											Line:   1,
											CharAt: 17,
										},
										Line:   1,
										CharAt: 15,
									},
									Operator: core.Operator{
										Symbol: "+",
										Line:   1,
										CharAt: 14,
									},
									Line:   1,
									CharAt: 7,
								},
								Right: &core.LiteralExpression{
									Type:   "number",
									Value:  "3",
									Line:   1,
									CharAt: 19,
								},
								Operator: core.Operator{
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
			name: "parse array expression a+([b])",
			in:   "a+([b])",
			want: &core.BinaryExpression{
				Left: &core.VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.ArrayExpression{
					Elements: []core.Expression{
						&core.VariableExpression{
							Name:   "b",
							Line:   1,
							CharAt: 5,
						},
					},
					Line:   1,
					CharAt: 4,
				},
				Operator: core.Operator{
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
			want: &core.MemberAccessExpression{
				Object: &core.VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				PropertyExpression: &core.VariableExpression{
					Name:   "b",
					Line:   1,
					CharAt: 3,
				},
				Compute: true,
				Line:    1,
				CharAt:  1,
			},
		},
		{
			name: "parse member access expression a[b][c]",
			in:   "a[b][c]",
			want: &core.MemberAccessExpression{
				Object: &core.MemberAccessExpression{
					Object: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					PropertyExpression: &core.VariableExpression{
						Name:   "b",
						Line:   1,
						CharAt: 3,
					},
					Compute: true,
					Line:    1,
					CharAt:  1,
				},
				PropertyExpression: &core.VariableExpression{
					Name:   "c",
					Line:   1,
					CharAt: 6,
				},
				Compute: true,
				Line:    1,
				CharAt:  1,
			},
		},
		{
			name: "parse member access expression ((a[b]))[c]",
			in:   "((a[b]))[c]",
			want: &core.MemberAccessExpression{
				Object: &core.MemberAccessExpression{
					Object: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 3,
					},
					PropertyExpression: &core.VariableExpression{
						Name:   "b",
						Line:   1,
						CharAt: 5,
					},
					Compute: true,
					Line:    1,
					CharAt:  3,
				},
				PropertyExpression: &core.VariableExpression{
					Name:   "c",
					Line:   1,
					CharAt: 10,
				},
				Compute: true,
				Line:    1,
				CharAt:  3,
			},
		},
		{
			name: "parse member access expression a.b[c]",
			in:   "a.b[c]",
			want: &core.MemberAccessExpression{
				Object: &core.MemberAccessExpression{
					Object: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					PropertyIdentifier: core.Identifier{
						Name:   "b",
						Line:   1,
						CharAt: 3,
					},
					Line:   1,
					CharAt: 1,
				},
				PropertyExpression: &core.VariableExpression{
					Name:   "c",
					Line:   1,
					CharAt: 5,
				},
				Compute: true,
				Line:    1,
				CharAt:  1,
			},
		},
		{
			name: "parse member access expression (a.b)[c]",
			in:   "(a.b)[c]",
			want: &core.MemberAccessExpression{
				Object: &core.MemberAccessExpression{
					Object: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 2,
					},
					PropertyIdentifier: core.Identifier{
						Name:   "b",
						Line:   1,
						CharAt: 4,
					},
					Line:   1,
					CharAt: 2,
				},
				PropertyExpression: &core.VariableExpression{
					Name:   "c",
					Line:   1,
					CharAt: 7,
				},
				Compute: true,
				Line:    1,
				CharAt:  2,
			},
		},
		{
			name: "parse member access expression a[b[c]]",
			in:   "a[b[c]]",
			want: &core.MemberAccessExpression{
				Object: &core.VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				PropertyExpression: &core.MemberAccessExpression{
					Object: &core.VariableExpression{
						Name:   "b",
						Line:   1,
						CharAt: 3,
					},
					PropertyExpression: &core.VariableExpression{
						Name:   "c",
						Line:   1,
						CharAt: 5,
					},
					Compute: true,
					Line:    1,
					CharAt:  3,
				},
				Compute: true,
				Line:    1,
				CharAt:  1,
			},
		},
		{
			name: "parse binary expression a[b[c]]+1",
			in:   "a[b[c]]+1",
			want: &core.BinaryExpression{
				Left: &core.MemberAccessExpression{
					Object: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					PropertyExpression: &core.MemberAccessExpression{
						Object: &core.VariableExpression{
							Name:   "b",
							Line:   1,
							CharAt: 3,
						},
						PropertyExpression: &core.VariableExpression{
							Name:   "c",
							Line:   1,
							CharAt: 5,
						},
						Compute: true,
						Line:    1,
						CharAt:  3,
					},
					Compute: true,
					Line:    1,
					CharAt:  1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 9,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.MemberAccessExpression{
					Object: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 3,
					},
					PropertyExpression: &core.MemberAccessExpression{
						Object: &core.VariableExpression{
							Name:   "b",
							Line:   1,
							CharAt: 5,
						},
						PropertyExpression: &core.VariableExpression{
							Name:   "c",
							Line:   1,
							CharAt: 7,
						},
						Compute: true,
						Line:    1,
						CharAt:  5,
					},
					Compute: true,
					Line:    1,
					CharAt:  3,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.BinaryExpression{
					Left: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 1,
					},
					Right: &core.MemberAccessExpression{
						Object: &core.VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 3,
						},
						PropertyExpression: &core.MemberAccessExpression{
							Object: &core.VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 5,
							},
							PropertyExpression: &core.VariableExpression{
								Name:   "c",
								Line:   1,
								CharAt: 7,
							},
							Compute: true,
							Line:    1,
							CharAt:  5,
						},
						Compute: true,
						Line:    1,
						CharAt:  3,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 2,
					},
					Line:   1,
					CharAt: 1,
				},
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "2",
					Line:   1,
					CharAt: 11,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.BinaryExpression{
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 11,
					},
					Left: &core.MemberAccessExpression{
						Object: &core.VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 3,
						},
						PropertyExpression: &core.MemberAccessExpression{
							Object: &core.VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 5,
							},
							PropertyExpression: &core.VariableExpression{
								Name:   "c",
								Line:   1,
								CharAt: 7,
							},
							Compute: true,
							Line:    1,
							CharAt:  5,
						},
						Compute: true,
						Line:    1,
						CharAt:  3,
					},
					Operator: core.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 10,
					},
					Line:   1,
					CharAt: 3,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.BinaryExpression{
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "2",
						Line:   1,
						CharAt: 14,
					},
					Left: &core.MemberAccessExpression{
						Object: &core.VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 5,
						},
						PropertyExpression: &core.MemberAccessExpression{
							Object: &core.VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 7,
							},
							PropertyExpression: &core.VariableExpression{
								Name:   "c",
								Line:   1,
								CharAt: 9,
							},
							Compute: true,
							Line:    1,
							CharAt:  7,
						},
						Compute: true,
						Line:    1,
						CharAt:  5,
					},
					Operator: core.Operator{
						Symbol: "+",
						Line:   1,
						CharAt: 13,
					},
					Group:  true,
					Line:   1,
					CharAt: 5,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.BinaryExpression{
					Left: &core.BinaryExpression{
						Right: &core.LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 14,
						},
						Left: &core.MemberAccessExpression{
							Object: &core.VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							PropertyExpression: &core.MemberAccessExpression{
								Object: &core.VariableExpression{
									Name:   "b",
									Line:   1,
									CharAt: 7,
								},
								PropertyExpression: &core.VariableExpression{
									Name:   "c",
									Line:   1,
									CharAt: 9,
								},
								Compute: true,
								Line:    1,
								CharAt:  7,
							},
							Compute: true,
							Line:    1,
							CharAt:  5,
						},
						Operator: core.Operator{
							Symbol: "+",
							Line:   1,
							CharAt: 13,
						},
						Group:  true,
						Line:   1,
						CharAt: 5,
					},
					Right: &core.LiteralExpression{
						Type:   "number",
						Value:  "3",
						Line:   1,
						CharAt: 17,
					},
					Operator: core.Operator{
						Symbol: "*",
						Line:   1,
						CharAt: 16,
					},
					Line:   1,
					CharAt: 5,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Left: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 1,
				},
				Right: &core.FunctionExpression{
					Params: []core.Identifier{},
					Body:   []core.Statement{},
					Line:   1,
					CharAt: 4,
				},
				Operator: core.Operator{
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
			want: &core.BinaryExpression{
				Right: &core.LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 13,
				},
				Left: &core.FunctionExpression{
					Params: []core.Identifier{},
					Body:   []core.Statement{},
					Line:   1,
					CharAt: 2,
				},
				Operator: core.Operator{
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
			want: &core.MemberAccessExpression{
				Object: &core.FunctionExpression{
					Params: []core.Identifier{},
					Body:   []core.Statement{},
					Line:   1,
					CharAt: 2,
				},
				PropertyIdentifier: core.Identifier{
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
		want []core.Statement
	}{
		{
			name: "parse variable declaration statement (without initialization)",
			in:   `var a`,
			want: []core.Statement{
				core.VariableDeclaration{
					Declarations: []core.VariableDeclarator{
						core.VariableDeclarator{
							ID: core.Identifier{
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
			want: []core.Statement{
				core.VariableDeclaration{
					Declarations: []core.VariableDeclarator{
						core.VariableDeclarator{
							ID: core.Identifier{
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
				core.VariableDeclaration{
					Declarations: []core.VariableDeclarator{
						core.VariableDeclarator{
							ID: core.Identifier{
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
			want: []core.Statement{
				core.VariableDeclaration{
					Declarations: []core.VariableDeclarator{
						core.VariableDeclarator{
							ID: core.Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init: &core.LiteralExpression{
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
			want: []core.Statement{
				core.VariableDeclaration{
					Declarations: []core.VariableDeclarator{
						core.VariableDeclarator{
							ID: core.Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init: &core.LiteralExpression{
								Type:   "string",
								Value:  "xxx",
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
			want: []core.Statement{
				core.VariableDeclaration{
					Declarations: []core.VariableDeclarator{
						core.VariableDeclarator{
							ID: core.Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init: &core.LiteralExpression{
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
			want: []core.Statement{
				core.VariableDeclaration{
					Declarations: []core.VariableDeclarator{
						core.VariableDeclarator{
							ID: core.Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init: &core.ArrayExpression{
								Elements: []core.Expression{
									&core.LiteralExpression{
										Type:   "number",
										Value:  "123",
										Line:   2,
										CharAt: 1,
									},
									&core.LiteralExpression{
										Type:   "string",
										Value:  "456",
										Line:   3,
										CharAt: 1,
									},
									&core.BinaryExpression{
										Left: &core.LiteralExpression{
											Type:   "number",
											Value:  "1",
											Line:   4,
											CharAt: 1,
										},
										Right: &core.LiteralExpression{
											Type:   "number",
											Value:  "1",
											Line:   4,
											CharAt: 3,
										},
										Operator: core.Operator{
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
			want: []core.Statement{
				core.VariableDeclaration{
					Declarations: []core.VariableDeclarator{
						core.VariableDeclarator{
							ID: core.Identifier{
								Name:   "c",
								Line:   1,
								CharAt: 5,
							},
							Init: &core.ObjectExpression{
								Properties: []core.ObjectProperty{
									core.ObjectProperty{
										KeyIdentifier: core.Identifier{
											Name:   "a",
											Line:   1,
											CharAt: 8,
										},
										Value: &core.LiteralExpression{
											Type:   "number",
											Value:  "1",
											Line:   1,
											CharAt: 10,
										},
										Line:   1,
										CharAt: 8,
									},
									core.ObjectProperty{
										KeyIdentifier: core.Identifier{
											Name:   "b",
											Line:   1,
											CharAt: 12,
										},
										Value: &core.LiteralExpression{
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
			want: []core.Statement{
				core.VariableDeclaration{
					Declarations: []core.VariableDeclarator{
						core.VariableDeclarator{
							ID: core.Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init: &core.LiteralExpression{
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
				core.VariableDeclaration{
					Declarations: []core.VariableDeclarator{
						core.VariableDeclarator{
							ID: core.Identifier{
								Name:   "b",
								Line:   2,
								CharAt: 5,
							},
							Init: &core.VariableExpression{
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
			want: []core.Statement{
				core.VariableDeclaration{
					Declarations: []core.VariableDeclarator{
						core.VariableDeclarator{
							ID: core.Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init: &core.LiteralExpression{
								Type:   "number",
								Value:  "1",
								Line:   1,
								CharAt: 9,
							},
							Line:   1,
							CharAt: 5,
						},
						core.VariableDeclarator{
							ID: core.Identifier{
								Name:   "b",
								Line:   1,
								CharAt: 7,
							},
							Init: &core.LiteralExpression{
								Type:   "string",
								Value:  "2",
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
			want: []core.Statement{
				core.VariableDeclaration{
					Declarations: []core.VariableDeclarator{
						core.VariableDeclarator{
							ID: core.Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init: &core.LiteralExpression{
								Type:   "number",
								Value:  "1",
								Line:   1,
								CharAt: 9,
							},
							Line:   1,
							CharAt: 5,
						},
						core.VariableDeclarator{
							ID: core.Identifier{
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
				core.VariableDeclaration{
					Declarations: []core.VariableDeclarator{
						core.VariableDeclarator{
							ID: core.Identifier{
								Name:   "c",
								Line:   2,
								CharAt: 5,
							},
							Init: &core.LiteralExpression{
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
		want []core.Statement
	}{
		{
			name: "parse function #1",
			in:   `func a(){}`,
			want: []core.Statement{
				core.FunctionDeclaration{
					ID: core.Identifier{
						Name:   "a",
						Line:   1,
						CharAt: 6,
					},
					Params: []core.Identifier{},
					Body:   []core.Statement{},
					Line:   1,
					CharAt: 1,
				},
			},
		},
		{
			name: "parse function #2",
			in:   `func a(b,c){}`,
			want: []core.Statement{
				core.FunctionDeclaration{
					ID: core.Identifier{
						Name:   "a",
						Line:   1,
						CharAt: 6,
					},
					Params: []core.Identifier{
						core.Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 8,
						},
						core.Identifier{
							Name:   "c",
							Line:   1,
							CharAt: 10,
						},
					},
					Body:   []core.Statement{},
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
			want: []core.Statement{
				core.FunctionDeclaration{
					ID: core.Identifier{
						Name:   "a",
						Line:   1,
						CharAt: 6,
					},
					Params: []core.Identifier{
						core.Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 8,
						},
						core.Identifier{
							Name:   "c",
							Line:   1,
							CharAt: 10,
						},
					},
					Body: []core.Statement{
						core.VariableDeclaration{
							Declarations: []core.VariableDeclarator{
								core.VariableDeclarator{
									ID: core.Identifier{
										Name:   "a",
										Line:   2,
										CharAt: 5,
									},
									Init: &core.LiteralExpression{
										Type:   "number",
										Value:  "1",
										Line:   2,
										CharAt: 9,
									},
									Line:   2,
									CharAt: 5,
								},
								core.VariableDeclarator{
									ID: core.Identifier{
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
						core.VariableDeclaration{
							Declarations: []core.VariableDeclarator{
								core.VariableDeclarator{
									ID: core.Identifier{
										Name:   "c",
										Line:   3,
										CharAt: 5,
									},
									Init: &core.LiteralExpression{
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
			want: []core.Statement{
				core.FunctionDeclaration{
					ID: core.Identifier{
						Name:   "a",
						Line:   1,
						CharAt: 6,
					},
					Params: []core.Identifier{
						core.Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 8,
						},
						core.Identifier{
							Name:   "c",
							Line:   1,
							CharAt: 10,
						},
					},
					Body: []core.Statement{
						core.VariableDeclaration{
							Declarations: []core.VariableDeclarator{
								core.VariableDeclarator{
									ID: core.Identifier{
										Name:   "a",
										Line:   2,
										CharAt: 5,
									},
									Init: &core.LiteralExpression{
										Type:   "number",
										Value:  "1",
										Line:   2,
										CharAt: 9,
									},
									Line:   2,
									CharAt: 5,
								},
								core.VariableDeclarator{
									ID: core.Identifier{
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
						core.VariableDeclaration{
							Declarations: []core.VariableDeclarator{
								core.VariableDeclarator{
									ID: core.Identifier{
										Name:   "c",
										Line:   3,
										CharAt: 5,
									},
									Init: &core.LiteralExpression{
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
						core.ReturnStatement{
							Argument: &core.BinaryExpression{
								Left: &core.VariableExpression{
									Name:   "a",
									Line:   4,
									CharAt: 8,
								},
								Right: &core.VariableExpression{
									Name:   "b",
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
		want core.Expression
	}{
		{
			name: "parse function expression #1",
			in:   `func (b,c){}`,
			want: &core.FunctionExpression{
				Params: []core.Identifier{
					core.Identifier{
						Name:   "b",
						Line:   1,
						CharAt: 7,
					},
					core.Identifier{
						Name:   "c",
						Line:   1,
						CharAt: 9,
					},
				},
				Body:   []core.Statement{},
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
			want: &core.FunctionExpression{
				Params: []core.Identifier{
					core.Identifier{
						Name:   "b",
						Line:   1,
						CharAt: 7,
					},
					core.Identifier{
						Name:   "c",
						Line:   1,
						CharAt: 9,
					},
				},
				Body: []core.Statement{
					core.VariableDeclaration{
						Declarations: []core.VariableDeclarator{
							core.VariableDeclarator{
								ID: core.Identifier{
									Name:   "a",
									Line:   2,
									CharAt: 5,
								},
								Init: &core.LiteralExpression{
									Type:   "number",
									Value:  "1",
									Line:   2,
									CharAt: 9,
								},
								Line:   2,
								CharAt: 5,
							},
							core.VariableDeclarator{
								ID: core.Identifier{
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
					core.VariableDeclaration{
						Declarations: []core.VariableDeclarator{
							core.VariableDeclarator{
								ID: core.Identifier{
									Name:   "c",
									Line:   3,
									CharAt: 5,
								},
								Init: &core.LiteralExpression{
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
					core.ReturnStatement{
						Argument: &core.BinaryExpression{
							Left: &core.VariableExpression{
								Name:   "a",
								Line:   4,
								CharAt: 8,
							},
							Right: &core.VariableExpression{
								Name:   "b",
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
		want core.Expression
	}{
		{
			name: "parse call expression #1",
			in:   "a()",
			want: &core.CallExpression{
				Callee: &core.VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				Arguments: []core.Expression{},
				Line:      1,
				CharAt:    1,
			},
		},
		{
			name: "parse call expression #2",
			in:   "((a))(1)",
			want: &core.CallExpression{
				Callee: &core.VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 3,
				},
				Arguments: []core.Expression{
					&core.LiteralExpression{
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
			want: &core.CallExpression{
				Callee: &core.VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 3,
				},
				Arguments: []core.Expression{
					&core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 7,
					},
					&core.BinaryExpression{
						Left: &core.LiteralExpression{
							Type:   "number",
							Value:  "2",
							Line:   1,
							CharAt: 10,
						},
						Right: &core.LiteralExpression{
							Type:   "number",
							Value:  "3",
							Line:   1,
							CharAt: 12,
						},
						Operator: core.Operator{
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
			want: &core.CallExpression{
				Callee: &core.CallExpression{
					Callee: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 3,
					},
					Arguments: []core.Expression{
						&core.LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 7,
						},
						&core.BinaryExpression{
							Left: &core.LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 10,
							},
							Right: &core.LiteralExpression{
								Type:   "number",
								Value:  "3",
								Line:   1,
								CharAt: 12,
							},
							Operator: core.Operator{
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
				Arguments: []core.Expression{
					&core.LiteralExpression{
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
			want: &core.CallExpression{
				Callee: &core.CallExpression{
					Callee: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 3,
					},
					Arguments: []core.Expression{
						&core.LiteralExpression{
							Type:   "number",
							Value:  "1",
							Line:   1,
							CharAt: 5,
						},
						&core.BinaryExpression{
							Left: &core.LiteralExpression{
								Type:   "number",
								Value:  "2",
								Line:   1,
								CharAt: 8,
							},
							Right: &core.LiteralExpression{
								Type:   "number",
								Value:  "3",
								Line:   1,
								CharAt: 10,
							},
							Operator: core.Operator{
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
				Arguments: []core.Expression{
					&core.LiteralExpression{
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
			want: &core.CallExpression{
				Callee: &core.MemberAccessExpression{
					Object: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					PropertyExpression: &core.LiteralExpression{
						Type:   "number",
						Value:  "1",
						Line:   1,
						CharAt: 3,
					},
					Compute: true,
					Line:    1,
					CharAt:  1,
				},
				Arguments: []core.Expression{
					&core.LiteralExpression{
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
			want: &core.ArrayExpression{
				Elements: []core.Expression{
					&core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 2,
					},
					&core.CallExpression{
						Callee: &core.VariableExpression{
							Name:   "b",
							Line:   1,
							CharAt: 4,
						},
						Arguments: []core.Expression{
							&core.VariableExpression{
								Name:   "a",
								Line:   1,
								CharAt: 6,
							},
						},
						Line:   1,
						CharAt: 4,
					},
					&core.VariableExpression{
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
			want: &core.ObjectExpression{
				Properties: []core.ObjectProperty{
					core.ObjectProperty{
						KeyIdentifier: core.Identifier{
							Name:   "a",
							Line:   1,
							CharAt: 2,
						},
						Value: &core.CallExpression{
							Callee: &core.VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 4,
							},
							Arguments: []core.Expression{
								&core.VariableExpression{
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
			want: &core.CallExpression{
				Callee: &core.VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 1,
				},
				Arguments: []core.Expression{
					&core.CallExpression{
						Callee: &core.VariableExpression{
							Name:   "b",
							Line:   1,
							CharAt: 3,
						},
						Arguments: []core.Expression{
							&core.CallExpression{
								Callee: &core.VariableExpression{
									Name:   "c",
									Line:   1,
									CharAt: 5,
								},
								Arguments: []core.Expression{},
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

func TestToAST_ExpressionStatement(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []core.Statement
	}{
		{
			name: "parse core.ExpressionStatement #1",
			in:   "a",
			want: []core.Statement{
				core.ExpressionStatement{
					Expression: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					Line:   1,
					CharAt: 1,
				},
			},
		},
		{
			name: "parse core.ExpressionStatement #2",
			in:   "var a a",
			want: []core.Statement{
				core.VariableDeclaration{
					Declarations: []core.VariableDeclarator{
						core.VariableDeclarator{
							ID: core.Identifier{
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
				core.ExpressionStatement{
					Expression: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 7,
					},
					Line:   1,
					CharAt: 7,
				},
			},
		},
		{
			name: "parse core.ExpressionStatement #3",
			in:   "var a,b=1 a",
			want: []core.Statement{
				core.VariableDeclaration{
					Declarations: []core.VariableDeclarator{
						core.VariableDeclarator{
							ID: core.Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init: &core.LiteralExpression{
								Type:   "number",
								Value:  "1",
								Line:   1,
								CharAt: 9,
							},
							Line:   1,
							CharAt: 5,
						},
						core.VariableDeclarator{
							ID: core.Identifier{
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
				core.ExpressionStatement{
					Expression: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 11,
					},
					Line:   1,
					CharAt: 11,
				},
			},
		},
		{
			name: "parse core.ExpressionStatement #4",
			in:   "func a(){b}",
			want: []core.Statement{
				core.FunctionDeclaration{
					ID: core.Identifier{
						Name:   "a",
						Line:   1,
						CharAt: 6,
					},
					Params: []core.Identifier{},
					Body: []core.Statement{
						core.ExpressionStatement{
							Expression: &core.VariableExpression{
								Name:   "b",
								Line:   1,
								CharAt: 10,
							},
							Line:   1,
							CharAt: 10,
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

func TestToAST_AssignmentStatement(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []core.Statement
	}{
		{
			name: "parse assignment statement #1",
			in:   "a=b",
			want: []core.Statement{
				core.AssignmentStatement{
					Left: &core.VariableExpression{
						Name:   "a",
						Line:   1,
						CharAt: 1,
					},
					Right: &core.VariableExpression{
						Name:   "b",
						Line:   1,
						CharAt: 3,
					},
					Line:   1,
					CharAt: 1,
				},
			},
		},
		{
			name: "parse call expression #2",
			in: `var a=1
			     a=(b)`,
			want: []core.Statement{
				core.VariableDeclaration{
					Declarations: []core.VariableDeclarator{
						core.VariableDeclarator{
							ID: core.Identifier{
								Name:   "a",
								Line:   1,
								CharAt: 5,
							},
							Init: &core.LiteralExpression{
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
				core.AssignmentStatement{
					Left: &core.VariableExpression{
						Name:   "a",
						Line:   2,
						CharAt: 1,
					},
					Right: &core.VariableExpression{
						Name:   "b",
						Line:   2,
						CharAt: 4,
					},
					Line:   2,
					CharAt: 1,
				},
			},
		},
		{
			name: "parse assignment statement #3",
			in:   "a.b=c",
			want: []core.Statement{
				core.AssignmentStatement{
					Left: &core.MemberAccessExpression{
						Object: &core.VariableExpression{
							Name:   "a",
							Line:   1,
							CharAt: 1,
						},
						PropertyIdentifier: core.Identifier{
							Name:   "b",
							Line:   1,
							CharAt: 3,
						},
						Line:   1,
						CharAt: 1,
					},
					Right: &core.VariableExpression{
						Name:   "c",
						Line:   1,
						CharAt: 5,
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
