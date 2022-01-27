

block dclone do
    dup -> tmpa
    swap
    dup -> tmpb
    swap
    tmpb
    tmpa
end

1 for dup 10 < do
    1 for dup 10 < do
        call dclone
        *
        if dup 10 < do
            " " print
        end
        print
        " " print
        inc
    end
    "\n" print
    drop
    inc
end

