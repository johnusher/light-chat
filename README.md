# light-chat
Chat-GPT enabled control of a strip of Arduino-controlled programmable LEDs using a text or voice prompt (with the Whisper API) to modify the existing light pattern or to generate a new pattern.

This is a Golang and Arduino project using the Go openAI interface library [github.com/sashabaranov/go-openai](github.com/sashabaranov/go-openai) and the FastLED library [https://github.com/FastLED](https://github.com/FastLED).

The project has only been tested with Windows OS.



# Useage
The main file is [light_chat_main.go](light_chat_main.go)

The program determines if an Arduino is connected via USB, what COM port it is connected to, and what board type it is (atmega328 or Uno).


Command line execution:  <br /> 
go run .\light_chat_main.go <br /> 

flags:<br /> 
-p "make lights kinda blue"= text input (if we do not use this flag then we default to use the microphone input.)<br /> 
-r (bool) reset reference the light pattern to [xmasTwinkleDuino.c](xmasTwinkleDuino.c) <br /> 
 -nb (bool) = no board. Use this flag if no Arduino is connected but you still want to generate a response.  <br /> 

 The ChatGPT responses are stored in [responseHistory](responseHistory) as a text file with filename the same as the user request for the light pattern.


# How it works
1. Obtain desired new or ammended light pattern - use text or listen to microphone and convert speech to text using the chatgpt whisper api
3. Create a new prompt to replace or update the light pattern accordingly.
4. Send new prompt to chatgpt.
5. Receive ChatGPT response, parse for valid Arduino C code.
6. Compile new light pattern and program Arduino (using Ardunio-CLI). 

# Setup the software environment

Download: <br /> 
golang (tested with 1.21.5 on Windows)

Arduino CLI

Install the go-openai library (provides unofficial Go clients for OpenAI API):

```
go get github.com/sashabaranov/go-openai
```



# ChatGPT

Pay $ for ChatGPT-API credit and store your private key somewhere eg ...//privateChatGPTAPIkey.txt

# Install Arduino libraries
Used library Version Path<br />
FastLED      3.6.0   C:\Users\jusher\Documents\Arduino\libraries\FastLED<br />

Used platform Version Path<br />
arduino:avr   1.8.6   C:\Users\jusher\AppData\Local\Arduino15\packages\arduino\hardware\avr\1.8.6<br />

# Useful files
[duinoHelp.txt](duinoCode/duinoHelp.md)


# Problems, TODOs, comments

We sometimes get invalid code: Use response from ```arduino-cli compile ``` to extract error messages and resend this to ChatGPT to try and get revised and working code.<br />
Some colours look bad with RGB LEDS- e.g. brown.<br />
It cost about 0.6 cents for each prompt using ChatGPT4.<br />
The winmm.dll is pretty janky and about 1 in 10 times the microphone recording will not work: need to switch to using core-audio (for Linux compatability).


# Credits

Nim for the idea of including example code in the response and encouragement:  [https://github.com/nimrod-gileadi](https://github.com/nimrod-gileadi) <br /> 
Some of the prompt "streamlining" suffix code was taken from a paper he worked on (Wenhao Yu, Nimrod Gileadi, Chuyuan Fu, Sean Kirmani, Kuang-Huei Lee, Montse Gonzalez Arenas, Hao-Tien Lewis Chiang, Tom Erez, Leonard Hasenclever, Jan Humplik, et al. Language to
rewards for robotic skill synthesis. arXiv preprint arXiv:2306.08647, 2023).<br />


Siggy for golang help: [https://sig.gy/](https://sig.gy/) <br />
The Arduino project (try to buy at least one legit product!) :[https://www.arduino.cc/](https://www.arduino.cc/) <br />  



