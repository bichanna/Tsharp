package main

import (
	"fmt"
	"bufio"
	"io"
	"os"
	"strconv"
	"unicode"
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
	TOKEN_EQUALS
	TOKEN_L_PAREN
	TOKEN_R_PAREN
	TOKEN_DO
	TOKEN_END
	TOKEN_BOOL
	TOKEN_COMMA
	TOKEN_MUL
	TOKEN_DIV
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_IS_EQUALS
	TOKEN_NONE
	TOKEN_ELSE
	TOKEN_LESS_THAN
	TOKEN_GREATER_THAN
	TOKEN_LESS_EQUALS
	TOKEN_GREATER_EQUALS
	TOKEN_NOT_EQUALS
)

var tokens = []string {
	TOKEN_EOF:            "TOKEN_EOF",
	TOKEN_ILLEGAL:        "TOKEN_ILLEGAL",
	TOKEN_ID:             "TOKEN_ID",
	TOKEN_STRING:         "TOKEN_STRING",
	TOKEN_INT:            "TOKEN_INT",
	TOKEN_EQUALS:         "TOKEN_EQUALS",
	TOKEN_L_PAREN:        "TOKEN_L_PAREN",
	TOKEN_R_PAREN:        "TOKEN_R_PAREN",
	TOKEN_DO:      		  "TOKEN_DO",
	TOKEN_END:            "TOKEN_END",
	TOKEN_BOOL:           "TOKEN_BOOL",
	TOKEN_COMMA:          "TOKEN_COMMA",
	TOKEN_MUL:            "TOKEN_MUL",
	TOKEN_DIV:            "TOKEN_DIV",
	TOKEN_PLUS:           "TOKEN_PLUS",
	TOKEN_MINUS:          "TOKEN_MINUS",
	TOKEN_IS_EQUALS:      "TOKEN_IS_EQUALS",
	TOKEN_NONE:           "TOKEN_NONE",
	TOKEN_ELSE:           "TOKEN_ELSE",
	TOKEN_LESS_THAN:      "TOKEN_LESS_THAN",
	TOKEN_GREATER_THAN:   "TOKEN_GREATER_THAN",
	TOKEN_LESS_EQUALS:    "TOKEN_LESS_EQUALS",
	TOKEN_GREATER_EQUALS: "TOKEN_GREATER_EQUALS",
	TOKEN_NOT_EQUALS:     "TOKEN_NOT_EQUALS",
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
			case '(': return lexer.pos, TOKEN_L_PAREN,      "("
			case ')': return lexer.pos, TOKEN_R_PAREN,      ")"
			case ',': return lexer.pos, TOKEN_COMMA,        ","
			case '*': return lexer.pos, TOKEN_MUL,          "*"
			case '/': return lexer.pos, TOKEN_DIV,          "/"
			case '+': return lexer.pos, TOKEN_PLUS,         "+"
			case '-': return lexer.pos, TOKEN_MINUS,        "-"
			case '{': return lexer.pos, TOKEN_DO,           "{"
			case '}': return lexer.pos, TOKEN_END,          "}"
			default:
				if unicode.IsSpace(r) {
					continue
				} else if r == '#' {
					for {
						r, _, err := lexer.reader.ReadRune()
						if r == '\n' {break}
						if err != nil {panic(err)}
						lexer.pos.column++
					}
					continue
				} else if r == '=' {
					r, _, err := lexer.reader.ReadRune()
					if err != nil {panic(err)}
					lexer.pos.column++
					if r == '=' {
						return lexer.pos, TOKEN_IS_EQUALS, "=="
					} else {
						return lexer.pos, TOKEN_EQUALS, "="
					}
				} else if r == '!' {
					r, _, err := lexer.reader.ReadRune()
					if err != nil {panic(err)}
					lexer.pos.column++
					if r == '=' {
						return lexer.pos, TOKEN_NOT_EQUALS, "!="
					} else {
						fmt.Println("Error:")
						os.Exit(0)
					}
				} else if r == '<' {
					r, _, err := lexer.reader.ReadRune()
					if err != nil {panic(err)}
					lexer.pos.column++
					if r == '=' {
						return lexer.pos, TOKEN_LESS_EQUALS, "<="
					} else {
						return lexer.pos, TOKEN_LESS_THAN, "<"
					}
				} else if r == '>' {
					r, _, err := lexer.reader.ReadRune()
					if err != nil {panic(err)}
					lexer.pos.column++
					if r == '=' {
						return lexer.pos, TOKEN_GREATER_EQUALS, ">="
					} else {
						return lexer.pos, TOKEN_GREATER_THAN, ">"
					}
				} else if unicode.IsDigit(r) {
					startPos := lexer.pos
					lexer.backup()
					val := lexer.lexInt()
					return startPos, TOKEN_INT, val
				} else if unicode.IsLetter(r) {
					startPos := lexer.pos
					lexer.backup()
					val := lexer.lexId()
					if val == "true" || val == "false" {
						return startPos, TOKEN_BOOL, val
					} else if val == "none" {
						return startPos, TOKEN_NONE, val
					} else if val == "else" {
						return startPos, TOKEN_ELSE, val
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
	ExprBool
	ExprNone
	ExprVar
	ExprVardef
	ExprFunCall
	ExprFunDef
	ExprReturn
	ExprIf
	ExprBinop
	ExprCompare
)

type Expr struct {
	Type ExprType
	AsInt int
	AsStr string
	AsNone string
	AsVar string
	AsBool bool
	AsVardef *Vardef
	AsFunCall *FunCall
	AsFunDef *FunDef
	AsReturn *Return
	AsBinop *Binop
	AsIf *If
	AsCompare *Compare
}

type Vardef struct {
	Name string
	Value Expr
	FunName *string
}

type FunCall struct {
	Name string
	Args []Expr
}

type FunDef struct {
	Name string
	Args []Expr
	Body []Expr
}

type Return struct {
	Value Expr
}

type If struct {
	Op Expr
	Body []Expr
	ElseBody []Expr
}

type Binop struct {
	Left Expr
	Right Expr
	Op int
}

type Compare struct {
	Left Expr
	Right Expr
	Op int
}

// -----------------------------
// ----------- Parse -----------
// -----------------------------

type Parser struct {
	CurrentTokenType Token
	CurrentTokenValue string
	lexer Lexer
	line int
	column int
}

func ParserInit(lexer *Lexer) *Parser {
	pos, tok, val := lexer.Lex()
	return &Parser{
		CurrentTokenType: tok,
		CurrentTokenValue: val,
		lexer: *lexer,
		line: pos.line,
		column: pos.column,
	}
}

func (parser *Parser) ParserEat(token Token) {
	if token != parser.CurrentTokenType {
		fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value '%s'", parser.line, parser.column, parser.CurrentTokenValue))
		os.Exit(0)
	}
	pos, tok, val := parser.lexer.Lex()
	parser.CurrentTokenType = tok
	parser.CurrentTokenValue = val
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

func ParserParseFunCall(parser *Parser, FunName string) (Expr) {
	expr := Expr{}
	expr.Type = ExprFunCall
	name := FunName
	parser.ParserEat(TOKEN_L_PAREN)
	args := []Expr{}
	if parser.CurrentTokenType != TOKEN_R_PAREN {
		for {
			arg := ParserParseExpr(parser)
			args = append(args, arg)
			if parser.CurrentTokenType == TOKEN_R_PAREN {break}
			parser.ParserEat(TOKEN_COMMA)
		}
	} else {
		args = nil
	}
	parser.ParserEat(TOKEN_R_PAREN)
	expr.AsFunCall = &FunCall {
		Name: name,
		Args: args,
	}
	return expr
}

func ParserParseVardef(parser *Parser, name string) (Expr) {
	expr := Expr{}
	expr.Type = ExprVardef
	parser.ParserEat(TOKEN_EQUALS)
	value := ParserParseExpr(parser)
	expr.AsVardef = &Vardef {
		Name: name,
		Value: value,
	}
	return expr
}

func ParserParseVariable(parser *Parser) (Expr) {
	expr := Expr{}
	expr.Type = ExprVar
	name := parser.CurrentTokenValue
	parser.ParserEat(TOKEN_ID)
	if parser.CurrentTokenType == TOKEN_L_PAREN {
		return ParserParseFunCall(parser, name)
	}
	if parser.CurrentTokenType == TOKEN_EQUALS {
		return ParserParseVardef(parser, name)
	}
	expr.AsVar = name
	return expr
}

func ParserParseFunDef(parser *Parser) (Expr) {
	expr := Expr{}
	expr.Type = ExprFunDef
	parser.ParserEat(TOKEN_ID)
	name := parser.CurrentTokenValue
	parser.ParserEat(TOKEN_ID)
	parser.ParserEat(TOKEN_L_PAREN)
	args := []Expr{}
	if parser.CurrentTokenType != TOKEN_R_PAREN {
		for {
			arg := ParserParseId(parser)
			args = append(args, arg)
			if parser.CurrentTokenType == TOKEN_R_PAREN {break}
			parser.ParserEat(TOKEN_COMMA)
		}
	} else {
		args = nil
	}
	parser.ParserEat(TOKEN_R_PAREN)
	parser.ParserEat(TOKEN_DO)
	FunBody := ParserParse(parser, &name)
	parser.ParserEat(TOKEN_END)
	expr.AsFunDef = &FunDef {
		Name: name,
		Body: FunBody,
		Args: args,
	}
	return expr
}

func ParserParseReturn(parser *Parser) (Expr) {
	expr := Expr{}
	parser.ParserEat(TOKEN_ID)
	expr.Type = ExprReturn
	ReturnValue := ParserParseExpr(parser)
	expr.AsReturn = &Return {
		Value: ReturnValue,
	}
	return expr
}

func ParserParseIf(parser *Parser) (Expr) {
	expr := Expr{}
	expr.Type = ExprIf
	parser.ParserEat(TOKEN_ID)
	Op := ParserParseExpr(parser)
	parser.ParserEat(TOKEN_DO)
	Body := ParserParse(parser, nil)
	parser.ParserEat(TOKEN_END)
	var ElseBody []Expr
	if parser.CurrentTokenType == TOKEN_ELSE {
		parser.ParserEat(TOKEN_ELSE)
		parser.ParserEat(TOKEN_DO)
		ElseBody = ParserParse(parser, nil)
		parser.ParserEat(TOKEN_END)
	} else {
		ElseBody = nil
	}
	expr.AsIf = &If {
		Op: Op,
		Body: Body,
		ElseBody: ElseBody,
	}
	return expr
}

func ParserParseId(parser *Parser) (Expr) {
	if parser.CurrentTokenValue == "fn" {
		return ParserParseFunDef(parser)
	} else if parser.CurrentTokenValue == "return" {
		return ParserParseReturn(parser)
	} else if parser.CurrentTokenValue == "if" {
		return ParserParseIf(parser)
	}
	return ParserParseVariable(parser)
}

func ParserParseString(parser *Parser) (Expr) {
	expr := Expr{}
	expr.Type = ExprStr
	expr.AsStr = parser.CurrentTokenValue
	parser.ParserEat(TOKEN_STRING)
	return expr
}

func ParserParseInt(parser *Parser) (Expr) {
	expr := Expr{}
	expr.Type = ExprInt
	expr.AsInt = StrToInt(parser.CurrentTokenValue)
	parser.ParserEat(TOKEN_INT)
	return expr
}

func ParserParseBool(parser *Parser) (Expr) {
	expr := Expr{}
	expr.Type = ExprBool
	var BoolValue bool
	if parser.CurrentTokenValue == "true" {
		BoolValue = true
	} else {
		BoolValue = false
	}
	expr.AsBool = BoolValue
	parser.ParserEat(TOKEN_BOOL)
	return expr
}

func ParserParseNone(parser *Parser) (Expr) {
	expr := Expr{}
	expr.Type = ExprNone
	expr.AsNone = "none"
	parser.ParserEat(TOKEN_NONE)
	return expr
}

func ParserParseTest(parser *Parser) (Expr) {
	LeftExpr := Expr{}
	if parser.CurrentTokenType == TOKEN_ID {
		LeftExpr = ParserParseId(parser)
	} else if parser.CurrentTokenType == TOKEN_STRING {
		LeftExpr = ParserParseString(parser)
	} else if parser.CurrentTokenType == TOKEN_INT {
		LeftExpr = ParserParseInt(parser)
	} else if parser.CurrentTokenType == TOKEN_BOOL {
		LeftExpr = ParserParseBool(parser)
	} else if parser.CurrentTokenType == TOKEN_NONE {
		LeftExpr = ParserParseNone(parser)
	}
	return LeftExpr
}

func ParserParseFactor(parser *Parser) (Expr) {
	LeftExpr := Expr{}
	if parser.CurrentTokenType == TOKEN_ID {
		LeftExpr = ParserParseId(parser)
	} else if parser.CurrentTokenType == TOKEN_STRING {
		LeftExpr = ParserParseString(parser)
	} else if parser.CurrentTokenType == TOKEN_INT {
		LeftExpr = ParserParseInt(parser)
	} else if parser.CurrentTokenType == TOKEN_BOOL {
		LeftExpr = ParserParseBool(parser)
	} else if parser.CurrentTokenType == TOKEN_NONE {
		LeftExpr = ParserParseNone(parser)
	}
	for {
		if parser.CurrentTokenType != TOKEN_MUL && parser.CurrentTokenType != TOKEN_DIV {break}
		op := int(parser.CurrentTokenType)
		parser.ParserEat(parser.CurrentTokenType)
		expr := Expr{}
		expr.Type = ExprBinop
		left := LeftExpr
		right := ParserParseTest(parser)
		expr.AsBinop = &Binop {
			Left: left,
			Right: right,
			Op: op,
		}
		LeftExpr = expr
	}
	return LeftExpr
}

func ParserParseTerm(parser *Parser) (Expr) {
	ExprLeft := ParserParseFactor(parser)
	for {
		if parser.CurrentTokenType != TOKEN_PLUS && parser.CurrentTokenType != TOKEN_MINUS { break }
		op := int(parser.CurrentTokenType)
		parser.ParserEat(parser.CurrentTokenType)
		expr := Expr{}
		expr.Type = ExprBinop
		left := ExprLeft
		right := ParserParseFactor(parser)
		expr.AsBinop = &Binop {
			Left: left,
			Right: right,
			Op: op,
		}
		ExprLeft = expr
	}
	if parser.CurrentTokenType == TOKEN_IS_EQUALS || parser.CurrentTokenType == TOKEN_LESS_THAN || parser.CurrentTokenType == TOKEN_GREATER_THAN || parser.CurrentTokenType == TOKEN_GREATER_EQUALS || parser.CurrentTokenType == TOKEN_LESS_EQUALS || parser.CurrentTokenType == TOKEN_NOT_EQUALS {
		op := int(parser.CurrentTokenType)
		parser.ParserEat(parser.CurrentTokenType)
		CompareExpr := Expr{}
		CompareExpr.Type = ExprCompare
		left := ExprLeft
		right := ParserParseFactor(parser)
		CompareExpr.AsCompare = &Compare {
			Left: left,
			Right: right,
			Op: op,
		}
		ExprLeft = CompareExpr
	}
	return ExprLeft
}

func ParserParseExpr(parser *Parser) (Expr) {
	return ParserParseTerm(parser)
}

func FunNotReturn(exprs []Expr) (bool) {
	for _, expr := range(exprs) {
		if expr.Type == ExprReturn {
			return false
		}
	}
	return true
}

func ParserParse(parser *Parser, FunName *string) ([]Expr) {
	exprs := []Expr{}
	for {
		if parser.CurrentTokenType == TOKEN_ID {
			expr := ParserParseId(parser)
			exprs = append(exprs, expr)
		} else if parser.CurrentTokenType == TOKEN_EOF || parser.CurrentTokenType == TOKEN_END || parser.CurrentTokenType == TOKEN_ELSE {
			if FunName != nil {
				if FunNotReturn(exprs) {
					fmt.Println(fmt.Sprintf("error: function '%s()' does not return anything", *FunName))
					os.Exit(0)
				}
			}
			return exprs
		} else {
			fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value '%s'", parser.line, parser.column, parser.CurrentTokenValue))
			os.Exit(0)
		}
	}
	return exprs
}


// -----------------------------
// -------- Visit Exprs --------
// -----------------------------

type FunVariableScope struct {
	Scope map[string]Expr
}

var FunScope = []Expr{}
var GlobalVariableScope = map[string]Expr{}

func VisitExprFuncCall(expr Expr, VariableScope *FunVariableScope) (Expr, bool, *FunVariableScope) {
	if expr.AsFunCall.Name == "print" {
		for _, arg := range(expr.AsFunCall.Args) {
			if arg.Type == ExprIf {
				VisitedIf, BoolValue, _ := VisitExprs(arg, false, VariableScope)
				if !BoolValue {
					NoneExpr := Expr{}
					NoneExpr.Type = ExprNone
					NoneExpr.AsNone = "none"
					arg = NoneExpr
				} else {
					arg = VisitedIf
				}
			}
			VisitedArg, _, _ := VisitExprs(arg, false, VariableScope)
			switch VisitedArg.Type {
				case ExprStr: fmt.Print(VisitedArg.AsStr)
				case ExprInt: fmt.Print(VisitedArg.AsInt)
				case ExprBool: fmt.Print(VisitedArg.AsBool)
				case ExprNone: fmt.Print(VisitedArg.AsNone)
				default:
					fmt.Println("Error: print error.")
					os.Exit(0)
			}
			fmt.Print(" ")
		}
		fmt.Println()
	} else if expr.AsFunCall.Name == "exit" {
		os.Exit(0)
	} else {
		for _, VisitedFun := range(FunScope) {
			if VisitedFun.AsFunDef.Name == expr.AsFunCall.Name {
				if len(VisitedFun.AsFunDef.Args) != len(expr.AsFunCall.Args) {
					if len(VisitedFun.AsFunDef.Args) > len(expr.AsFunCall.Args) {
						fmt.Println(fmt.Sprintf("Error: not enough arguments to call function %s()", expr.AsFunCall.Name))
					} else {
						fmt.Println(fmt.Sprintf("Error: too many arguments to call function %s()", expr.AsFunCall.Name))
					}
					os.Exit(0)
				}
				VariableScopeMap := map[string]Expr{}
				FunVariableScopeInit := FunVariableScope{
					Scope: VariableScopeMap,
				}
				for i := 0; i < len(expr.AsFunCall.Args); i++ {
					VarExpr := Expr{}
					VarExpr.Type = ExprVardef
					if expr.AsFunCall.Args[i].Type == ExprIf {
						VisitedIf, BoolValue, _ := VisitExprs(expr.AsFunCall.Args[i], false, VariableScope)
						if !BoolValue {
							NoneExpr := Expr{}
							NoneExpr.Type = ExprNone
							NoneExpr.AsNone = "none"
							expr.AsFunCall.Args[i] = NoneExpr
						} else {
							expr.AsFunCall.Args[i] = VisitedIf
						}
					}
					VarExpr.AsVardef = &Vardef {
						Name: VisitedFun.AsFunDef.Args[i].AsVar,
						Value: expr.AsFunCall.Args[i],
					}
					
					VisitExprVardefFunArg(VarExpr, &FunVariableScopeInit, VariableScope)
				}
				RetExpr, RetBool, _ := Exprs(VisitedFun.AsFunDef.Body, &FunVariableScopeInit)
				return RetExpr, RetBool, VariableScope
			}
		}
		fmt.Println(fmt.Sprintf("Error: undifined function '%s'.", expr.AsFunCall.Name))
		os.Exit(0)
	}
	return expr, false, VariableScope
}

func VisitExprFunDef(expr Expr) {
	for _, VisitedFun := range(FunScope) {
		if VisitedFun.AsFunDef.Name == expr.AsFunDef.Name {
			fmt.Println(fmt.Sprintf("Error: function '%s' is already defined.", expr.AsFunDef.Name)); os.Exit(0);
		}
	}
	FunScope = append(FunScope, expr)
}


func VisitExprVardefFunArg(expr Expr, VariableScope *FunVariableScope, PrevVariableScope *FunVariableScope) {
	VisitedVar, BoolValue, _ := VisitExprs(expr.AsVardef.Value, false, PrevVariableScope)
	if expr.AsVardef.Value.Type == ExprIf {
		if !BoolValue {
			NoneExpr := Expr{}
			NoneExpr.Type = ExprNone
			NoneExpr.AsNone = "none"
			VisitedVar = NoneExpr
		}
	}
	if VariableScope == nil {
		GlobalVariableScope[expr.AsVardef.Name] = VisitedVar
	} else {
		VariableScope.Scope[expr.AsVardef.Name] = VisitedVar
	}
}

func VisitExprVardef(expr Expr, VariableScope *FunVariableScope) {
	VisitedVar, BoolValue, _ := VisitExprs(expr.AsVardef.Value, false, VariableScope)
	if expr.AsVardef.Value.Type == ExprIf {
		if !BoolValue {
			NoneExpr := Expr{}
			NoneExpr.Type = ExprNone
			NoneExpr.AsNone = "none"
			VisitedVar = NoneExpr
		}
	}
	if VariableScope == nil {
		GlobalVariableScope[expr.AsVardef.Name] = VisitedVar
	} else {
		VariableScope.Scope[expr.AsVardef.Name] = VisitedVar
	}
}

func VisitExprVar(expr Expr, VariableScope *FunVariableScope) (Expr) {
	var VisitedVar Expr

	if VariableScope == nil {
		// global
		if _, ok := GlobalVariableScope[expr.AsVar]; ok {
			VisitedVar = GlobalVariableScope[expr.AsVar]
		} else {
			fmt.Println(fmt.Sprintf("Error: variable '%s' is not defined", expr.AsVar)); os.Exit(0);
		}
	} else {
		// function variable scope
		if _, ok := VariableScope.Scope[expr.AsVar]; ok {
			VisitedVar = VariableScope.Scope[expr.AsVar]
		} else {
			if _, ok := GlobalVariableScope[expr.AsVar]; ok {
				VisitedVar = GlobalVariableScope[expr.AsVar]
			} else {
				fmt.Println(fmt.Sprintf("Error: variable '%s' is not defined", expr.AsVar)); os.Exit(0);
			}
		}
	}
	return VisitedVar
}

func VisitExprIf(expr Expr, VariableScope *FunVariableScope) (Expr, bool, *FunVariableScope) {
	VisitedBool, _, _ := VisitExprs(expr.AsIf.Op, false, VariableScope)
	if VisitedBool.AsBool {
		RetExpr, BoolValue, VariableScope := Exprs(expr.AsIf.Body, VariableScope)
		return RetExpr, BoolValue, VariableScope
	} else {
		RetExpr, BoolValue, VariableScope := Exprs(expr.AsIf.ElseBody, VariableScope)
		return RetExpr, BoolValue, VariableScope
	}
	return expr, false, VariableScope
}

func VisitExprBinop(expr Expr, VariableScope *FunVariableScope) (Expr) {
	VisitedLeft, BoolLeft, _ := VisitExprs(expr.AsBinop.Left, false, VariableScope)
	VisitedRight, BoolRight, _ := VisitExprs(expr.AsBinop.Right, false, VariableScope)
	if expr.AsBinop.Left.Type == ExprIf {
		if !BoolLeft {
			fmt.Println("Error: binop left value if statement didn't return.")
			os.Exit(0)
		}
	}
	if expr.AsBinop.Right.Type == ExprIf {
		if !BoolRight {
			fmt.Println("Error: binop right value if statement didn't return.")
			os.Exit(0)
		}
	}
	if VisitedLeft.Type != ExprInt || VisitedRight.Type != ExprInt {
		fmt.Println("Error: binary operation expected type <int>.")
		os.Exit(0)
	}
	var Value int
	if expr.AsBinop.Op == TOKEN_PLUS {
		Value = VisitedLeft.AsInt + VisitedRight.AsInt
	} else if expr.AsBinop.Op == TOKEN_MINUS {
		Value = VisitedLeft.AsInt - VisitedRight.AsInt
	} else if expr.AsBinop.Op == TOKEN_MUL {
		Value = VisitedLeft.AsInt * VisitedRight.AsInt
	}	else if expr.AsBinop.Op == TOKEN_DIV {
		Value = VisitedLeft.AsInt / VisitedRight.AsInt
	}
	IntExpr := Expr{}
	IntExpr.Type = ExprInt
	IntExpr.AsInt = Value
	return IntExpr
}

func VisitExprCompare(expr Expr, VariableScope *FunVariableScope) (Expr) {
	VisitedLeft, BoolLeft, _ := VisitExprs(expr.AsCompare.Left, false, VariableScope)
	VisitedRight, BoolRight, _ := VisitExprs(expr.AsCompare.Right, false, VariableScope)
	if expr.AsCompare.Left.Type == ExprIf {
		if !BoolLeft {
			NoneExpr := Expr{}
			NoneExpr.Type = ExprNone
			NoneExpr.AsNone = "none"
			VisitedLeft = NoneExpr
		}
	}
	if expr.AsCompare.Right.Type == ExprIf {
		if !BoolRight {
			NoneExpr := Expr{}
			NoneExpr.Type = ExprNone
			NoneExpr.AsNone = "none"
			VisitedRight = NoneExpr
		}
	}
	var Value bool
	if expr.AsCompare.Op == TOKEN_IS_EQUALS {
		if VisitedLeft.Type != VisitedRight.Type {
			Value = false
		} else if VisitedLeft.Type == ExprStr {
			Value = VisitedLeft.AsStr == VisitedRight.AsStr
		} else if VisitedLeft.Type == ExprInt {
			Value = VisitedLeft.AsInt == VisitedRight.AsInt
		} else if VisitedLeft.Type == ExprNone {
			Value = true
		}
	} else if expr.AsCompare.Op == TOKEN_NOT_EQUALS {
		if VisitedLeft.Type != VisitedRight.Type {
			Value = true
		} else if VisitedLeft.Type == ExprStr {
			Value = VisitedLeft.AsStr != VisitedRight.AsStr
		} else if VisitedLeft.Type == ExprInt {
			Value = VisitedLeft.AsInt != VisitedRight.AsInt
		} else if VisitedLeft.Type == ExprNone {
			Value = false
		}
	} else if expr.AsCompare.Op == TOKEN_LESS_THAN {
		if VisitedLeft.Type == ExprInt && VisitedRight.Type == ExprInt {
			Value = VisitedLeft.AsInt < VisitedRight.AsInt
		} else {
			fmt.Println("Error: '<' expected type <int>")
			os.Exit(0)
		}
	} else if expr.AsCompare.Op == TOKEN_GREATER_THAN {
		if VisitedLeft.Type == ExprInt && VisitedRight.Type == ExprInt {
			Value = VisitedLeft.AsInt > VisitedRight.AsInt
		} else {
			fmt.Println("Error: '>' expected type <int>")
			os.Exit(0)
		}
	} else if expr.AsCompare.Op == TOKEN_GREATER_EQUALS {
		if VisitedLeft.Type == ExprInt && VisitedRight.Type == ExprInt {
			Value = VisitedLeft.AsInt >= VisitedRight.AsInt
		} else {
			fmt.Println("Error: '>=' expected type <int>")
			os.Exit(0)
		}
	} else if expr.AsCompare.Op == TOKEN_LESS_EQUALS {
		if VisitedLeft.Type == ExprInt && VisitedRight.Type == ExprInt {
			Value = VisitedLeft.AsInt <= VisitedRight.AsInt
		} else {
			fmt.Println("Error: '<=' expected type <int>")
			os.Exit(0)
		}
	}
	BoolExpr := Expr{}
	BoolExpr.Type = ExprBool
	BoolExpr.AsBool = Value
	return BoolExpr
}

func VisitExprs(expr Expr, RetBool bool, VariableScope *FunVariableScope) (Expr, bool, *FunVariableScope) {
	switch expr.Type {
		case ExprFunCall:
			RetExpr, _, VariableScope := VisitExprFuncCall(expr, VariableScope)
			return RetExpr, RetBool, VariableScope
		case ExprVardef:
			VisitExprVardef(expr, VariableScope)
		case ExprStr:
			return expr, RetBool, VariableScope
		case ExprInt:
			return expr, RetBool, VariableScope
		case ExprBool:
			return expr, RetBool, VariableScope
		case ExprNone:
			return expr, RetBool, VariableScope
		case ExprVar:
			return VisitExprVar(expr, VariableScope), RetBool, VariableScope
		case ExprReturn:
			return VisitExprs(expr.AsReturn.Value, true, VariableScope)
		case ExprFunDef:
			VisitExprFunDef(expr)
		case ExprIf:
			RetExpr, RetBoolIf, VariableScope := VisitExprIf(expr, VariableScope)
			return RetExpr, RetBoolIf, VariableScope
		case ExprBinop:
			return VisitExprBinop(expr, VariableScope), RetBool, VariableScope
		case ExprCompare:
			return VisitExprCompare(expr, VariableScope), RetBool, VariableScope
	}
	return expr, RetBool, VariableScope
}

func Exprs(exprs []Expr, VariableScope *FunVariableScope) (Expr, bool, *FunVariableScope) {
	var RetExpr Expr
	var RetBool bool
	for _, expr := range exprs {
		RetExpr, RetBool, VariableScope = VisitExprs(expr, false, VariableScope)
		if RetBool {
			break
		}
	}
	return RetExpr, RetBool, VariableScope
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
	if len(os.Args) != 2 {
		Usage()
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error: file '" + os.Args[1] + "' does not exist")
		os.Exit(0)
	}

	lexer := LexerInit(file)
	parser := ParserInit(lexer)
	exprs := ParserParse(parser, nil)
	Exprs(exprs, nil)
}

