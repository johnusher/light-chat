package main

// light_chat_main.go

// Chat-GPT enabled control of a strip of programmable LEDs using a text or voice prompt.
// For speech to text we use the OpenAI Whisper API.
// To program the LEDs we use an Arduino.
// The user prompt will generate new Arduino code to modify the existing light pattern or to generate a new pattern.

// Useage:
// go run .\light_chat_main.go
// flags:
// -p "make lights kinda blue"= text input (default is use microphone)
// -r (bool) reset reference light pattern to xmasTwinkleDuino.c
// -nb (bool) = no-board: (no Arduino connected). A response is still generated in

// Using a voice or text prompt, we program a strip of Arduino-controlled LEDs.

// 1. Obtain desired new or ammended light pattern - use text or ...
//    ... listen to microphone and convert speech to text using online chatgpt whisper api
// 3. create new prompt to replace or update light pattern according.
// 4. Send new prompt to chatgpt
// 5. Receive ChatGPT response, parse for valid Arduino C code.
// 6. Compile new light pattern and program Arduino (using Ardunio-CLI).

import (
	"C"
	"bufio"
	"bytes"
	"context"
	"flag"
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
)

// For Windows OS, we use winmm to record microphone signal
var (
	winmm         = syscall.MustLoadDLL("winmm.dll")
	mciSendString = winmm.MustFindProc("mciSendStringW")
)

// This sets up response for key to complete voice prompt:
var reader = bufio.NewReader(os.Stdin)

func readKey(input chan rune) {
	char, _, err := reader.ReadRune()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)
	input <- char
}

func MCIWorker(lpstrCommand string, lpstrReturnString string, uReturnLength int, hwndCallback int) uintptr {
	i, _, _ := mciSendString.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpstrCommand))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpstrReturnString))),
		uintptr(uReturnLength), uintptr(hwndCallback))
	return i
}

