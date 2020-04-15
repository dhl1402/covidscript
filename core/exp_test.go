package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvaluate_LiteralExpression(t *testing.T) {
	cases := []struct {
		name string
		ec   *ExecutionContext
		exp  Expression
		want Expression
		err  error
	}{
		{
			name: "evaluate literal expression",
			ec:   &ExecutionContext{},
			exp: &LiteralExpression{
				Type:  "number",
				Value: "1",
			},
			want: &LiteralExpression{
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
		ec   *ExecutionContext
		exp  Expression
		want Expression
		err  error
	}{
		{
			name: "evaluate variable expression #1",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &LiteralExpression{
						Type:   "number",
						Value:  "3",
						Line:   1,
						CharAt: 1,
					},
				},
			},
			exp: &VariableExpression{
				Name:   "a",
				Line:   2,
				CharAt: 1,
			},
			want: &LiteralExpression{
				Type:   "number",
				Value:  "3",
				Line:   2,
				CharAt: 1,
			},
			err: nil,
		},
		{
			name: "evaluate variable expression #2",
			ec:   &ExecutionContext{},
			exp: &VariableExpression{
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

func TestEvaluate_BinaryExpression(t *testing.T) {
	cases := []struct {
		name string
		ec   *ExecutionContext
		exp  Expression
		want Expression
		err  error
	}{
		{
			name: "evaluate binary expression #1",
			ec: &ExecutionContext{
				Variables: map[string]Expression{},
			},
			exp: &BinaryExpression{
				Left: &LiteralExpression{
					Type:  "number",
					Value: "1",
				},
				Right: &LiteralExpression{
					Type:  "number",
					Value: "2",
				},
				Operator: Operator{
					Symbol: "+",
				},
			},
			want: &LiteralExpression{
				Type:  "number",
				Value: "3",
			},
			err: nil,
		},
		{
			name: "evaluate binary expression #2",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &LiteralExpression{
						Type:  "number",
						Value: "3",
					},
				},
			},
			exp: &BinaryExpression{
				Left: &LiteralExpression{
					Type:  "number",
					Value: "1",
				},
				Right: &BinaryExpression{
					Left: &LiteralExpression{
						Type:  "number",
						Value: "2",
					},
					Right: &VariableExpression{
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
			want: &LiteralExpression{
				Type:  "number",
				Value: "6",
			},
			err: nil,
		},
		{
			name: "evaluate binary expression #3",
			ec: &ExecutionContext{
				Variables: map[string]Expression{},
			},
			exp: &BinaryExpression{
				Left: &LiteralExpression{
					Type:  "number",
					Value: "1",
				},
				Right: &LiteralExpression{
					Type:  "string",
					Value: "2",
				},
				Operator: Operator{
					Symbol: "+",
				},
			},
			want: &LiteralExpression{
				Type:  "string",
				Value: "12",
			},
			err: nil,
		},
		{
			name: "evaluate binary expression #4",
			ec: &ExecutionContext{
				Variables: map[string]Expression{},
			},
			exp: &BinaryExpression{
				Left: &LiteralExpression{
					Type:  "string",
					Value: "abc",
				},
				Right: &LiteralExpression{
					Type:  "string",
					Value: "xyz",
				},
				Operator: Operator{
					Symbol: "+",
				},
			},
			want: &LiteralExpression{
				Type:  "string",
				Value: "abcxyz",
			},
			err: nil,
		},
		{
			name: "evaluate binary expression #5",
			ec: &ExecutionContext{
				Variables: map[string]Expression{},
			},
			exp: &BinaryExpression{
				Left: &LiteralExpression{
					Type:  "number",
					Value: "1",
				},
				Right: &LiteralExpression{
					Type:  "string",
					Value: "xyz",
				},
				Operator: Operator{
					Symbol: "-",
					Line:   1,
					CharAt: 2,
				},
			},
			want: nil,
			err:  fmt.Errorf("Cannot use '-' operator with string. [1,2]"),
		},
		{
			name: "evaluate binary expression #6",
			ec: &ExecutionContext{
				Variables: map[string]Expression{},
			},
			exp: &BinaryExpression{
				Left: &LiteralExpression{
					Type:  "number",
					Value: "1",
				},
				Right: &ArrayExpression{
					Elements: []Expression{
						&LiteralExpression{
							Type:  "number",
							Value: "1",
						},
					},
				},
				Operator: Operator{
					Symbol: "+",
					Line:   1,
					CharAt: 2,
				},
			},
			want: nil,
			err:  fmt.Errorf("Cannot use '+' operator with array. [1,2]"),
		},
		{
			name: "evaluate binary expression #6",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &LiteralExpression{
						Type:   "number",
						Value:  "0",
						Line:   1,
						CharAt: 1,
					},
				},
			},
			exp: &BinaryExpression{
				Left: &LiteralExpression{
					Type:  "number",
					Value: "1",
				},
				Right: &VariableExpression{
					Name:   "a",
					Line:   1,
					CharAt: 3,
				},
				Operator: Operator{
					Symbol: "/",
				},
			},
			want: nil,
			err:  fmt.Errorf("Cannot divide by zero. [1,3]"),
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
		ec   *ExecutionContext
		exp  Expression
		want Expression
		err  error
	}{
		{
			name: "evaluate array expression #1",
			ec: &ExecutionContext{
				Variables: map[string]Expression{},
			},
			exp: &ArrayExpression{
				Elements: []Expression{
					&LiteralExpression{
						Type:  "number",
						Value: "1",
					},
				},
			},
			want: &ArrayExpression{
				Elements: []Expression{
					&LiteralExpression{
						Type:  "number",
						Value: "1",
					},
				},
			},
			err: nil,
		},
		{
			name: "evaluate array expression #2",
			ec: &ExecutionContext{
				Variables: map[string]Expression{},
			},
			exp: &ArrayExpression{
				Elements: []Expression{
					&BinaryExpression{
						Left: &LiteralExpression{
							Type:  "number",
							Value: "1",
						},
						Right: &LiteralExpression{
							Type:  "number",
							Value: "1",
						},
						Operator: Operator{
							Symbol: "+",
						},
					},
				},
			},
			want: &ArrayExpression{
				Elements: []Expression{
					&LiteralExpression{
						Type:  "number",
						Value: "2",
					},
				},
			},
			err: nil,
		},
		{
			name: "evaluate array expression #3",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &LiteralExpression{
						Type:  "number",
						Value: "3",
					},
				},
			},
			exp: &ArrayExpression{
				Elements: []Expression{
					&BinaryExpression{
						Left: &LiteralExpression{
							Type:  "number",
							Value: "1",
						},
						Right: &BinaryExpression{
							Left: &LiteralExpression{
								Type:  "number",
								Value: "2",
							},
							Right: &VariableExpression{
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
					&VariableExpression{
						Name: "a",
					},
				},
			},
			want: &ArrayExpression{
				Elements: []Expression{
					&LiteralExpression{
						Type:  "number",
						Value: "6",
					},
					&LiteralExpression{
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
		ec   *ExecutionContext
		exp  Expression
		want Expression
		err  error
	}{
		{
			name: "evaluate object expression #1",
			ec: &ExecutionContext{
				Variables: map[string]Expression{},
			},
			exp: &ObjectExpression{
				Properties: []ObjectProperty{
					{
						KeyIdentifier: Identifier{
							Name: "a",
						},
						Value: &LiteralExpression{
							Type:  "string",
							Value: "xxx",
						},
					},
				},
			},
			want: &ObjectExpression{
				Properties: []ObjectProperty{
					{
						KeyIdentifier: Identifier{
							Name: "a",
						},
						Value: &LiteralExpression{
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
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &LiteralExpression{
						Type:  "number",
						Value: "3",
					},
				},
			},
			exp: &ObjectExpression{
				Properties: []ObjectProperty{
					{
						KeyExpression: &BinaryExpression{
							Left: &LiteralExpression{
								Type:  "string",
								Value: "a",
							},
							Right: &LiteralExpression{
								Type:  "string",
								Value: "b",
							},
							Operator: Operator{
								Symbol: "+",
							},
						},
						Value: &BinaryExpression{
							Left: &LiteralExpression{
								Type:  "number",
								Value: "2",
							},
							Right: &VariableExpression{
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
			want: &ObjectExpression{
				Properties: []ObjectProperty{
					{
						KeyExpression: &LiteralExpression{
							Type:  "string",
							Value: "ab",
						},
						Value: &LiteralExpression{
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

func TestEvaluate_MemberAccessExpression(t *testing.T) {
	cases := []struct {
		name string
		ec   *ExecutionContext
		exp  Expression
		want Expression
		err  error
	}{
		{
			name: "evaluate member access expression #1",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &ObjectExpression{
						Properties: []ObjectProperty{
							{
								KeyExpression: &LiteralExpression{
									Type:  "string",
									Value: "b",
								},
								Value: &LiteralExpression{
									Type:  "boolean",
									Value: "true",
								},
								Computed: true,
							},
						},
					},
				},
			},
			exp: &MemberAccessExpression{
				Object: &VariableExpression{
					Name: "a",
				},
				PropertyIdentifier: Identifier{
					Name: "b",
				},
			},
			want: &LiteralExpression{
				Type:  "boolean",
				Value: "true",
			},
			err: nil,
		},
		{
			name: "evaluate member access expression #2",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &ObjectExpression{
						Properties: []ObjectProperty{
							{
								KeyExpression: &LiteralExpression{
									Type:  "string",
									Value: "b",
								},
								Value: &LiteralExpression{
									Type:  "boolean",
									Value: "true",
								},
								Computed: true,
							},
						},
					},
				},
			},
			exp: &MemberAccessExpression{
				Object: &VariableExpression{
					Name: "a",
				},
				PropertyExpression: &LiteralExpression{
					Type:  "string",
					Value: "b",
				},
				Compute: true,
			},
			want: &LiteralExpression{
				Type:  "boolean",
				Value: "true",
			},
			err: nil,
		},
		{
			name: "evaluate member access expression #3",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &ObjectExpression{
						Properties: []ObjectProperty{
							{
								KeyIdentifier: Identifier{
									Name: "b",
								},
								Value: &LiteralExpression{
									Type:  "boolean",
									Value: "true",
								},
							},
						},
					},
				},
			},
			exp: &MemberAccessExpression{
				Object: &VariableExpression{
					Name: "a",
				},
				PropertyExpression: &LiteralExpression{
					Type:  "string",
					Value: "b",
				},
				Compute: true,
			},
			want: &LiteralExpression{
				Type:  "boolean",
				Value: "true",
			},
			err: nil,
		},
		{
			name: "evaluate member access expression #4",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &ObjectExpression{},
				},
			},
			exp: &MemberAccessExpression{
				Object: &VariableExpression{
					Name: "a",
				},
				PropertyIdentifier: Identifier{
					Name: "c",
				},
			},
			want: &LiteralExpression{
				Type: "undefined",
			},
			err: nil,
		},
		{
			name: "evaluate member access expression #5",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &ObjectExpression{},
				},
			},
			exp: &MemberAccessExpression{
				Object: &VariableExpression{
					Name: "a",
				},
				PropertyExpression: &LiteralExpression{
					Type:   "boolean",
					Value:  "true",
					Line:   1,
					CharAt: 4,
				},
				Compute: true,
			},
			want: nil,
			err:  fmt.Errorf("Property key of type boolean is not supported. [1,4]"),
		},
		{
			name: "evaluate member access expression #6",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &ArrayExpression{
						Elements: []Expression{
							&LiteralExpression{
								Type:  "number",
								Value: "1",
							},
						},
					},
				},
			},
			exp: &MemberAccessExpression{
				Object: &VariableExpression{
					Name: "a",
				},
				PropertyExpression: &LiteralExpression{
					Type:  "number",
					Value: "0",
				},
				Compute: true,
				Line:    1,
				CharAt:  1,
			},
			want: &LiteralExpression{
				Type:   "number",
				Value:  "1",
				Line:   1,
				CharAt: 1,
			},
			err: nil,
		},
		{
			name: "evaluate member access expression #7",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &ArrayExpression{
						Elements: []Expression{
							&LiteralExpression{
								Type:  "number",
								Value: "1",
							},
						},
					},
				},
			},
			exp: &MemberAccessExpression{
				Object: &VariableExpression{
					Name: "a",
				},
				PropertyIdentifier: Identifier{
					Name: "xxx",
				},
				Line:   1,
				CharAt: 1,
			},
			want: &LiteralExpression{
				Type:   "undefined",
				Line:   1,
				CharAt: 1,
			},
			err: nil,
		},
		{
			name: "evaluate member access expression #8",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &ArrayExpression{
						Elements: []Expression{
							&LiteralExpression{
								Type:  "number",
								Value: "1",
							},
						},
					},
				},
			},
			exp: &MemberAccessExpression{
				Object: &VariableExpression{
					Name: "a",
				},
				PropertyExpression: &LiteralExpression{
					Type:   "number",
					Value:  "1",
					Line:   1,
					CharAt: 3,
				},
				Compute: true,
				Line:    1,
				CharAt:  1,
			},
			want: nil,
			err:  fmt.Errorf("Index is out of range. [1.3]"),
		},
		{
			name: "evaluate member access expression #9",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &ArrayExpression{
						Elements: []Expression{
							&LiteralExpression{
								Type:  "number",
								Value: "1",
							},
						},
					},
				},
			},
			exp: &MemberAccessExpression{
				Object: &VariableExpression{
					Name: "a",
				},
				PropertyExpression: &LiteralExpression{
					Type:   "string",
					Value:  "1",
					Line:   1,
					CharAt: 3,
				},
				Compute: true,
				Line:    1,
				CharAt:  1,
			},
			want: nil,
			err:  fmt.Errorf("Index must be number. [1,3]"),
		},
		{
			name: "evaluate member access expression #10",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &ArrayExpression{
						Elements: []Expression{
							&LiteralExpression{
								Type:  "number",
								Value: "1",
							},
						},
					},
				},
			},
			exp: &MemberAccessExpression{
				Object: &LiteralExpression{
					Type:  "number",
					Value: "0",
				},
				PropertyExpression: &LiteralExpression{
					Type:   "string",
					Value:  "1",
					Line:   1,
					CharAt: 3,
				},
				Compute: true,
				Line:    1,
				CharAt:  1,
			},
			want: nil,
			err:  fmt.Errorf("Can't access property of type number. [1,1]"),
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

func TestEvaluate_FunctionExpression(t *testing.T) {
	cases := []struct {
		name string
		ec   *ExecutionContext
		exp  Expression
		want Expression
		err  error
	}{
		{
			name: "evaluate function expression #1",
			ec:   nil,
			exp: &FunctionExpression{
				Params: []Identifier{
					{
						Name: "a",
					},
					{
						Name: "b",
					},
				},
				Body: []Statement{},
				EC:   &ExecutionContext{},
			},
			want: &FunctionExpression{
				Params: []Identifier{
					{
						Name: "a",
					},
					{
						Name: "b",
					},
				},
				Body: []Statement{},
				EC: &ExecutionContext{
					Variables: map[string]Expression{
						"a": &LiteralExpression{
							Type: "undefined",
						},
						"b": &LiteralExpression{
							Type: "undefined",
						},
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

func TestEvaluate_CallExpression(t *testing.T) {
	cases := []struct {
		name string
		ec   *ExecutionContext
		exp  Expression
		want Expression
		err  error
	}{
		{
			name: "evaluate call expression #1",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &FunctionExpression{
						Params: []Identifier{
							{
								Name: "a",
							},
							{
								Name: "b",
							},
						},
						Body: []Statement{},
						EC:   &ExecutionContext{},
					},
				},
			},
			exp: &CallExpression{
				Callee: &VariableExpression{
					Name: "a",
				},
				Arguments: []Expression{},
			},
			want: &LiteralExpression{
				Type: "undefined",
			},
			err: nil,
		},
		{
			name: "evaluate call expression #2",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &FunctionExpression{
						Params: []Identifier{
							{
								Name: "a",
							},
							{
								Name: "b",
							},
						},
						Body: []Statement{
							ReturnStatement{
								Argument: &BinaryExpression{
									Left: &VariableExpression{
										Name: "a",
									},
									Right: &VariableExpression{
										Name: "b",
									},
									Operator: Operator{
										Symbol: "+",
									},
								},
							},
						},
						EC: &ExecutionContext{},
					},
				},
			},
			exp: &CallExpression{
				Callee: &VariableExpression{
					Name: "a",
				},
				Arguments: []Expression{
					&LiteralExpression{
						Type:  "number",
						Value: "1",
					},
					&LiteralExpression{
						Type:  "number",
						Value: "2",
					},
				},
			},
			want: &LiteralExpression{
				Type:  "number",
				Value: "3",
			},
			err: nil,
		},
		{
			name: "evaluate call expression #1",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &LiteralExpression{
						Type:  "number",
						Value: "1",
					},
				},
			},
			exp: &CallExpression{
				Callee: &VariableExpression{
					Name: "a",
				},
				Arguments: []Expression{},
				Line:      1,
				CharAt:    1,
			},
			want: nil,
			err:  fmt.Errorf("a is not a function. [1,1]"), // TODO: e.Callee.ToString()
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

func TestEvaluate_TMP(t *testing.T) {
	cases := []struct {
		name string
		ec   *ExecutionContext
		exp  Expression
		want Expression
		err  error
	}{
		{
			name: "evaluate member access expression #6",
			ec: &ExecutionContext{
				Variables: map[string]Expression{
					"a": &ArrayExpression{
						Elements: []Expression{
							&LiteralExpression{
								Type:  "number",
								Value: "1",
							},
						},
					},
				},
			},
			exp: &MemberAccessExpression{
				Object: &VariableExpression{
					Name: "a",
				},
				PropertyExpression: &LiteralExpression{
					Type:  "number",
					Value: "0",
				},
				Compute: true,
				Line:    1,
				CharAt:  1,
			},
			want: &LiteralExpression{
				Type:   "number",
				Value:  "1",
				Line:   1,
				CharAt: 1,
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
