# JSON Difference Calculator

This project calculates the difference between two JSON files with 1.0 being exactly similar
and 0.0 being completely different

## Getting Started

Make sure you have Go installed. Then, if you have GOROOT="/usr/local/go",
you can do the following commands in a terminal from any directory:

```
git clone https://github.com/cadelaney3/jsondiff
cd jsondiff
go run main.go <file1.json> <file2.json>
```
You can also run from the terminal:

```
go get github.com/cadelaney3/jsondiff
cd $GOPATH/src/github.com/cadelaney3/jsondiff
go run main.go <file1.json> <file2.json>
```

### Program Info

The scoring is based on the number of keys that are the same, and whether or not those keys
have the same relationship to other keys in the JSON. The following example would be scored
as 0.2:

```
{
    "ex1": {
        "sub1": "kjdl"
    }
}

{
    "ex2": {
        "sub1": "iojkl"
    }
}
```

The data count for the first json is "ex1" count + "ex1" parent count + "ex1" child count +
"sub1" count + "sub1" parent count + ... = 5 ("kjdl" is not a key, so does not appear in data map).
Since data count for second json is also 5, and number of equal items between the two files is 2,
2 / (5+5) = 0.2.

This program assumes values are the same if they represent the same value in string format.
For example, 31 and "31" are equal, as well as the boolean value false and the string "false".

The program can also be further improved if keys were examined to be roughly equal, such as
"animals" being rougly equal to "critters" and thus be scored as roughly equal. It would also be an
improvement to use word stems for this purpose, because "beers" and "beer-list" can represent the same
thing and so looking for "beer" in these keys could help determine if they are the same. I left these 
improvements out of the program because of the extra time it would take to implement and the because it
would have extra computational cost.
