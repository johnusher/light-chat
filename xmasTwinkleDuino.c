#include <FastLED.h>

#define LED_PIN     6  // do not change
#define COLOR_ORDER GRB  // do not change
#define CHIPSET     WS2812B   // LED chipset. do not change
#define NUM_LEDS    30  // length of the light strip, also known as lightstrip and LED strip or LEDstrip. 
                        // NUM_LEDS is the number of LEDs in the LED strip. do not change this value.

#define BRIGHTNESS  150    // set brightness of all LEDs. 150 is default
#define FRAMES_PER_SECOND 5   // 20 is very fast, 10 is normal, 5 is slow. Above 30 is too fast.

CRGB leds[NUM_LEDS];

void setup() {
  delay(1000); // delay before we start - delay units are milliseconds
  FastLED.addLeds<CHIPSET, LED_PIN, COLOR_ORDER>(leds, NUM_LEDS).setCorrection(TypicalLEDStrip);
  FastLED.setBrightness(BRIGHTNESS);  // set brightness of all LEDs
}

void loop() {
  ChristmasLights(); // run simulation frame
  
  FastLED.show(); // display this frame
  FastLED.delay(1000 / FRAMES_PER_SECOND);
}

void ChristmasLights() {
  // Fill the LEDs with a repeating pattern of red, green, and white
  for (int i = 0; i < NUM_LEDS; i += 3) {
    leds[i] = CRGB::Red; 
    if (i + 1 < NUM_LEDS) leds[i + 1] = CRGB::Green;
    if (i + 2 < NUM_LEDS) leds[i + 2] = CRGB::White;
  }

  // Optional: Add a twinkling effect
  int brightness = 64;
  for (int j = 0; j < NUM_LEDS; j++) {
    if (random8() < 50) { // adjust the probability to control twinkling
      leds[j] = CRGB::Black; // 'turn off' the LED randomly
//      leds[j].maximizeBrightness(brightness); // do not use maximizeBrightness!
    }
  }
}