
# Example
# 'Hello World!' print


# ------------------------
# -------- TOKENS --------
# ------------------------

0 -> EnumIndex
block Enum do
    try
        if dup typeof bool == do
            -> BoolValue
            if BoolValue do
                0 -> EnumIndex
                EnumIndex
            else
                EnumIndex
            end
        else
            EnumIndex
        end
    except StackOutRange do
        EnumIndex
    end
    EnumIndex 1 + -> EnumIndex
end

true call Enum -> TOKEN_STRING
call Enum -> TOKEN_ID
call Enum -> TOKEN_EOF
call Enum -> TOKEN_PLUS
call Enum -> TOKEN_MINUS
call Enum -> TOKEN_MUL
call Enum -> TOKEN_DIV


# ------------------------
# -------- Lexer ---------
# ------------------------

block LexerInit do
    PROGRAM len -> Length drop
    0 -> Index
    [] -> Tokens
end

block LexerAdvance do
    Index inc -> Index
end

block Lexer do
    call LexerInit
    for true do
        if Index Length == do
            [] TOKEN_EOF append "TOKEN_EOF" append -> EOFArr
            Tokens EOFArr append -> Tokens
            break
        end
        PROGRAM Index read -> char drop
        if char " " == do
            call LexerAdvance
        elif char "'" == do
            "" -> STRINGValue
            call LexerAdvance
            for true do
                PROGRAM Index read -> char drop
                if char "'" == do
                    call LexerAdvance
                    break
                end
                STRINGValue char + -> STRINGValue
                call LexerAdvance
            end
            [] TOKEN_STRING append STRINGValue append -> STRINGArr
            Tokens STRINGArr append -> Tokens
        elif char "+" == do
            [] TOKEN_PLUS append "+" append -> PLUSArr
            Tokens PLUSArr append -> Tokens
            call LexerAdvance
        elif char "-" == do
            [] TOKEN_MINUS append "-" append -> MINUSArr
            Tokens MINUSArr append -> Tokens
            call LexerAdvance
        elif char "*" == do
            [] TOKEN_MUL append "*" append -> MULArr
            Tokens MULArr append -> Tokens
            call LexerAdvance
        elif char "/" == do
            [] TOKEN_DIV append "/" append -> DIVArr
            Tokens DIVArr append -> Tokens
            call LexerAdvance
        else
            "" -> IDValue
            IDValue char + -> IDValue
            call LexerAdvance
            for true do
                if Index Length == do
                    break
                end
                PROGRAM Index read -> char drop
                if char " " == do
                    break
                end
                IDValue char + -> IDValue
                call LexerAdvance
            end
            [] TOKEN_ID append IDValue append -> IDArr
            Tokens IDArr append -> Tokens
        end
    end
end


# ------------------------
# -------- Stack ---------
# ------------------------

[] -> stack
0 -> stack_len

block op_push do
    stack_len 1 + -> stack_len
    stack swap append -> stack
end

block op_print do
    if stack_len 0 == do
        "Error: the stack is empty." print
    else
        stack len 1 - read print drop
        call op_drop
    end
end

block op_drop do
    if stack_len 0 == do
        "Error: the stack is empty." print
    else
        stack_len 1 - -> stack_len
        stack len 1 - remove -> stack
    end
end

block op_swap do
    if stack_len 2 < do
        "Error: `swap` expected more than 2 elementes in stack." print
    else
        stack len 1 - read -> first
        len 2 - read       -> second
        call op_drop
        call op_drop
        drop
        first call op_push
        second call op_push
    end
end


# ------------------------
# -------- Parser --------
# ------------------------

block ParserInit do
    0 -> ParserIndex
    Tokens ParserIndex read
    swap drop
    0 read -> parser_current_token_type
    1 read -> parser_current_token_value
    drop
end

block parser_parse_string do
    if parser_current_token_value "exit" == do
        TOKEN_ID call ParserEat
        exit
    elif parser_current_token_value "print" == do
        call op_print
        TOKEN_ID call ParserEat
    elif parser_current_token_value "swap" == do
        call op_swap
        TOKEN_ID call ParserEat
    elif parser_current_token_value "drop" == do
        call op_drop
        TOKEN_ID call ParserEat
    else
        "Syntax Error: unexpected token value '" parser_current_token_value + "'" + print
        exit
    end
end

block Parser do
    for true do
        if parser_current_token_type TOKEN_STRING == do
            parser_current_token_value call op_push
            TOKEN_STRING call ParserEat
        elif parser_current_token_type TOKEN_ID == do
            call parser_parse_string
        elif parser_current_token_type TOKEN_PLUS == do
            "Error: parsing `+` are not implemented." print
            TOKEN_PLUS call ParserEat
        elif parser_current_token_type TOKEN_MINUS == do
            "Error: parsing `-` are not implemented." print
            TOKEN_MINUS call ParserEat
        elif parser_current_token_type TOKEN_MUL == do
            "Error: parsing `*` are not implemented." print
            TOKEN_MUL call ParserEat
        elif parser_current_token_type TOKEN_DIV == do
            "Error: parsing `/` are not implemented." print
            TOKEN_DIV call ParserEat
        elif parser_current_token_type TOKEN_EOF == do
            break
        else
            "Syntax Error: unexpected token value '" parser_current_token_value + "'" + print
            exit
        end
    end
end

block ParserEat do
    if parser_current_token_type != do
        "Syntax Error: unexpected token value '" parser_current_token_value + "'" + print
        exit
    end
    ParserIndex inc -> ParserIndex

    Tokens ParserIndex read
    swap drop
    0 read -> parser_current_token_type
    1 read -> parser_current_token_value
    drop
end


# ------------------------
# --------- Main ---------
# ------------------------

block Help do
    "documentation: https://github.com/Tsharp-lang/Tsharp/blob/main/DOC/docs.md" print
end

block Main do
    "T# interpreter" print
    "Type 'help' for more information, 'exit' to quit." print
    for true do
        ">>> " puts input -> PROGRAM
        if PROGRAM "help" == do
            call Help
        else
            call Lexer
            call ParserInit
            call Parser
        end
    end
end

call Main


