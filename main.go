package main

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
	"os"
	"strconv"
	"reflect"
	"math"
	"github.com/fatih/color"
)


// -----------------------------
// ----------- Lexer -----------
// -----------------------------

type Token int
const (
	TOKEN_EOF = iota
	TOKEN_ILLEGAL
	TOKEN_ID
	TOKEN_STRING
	TOKEN_INT
	TOKEN_TYPE
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_END
	TOKEN_DO
	TOKEN_BOOL
	TOKEN_ELIF
	TOKEN_ELSE
	TOKEN_DIV
	TOKEN_MUL
	TOKEN_EQUALS
	TOKEN_IS_EQUALS
	TOKEN_NOT_EQUALS
	TOKEN_LESS_THAN
	TOKEN_GREATER_THAN
	TOKEN_LESS_EQUALS
	TOKEN_GREATER_EQUALS
	TOKEN_REM
	TOKEN_L_BRACKET
	TOKEN_R_BRACKET
	TOKEN_DOT
	TOKEN_COMMA
)

var tokens = []string{
	TOKEN_EOF:            "TOKEN_EOF",
	TOKEN_ILLEGAL:        "TOKEN_ILLEGAL",
	TOKEN_ID:             "TOKEN_ID",
	TOKEN_STRING:         "TOKEN_STRING",
	TOKEN_INT:            "TOKEN_INT",
	TOKEN_PLUS:           "TOKEN_PLUS",
	TOKEN_MINUS:          "TOKEN_MINUS",
	TOKEN_END:            "TOKEN_END",
	TOKEN_DO:             "TOKEN_DO",
	TOKEN_BOOL:           "TOKEN_BOOL",
	TOKEN_ELIF:           "TOKEN_ELIF",
	TOKEN_ELSE:           "TOKEN_ELSE",
	TOKEN_DIV:            "TOKEN_DIV",
	TOKEN_MUL:            "TOKEN_MUL",
	TOKEN_EQUALS:         "TOKEN_EQUALS",
	TOKEN_IS_EQUALS:      "TOKEN_IS_EQUALS",
	TOKEN_NOT_EQUALS:     "TOKEN_NOT_EQUALS",
	TOKEN_LESS_THAN:      "TOKEN_LESS_THAN",
	TOKEN_GREATER_THAN:   "TOKEN_GREATER_THAN",
	TOKEN_LESS_EQUALS:    "TOKEN_LESS_EQUALS",
	TOKEN_GREATER_EQUALS: "TOKEN_GREATER_EQUALS",
	TOKEN_REM:            "TOKEN_REM",
	TOKEN_L_BRACKET:      "TOKEN_L_BRACKET",
	TOKEN_R_BRACKET:      "TOKEN_R_BRACKET",
	TOKEN_DOT:            "TOKEN_DOT",
	TOKEN_COMMA:          "TOKEN_COMMA",
}

func (token Token) String() string {
	return tokens[token]
}

type Position struct {
	line int
	column int
}

type Lexer struct {
	pos Position
	reader *bufio.Reader
}

func LexerInit(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position {line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

func (lexer *Lexer) Lex() (Position, Token, string) {
	for {
		r, _, err := lexer.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				err = nil
				return lexer.pos, TOKEN_EOF, "EOF"
			}
			panic(err)
		}
		lexer.pos.column++
		switch r {
			case '\n': lexer.resetPosition()
			case '+': return lexer.pos, TOKEN_PLUS, "+"
			case '/': return lexer.pos, TOKEN_DIV, "/"
			case '*': return lexer.pos, TOKEN_MUL, "*"
			case '%': return lexer.pos, TOKEN_REM, "%"
			case '[': return lexer.pos, TOKEN_L_BRACKET, "["
			case ']': return lexer.pos, TOKEN_R_BRACKET, "]"
			case ',': return lexer.pos, TOKEN_COMMA, ","
			case '.': return lexer.pos, TOKEN_DOT, "."
			default:
				if unicode.IsSpace(r) {
					continue
				} else if r == '=' {
					r, _, err := lexer.reader.ReadRune()
					if err != nil {
						panic(err)
					}
					lexer.pos.column++
					if r == '=' {
						return lexer.pos, TOKEN_IS_EQUALS, "=="
					}
				} else if r == '-' {
					r, _, err := lexer.reader.ReadRune()
					if err != nil {
						if err == io.EOF {
							return lexer.pos, TOKEN_MINUS, "-"
						}
						panic(err)
					}
					lexer.pos.column++
					if r == '>' {
						return lexer.pos, TOKEN_EQUALS, "->"
					} else {
						return lexer.pos, TOKEN_MINUS, "-"
					}
				} else if r == '<' {
					r, _, err := lexer.reader.ReadRune()
					if err != nil {
						if err == io.EOF {
							return lexer.pos, TOKEN_LESS_THAN, "<"
						}
						panic(err)
					}
					if r == '=' {
						lexer.pos.column++
						return lexer.pos, TOKEN_LESS_EQUALS, "<="
					} else {
						return lexer.pos, TOKEN_LESS_THAN, "<"
					}
				} else if r == '>' {
					r, _, err := lexer.reader.ReadRune()
					if err != nil {
						if err == io.EOF {
							return lexer.pos, TOKEN_GREATER_THAN, ">"
						}
						panic(err)
					}
					if r == '=' {
						lexer.pos.column++
						return lexer.pos, TOKEN_GREATER_EQUALS, ">="
					} else {
						return lexer.pos, TOKEN_GREATER_THAN, ">"
					}
				} else if r == '!' {
					r, _, err := lexer.reader.ReadRune()
					if err != nil {panic(err)}
					lexer.pos.column++
					if r == '=' {
						return lexer.pos, TOKEN_NOT_EQUALS, "!="
					}
				} else if r == '#' {
					for {
						r, _, err := lexer.reader.ReadRune()
						if r == '\n' {break}
						if err != nil {panic(err)}
						lexer.pos.column++
					}
					continue
				} else if unicode.IsDigit(r) {
					startPos := lexer.pos
					lexer.backup()
					val := lexer.lexInt()
					return startPos, TOKEN_INT, val
				} else if unicode.IsLetter(r) {
					startPos := lexer.pos
					lexer.backup()
					val := lexer.lexId()
					if val == "end" {
						return startPos, TOKEN_END, val
					} else if val == "do" {
						return startPos, TOKEN_DO, val
					} else if val == "true" || val == "false" {
						return startPos, TOKEN_BOOL, val
					} else if val == "string" || val == "int" || val == "bool" || val == "type" /*|| val == "list"*/ {
						return startPos, TOKEN_TYPE, val
					} else if val == "else" {
						return startPos, TOKEN_ELSE, val
					} else if val == "elif" {
						return startPos, TOKEN_ELIF, val
					}
					return startPos, TOKEN_ID, val
				} else if r == '"' {
					startPos := lexer.pos
					lexer.backup()
					val := lexer.lexString()
					r, _, err = lexer.reader.ReadRune()
					return startPos, TOKEN_STRING, val
				}
        }
	}
}

