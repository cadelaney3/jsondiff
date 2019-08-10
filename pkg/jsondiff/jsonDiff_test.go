package jsondiff

import (
	"encoding/json"
	"testing"
	"reflect"
)


func TestCheckExactEquals(t *testing.T) {
	goodJSON := `{"example": 1}`
	badJSON := `{"example":2:]}}`
	goodJSON2 := `{"example: "1"}`
	t.Run("good & bad", testCheckExactEquals([]byte(goodJSON), []byte(badJSON), false))
	t.Run("good & good", testCheckExactEquals([]byte(goodJSON), []byte(goodJSON2), false))
}

func testCheckExactEquals(b1, b2 []byte, expected bool) func (t *testing.T) {
	return func (t *testing.T) {
		got := CheckExactEquals(b1, b2)
		if got != expected {
			t.Errorf("CheckExactEquals(b1, b2) = %v, want %v", got, expected)
		}
	}
}

func TestCheckDeepEquals(t *testing.T) {
	var JSON1 File
	var JSON2 File
	var JSON3 File

	_ = json.Unmarshal([]byte(`{"example": 1, "test": "two"}`), &JSON1)
	_ = json.Unmarshal([]byte(`{"example: "1", "test": "two"}`), &JSON2)
	_ = json.Unmarshal([]byte(`{"test": "two", "example": 1}`), &JSON3)

	t.Run("JSON1 & JSON2", testCheckDeepEquals(JSON1, JSON2, false))
	t.Run("JSON1 & JSON3", testCheckDeepEquals(JSON1, JSON3, true))

}

func testCheckDeepEquals(f1, f2 File, expected bool) func (t *testing.T) {
	return func (t *testing.T) {
		got := CheckDeepEquals(f1, f2)
		if got != expected {
			t.Errorf("CheckDeepEquals(f1, f2) = %v, want %v", got, expected)
		}	
	}
}

func TestFindKeys(t *testing.T) {
	var JSON File
	_ = json.Unmarshal([]byte(`{"example": 3, "parent": {"child1": 1, "child2": [2, 4]}}`), &JSON)

	expected := map[string]*Key{"example": &Key{Name: "example", 
		Children: map[string]int{"3": 1}, Count: 1}, "parent": &Key{Name: "parent",
		Children: map[string]int{"child1": 1, "child2": 1}, Count: 1}, "child1": &Key{Name: "child1",
		Parent: map[string]int{"parent": 1}, Children: map[string]int{"1": 1}, Count: 1}, "child2": &Key{Name: "child2",
		Parent: map[string]int{"parent": 1}, Children: map[string]int{"2": 1, "4": 1}, Count: 1}}

	got := FindKeys(JSON)

	for k, v := range got {
		if !reflect.DeepEqual(v, expected[k]) {
			t.Errorf("FindKeys(JSON) = %v, want %v", v, expected[k])
		}
	}
}

func TestCompare(t *testing.T) {
	var JSON File
	_ = json.Unmarshal([]byte(`{"example": 3, "parent": {"child1": 1, "child2": [2, 4]}}`), &JSON)

	data1 := map[string]*Key{"example": &Key{Name: "example", 
		Children: map[string]int{"3": 1}, Count: 1}, "parent": &Key{Name: "parent",
		Children: map[string]int{"child1": 1, "child2": 1}, Count: 1}, "child1": &Key{Name: "child1",
		Parent: map[string]int{"parent": 1}, Children: map[string]int{"1": 1}, Count: 1}, "child2": &Key{Name: "child2",
		Parent: map[string]int{"parent": 1}, Children: map[string]int{"2": 1, "4": 1}, Count: 1}}

	data2 := FindKeys(JSON)

	result := Compare(data1, data2)
	if result < 1.0 {
		t.Errorf("Compare(data1, data2) = %f, want == 1.0", result)
	}
}