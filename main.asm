; this is for https://github.com/Tsharp-lang/TsharpTsharp

GLOBAL _main

section .text

print:
    sub     rsp, 40
    mov     BYTE [rsp+31], 0xA
    test    rdi, rdi
    jne     .L2
    mov     BYTE [rsp+30], 48
    mov     edx, 30
    mov     r8d, 2

.L3:
    lea     rsi, [rsp+rdx]
    mov     edi, 1
    mov     rdx, r8
    mov     rax, 0x2000004
    syscall
    add     rsp, 40
    ret

.L2:
    mov  r9, -3689348814741910323
    lea     rcx, [rsp+30]

.L4:
    mov     rax, rdi
    lea     r8, [rsp+32]
    mul     r9
    mov     rax, rdi
    sub     r8, rcx
    shr     rdx, 3
    lea     rsi, [rdx+rdx*4]
    add     rsi, rsi
    sub     rax, rsi
    add     eax, 48
    mov     BYTE [rcx], al
    mov     rax, rdi
    mov     rdi, rdx
    mov     rdx, rcx
    sub     rcx, 1
    cmp     rax, 9
    ja      .L4
    lea     rax, [rsp+32]
    sub     rdx, rax
    add     rdx, 32
    jmp     .L3

global start

_main:
    ; -- push --
    push    1000
    push    35
    push    20

    ; -- plus --
    ;pop    rax
    ;pop    rbx
    ;add    rbx, rax
    ;push   rbx

    ; -- minus --
    ;pop    rax
    ;pop    rbx
    ;sub    rbx, rax
    ;push   rbx

    ; -- mul --
    ;pop    rax
    ;pop    rbx
    ;mul    rbx, rax
    ;push   rbx

    ; -- div --
    ;pop    rax
    ;pop    rbx
    ;div    rbx, rax
    ;push   rbx

    ; -- print ---
    pop    rdi
    call   print

    ; -- print ---
    pop    rdi
    call   print

    ; -- print ---
    pop    rdi
    call   print

    mov    rax, 0x2000001 ; exit
    mov    rdi, 0
    syscall

section .data