func main() {

	log.Printf(" starting light-chat")

	// setup optional command-line flags:
	resetChat := flag.Bool("r", false, "reset chat history to xmas example")
	textPromptFlag := flag.String("p", "mic", "use text or mic input to generate prompt")
	noBoardFlag := flag.Bool("nb", false, "false if we want to run without connected board")

	flag.Parse()

	if *textPromptFlag != "mic" {
		log.Printf(" using text input to generate prompt\n")
	} else {
		log.Printf(" use speech prompt for prompt \n")
	}

	comPort := "" // USB port to which Arduino is connected
	boardFlagType := 0
	index := -1

	if (runtime.GOOS != "windows") && (!*noBoardFlag) {
		log.Printf("only tested with windows OS! \n")
	}

	if *noBoardFlag {
		log.Printf("running without board attached \n")
	}

	if (runtime.GOOS == "windows") && (!*noBoardFlag) {
		// find the USB COM port the duino is attached to
		// NB this is not robust!! manually check using
		// arduino-cli.exe board list
		boards, err := exec.Command("arduino-cli.exe", "board", "list").Output()
		if err != nil {
			log.Printf(" COM problem \n")
			log.Fatal(err)
		}

		boardsS := string(boards)
		// log.Printf("boards: \n %v \n", boardsS)

		charS := "Arduino"
		index = strings.Index(boardsS, charS)

		// check if a duino is connected:
		if index == -1 {
			fmt.Printf("Error- Arduino not connected!\n", charS)
			os.Exit(3)
		}

		// if so, check if it is a Uno type board:
		ifUno := strings.Index(boardsS, "Uno")

		if ifUno == -1 {
			fmt.Printf("Uno not connected- assume atmega328 \n")
			boardFlagType = 1
		}

		// now fine which USB COM port duino is connected with:
		for _, line := range strings.Split(strings.TrimRight(boardsS, "\n"), "\n") {
			index = strings.Index(line, charS)
			if index != -1 {
				comPort = line[0:4]
			}
		}

		log.Printf("Duino connected with USB: %v \n", comPort)
	}

	var duinoExamplePrompt []byte
	var err error

	//  load example duino code for prompt:

	if *resetChat {
		log.Printf("reset flag used: using xmasTwinkleDuino.c")
		// reset flag used- use the default c code as a reference for the prompt
		duinoExamplePrompt, err = os.ReadFile("xmasTwinkleDuino.c")

	} else {
		// modify the existing light program:
		duinoExamplePrompt, err = os.ReadFile("duinoCode//duinoCode.ino")
		log.Printf("  reset not used. using duinoCode.ino as reference code.")
	}

	if err != nil {
		fmt.Print(err)
		os.Exit(3)
	}

	newPrompt := ""
	duinoExamplePromptStr := string(duinoExamplePrompt)

	// now setup ChatGPT-API connection:

	secretAPI, err := os.ReadFile("..//..//chatgpt.txt") // load your secret chatgpt apiKey
	if err != nil {
		fmt.Print(err)
		os.Exit(3)
	}
	str := string(secretAPI)
	apiKey := str

	if err != nil {
		fmt.Print(err)
		os.Exit(3)
	}

	client := openai.NewClient(apiKey)
	ctx := context.Background()

	if *textPromptFlag == "mic" {

		// we are using speech input with the microphone
		// only works on windows OS!

		fmt.Println("Use winmm.dll to record Audio to .wav file")

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

		// i = MCIWorker("set capture bytespersec 32000", "", 0, 0)  // doesnt work?
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

		// following wait method is unreliable for saving wav file:

		input := make(chan rune, 1)
		fmt.Println("recording mic \n")
		fmt.Println("hit keyboard return to finish \n")

		// wait for enter key to finish microphone recording
		go readKey(input)
		select {
		case <-input:

		case <-time.After(18000 * time.Millisecond):
			fmt.Println("Time out!")
		}

		fmt.Println("hit key")
		n := 1
		for n < 850 { // try many times to save file
			i = MCIWorker("save capture mic.wav", "", 0, 0)
			if i != 0 {
				// log.Fatal("Error Code C: ", i)
				n = n + 1
				time.Sleep(50 * time.Millisecond)
			} else {
				break
			}
		}
		fmt.Printf("n= %v \n", n)
		if n == 850 {
			log.Fatal("Error Code C: ", i)
		}

		i = MCIWorker("close capture", "", 0, 0)
		if i != 0 {
			log.Fatal("Error Code D: ", i)
		}

		fmt.Println("Audio saved to mic.wav")

		//--------------------
		// 1.read in audio file (can be wav or mp3)
		// audioFn := "audio//moreMoody.mp3"
		// audioFn := "audio//kafka.mp3"
		audioFn := "mic.wav"
		req := openai.AudioRequest{
			Model:    openai.Whisper1,
			FilePath: audioFn,
		}

		// convert speech in audio file to text using online chatgpt whisper api:
		resp, err := client.CreateTranscription(ctx, req)
		if err != nil {
			fmt.Printf("Transcription error: %v\n", err)
			return
		}
		fmt.Println(resp.Text)

		newPrompt = resp.Text

	} else {
		// use text input not speech:
		newPrompt = *textPromptFlag

		if newPrompt[len(newPrompt)-1] != '.' {
			newPrompt = newPrompt + "." // append fullstop
		}
		// fmt.Println(newPrompt)
	}

	newPromptStr := "" // create user string prompt

	//--------------------

	var promptIntro string
	if *resetChat {
		promptIntro = "Here is an example C code for an Arduino using the FastLED library to create a christmas effect on an LED strip. This is the existing pattern: \n ``` \n "
	} else {
		promptIntro = "Here is an example C code for an Arduino using the FastLED library to create light patterns for an LED strip.  This is the existing pattern: \n ``` \n"
	}

	req2 := openai.ChatCompletionRequest{
		// Model: openai.GPT3Dot5Turbo,   // select GPT3.5: not so great
		Model: openai.GPT4, // GPT4 is default
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: promptIntro,
			},
		},
	}
	newPromptStr = promptIntro

	req2.Messages = append(req2.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: duinoExamplePromptStr, // add example code
	})

	newPromptStr = newPromptStr + duinoExamplePromptStr

	// duinoExamplePromptStr = "please update the code to make it look like disco lights flashing at 120 beats-per-minute."   // sanity check
	duinoExamplePromptStr = "\n ``` \n Please generate new C code to satisfy the following request and output the modified code in its entirety beginning with #include <FastLED.h> and include the line: #define COLOR_ORDER GRB. "

	req2.Messages = append(req2.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: duinoExamplePromptStr,
	})

	newPromptStr = newPromptStr + duinoExamplePromptStr

	req2.Messages = append(req2.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: newPrompt,
	})

	duinoExamplePromptStr = "Please provide a response in C code. After the code please briefly explain the decision for your response. Here is the request: \n "
	req2.Messages = append(req2.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: duinoExamplePromptStr,
	})

	// Response rules:
	duinoExamplePromptStrFinal := `
	Remember: 1. Always format the code in code blocks.
	2. Do not leave unimplemented code blocks in your response. Begin and end the new C code with the code block symbol of 3 backticks.
	3. The only allowed library is fastLED. Do not import or use any other library.
	4. If you are not sure what value to use, just use your best judge. Do not use None for anything.
	5. Only if you can not provide a valid response, please flash all LEDs red twice (do this just once) and then display the existing pattern.
	`

	req2.Messages = append(req2.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: duinoExamplePromptStrFinal,
	})

	newPromptStr = newPromptStr + duinoExamplePromptStr + newPrompt + duinoExamplePromptStrFinal

	f, err := os.Create("newPrompt.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = f.WriteString(newPromptStr)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// look at https://pkg.go.dev/github.com/sashabaranov/go-openai@v1.17.9#section-readme

	resp2, err := client.CreateChatCompletion(context.Background(), req2)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
	}

	log.Printf(" chatGPT response received")

	//fmt.Println(resp)
	// fmt.Println(resp.Choices[0].Message.Content)

	responseContent := resp2.Choices[0].Message.Content
	str = string(responseContent) // convert content to a 'string'

	responseSn := newPrompt
	responseSn = strings.Replace(responseSn, ".", "-", -1)
	if len(responseSn) > 75 {
		responseSn = responseSn[:75]
	}
	responseSn = responseSn + ".txt"
	responseSn = "responseHistory" + "//" + responseSn

	f, err = os.Create(responseSn)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(str)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	char := "#include" // assume this is first line of C!

	index = strings.Index(str, char)
	fmt.Println(index)
	if index == -1 {
		fmt.Printf("Character '%s' not found in the string.\n", char)
		index = 1
		return
	}
	trimmedStr := str[index:]

	char = "```" // C code is wrapped in this code identifier like md format
	index2 := strings.Index(trimmedStr, char)

	if index2 > 0 {
		trimmedStr = trimmedStr[:index2]
	}

	// fmt.Printf("trimmedStr: %v\n", trimmedStr)

	f, err = os.Create("duinoCode//duinoCode.ino")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(trimmedStr)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// now compile and program the arduino board.
	// the compile step is useful to reveal problems with the code.
	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		// select which arduino board are using:
		if boardFlagType == 0 {
			log.Printf("board is uno \n")
			c = exec.Command("arduino-cli.exe", "compile", "--fqbn", "arduino:avr:uno", "duinoCode\\duinoCode.ino")
		} else {
			log.Printf("board is atmega328  \n")
			c = exec.Command("arduino-cli.exe", "compile", "--fqbn", "arduino:avr:diecimila:cpu=atmega328", "duinoCode\\duinoCode.ino")
		}

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

		if boardFlagType == 0 {
			c = exec.Command("arduino-cli.exe", "upload", "-p", comPort, "-b", "arduino:avr:uno", "duinoCode\\duinoCode.ino", "-v")
		} else {
			c = exec.Command("arduino-cli.exe", "upload", "-p", comPort, "-b", "arduino:avr:diecimila:cpu=atmega328", "duinoCode\\duinoCode.ino", "-v")
		}

		c.Stdout = &out
		c.Stderr = &stderr

		fmt.Println(c)

		err = c.Run()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return
		}
		fmt.Println("Result: " + out.String())
		os.Exit(3) // peace out

	}

}
