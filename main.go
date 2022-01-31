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
	"time"
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
					lexer.pos.column++
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
					lexer.pos.column++
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
					lexer.pos.column++
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
					lexer.pos.column++
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
					lexer.pos.column++
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
					} else if val == "NameError" || val == "StackUnderflowError" || val == "TypeError" || val == "ImportError" || val == "IndexError" || val == "AssertionError" {
						return startPos, TOKEN_ERROR, val, lexer.FileName
					} else if val == "except" {
						return startPos, TOKEN_EXCEPT, val, lexer.FileName
					}
					return startPos, TOKEN_ID, val, lexer.FileName
				} else if r == '"' {
					startPos := lexer.pos
					lexer.backup()
					lexer.pos.column++
					lexer.pos.column++
					val, _ := strconv.Unquote(lexer.lexString())
					r, _, err = lexer.reader.ReadRune()
					return startPos, TOKEN_STRING, val, lexer.FileName
				} else if string(r) == "'" {
					startPos := lexer.pos
					lexer.backup()
					lexer.pos.column++
					lexer.pos.column++
					val, _ := strconv.Unquote(lexer.lexStringSingle())
					r, _, err = lexer.reader.ReadRune()
					return startPos, TOKEN_STRING, val, lexer.FileName
				} else {
					file, err := os.Open(lexer.FileName)
					if err != nil {
						panic(err)
					}
					lexer := LexerInit(file, lexer.FileName)
					fmt.Println(fmt.Sprintf("./%s:SyntaxError:%d:%d: unexpected token value `%s`.", lexer.FileName, lexer.pos.line, lexer.pos.column, string(r)))
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
	val = val + string(r)
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
			val = val + string(r)
			lexer.backup()
			return val
		}
	}
}

func (lexer *Lexer) lexStringSingle() string {
	var val string
	r, _, err := lexer.reader.ReadRune()
	val = val + string("\"")
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
			val = val + "\""
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
	StackUnderflowError
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

type NodePosition struct {
	FileName string
	Line int
	Column int
}

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
	Position NodePosition
}

func (node AsId) node() {}

type Import struct {
	FileName string
	Position NodePosition
}

func (node Import) node() {}

type Assert struct {
	Position NodePosition
	Message string
}

func (node Assert) node() {}

type Compare struct {
	op uint8
	Position NodePosition
}

func (node Compare) node() {}

type AsError struct {
	err ErrorType
}

func (node AsError) node() {}

type AsBinop struct {
	op uint8
	Position NodePosition
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
	Position NodePosition
}

func (node Vardef) node() {}

type Var struct {
	Name string
	Position NodePosition
}

func (node Var) node() {}

type Blockdef struct {
	Name string
	BlockBody AST
}

func (node Blockdef) node() {}

type CallBlock struct {
	Name string
	Position NodePosition
}

func (node CallBlock) node() {}

type If struct {
	IfOp AST
	Position NodePosition
	IfBody AST
	ElifOps []AST
	ElifPositions []NodePosition
	ElifBodys []AST
	ElseBody AST
}

func (node If) node() {}

type For struct {
	ForOp AST
	Position NodePosition
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
		fmt.Println(fmt.Sprintf("./%s:SyntaxError:%d:%d: unexpected token value `%s`.", parser.FileName, parser.line, parser.column, parser.current_token_value))
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

func RetNodePosition(parser *Parser) NodePosition {
	return NodePosition{
		Line: parser.line,
		Column: parser.column,
		FileName: parser.FileName,
	}
}

func ParserParseError(parser *Parser) AST {
	var err ErrorType
	if parser.current_token_value == "StackUnderflowError" {
		err = StackUnderflowError
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
				Name: parser.current_token_value,
				Position: RetNodePosition(parser),
			}
			parser.ParserEat(TOKEN_ID)
		case TOKEN_TYPE:
			expr = AsType {
				parser.current_token_value,
			}
			parser.ParserEat(TOKEN_TYPE)
		default:
			fmt.Println(fmt.Sprintf("./%s:SyntaxError:%d:%d: unexpected token value `%s`.", parser.FileName, parser.line, parser.column, parser.current_token_value))
			os.Exit(0)
	}

	return expr
}

