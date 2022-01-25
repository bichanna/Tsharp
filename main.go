package main

import (
	"fmt"
	"os"
	"bufio"
	"io"
	"unicode"
	"reflect"
	_"unsafe"
	"strconv"
	"github.com/fatih/color"
)


// TODO: rewrite later...
func (lexer *Lexer) PrintErrorLineAsString(line int, column int) {
	var context string
	var FirstContext string
	var ThirdContext string
	var FirstLineString string = fmt.Sprintf("%d | ", line-1)
	var LineString string = fmt.Sprintf("%d | ", line)
	var ThirdLineString string = fmt.Sprintf("%d | ", line+1)
	for {
		r, _, err := lexer.reader.ReadRune()
		if r == '\n' {
			lexer.resetPosition()
		} else {
			if line == lexer.pos.line {
				context = context + string(r)
			} 
			if line - 1 != 0 {
				if line - 1 == lexer.pos.line {
					FirstContext = FirstContext + string(r)
				}
			}
			if line + 1 == lexer.pos.line {
				ThirdContext = ThirdContext + string(r)
			}
		}
		if line + 2 == lexer.pos.line {
			break
		}
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			panic(err)
		}
	}
	first := FirstLineString + FirstContext
	middle := LineString + context
	third := ThirdLineString + ThirdContext
	fmt.Print("  ")
	fmt.Println(first)
	fmt.Print("  ")
	fmt.Println(middle)
	column += 5
	for i := 0; i < column; i++ {
		fmt.Print(" ")
	}
	color.Red("^")
	fmt.Print("  ")
	fmt.Println(third)
}


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
	TOKEN_ERROR
	TOKEN_EXCEPT
	TOKEN_OR
	TOKEN_AND
)

var tokens = []string{
	TOKEN_PLUS:           "+",
	TOKEN_MINUS:          "-",
	TOKEN_DIV:            "/",
	TOKEN_MUL:            "*",
	TOKEN_IS_EQUALS:      "==",
	TOKEN_NOT_EQUALS:     "!=",
	TOKEN_LESS_THAN:      "<",
	TOKEN_GREATER_THAN:   ">",
	TOKEN_LESS_EQUALS:    "<=",
	TOKEN_GREATER_EQUALS: ">=",
	TOKEN_REM:            "%",
	TOKEN_OR:             "||",
	TOKEN_AND:            "&&",
}

type Position struct {
	line int
	column int
}

type Lexer struct {
	pos Position
	reader *bufio.Reader
	FileName string
}

func LexerInit(reader io.Reader, FileName string) *Lexer {
	return &Lexer{
		pos:    Position {line: 1, column: 0},
		reader: bufio.NewReader(reader),
		FileName: FileName,
	}
}

