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
					} else if val == "string" || val == "int" || val == "bool" || val == "type" || val == "list" || val == "error" {
						return startPos, TOKEN_TYPE, val
					} else if val == "else" {
						return startPos, TOKEN_ELSE, val
					} else if val == "elif" {
						return startPos, TOKEN_ELIF, val
					} else if val == "NameError" || val == "StackIndexError" || val == "TypeError" || val == "ImportError" {
						return startPos, TOKEN_ERROR, val
					} else if val == "except" {
						return startPos, TOKEN_EXCEPT, val
					}
					return startPos, TOKEN_ID, val
				} else if r == '"' {
					startPos := lexer.pos
					lexer.backup()
					val := lexer.lexString()
					r, _, err = lexer.reader.ReadRune()
					return startPos, TOKEN_STRING, val
				} else if string(r) == "'" {
					startPos := lexer.pos
					lexer.backup()
					val := lexer.lexStringSingle()
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
	ImportError
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

type AsId struct {
	name string
}

func (node AsId) node() {}

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
		fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value `%s`.", parser.line, parser.column, parser.current_token_value))
		os.Exit(0)
	}
	pos, tok, val := parser.lexer.Lex()
	parser.current_token_type = tok
	parser.current_token_value = val
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
		default:
			fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value `%s`.", parser.line, parser.column, parser.current_token_value))
			os.Exit(0)
	}

	return expr
}

func ParserParse(parser *Parser) AST {
	var Statements AsStatements
	if  parser.current_token_type == TOKEN_DO || parser.current_token_type == TOKEN_END || parser.current_token_type == TOKEN_ELIF || parser.current_token_type == TOKEN_ELSE || parser.current_token_type == TOKEN_EXCEPT {
		fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: the body is empty, unexpected token value `%s`.", parser.line, parser.column, parser.current_token_value))
		os.Exit(0)
	}
	for {
		if parser.current_token_type == TOKEN_ID {
			if parser.current_token_value == "print" || parser.current_token_value == "break" ||
			   parser.current_token_value == "drop"  || parser.current_token_value == "dup" || parser.current_token_value == "inc" || parser.current_token_value == "dec"  {
				name := parser.current_token_value
				IdExpr := AsId{name}
				parser.ParserEat(TOKEN_ID)
				Statements = append(Statements, IdExpr)
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
				fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value `%s`.", parser.line, parser.column, parser.current_token_value))
				os.Exit(0)
			}
		} else if parser.current_token_type == TOKEN_INT  || parser.current_token_type == TOKEN_STRING ||
		    parser.current_token_type == TOKEN_BOOL || parser.current_token_type == TOKEN_ERROR {
			expr := ParserParseExpr(parser)
			PushExpr := AsPush{
				value: expr,
			}
			Statements = append(Statements, PushExpr)
		} else if parser.current_token_type == TOKEN_PLUS || parser.current_token_type == TOKEN_MINUS || parser.current_token_type == TOKEN_MUL || parser.current_token_type == TOKEN_DIV {
			BinopExpr := AsBinop {
				op: uint8(parser.current_token_type),
			}
			parser.ParserEat(parser.current_token_type)
			Statements = append(Statements, BinopExpr)
		} else if parser.current_token_type == TOKEN_LESS_EQUALS || parser.current_token_type == TOKEN_GREATER_EQUALS || parser.current_token_type == TOKEN_LESS_THAN || parser.current_token_type == TOKEN_GREATER_THAN || parser.current_token_type == TOKEN_IS_EQUALS || parser.current_token_type == TOKEN_NOT_EQUALS {
			CompareExpr := Compare {
				op: uint8(parser.current_token_type),
			}
			parser.ParserEat(parser.current_token_type)
			Statements = append(Statements, CompareExpr)
		} else if parser.current_token_type == TOKEN_EOF || parser.current_token_type == TOKEN_DO ||
		    parser.current_token_type == TOKEN_END || parser.current_token_type == TOKEN_ELIF ||
			parser.current_token_type == TOKEN_ELSE || parser.current_token_type == TOKEN_EXCEPT {
			break
		} else {
			fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value `%s`.", parser.line, parser.column, parser.current_token_value))
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
    Variables map[string]AST
    Blocks map[string][]AST
}

func InitScope() *Scope {
	stack := []AST{}
	variables := map[string]AST{}
	blocks := map[string][]AST{}
	return &Scope{
		Stack: stack,
		Variables: variables,
		Blocks: blocks,
	}
}

func (scope *Scope) OpPush(node AST) {
	scope.Stack = append(scope.Stack, node)
}

func (scope *Scope) OpDrop() (*Error) {
	if len(scope.Stack) < 1 {
		err := Error{}
		err.message = "StackIndexError: `drop` expected more than one element in the stack."
		err.Type = StackIndexError
		return &err
	}
	scope.Stack = scope.Stack[:len(scope.Stack)-1]
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
	_, ok := first.(AsInt);
	_, ok2 := second.(AsInt);
	if !ok || !ok2 {
		err := Error{}
		err.message = fmt.Sprintf("StackIndexError: `%s` expected type <int>.", RetTokenAsStr(op))
		err.Type = TypeError
		return &err
	}
	var val int
	switch op {
		case TOKEN_PLUS: val = first.(AsInt).IntValue + second.(AsInt).IntValue
		case TOKEN_MINUS:  val = second.(AsInt).IntValue - first.(AsInt).IntValue
		case TOKEN_MUL: val = second.(AsInt).IntValue * first.(AsInt).IntValue
		case TOKEN_DIV: val = second.(AsInt).IntValue / first.(AsInt).IntValue
	}
	scope.OpDrop()
	scope.OpDrop()
	expr := AsInt {
		val,
	}
	scope.OpPush(expr)
	return nil
}

func (scope *Scope) OpCompare(op uint8) (*Error) {
	if len(scope.Stack) < 2 {
		err := Error{}
		err.message = fmt.Sprintf("StackIndexError: `%s` expected more than 2 <int> type elements in the stack.", RetTokenAsStr(op))
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
			}
		}
	} else {
		_, ok := first.(AsInt);
		_, ok2 := second.(AsInt);
		if !ok || !ok2 {
			err := Error{}
			err.message = fmt.Sprintf("StackIndexError: `%s` expected <int> type.", RetTokenAsStr(op))
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
	scope.OpPush(expr)
	return nil
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
		case AsError:
			switch expr.(AsError).err {
				case NameError: fmt.Println("<error 'NameError'>")
				case StackIndexError: fmt.Println("<error 'StackIndexError'>")
				case ImportError: fmt.Println("<error 'ImportError'>")
				default: fmt.Println(fmt.Sprintf("unexpected error <%d>", expr.(AsError).err))
			}
	}
	scope.OpDrop()
	return nil
}

func (scope *Scope) OpIf(node AST, IsTry bool) (bool, *Error) {
	scope.VisitorVisit(node.(If).IfOp, IsTry)
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
		BreakValue, err = scope.VisitorVisit(node.(If).IfBody, IsTry)
		if err != nil {
			return BreakValue, err
		}
		return BreakValue, nil
	}
	for i := 0; i < len(node.(If).ElifOps); i++ {
		scope.VisitorVisit(node.(If).ElifOps[i], IsTry)
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
			BreakValue, err = scope.VisitorVisit(node.(If).ElifBodys[i], IsTry)
			return BreakValue, err
		}
	}
	if node.(If).ElseBody != nil {
		BreakValue, err = scope.VisitorVisit(node.(If).ElseBody, IsTry)
	}
	return BreakValue, err
}

