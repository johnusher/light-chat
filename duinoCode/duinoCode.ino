#include <FastLED.h>

#define LED_PIN     6  // do not change
#define COLOR_ORDER GRB  // do not change
#define CHIPSET     WS2812B   // do not change
#define NUM_LEDS    30  // do not change

#define BRIGHTNESS  150    // set brightness of all LEDs. 150 is default
#define FRAMES_PER_SECOND 5   // 20 is very fast, 10 is normal, 5 is slow. Above 30 is too fast.

CRGB leds[NUM_LEDS];

void setup() {
  delay(1000); // sanity delay
  FastLED.addLeds<CHIPSET, LED_PIN, COLOR_ORDER>(leds, NUM_LEDS).setCorrection(TypicalLEDStrip);
  FastLED.setBrightness(BRIGHTNESS);  // set brightness of all LEDs
}

void loop() {
  PartyLights(); // run simulation frame
  
  FastLED.show(); // display this frame
  FastLED.delay(1000 / FRAMES_PER_SECOND);
}

void PartyLights() {
  // Fill the LEDs with a repeating pattern of different shades of blue and green
  for (int i = 0; i < NUM_LEDS; i += 3) {
    leds[i] = CRGB(255, 0, 0);  //Red
    if (i + 1 < NUM_LEDS) leds[i + 1] = CRGB(0, 255, 0);  //Green
    if (i + 2 < NUM_LEDS) leds[i + 2] = CRGB(180, 0, 255);  //Purple
  }

  // Optional: Add a twinkling effect
  int brightness = 64;
  for (int j = 0; j < NUM_LEDS; j++) {
    if (random8() < 50) { // adjust the probability to control twinkling
      leds[j] = CRGB::Black; // 'turn off' the LED randomly
    }
  }
}
