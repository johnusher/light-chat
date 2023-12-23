package main

// 1.read in mp3
// 2. convert speech to text using online chatgpt whisper api
// 3. create new prompt to change lights according to speech desire
// 4. send new prompt to chatgpt and generate arduinno code

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func main() {

	log.Printf(" starting ...")

	secretAPI, err := os.ReadFile("..//..//chatgpt.txt") // load my secret chatgpt apiKey
	if err != nil {
		fmt.Print(err)
		os.Exit(3)
	}
	str := string(secretAPI) // convert content to a 'string'
	// fmt.Println(str) // print my secret chatgpt apiKey
	apiKey := str

	// load example duno code for prompt:

	// duinoExamplePrompt, err := os.ReadFile("xmasTwinkleDuino.c")
	duinoExamplePrompt, err := os.ReadFile("duinoCode//duinoCode.ino")

	if err != nil {
		fmt.Print(err)
		os.Exit(3)

	}

	duinoExamplePromptStr := string(duinoExamplePrompt)

	// prompt := {map[string]interface{}{"role": "user", "content": "Here is an example C code for an Arduino using the FastLED library to create a christmas effect on an LED strip: ."},
	// 	map[string]interface{}{"role": "user", "content": duinoExamplePromptStr},
	// 	map[string]interface{}{"role": "user", "content": "please update the code to make it look like disco lights flashing at 120 beats-per-minute. Please providing a response in C code."}}

	client := openai.NewClient(apiKey)
	ctx := context.Background()
	//--------------------
	// 1.read in mp3
	// audioFn := "audio//moreMoody.mp3"
	audioFn := "audio//blue.mp3"
	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: audioFn,
	}

	// 2. convert speech to text using online chatgpt whisper api
	resp, err := client.CreateTranscription(ctx, req)
	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return
	}
	fmt.Println(resp.Text)

	newPrompt := resp.Text
	//--------------------

	req2 := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "Here is an example C code for an Arduino using the FastLED library to create a christmas effect on an LED strip:",
				// Content: "My name is Bob.",
			},
		},
	}

	// duinoExamplePromptStr = "what is my name?"

	req2.Messages = append(req2.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: duinoExamplePromptStr,
	})

	// duinoExamplePromptStr = "please update the code to make it look like disco lights flashing at 120 beats-per-minute."
	duinoExamplePromptStr = "please update the code to satisfy the following request:"

	req2.Messages = append(req2.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: duinoExamplePromptStr,
	})

	req2.Messages = append(req2.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: newPrompt,
	})

	duinoExamplePromptStr = "Please provide a response in C code"
	req2.Messages = append(req2.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: duinoExamplePromptStr,
	})

	// log.Printf("prompt : %+v\n", req2.Messages) // how do i print this??
	// var req2 openai.ChatCompletionRequest
	// look at https://pkg.go.dev/github.com/sashabaranov/go-openai@v1.17.9#section-readme

	resp2, err := client.CreateChatCompletion(context.Background(), req2)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		// continue
	}

	log.Printf(" request ok...")

	//fmt.Println(resp)
	// fmt.Println(resp.Choices[0].Message.Content)

	responseContent := resp2.Choices[0].Message.Content
	str = string(responseContent) // convert content to a 'string'

	fmt.Println(str)

	char := "#include <FastLED.h>" // assume this is first line of C!

	index := strings.Index(str, char)
	fmt.Println(index)
	if index == -1 {
		fmt.Printf("Character '%s' not found in the string.\n", char)
		return
	}
	trimmedStr := str[index:]

	// fmt.Println(trimmedStr)
	// fmt.Printf("trimmedStr: %v\n", trimmedStr)

	char = "```"

	index2 := strings.Index(trimmedStr, char)

	// fmt.Println(index2)

	if index2 > 0 {
		trimmedStr = trimmedStr[:index2]
		fmt.Printf("index2: %v\n", index2)
	}

	// fmt.Printf("trimmedStr: %v\n", trimmedStr)

	f, err := os.Create("duinoCode//duinoCode.ino")
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(trimmedStr)
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

	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.Command("arduino-cli.exe", "compile", "--fqbn", "arduino:avr:diecimila:cpu=atmega328", "duinoCode\\duinoCode.ino")
		fmt.Println(c)

		var out bytes.Buffer
		var stderr bytes.Buffer
		c.Stdout = &out
		c.Stderr = &stderr

		err := c.Run()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return
		}
		fmt.Println("Result: " + out.String())

		c = exec.Command("arduino-cli.exe", "upload", "-p", "COM5", "-b", "arduino:avr:diecimila:cpu=atmega328", "duinoCode\\duinoCode.ino", "-v")

		// var out bytes.Buffer
		// var stderr bytes.Buffer
		c.Stdout = &out
		c.Stderr = &stderr

		fmt.Println(c)

		err = c.Run()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return
		}
		fmt.Println("Result: " + out.String())
		// if err := c.Run(); err != nil {
		// 	fmt.Println("Error: ", err)
		// }

		os.Exit(3)

	}

}
