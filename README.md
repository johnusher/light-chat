# light-chat
* Light-chat is a Chat-GPT enabled system to control a strip of programmable LEDs using a text or voice prompt. For speech to text we use the OpenAI Whisper API. To program the LEDs we use an Arduino. The user prompt will generate new Arduino code to modify the existing light pattern or to generate a new pattern.

* We sometimes get invalid code from the ChatGPT response. We use the response from the Arduino compiler ```arduino-cli compile ``` to extract error messages and resend the invalid code and the error message to ChatGPT to get new and (hopefully) valid code.<br />

* This is a Golang and Arduino project using the Go openAI interface library [https://github.com/sashabaranov/go-openai](https://github.com/sashabaranov/go-openai) and the FastLED library [https://github.com/FastLED](https://github.com/FastLED).




# Useage
The main file is [light_chat_main.go](light_chat_main.go)

It determines if an Arduino is connected via USB, what COM port it is connected to, and what board type it is (atmega328 or Uno).

Command line execution:  <br /> 
```
go run .\light_chat_main.go
``` 
<br /> 

flags:<br /> 
-p "make lights kinda blue"= text prompt input. If we do not use this p flag then we default to use voice input with the default microphone.<br /> 
-r (bool) reset the reference (provided example) light pattern code to [xmasTwinkleDuino.c](xmasTwinkleDuino.c). This file contains information about the programable LED strip chipset (currently WS2812B), the number of LEDs on the strip, and which Arduino pin to use for the data output.<br /> 
If the flag is not used, the reference (provided example) light pattern code is [duinoCode.ino](duinoCode/duinoCode.ino).<br /> 
In summary: Use this -r flag when you want to reset the light pattern and make a new one (e.g. "make a new calm light pattern that is a poem by John Donne"). Otherwise do not use the flag to modify whatever the exisiting pattern is (e.g. "make the light pattern less romantic"). <br /> 
 -nb (bool) = no board. Use this flag if no Arduino is connected but you still want to generate a response that is saved in directory [responseHistory](responseHistory).  <br /> 

# Examples

### 1. Barbie's dream house

[![Watch the video]![Alt text](video/bdh.png)](https://youtube.com/shorts/Lj5bpjEMJQU?feature=shared)

(In the video, the prompt was repeated just before the Arduino updated the lights.)

Response:<br /> 

[Barbie's dream house](<responseHistory/Please it Barbie's dream house-.txt>)


```
Here is a modified code that displays a Barbie's Dream House themed light pattern by alternating between pink, cyan, and white lights, which imitates the traditional color scheme of Barbie. ++ It briefly flashes all LEDs red twice at the initial setup as requested ++ . 
```
Note this is an incorrect interpretation of rule #5:
```
	5. Only if you can not provide a valid response, please flash all LEDs red twice (do this just once) and then display the existing pattern.

```



### 2. 420 Party.
```
go run .\light_chat_main.go -p "Make the colours suitable for a party 
at 420 wink-wink nudge-nudge say no more" -nb -r 
```
This uses a text prompt input, with no Arduino board connected, and using the default reference code in xmasTwinkleDuino.c.

 The ChatGPT response is stored in directory [responseHistory](responseHistory) as a text file with filename the same as the user prompt for the light pattern -  [example for above](<responseHistory/Make the colours suitable for an adult 4-20 party nudge-nudge wink-wink say.txt>).<br /> 

 We add "response rules" to ask ChatGPT to describe how it interpretted the prompt. For the above example, this is the comment it generated on the code:<br /> 
```Your ChristmasLights function was replaced by the PartyLights function. The colours were changed from a repeating pattern of red, green, white to green and purple. For a 4-20 themed party these colours are more suitable as green often represents marijuana and purple represents the feeling of relaxation and royalty.```

### 3. More Romantic.
 ``` 
 go run .\light_chat_main.go -p "Now make the lights more romantic." -nb
 ```
This uses a text prompt input without the -r flag to modify the code generated by the previous response in the 1st example.<br /> 
This is the ChatGPT comment on the code:<br /> 

```
 The above code will create a "romantic" atmosphere by decreasing the brightness to a soft glow (80 out of 255). It also slows down the frame rate to a calm 2 frames per second. Most importantly, it sets all the LEDs to red (CHSV(0, 255, 255) is a hue-saturation-value representation of the color red), which is traditionally associated with romance.
 ```

### 4. Age of Aquarius.
```
I would like a monochrome black and white interpretation of the song Age of Aquarius
```

Response: [here](<responseHistory/I would like a monochrome black and white interpretation of the song Age of Aquarius.txt>)<br /> 
The C code did not compile first time- it used an undefined parameter HALF. But the response did not lack in ambition nor inventive interpretation :)

the function ```createNewPromptFromBadCode``` uses ChatGPT to modify such non-compilable "bad" code: the function takes the error message from the Arduino compiler and sends this to ChatGPT with the "bad" code, and asks it to correct the bad code so it compiles. If we again get "bad" code, we abort.


### 5. ... it's kinda complicated...
```Please make the lights display a single blue dot that travels down the light strip starting slow and speeding up as it travels. When the blue dot reaches the end turn it to a single yellow light that reflects and returns along the light strip slowly and flashing. When it reaches the end flash all lights bright white 4 times and then repeat this pattern but with the colours becoming more and more random.```

This took 1.4 minutes to receive a response and the code compiled first time. 



# How it works
1. Obtain desired new or ammended light pattern using a text or microphone recording and convert speech to text using the OpenAI Whisper api.
3. Create a new prompt to replace or update the light pattern accordingly. <br /> 
We concat the following text:<br /> 
    "Here is an example of C code for an Arduino using the FastLED library to create light patterns for an LED strip.  This is the existing pattern:"<br /> 
    &lt; insert C code from example reference &gt; <br /> 
    " Please generate new C code to satisfy the following request and output the modified code in its entirety: "<br /> 
    &lt; user prompt &gt; <br /> 
    "These are the rules:"<br /> 
    &lt; 1. Always format the code in code blocks. etc...&gt;<br /> 
4. Send this new prompt to ChatGPT.
5. Receive ChatGPT response and parse for valid C code.
6. Compile new light pattern.
7. If the code does not compile, there is a function ```createNewPromptFromBadCode``` to use ChatGPT to amend such bad code : the function takes the error message from the Arduino compiler and sends this to ChatGPT with the "bad" code and we prompt ChatGPT to correct the bad code so it compiles. If we again get "bad" code, we abort.
8. We program the Arduino using the valid code from step 6/7.

# Setup the software environment

Download: <br /> 
golang (tested with 1.21.5 on Windows)

Arduino CLI

Install the go-openai library (provides unofficial Go clients for OpenAI API):

```
go get github.com/sashabaranov/go-openai
```


## Install Arduino libraries
Used library Version Path<br />
FastLED      3.6.0   C:\Users\jusher\Documents\Arduino\libraries\FastLED<br />

Used platform Version Path<br />
arduino:avr   1.8.6   C:\Users\jusher\AppData\Local\Arduino15\packages\arduino\hardware\avr\1.8.6<br />


# ChatGPT

* Pay $ for ChatGPT-API credit and store your private key somewhere eg ...//privateChatGPTAPIkey.txt <br />
* We use ChatGPT4 because it gives better code responses than ChatGPT3.5 Turbo (less errors in the C code). For a typical prompt+response it costs 1-4 cents. 

* Install Arduino libraries
Used library Version Path<br />
FastLED      3.6.0   C:\Users\jusher\Documents\Arduino\libraries\FastLED<br />

* Used platform Version Path<br />
arduino:avr   1.8.6   C:\Users\jusher\AppData\Local\Arduino15\packages\arduino\hardware\avr\1.8.6<br />

# Useful files
[duinoHelp.txt](duinoCode/duinoHelp.md)


# Problems, TODOs, comments


* It take typically 30-50 seconds to receive the response using ChatGPT4.
* Some colours look bad with RGB LEDS- e.g. brown and yellow, and white if no white LED.<br />
* It costs 1-4 cents for each prompt using ChatGPT4. The WhisperAPI is <0.1C per ~10 second prompt.<br />
* The winmm.dll is pretty janky and about 1 in 10 times the microphone recording will not save: need to switch to using core-audio (for Linux compatability).


# Credits

* Nimrod Gileadi for the idea of including example code in the response and encouragement:  [https://github.com/nimrod-gileadi](https://github.com/nimrod-gileadi) <br /> 

Some of the prompt engineering suffix (the "response rules") were taken from a paper he worked on (Wenhao Yu, Nimrod Gileadi, et al. Language to rewards for robotic skill synthesis. arXiv preprint arXiv:2306.08647, 2023).<br />

* Siggy for golang help: [https://sig.gy/](https://sig.gy/) <br />

* The Arduino project (try to buy at least one legit product!) :[https://www.arduino.cc/](https://www.arduino.cc/) <br />  



