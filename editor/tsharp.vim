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
syntax keyword tsharpKeywords import block do end if elif else for int string bool type

" Comments
syntax region tsharpCommentLine start="//" end="$"   contains=tsharpTodos
           
" Strings
syntax region tsharpString start=/\v"/ skip=/\v\\./ end=/\v"/
syntax region tsharpString start=/\v'/ skip=/\v\\./ end=/\v'/

" Set highlights
highlight default link tsharpKeywords Identifier
highlight default link tsharpCommentLine Comment
highlight default link tsharpString String

let b:current_syntax = "tsharp"