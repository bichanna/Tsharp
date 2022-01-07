
block LexerInit do
    0 -> Index
    [] -> Tokens
end

block LexerAdvance do
    Index 1 + -> Index
end

block Lexer do
    call LexerInit
    for true do
        if Index Length == do
            break
        end
        PROGRAM Index read -> char drop
        if char " " == do
            call LexerAdvance
        elif char "'" == do
            # "TOKEN_STRING" print
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
            [] "TOKEN_STRING" append STRINGValue append -> STRINGArr
            Tokens STRINGArr append -> Tokens
            # STRINGValue print
        elif char "+" == do
            # "TOKEN_PLUS" print
            Tokens "TOKEN_PLUS" append -> Tokens
            call LexerAdvance
        elif char "-" == do
            # "TOKEN_MINUS" print
            Tokens "TOKEN_MINUS" append -> Tokens
            call LexerAdvance
        elif char "*" == do
            # "TOKEN_MUL" print
            Tokens "TOKEN_MUL" append -> Tokens
            call LexerAdvance
        elif char "/" == do
            # "TOKEN_DIV" print
            Tokens "TOKEN_DIV" append -> Tokens
            call LexerAdvance
        else
            # "TOKEN_ID" print
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
            [] "TOKEN_ID" append IDValue append -> IDArr
            Tokens IDArr append -> Tokens
            # IDValue print
        end
    end
end

block Main do
    "print 'Hello World' num +" -> PROGRAM # program

    PROGRAM len -> Length drop # Length of the program

    call Lexer

    Tokens print
end

call Main