func (lexer *Lexer) Lex() (Position, Token, string, string) {
	for {
		r, _, err := lexer.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				err = nil
				return lexer.pos, TOKEN_EOF, "EOF", lexer.FileName
			}
			panic(err)
		}
		lexer.pos.column++
		switch r {
			case '\n': lexer.resetPosition()
			case '+': return lexer.pos, TOKEN_PLUS, "+", lexer.FileName
			case '/': return lexer.pos, TOKEN_DIV, "/", lexer.FileName
			case '*': return lexer.pos, TOKEN_MUL, "*", lexer.FileName
			case '%': return lexer.pos, TOKEN_REM, "%", lexer.FileName
			case '{': return lexer.pos, TOKEN_L_BRACKET, "{", lexer.FileName
			case '}': return lexer.pos, TOKEN_R_BRACKET, "}", lexer.FileName
			case ',': return lexer.pos, TOKEN_COMMA, ",", lexer.FileName
			case '.': return lexer.pos, TOKEN_DOT, ".", lexer.FileName
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
						return lexer.pos, TOKEN_IS_EQUALS, "==", lexer.FileName
					}
				} else if r == '-' {
					r, _, err := lexer.reader.ReadRune()
					if err != nil {
						if err == io.EOF {
							return lexer.pos, TOKEN_MINUS, "-", lexer.FileName
						}
						panic(err)
					}
					lexer.pos.column++
					if r == '>' {
						return lexer.pos, TOKEN_EQUALS, "->", lexer.FileName
					} else {
						return lexer.pos, TOKEN_MINUS, "-", lexer.FileName
					}
				} else if r == '<' {
					r, _, err := lexer.reader.ReadRune()
					if err != nil {
						if err == io.EOF {
							return lexer.pos, TOKEN_LESS_THAN, "<", lexer.FileName
						}
						panic(err)
					}
					if r == '=' {
						lexer.pos.column++
						return lexer.pos, TOKEN_LESS_EQUALS, "<=", lexer.FileName
					} else {
						return lexer.pos, TOKEN_LESS_THAN, "<", lexer.FileName
					}
				} else if r == '|' {
					r, _, err := lexer.reader.ReadRune()
					if err != nil {
						if err == io.EOF {
							fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value `%s`.", lexer.pos.line, lexer.pos.column, string(r)))
							os.Exit(0)
						}
						err = nil
						fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value `%s`.", lexer.pos.line, lexer.pos.column, string(r)))
						os.Exit(0)
						panic(err)
					}
					if r == '|' {
						lexer.pos.column++
						return lexer.pos, TOKEN_OR, "||", lexer.FileName
					} else {
						fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value `%s`.", lexer.pos.line, lexer.pos.column, string(r)))
						os.Exit(0)
					}
				} else if r == '&' {
					r, _, err := lexer.reader.ReadRune()
					if err != nil {
						if err == io.EOF {
							fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value `%s`.", lexer.pos.line, lexer.pos.column, string(r)))
							os.Exit(0)
						}
						err = nil
						fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value `%s`.", lexer.pos.line, lexer.pos.column, string(r)))
						os.Exit(0)
						panic(err)
					}
					if r == '&' {
						lexer.pos.column++
						return lexer.pos, TOKEN_AND, "&&", lexer.FileName
					} else {
						fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value `%s`.", lexer.pos.line, lexer.pos.column, string(r)))
						os.Exit(0)
					}
				} else if r == '>' {
					r, _, err := lexer.reader.ReadRune()
					if err != nil {
						if err == io.EOF {
							return lexer.pos, TOKEN_GREATER_THAN, ">", lexer.FileName
						}
						panic(err)
					}
					if r == '=' {
						lexer.pos.column++
						return lexer.pos, TOKEN_GREATER_EQUALS, ">=", lexer.FileName
					} else {
						return lexer.pos, TOKEN_GREATER_THAN, ">", lexer.FileName
					}
				} else if r == '!' {
					r, _, err := lexer.reader.ReadRune()
					if err != nil {panic(err)}
					lexer.pos.column++
					if r == '=' {
						return lexer.pos, TOKEN_NOT_EQUALS, "!=", lexer.FileName
					}
				} else if r == '#' {
					for {
						r, _, err := lexer.reader.ReadRune()
						if err != nil {
							if err == io.EOF {
								err = nil
								return lexer.pos, TOKEN_EOF, "EOF", lexer.FileName
							}
							panic(err)
						}
						if r == '\n' {break}
						if err != nil {panic(err)}
						lexer.pos.column++
					}
					continue
				} else if unicode.IsDigit(r) {
					startPos := lexer.pos
					lexer.backup()
					val := lexer.lexInt()
					return startPos, TOKEN_INT, val, lexer.FileName
				} else if unicode.IsLetter(r) {
					startPos := lexer.pos
					lexer.backup()
					val := lexer.lexId()
					if val == "end" {
						return startPos, TOKEN_END, val, lexer.FileName
					} else if val == "do" {
						return startPos, TOKEN_DO, val, lexer.FileName
					} else if val == "true" || val == "false" {
						return startPos, TOKEN_BOOL, val, lexer.FileName
					} else if val == "string" || val == "int" || val == "bool" || val == "type" || val == "list" || val == "error" {
						return startPos, TOKEN_TYPE, val, lexer.FileName
					} else if val == "else" {
						return startPos, TOKEN_ELSE, val, lexer.FileName
					} else if val == "elif" {
						return startPos, TOKEN_ELIF, val, lexer.FileName
					} else if val == "NameError" || val == "StackIndexError" || val == "TypeError" || val == "ImportError" || val == "IndexError" || val == "AssertionError" {
						return startPos, TOKEN_ERROR, val, lexer.FileName
					} else if val == "except" {
						return startPos, TOKEN_EXCEPT, val, lexer.FileName
					}
					return startPos, TOKEN_ID, val, lexer.FileName
				} else if r == '"' {
					startPos := lexer.pos
					lexer.backup()
					lexer.pos.column++
					val := lexer.lexString()
					lexer.pos.column++
					r, _, err = lexer.reader.ReadRune()
					return startPos, TOKEN_STRING, val, lexer.FileName
				} else if string(r) == "'" {
					startPos := lexer.pos
					lexer.backup()
					lexer.pos.column++
					val := lexer.lexStringSingle()
					lexer.pos.column++
					r, _, err = lexer.reader.ReadRune()
					return startPos, TOKEN_STRING, val, lexer.FileName
				} else {
					line := lexer.pos.line
					col := lexer.pos.column
					file, err := os.Open(lexer.FileName)
					if err != nil {
						panic(err)
					}
					lexer := LexerInit(file, lexer.FileName)
					fmt.Println(fmt.Sprintf("%s SyntaxError:%d:%d: unexpected token value `%s`.", lexer.FileName, lexer.pos.line, lexer.pos.column, string(r)))
					lexer.PrintErrorLineAsString(line, col)
					os.Exit(0)
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
		} else if unicode.IsDigit(r) {
			val = val + string(r)
		} else if unicode.IsSpace(r) {
			lexer.backup()
			return val
		} else if string(r) == "_" {
			val = val + string(r)
		} else if string(r) == "-" {
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

func (lexer *Lexer) lexStringSingle() string {
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
		if string(r) != "'" {
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
// ---------- Errors -----------
// -----------------------------

type ErrorType int
const (
	ErrorVoid ErrorType = iota
	StackIndexError
	NameError
	TypeError
	IndexError
	ImportError
	AssertionError
)

type Error struct {
    message string
	Type ErrorType
}


// -----------------------------
// ------------ AST ------------
// -----------------------------

type AsStr struct {
	StringValue string
}

func (node AsStr) node() {}

type AsInt struct {
	IntValue int
}

func (node AsInt) node() {}

type AsBool struct {
	BoolValue bool
}

func (node AsBool) node() {}

type NewList struct {
	ListBody AST
}

func (node NewList) node() {}

type AsList struct {
	ListArgs []AST
}

func (node AsList) node() {}

type AsId struct {
	name string
}

func (node AsId) node() {}

type Import struct {
	FileName string
}

func (node Import) node() {}

type Assert struct {
	Line int
	Col int
	Message string
}

func (node Assert) node() {}

type Compare struct {
	op uint8
}

func (node Compare) node() {}

type AsError struct {
	err ErrorType
}

func (node AsError) node() {}

type AsBinop struct {
	op uint8
}

func (node AsBinop) node() {}

type AsPush struct {
	value AST
}

func (node AsPush) node() {}

type AsType struct {
	TypeValue string
}

func (node AsType) node() {}

type Vardef struct {
	Name string
}

func (node Vardef) node() {}

type Var struct {
	Name string
}

func (node Var) node() {}

type Blockdef struct {
	Name string
	BlockBody AST
}

func (node Blockdef) node() {}

type CallBlock struct {
	Name string
}

func (node CallBlock) node() {}

type If struct {
	IfOp AST
	IfBody AST
	ElifOps []AST
	ElifBodys []AST
	ElseBody AST
}

func (node If) node() {}

type For struct {
	ForOp AST
	ForBody AST
}

func (node For) node() {}

type Try struct {
	TryBody AST
	ExceptErrors []AST
	ExceptBodys []AST
}

func (node Try) node() {}

type AsStatements []AST

func (node AsStatements) node() {}

type AST interface {
	node()
}

// -----------------------------
// ---------- Parser -----------
// -----------------------------

type Parser struct {
	current_token_type Token
	current_token_value string
	FileName string
	lexer Lexer
	line int
	column int
}


func ParserInit(lexer *Lexer) *Parser {
	pos, tok, val, file := lexer.Lex()
	return &Parser{
		current_token_type: tok,
		current_token_value: val,
		FileName: file,
		lexer: *lexer,
		line: pos.line,
		column: pos.column,
	}
}

func (parser *Parser) ParserEat(token Token) {
	if token != parser.current_token_type {
		line := parser.line
		col := parser.column
		file, err := os.Open(parser.FileName)
		if err != nil {
			panic(err)
		}
		lexer := LexerInit(file, parser.FileName)
		fmt.Println(fmt.Sprintf("%s SyntaxError:%d:%d: unexpected token value `%s`.", parser.FileName, parser.line, parser.column, parser.current_token_value))
		lexer.PrintErrorLineAsString(line, col)
		os.Exit(0)
	}
	pos, tok, val, file := parser.lexer.Lex()
	parser.current_token_type = tok
	parser.current_token_value = val
	parser.FileName = file
	parser.line = pos.line
	parser.column = pos.column
}

func StrToInt(num string) int {
	i, err := strconv.Atoi(num)
	if err != nil{
		panic(err)
	}
	return i
}

func ParserParseError(parser *Parser) AST {
	var err ErrorType
	if parser.current_token_value == "StackIndexError" {
		err = StackIndexError
	} else if parser.current_token_value == "NameError" {
		err = NameError
	} else if parser.current_token_value == "ImportError" {
		err = ImportError
	} else if parser.current_token_value == "TypeError" {
		err = TypeError
	} else if parser.current_token_value == "IndexError" {
		err = IndexError
	} else if parser.current_token_value == "AssertionError" {
		err = AssertionError
	}
	parser.ParserEat(TOKEN_ERROR)
	ErrorExpr := AsError {
		err: err,
	}
	return ErrorExpr
}

func ParserParseExpr(parser *Parser) AST {
	var expr AST
	switch parser.current_token_type {
		case TOKEN_INT:
			expr = AsInt {
				StrToInt(parser.current_token_value),
			}
			parser.ParserEat(TOKEN_INT)
		case TOKEN_STRING:
			expr = AsStr {
				parser.current_token_value,
			}
			parser.ParserEat(TOKEN_STRING)
		case TOKEN_BOOL:
			BoolValue := parser.current_token_value == "true"
			expr = AsBool {
				BoolValue,
			}
			parser.ParserEat(TOKEN_BOOL)
		case TOKEN_ERROR:
			expr = ParserParseError(parser)
		case TOKEN_L_BRACKET:
			parser.ParserEat(TOKEN_L_BRACKET)
			var ListBody AST
			if parser.current_token_type != TOKEN_R_BRACKET {
				ListBody = ParserParse(parser)
			}
			expr = NewList {
				ListBody,
			}
			parser.ParserEat(TOKEN_R_BRACKET)
		case TOKEN_ID:
			expr = Var {
				parser.current_token_value,
			}
			parser.ParserEat(TOKEN_ID)
		case TOKEN_TYPE:
			expr = AsType {
				parser.current_token_value,
			}
			parser.ParserEat(TOKEN_TYPE)
		default:
			line := parser.line
			col := parser.column
			file, err := os.Open(parser.FileName)
			if err != nil {
				panic(err)
			}
			lexer := LexerInit(file, parser.FileName)
			fmt.Println(fmt.Sprintf("%s SyntaxError:%d:%d: unexpected token value `%s`.", parser.FileName, parser.line, parser.column, parser.current_token_value))
			lexer.PrintErrorLineAsString(line, col)
			os.Exit(0)
	}

	return expr
}

func ParserParse(parser *Parser) AST {
	var Statements AsStatements
	if  parser.current_token_type == TOKEN_DO || parser.current_token_type == TOKEN_END || parser.current_token_type == TOKEN_ELIF || parser.current_token_type == TOKEN_ELSE || parser.current_token_type == TOKEN_EXCEPT {
		line := parser.line
		col := parser.column
		file, err := os.Open(parser.FileName)
		if err != nil {
			panic(err)
		}
		lexer := LexerInit(file, parser.FileName)
		fmt.Println(fmt.Sprintf("%s SyntaxError:%d:%d: the body is empty, unexpected token value `%s`.", parser.FileName, parser.line, parser.column, parser.current_token_value))
		lexer.PrintErrorLineAsString(line, col)
		os.Exit(0)
	}
	for {
		if parser.current_token_type == TOKEN_ID {
			if parser.current_token_value == "print" || parser.current_token_value == "break" || parser.current_token_value == "append" || parser.current_token_value == "remove" || parser.current_token_value == "swap" || parser.current_token_value == "in" || parser.current_token_value == "typeof" || parser.current_token_value == "rot" || parser.current_token_value == "len" || parser.current_token_value == "input" ||
			   parser.current_token_value == "drop"  || parser.current_token_value == "dup" || parser.current_token_value == "inc" || parser.current_token_value == "dec" || parser.current_token_value == "replace" || parser.current_token_value == "read" || parser.current_token_value == "puts" || parser.current_token_value == "over" || parser.current_token_value == "printS" || parser.current_token_value == "exit" {
				name := parser.current_token_value
				IdExpr := AsId{name}
				parser.ParserEat(TOKEN_ID)
				Statements = append(Statements, IdExpr)
			} else if parser.current_token_value == "assert" {
				Line := parser.line
				Col := parser.column
				parser.ParserEat(TOKEN_ID)
				AssertExpr := Assert {
					Line: Line,
					Col: Col,
					Message: parser.current_token_value,
				}
				parser.ParserEat(TOKEN_STRING)
				Statements = append(Statements, AssertExpr)
			} else if parser.current_token_value == "block" {
				parser.ParserEat(TOKEN_ID)
				name := parser.current_token_value
				parser.ParserEat(TOKEN_ID)
				parser.ParserEat(TOKEN_DO)
				BlockBody := ParserParse(parser)
				parser.ParserEat(TOKEN_END)
				BlockdefExpr := Blockdef {
					Name: name,
					BlockBody: BlockBody,
				}
				Statements = append(Statements, BlockdefExpr)
			} else if parser.current_token_value == "call" {
				parser.ParserEat(TOKEN_ID)
				name := parser.current_token_value
				parser.ParserEat(TOKEN_ID)
				CallBlockExpr := CallBlock {
					name,
				}
				Statements = append(Statements, CallBlockExpr)
			} else if parser.current_token_value == "import" {
				parser.ParserEat(TOKEN_ID)
				ImportExpr := Import {
					parser.current_token_value,
				}
				parser.ParserEat(TOKEN_STRING)
				Statements = append(Statements, ImportExpr)
			} else if parser.current_token_value == "if" {
				parser.ParserEat(TOKEN_ID)
				IfOp := ParserParse(parser)
				parser.ParserEat(TOKEN_DO)
				IfBody := ParserParse(parser)
				var ElifOps []AST
				var ElifBodys []AST
				for {
					if parser.current_token_type != TOKEN_ELIF {
						break
					}
					parser.ParserEat(TOKEN_ELIF)
					ElifOp := ParserParse(parser)
					ElifOps = append(ElifOps, ElifOp)
					parser.ParserEat(TOKEN_DO)
					ElifBody := ParserParse(parser)
					ElifBodys = append(ElifBodys, ElifBody)
				}
				var ElseBody AST = nil
				for {
					if parser.current_token_type != TOKEN_ELSE {
						break
					}
					parser.ParserEat(TOKEN_ELSE)
					ElseBody = ParserParse(parser)
				}
				parser.ParserEat(TOKEN_END)
				IfExpr := If {
					IfOp: IfOp,
					IfBody: IfBody,
					ElifOps: ElifOps,
					ElifBodys: ElifBodys,
					ElseBody: ElseBody,
				}
				Statements = append(Statements, IfExpr)
			} else if parser.current_token_value == "for" {
				parser.ParserEat(TOKEN_ID)
				ForOp := ParserParse(parser)
				parser.ParserEat(TOKEN_DO)
				ForBody := ParserParse(parser)
				parser.ParserEat(TOKEN_END)
				ForExpr := For {
					ForOp: ForOp,
					ForBody: ForBody,
				}
				Statements = append(Statements, ForExpr)
			} else if parser.current_token_value == "try" {
				parser.ParserEat(TOKEN_ID)
				TryBody := ParserParse(parser)
				var ExceptErrors []AST
				var ExceptBodys []AST
				for {
					if parser.current_token_type != TOKEN_EXCEPT {
						break
					}
					parser.ParserEat(TOKEN_EXCEPT)
					ExceptError := ParserParseError(parser)
					ExceptErrors = append(ExceptErrors, ExceptError)
					parser.ParserEat(TOKEN_DO)
					ExceptBody := ParserParse(parser)
					ExceptBodys = append(ExceptBodys, ExceptBody)
				}
				parser.ParserEat(TOKEN_END)
				TryExpr := Try {
					TryBody: TryBody,
					ExceptErrors: ExceptErrors,
					ExceptBodys: ExceptBodys,
				}
				Statements = append(Statements, TryExpr)
			} else {
				expr := ParserParseExpr(parser)
				PushExpr := AsPush{
					value: expr,
				}
				Statements = append(Statements, PushExpr)
			}
		} else if parser.current_token_type == TOKEN_INT  || parser.current_token_type == TOKEN_STRING ||
		    parser.current_token_type == TOKEN_BOOL || parser.current_token_type == TOKEN_ERROR || parser.current_token_type == TOKEN_L_BRACKET || parser.current_token_type == TOKEN_TYPE {
			expr := ParserParseExpr(parser)
			PushExpr := AsPush{
				value: expr,
			}
			Statements = append(Statements, PushExpr)
		} else if parser.current_token_type == TOKEN_EQUALS {
			parser.ParserEat(TOKEN_EQUALS)
			VardefExpr := Vardef {
				Name: parser.current_token_value,
			}
			parser.ParserEat(TOKEN_ID)
			Statements = append(Statements, VardefExpr)
		} else if parser.current_token_type == TOKEN_PLUS || parser.current_token_type == TOKEN_MINUS || parser.current_token_type == TOKEN_MUL || parser.current_token_type == TOKEN_DIV || parser.current_token_type == TOKEN_REM {
			BinopExpr := AsBinop {
				op: uint8(parser.current_token_type),
			}
			parser.ParserEat(parser.current_token_type)
			Statements = append(Statements, BinopExpr)
		} else if parser.current_token_type == TOKEN_LESS_EQUALS || parser.current_token_type == TOKEN_GREATER_EQUALS || parser.current_token_type == TOKEN_LESS_THAN || parser.current_token_type == TOKEN_GREATER_THAN || parser.current_token_type == TOKEN_IS_EQUALS || parser.current_token_type == TOKEN_NOT_EQUALS || parser.current_token_type == TOKEN_OR || parser.current_token_type == TOKEN_AND {
			CompareExpr := Compare {
				op: uint8(parser.current_token_type),
			}
			parser.ParserEat(parser.current_token_type)
			Statements = append(Statements, CompareExpr)
		} else if parser.current_token_type == TOKEN_EOF || parser.current_token_type == TOKEN_DO ||
		    parser.current_token_type == TOKEN_END || parser.current_token_type == TOKEN_ELIF ||
			parser.current_token_type == TOKEN_ELSE || parser.current_token_type == TOKEN_EXCEPT || parser.current_token_type == TOKEN_R_BRACKET {
			break
		} else {
			line := parser.line
			col := parser.column
			file, err := os.Open(parser.FileName)
			if err != nil {
				panic(err)
			}
			lexer := LexerInit(file, parser.FileName)
			fmt.Println(fmt.Sprintf("%s SyntaxError:%d:%d: unexpected token value `%s`.", parser.FileName, parser.line, parser.column, parser.current_token_value))
			lexer.PrintErrorLineAsString(line, col)
			os.Exit(0)
		}
	}
	return Statements
}


// -----------------------------
// ----------- Stack -----------
// -----------------------------

type Scope struct {
    Stack []AST
}

var Variables = map[string]AST{}
var Blocks = map[string]AST{}

func InitScope() *Scope {
	return &Scope{
		[]AST{},
	}
}

func (scope *Scope) OpPush(node AST, VariableScope *map[string]AST) (*Error) {
	_, IsList := node.(NewList);
	_, IsVar := node.(Var);
	if IsList {
		ListScope := InitScope()
		if node.(NewList).ListBody != nil {
			ListScope.VisitorVisit(node.(NewList).ListBody, false, VariableScope)
		}
		var expr AST = AsList {
			ListScope.Stack,
		}
		scope.Stack = append(scope.Stack, expr)
	} else if IsVar {
		if VariableScope == nil {
			if _, ok := Variables[node.(Var).Name]; ok {
				VarValue := Variables[node.(Var).Name]
				scope.Stack = append(scope.Stack, VarValue)
			} else {
				err := Error{}
				err.message = fmt.Sprintf("NameError: undefined variable `%s`", node.(Var).Name)
				err.Type = NameError
				return &err
			}
		} else {
			VariablesScope := *VariableScope
			if _, ok := VariablesScope[node.(Var).Name]; ok {
				VarValue := VariablesScope[node.(Var).Name]
				scope.Stack = append(scope.Stack, VarValue)
			} else {
				if _, ok := Variables[node.(Var).Name]; ok {
					VarValue := Variables[node.(Var).Name]
					scope.Stack = append(scope.Stack, VarValue)
				} else {
					err := Error{}
					err.message = fmt.Sprintf("NameError: undefined variable `%s`", node.(Var).Name)
					err.Type = NameError
					return &err
				}
			}
		}
	} else {
		scope.Stack = append(scope.Stack, node)
	}
	return nil
}

func (scope *Scope) OpDrop() (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = "StackIndexError: `drop` expected one or more element in the stack."
		err.Type = StackIndexError
		return &err
	}
	scope.Stack = scope.Stack[:len(scope.Stack)-1]
	return nil
}

func (scope *Scope) OpSwap() (*Error) {
	if len(scope.Stack) < 2 {
		err := Error{}
		err.message = "StackIndexError: `swap` expected two or more element in the stack."
		err.Type = StackIndexError
		return &err
	}
	first := scope.Stack[len(scope.Stack)-1]
	second := scope.Stack[len(scope.Stack)-2]
	scope.OpDrop()
	scope.OpDrop()
	scope.OpPush(first, nil)
	scope.OpPush(second, nil)
	return nil
}

func RetTokenAsStr(token uint8) string {
	return tokens[token]
}

func (scope *Scope) OpBinop(op uint8) (*Error) {
	if len(scope.Stack) < 2 {
		err := Error{}
		err.message = fmt.Sprintf("StackIndexError: `%s` expected more than 2 <int> type elements in the stack.", RetTokenAsStr(op))
		err.Type = StackIndexError
		return &err
	}
	first := scope.Stack[len(scope.Stack)-1]
	second := scope.Stack[len(scope.Stack)-2]
	scope.OpDrop()
	scope.OpDrop()
	_, ok := first.(AsInt);
	_, ok2 := second.(AsInt);

	_, IsStr := first.(AsStr);
	_, IsStr2 := second.(AsStr);

	if IsStr && IsStr2 {
		StrVal := second.(AsStr).StringValue + first.(AsStr).StringValue
		expr := AsStr {
			StrVal,
		}
		scope.OpPush(expr, nil)
		return nil
	}

	if !ok || !ok2 {
		err := Error{}
		err.message = fmt.Sprintf("StackIndexError: `%s` expected 2 <int> type or 2 <string> type elements in the stack.", RetTokenAsStr(op))
		err.Type = TypeError
		return &err
	}
	var val int
	switch op {
		case TOKEN_PLUS: val = first.(AsInt).IntValue + second.(AsInt).IntValue
		case TOKEN_MINUS:  val = second.(AsInt).IntValue - first.(AsInt).IntValue
		case TOKEN_MUL: val = second.(AsInt).IntValue * first.(AsInt).IntValue
		case TOKEN_DIV: val = second.(AsInt).IntValue / first.(AsInt).IntValue
		case TOKEN_REM: val = second.(AsInt).IntValue % first.(AsInt).IntValue
	}
	expr := AsInt {
		val,
	}
	scope.OpPush(expr, nil)
	return nil
}

func (scope *Scope) OpCompare(op uint8) (*Error) {
	if len(scope.Stack) < 2 {
		err := Error{}
		err.message = fmt.Sprintf("StackIndexError: `%s` expected more than 2 elements in the stack.", RetTokenAsStr(op))
		err.Type = StackIndexError
		return &err
	}
	first := scope.Stack[len(scope.Stack)-1]
	second := scope.Stack[len(scope.Stack)-2]
	var val bool
	if op == TOKEN_IS_EQUALS {
		if reflect.TypeOf(first) != reflect.TypeOf(second) {
			val = false
		} else {
			switch first.(type) {
				case AsStr: val = first.(AsStr).StringValue == second.(AsStr).StringValue
				case AsInt:  val = second.(AsInt).IntValue == first.(AsInt).IntValue
				case AsBool: val = second.(AsBool).BoolValue == first.(AsBool).BoolValue
				case AsType: val = second.(AsType).TypeValue == first.(AsType).TypeValue
				case AsError: val = second.(AsError).err == first.(AsError).err
				case AsList: val = reflect.DeepEqual(second.(AsList).ListArgs, first.(AsList).ListArgs)
			}
		}
	} else if op == TOKEN_NOT_EQUALS {
		if reflect.TypeOf(first) != reflect.TypeOf(second) {
			val = true
		} else {
			switch first.(type) {
				case AsStr: val = first.(AsStr).StringValue != second.(AsStr).StringValue
				case AsInt:  val = second.(AsInt).IntValue != first.(AsInt).IntValue
				case AsBool: val = second.(AsBool).BoolValue != first.(AsBool).BoolValue
				case AsType: val = second.(AsType).TypeValue != first.(AsType).TypeValue
				case AsError: val = second.(AsError).err != first.(AsError).err
				case AsList: val = !reflect.DeepEqual(second.(AsList).ListArgs, first.(AsList).ListArgs)
			}
		}
	} else if op == TOKEN_OR || op == TOKEN_AND {
		_, ok := first.(AsBool);
		_, ok2 := second.(AsBool);
		if !ok || !ok2 {
			err := Error{}
			err.message = fmt.Sprintf("StackIndexError: `%s` expected 2 <bool> type elements in the stack.", RetTokenAsStr(op))
			err.Type = TypeError
			return &err
		}
		switch op {
			case TOKEN_OR: val = second.(AsBool).BoolValue || first.(AsBool).BoolValue
			case TOKEN_AND: val = second.(AsBool).BoolValue && first.(AsBool).BoolValue
		}
	} else {
		_, ok := first.(AsInt);
		_, ok2 := second.(AsInt);
		if !ok || !ok2 {
			err := Error{}
			err.message = fmt.Sprintf("StackIndexError: `%s` expected 2 <int> type elements in the stack.", RetTokenAsStr(op))
			err.Type = TypeError
			return &err
		}
		switch op {
			case TOKEN_LESS_THAN: val = second.(AsInt).IntValue < first.(AsInt).IntValue
			case TOKEN_LESS_EQUALS: val = second.(AsInt).IntValue <= first.(AsInt).IntValue
			case TOKEN_GREATER_THAN: val = second.(AsInt).IntValue > first.(AsInt).IntValue
			case TOKEN_GREATER_EQUALS: val = second.(AsInt).IntValue >= first.(AsInt).IntValue
		}
	}
	scope.OpDrop()
	scope.OpDrop()
	expr := AsBool {
		val,
	}
	scope.OpPush(expr, nil)
	return nil
}

func PrintAsList(node AST) {
	print("{")
	for i := 0; i < len(node.(AsList).ListArgs); i++ {
		switch node.(AsList).ListArgs[i].(type) {
			case AsStr:
				print(node.(AsList).ListArgs[i].(AsStr).StringValue)
			case AsInt:
				print(strconv.Itoa(node.(AsList).ListArgs[i].(AsInt).IntValue))
			case AsBool:
				print(node.(AsList).ListArgs[i].(AsBool).BoolValue)
			case AsList:
				PrintAsList(node.(AsList).ListArgs[i])
		}
		if i < len(node.(AsList).ListArgs)-1 {
			print(", ")
		}
	}
	print("}")
}

func (scope *Scope) OpPrint() (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = "StackIndexError: `print` the stack is empty."
		err.Type = StackIndexError
		return &err
	}
	expr := scope.Stack[len(scope.Stack)-1]
	switch expr.(type) {
		case AsStr: fmt.Println(expr.(AsStr).StringValue)
		case AsInt: fmt.Println(expr.(AsInt).IntValue)
		case AsBool: fmt.Println(expr.(AsBool).BoolValue)
		case AsType: fmt.Println(fmt.Sprintf("<%s>" ,expr.(AsType).TypeValue))
		case AsError:
			switch expr.(AsError).err {
				case NameError: fmt.Println("<error 'NameError'>")
				case StackIndexError: fmt.Println("<error 'StackIndexError'>")
				case ImportError: fmt.Println("<error 'ImportError'>")
				case IndexError: fmt.Println("<error 'IndexError'>")
				case TypeError: fmt.Println("<error 'TypeError'>")
				default: fmt.Println(fmt.Sprintf("unexpected error <%d>", expr.(AsError).err))
			}
		case AsList:
			PrintAsList(expr)
			fmt.Println()
	}
	scope.OpDrop()
	return nil
}

func (scope *Scope) OpPuts() (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = "StackIndexError: `puts` the stack is empty."
		err.Type = StackIndexError
		return &err
	}
	expr := scope.Stack[len(scope.Stack)-1]
	switch expr.(type) {
		case AsStr: fmt.Print(expr.(AsStr).StringValue)
		case AsInt: fmt.Print(expr.(AsInt).IntValue)
		case AsBool: fmt.Print(expr.(AsBool).BoolValue)
		case AsType: fmt.Print(fmt.Sprintf("<%s>" ,expr.(AsType).TypeValue))
		case AsError:
			switch expr.(AsError).err {
				case NameError: fmt.Print("<error 'NameError'>")
				case StackIndexError: fmt.Print("<error 'StackIndexError'>")
				case ImportError: fmt.Print("<error 'ImportError'>")
				case IndexError: fmt.Print("<error 'IndexError'>")
				case TypeError: fmt.Print("<error 'TypeError'>")
				default: fmt.Print(fmt.Sprintf("unexpected error <%d>", expr.(AsError).err))
			}
		case AsList:
			PrintAsList(expr)
	}
	scope.OpDrop()
	return nil
}

