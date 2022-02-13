<div align="center">
    <br>
    <img width="300px" src="https://user-images.githubusercontent.com/81926489/150945785-a4b40a2c-e68b-4bf8-b68c-009cc33985ba.PNG">
    <h1> The T# Programming Language</h1>


[Docs](https://github.com/Tsharp-lang/Tsharp/blob/main/DOC/docs.md) |
[Docs(日本語)](https://github.com/Tsharp-lang/Tsharp/blob/main/DOC/ドキュメント.md) |
[Vim](https://github.com/Tsharp-lang/Tsharp/blob/main/editor/tsharp.vim) | 
[VSCode](https://marketplace.visualstudio.com/items?itemName=akamurasaki.tsharplanguage-color)
</div>

<div align="center"> 

[![Ubuntu](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-ubuntu.yml/badge.svg)](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-ubuntu.yml)
[![Windows](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-windows.yml/badge.svg)](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-windows.yml)
[![CodeQL](https://github.com/Tsharp-lang/Tsharp/actions/workflows/codeql-analysis.yml/badge.svg?branch=main)](https://github.com/Tsharp-lang/Tsharp/actions/workflows/codeql-analysis.yml)

</div>
WARNING! THIS LANGUAGE IS A WORK IN PROGRESS! ANYTHING CAN CHANGE AT ANY MOMENT WITHOUT ANY NOTICE!
<br><br>
It's like Forth, but written in Go ʕ◔ϖ◔ʔ.
<a href="https://en.wikipedia.org/wiki/Stack-oriented_programming">Stack-oriented programming</a>

### Install (build) && Run

> Install
```
go build main.go
```

> Run
```
$ ./main <filename>.tsp

or

$ ./main.exe <filename>.tsp
```

# Hello World
```pascal
"Hello World" println
```
```Crystal
block main do
    {"game" "web" "tools" "science" "systems" "embedded" "drivers" "GUI" "mobile"} -> areas

  0 for dup areas len < do -> i
        "Hello, " print areas i read print " developers!\n" print
        i inc
    end
end

main
```


# Sudoku Solver

Sudoku solver in T#
https://github.com/Tsharp-lang/SudokuSolver

<img width="1230" alt="Screen Shot 2022-02-13 at 13 18 13" src="https://user-images.githubusercontent.com/81926489/153741735-aac281fc-1b5d-4608-bca3-3745b1ff3dc6.png">

# Tic Tac Toe

tic tac toe written in T#
https://github.com/Tsharp-lang/tictactoe

![tictactoe](https://user-images.githubusercontent.com/81926489/150774403-dafdb578-ca0d-497a-b123-dcd1639654e8.gif)


### Contributing
Welcome! 💕
