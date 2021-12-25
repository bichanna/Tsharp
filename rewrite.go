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
			case '(': return lexer.pos, TOKEN_L_PAREN, "TOKEN_L_PAREN"
			case ')': return lexer.pos, TOKEN_R_PAREN, "TOKEN_R_PAREN"
			case ',': return lexer.pos, TOKEN_COMMA, "TOKEN_COMMA"
			case '*': return lexer.pos, TOKEN_MUL, "TOKEN_MUL"
			case '/': return lexer.pos, TOKEN_DIV, "TOKEN_DIV"
			case '+': return lexer.pos, TOKEN_PLUS, "TOKEN_PLUS"
			case '-': return lexer.pos, TOKEN_MINUS, "TOKEN_MINUS"
			case '{': return lexer.pos, TOKEN_DO,    "TOKEN_DO"
			case '}': return lexer.pos, TOKEN_END,    "TOKEN_END"
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

func VisitExprFuncCall(expr Expr) (Expr, bool) {
	if expr.AsFunCall.Name == "print" {
		for _, arg := range(expr.AsFunCall.Args) {
			if arg.Type == ExprIf {
				VisitedIf, BoolValue := VisitExprs(arg, false)
				if !BoolValue {
					NoneExpr := Expr{}
					NoneExpr.Type = ExprNone
					NoneExpr.AsNone = "none"
					arg = NoneExpr
				} else {
					arg = VisitedIf
				}
			}
			VisitedArg, _ := VisitExprs(arg, false)
			switch VisitedArg.Type {
				case ExprStr: fmt.Println(VisitedArg.AsStr)
				case ExprInt: fmt.Println(VisitedArg.AsInt)
				case ExprBool: fmt.Println(VisitedArg.AsBool)
				case ExprNone: fmt.Println(VisitedArg.AsNone)
			}
		}
	} else {
		for _, VisitedFun := range(FunScope) {
			if VisitedFun.AsFunDef.Name == expr.AsFunCall.Name {
				for i := 0; i < len(expr.AsFunCall.Args); i++ {
					VarExpr := Expr{}
					VarExpr.Type = ExprVardef
					if expr.AsFunCall.Args[i].Type == ExprIf {
						VisitedIf, BoolValue := VisitExprs(expr.AsFunCall.Args[i], false)
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
					VisitExprs(VarExpr, false)
				}
				RetExpr, RetBool := Exprs(VisitedFun.AsFunDef.Body)
				return RetExpr, RetBool
			}
		}
		fmt.Println(fmt.Sprintf("Error: undifined function '%s'.", expr.AsFunCall.Name))
		os.Exit(0)
	}
	return expr, false
}

var VariableScope = []Expr{}

func VisitExprVardef(expr Expr) {
	if expr.AsVardef.Value.Type == ExprFunCall {
		VisitedFunCall, _ := VisitExprs(expr.AsVardef.Value, false)
		expr.AsVardef.Value = VisitedFunCall
	} else if expr.AsVardef.Value.Type == ExprIf {
		VisitedIf, BoolValue := VisitExprs(expr.AsVardef.Value, false)
		if !BoolValue {
			NoneExpr := Expr{}
			NoneExpr.Type = ExprNone
			NoneExpr.AsNone = "none"
			expr.AsVardef.Value = NoneExpr
		} else {
			expr.AsVardef.Value = VisitedIf
		}
	}
	for i := 0; i < len(VariableScope); i++ {
		if VariableScope[i].AsVardef.Name == expr.AsVardef.Name {
			VariableScope[i] = expr
			return
		}
	}
	VariableScope = append(VariableScope, expr)
}

func VisitExprVar(expr Expr) (Expr) {
	for _, VisitedVar := range(VariableScope) {
		if VisitedVar.AsVardef.Name == expr.AsVar {
			VisitedVariableValue, _ := VisitExprs(VisitedVar.AsVardef.Value, false)
			return VisitedVariableValue
		}
	}
	fmt.Println("Error: undefined variable '" + expr.AsVar + "'"); os.Exit(0);
	return expr
}

var FunScope = []Expr{}

func VisitExprFunDef(expr Expr) {
	for _, VisitedFun := range(FunScope) {
		if VisitedFun.AsFunDef.Name == expr.AsFunDef.Name {
			fmt.Println(fmt.Sprintf("Error: function '%s' is already defined.", expr.AsFunDef.Name)); os.Exit(0);
		}
	}
	FunScope = append(FunScope, expr)
}

func VisitExprIf(expr Expr) (Expr, bool) {
	VisitedBool, _ := VisitExprs(expr.AsIf.Op, false)
	if VisitedBool.AsBool {
		RetExpr, BoolValue := Exprs(expr.AsIf.Body)
		return RetExpr, BoolValue
	} else {
		RetExpr, BoolValue := Exprs(expr.AsIf.ElseBody)
		return RetExpr, BoolValue
	}
	return expr, false
}

func VisitExprBinop(expr Expr) (Expr) {
	if expr.AsBinop.Left.Type == ExprIf {
		VisitedIf, BoolValue := VisitExprs(expr.AsBinop.Left, false)
		if !BoolValue {
			NoneExpr := Expr{}
			NoneExpr.Type = ExprNone
			NoneExpr.AsNone = "none"
			expr.AsBinop.Left = NoneExpr
		} else {
			expr.AsBinop.Left = VisitedIf
		}
	}
	if expr.AsBinop.Right.Type == ExprIf {
		VisitedIf, BoolValue := VisitExprs(expr.AsBinop.Right, false)
		if !BoolValue {
			NoneExpr := Expr{}
			NoneExpr.Type = ExprNone
			NoneExpr.AsNone = "none"
			expr.AsBinop.Right = NoneExpr
		} else {
			expr.AsBinop.Right = VisitedIf
		}
	}
	VisitedLeft, _ := VisitExprs(expr.AsBinop.Left, false)
	VisitedRight, _ := VisitExprs(expr.AsBinop.Right, false)
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

func VisitExprs(expr Expr, RetBool bool) (Expr, bool) {
	switch expr.Type {
		case ExprFunCall:
			RetExpr, RetBoolFun := VisitExprFuncCall(expr)
			return RetExpr, RetBoolFun
		case ExprVardef:
			VisitExprVardef(expr)
		case ExprStr:
			return expr, RetBool
		case ExprInt:
			return expr, RetBool
		case ExprNone:
			return expr, RetBool
		case ExprVar:
			return VisitExprVar(expr), RetBool
		case ExprReturn:
			return VisitExprs(expr.AsReturn.Value, true)
		case ExprFunDef:
			VisitExprFunDef(expr)
		case ExprIf:
			RetExpr, RetBoolIf := VisitExprIf(expr)
			return RetExpr, RetBoolIf
		case ExprBinop:
			return VisitExprBinop(expr), RetBool
	}
	return expr, RetBool
}

func Exprs(exprs []Expr) (Expr, bool) {
	var RetExpr Expr
	var RetBool bool
	for _, expr := range exprs {
		RetExpr, RetBool = VisitExprs(expr, false)
		if RetBool {
			break
		}
	}
	return RetExpr, RetBool
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
	Exprs(exprs)
}


