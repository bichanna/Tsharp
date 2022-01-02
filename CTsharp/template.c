#include <stdio.h>
#include <stdlib.h>

typedef struct EXPR_STRUCT {
    enum {
        EXPR_INT,
        EXPR_STRING,
    }type;
    char* string_value;
    int   int_value;
}EXPR_T;

typedef struct STACK_STRUCT
{
    EXPR_T** expr_stack;
    size_t  expr_stack_size;
} STACK_T;

STACK_T* init_stack()
{
    STACK_T* stack         = calloc(1, sizeof(struct STACK_STRUCT));
    stack->expr_stack      = (void*)0;
    stack->expr_stack_size = 0;
    return stack;
}

STACK_T* stack_push(STACK_T* stack, EXPR_T* expr)
{
    stack->expr_stack_size += 1;

    if (stack->expr_stack == (void*) 0)
    {
        stack->expr_stack = calloc(1, sizeof(struct EXPR_STRUCT*));
    }
    else
    {
        stack->expr_stack = 
            realloc(
                stack->expr_stack,
                stack->expr_stack_size * sizeof(struct EXPR_STRUCT**)
            );
    }

    stack->expr_stack[stack->expr_stack_size-1] = expr;
    return (void*)0;
}

STACK_T* stack_drop(STACK_T* stack)
{
    if (stack->expr_stack == (void*) 0)
    {
        printf("Error: empty stack.");
        exit(0);
    }
    else
    {
        free(stack->expr_stack[stack->expr_stack_size-1]);
    }

    stack->expr_stack_size -= 1;
    return (void*)0;
}

EXPR_T* init_expr()
{
    EXPR_T* expr       = calloc(1, sizeof(struct EXPR_STRUCT));
    expr->string_value = (void*)0;
    expr->int_value    = 0;
    return expr;
}

void clean_up(EXPR_T** exprs)
{
    int length = sizeof(exprs) / sizeof(struct EXPR_STRUCT);
    for (int i = 0; i < length; i ++) {
        free(exprs[i]);
    }
}

int main(int argc, char* argv[])
{
    STACK_T* stack = init_stack();
    EXPR_T* expr   = init_expr();
    // stack_push(stack, expr);
    // stack_drop(stack);
    // clean_up();
}