func (lexer *Lexer) backup() {
	if err := lexer.reader.UnreadRune(); err != nil {
		panic(err)
	}
	lexer.pos.column--
}

func (lexer *Lexer) lexId() string {
	var val string
	for {
		r, _, err := lexer.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return val
			}
		}
        lexer.pos.column++
		if unicode.IsLetter(r) {
			val = val + string(r)
		} else {
			lexer.backup()
			return val
		}
	}
}

func (lexer *Lexer) lexInt() string {
	var val string
	for {
		r, _, err := lexer.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return val
			}
		}
		lexer.pos.column++
		if unicode.IsDigit(r) {
			val = val + string(r)
		} else {
			lexer.backup()
			return val
		}
	}
}

func (lexer *Lexer) lexString() string {
	var val string
	r, _, err := lexer.reader.ReadRune()
	for {
		r, _, err = lexer.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return val
			}
		}
		lexer.pos.column++
		if r != '"' {
			val = val + string(r)
		} else {
			lexer.backup()
			return val
		}
	}
}

func (lexer *Lexer) resetPosition() {
	lexer.pos.line++
	lexer.pos.column = 0
}


// -----------------------------
// ------------ AST ------------
// -----------------------------

type ExprType int
const (
	ExprVoid ExprType = iota
	ExprInt
	ExprStr
	ExprId
	ExprArr
	ExprAppend
	ExprReplace
	ExprRead
	ExprTypeType
	ExprPush
	ExprBlockdef
	ExprPrint
	ExprPrintS
	ExprPrintV
	ExprPuts
	ExprInput
	ExprOver
	ExprRot
	ExprInc
	ExprDec
	ExprLen
	ExprTypeOf
	ExprBreak
	ExprSwap
	ExprImport
	ExprCall
	ExprBool
	ExprIf
	ExprDup
	ExprDrop
	ExprExit
	ExprFor	
	ExprBinop // + - * / %
	ExprCompare // < > == !=
	ExprVardef
)

type Expr struct {
	Type ExprType
	AsInt float64
	AsStr string
	AsId *Id
	AsArr []Expr
	AsType string
	AsPush *Push
	AsBlockdef *Blockdef
	AsCall *Call
	AsBool bool
	AsIf *If
	AsFor *For
	AsBiniop int
	AsCompare int
	AsImport string
	AsVardef *Vardef
}

type Push struct {
	Arg Expr
}

type Call struct {
	Value string
}

type Blockdef struct {
	Name string
	Body []Expr
}

type If struct {
	Op []Expr
	Body []Expr
	ElifBodys [][]Expr
	ElifOps [][]Expr
	ElseBody []Expr
}

type For struct {
	Op []Expr
	Body []Expr
}

type Vardef struct {
	Name string
}

type Id struct {
	Name string
}


// -----------------------------
// ----------- Parse -----------
// -----------------------------

type Parser struct {
	current_token_type Token
	current_token_value string
	lexer Lexer
	line int
	column int
}

func ParserInit(lexer *Lexer) *Parser {
	pos, tok, val := lexer.Lex()
	return &Parser{
		current_token_type: tok,
		current_token_value: val,
		lexer: *lexer,
		line: pos.line,
		column: pos.column,
	}
}

