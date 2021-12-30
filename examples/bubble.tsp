block BubbleSort do
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
end

block Main do
    [] 19 append 13 append 6  append 2  append 18 append 8 append 1 append dup -> arr

    len -> length

    "before:      " puts arr print

    call BubbleSort

    "sorted list: " puts print
end

call Main
