
block fmt do
    try
        if dup typeof list != do
            "`fmt` expected type <list> but got " puts puts "." print drop
            exit 
        end
    except StackIndexError do
        "`fmt` the stack is empty." print
        exit
    end
    dup len 1 - -> fmt_list_length

    -> fmt_list

    fmt_list 0 read -> fmt_string

    fmt_string len -> fmt_string_len

    "" -> fmt_string_val
    
    1 -> fmt_count

    0 for dup fmt_string_len < do -> fmt_string_index
        fmt_string fmt_string_index read -> fmt_string_char
        if fmt_string_char "$" == do
            try
                fmt_string_val fmt_list fmt_count read + -> fmt_string_val
            except IndexError do
                "`fmt` index error" print
                exit
            end
            fmt_count inc -> fmt_count
        else
            fmt_string_val fmt_string_char + -> fmt_string_val
        end
        fmt_string_index inc
    end drop

    fmt_string_val
end


# "game" -> N
# { "Hello $ developers" N } call fmt print