func (scope *Scope) OpPrintS() {
	fmt.Print(fmt.Sprintf("<%d> ", len(scope.Stack)))
	for i:=len(scope.Stack); i > 0; i-- {
		expr := scope.Stack[len(scope.Stack)-i]
		switch expr.(type) {
			case AsStr: fmt.Print(expr.(AsStr).StringValue)
			case AsInt: fmt.Print(expr.(AsInt).IntValue)
			case AsBool: fmt.Print(expr.(AsBool).BoolValue)
			case AsType: fmt.Print(fmt.Sprintf("<%s>" ,expr.(AsType).TypeValue))
			case AsError:
				switch expr.(AsError).err {
					case NameError: fmt.Print("<error 'NameError'>")
					case StackIndexError: fmt.Print("<error 'StackIndexError'>")
					case ImportError: fmt.Print("<error 'ImportError'>")
					case IndexError: fmt.Print("<error 'IndexError'>")
					case TypeError: fmt.Print("<error 'TypeError'>")
					default: fmt.Print(fmt.Sprintf("unexpected error <%d>", expr.(AsError).err))
				}
			case AsList:
				PrintAsList(expr)
		}
		print(" ")
	}
}

func (scope *Scope) OpIf(node AST, IsTry bool, VariableScope *map[string]AST) (bool, *Error) {
	scope.VisitorVisit(node.(If).IfOp, IsTry, VariableScope)
	var BreakValue bool = false
	var err *Error = nil
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = "StackIndexError: if statement expected one or more <bool> type element in the stack."
		err.Type = StackIndexError
		return BreakValue, &err
	}
	expr := scope.Stack[len(scope.Stack)-1]
	if _, ok := expr.(AsBool); !ok {
		err := Error{}
		err.message = "TypeError: if statement expected one or more <bool> type element in the stack."
		err.Type = StackIndexError
		return BreakValue, &err
	}
	scope.OpDrop()
	if expr.(AsBool).BoolValue {
		BreakValue, err, _ = scope.VisitorVisit(node.(If).IfBody, IsTry, VariableScope)
		if err != nil {
			return BreakValue, err
		}
		return BreakValue, nil
	}
	for i := 0; i < len(node.(If).ElifOps); i++ {
		scope.VisitorVisit(node.(If).ElifOps[i], IsTry, VariableScope)
		if len(scope.Stack) < 1 {
			err := Error{}
			err.message = "StackIndexError: if statement expected one or more <bool> type element in the stack."
			err.Type = StackIndexError
			return BreakValue, &err
		}
		expr := scope.Stack[len(scope.Stack)-1]
		if _, ok := expr.(AsBool); !ok {
			err := Error{}
			err.message = "TypeError: if statement expected one or more <bool> type element in the stack."
			err.Type = StackIndexError
			return BreakValue, &err
		}
		scope.OpDrop()
		if expr.(AsBool).BoolValue {
			BreakValue, err, _ = scope.VisitorVisit(node.(If).ElifBodys[i], IsTry, VariableScope)
			return BreakValue, err
		}
	}
	if node.(If).ElseBody != nil {
		BreakValue, err, _ = scope.VisitorVisit(node.(If).ElseBody, IsTry, VariableScope)
	}
	return BreakValue, err
}

