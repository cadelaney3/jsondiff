package main

import (
	"fmt"
	"quiq/pkg/jsondiff"
	"os"
	"io/ioutil"
	"encoding/json"
)

func main() {
	if len(os.Args[1:]) != 2 {
		fmt.Println("Error: Invalid number of arguments")
		os.Exit(1)
	}
	files, err := readFiles(os.Args[1:])
	if err != nil {
		panic(err)
	}

	if !json.Valid(files[0]) || !json.Valid(files[1]) {
		fmt.Println("Error: Invalid JSON in one or both files")
		os.Exit(1)
	}

	jsonFiles, err := jsondiff.LoadJSON(files)
	if err != nil {
		panic(err)
	}

	var f1, f2 jsondiff.File
	err = json.Unmarshal(files[0], &f1)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(files[1], &f2)
	if err != nil {
		panic(err)
	}

	fmt.Print("Score ranges between 0.0 and 1.0. Lower scores indicate greater difference \n" +
	"(0.0 completely different) and higher scores greater similarity (1.0 exactly equal).\n\n")

	if jsondiff.CheckExactEquals(files[0], files[1]) {
		fmt.Printf("Score: %.1f\nFiles are exactly the same\n", 1.0)
		return
	}

	if jsondiff.CheckDeepEquals(jsonFiles[0], jsonFiles[1]) {
		fmt.Printf("Score: %.2f\nSame content, different order\n", 0.99)
		return
	}

	data1 := jsondiff.FindKeys(f1)
	data2 := jsondiff.FindKeys(f2)

	result := jsondiff.Compare(data1, data2)

	fmt.Println("Score: ", result)
}

func readFiles(paths []string) ([][]byte, error) {
	files := make([][]byte, 2)
	if len(paths) != 2 {
		return nil, fmt.Errorf("Please enter exactly two files")
	}
	for i, file := range paths {
		f, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("Error reading file %s: %v", file, err)
		}
		files[i] = f
	}
	return files, nil
}