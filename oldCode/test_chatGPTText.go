package main

// not this!

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

// 1 token ~= 4 chars in English
// 1 token ~= ¾ words
// 100 tokens ~= 75 words

const (
	apiEndpoint = "https://api.openai.com/v1/chat/completions"
)

func mainxxxx() {

	secretAPI, err := os.ReadFile("..//..//chatgpt.txt") // load my secret chatgpt apiKey
	if err != nil {
		fmt.Print(err)
	}
	str := string(secretAPI) // convert content to a 'string'
	// fmt.Println(str) // print my secret chatgpt apiKey
	apiKey := str

	// load example duno code for prompt:

	duinoExamplePrompt, err := os.ReadFile("duinoCode//xmasTwinkleDuino.c")
	if err != nil {
		fmt.Print(err)
	}

	duinoExamplePromptStr := string(duinoExamplePrompt)
	// fmt.Println(duinoExamplePromptStr)

	// os.Exit(3)

	client := resty.New()

	// messages:=[
	//     {"role": "system", "content": "You are a helpful assistant."},
	//     {"role": "user", "content": "Who won the world series in 2020?"},
	//     {"role": "assistant", "content": "The Los Angeles Dodgers won the World Series in 2020."},
	//     {"role": "user", "content": "Where was it played?"}
	// ]

	// prompt :=
	// 	"messages": []interface{}{map[string]interface{}{"role": "system", "content": "Please tell me a joke about a dog and a cat."}},
	// "messages": []interface{}{map[string]interface{}{"role": "system", "content": "You are an assistant that speaks like Shakespeare."},
	// 	map[string]interface{}{"role": "user", "content": "tell me a joke"}},

	// prompt := []interface{}{map[string]interface{}{"role": "system", "content": "You are an assistant that speaks like Shakespeare."},
	// 	map[string]interface{}{"role": "user", "content": "tell me a joke"}}

	prompt := []interface{}{map[string]interface{}{"role": "user", "content": "Here is an example C code for an Arduino using the FastLED library to create a christmas effect on an LED strip: ."},
		map[string]interface{}{"role": "user", "content": duinoExamplePromptStr},
		map[string]interface{}{"role": "user", "content": "please update the code to make it look like disco lights flashing at 120 beats-per-minute. Please providing a response in C code."}}

	response, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			// "model": "gpt-3.5-turbo",
			"model":    "gpt-4",
			"messages": prompt,
			// "max_tokens": 50,
			// "max_tokens": null,
		}).
		Post(apiEndpoint)

	if err != nil {
		log.Fatalf("Error while sending send the request: %v", err)
	}

	body := response.Body()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error while decoding JSON response:", err)
		return
	}

	log.Printf("sent request ...")
	fmt.Println("sent request")

	// Extract the content from the JSON response
	content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)

	fmt.Println(content)

	f, err := os.Create("response.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(content)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

}
