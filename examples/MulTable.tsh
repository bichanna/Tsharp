# https://twitter.com/mio2561 created this program.

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
            " " puts
        end
        puts
        " " puts
        inc
    end
    "" print
    drop
    inc
end

