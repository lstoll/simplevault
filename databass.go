package main

import "encoding/json"

func databassDecode(j []byte) (map[string]string, error) {
	var data map[string]string
	if err := json.Unmarshal(j, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func databassEncode(d map[string]string) ([]byte, error) {
	data, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return data, nil
}
