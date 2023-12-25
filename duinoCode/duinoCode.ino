#include <FastLED.h>

#define LED_PIN     6
#define COLOR_ORDER GRB
#define CHIPSET     WS2812B
#define NUM_LEDS    30

#define BRIGHTNESS  200
#define FRAMES_PER_SECOND 10 // changed to twinkle slower

CRGB leds[NUM_LEDS];

void setup() {
  delay(3000); // sanity delay
  FastLED.addLeds<CHIPSET, LED_PIN, COLOR_ORDER>(leds, NUM_LEDS).setCorrection(TypicalLEDStrip);
  FastLED.setBrightness(BRIGHTNESS);
}

void loop() {
  TwinkleLights(); // run twinkle simulation frame
  
  FastLED.show(); // display this frame
  FastLED.delay(1000 / FRAMES_PER_SECOND);
}

void TwinkleLights() {
  // Make some LEDs randomly twinkle by setting them to a random brightness level
  for (int i = 0; i < NUM_LEDS; i++) {
    if (random8() < 64) { // 25% chance of twinkle
      leds[i] = CRGB::White;
    } else {
      leds[i] = CRGB::Black;
    }
  }
}