func (scope *Scope) OpFor(node AST, IsTry bool, VariableScope *map[string]AST) (*Error) {
	var BreakValue bool
	for {
		_, err, _ := scope.VisitorVisit(node.(For).ForOp, IsTry, VariableScope)
		if err != nil {
			return err
		}
		if len(scope.Stack) < 1 {
			err := Error{}
			err.message = "StackIndexError: for loop expected one or more <bool> type element in the stack."
			err.Type = StackIndexError
			return &err
		}
		expr := scope.Stack[len(scope.Stack)-1]
		if _, ok := expr.(AsBool); !ok {
			err := Error{}
			err.message = "TypeError: for loop expected one or more <bool> type element in the stack."
			err.Type = StackIndexError
			return &err
		}
		scope.OpDrop()
		if !expr.(AsBool).BoolValue {
			return nil
		}
		BreakValue, err, _ = scope.VisitorVisit(node.(For).ForBody, IsTry, VariableScope)
		if err != nil {
			return err
		}
		if BreakValue {
			return nil
		}
	}
	return nil
}

func (scope *Scope) OpTry(node AST, VariableScope *map[string]AST) (*Error) {
	_, err, _ := scope.VisitorVisit(node.(Try).TryBody, true, nil)
	if err != nil {
		for i := 0; i < len(node.(Try).ExceptErrors); i++ {
			if node.(Try).ExceptErrors[i].(AsError).err == err.Type {
				scope.VisitorVisit(node.(Try).ExceptBodys[i], false, VariableScope)
				return nil
			}
		}
	}
	return err
}

