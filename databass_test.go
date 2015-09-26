package main

import (
	"reflect"
	"testing"
)

var (
	testJson = []byte(`{"key1":"val1","key2":"val2"}`)
	testMap  = map[string]string{"key1": "val1", "key2": "val2"}
)

func TestDatabassDecode(t *testing.T) {
	data, err := databassDecode(testJson)
	if err != nil {
		t.Errorf("%v", err)
	}
	if !reflect.DeepEqual(data, testMap) {
		t.Errorf("expected %v got %v", testMap, data)
	}
}

func TestDatabassEncode(t *testing.T) {
	data, err := databassEncode(testMap)
	if err != nil {
		t.Errorf("%v", err)
	}
	if string(data) != string(testJson) {
		t.Errorf("expected `%s` got `%s`", testJson, data)
	}
}
