block whiteblock do
    "■ " print
end

block blackblock do
    "□ " print
end

block natural do
    0 for dup 10 < do
        dup -> i
        0 for dup 10 < do
            dup -> j
            if i j < do
                call whiteblock
            else
                call blackblock
            end
        inc end drop
        " \n" print
    inc end drop
end

call natural

"\n" print

block horizonal do
    0 for dup 10 < do
        dup -> i
        0 for dup 10 < do
            dup -> j
            if j 2 % 0 == do
                call whiteblock
            else
                call blackblock
            end
        inc end drop
        " \n" print
    inc end drop
end

call horizonal

"\n" print

block vertical do
    0 for dup 10 < do
        dup -> i
        0 for dup 10 < do
            dup -> j
            if i 2 % 0 == do
                call whiteblock
            else
                call blackblock
            end
        inc end drop
        " \n" print
    inc end drop
end

call vertical

"\n" print

block cross do
    0 for dup 10 < do
        dup -> i
        0 for dup 10 < do
            dup -> j
            if i j + 2 % 0 == do
                call whiteblock
            else
                call blackblock
            end
        inc end drop
        " \n" print
    inc end drop
end

call cross

"\n" print

block triangle do
    10 -> y
    y y 1 - + -> x
    0 for dup y < do
        dup -> i
        0 for dup x < do
            dup -> j
            if j 1 + y i - < do
                call whiteblock
            elif x j - y i - < do
                call whiteblock
            else
                call blackblock
            end
        inc end drop
        " \n" print
    inc end drop
end

call triangle

"\n" print

block inversetriangle do
    10 -> y
    y y 1 - + -> x
    0 for dup y < do
        dup -> i
        0 for dup x < do
            dup -> j
            if j i < do
                call whiteblock
            elif x j - i > do
                call blackblock
            else
                call whiteblock
            end
        inc end drop
        " \n" print
    inc end drop
end

call inversetriangle