func (scope *Scope) OpInc() (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = "StackIndexError: `inc` expected one or more <int> type element in the stack."
		err.Type = StackIndexError
		return &err
	}
	first := scope.Stack[len(scope.Stack)-1]
	_, ok := first.(AsInt);
	if !ok {
		err := Error{}
		err.message = "StackIndexError: `inc` expected one or more <int> type element in the stack."
		err.Type = StackIndexError
		return &err
	}
	scope.OpDrop()
	val := first.(AsInt).IntValue + 1
	expr := AsInt {
		val,
	}
	scope.OpPush(expr, nil)
	return nil
}

func (scope *Scope) OpDec() (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = "StackIndexError: `dec` expected one or more <int> type element in the stack."
		err.Type = StackIndexError
		return &err
	}
	first := scope.Stack[len(scope.Stack)-1]
	_, ok := first.(AsInt);
	if !ok {
		err := Error{}
		err.message = "StackIndexError: `dec` expected one or more <int> type element in the stack."
		err.Type = StackIndexError
		return &err
	}
	scope.OpDrop()
	val := first.(AsInt).IntValue - 1
	expr := AsInt {
		val,
	}
	scope.OpPush(expr, nil)
	return nil
}

func (scope *Scope) OpDup() (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = "StackIndexError: `dup` expected one or more <int> type element in the stack."
		err.Type = StackIndexError
		return &err
	}
	first := scope.Stack[len(scope.Stack)-1]
	scope.OpPush(first, nil)
	return nil
}

