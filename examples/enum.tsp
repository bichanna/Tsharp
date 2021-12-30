
0 -> EnumIndex
block Enum do
    if dup typeof bool == do
        -> BoolValue
        if BoolValue do
            0 -> EnumIndex
            EnumIndex
        else
            EnumIndex
        end
    else
        "Error: block 'Enum' expected type <bool>" print
    end
    EnumIndex 1 + -> EnumIndex
end


true  call Enum -> MON
false call Enum -> TUE
false call Enum -> WED
false call Enum -> THU
false call Enum -> FRI
false call Enum -> SAT
false call Enum -> SUN

MON print
SUN print
