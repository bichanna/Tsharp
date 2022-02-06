# T# Documentation

## Introduction
T# is a dynamic typed stack-oriented programming language designed for building software.

It's similar to Forth, Porth but written in Go.

`WARNING! THIS LANGUAGE IS A WORK IN PROGRESS! ANYTHING CAN CHANGE AT ANY MOMENT WITHOUT ANY NOTICE!`

## Build
```shell
$ git clone https://github.com/Tsharp-lang/Tsharp
$ cd tsharp
$ go build main.go
$ ./main examples/main.tsp
or
$ ./main.exe examples/main.tsp
```

## Hello World!
```
"Hello World!\n" print
```
`print` print the element on top of the stack and remove it from the stack.

## Built-in Words
| name | stack | description |
| ---- | --------- | ----------- |
| `dup` | `a -- a a` | duplicate an element on top of the stack. |
| `drop` | ` a --  ` | drops the top element of the stack. |
| `swap` | `a b -- b a` | swap 2 elements on the top of the stack. |
| `print` | `a -- ` | print the element on top of the stack and remove it from the stack. |
| `println` | `a -- ` | `print` with a new line. |
| `rot` | ` a b c -- b c a ` | rotate the top three stack elements. |
| `over` | ` a b -- a b a ` | duplicate the second value on the stack. |
| `input` | ` -- <input value> ` | user input. |
| `exit` | ` -- ` | exit |
| `free` | ` a b c -- ` | drop all elements of the stack. |
| `isdigit` | ` <string value> -- <bool value> ` | check the top string type element is digit. push the bool value. |

## Arithmetic Operators
```
34 35 + println
```
push `34` and `35` to the stack.
`+` plus two elements on the stack and push it back to the stack.
`println` will print the top element on the stack.

## Comments
```python
# comment...
```

## Data-types
```
int      # 1 2 3 4
string   # "Hello World!"
bool     # true false
list     # { 1 2 3 4 }
error    # NameError...
type     # int string bool list...
```

## Block
```
block Main do
    "Hello World!" println
end

Main
```
In T# a block is defined using the `block` keyword.

## Variables
```
"Hello World!" -> N
N println
```

## Variable scope
```
10 -> N # Global variable

block Main do
    N println
    100 -> A # `A` can be only used in the Main block.

    if true do
        A println
    end

    # `i` can be only used in the Main block.
    0 for dup 2 < do -> i
        i println
        i inc
    end
end

Main

N println

A println # error

i println # error
```

## If statements
```
if true do
    "Hello World!" println
elif true do
    "Hello World!" println
else
    "Hello World!" println
end
```

```
2 2 == println
2 3 != println
2 3 < println
3 2 > println
2 3 <= println
3 2 >= println
```

```
11 -> N

N { 20 30 11 42 28 91 } in
```

## For loop
```
for true do
     "Hello World!" println
end
```

```
0
for dup 100 < do
    dup println
    1 +
end
```
The code above is a program that outputs from 1 to 100.<br>
First, push `0` onto the stack.<br>
`dup` duplicates the top element of the stack.<br>
Now there are two `0` in the stack. ( 0  0 )<br>
push `100` onto the stack ( 0  0  100 )<br>
`<` pushes a `bool` type to the stack using the top and second element of the stack. 0 < 100.<br>
The stack will look like this. ( 0  true )<br>
`do` checks for `true` or `false`. If `true`, it will run the for statement.<br>
The used `bool` type will be removed from the stack.<br>
If it becomes `false`, the loop will stop.<br>

```
0
for dup 100 < do
    -> i
    i println
    i 1 +
end
```
By the way, this is how I write the loop process.

## List
```
{ 1 2 3 4 5 6 7 8 9 10 } println
```

### Append
```
{ 1 2 3 4 5 6 7 8 9 10 } 11 append

# <list> <index> append
```

### Read
```
{ 1 2 3 4 5 6 7 8 9 10 } 0 read println

# <list> <index> read
# or
# <string> <index> read
```

### Replace
```
{ 1 2 3 4 5 6 7 8 9 10 } "Hello World!" 0 replace println

# <list> <replace value> <index> replace
```

### Remove
```
{ 1 2 3 4 5 6 7 8 9 10 } 0 remove println

# <list> <index> remove
```

### Len
```
{ 1 2 3 4 5 6 7 8 9 10 } len println

# <list> len
# or
# <string> len
```

## File operations
```
"main.asm" fopen -> F

"; Hello World!\n" F fwrite

F fread -> context

context println

F fclose
```

## Error handling
```
try
    println
except StackUnderflowError do
    # do something...
end
```

`StackUnderflowError` When you try to use a value from the stack when there is nothing on the stack.<br>
`TypeError` When the type is different.<br>
`IndexError` When the indexes of arrays and strings are different.<br>
`IncludeError` When you try to include an invalid file.<br>
`NameError` When you use a variable that does not exist.<br>
`AssertionError` Assertion.<br>
`FileNotFoundError` file not found.<br>

## Assertion
```
false assert "assertion error message..."
```

## Include
```
include "main.tsp"
```

## Built in T#
### tic tac toe game 
<a href="https://github.com/Tsharp-lang/tictactoe"><img src="https://github-readme-stats.vercel.app/api/pin/?username=Tsharp-lang&repo=tictactoe"/></a>

## T# highlighter
https://twitter.com/m0k1m0k1 created this!

https://marketplace.visualstudio.com/items?itemName=akamurasaki.tsharplanguage-color