func (scope *Scope) OpAppend() (*Error) {
	if len(scope.Stack) < 2 {
		err := Error{}
		err.message = "StackIndexError: `append` expected two or more element in the stack."
		err.Type = StackIndexError
		return &err
	}
	if _, ok := scope.Stack[len(scope.Stack)-2].(AsList); !ok {
		err := Error{}
		err.message = "TypeError: `append` expected <list> type element in the stack."
		err.Type = TypeError
		return &err
	}
	a := append(scope.Stack[len(scope.Stack)-2].(AsList).ListArgs, scope.Stack[len(scope.Stack)-1])
	scope.OpDrop()
	scope.OpDrop()
	var NewList AST = AsList {
		a,
	}
	scope.OpPush(NewList, nil)
	return nil
}

func (scope *Scope) OpRead() (*Error) {
	if len(scope.Stack) < 2 {
		err := Error{}
		err.message = "StackIndexError: `read` expected two or more element in the stack."
		err.Type = StackIndexError
		return &err
	}
	visitedList := scope.Stack[len(scope.Stack)-2]
	visitedIndex := scope.Stack[len(scope.Stack)-1]
	if _, ok := visitedIndex.(AsInt); !ok {
		err := Error{}
		err.message = "TypeError: `read` index expected <int> type element in the stack."
		err.Type = TypeError
		return &err
	}
	_, ok := visitedList.(AsStr);
	_, ok2 := visitedList.(AsList);
	if !ok && !ok2 {
		err := Error{}
		err.message = "TypeError: `read` expected <list> or <string> type element in the stack."
		err.Type = TypeError
		return &err
	}
	scope.OpDrop()
	scope.OpDrop()
	if _, ok := visitedList.(AsList); ok {
		if len(visitedList.(AsList).ListArgs) <= visitedIndex.(AsInt).IntValue {
			err := Error{}
			err.message = "IndexError: `read` type <list> element index out of range."
			err.Type = IndexError
			return &err
		}
		scope.OpPush(visitedList.(AsList).ListArgs[int(visitedIndex.(AsInt).IntValue)], nil)
	} else {
		if len(visitedList.(AsStr).StringValue) <= visitedIndex.(AsInt).IntValue {
			err := Error{}
			err.message = "IndexError: `read` type <string> element index out of range."
			err.Type = IndexError
			return &err
		}
		StringValue := string([]rune(visitedList.(AsStr).StringValue)[int(visitedIndex.(AsInt).IntValue)])
		var StrExpr AST = AsStr {
			StringValue,
		}
		scope.OpPush(StrExpr, nil)
	}
	return nil
}

