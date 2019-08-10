# JSON Difference Calculator

This project calculates the difference between two JSON files with 1.0 being exactly similar
and 0.0 being completely different

## Getting Started

Make sure you have Go installed. Then, if you have GOROOT="/usr/local/go",
you can do the following commands in a terminal from any directory:

```
git clone https://github.com/cadelaney3/jsondiff
cd jsondiff
go run <file1.json> <file2.json>
```
You can also run from the terminal:

```
go get github.com/cadelaney3/jsondiff
cd $GOPATH/src/github.com/cadelaney3/jsondiff
go run <file1.json> <file2.json>
```