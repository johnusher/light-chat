#include <FastLED.h>

#define LED_PIN     6  
#define COLOR_ORDER GRB  
#define CHIPSET     WS2812B   
#define NUM_LEDS    30  

#define BRIGHTNESS  150   
#define FRAMES_PER_SECOND 20   // decreased for slower motion

CRGB leds[NUM_LEDS];

void setup() {
  FastLED.addLeds<CHIPSET, LED_PIN, COLOR_ORDER>(leds, NUM_LEDS).setCorrection(TypicalLEDStrip);
  FastLED.setBrightness(BRIGHTNESS);  
}

void loop() {
  TripZone(); // run the Trip Zone pattern
  FastLED.delay(1000 / FRAMES_PER_SECOND);
}

void TripZone() {
  // random number generator
  random16_add_entropy(random());

  // Set up a color palette of psychedelic colors
  CRGBPalette16 palette = PartyColors_p;
  
  // Step through the LEDs, using the palette to get the color
  for (int i = 0; i < NUM_LEDS; ++i) {
    // The randomness is added to the triwave function to make the transition more random
    leds[i] = ColorFromPalette(palette, triwave8(random8() + (i * 256 / NUM_LEDS) + (millis() / 8)));
  }

  FastLED.show(); // display this frame
}
