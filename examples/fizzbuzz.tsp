
1 
for dup 100 <= do
    if dup 15 % 0 == do
        "FizzBuzz" print
    elif dup 3 % 0 == do
        "Fizz" print
    elif dup 5 % 0 == do
        "Buzz" print
    else
        dup print
    end
    inc
end drop

