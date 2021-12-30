" Vim syntax file
" Language: Tsharp

" Usage Instructions
" Put this file in .vim/syntax/tsharp.vim
" and add in your .vimrc file the next line:
" autocmd BufRead,BufNewFile *.tsp set filetype=tsharp

if exists("b:current_syntax")
  finish
endif

" Language keywords
syntax keyword tsharpKeywords import block do end if elif else for try except break

" Type keywords
syntax keyword tsharpType int string type bool list

" Boolean keywords
syntax keyword tsharpBoolean true false

" Comments
syntax region tsharpCommentLine start="//" end="$"   contains=tsharpTodos
           
" Strings
syntax region tsharpString start=/\v"/ skip=/\v\\./ end=/\v"/
syntax region tsharpString start=/\v'/ skip=/\v\\./ end=/\v'/

" Numbers
syntax match tsharpNumbers '\d\+'

" Exceptions
syntax keyword tsharpExceptions try except

" Set highlights
highlight default link tsharpKeywords Repeat
highlight default link tsharpCommentLine Comment
highlight default link tsharpString String
highlight default link tsharpNumbers Number
highlight default link tsharpType Type
highlight default link tsharpBoolean Boolean
highlight default link tsharpExceptions Exception

let b:current_syntax = "tsharp"