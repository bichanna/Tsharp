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
                whiteblock
            else
                blackblock
            end
        inc end drop
        " \n" print
    inc end drop
end

natural

"\n" print

block horizonal do
    0 for dup 10 < do
        dup -> i
        0 for dup 10 < do
            dup -> j
            if j 2 % 0 == do
                whiteblock
            else
                blackblock
            end
        inc end drop
        " \n" print
    inc end drop
end

horizonal

"\n" print

block vertical do
    0 for dup 10 < do
        dup -> i
        0 for dup 10 < do
            dup -> j
            if i 2 % 0 == do
                whiteblock
            else
                blackblock
            end
        inc end drop
        " \n" print
    inc end drop
end

vertical

"\n" print

block cross do
    0 for dup 10 < do
        dup -> i
        0 for dup 10 < do
            dup -> j
            if i j + 2 % 0 == do
                whiteblock
            else
                blackblock
            end
        inc end drop
        " \n" print
    inc end drop
end

cross

"\n" print

block triangle do
    10 -> y
    y y 1 - + -> x
    0 for dup y < do
        dup -> i
        0 for dup x < do
            dup -> j
            if j 1 + y i - < do
                whiteblock
            elif x j - y i - < do
                whiteblock
            else
                blackblock
            end
        inc end drop
        " \n" print
    inc end drop
end

triangle

"\n" print

block inversetriangle do
    10 -> y
    y y 1 - + -> x
    0 for dup y < do
        dup -> i
        0 for dup x < do
            dup -> j
            if j i < do
                whiteblock
            elif x j - i > do
                blackblock
            else
                whiteblock
            end
        inc end drop
        " \n" print
    inc end drop
end

inversetriangle
