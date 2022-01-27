
1 
for dup 100 <= do
    if dup 15 % 0 == do
        "FizzBuzz\n" print
    elif dup 3 % 0 == do
        "Fizz\n" print
    elif dup 5 % 0 == do
        "Buzz\n" print
    else
        dup println
    end
    inc
end drop

