package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

func main() {
	user := User{}

	// Marshal the user struct to JSON.
	json, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	// Print the JSON output.
	fmt.Println(string(json))
}