func (parser *Parser) ParserEat(token Token) {
	if token != parser.current_token_type {
		fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value '%s'", parser.line, parser.column, parser.current_token_value))
		os.Exit(0)
	}
	pos, tok, val := parser.lexer.Lex()
	parser.current_token_type = tok
	parser.current_token_value = val
	parser.line = pos.line
	parser.column = pos.column
}

func StrToInt(num string) float64 {
	i, err := strconv.ParseFloat(num, 64)
	if err != nil{
		panic(err)
	}
	return i
}

func ParserParseExpr(parser *Parser) (Expr) {
	expr := Expr{}
	switch parser.current_token_type {
		case TOKEN_INT:
			expr.Type = ExprInt
			expr.AsInt = StrToInt(parser.current_token_value)
			parser.ParserEat(TOKEN_INT)
		case TOKEN_STRING:
			expr.Type = ExprStr
			expr.AsStr = parser.current_token_value
			parser.ParserEat(TOKEN_STRING)
		case TOKEN_BOOL:
			expr.Type = ExprBool
			if parser.current_token_value == "true" {
				expr.AsBool = true
			} else {
				expr.AsBool = false
			}
			parser.ParserEat(TOKEN_BOOL)
		case TOKEN_TYPE:
			expr.Type = ExprTypeType
			expr.AsType = parser.current_token_value
			parser.ParserEat(TOKEN_TYPE)
		case TOKEN_ID:
			expr.Type = ExprId
			vname := parser.current_token_value
			parser.ParserEat(TOKEN_ID)
			expr.AsId = &Id {
				Name: vname,
			}
		case TOKEN_L_BRACKET:
			parser.ParserEat(TOKEN_L_BRACKET)
			expr.Type = ExprArr
			parser.ParserEat(TOKEN_R_BRACKET)
		default:
			fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value '%s'", parser.line, parser.column, parser.current_token_value))
			os.Exit(0)
	}
	return expr
}

