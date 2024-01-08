This is an overview of the arduino cli tool:

https://arduino.github.io/arduino-cli/0.21/getting-started/


1. compile a C file (which for Arduino uses the .ino extension, not .c)<br />
cd to folder containing the .ino (this folder must have the same name as the .ino)<br />

./arduino-cli compile --fqbn arduino:avr:diecimila:cpu=atmega328 .\duinoCode.ino<br />
or <br />
arduino-cli.exe compile --fqbn arduino:avr:uno duinoCode.ino

update libraries: <br />
arduino-cli.exe  core update-index

2. upload .ino to the board: <br />
1. check what boards are connected and to which COM port: <br />
arduino-cli.exe board list
now upload: <br />
arduino-cli.exe upload -p COM5 --fqbn arduino:avr:diecimila:cpu=atmega328 duinoCode\\duinoCode.ino" -v <br />
arduino-cli.exe upload -p COM4 --fqbn arduino:avr:uno duinoCode.ino -v






----------------------


Connect the board to your PC¶
The first thing to do upon a fresh install is to update the local cache of available platforms and libraries by running:


$ arduino-cli core update-index

arduino-cli.exe board list
./arduino-cli board listall mkr

./arduino-cli.exe core install arduino:avr


./arduino-cli compile --fqbn arduino:avr:diecimila:cpu=atmega328 xmasTwinkleDuino.ino


./arduino-cli compile --fqbn arduino:avr:diecimila:cpu=atmega328 

Now verify we have installed the core properly by running:


./arduino-cli core list


in duino_src:

compile and flash:
Uno:
arduino-cli compile --fqbn arduino:avr:uno duino_src

arduino-cli upload -p /dev/ttyACM0 --fqbn arduino:avr:uno duino_src



./arduino-cli upload -p COM5 --fqbn arduino:avr:uno duino_src.ino

arduino-cli.exe upload -p COM5 --fqbn arduino:avr:uno 