func (scope *Scope) OpFor(node AST, IsTry bool) (*Error) {
	var BreakValue bool
	for {
		_, err := scope.VisitorVisit(node.(For).ForOp, IsTry)
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
		BreakValue, err = scope.VisitorVisit(node.(For).ForBody, IsTry)
		if err != nil {
			return err
		}
		if BreakValue {
			return nil
		}
	}
	return nil
}

func (scope *Scope) OpTry(node AST) (*Error) {
	_, err := scope.VisitorVisit(node.(Try).TryBody, true)
	if err != nil {
		for i := 0; i < len(node.(Try).ExceptErrors); i++ {
			if node.(Try).ExceptErrors[i].(AsError).err == err.Type {
				scope.VisitorVisit(node.(Try).ExceptBodys[i], false)
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
	scope.OpPush(expr)
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
	scope.OpPush(expr)
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
	scope.OpPush(first)
	return nil
}


// -----------------------------
// --------- Visitor -----------
// -----------------------------

func (scope *Scope) VisitorVisit(node AST, IsTry bool) (bool, *Error) {
	BreakValue := false
	var err *Error
	for i := 0; i < len(node.(AsStatements)); i++ {
		node := node.(AsStatements)[i]
		switch node.(type) {
			case AsPush:
				scope.OpPush(node.(AsPush).value)
			case AsId:
				if node.(AsId).name == "print" {
					err = scope.OpPrint()
				} else if node.(AsId).name == "break" {
					BreakValue = true
				} else if node.(AsId).name == "drop" {
					err = scope.OpDrop()
				} else if node.(AsId).name == "inc" {
					err = scope.OpInc()
				} else if node.(AsId).name == "dec" {
					err = scope.OpDec()
				} else if node.(AsId).name == "dup" {
					err = scope.OpDup()
				}
			case AsBinop:
				err = scope.OpBinop(node.(AsBinop).op)
			case Compare:
				err = scope.OpCompare(node.(Compare).op)
			case AsStatements:
				scope.VisitorVisit(node.(AsStatements), IsTry)
			case If:
				BreakValue, err = scope.OpIf(node.(If), IsTry)
			case For:
				err = scope.OpFor(node.(For), IsTry)
			case Try:
				err = scope.OpTry(node.(Try))
			default:
				fmt.Println("Error: unexpected node type.")
		}
		if err != nil {
			if !IsTry {
				fmt.Println(err.message)
				os.Exit(0)
			}
			return BreakValue, err
		}
		if BreakValue {
			break
		}
	}
	return BreakValue, err
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
		fmt.Println(fmt.Sprintf("Error: invalid file name '%s'.", os.Args[1]))

		whilte := color.New(color.FgWhite)

		fmt.Print("Run ")
		boldWhite := whilte.Add(color.BgCyan)
		boldWhite.Print(" tsh help ")
		fmt.Println(" for usage")

		os.Exit(0)
	}

	lexer := LexerInit(file)
	parser := ParserInit(lexer)
	ast := ParserParse(parser)
	scope := InitScope()
	scope.VisitorVisit(ast, false)
}

