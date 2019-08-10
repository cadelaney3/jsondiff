package jsondiff

import (
	"bytes"
	"encoding/json"
	"math"
	"strconv"
	"reflect"
	"strings"
)

// File defines an arbitrary JSON file
type File map[string]interface{}

// Key is a struct that keeps track of properties of a JSON key
type Key struct {
	Name string // name of key
	Children map[string]int // maps child names and counts
	Parent map[string]int // maps parent names and counts
	Count int // counts how many times the key shows up
}

// LoadJSON loads a a slice of byte slices into a File type and returns a slice of Files
func LoadJSON(files [][]byte) ([]File, error) {
	fileSlice := make([]File, len(files))

	for i, file := range files {
		json.Unmarshal(file, &fileSlice[i])
	}
	return fileSlice, nil
}

// FindKeys finds all the keys in a JSON and stores the data for each key in a Key struct.
// Returns a map of Key structs containing appropriate data
func FindKeys(file File) map[string]*Key {

	// init data map
	data := make(map[string]*Key)

	// initial File contains keys, so loop thru initial key, value pairs
	for k, v := range file {
		// convert to lower case to make future comparisons easier
		k = strings.ToLower(k)
		// initialize map key
		data[k] = &Key{ Name: k, Count: 1 }
		// go to lower level of JSON
		findKeysHelper(data, k, v)
	}

	return data
}

// findKeysHelper uses recursion to go thru keys of JSON and store the data for each key
func findKeysHelper(data map[string]*Key, parent string, jsonItem interface{}) {
	// initial children map of parent item
	if data[parent].Children == nil {
		data[parent].Children = make(map[string]int)
	}

	switch jsonItem.(type) {
	case []interface{}:
		for _, item := range jsonItem.([]interface{}) {
			findKeysHelper(data, parent, item)
		}
		break

	case map[string]interface{}:
		for k, v := range jsonItem.(map[string]interface{}) {
			k = strings.ToLower(k)

			// increment this child for the parent key
			data[parent].Children[k]++

			// if this key does not exist, initialize
			if _, ok := data[k]; !ok {
				data[k] = &Key{ Name: k }
			}

			if data[k].Parent == nil {
				data[k].Parent = make(map[string]int)
			}

			// increment the parent for this key
			data[k].Parent[parent]++

			// increment times this key has appeared
			data[k].Count++

			// recursive call to go down a level in JSON
			findKeysHelper(data, k, v)
		}
		
	// these are for final values (have no lower levels). Add these values as children to parent keys
	case string:
		temp := strings.ToLower(jsonItem.(string))
		if _, ok := data[temp]; ok {
			data[temp].Parent[parent]++
			data[temp].Count++
		}
		data[parent].Children[temp]++
		break

	case bool:
		temp := strings.ToLower(strconv.FormatBool(jsonItem.(bool)))
		if _, ok := data[temp]; ok {
			data[temp].Parent[parent]++
			data[temp].Count++
		}
		data[parent].Children[temp]++
		break

	case int:
		if _, ok := data[strconv.Itoa(jsonItem.(int))]; ok {
			data[strconv.Itoa(jsonItem.(int))].Parent[parent]++
		}
		data[parent].Children[strconv.Itoa(jsonItem.(int))]++
		break

	case float64:
		if _, ok := data[strconv.FormatFloat(jsonItem.(float64), 'f', -1, 32)]; ok {
			data[strconv.FormatFloat(jsonItem.(float64), 'f', -1, 32)].Parent[parent]++
		}
		data[parent].Children[strconv.FormatFloat(jsonItem.(float64), 'f', -1, 32)]++
		break

	case nil:
		if _, ok := data["null"]; ok {
			data["null"].Parent[parent]++
		}
		data[parent].Children["null"]++
		break

	default:
		break
	}
}

// Compare returns a score that is calculated by diving the total number of keys and final values of each JSON
// by the number of items that are equal in the two JSON. Then the two values are added together and divided by
// two to get an average between 0 and 1.
func Compare(dat1 map[string]*Key, dat2 map[string]*Key) float64 {
	dat1Count := 0
	dat2Count := 0
	sameCount := 0.0
	for k, v := range dat1 {

		// get the data counts from each count in key, parent map, and child map
		dat1Count += v.Count
		for _, val := range v.Parent {
			dat1Count += val
		}
		for _, val := range v.Children {
			dat1Count += val
		}

		// check if key exists in both data maps
		if v2, ok := dat2[k]; ok {

			sameCount += math.Min(float64(v.Count), float64(v2.Count))

			for key, val := range v.Parent {
				if _, ok := v2.Parent[key]; ok {
					// add min of the parents that are the same to sameCount
					sameCount += math.Min(float64(val), float64(v2.Parent[key]))
				}
			}

			for key, val := range v.Children {
				if _, ok := v2.Children[key]; ok {
					// add min of the children that are the same to sameCount
					sameCount += math.Min(float64(val), float64(v2.Children[key]))
				}
			}
		}
	}

	// get the data count for the other json
	for _, v := range dat2 {
		dat2Count += v.Count
		for _, val := range v.Parent {
			dat2Count += val
		}
		for _, val := range v.Children {
			dat2Count += val
		}
	}

	avg1 := sameCount / float64(dat1Count)
	avg2 := sameCount / float64(dat2Count)
	return (avg1+avg2) / 2
}


// CheckExactEquals checks if the two files are exactly identicle
func CheckExactEquals(file1, file2 []byte) bool {
	return bytes.Equal(file1, file2)
}

// CheckDeepEquals checks if the two files contains the same data.
// It will return true if the files have the same data but in a 
// different order
func CheckDeepEquals(file1, file2 File) bool {
	return reflect.DeepEqual(file1, file2)
}