func ParserParse(parser *Parser) AST {
	var Statements AsStatements
	if  parser.current_token_type == TOKEN_DO || parser.current_token_type == TOKEN_END || parser.current_token_type == TOKEN_ELIF || parser.current_token_type == TOKEN_ELSE || parser.current_token_type == TOKEN_EXCEPT {
		fmt.Println(fmt.Sprintf("./%s:SyntaxError:%d:%d: the body is empty, unexpected token value `%s`.", parser.FileName, parser.line, parser.column, parser.current_token_value))
		os.Exit(0)
	}
	for {
		if parser.current_token_type == TOKEN_ID {
			if parser.current_token_value == "print" || parser.current_token_value == "break" || parser.current_token_value == "append" || parser.current_token_value == "remove" || parser.current_token_value == "swap" || parser.current_token_value == "in" || parser.current_token_value == "typeof" || parser.current_token_value == "rot" || parser.current_token_value == "len" || parser.current_token_value == "input" || parser.current_token_value == "sleep" ||
			   parser.current_token_value == "drop"  || parser.current_token_value == "dup" || parser.current_token_value == "inc" || parser.current_token_value == "dec" || parser.current_token_value == "replace" || parser.current_token_value == "read" || parser.current_token_value == "println" || parser.current_token_value == "over" || parser.current_token_value == "printS" || parser.current_token_value == "exit" {
				name := parser.current_token_value
				position := RetNodePosition(parser)
				IdExpr := AsId{
					name,
					position,
				}
				parser.ParserEat(TOKEN_ID)
				Statements = append(Statements, IdExpr)
			} else if parser.current_token_value == "assert" {
				position := RetNodePosition(parser)
				parser.ParserEat(TOKEN_ID)
				AssertExpr := Assert {
					Position: position,
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
				position := RetNodePosition(parser)
				parser.ParserEat(TOKEN_ID)
				name := parser.current_token_value
				parser.ParserEat(TOKEN_ID)
				CallBlockExpr := CallBlock {
					name,
					position,
				}
				Statements = append(Statements, CallBlockExpr)
			} else if parser.current_token_value == "import" {
				parser.ParserEat(TOKEN_ID)
				ImportExpr := Import {
					parser.current_token_value,
					RetNodePosition(parser),
				}
				parser.ParserEat(TOKEN_STRING)
				Statements = append(Statements, ImportExpr)
			} else if parser.current_token_value == "if" {
				parser.ParserEat(TOKEN_ID)
				position := RetNodePosition(parser)
				IfOp := ParserParse(parser)
				parser.ParserEat(TOKEN_DO)
				IfBody := ParserParse(parser)
				var ElifOps []AST
				var ElifBodys []AST
				var ElifPositions []NodePosition
				for {
					if parser.current_token_type != TOKEN_ELIF {
						break
					}
					parser.ParserEat(TOKEN_ELIF)
					ElifPosition := RetNodePosition(parser);
					ElifOp := ParserParse(parser)
					ElifOps = append(ElifOps, ElifOp)
					ElifPositions = append(ElifPositions, ElifPosition)
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
					Position: position,
					IfBody: IfBody,
					ElifOps: ElifOps,
					ElifPositions: ElifPositions,
					ElifBodys: ElifBodys,
					ElseBody: ElseBody,
				}
				Statements = append(Statements, IfExpr)
			} else if parser.current_token_value == "for" {
				parser.ParserEat(TOKEN_ID)
				position := RetNodePosition(parser);
				ForOp := ParserParse(parser)
				parser.ParserEat(TOKEN_DO)
				ForBody := ParserParse(parser)
				parser.ParserEat(TOKEN_END)
				ForExpr := For {
					ForOp: ForOp,
					Position: position,
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
			position := RetNodePosition(parser)
			parser.ParserEat(TOKEN_EQUALS)
			VardefExpr := Vardef {
				Name: parser.current_token_value,
				Position: position,
			}
			parser.ParserEat(TOKEN_ID)
			Statements = append(Statements, VardefExpr)
		} else if parser.current_token_type == TOKEN_PLUS || parser.current_token_type == TOKEN_MINUS || parser.current_token_type == TOKEN_MUL || parser.current_token_type == TOKEN_DIV || parser.current_token_type == TOKEN_REM {
			BinopExpr := AsBinop {
				op: uint8(parser.current_token_type),
				Position: RetNodePosition(parser),
			}
			parser.ParserEat(parser.current_token_type)
			Statements = append(Statements, BinopExpr)
		} else if parser.current_token_type == TOKEN_LESS_EQUALS || parser.current_token_type == TOKEN_GREATER_EQUALS || parser.current_token_type == TOKEN_LESS_THAN || parser.current_token_type == TOKEN_GREATER_THAN || parser.current_token_type == TOKEN_IS_EQUALS || parser.current_token_type == TOKEN_NOT_EQUALS || parser.current_token_type == TOKEN_OR || parser.current_token_type == TOKEN_AND {
			CompareExpr := Compare {
				op: uint8(parser.current_token_type),
				Position: RetNodePosition(parser),
			}
			parser.ParserEat(parser.current_token_type)
			Statements = append(Statements, CompareExpr)
		} else if parser.current_token_type == TOKEN_EOF || parser.current_token_type == TOKEN_DO ||
		    parser.current_token_type == TOKEN_END || parser.current_token_type == TOKEN_ELIF ||
			parser.current_token_type == TOKEN_ELSE || parser.current_token_type == TOKEN_EXCEPT || parser.current_token_type == TOKEN_R_BRACKET {
			break
		} else {
			fmt.Println(fmt.Sprintf("./%s:SyntaxError:%d:%d: unexpected token value `%s`.", parser.FileName, parser.line, parser.column, parser.current_token_value))
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
		scope.Stack = append(scope.Stack, AsList{ListScope.Stack})
	} else if IsVar {
		if VariableScope == nil {
			if _, ok := Variables[node.(Var).Name]; ok {
				scope.Stack = append(scope.Stack, Variables[node.(Var).Name])
			} else {
				err := Error{}
				err.message = fmt.Sprintf("./%s:NameError:%d:%d: undefined variable `%s`", node.(Var).Position.FileName, node.(Var).Position.Line, node.(Var).Position.Column,node.(Var).Name)
				err.Type = NameError
				return &err
			}
		} else {
			if _, ok := (*VariableScope)[node.(Var).Name]; ok {
				scope.Stack = append(scope.Stack, (*VariableScope)[node.(Var).Name])
			} else {
				if _, ok := Variables[node.(Var).Name]; ok {
					scope.Stack = append(scope.Stack, Variables[node.(Var).Name])
				} else {
					err := Error{}
					err.message = fmt.Sprintf("./%s:NameError:%d:%d: undefined variable `%s`", node.(Var).Position.FileName, node.(Var).Position.Line, node.(Var).Position.Column,node.(Var).Name)
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

func (scope *Scope) OpDrop(node AST) (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `drop` expected one or more element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
		return &err
	}
	scope.Stack = scope.Stack[:len(scope.Stack)-1]
	return nil
}

func (scope *Scope) OpSwap(node AST) (*Error) {
	if len(scope.Stack) < 2 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `swap` expected two or more element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
		return &err
	}
	first := scope.Stack[len(scope.Stack)-1]
	second := scope.Stack[len(scope.Stack)-2]
	scope.Stack = scope.Stack[:len(scope.Stack)-2]
	scope.OpPush(first, nil)
	scope.OpPush(second, nil)
	return nil
}

func RetTokenAsStr(token uint8) string {
	return tokens[token]
}

func (scope *Scope) OpBinop(op uint8, position NodePosition) (*Error) {
	if len(scope.Stack) < 2 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `%s` expected more than 2 <int> type elements in the stack.", position.FileName, position.Line, position.Column, RetTokenAsStr(op))
		err.Type = StackUnderflowError
		return &err
	}
	first := scope.Stack[len(scope.Stack)-1]
	second := scope.Stack[len(scope.Stack)-2]
	scope.Stack = scope.Stack[:len(scope.Stack)-2]
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
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `%s` expected 2 <int> type or 2 <string> type elements in the stack.", position.FileName, position.Line, position.Column, RetTokenAsStr(op))
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

func (scope *Scope) OpCompare(op uint8, position NodePosition) (*Error) {
	if len(scope.Stack) < 2 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `%s` expected more than 2 elements in the stack.", position.FileName, position.Line, position.Column, RetTokenAsStr(op))
		err.Type = StackUnderflowError
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
			err.message = fmt.Sprintf("./%s:TypeError:%d:%d: `%s` expected 2 <bool> type elements in the stack.", position.FileName, position.Line, position.Column, RetTokenAsStr(op))
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
			err.message = fmt.Sprintf("./%s:TypeError:%d:%d: `%s` expected 2 <int> type elements in the stack.", position.FileName, position.Line, position.Column, RetTokenAsStr(op))
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
	scope.Stack = scope.Stack[:len(scope.Stack)-2]
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

func (scope *Scope) OpPrintln(node AST) (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `println` the stack is empty.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
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
				case StackUnderflowError: fmt.Println("<error 'StackUnderflowError'>")
				case ImportError: fmt.Println("<error 'ImportError'>")
				case IndexError: fmt.Println("<error 'IndexError'>")
				case TypeError: fmt.Println("<error 'TypeError'>")
				default: fmt.Println(fmt.Sprintf("unexpected error <%d>", expr.(AsError).err))
			}
		case AsList:
			PrintAsList(expr)
			fmt.Println()
	}
	scope.Stack = scope.Stack[:len(scope.Stack)-1]
	return nil
}

func (scope *Scope) OpPrint(node AST) (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `print` the stack is empty.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
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
				case StackUnderflowError: fmt.Print("<error 'StackUnderflowError'>")
				case ImportError: fmt.Print("<error 'ImportError'>")
				case IndexError: fmt.Print("<error 'IndexError'>")
				case TypeError: fmt.Print("<error 'TypeError'>")
				default: fmt.Print(fmt.Sprintf("unexpected error <%d>", expr.(AsError).err))
			}
		case AsList:
			PrintAsList(expr)
	}
	scope.Stack = scope.Stack[:len(scope.Stack)-1]
	return nil
}

func (scope *Scope) OpIf(node AST, IsTry bool, VariableScope *map[string]AST) (bool, *Error) {
	scope.VisitorVisit(node.(If).IfOp, IsTry, VariableScope)
	var BreakValue bool = false
	var err *Error = nil
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: if statement expected one or more <bool> type element in the stack.", node.(If).Position.FileName, node.(If).Position.Line, node.(If).Position.Column)
		err.Type = StackUnderflowError
		return BreakValue, &err
	}
	expr := scope.Stack[len(scope.Stack)-1]
	if _, ok := expr.(AsBool); !ok {
		err := Error{}
		err.message = fmt.Sprintf("./%s:TypeError:%d:%d: if statement expected one or more <bool> type element in the stack.", node.(If).Position.FileName, node.(If).Position.Line, node.(If).Position.Column)
		err.Type = StackUnderflowError
		return BreakValue, &err
	}
	scope.Stack = scope.Stack[:len(scope.Stack)-1]
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
			err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: if statement expected one or more <bool> type element in the stack.", node.(If).ElifPositions[i].FileName, node.(If).ElifPositions[i].Line, node.(If).ElifPositions[i].Column)
			err.Type = StackUnderflowError
			return BreakValue, &err
		}
		expr := scope.Stack[len(scope.Stack)-1]
		if _, ok := expr.(AsBool); !ok {
			err := Error{}
			err.message = fmt.Sprintf("./%s:TypeError:%d:%d: if statement expected one or more <bool> type element in the stack.", node.(If).ElifPositions[i].FileName, node.(If).ElifPositions[i].Line, node.(If).ElifPositions[i].Column)
			err.Type = StackUnderflowError
			return BreakValue, &err
		}
		scope.Stack = scope.Stack[:len(scope.Stack)-1]
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
			err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: for loop expected one or more <bool> type element in the stack.", node.(For).Position.FileName, node.(For).Position.Line, node.(For).Position.Column)
			err.Type = StackUnderflowError
			return &err
		}
		expr := scope.Stack[len(scope.Stack)-1]
		if _, ok := expr.(AsBool); !ok {
			err := Error{}
			err.message = fmt.Sprintf("./%s:TypeError:%d:%d: for loop expected one or more <bool> type element in the stack.", node.(For).Position.FileName, node.(For).Position.Line, node.(For).Position.Column)
			err.Type = TypeError
			return &err
		}
		scope.Stack = scope.Stack[:len(scope.Stack)-1]
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
	_, err, _ := scope.VisitorVisit(node.(Try).TryBody, true, VariableScope)
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

func (scope *Scope) OpInc(node AST) (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `inc` expected one or more <int> type element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
		return &err
	}
	first := scope.Stack[len(scope.Stack)-1]
	_, ok := first.(AsInt);
	if !ok {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `inc` expected one or more <int> type element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
		return &err
	}
	scope.Stack = scope.Stack[:len(scope.Stack)-1]
	val := first.(AsInt).IntValue + 1
	expr := AsInt {
		val,
	}
	scope.OpPush(expr, nil)
	return nil
}

func (scope *Scope) OpDec(node AST) (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `dec` expected one or more <int> type element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
		return &err
	}
	first := scope.Stack[len(scope.Stack)-1]
	_, ok := first.(AsInt);
	if !ok {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `dec` expected one or more <int> type element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
		return &err
	}
	scope.Stack = scope.Stack[:len(scope.Stack)-1]
	val := first.(AsInt).IntValue - 1
	expr := AsInt {
		val,
	}
	scope.OpPush(expr, nil)
	return nil
}

