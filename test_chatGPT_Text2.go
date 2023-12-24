package main

// 1.read in mp3
// 2. convert speech to text using online chatgpt whisper api
// 3. create new prompt to change lights according to speech desire
// 4. send new prompt to chatgpt and generate arduinno code

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"

	openai "github.com/sashabaranov/go-openai"
	"github.com/schollz/progressbar/v3"
)

var (
	winmm         = syscall.MustLoadDLL("winmm.dll")
	mciSendString = winmm.MustFindProc("mciSendStringW")
)

var reader = bufio.NewReader(os.Stdin)

func readKey(input chan rune) {
	char, _, err := reader.ReadRune()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)
	input <- char
}

// func Rename(old_path, new_path string) error

func MCIWorker(lpstrCommand string, lpstrReturnString string, uReturnLength int, hwndCallback int) uintptr {
	i, _, _ := mciSendString.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpstrCommand))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpstrReturnString))),
		uintptr(uReturnLength), uintptr(hwndCallback))
	return i
}

func main() {

	log.Printf(" starting lightChat3000")

	// load example duno code for prompt:

	duinoExamplePrompt, err := os.ReadFile("xmasTwinkleDuino.c")
	// duinoExamplePrompt, err := os.ReadFile("duinoCode//duinoCode.ino")

	// Rename and Remove a file
	// Using Rename() function
	// Original_Path := "mic.wav"
	// New_Path := "old_mic.wav.txt"
	// e := os.Rename(Original_Path, New_Path)
	// if e != nil {
	// 	fmt.Printf("e : %v\n", e)
	// }

	fmt.Println("winmm.dll Record Audio to .wav file")

	i := MCIWorker("open new type waveaudio alias capture", "", 0, 0)
	if i != 0 {
		log.Fatal("Error Code A: ", i)
	}

	i = MCIWorker("set capture bitspersample 16", "", 0, 0)
	if i != 0 {
		log.Fatal("Error Code a1: ", i)
	}

	i = MCIWorker("set capture channels 1", "", 0, 0)
	if i != 0 {
		log.Fatal("Error Code a2: ", i)
	}

	// i = MCIWorker("set capture bytespersec 32000", "", 0, 0)
	// if i != 0 {
	// 	log.Fatal("Error Code a3: ", i)
	// }

	i = MCIWorker("set capture alignment 2", "", 0, 0)
	if i != 0 {
		log.Fatal("Error Code a3: ", i)
	}

	i = MCIWorker("record capture", "", 0, 0)
	if i != 0 {
		log.Fatal("Error Code B: ", i)
	}

	// time.Sleep(200 * time.Millisecond)

	// fmt.Println("Listening...")

	// time.Sleep(5 * time.Second)

	bar := progressbar.Default(25)
	for i := 0; i < 25; i++ {
		bar.Add(1)
		time.Sleep(200 * time.Millisecond)
	}

	// I want to replace the above sleep command with below wait for keyboard but it is unreliable!

	// input := make(chan rune, 1)
	// var iru rune
	// fmt.Println("recording + waiting keyboard input...")
	// go readKey(input)
	// select {
	// case iru = <-input:
	// 	// fmt.Printf("Input : %v\n", i)

	// 	// close(input)

	// case <-time.After(5000 * time.Millisecond):
	// 	fmt.Println("Time out!")
	// 	// i = MCIWorker("save capture mic.wav", "", 0, 0)
	// 	// if i != 0 {
	// 	// 	log.Fatal("Error Code C: ", i)
	// 	// }

	// 	// i = MCIWorker("close capture", "", 0, 0)
	// 	// if i != 0 {
	// 	// 	log.Fatal("Error Code D: ", i)
	// 	// }

	// 	// fmt.Println("Audio saved to mic.wav")

	// }

	// // fmt.Printf("Input : %v\n", iru)
	// close(input)

	i = MCIWorker("save capture mic.wav", "", 0, 0)
	if i != 0 {
		log.Fatal("Error Code C: ", i)
	}

	i = MCIWorker("close capture", "", 0, 0)
	if i != 0 {
		log.Fatal("Error Code D: ", i)
	}

	fmt.Println("Audio saved to mic.wav")

	// fmt.Println("Time out!")
	// i = MCIWorker("save capture mic.wav", "", 0, 0)
	// if i != 0 {
	// 	log.Fatal("Error Code C: ", i)
	// }

	// i = MCIWorker("save capture mic2.wav", "", 0, 0)
	// if i != 0 {
	// 	log.Fatal("Error Code C: ", i)
	// }

	// os.Exit(3)
	///////////////////////////////////////

	secretAPI, err := os.ReadFile("..//..//chatgpt.txt") // load my secret chatgpt apiKey
	if err != nil {
		fmt.Print(err)
		os.Exit(3)
	}
	str := string(secretAPI) // convert content to a 'string'
	// fmt.Println(str) // print my secret chatgpt apiKey
	apiKey := str

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
	// audioFn := "audio//kafka.mp3"
	audioFn := "mic.wav"
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
		// Model: openai.GPT4,
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
	duinoExamplePromptStr = "please generate new code to satisfy the following request and output the modified code in its entirety beginning with #include <FastLED.h> and include the line: #define COLOR_ORDER GRB"

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

	char := "#include" // assume this is first line of C!

	index := strings.Index(str, char)
	fmt.Println(index)
	if index == -1 {
		fmt.Printf("Character '%s' not found in the string.\n", char)

		index = 1

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
