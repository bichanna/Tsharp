
block Factorial do
    -> n
    1 -> x
    for n 1 >= do
        x n * -> x
        n 1 - -> n
    end
    x
end

5
call Factorial
print
