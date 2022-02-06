

block Fib do -> N
    if N 2 < do
        N
    else
        N 1 - call Fib
        N 2 - call Fib
        +
    end
end

10 call Fib println
