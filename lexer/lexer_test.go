package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLex_Value(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{
			name: "lex variable declaration",
			in:   `var a=123   `,
			want: []string{"var", "a", "=", "123"},
		},
		{
			name: "lex variable declaration",
			in:   `var a,b=123,456`,
			want: []string{"var", "a", ",", "b", "=", "123", ",", "456"},
		},
		{
			name: "lex string variable declaration",
			in:   "var a=`ab+c`",
			want: []string{"var", "a", "=", "`ab+c`"},
		},
		{
			name: "lex boolean variable declaration",
			in:   `var a, b=false,true`,
			want: []string{"var", "a", ",", "b", "=", "false", ",", "true"},
		},
		{
			name: "lex array variable declaration",
			in:   `var a=[1,"2",3]`,
			want: []string{"var", "a", "=", "[", "1", ",", `"2"`, ",", "3", "]"},
		},
		{
			name: "lex object variable declaration",
			in:   `var a={b:1,c:"2"}`,
			want: []string{"var", "a", "=", "{", "b", ":", "1", ",", "c", ":", `"2"`, "}"},
		},
		{
			name: "lex variable declaration",
			in: `
			var a =123
			b:=a
			var a,b = 1,2
			a,b:=1,2
			var a
			var a,b`,
			want: []string{"var", "a", "=", "123", "b", ":", "=", "a", "var", "a", ",", "b", "=", "1", ",", "2", "a", ",", "b", ":", "=", "1", ",", "2", "var", "a", "var", "a", ",", "b"},
		},
		{
			name: "lex function declaration",
			in: `
			func a(b,c){
				var d=1
			}
			`,
			want: []string{"func", "a", "(", "b", ",", "c", ")", "{", "var", "d", "=", "1", "}"},
		},
		{
			name: "lex function declaration",
			in: `
			var a = func (b,c){
				var d=1
			}
			`,
			want: []string{"var", "a", "=", "func", "(", "b", ",", "c", ")", "{", "var", "d", "=", "1", "}"},
		},
		{
			name: "lex function declaration",
			in: `
			var a = {
				b: func (b,c) {
					var d=1
				}
			}
			`,
			want: []string{"var", "a", "=", "{", "b", ":", "func", "(", "b", ",", "c", ")", "{", "var", "d", "=", "1", "}", "}"},
		},
		{
			name: "lex function declaration",
			in: `setTimeout(func(a,b){
				console.log(a,b)
			}, 3000)`,
			want: []string{"setTimeout", "(", "func", "(", "a", ",", "b", ")", "{", "console", ".", "log", "(", "a", ",", "b", ")", "}", ",", "3000", ")"},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actual := []string{}
			tokens := Lex(tt.in)
			for _, tk := range tokens {
				actual = append(actual, tk.Value)
			}
			require.Equal(t, tt.want, actual)
		})
	}
}

