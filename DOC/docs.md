# T# Documentation

## Introduction

T# is a Stack-based programming language designed for building software.
It's similar to Porth and Forth.

## Install & Run
```bash
$ git clone https://github.com/Tsharp-lang/Tsharp
$ cd tsharp
$ go build main.go
$ ./main examples/main.tsp
or
$ ./main.exe examples/main.tsp
```

## Table of Contents
<table>
    <tr><td width=33% valign=top>

* [Hello World](#hello-world)
* [Comments](#comments)
* [Block](#block)
* [Arithmetic](#arithmetic)
* [If Statement](#if-statement)
* [For loop](#for-loop)
* [Error handling](#error-handling)
* [Import](#import)

</td><td width=33% valign=top>

* [Dup](#dup)
* [Drop](#drop)
* [Over](#over)
* [Rot](#rot)

</td><td valign=top>

</td></tr>
</table>

## Hello World
```pascal
"Hello World!" print
```
'print' will print the top element of the stack, then remove it.

## Comments
```python
# Sample comment
```

## Import
```python
import "main.t#"
```

## Block
```pascal
block main do
    "Hello World!" print
end

call main
```

'block' is like Function in other languages.

## If Statement
```pascal
if false do
    "Hello World" print
elif false do
    "elif body!" print
else
    "else body!" print
end

10 10 == print
20 10 != print
2 10 < print
10 2 > print
```

## Dup
```pascal
"Hello World" dup print print
```
'dup' duplicate element on top of the stack.

## Drop
```pascal
"Hello World" "T# Programming Language" drop print
```
'drop' drops the top element of the stack.

## PrintS
```python
1 2 "Hello World!"

printS

# stack length  
#     ↓ 
#    <3>  1 2 'Hello World' <- top
```
'printS' print all stack values. 'printS' won't drop stack value after print.

## PrintV
```python
10 -> N
"Hello World!" -> text

printV

# N : 10
# text : "Hello World!"

```
'printV' print all variables.

## For loop
```pascal
for true do
    "Hello World!" print
    break
end
```

## Error handling
```python

# errors
# StackOutRange
# TypeError
# ImportError

try 
    dup
except StackOutRange do
    "Error" print
end
```

## Arithmetic
```pascal
34 35 + print

100 40 - print

200 5 / print

10 2 * print
```

## Variable
```pascal
10 -> x

x -> y

y print
```

## Type
```python
int # 12345
string # "Hello World!"
bool # true false
type # int string bool type
```

## Typeof
```python
"Hello World" dup typeof print
```

## Rot
```python
1 2 3 rot print print print
```
'rot' rotate top three stack elements.

## Over
```python
1 2 over print print print
```
'over' copy the element below the top of the stack

## append string
```python
"Hello " "World!" + print 
```

## Inc
```python
1 inc print
```
'inc' increment the top element of the stack

## Dec
```python
10 dec print
```
'dec' decrement the top element of the stack

## Exit
```python
"Hello World"
exit
print
```
'exit' will exit the program.


## List
```python
[] # push empty list

"Hello World!" append # append string "Hello World!"

34 append # append int 34

68 append # append int 68

"T# Programming Language" 1 replace # replace list index '1' to string "T# Programming Language"

print
```


## FizzBuzz
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
