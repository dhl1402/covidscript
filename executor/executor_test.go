package executor

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gs/lexer"
	"gs/parser"
)

func Test(t *testing.T) {
	cases := []struct {
		name string
		in   string
		// want Expression
		want error
	}{
		{
			name: "",
			in: `var a=[1,2]
				var b = a.ac
			`,
			want: nil,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ast, _ := parser.ToAST(lexer.Lex(tt.in))
			require.Equal(t, tt.want, Execute(ast))
		})
	}
}
