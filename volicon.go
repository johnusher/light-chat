package main

// voice light control

// 1.read in mp3
// 2. convert speech to text using online chatgpt whisper api
// 3. create new prompt to change lights according to speech desire
// 4. send new prompt to chatgpt and generate arduinno code

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/johnusher/light-chat/pkg/keyboard"

	openai "github.com/sashabaranov/go-openai"
)

func main() {

	log.Printf(" starting ...")

	// Setup keyboard input:
	// stop := make(chan os.Signal, 1)
	// signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	keys := make(chan rune)

	kb, err := keyboard.Init(keys)
	if err != nil {
		// log.Errorf("failed to initialize keyboard: %s", err)
		fmt.Print(err)
		return
	}

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

	// prompt := {map[string]interface{}{"role": "user", "content": "Here is an example C code for an Arduino using the FastLED library to create a christmas effect on an LED strip: ."},
	// 	map[string]interface{}{"role": "user", "content": duinoExamplePromptStr},
	// 	map[string]interface{}{"role": "user", "content": "please update the code to make it look like disco lights flashing at 120 beats-per-minute. Please providing a response in C code."}}

	client := openai.NewClient(apiKey)
	ctx := context.Background()
	//--------------------
	// 1.read in mp3
	audioFn := "audio//moreMoody.mp3"

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

	f, err := os.Create("response2.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(responseContent)
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

	//-----------
	// go forth
	go kb.Run()

	errs := make(chan error)

	go func() {
		errs <- VoliLoop(keys)
	}()

	// block until ctrl-c or one of the loops returns an error
	select {
	case <-errs:
	}

}

func VoliLoop(keys <-chan rune) error {

	// more := false
	for {

		select {

		case _, more := <-keys:
			// received local key press
			// todo: replace/ augment this with a GPIO button press

			if !more {
				log.Printf("keyboard listener closed, exiting...")
				return nil
			}

		}
	}
}