func TestLex_Token(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []Token
	}{
		{
			name: "lex variable declaration",
			in:   `var abc=123   `,
			want: []Token{
				Token{Value: "var", Line: 1, CharAt: 1},
				Token{Value: "abc", Line: 1, CharAt: 5},
				Token{Value: "=", Line: 1, CharAt: 8},
				Token{Value: "123", Line: 1, CharAt: 9},
			},
		},
		{
			name: "lex string variable declaration",
			in:   `var a,b="123","456"   `,
			want: []Token{
				Token{Value: "var", Line: 1, CharAt: 1},
				Token{Value: "a", Line: 1, CharAt: 5},
				Token{Value: ",", Line: 1, CharAt: 6},
				Token{Value: "b", Line: 1, CharAt: 7},
				Token{Value: "=", Line: 1, CharAt: 8},
				Token{Value: `"123"`, Line: 1, CharAt: 9},
				Token{Value: ",", Line: 1, CharAt: 14},
				Token{Value: `"456"`, Line: 1, CharAt: 15},
			},
		},
		{
			name: "lex multiline string declaration",
			in: `var a = 456
				 var b = "
1
2
3
"
				var c = 789
			`,
			want: []Token{
				Token{Value: "var", Line: 1, CharAt: 1},
				Token{Value: "a", Line: 1, CharAt: 5},
				Token{Value: "=", Line: 1, CharAt: 7},
				Token{Value: "456", Line: 1, CharAt: 9},
				Token{Value: "var", Line: 2, CharAt: 1},
				Token{Value: "b", Line: 2, CharAt: 5},
				Token{Value: "=", Line: 2, CharAt: 7},
				Token{Value: `"
1
2
3
"`,
					Line: 2, CharAt: 9},
				Token{Value: "var", Line: 3, CharAt: 1},
				Token{Value: "c", Line: 3, CharAt: 5},
				Token{Value: "=", Line: 3, CharAt: 7},
				Token{Value: "789", Line: 3, CharAt: 9},
			},
		},
		{
			name: "lex variable declaration",
			in: `var abc = 123
				 var xyz = 456
			    `,
			want: []Token{
				Token{Value: "var", Line: 1, CharAt: 1},
				Token{Value: "abc", Line: 1, CharAt: 5},
				Token{Value: "=", Line: 1, CharAt: 9},
				Token{Value: "123", Line: 1, CharAt: 11},
				Token{Value: "var", Line: 2, CharAt: 1},
				Token{Value: "xyz", Line: 2, CharAt: 5},
				Token{Value: "=", Line: 2, CharAt: 9},
				Token{Value: "456", Line: 2, CharAt: 11},
			},
		},
		{
			name: "lex function declaration",
			in: `setTimeout(func(a,b){
				console.log(a,b)
			}, 3000)`,
			want: []Token{
				Token{Value: "setTimeout", Line: 1, CharAt: 1},
				Token{Value: "(", Line: 1, CharAt: 11},
				Token{Value: "func", Line: 1, CharAt: 12},
				Token{Value: "(", Line: 1, CharAt: 16},
				Token{Value: "a", Line: 1, CharAt: 17},
				Token{Value: ",", Line: 1, CharAt: 18},
				Token{Value: "b", Line: 1, CharAt: 19},
				Token{Value: ")", Line: 1, CharAt: 20},
				Token{Value: "{", Line: 1, CharAt: 21},
				Token{Value: "console", Line: 2, CharAt: 1},
				Token{Value: ".", Line: 2, CharAt: 8},
				Token{Value: "log", Line: 2, CharAt: 9},
				Token{Value: "(", Line: 2, CharAt: 12},
				Token{Value: "a", Line: 2, CharAt: 13},
				Token{Value: ",", Line: 2, CharAt: 14},
				Token{Value: "b", Line: 2, CharAt: 15},
				Token{Value: ")", Line: 2, CharAt: 16},
				Token{Value: "}", Line: 3, CharAt: 1},
				Token{Value: ",", Line: 3, CharAt: 2},
				Token{Value: "3000", Line: 3, CharAt: 4},
				Token{Value: ")", Line: 3, CharAt: 8},
			},
		},
		{
			name: "lex object declaration",
			in: `var a={
					b:'123',
					c:456
				}`,
			want: []Token{
				Token{Value: "var", Line: 1, CharAt: 1},
				Token{Value: "a", Line: 1, CharAt: 5},
				Token{Value: "=", Line: 1, CharAt: 6},
				Token{Value: "{", Line: 1, CharAt: 7},
				Token{Value: "b", Line: 2, CharAt: 1},
				Token{Value: ":", Line: 2, CharAt: 2},
				Token{Value: `'123'`, Line: 2, CharAt: 3},
				Token{Value: ",", Line: 2, CharAt: 8},
				Token{Value: "c", Line: 3, CharAt: 1},
				Token{Value: ":", Line: 3, CharAt: 2},
				Token{Value: "456", Line: 3, CharAt: 3},
				Token{Value: "}", Line: 4, CharAt: 1},
			},
		},
		{
			name: "lex multiline expression",
			in: `var a=1 
			+
			2`,
			want: []Token{
				Token{Value: "var", Line: 1, CharAt: 1},
				Token{Value: "a", Line: 1, CharAt: 5},
				Token{Value: "=", Line: 1, CharAt: 6},
				Token{Value: "1", Line: 1, CharAt: 7},
				Token{Value: "+", Line: 2, CharAt: 1},
				Token{Value: "2", Line: 3, CharAt: 1},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, Lex(tt.in))
		})
	}
}
