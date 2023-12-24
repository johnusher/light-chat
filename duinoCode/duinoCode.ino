#include <FastLED.h>

#define LED_PIN     6
#define COLOR_ORDER GRB
#define CHIPSET     WS2812B
#define NUM_LEDS    144

#define BRIGHTNESS  200
#define FRAMES_PER_SECOND 1

CRGB leds[NUM_LEDS];

void setup() {
  delay(3000); // sanity delay
  FastLED.addLeds<CHIPSET, LED_PIN, COLOR_ORDER>(leds, NUM_LEDS).setCorrection(TypicalLEDStrip);
  FastLED.setBrightness(BRIGHTNESS);
}

void loop() {
  ChristmasLights(); // run simulation frame
  
  FastLED.show(); // display this frame
  FastLED.delay(1000 / FRAMES_PER_SECOND);
}

void ChristmasLights() {
  // Alternate between red, yellow, and green colors
  for (int i = 0; i < NUM_LEDS; i += 3) {
    leds[i] = CRGB::Red;
    if (i + 1 < NUM_LEDS) leds[i + 1] = CRGB::Yellow;
    if (i + 2 < NUM_LEDS) leds[i + 2] = CRGB::Green;
  }

  // Optional: Add a flashing effect
  for (int j = 0; j < NUM_LEDS; j++) {
    if (random8() < 50) { // adjust the probability to control flashing
      leds[j] = CRGB::Black; // 'turn off' the LED randomly
    }
  }
}