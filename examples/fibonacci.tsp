

block Fib do -> N
    if N 2 < do
        N
    else
        N 1 - Fib
        N 2 - Fib
        +
    end
end

10 Fib println

