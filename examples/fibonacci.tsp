
block Fib do -> n
    if n 2 < do
        n
    else
        n 1 - call Fib
        n 2 - call Fib
        +
    end
end

10 call Fib print
