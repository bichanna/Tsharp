# T# Documentation

## Introduction

T# is a Stack-based programming language designed for building software.

It's like Forth but written in Go ʕ◔ϖ◔ʔ.

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

* [Variables](#variables)
    * [Variable Scopes](#variable-scopes)

* [Data types](#data-types)
    * [Typeof](#typeof)

* [Arithmetic](#arithmetic)
* [Block](#block)

* [List](#list)
    * [Append](#append)
    * [Read](#read)
    * [Remove](#remove)
    * [Len](#len)

* [If Statement](#if-statement)

* [For loop](#for-loop)
    * [Loop 100 times](#loop-100-times)
    * [Break](#break)
 
* [Error handling](#error-handling)
* [Assertions](#assertions)
* [Import](#import)

</td><td width=33% valign=top>

* [Built-in Words](#built-in-words)
    * [Dup](#dup)
    * [Drop](#drop)
    * [Over](#over)
    * [Rot](#rot)
    * [Swap](#swap)
    * [Inc](#inc)
    * [Dec](#dec)
    * [Input](#input)
    * [Exit](#exit)
    * [PrintS](#prints)

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
`print` will print the top element of the stack, then remove it.

## Comments
```python
# comment...
```

## Variables
```pascal
10 -> x

x -> y

y print
```
### Variable Scopes
```crystal
10 -> N # Global variable

block Main do
    N print
    100 -> A # `A` can only be used within the Main function.

    if true do
        A print
    end

    # `i` can only be used within the Main function.
    0 for dup 2 < do -> i
        i print
        i inc
    end
end

call Main

try
    A print # error
except NameError do
    "`A` var error" print
end

try
    i print # error
except NameError do
    "`i` var error" print
end
```

## Data types
```python
int # 12345
string # "Hello World!"
bool # true false
list # { 1 2 3 4 5 }
type # int string bool type
error # NameError, StackIndexError...
```

### Typeof
```python
"Hello World" typeof print
```
`typeof` push the type of the top element of the stack. The element that `typeof` used will be dropped.

## Arithmetic
```pascal
34 35 + print

100 40 - print

200 5 / print

10 2 * print
```

## Block
```pascal
block main do
    "Hello World!" print
end

call main
```

`block` function in other languages.

## List
```python
{ 1 2 3 4 5 6 7 8 9 10 } print
```
### Append
```python
{ 1 2 3 4 5 6 7 8 9 10 } 11 append
```
### Read
```python
{ 1 2 3 4 5 6 7 8 9 10 } 0 read print
```
### Replace
```python
{ 1 2 3 4 5 6 7 8 9 10 } "Hello World!" 0 replace print
```
### Remove
```python
{ 1 2 3 4 5 6 7 8 9 10 } 0 remove print
```
### Len
```python
{ 1 2 3 4 5 6 7 8 9 10 } len print
```

## If Statement
```pascal
if true do
    "Hello World" print
end
```
```pascal
if false do
    "Hello World" print
else
    "Hello John Doe" print
end
```
```pascal
if false do
    "Hello World" print
elif true
    "Hello John Doe" print
end
```
```pascal
if true false || do
    "Hello World" print
end
```
```pascal
if true true && do
    "Hello World" print
end
```
```pascal
2 2 == print
2 3 != print
2 3 < print
3 2 > print
2 3 <= print
3 2 >= print
```
```pascal
11 -> N

if N { 20 30 11 42 28 91 } in do
    "Hello World!" print
end
```

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

write loop like other languages
```
# T#                       |    // JavaScript
0 for dup 100 < do -> i    |    var i = 0
    i print                |    while (i < 100) {
    i inc -> i             |       console.log(i);
end                        |       i++;
                           |    }
```

## Error handling

`StackIndexError`
`TypeError`
`ImportError`
`NameError`
`AssertionError`
`IndexError`
```python
try 
    dup
except StackIndexError do
    "Error..." print
except IndexError do
    "Error..." print
except NameError do
    "Error..." print
end
```

## Assertions
```python
false assert "assertion error..."

# AssertionError:1:7: assertion error...
```

## Import
```python
import "main.tsp"
```

## Built-in Words

### Dup
```pascal
"Hello World" dup print print
```
`dup` duplicate element on top of the stack. ( a -- a a )

### Drop
```pascal
"Hello World" "T# Programming Language" drop print
```
`drop` drops the top element of the stack. ( a --  )

### Rot
```python
1 2 3 rot print print print
```
`rot` rotate top three stack elements. ( a b c -- b c a )

### Swap
```pascal
1 2 swap print print
```
`swap` swaps two values in stack. ( a b -- b a )

### Over
```python
1 2 over print print print
```
`over` copy the element below the top of the stack. ( a b -- a b a )

### Inc
```python
1 inc print
```
`inc` increment the top element of the stack

### Dec
```python
10 dec print
```
`dec` decrement the top element of the stack

### input
```python
input print
```

### Exit
```python
"Hello World"
exit
print
```
`exit` will exit the program.

### PrintS
```python
1 2 "Hello World!"

printS

# stack length  
#     ↓ 
#    <3>  1 2 'Hello World' <- top
```
`printS` print all stack values. 'printS' won't drop stack value after print.

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
```factor
{ 19 13 6 2 18 8 1 4 11 9 100 30 4 } -> arr arr print

arr len -> length

0 for dup length <= do
    0 for dup length 1 - < do
        dup -> j
        j 1 + -> i
        if arr j read arr i read > do
            arr j read -> x
            arr i read -> y
            arr y j replace dup -> arr
            x i replace -> arr
        end 
        inc
    end drop
    inc
end drop

arr print
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


