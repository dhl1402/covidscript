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
			in:   `var a, b=#f,#t`,
			want: []string{"var", "a", ",", "b", "=", "#f", ",", "#t"},
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
			want: []string{"var", "a", "=", "123", "b", ":=", "a", "var", "a", ",", "b", "=", "1", ",", "2", "a", ",", "b", ":=", "1", ",", "2", "var", "a", "var", "a", ",", "b"},
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
			tokens, err := Lex(tt.in)
			require.Equal(t, err, nil)
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
				{Value: "var", Line: 1, CharAt: 1},
				{Value: "abc", Line: 1, CharAt: 5},
				{Value: "=", Line: 1, CharAt: 8},
				{Value: "123", Line: 1, CharAt: 9},
			},
		},
		{
			name: "lex string variable declaration",
			in:   `var a,b="123","456"   `,
			want: []Token{
				{Value: "var", Line: 1, CharAt: 1},
				{Value: "a", Line: 1, CharAt: 5},
				{Value: ",", Line: 1, CharAt: 6},
				{Value: "b", Line: 1, CharAt: 7},
				{Value: "=", Line: 1, CharAt: 8},
				{Value: `"123"`, Line: 1, CharAt: 9},
				{Value: ",", Line: 1, CharAt: 14},
				{Value: `"456"`, Line: 1, CharAt: 15},
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
				{Value: "var", Line: 1, CharAt: 1},
				{Value: "a", Line: 1, CharAt: 5},
				{Value: "=", Line: 1, CharAt: 7},
				{Value: "456", Line: 1, CharAt: 9},
				{Value: "var", Line: 2, CharAt: 1},
				{Value: "b", Line: 2, CharAt: 5},
				{Value: "=", Line: 2, CharAt: 7},
				{Value: `"
1
2
3
"`,
					Line: 2, CharAt: 9},
				{Value: "var", Line: 3, CharAt: 1},
				{Value: "c", Line: 3, CharAt: 5},
				{Value: "=", Line: 3, CharAt: 7},
				{Value: "789", Line: 3, CharAt: 9},
			},
		},
		{
			name: "lex variable declaration",
			in: `var abc = 123
				 var xyz = 456
			    `,
			want: []Token{
				{Value: "var", Line: 1, CharAt: 1},
				{Value: "abc", Line: 1, CharAt: 5},
				{Value: "=", Line: 1, CharAt: 9},
				{Value: "123", Line: 1, CharAt: 11},
				{Value: "var", Line: 2, CharAt: 1},
				{Value: "xyz", Line: 2, CharAt: 5},
				{Value: "=", Line: 2, CharAt: 9},
				{Value: "456", Line: 2, CharAt: 11},
			},
		},
		{
			name: "lex function declaration",
			in: `setTimeout(func(a,b){
				console.log(a,b)
			}, 3000)`,
			want: []Token{
				{Value: "setTimeout", Line: 1, CharAt: 1},
				{Value: "(", Line: 1, CharAt: 11},
				{Value: "func", Line: 1, CharAt: 12},
				{Value: "(", Line: 1, CharAt: 16},
				{Value: "a", Line: 1, CharAt: 17},
				{Value: ",", Line: 1, CharAt: 18},
				{Value: "b", Line: 1, CharAt: 19},
				{Value: ")", Line: 1, CharAt: 20},
				{Value: "{", Line: 1, CharAt: 21},
				{Value: "console", Line: 2, CharAt: 1},
				{Value: ".", Line: 2, CharAt: 8},
				{Value: "log", Line: 2, CharAt: 9},
				{Value: "(", Line: 2, CharAt: 12},
				{Value: "a", Line: 2, CharAt: 13},
				{Value: ",", Line: 2, CharAt: 14},
				{Value: "b", Line: 2, CharAt: 15},
				{Value: ")", Line: 2, CharAt: 16},
				{Value: "}", Line: 3, CharAt: 1},
				{Value: ",", Line: 3, CharAt: 2},
				{Value: "3000", Line: 3, CharAt: 4},
				{Value: ")", Line: 3, CharAt: 8},
			},
		},
		{
			name: "lex object declaration",
			in: `var a={
					b:'123',
					c:456
				}`,
			want: []Token{
				{Value: "var", Line: 1, CharAt: 1},
				{Value: "a", Line: 1, CharAt: 5},
				{Value: "=", Line: 1, CharAt: 6},
				{Value: "{", Line: 1, CharAt: 7},
				{Value: "b", Line: 2, CharAt: 1},
				{Value: ":", Line: 2, CharAt: 2},
				{Value: `'123'`, Line: 2, CharAt: 3},
				{Value: ",", Line: 2, CharAt: 8},
				{Value: "c", Line: 3, CharAt: 1},
				{Value: ":", Line: 3, CharAt: 2},
				{Value: "456", Line: 3, CharAt: 3},
				{Value: "}", Line: 4, CharAt: 1},
			},
		},
		{
			name: "lex multiline expression",
			in: `var a=1 
			+
			2`,
			want: []Token{
				{Value: "var", Line: 1, CharAt: 1},
				{Value: "a", Line: 1, CharAt: 5},
				{Value: "=", Line: 1, CharAt: 6},
				{Value: "1", Line: 1, CharAt: 7},
				{Value: "+", Line: 2, CharAt: 1},
				{Value: "2", Line: 3, CharAt: 1},
			},
		},
		{
			name: "lex variable declaration with :=",
			in:   `a:=b`,
			want: []Token{
				{Value: "a", Line: 1, CharAt: 1},
				{Value: ":=", Line: 1, CharAt: 2},
				{Value: "b", Line: 1, CharAt: 4},
			},
		},
		{
			name: "lex multi char operator ||",
			in:   `a||b`,
			want: []Token{
				{Value: "a", Line: 1, CharAt: 1},
				{Value: "||", Line: 1, CharAt: 2},
				{Value: "b", Line: 1, CharAt: 4},
			},
		},
		{
			name: "lex multi char operator &&",
			in:   `a&&b`,
			want: []Token{
				{Value: "a", Line: 1, CharAt: 1},
				{Value: "&&", Line: 1, CharAt: 2},
				{Value: "b", Line: 1, CharAt: 4},
			},
		},
		{
			name: "lex multi char operator >=",
			in:   `a>=b`,
			want: []Token{
				{Value: "a", Line: 1, CharAt: 1},
				{Value: ">=", Line: 1, CharAt: 2},
				{Value: "b", Line: 1, CharAt: 4},
			},
		},
		{
			name: "lex multi char operator <=",
			in:   `a<=b`,
			want: []Token{
				{Value: "a", Line: 1, CharAt: 1},
				{Value: "<=", Line: 1, CharAt: 2},
				{Value: "b", Line: 1, CharAt: 4},
			},
		},
		{
			name: "lex multi char operator ==",
			in:   `a==b`,
			want: []Token{
				{Value: "a", Line: 1, CharAt: 1},
				{Value: "==", Line: 1, CharAt: 2},
				{Value: "b", Line: 1, CharAt: 4},
			},
		},
		{
			name: "lex multi char operator !=",
			in:   `a!=b`,
			want: []Token{
				{Value: "a", Line: 1, CharAt: 1},
				{Value: "!=", Line: 1, CharAt: 2},
				{Value: "b", Line: 1, CharAt: 4},
			},
		},
		{
			name: "lex multi char operator ===",
			in:   `a===b`,
			want: []Token{
				{Value: "a", Line: 1, CharAt: 1},
				{Value: "===", Line: 1, CharAt: 2},
				{Value: "b", Line: 1, CharAt: 5},
			},
		},
		{
			name: "lex multi char operator !==",
			in:   `a!==b`,
			want: []Token{
				{Value: "a", Line: 1, CharAt: 1},
				{Value: "!==", Line: 1, CharAt: 2},
				{Value: "b", Line: 1, CharAt: 5},
			},
		},
		{
			name: "lex float number #1",
			in:   `123.123`,
			want: []Token{
				{Value: "123.123", Line: 1, CharAt: 1},
			},
		},
		{
			name: "lex float number #2",
			in:   `1.1==1.2`,
			want: []Token{
				{Value: "1.1", Line: 1, CharAt: 1},
				{Value: "==", Line: 1, CharAt: 4},
				{Value: "1.2", Line: 1, CharAt: 6},
			},
		},
		{
			name: "lex float number #3",
			in:   `1.1==1.2.3`,
			want: []Token{
				{Value: "1.1", Line: 1, CharAt: 1},
				{Value: "==", Line: 1, CharAt: 4},
				{Value: "1.2", Line: 1, CharAt: 6},
				{Value: ".", Line: 1, CharAt: 9},
				{Value: "3", Line: 1, CharAt: 10},
			},
		},
		{
			name: "lex unary token #1",
			in:   `!a`,
			want: []Token{
				{Value: "!", Line: 1, CharAt: 1},
				{Value: "a", Line: 1, CharAt: 2},
			},
		},
		{
			name: "lex unary token #2",
			in:   `a !b`,
			want: []Token{
				{Value: "a", Line: 1, CharAt: 1},
				{Value: "!", Line: 1, CharAt: 3},
				{Value: "b", Line: 1, CharAt: 4},
			},
		},
		{
			name: "lex unary token #3",
			in:   `!a=!b`,
			want: []Token{
				{Value: "!", Line: 1, CharAt: 1},
				{Value: "a", Line: 1, CharAt: 2},
				{Value: "=", Line: 1, CharAt: 3},
				{Value: "!", Line: 1, CharAt: 4},
				{Value: "b", Line: 1, CharAt: 5},
			},
		},
		{
			name: "lex unary token #3",
			in:   `a=2!`,
			want: []Token{
				{Value: "a", Line: 1, CharAt: 1},
				{Value: "=", Line: 1, CharAt: 2},
				{Value: "2", Line: 1, CharAt: 3},
				{Value: "!", Line: 1, CharAt: 4},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := Lex(tt.in)
			require.Equal(t, err, nil)
			require.Equal(t, tt.want, tokens)
		})
	}
}

func TestLex_TMP(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []Token
	}{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := Lex(tt.in)
			require.Equal(t, err, nil)
			require.Equal(t, tt.want, tokens)
		})
	}
}
