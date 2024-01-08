package main

// remove non C code from response

import (
	"fmt"
	"os"
	"strings"
)

func mainx() {

	response, err := os.ReadFile("response.txt") // saved response from chatgpt
	if err != nil {
		fmt.Print(err)
	}
	str := string(response) // convert content to a 'string'
	// fmt.Println(str)        // print my secret chatgpt apiKey
	// apiKey := str

	char := "#include <FastLED.h>" // assume this is first line of C!

	index := strings.Index(str, char)
	fmt.Println(index)
	if index == -1 {
		fmt.Printf("Character '%s' not found in the string.\n", char)
		return
	}
	trimmedStr := str[index:]

	// fmt.Println(trimmedStr)

	char = "```"

	index2 := strings.Index(trimmedStr, char)

	trimmedStr = trimmedStr[:index2]

	// fmt.Println(index2)
	fmt.Println(trimmedStr)
}