func (scope *Scope) OpReplace() (*Error) {
	if len(scope.Stack) < 3 {
		err := Error{}
		err.message = "StackIndexError: `replace` expected three or more element in the stack."
		err.Type = StackIndexError
		return &err
	}
	visitedList := scope.Stack[len(scope.Stack)-3]
	visitedValue := scope.Stack[len(scope.Stack)-2]
	visitedIndex := scope.Stack[len(scope.Stack)-1]
	if _, ok := visitedIndex.(AsInt); !ok {
		err := Error{}
		err.message = "TypeError: `replace` index expected <int> type element in the stack."
		err.Type = TypeError
		return &err
	}
	if _, ok := visitedList.(AsList); !ok {
		err := Error{}
		err.message = "TypeError: `replace` expected <list> type element in the stack."
		err.Type = TypeError
		return &err
	}
	if len(visitedList.(AsList).ListArgs) <= visitedIndex.(AsInt).IntValue {
		err := Error{}
		err.message = "IndexError: `replace` type <list> element index out of range."
		err.Type = IndexError
		return &err
	}
	visitedList.(AsList).ListArgs[int(visitedIndex.(AsInt).IntValue)] = visitedValue
	scope.OpDrop()
	scope.OpDrop()
	scope.OpDrop()
	scope.OpPush(visitedList, nil)
	return nil
}

func (scope *Scope) OpRemove() (*Error) {
	if len(scope.Stack) < 2 {
		err := Error{}
		err.message = "StackIndexError: `remove` expected two or more element in the stack."
		err.Type = StackIndexError
		return &err
	}
	visitedList := scope.Stack[len(scope.Stack)-2]
	visitedIndex := scope.Stack[len(scope.Stack)-1]
	if _, ok := visitedIndex.(AsInt); !ok {
		err := Error{}
		err.message = "TypeError: `remove` index expected <int> type element in the stack."
		err.Type = TypeError
		return &err
	}
	if _, ok := visitedList.(AsList); !ok {
		err := Error{}
		err.message = "TypeError: `remove` expected <list> type element in the stack."
		err.Type = TypeError
		return &err
	}
	if len(visitedList.(AsList).ListArgs) <= visitedIndex.(AsInt).IntValue {
		err := Error{}
		err.message = "IndexError: `remove` type <list> element index out of range."
		err.Type = IndexError
		return &err
	}

	NewList := append(visitedList.(AsList).ListArgs[:int(visitedIndex.(AsInt).IntValue)], visitedList.(AsList).ListArgs[int(visitedIndex.(AsInt).IntValue)+1:]...)
    var ListExpr AST = AsList {
		NewList,
	}
	visitedList = nil
	scope.OpDrop()
	scope.OpDrop()
	scope.OpPush(ListExpr, nil)
	return nil
}

func (scope *Scope) OpIn() (*Error) {
	if len(scope.Stack) < 2 {
		err := Error{}
		err.message = "StackIndexError: `in` expected two or more element in the stack."
		err.Type = StackIndexError
		return &err
	}
	visitedVal := scope.Stack[len(scope.Stack)-2]
	visitedList := scope.Stack[len(scope.Stack)-1]
	scope.OpDrop()
	scope.OpDrop()
	if _, ok := visitedList.(AsList); !ok {
		err := Error{}
		err.message = "TypeError: `in` expected <list> type element in the stack."
		err.Type = TypeError
		return &err
	}
	for i := 0; i < len(visitedList.(AsList).ListArgs); i++ {
		var val AST
		switch visitedVal.(type) {
			case AsStr: val = visitedVal.(AsStr)
			case AsInt: val = visitedVal.(AsInt)
			case AsList: val = visitedVal.(AsList)
		}
		if visitedList.(AsList).ListArgs[i] == val {
			expr := AsBool {
				true,
			}
			scope.OpPush(expr, nil)
			return nil
		}
	}
	expr := AsBool {
		false,
	}
	scope.OpPush(expr, nil)
	return nil
}

func (scope *Scope) OpLen() (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = "StackIndexError: `len` expected one or more element in the stack."
		err.Type = StackIndexError
		return &err
	}
	visitedExpr := scope.Stack[len(scope.Stack)-1]
	_, ok := visitedExpr.(AsList);
	_, ok2 := visitedExpr.(AsStr);
	if !ok && !ok2 {
		err := Error{}
		err.message = "TypeError: `len` expected <list> or <string> type element in the stack."
		err.Type = TypeError
		return &err
	}
	IntExpr := AsInt {}
	if ok {
		IntExpr.IntValue = len(visitedExpr.(AsList).ListArgs)
	} else {
		IntExpr.IntValue = len(visitedExpr.(AsStr).StringValue)
	}
	scope.OpDrop()
	scope.OpPush(IntExpr, nil)
	return nil
}

func (scope *Scope) OpTypeOf() (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = "StackIndexError: `typeof` expected one or more element in the stack."
		err.Type = StackIndexError
		return &err
	}
	visitedVal := scope.Stack[len(scope.Stack)-1]
	var TypeVal string
	switch visitedVal.(type) {
		case AsStr: TypeVal = "string"
		case AsInt: TypeVal = "int"
		case AsList: TypeVal = "list"
		case AsBool: TypeVal = "bool"
		case AsType: TypeVal = "type"
		case AsError: TypeVal = "error"
	}
	expr := AsType {
		TypeVal,
	}
	scope.OpDrop()
	scope.OpPush(expr, nil)
	return nil
}

func (scope *Scope) OpRot() (*Error) {
	if len(scope.Stack) < 3 {
		err := Error{}
		err.message = "StackIndexError: `rot` expected more than three elements in the stack."
		err.Type = StackIndexError
		return &err
	}
	visitedExpr := scope.Stack[len(scope.Stack)-1]
	visitedExprSecond := scope.Stack[len(scope.Stack)-2]
	visitedExprThird := scope.Stack[len(scope.Stack)-3]
	scope.OpDrop()
	scope.OpDrop()
	scope.OpDrop()
	scope.OpPush(visitedExprSecond, nil)
	scope.OpPush(visitedExpr, nil)
	scope.OpPush(visitedExprThird, nil)
	return nil
}

