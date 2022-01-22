
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


