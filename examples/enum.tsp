
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

true call Enum -> MON
call Enum -> TUE
call Enum -> WED
call Enum -> THU
call Enum -> FRI
call Enum -> SAT
call Enum -> SUN

MON print
TUE print
WED print
THU print
FRI print
SAT print
SUN print