func (scope *Scope) OpOver() (*Error) {
	if len(scope.Stack) < 2 {
		err := Error{}
		err.message = "StackIndexError: `over` expected more than two elements in the stack."
		err.Type = StackIndexError
		return &err
	}
	scope.OpPush(scope.Stack[len(scope.Stack)-2], nil)
	return nil
}

func (scope *Scope) OpVardef(name string, VariableScope *map[string]AST) (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = fmt.Sprintf("StackIndexError: variable `%s` definiton expected one or more element in the stack.", name)
		err.Type = StackIndexError
		return &err
	}
	if VariableScope == nil {
		VarValue := scope.Stack[len(scope.Stack)-1]
		Variables[name] = VarValue
		scope.OpDrop()
	} else {
		if _, ok := Variables[name]; ok {
			VarValue := scope.Stack[len(scope.Stack)-1]
			Variables[name] = VarValue
			scope.OpDrop()
		} else {
			VariablesScope := *VariableScope
			VarValue := scope.Stack[len(scope.Stack)-1]
			VariablesScope[name] = VarValue
			scope.OpDrop()
		}
	}
	return nil
}

func (scope *Scope) OpBlockdef(node AST) (*Error) {
	if _, ok := Blocks[node.(Blockdef).Name]; ok {
		err := Error{}
		err.message = fmt.Sprintf("NameError: block `%s` is already defined.", node.(Blockdef).Name)
		err.Type = NameError
		return &err
	}
	Blocks[node.(Blockdef).Name] = node.(Blockdef).BlockBody
	return nil
}

func (scope *Scope) OpCallBlock(name string) *Error {
	if _, ok := Blocks[name]; !ok {
		err := Error{}
		err.message = fmt.Sprintf("NameError: undefined block `%s`.", name)
		err.Type = NameError
		return &err
	}
	VariableScope := map[string]AST{}
	scope.VisitorVisit(Blocks[name], false, &VariableScope)
	return nil
}

func (scope *Scope) OpImport(FileName string) (*Error) {
	if _, err := os.Stat(FileName); os.IsNotExist(err) {
		err := Error{}
		err.message = fmt.Sprintf("ImportError: invalid file name `%s`.", FileName)
		err.Type = ImportError
		return &err
	}
	file, err := os.Open(FileName)
	if err != nil {
		panic(err)
	}
	lexer := LexerInit(file, FileName)
	parser := ParserInit(lexer)
	ast := ParserParse(parser)
	scope.VisitorVisit(ast, false, nil)
	return nil
}

func (scope *Scope) OpAssert(node AST) (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = "StackIndexError: `assert` expected one or more <bool> type element in the stack."
		err.Type = StackIndexError
		return &err
	}
	BoolValue := scope.Stack[len(scope.Stack)-1]
	if _, ok := BoolValue.(AsBool); !ok {
		err := Error{}
		err.message = "TypeError: `assert` expected <bool> type element in the stack."
		err.Type = TypeError
		return &err
	}
	scope.OpDrop()
	if !BoolValue.(AsBool).BoolValue {
		err := Error{}
		err.message = fmt.Sprintf("AssertionError:%d:%d: %s", node.(Assert).Line, node.(Assert).Col, node.(Assert).Message)
		err.Type = AssertionError
		return &err
	}
	return nil
}

func (scope *Scope) OpInput() {
	inputReader := bufio.NewReader(os.Stdin)
	input, _ := inputReader.ReadString('\n')
	input = input[:len(input)-1]
	
	StrExpr := AsStr {
		input,
	}
	scope.OpPush(StrExpr, nil)
}

// -----------------------------
// --------- Visitor -----------
// -----------------------------

func (scope *Scope) VisitorVisit(node AST, IsTry bool, VariableScope *map[string]AST) (bool, *Error, *map[string]AST) {
	BreakValue := false
	var err *Error
	for i := 0; i < len(node.(AsStatements)); i++ {
		node := node.(AsStatements)[i]
		switch node.(type) {
			case AsPush:
				err = scope.OpPush(node.(AsPush).value, VariableScope)
			case AsId:
				switch node.(AsId).name {
					case "print": err = scope.OpPrint()
					case "puts": err = scope.OpPuts()
					case "printS": scope.OpPrintS()
					case "break": BreakValue = true
					case "drop": err = scope.OpDrop()
					case "swap": err = scope.OpSwap()
					case "inc": err = scope.OpInc()
					case "dec": err = scope.OpDec()
					case "dup": err = scope.OpDup()
					case "append": err = scope.OpAppend()
					case "read": err = scope.OpRead()
					case "replace": err = scope.OpReplace()
					case "remove": err = scope.OpRemove()
					case "in": err = scope.OpIn()
					case "len": err = scope.OpLen()
					case "typeof": err = scope.OpTypeOf()
					case "rot": err = scope.OpRot()
					case "over": err = scope.OpOver()
					case "exit": os.Exit(0)
					case "input": scope.OpInput()
					default: panic("unreachable")
				}
			case AsBinop:
				err = scope.OpBinop(node.(AsBinop).op)
			case Vardef:
				err = scope.OpVardef(node.(Vardef).Name, VariableScope)
			case Blockdef:
				err = scope.OpBlockdef(node)
			case CallBlock:
				err = scope.OpCallBlock(node.(CallBlock).Name)
			case Import:
				err = scope.OpImport(node.(Import).FileName)
			case Compare:
				err = scope.OpCompare(node.(Compare).op)
			case AsStatements:
				scope.VisitorVisit(node.(AsStatements), IsTry, VariableScope)
			case If:
				BreakValue, err = scope.OpIf(node.(If), IsTry, VariableScope)
			case For:
				err = scope.OpFor(node.(For), IsTry, VariableScope)
			case Try:
				err = scope.OpTry(node.(Try), VariableScope)
			case Assert:
				err = scope.OpAssert(node)
			default:
				fmt.Println("Error: unexpected node type.")
		}
		if err != nil {
			if !IsTry {
				fmt.Println(err.message)
				os.Exit(0)
			}
			return BreakValue, err, VariableScope
		}
		if BreakValue {
			break
		}
	}
	return BreakValue, err, VariableScope
}


// -----------------------------
// ----------- Main ------------
// -----------------------------

func Usage() {
	fmt.Println("Usage:")
	fmt.Println("  tsh <filename>.tsp")
	os.Exit(0)
}

func main() {
	if len(os.Args) != 2 || os.Args[1] == "help" {
		Usage()
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(fmt.Sprintf("Error: invalid file name `%s`.", os.Args[1]))

		whilte := color.New(color.FgWhite)

		fmt.Print("Run ")
		boldWhite := whilte.Add(color.BgCyan)
		boldWhite.Print(" tsh help ")
		fmt.Println(" for usage")

		os.Exit(0)
	}

	lexer := LexerInit(file, os.Args[1])
	parser := ParserInit(lexer)
	ast := ParserParse(parser)
	scope := InitScope()
	scope.VisitorVisit(ast, false, nil)
}

