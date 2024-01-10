package main

// remove non C code from response

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

func makeInoFromResponse(str string) string {

	// generate .ino arduino file from ChatGPT response
	// first we parse for valid Ccode, then we save this text as a .iso

	char := "#include" // assume this is first line of C

	index := strings.Index(str, char)
	fmt.Println(index)
	if index == -1 {
		fmt.Printf("Character '%s' not found in the string.\n", char)
		index = 1
		return "bad"
	}
	trimmedStr := str[index:]

	char = "```" // C code is wrapped in this code identifier like md format
	index2 := strings.Index(trimmedStr, char)

	if index2 > 0 {
		trimmedStr = trimmedStr[:index2]
	}

	f, err := os.Create("duinoCode//duinoCode.ino")
	if err != nil {
		fmt.Println(err)
		return "bad"
	}
	_, err = f.WriteString(trimmedStr)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return "bad"
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return "bad"
	}

	return trimmedStr

}

func createNewPromptFromBadCode(trimmedStr string, duinoCompErr string) string {
	errorchar := "error:" // C code is wrapped in this code identifier like md format

	// now fine which line contains an error message::
	errorMsg := ""
	numberOfErrors := 0
	for _, line := range strings.Split(strings.TrimRight(duinoCompErr, "\n"), "\n") {
		index := strings.Index(line, errorchar)
		if index != -1 {
			// log.Printf("line= %v", line)
			indexError := strings.Index(line, errorchar)
			if indexError != -1 {
				errorMsg = errorMsg + "-" + line[indexError+len(errorchar):] // TODO: should use a map here
				numberOfErrors = numberOfErrors + 1
			}
		}
	}
	log.Printf("errorMsg= %v", errorMsg)
	log.Printf("numberOfErrors= %v", numberOfErrors)

	errorPrompt1 := "Here is some C code for an Arduino using the FastLED library to create light patterns for an LED strip: \n ``` \n"

	errorPrompt2 := ""
	if numberOfErrors == 1 {
		errorPrompt2 = "\n ``` \n The code does not compile. The compiler creates the following error message:\n"
	} else {
		errorPrompt2 = "\n ``` \n The code does not compile. The compiler creates the following " + string(numberOfErrors) + " error messages: \n"
	}

	// insert errorMsg

	errorPrompt3 := "\n Please use this error message to modify the provided code so the modified code compiles. "

	errorPrompt4 := "Please provide a response in C code. After the code please briefly explain the decision for your response.  \n "

	errorPrompt5 := `
Remember: 1. Always format the code in code blocks. Begin and end the new C code with the code block symbol of 3 backticks.
2. Do not leave unimplemented code blocks in your response. 
3. The only allowed library is fastLED. Do not import or use any other library.
4. If you are not sure what value to use, just use your best judge. Do not use None for anything.
`

	newErrorPromptStr := errorPrompt1 + trimmedStr + errorPrompt2 + errorMsg + errorPrompt3 + errorPrompt4 + errorPrompt5

	f, err := os.Create("newErrorPromptStr.txt")
	if err != nil {
		fmt.Println(err)
		return "bad"
	}

	_, err = f.WriteString(newErrorPromptStr)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return "bad"
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return "bad"
	}

	return newErrorPromptStr
}

func main() {

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

	fn := "I would like a monochrome black and white interpretation of the song Age of Aquarius.txt"
	fn = "responseHistory//" + fn
	response, err := os.ReadFile(fn) // saved response from chatgpt
	if err != nil {
		fmt.Print(err)
	}
	str = string(response)                 // convert content to a 'string'
	trimmedStr := makeInoFromResponse(str) // saves new .ino

	if trimmedStr == "bad" {
		fmt.Println("problem saving .ino")
	} else {
		fmt.Println("saved new .ino")
	}

	boardFlagType := 0

	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		// now compile and program the arduino board.
		// the compile step is useful to reveal problems with the code.
		// select which arduino board are using:
		if boardFlagType == 0 {
			log.Printf("board is uno \n")
			// c = exec.Command("arduino-cli.exe", "compile", "--fqbn", "arduino:avr:uno", "duinoCode\\duinoCode.ino")
			c = exec.Command("arduino-cli.exe", "compile", "--fqbn", "arduino:avr:uno", "duinoCode\\duinoCode.ino", "-v")

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
			log.Printf("we have an error in the response- Duino code does not compile!")
			// ask ChatGPT to modify the bad code to try and generate good code,
			// according to the error message from Arduino-CLI:
			duinoCompErr := stderr.String()
			newErrorPromptStr := createNewPromptFromBadCode(trimmedStr, duinoCompErr)
			// fmt.Println(newErrorPromptStr)

			req_postErr := openai.ChatCompletionRequest{
				// Model: openai.GPT3Dot5Turbo,   // select GPT3.5: not so great
				Model: openai.GPT4, // GPT4 is default
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: newErrorPromptStr,
					},
				},
			}

			fmt.Println("sending new prompt to correct error")
			timeA := time.Now()
			resp3, err := client.CreateChatCompletion(context.Background(), req_postErr)
			if err != nil {
				fmt.Printf("ChatCompletion error: %v\n", err)
			}
			// elapsed := timeA - time.Now()

			elapsedTime := time.Now().Sub(timeA)

			fmt.Println("ChatGPT response received- hopefully with error corrected")
			log.Printf("elapsed= %v", elapsedTime)
			responseContentPostErr := resp3.Choices[0].Message.Content
			ErrStr := string(responseContentPostErr) // convert content to a 'string'

			str = string(ErrStr)                   // convert content to a 'string'
			trimmedStr := makeInoFromResponse(str) // saves new .ino

			if trimmedStr == "bad" {
				fmt.Println("problem saving error-corrected .ino")
			} else {
				fmt.Println("saved new error-corrected .ino")
			}

			// now compile and program the arduino board with error-corrected .ino

			// select which arduino board are using:
			if boardFlagType == 0 {
				log.Printf("board is uno \n")
				// c = exec.Command("arduino-cli.exe", "compile", "--fqbn", "arduino:avr:uno", "duinoCode\\duinoCode.ino")
				c = exec.Command("arduino-cli.exe", "compile", "--fqbn", "arduino:avr:uno", "duinoCode\\duinoCode.ino", "-v")

			} else {
				log.Printf("board is atmega328  \n")
				c = exec.Command("arduino-cli.exe", "compile", "--fqbn", "arduino:avr:diecimila:cpu=atmega328", "duinoCode\\duinoCode.ino")
			}

			fmt.Println(c)

			var out bytes.Buffer
			var stderr bytes.Buffer
			c.Stdout = &out
			c.Stderr = &stderr

			err = c.Run()
			if err != nil {
				log.Printf("we have another error in the response- Duino code does not compile!")
				// ask ChatGPT to modify the bad code to try and generate good code,
				// according to the error message from Arduino-CLI:
				// duinoCompErr := stderr.String()
				os.Exit(3) // peace out
			}

			fmt.Println("new error-corrected code compiles!")

		} else {
			// fmt.Println("program the board")

		}

		fmt.Println("program the board")

		os.Exit(3) // peace out

	} // if windows

}
