# T# Documentation

## Introduction

T# is a Stack-based programming language designed for building software.
It's like Porth, Forth.

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
* [Variable](#variable)
* [Arithmetic](#arithmetic)
* [If Statement](#if-statement)

* [For loop](#for-loop)
    * [Loop 100 times](#loop-100-times)
    * [Break](#break)
 
* [Error handling](#error-handling)
* [Import](#import)

</td><td width=33% valign=top>

* [Dup](#dup)
* [Drop](#drop)
* [Over](#over)
* [Rot](#rot)
* [Swap](#swap)
* [Inc](#inc)
* [Dec](#dec)
* [printS](#prints)

</td><td valign=top>

* [Examples](#examples)
    * [FizzBuzz](#fizzbuzz)
    * [Factorial](#factorial)
    * [Bubble Sort](#bubble-sort)
    * [Fibonacci](#fibonacci)
        
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
import "main.tsp"
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


## For loop

### Loop 100 times
```pascal
0 for dup 100 <=  do
    dup print
    inc
end
```

### Break
```pascal
0 for dup 100 <= do
    dup print
    if dup 10 == do
        break
    end
end
```

## Error handling
```python

# errors
# StackIndexError
# TypeError
# ImportError
# NameError

try 
    dup
except StackIndexError do
    "Error..." print
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

## Swap
```pascal
1 2 swap print print
```
'swap' swaps two values in stack.

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
{} # push empty list

"Hello World!" append # append string "Hello World!"

34 append # append int 34

68 append # append int 68

"T# Programming Language" 1 replace # replace list index '1' to string "T# Programming Language"

print

{ 1 2 3 4 5 6 7 8 9 10 } print
```

## Examples
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

### Factorial
```pascal
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

### Fibonacci
```pascal
10000 -> n

0 1 for over n < do
  over puts " " puts
  swap over +
end
drop drop

"" print
```


