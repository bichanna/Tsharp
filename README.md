<div align="center">    
    <img width="100px" src="https://user-images.githubusercontent.com/81926489/143374038-059715ef-a83d-479d-a8c3-56ea57b8cc8e.PNG">
    <h1> The T# Programming Language</h1>
    <a href="https://github.com/Tsharp-lang/Tsharp/blob/main/DOC/docs.md">Doc</a>
    |
    <a href="https://github.com/Tsharp-lang/Tsharp/blob/main/DOC/ドキュメント.md">Doc(日本語)</a>
    |
    <a href="https://github.com/Tsharp-lang/Tsharp/blob/main/editor/tsharp.vim">Vim</a>
    |
    <a href="https://marketplace.visualstudio.com/items?itemName=akamurasaki.tsharplanguage-color">VSCode</a>
</div>

[![Ubuntu](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-ubuntu.yml/badge.svg)](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-ubuntu.yml)
[![Windows](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-windows.yml/badge.svg)](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-windows.yml)
[![CodeQL](https://github.com/Tsharp-lang/Tsharp/actions/workflows/codeql-analysis.yml/badge.svg?branch=main)](https://github.com/Tsharp-lang/Tsharp/actions/workflows/codeql-analysis.yml)

WARNING! THIS LANGUAGE IS A WORK IN PROGRESS! ANYTHING CAN CHANGE AT ANY MOMENT WITHOUT ANY NOTICE!

It's like Forth and Porth, but written in Go.
<a href="https://en.wikipedia.org/wiki/Stack-oriented_programming">Stack-oriented programming</a>

### TODO
- [ ] Compile to C
- [ ] Self-hosted

### Install

> Install
```
go build main.go
```

### Run

> Run
```
$ ./main <filename>.tsp

or

$ ./main.exe <filename>.tsp
```

### Hello World
```pascal
"Hello World" print
```

### Fibonacci Sequence
```pascal
10000 -> N

0 1 for over N < do
  over puts " " puts
  swap over +
end
drop drop

"" print
```

### Bubble Sort
```pascal
{ 19 13 6 2 18 8 1 4 11 9 100 30 4 } dup dup print -> arr

13 -> length

0 for dup length <= do
    0 for dup length 1 - < do
        dup -> j
        j 1 + -> i
        if arr j read swap i read swap drop > do
            arr j read -> x
            i read -> y
            y j replace
            x i replace
            drop
        end 
        inc
    end drop
    inc
end drop

print
```

### Factorial
```pascal
block Factorial do
    -> N
    1 -> x
    for N 1 >= do
        x N * -> x
        N 1 - -> N
    end
    x
end

5
call Factorial
print
```

### FizzBuzz
```pascal
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
```

### Multiplication table
```pascal
block dclone do
    dup -> tmpa
    swap
    dup -> tmpb
    swap
    tmpb
    tmpa
end

1 for dup 10 < do
    1 for dup 10 < do
        call dclone
        *
        if dup 10 < do
            " " puts
        end
        puts
        " " puts
        inc
    end
    " " print
    drop
    inc
end
```

### Contributing
Welcome! 💕