func (scope *Scope) OpDup(node AST) (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `dup` expected one or more <int> type element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
		return &err
	}
	first := scope.Stack[len(scope.Stack)-1]
	scope.OpPush(first, nil)
	return nil
}

func (scope *Scope) OpAppend(node AST) (*Error) {
	if len(scope.Stack) < 2 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `append` expected two or more element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
		return &err
	}
	if _, ok := scope.Stack[len(scope.Stack)-2].(AsList); !ok {
		err := Error{}
		err.message = fmt.Sprintf("./%s:TypeError:%d:%d: `append` expected <list> type element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = TypeError
		return &err
	}
	a := append(scope.Stack[len(scope.Stack)-2].(AsList).ListArgs, scope.Stack[len(scope.Stack)-1])
	scope.Stack = scope.Stack[:len(scope.Stack)-2]
	var NewList AST = AsList {
		a,
	}
	scope.OpPush(NewList, nil)
	return nil
}

func (scope *Scope) OpRead(node AST) (*Error) {
	if len(scope.Stack) < 2 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `read` expected two or more element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
		return &err
	}
	visitedList := scope.Stack[len(scope.Stack)-2]
	visitedIndex := scope.Stack[len(scope.Stack)-1]
	if _, ok := visitedIndex.(AsInt); !ok {
		err := Error{}
		err.message = fmt.Sprintf("./%s:TypeError:%d:%d: `read` index expected <int> type element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = TypeError
		return &err
	}
	_, ok := visitedList.(AsStr);
	_, ok2 := visitedList.(AsList);
	if !ok && !ok2 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:TypeError:%d:%d: `read` expected <list> or <string> type element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = TypeError
		return &err
	}
	scope.Stack = scope.Stack[:len(scope.Stack)-2]
	if _, ok := visitedList.(AsList); ok {
		if len(visitedList.(AsList).ListArgs) <= visitedIndex.(AsInt).IntValue {
			err := Error{}
			err.message = fmt.Sprintf("./%s:IndexError:%d:%d: `read` type <list> element index out of range.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
			err.Type = IndexError
			return &err
		}
		scope.OpPush(visitedList.(AsList).ListArgs[int(visitedIndex.(AsInt).IntValue)], nil)
	} else {
		if len(visitedList.(AsStr).StringValue) <= visitedIndex.(AsInt).IntValue {
			err := Error{}
			err.message = fmt.Sprintf("./%s:IndexError:%d:%d: `read` type <string> element index out of range.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
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

func (scope *Scope) OpReplace(node AST) (*Error) {
	if len(scope.Stack) < 3 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `replace` expected three or more element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
		return &err
	}
	visitedList := scope.Stack[len(scope.Stack)-3]
	visitedValue := scope.Stack[len(scope.Stack)-2]
	visitedIndex := scope.Stack[len(scope.Stack)-1]
	if _, ok := visitedIndex.(AsInt); !ok {
		err := Error{}
		err.message = fmt.Sprintf("./%s:TypeError:%d:%d: `replace` index expected <int> type element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = TypeError
		return &err
	}
	if _, ok := visitedList.(AsList); !ok {
		err := Error{}
		err.message = fmt.Sprintf("./%s:TypeError:%d:%d: `replace` expected <list> type element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = TypeError
		return &err
	}
	if len(visitedList.(AsList).ListArgs) <= visitedIndex.(AsInt).IntValue {
		err := Error{}
		err.message = fmt.Sprintf("./%s:IndexError:%d:%d: `replace` type <list> element index out of range.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = IndexError
		return &err
	}
	visitedList.(AsList).ListArgs[int(visitedIndex.(AsInt).IntValue)] = visitedValue
	scope.Stack = scope.Stack[:len(scope.Stack)-3]
	scope.OpPush(visitedList, nil)
	return nil
}

func (scope *Scope) OpRemove(node AST) (*Error) {
	if len(scope.Stack) < 2 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `remove` expected two or more element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
		return &err
	}
	visitedList := scope.Stack[len(scope.Stack)-2]
	visitedIndex := scope.Stack[len(scope.Stack)-1]
	if _, ok := visitedIndex.(AsInt); !ok {
		err := Error{}
		err.message = fmt.Sprintf("./%s:TypeError:%d:%d: `remove` index expected <int> type element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = TypeError
		return &err
	}
	if _, ok := visitedList.(AsList); !ok {
		err := Error{}
		err.message = fmt.Sprintf("./%s:TypeError:%d:%d: `remove` expected <list> type element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = TypeError
		return &err
	}
	if len(visitedList.(AsList).ListArgs) <= visitedIndex.(AsInt).IntValue {
		err := Error{}
		err.message = fmt.Sprintf("./%s:IndexError:%d:%d: `remove` type <list> element index out of range.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = IndexError
		return &err
	}

	NewList := append(visitedList.(AsList).ListArgs[:int(visitedIndex.(AsInt).IntValue)], visitedList.(AsList).ListArgs[int(visitedIndex.(AsInt).IntValue)+1:]...)
    var ListExpr AST = AsList {
		NewList,
	}
	visitedList = nil
	scope.Stack = scope.Stack[:len(scope.Stack)-2]
	scope.OpPush(ListExpr, nil)
	return nil
}

func (scope *Scope) OpIn(node AST) (*Error) {
	if len(scope.Stack) < 2 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `in` expected two or more element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
		return &err
	}
	visitedVal := scope.Stack[len(scope.Stack)-2]
	visitedList := scope.Stack[len(scope.Stack)-1]
	scope.Stack = scope.Stack[:len(scope.Stack)-2]
	if _, ok := visitedList.(AsList); !ok {
		err := Error{}
		err.message = fmt.Sprintf("./%s:TypeError:%d:%d: `in` expected <list> type element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
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

func (scope *Scope) OpLen(node AST) (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `len` expected one or more element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
		return &err
	}
	visitedExpr := scope.Stack[len(scope.Stack)-1]
	_, ok := visitedExpr.(AsList);
	_, ok2 := visitedExpr.(AsStr);
	if !ok && !ok2 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:TypeError:%d:%d: `len` expected <list> or <string> type element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = TypeError
		return &err
	}
	IntExpr := AsInt {}
	if ok {
		IntExpr.IntValue = len(visitedExpr.(AsList).ListArgs)
	} else {
		IntExpr.IntValue = len(visitedExpr.(AsStr).StringValue)
	}
	scope.Stack = scope.Stack[:len(scope.Stack)-1]
	scope.OpPush(IntExpr, nil)
	return nil
}

func (scope *Scope) OpTypeOf(node AST) (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `typeof` expected one or more element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
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
	scope.Stack = scope.Stack[:len(scope.Stack)-1]
	scope.OpPush(expr, nil)
	return nil
}

func (scope *Scope) OpRot(node AST) (*Error) {
	if len(scope.Stack) < 3 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `rot` expected more than three elements in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
		return &err
	}
	visitedExpr := scope.Stack[len(scope.Stack)-1]
	visitedExprSecond := scope.Stack[len(scope.Stack)-2]
	visitedExprThird := scope.Stack[len(scope.Stack)-3]
	scope.Stack = scope.Stack[:len(scope.Stack)-3]
	scope.OpPush(visitedExprSecond, nil)
	scope.OpPush(visitedExpr, nil)
	scope.OpPush(visitedExprThird, nil)
	return nil
}

func (scope *Scope) OpOver(node AST) (*Error) {
	if len(scope.Stack) < 2 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `over` expected more than two elements in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
		return &err
	}
	scope.OpPush(scope.Stack[len(scope.Stack)-2], nil)
	return nil
}

func (scope *Scope) OpVardef(name string, position NodePosition, VariableScope *map[string]AST) (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: variable `%s` definiton expected one or more element in the stack.", position.FileName, position.Line, position.Column, name)
		err.Type = StackUnderflowError
		return &err
	}
	if VariableScope == nil {
		VarValue := scope.Stack[len(scope.Stack)-1]
		Variables[name] = VarValue
		scope.Stack = scope.Stack[:len(scope.Stack)-1]
	} else {
		if _, ok := Variables[name]; ok {
			VarValue := scope.Stack[len(scope.Stack)-1]
			Variables[name] = VarValue
			scope.Stack = scope.Stack[:len(scope.Stack)-1]
		} else {
			VarValue := scope.Stack[len(scope.Stack)-1]
			(*VariableScope)[name] = VarValue
			scope.Stack = scope.Stack[:len(scope.Stack)-1]
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

func (scope *Scope) OpCallBlock(name string, position NodePosition) *Error {
	if _, ok := Blocks[name]; !ok {
		err := Error{}
		err.message = fmt.Sprintf("./%s:NameError:%d:%d: undefined block `%s`.", position.FileName, position.Line, position.Column, name)
		err.Type = NameError
		return &err
	}
	VariableScope := map[string]AST{}
	scope.VisitorVisit(Blocks[name], false, &VariableScope)
	return nil
}

func (scope *Scope) OpImport(FileName string, position NodePosition) (*Error) {
	if _, err := os.Stat(FileName); os.IsNotExist(err) {
		err := Error{}
		err.message = fmt.Sprintf("./%s:ImportError:%d:%d: invalid file name `%s`.", position.FileName, position.Line, position.Column, FileName)
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
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `assert` expected one or more <bool> type element in the stack.", node.(Assert).Position.FileName, node.(Assert).Position.Line, node.(Assert).Position.Column)
		err.Type = StackUnderflowError
		return &err
	}
	BoolValue := scope.Stack[len(scope.Stack)-1]
	if _, ok := BoolValue.(AsBool); !ok {
		err := Error{}
		err.message = fmt.Sprintf("./%s:TypeError:%d:%d: `assert` expected <bool> type element in the stack.", node.(Assert).Position.FileName, node.(Assert).Position.Line, node.(Assert).Position.Column)
		err.Type = TypeError
		return &err
	}
	scope.Stack = scope.Stack[:len(scope.Stack)-1]
	if !BoolValue.(AsBool).BoolValue {
		err := Error{}
		err.message = fmt.Sprintf("./%s:AssertionError:%d:%d: %s", node.(Assert).Position.FileName, node.(Assert).Position.Line, node.(Assert).Position.Column, node.(Assert).Message)
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

func (scope *Scope) OpSleep(node AST) (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = fmt.Sprintf("./%s:StackUnderflowError:%d:%d: `sleep` expected one or more element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = StackUnderflowError
		return &err
	}
	IntValue := scope.Stack[len(scope.Stack)-1]
	if _, ok := IntValue.(AsInt); !ok {
		err := Error{}
		err.message = fmt.Sprintf("./%s:TypeError:%d:%d: `sleep` expected <int> type element in the stack.", node.(AsId).Position.FileName, node.(AsId).Position.Line, node.(AsId).Position.Column)
		err.Type = TypeError
		return &err
	}
	scope.Stack = scope.Stack[:len(scope.Stack)-1]
	time.Sleep(time.Duration(IntValue.(AsInt).IntValue) * time.Millisecond)
	return nil
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
					case "println": err = scope.OpPrintln(node)
					case "print": err = scope.OpPrint(node)
					case "break": BreakValue = true
					case "drop": err = scope.OpDrop(node)
					case "swap": err = scope.OpSwap(node)
					case "inc": err = scope.OpInc(node)
					case "dec": err = scope.OpDec(node)
					case "dup": err = scope.OpDup(node)
					case "append": err = scope.OpAppend(node)
					case "read": err = scope.OpRead(node)
					case "replace": err = scope.OpReplace(node)
					case "remove": err = scope.OpRemove(node)
					case "in": err = scope.OpIn(node)
					case "len": err = scope.OpLen(node)
					case "typeof": err = scope.OpTypeOf(node)
					case "rot": err = scope.OpRot(node)
					case "over": err = scope.OpOver(node)
					case "exit": os.Exit(0)
					case "input": scope.OpInput()
					case "sleep": scope.OpSleep(node)
					default: panic("unreachable")
				}
			case AsBinop:
				err = scope.OpBinop(node.(AsBinop).op, node.(AsBinop).Position)
			case Vardef:
				err = scope.OpVardef(node.(Vardef).Name, node.(Vardef).Position, VariableScope)
			case Blockdef:
				err = scope.OpBlockdef(node)
			case CallBlock:
				err = scope.OpCallBlock(node.(CallBlock).Name, node.(CallBlock).Position)
			case Import:
				err = scope.OpImport(node.(Import).FileName, node.(Import).Position)
			case Compare:
				err = scope.OpCompare(node.(Compare).op, node.(Compare).Position)
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
				panic("unreachable")
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

