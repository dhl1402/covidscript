package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvaluate_LiteralExpression(t *testing.T) {
	cases := []struct {
		name string
		ec   ExecutionContext
		exp  Expression
		want Expression
		err  error
	}{
		{
			name: "evaluate literal expression",
			ec:   ExecutionContext{},
			exp: LiteralExpression{
				Type:  "number",
				Value: "1",
			},
			want: LiteralExpression{
				Type:  "number",
				Value: "1",
			},
			err: nil,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp, err := tt.exp.Evaluate(tt.ec)
			require.Equal(t, tt.want, exp)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestEvaluate_VariableExpression(t *testing.T) {
	cases := []struct {
		name string
		ec   ExecutionContext
		exp  Expression
		want Expression
		err  error
	}{
		{
			name: "evaluate array expression #1",
			ec: ExecutionContext{
				Variables: map[string]Expression{
					"a": LiteralExpression{
						Type:  "number",
						Value: "3",
					},
				},
			},
			exp: VariableExpression{
				Name:   "a",
				Line:   1,
				CharAt: 1,
			},
			want: LiteralExpression{
				Type:  "number",
				Value: "3",
			},
			err: nil,
		},
		{
			name: "evaluate variable expression #2",
			ec:   ExecutionContext{},
			exp: VariableExpression{
				Name:   "a",
				Line:   1,
				CharAt: 1,
			},
			want: nil,
			err:  fmt.Errorf("a is not defined. [1,1]"),
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp, err := tt.exp.Evaluate(tt.ec)
			require.Equal(t, tt.want, exp)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestEvaluate_ArrayExpression(t *testing.T) {
	cases := []struct {
		name string
		ec   ExecutionContext
		exp  Expression
		want Expression
		err  error
	}{
		{
			name: "evaluate array expression #1",
			ec: ExecutionContext{
				Variables: map[string]Expression{},
			},
			exp: ArrayExpression{
				Elements: []Expression{
					LiteralExpression{
						Type:  "number",
						Value: "1",
					},
				},
			},
			want: ArrayExpression{
				Elements: []Expression{
					LiteralExpression{
						Type:  "number",
						Value: "1",
					},
				},
			},
			err: nil,
		},
		{
			name: "evaluate array expression #2",
			ec: ExecutionContext{
				Variables: map[string]Expression{},
			},
			exp: ArrayExpression{
				Elements: []Expression{
					BinaryExpression{
						Left: LiteralExpression{
							Type:  "number",
							Value: "1",
						},
						Right: LiteralExpression{
							Type:  "number",
							Value: "1",
						},
						Operator: Operator{
							Symbol: "+",
						},
					},
				},
			},
			want: ArrayExpression{
				Elements: []Expression{
					LiteralExpression{
						Type:  "number",
						Value: "2",
					},
				},
			},
			err: nil,
		},
		{
			name: "evaluate array expression #3",
			ec: ExecutionContext{
				Variables: map[string]Expression{
					"a": LiteralExpression{
						Type:  "number",
						Value: "3",
					},
				},
			},
			exp: ArrayExpression{
				Elements: []Expression{
					BinaryExpression{
						Left: LiteralExpression{
							Type:  "number",
							Value: "1",
						},
						Right: BinaryExpression{
							Left: LiteralExpression{
								Type:  "number",
								Value: "2",
							},
							Right: VariableExpression{
								Name: "a",
							},
							Operator: Operator{
								Symbol: "+",
							},
						},
						Operator: Operator{
							Symbol: "+",
						},
					},
					VariableExpression{
						Name: "a",
					},
				},
			},
			want: ArrayExpression{
				Elements: []Expression{
					LiteralExpression{
						Type:  "number",
						Value: "6",
					},
					LiteralExpression{
						Type:  "number",
						Value: "3",
					},
				},
			},
			err: nil,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp, err := tt.exp.Evaluate(tt.ec)
			require.Equal(t, tt.want, exp)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestEvaluate_ObjectExpression(t *testing.T) {
	cases := []struct {
		name string
		ec   ExecutionContext
		exp  Expression
		want Expression
		err  error
	}{
		{
			name: "evaluate object expression #1",
			ec: ExecutionContext{
				Variables: map[string]Expression{},
			},
			exp: ObjectExpression{
				Properties: []ObjectProperty{
					ObjectProperty{
						KeyIdentifier: Identifier{
							Name: "a",
						},
						Value: LiteralExpression{
							Type:  "string",
							Value: "xxx",
						},
					},
				},
			},
			want: ObjectExpression{
				Properties: []ObjectProperty{
					ObjectProperty{
						KeyIdentifier: Identifier{
							Name: "a",
						},
						Value: LiteralExpression{
							Type:  "string",
							Value: "xxx",
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "evaluate object expression #2",
			ec: ExecutionContext{
				Variables: map[string]Expression{
					"a": LiteralExpression{
						Type:  "number",
						Value: "3",
					},
				},
			},
			exp: ObjectExpression{
				Properties: []ObjectProperty{
					ObjectProperty{
						KeyExpression: BinaryExpression{
							Left: LiteralExpression{
								Type:  "string",
								Value: "a",
							},
							Right: LiteralExpression{
								Type:  "string",
								Value: "b",
							},
							Operator: Operator{
								Symbol: "+",
							},
						},
						Value: BinaryExpression{
							Left: LiteralExpression{
								Type:  "number",
								Value: "2",
							},
							Right: VariableExpression{
								Name: "a",
							},
							Operator: Operator{
								Symbol: "+",
							},
						},
						Computed: true,
					},
				},
			},
			want: ObjectExpression{
				Properties: []ObjectProperty{
					ObjectProperty{
						KeyExpression: LiteralExpression{
							Type:  "string",
							Value: "ab",
						},
						Value: LiteralExpression{
							Type:  "number",
							Value: "5",
						},
						Computed: true,
					},
				},
			},
			err: nil,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp, err := tt.exp.Evaluate(tt.ec)
			require.Equal(t, tt.want, exp)
			require.Equal(t, tt.err, err)
		})
	}
}