func ParserParse(parser *Parser)  ([]Expr) {
	exprs := []Expr{}

	for {
		expr := Expr{}
		if parser.current_token_type == TOKEN_ID {
			if parser.current_token_value == "print" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprPrint
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "printS" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprPrintS
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "printV" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprPrintV
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "input" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprInput
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "len" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprLen
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "puts" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprPuts
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "typeof" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprTypeOf
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "swap" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprSwap
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "over" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprOver
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "rot" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprRot
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "inc" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprInc
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "dec" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprDec
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "import" {
				parser.ParserEat(TOKEN_ID)
				if parser.current_token_type != TOKEN_STRING {
					fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value '%s'", parser.line, parser.column, parser.current_token_value))
					os.Exit(0)
				}
				expr.Type = ExprImport
				expr.AsImport = parser.current_token_value
				parser.ParserEat(TOKEN_STRING)
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "dup" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprDup
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "drop" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprDrop
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "exit" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprExit
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "block" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprBlockdef
				if parser.current_token_type != TOKEN_ID {
					fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value '%s'", parser.line, parser.column, parser.current_token_value))
					os.Exit(0)
				}
				name := parser.current_token_value
				parser.ParserEat(TOKEN_ID)
				parser.ParserEat(TOKEN_DO)
				if parser.current_token_type == TOKEN_END {
					fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: block '%s' body is empty", parser.line, parser.column, name))
					os.Exit(0)
				}
				body := ParserParse(parser)
				expr.AsBlockdef = &Blockdef{
					Name: name,
					Body: body,
				}
				parser.ParserEat(TOKEN_END)
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "for" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprFor
				op := ParserParse(parser)
				parser.ParserEat(TOKEN_DO)
				if parser.current_token_type == TOKEN_END {
					fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: for loop body is empty", parser.line, parser.column))
					os.Exit(0)
				}
				body := ParserParse(parser)
				parser.ParserEat(TOKEN_END)
				expr.AsFor = &For{
					Op: op,
					Body: body,
				}
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "if" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprIf
				op := ParserParse(parser)
				parser.ParserEat(TOKEN_DO)
				if parser.current_token_type == TOKEN_ELSE || parser.current_token_type == TOKEN_END {
					fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: if statement body is empty", parser.line, parser.column))
					os.Exit(0)
				}
				body := ParserParse(parser)
				ElifBodys := [][]Expr{}
				ElifOps := [][]Expr{}
				if parser.current_token_type == TOKEN_ELIF {
					for {
						parser.ParserEat(TOKEN_ELIF)
						if parser.current_token_type == TOKEN_ELSE || parser.current_token_type == TOKEN_END || parser.current_token_type == TOKEN_DO {
							fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: if statement invalid syntax", parser.line, parser.column))
							os.Exit(0)
						}
						ElifOp := ParserParse(parser)
						parser.ParserEat(TOKEN_DO)
						if parser.current_token_type == TOKEN_ELSE || parser.current_token_type == TOKEN_END {
							fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: if statement elif body is empty", parser.line, parser.column))
							os.Exit(0)
						}
						ElifBody := ParserParse(parser)
						ElifBodys = append(ElifBodys, ElifBody)
						ElifOps = append(ElifOps, ElifOp)
						if parser.current_token_type != TOKEN_ELIF {
							break
						}
					}
				} else {
					ElifBodys = nil
					ElifOps = nil
				}
				if parser.current_token_type == TOKEN_ELSE {
					parser.ParserEat(TOKEN_ELSE)
					if parser.current_token_type == TOKEN_ELSE || parser.current_token_type == TOKEN_END {
						fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: if statement body is empty", parser.line, parser.column))
						os.Exit(0)
					}
					ElseBody := ParserParse(parser)
					parser.ParserEat(TOKEN_END)
					expr.AsIf = &If{
						Op: op,
						Body: body,
						ElifBodys: ElifBodys,
						ElifOps: ElifOps,
						ElseBody: ElseBody,
					}
					exprs = append(exprs, expr)
				} else {
					parser.ParserEat(TOKEN_END)
					expr.AsIf = &If{
						Op: op,
						Body: body,
						ElifBodys: ElifBodys,
						ElifOps: ElifOps,
					}
					exprs = append(exprs, expr)
				}
			} else if parser.current_token_value == "call" {
				parser.ParserEat(TOKEN_ID)
				if parser.current_token_type != TOKEN_ID {
					fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value '%s'", parser.line, parser.column, parser.current_token_value))
					os.Exit(0)
				}
				expr.Type = ExprCall
				expr.AsCall = &Call{
					Value: parser.current_token_value,
				}
				parser.ParserEat(TOKEN_ID)
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "break" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprBreak
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "append" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprAppend
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "replace" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprReplace
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "read" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprRead
				exprs = append(exprs, expr)
			} else {
				expr.Type = ExprPush
				expr.AsPush = &Push{
					Arg: ParserParseExpr(parser),
				}
				exprs = append(exprs, expr)
			}
		} else if parser.current_token_type == TOKEN_PLUS {
			expr.Type = ExprBinop
			expr.AsBiniop = TOKEN_PLUS
			parser.ParserEat(TOKEN_PLUS)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_MINUS {
			expr.Type = ExprBinop
			expr.AsBiniop = TOKEN_MINUS
			parser.ParserEat(TOKEN_MINUS)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_DIV {
			expr.Type = ExprBinop
			expr.AsBiniop = TOKEN_DIV
			parser.ParserEat(TOKEN_DIV)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_MUL {
			expr.Type = ExprBinop
			expr.AsBiniop = TOKEN_MUL
			parser.ParserEat(TOKEN_MUL)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_REM {
			expr.Type = ExprBinop
			expr.AsBiniop = TOKEN_REM
			parser.ParserEat(TOKEN_REM)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_IS_EQUALS {
			expr.Type = ExprCompare
			expr.AsCompare = TOKEN_IS_EQUALS
			parser.ParserEat(TOKEN_IS_EQUALS)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_NOT_EQUALS {
			expr.Type = ExprCompare
			expr.AsCompare = TOKEN_NOT_EQUALS
			parser.ParserEat(TOKEN_NOT_EQUALS)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_LESS_THAN {
			expr.Type = ExprCompare
			expr.AsCompare = TOKEN_LESS_THAN
			parser.ParserEat(TOKEN_LESS_THAN)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_GREATER_THAN {
			expr.Type = ExprCompare
			expr.AsCompare = TOKEN_GREATER_THAN
			parser.ParserEat(TOKEN_GREATER_THAN)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_GREATER_EQUALS {
			expr.Type = ExprCompare
			expr.AsCompare = TOKEN_GREATER_EQUALS
			parser.ParserEat(TOKEN_GREATER_EQUALS)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_LESS_EQUALS {
			expr.Type = ExprCompare
			expr.AsCompare = TOKEN_LESS_EQUALS
			parser.ParserEat(TOKEN_LESS_EQUALS)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_EQUALS {
			parser.ParserEat(TOKEN_EQUALS)
			expr.Type = ExprVardef
			expr.AsVardef = &Vardef {
				Name: parser.current_token_value,
			}
			parser.ParserEat(TOKEN_ID)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_INT || parser.current_token_type == TOKEN_STRING || parser.current_token_type == TOKEN_L_BRACKET || parser.current_token_type == TOKEN_TYPE || parser.current_token_type == TOKEN_BOOL {
			expr.Type = ExprPush
			expr.AsPush = &Push{
				Arg: ParserParseExpr(parser),
			}
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_END || parser.current_token_type == TOKEN_ELSE || parser.current_token_type == TOKEN_DO || parser.current_token_type == TOKEN_EOF || parser.current_token_type == TOKEN_ELIF {
			return exprs
		} else {
			fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value '%s'", parser.line, parser.column, parser.current_token_value))
			os.Exit(0)
		}
	}

	return exprs
}


// -----------------------------
// ----------- Stack -----------
// -----------------------------

var Stack = []Expr{}

func VisitVar(VarName string, expr Expr) (Expr) {
	var VisitedVar Expr
	if _, ok := VariableScope[VarName]; ok {
		VisitedVar = VariableScope[VarName]
	} else {
		fmt.Println("Error: undefined variable '" + VarName + "'"); os.Exit(0);
	}
	return VisitedVar
}

func OpPush(item Expr) {
	if item.Type == ExprId {
		item = VisitVar(item.AsId.Name, item)
	}
	Stack = append(Stack, item)
}

func OpDrop() {
	if len(Stack)-1 < 0 {
		fmt.Println("DropError: the stack is empty.")
		os.Exit(0)
	}

	Stack = Stack[:len(Stack)-1]
}

func OpDup() {
	if len(Stack) < 1 {
		fmt.Println("Error: 'dup' expected more than one element in stack")
		os.Exit(0)
	}

	visitedExpr := Stack[len(Stack)-1]
	Stack = append(Stack, visitedExpr)
}

func OpSwap() {
	if len(Stack) < 2 {
		fmt.Println("SwapError: expected more than two elements in stack")
		os.Exit(0)
	}
	visitedExpr := Stack[len(Stack)-1]
	visitedExprSecond := Stack[len(Stack)-2]
	OpDrop()
	OpDrop()
	OpPush(visitedExpr)
	OpPush(visitedExprSecond)
}

func OpOver() {
	if len(Stack) < 2 {
		fmt.Println("OverError: expected more than two elements in stack.")
		os.Exit(0)
	}
	visitedExprSecond := Stack[len(Stack)-2]
	OpPush(visitedExprSecond)
}

func OpRot() {
	if len(Stack) < 3 {
		fmt.Println("Error: 'rot' expected more than three elements in stack.")
		os.Exit(0)
	}
	visitedExpr := Stack[len(Stack)-1]
	visitedExprSecond := Stack[len(Stack)-2]
	visitedExprThird := Stack[len(Stack)-3]
	OpDrop()
	OpDrop()
	OpDrop()
	OpPush(visitedExprSecond)
	OpPush(visitedExpr)
	OpPush(visitedExprThird)
}

func OpInc() {
	if len(Stack) < 1 {
		fmt.Println("Error: 'inc' expected more than one element in stack.")
		os.Exit(0)
	}
	visitedExpr := Stack[len(Stack)-1]
	if visitedExpr.Type != ExprInt {
		fmt.Println("TypeError: 'inc' expected type int")
		os.Exit(0)
	}
	visitedExpr.AsInt++
	OpDrop()
	OpPush(visitedExpr)
}

func OpDec() {
	if len(Stack) < 1 {
		fmt.Println("Error: 'dec' expected more than one element in stack.")
		os.Exit(0)
	}
	visitedExpr := Stack[len(Stack)-1]
	if visitedExpr.Type != ExprInt {
		fmt.Println("TypeError: 'dec' expected type int")
		os.Exit(0)
	}
	visitedExpr.AsInt--
	OpDrop()
	OpPush(visitedExpr)
}

func PrintArray(visitedExpr Expr) {
	fmt.Print("[")
	for i := 0; i < len(visitedExpr.AsArr); i++ {
		if i != 0 {
			fmt.Print(", ")
		}
		switch (visitedExpr.AsArr[i].Type) {
			case ExprInt: fmt.Print(visitedExpr.AsArr[i].AsInt)
			case ExprStr: fmt.Print(fmt.Sprintf("'%s'", visitedExpr.AsArr[i].AsStr))
			case ExprTypeType: fmt.Print(visitedExpr.AsArr[i].AsType)
			case ExprBool: fmt.Print(visitedExpr.AsArr[i].AsBool)
			case ExprArr: PrintArray(visitedExpr.AsArr[i])
		}
	}
	fmt.Print("]")
}

func OpPuts() {
	if len(Stack) < 1 {
		fmt.Println("Error: 'print' expected more than one element in stack.")
		os.Exit(0)
	}

	visitedExpr := Stack[len(Stack)-1]
	switch (visitedExpr.Type) {
		case ExprInt: fmt.Print(visitedExpr.AsInt)
		case ExprStr: fmt.Print(visitedExpr.AsStr)
		case ExprBool: fmt.Print(visitedExpr.AsBool)
		case ExprTypeType: fmt.Print(fmt.Sprintf("<%s>",visitedExpr.AsType))
		case ExprArr: PrintArray(visitedExpr)
	}
	OpDrop()
}

func OpPrint() {
	OpPuts()
	fmt.Println()
}

func OpPrintS() {
	fmt.Print(fmt.Sprintf("<%d> ", len(Stack)))
	for i:=len(Stack); i > 0; i-- {
		visitedExpr := Stack[len(Stack)-i]
		switch (visitedExpr.Type) {
			case ExprInt: fmt.Print(visitedExpr.AsInt)
			case ExprStr: fmt.Print(visitedExpr.AsStr)
			case ExprBool: fmt.Print(visitedExpr.AsBool)
			case ExprTypeType: fmt.Print(fmt.Sprintf("<%s>",visitedExpr.AsType))
			case ExprArr: PrintArray(visitedExpr)
		}
		fmt.Print(" ")
	}
	fmt.Println("← top")
}

func OpPrintV() {
	for name, value := range VariableScope {
		fmt.Print(fmt.Sprintf("%s : ", name))
		switch (value.Type) {
			case ExprInt: fmt.Print(value.AsInt)
			case ExprStr: fmt.Print(value.AsStr)
			case ExprBool: fmt.Print(value.AsBool)
			case ExprTypeType: fmt.Print(fmt.Sprintf("<%s>",value.AsType))
			case ExprArr: PrintArray(value)
		}
		fmt.Println()
	}
}

func OpInput() {
	var input string
	fmt.Scanln(&input)
	inpExpr := Expr{}
	inpExpr.Type = ExprStr
	inpExpr.AsStr = input
	OpPush(inpExpr)
}


func OpTypeOf() {
	if len(Stack) == 0 {
		fmt.Println("Error: 'typeof' expected more than one element in stack")
		os.Exit(0)
	}

	visitedExpr := Stack[len(Stack)-1]
	OpDrop()
	TypeExpr := Expr{}
	TypeExpr.Type = ExprTypeType
	var type_value string
	if visitedExpr.Type == ExprStr {
		type_value = "string"
	} else if visitedExpr.Type == ExprInt {
		type_value = "int"
	} else if visitedExpr.Type == ExprBool {
		type_value = "bool"
	} else if visitedExpr.Type == ExprTypeType {
		type_value = "type"
	} else if  visitedExpr.Type == ExprArr {
		type_value = "list"
	}
	TypeExpr.AsType = type_value
	OpPush(TypeExpr)
}

func OpCompare(value int) (bool) {
	if len(Stack) < 2 {
		fmt.Println("Error: expected more than two elements in stack.")
		os.Exit(0)
	}

	visitedExpr := Stack[len(Stack)-1]
	visitedExprSecond := Stack[len(Stack)-2]

	OpDrop()
	OpDrop()

	if value == TOKEN_IS_EQUALS {
		if visitedExpr.Type != visitedExprSecond.Type {
			return false
		}

		if visitedExpr.Type == ExprInt {
			return visitedExpr.AsInt == visitedExprSecond.AsInt
		}

		if visitedExpr.Type == ExprStr {
			return visitedExpr.AsStr == visitedExprSecond.AsStr
		}

		if visitedExpr.Type == ExprBool {
			return visitedExpr.AsBool == visitedExprSecond.AsBool
		}

		if visitedExpr.Type == ExprTypeType {
			return visitedExpr.AsType == visitedExprSecond.AsType
		}

		if visitedExpr.Type == ExprArr {
			return reflect.DeepEqual(visitedExpr.AsArr, visitedExprSecond.AsArr)
		}
	}

	if value == TOKEN_NOT_EQUALS {
		if visitedExpr.Type != visitedExprSecond.Type {
			return true
		}

		if visitedExpr.Type == ExprInt {
			return visitedExpr.AsInt != visitedExprSecond.AsInt
		}

		if visitedExpr.Type == ExprStr {
			return visitedExpr.AsStr != visitedExprSecond.AsStr
		}

		if visitedExpr.Type == ExprBool {
			return visitedExpr.AsBool != visitedExprSecond.AsBool
		}

		if visitedExpr.Type == ExprTypeType {
			return visitedExpr.AsType != visitedExprSecond.AsType
		}

		if visitedExpr.Type == ExprArr {
			return !reflect.DeepEqual(visitedExpr.AsArr, visitedExprSecond.AsArr)
		}
	}
    
	if visitedExpr.Type != ExprInt || visitedExprSecond.Type != ExprInt {
		fmt.Println("TypeError: '<' expected type int")
		os.Exit(0)
	}

	if value == TOKEN_LESS_THAN {
		return visitedExprSecond.AsInt < visitedExpr.AsInt
	}

	if value == TOKEN_GREATER_THAN {
		return visitedExprSecond.AsInt > visitedExpr.AsInt
	}

	if value == TOKEN_GREATER_EQUALS {
		return visitedExprSecond.AsInt >= visitedExpr.AsInt
	}

	if value == TOKEN_LESS_EQUALS {
		return visitedExprSecond.AsInt <= visitedExpr.AsInt
	}

	return false
}

func OpLen() {
	if len(Stack) < 1 {
		fmt.Println("Error: 'len' expected more than one elements in stack.")
		os.Exit(0)
	}

	visitedExpr := Stack[len(Stack)-1]

	if visitedExpr.Type != ExprArr {
		fmt.Println("TypeError: 'len' expected type <list>")
		os.Exit(0)
	}
	
	IntExpr := Expr{}
	IntExpr.Type = ExprInt
	IntExpr.AsInt = float64(len(visitedExpr.AsArr))
	OpPush(IntExpr)
}

func RetBool() (bool) {
	if len(Stack)-1 < 0 {
		fmt.Println("Error: the stack is empty, couldn't find bool")
		os.Exit(0)
	}

	visitedExpr := Stack[len(Stack)-1]
	if visitedExpr.Type != ExprBool {
		fmt.Println("Error: if op should be bool")
		os.Exit(0)
	}
	OpDrop()
	return visitedExpr.AsBool
}

func OpIf(expr Expr) (bool) {
	VisitExpr(expr.AsIf.Op)
	BoolValue := RetBool()
	if BoolValue {
		return VisitExpr(expr.AsIf.Body)
	} else {
		if expr.AsIf.ElifBodys != nil {
			i := 0
			for _, elifOp := range(expr.AsIf.ElifOps) {
				VisitExpr(elifOp)
				BoolValue = RetBool()
				if BoolValue {
					return VisitExpr(expr.AsIf.ElifBodys[i])
				}
				i++
			}
		}
		if expr.AsIf.ElseBody != nil {
			return VisitExpr(expr.AsIf.ElseBody)
		}
	}
	return false
}

func OpCondition(expr Expr) {
	bool_value := OpCompare(expr.AsCompare)
	BoolExpr := Expr{}
	BoolExpr.Type = ExprBool
	BoolExpr.AsBool = bool_value
	OpPush(BoolExpr)
}

func OpBinop(value int) {
	if len(Stack) < 2 {
		fmt.Print("Error: ")
		switch (value) {
			case TOKEN_PLUS: fmt.Print("'+'")
			case TOKEN_MINUS: fmt.Print("'-'")
			case TOKEN_DIV: fmt.Print("'/'")
			case TOKEN_REM: fmt.Print("'%'")
			case TOKEN_MUL: fmt.Print("'*'")

		}
		fmt.Println(" expected more than two elements in stack")
		os.Exit(0)
	}

	visitedExpr := Stack[len(Stack)-1]
	visitedExprSecond := Stack[len(Stack)-2]
	OpDrop()
	OpDrop()

	ValueExpr := Expr{}
	if value == TOKEN_PLUS {
		if visitedExpr.Type == ExprStr && visitedExprSecond.Type == ExprStr {
			ValueExpr.Type = ExprStr
			ValueExpr.AsStr =  visitedExprSecond.AsStr + visitedExpr.AsStr
		} else if visitedExpr.Type == ExprInt && visitedExprSecond.Type == ExprInt {
			ValueExpr.Type = ExprInt
			ValueExpr.AsInt = visitedExpr.AsInt + visitedExprSecond.AsInt
		} else {
			fmt.Println("TypeError: binary operation expected type int")
			os.Exit(0)
		}
	} else if visitedExpr.Type != ExprInt && visitedExprSecond.Type != ExprInt {
		fmt.Println("TypeError: binary operation expected type int")
		os.Exit(0)
	} else {
		ValueExpr.Type = ExprInt
		if value == TOKEN_MINUS {
			ValueExpr.AsInt = visitedExprSecond.AsInt - visitedExpr.AsInt
		} else if value == TOKEN_MUL {
			ValueExpr.AsInt = visitedExpr.AsInt * visitedExprSecond.AsInt
		} else if value == TOKEN_DIV {
			ValueExpr.AsInt = visitedExprSecond.AsInt / visitedExpr.AsInt
		} else if value == TOKEN_REM {
			if visitedExprSecond.AsInt == math.Trunc(visitedExprSecond.AsInt) && visitedExpr.AsInt == math.Trunc(visitedExpr.AsInt) {
				ValueExpr.AsInt = float64(int(visitedExprSecond.AsInt) % int(visitedExpr.AsInt))
			} else {
				fmt.Println("Error: operator '%' not defined on float")
				os.Exit(0)
			}
		}
	}

	OpPush(ValueExpr)
}

func OpImport(expr Expr) {
	if _, err := os.Stat(expr.AsImport); os.IsNotExist(err) {
		fmt.Println(fmt.Sprintf("ImportError: file path '%s' does not exist.",expr.AsImport))
		os.Exit(0)
	}
	file, err := os.Open(expr.AsImport)
	if err != nil {
		panic(err)
	}
	lexer := LexerInit(file)
	parser := ParserInit(lexer)
	exprs := ParserParse(parser)
	VisitExpr(exprs)
}

func OpFor(expr Expr) {
	VisitExpr(expr.AsFor.Op)
	for RetBool() {
		BreakValue := VisitExpr(expr.AsFor.Body)
		if BreakValue == true {break}
		VisitExpr(expr.AsFor.Op)
	}
}

func OpAppend(expr Expr) {
	if len(Stack) < 2 {
		fmt.Println("Error: 'append' expected more than two element in stack."); os.Exit(0);
	}
	visitedList := Stack[len(Stack)-2]
	visitedExpr := Stack[len(Stack)-1]
	if visitedList.Type != ExprArr {
		fmt.Println("TypeError: 'append' expected type <list>"); os.Exit(0);
	}
	visitedList.AsArr = append(visitedList.AsArr, visitedExpr)
	OpDrop()
	OpDrop()
	OpPush(visitedList)
}

func OpReplace() {
	if len(Stack) < 3 {
		fmt.Println("Error: 'replace' expected more than three element in stack."); os.Exit(0);
	}
	visitedList := Stack[len(Stack)-3]
	visitedValue := Stack[len(Stack)-2]
	visitedIndex := Stack[len(Stack)-1]
	if visitedIndex.Type != ExprInt {
		fmt.Println("TypeError: 'replace' index expected type <int>"); os.Exit(0);
	}
	if visitedIndex.AsInt != math.Trunc(visitedIndex.AsInt) {
		fmt.Println("Error: list index must be integer not float"); os.Exit(0);
	}
	if int(visitedIndex.AsInt) >= len(visitedList.AsArr) {
		fmt.Println("Error: 'replace' index out of range."); os.Exit(0);
	}
	if visitedList.Type != ExprArr {
		fmt.Println("TypeError: 'replace' expected type <list>"); os.Exit(0);
	}
	visitedList.AsArr[int(visitedIndex.AsInt)] = visitedValue
	OpDrop()
	OpDrop()
	OpDrop()
	OpPush(visitedList)
}

func OpRead() {
	if len(Stack) < 2 {
		fmt.Println("Error: 'read' expected more than two element in stack."); os.Exit(0);
	}
	visitedList := Stack[len(Stack)-2]
	visitedIndex := Stack[len(Stack)-1]
	if visitedIndex.Type != ExprInt {
		fmt.Println("TypeError: 'replace' index expected type <int>"); os.Exit(0);
	}
	if visitedList.Type != ExprArr {
		fmt.Println("TypeError: 'replace' expected type <list>"); os.Exit(0);
	}
	if visitedIndex.AsInt != math.Trunc(visitedIndex.AsInt) {
		fmt.Println("Error: list index must be integer not float"); os.Exit(0);
	}
	OpDrop()
	ExprValue := visitedList.AsArr[int(visitedIndex.AsInt)]
	OpPush(ExprValue)
}


// -----------------------------
// ---------- Variable ---------
// -----------------------------

var VariableScope = map[string]Expr{}

func OpVardef(expr Expr) {
	if len(Stack) < 1 {
		fmt.Println("Error: variable definition expected more than one element in stack.")
		os.Exit(0)
	}
	exprValue := Stack[len(Stack)-1]
	VariableScope[expr.AsVardef.Name] = exprValue
	OpDrop()
}


// -----------------------------
// ----------- Block -----------
// -----------------------------

var BlockScope = map[string][]Expr{}

func OpBlockdef(expr Expr) {
	if _, ok := BlockScope[expr.AsBlockdef.Name]; ok {
		fmt.Println(fmt.Sprintf("Error: block '%s' is already defined", expr.AsBlockdef.Name))
		os.Exit(0)
	}
	BlockScope[expr.AsBlockdef.Name] = expr.AsBlockdef.Body
}

func OpCallBlock(expr Expr) {
	if _, ok := BlockScope[expr.AsCall.Value]; ok {
		BlockBody := BlockScope[expr.AsCall.Value]
		VisitExpr(BlockBody)
	} else {
		fmt.Println("Error: undefined block '" + expr.AsCall.Value + "'")
		os.Exit(0)
	}
}


// -----------------------------
// -------- Visit Exprs --------
// -----------------------------

func VisitExpr(exprs []Expr) (bool) {
	BreakValue := false
	for _, expr := range exprs {
		switch expr.Type {
			case ExprPush:
				OpPush(expr.AsPush.Arg)
			case ExprPrint:
				OpPrint()
			case ExprInput:
				OpInput()
			case ExprPuts:
				OpPuts()
			case ExprPrintS:
				OpPrintS()
			case ExprPrintV:
				OpPrintV()
			case ExprAppend:
				OpAppend(expr)
			case ExprReplace:
				OpReplace()
			case ExprRead:
				OpRead()
			case ExprTypeOf:
				OpTypeOf()
			case ExprSwap:
				OpSwap()
			case ExprOver:
				OpOver()
			case ExprRot:
				OpRot()
			case ExprInc:
				OpInc()
			case ExprDec:
				OpDec()
			case ExprImport:
				OpImport(expr)
			case ExprDup:
				OpDup()
			case ExprDrop:
				OpDrop()
			case ExprLen:
				OpLen()
			case ExprExit:
				os.Exit(0)
			case ExprBinop:
				OpBinop(expr.AsBiniop)
			case ExprCompare:
				OpCondition(expr)
			case ExprBlockdef:
				OpBlockdef(expr)
			case ExprCall:
				OpCallBlock(expr)
			case ExprIf:
				BreakValue = OpIf(expr)
			case ExprFor:
				OpFor(expr)
			case ExprVardef:
				OpVardef(expr)
			case ExprBreak:
				BreakValue = true
		}
		if BreakValue {
			break
		}
	}
	return BreakValue
}


// -----------------------------
// ----------- Main ------------
// -----------------------------

func Usage() {
	fmt.Println("Usage:")
	fmt.Println("  tsh <filename>.t#")
	os.Exit(0)
}


func main() {
	if len(os.Args) != 2 || os.Args[1] == "help" {
		Usage()
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error: file '" + os.Args[1] + "' does not exist")

		whilte := color.New(color.FgWhite)

		fmt.Print("Run ")
		boldWhite := whilte.Add(color.BgCyan)
		boldWhite.Print(" tsh help ")
		fmt.Println(" for usage")

		os.Exit(0)
	}

	lexer := LexerInit(file)
	parser := ParserInit(lexer)
	exprs := ParserParse(parser)
	VisitExpr(exprs)
}
