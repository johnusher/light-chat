#include <FastLED.h>

#define LED_PIN     6  
#define COLOR_ORDER GRB  
#define CHIPSET     WS2812B 
#define NUM_LEDS    30  
#define BRIGHTNESS  150
#define FRAMES_PER_SECOND 60

CRGB leds[NUM_LEDS];

void setup() {
  delay(1000);
  FastLED.addLeds<CHIPSET, LED_PIN, COLOR_ORDER>(leds, NUM_LEDS).setCorrection(TypicalLEDStrip);
  FastLED.setBrightness(BRIGHTNESS); 
  FastLED.clear();
}

void loop() {
  RainbowCircle();
}

void RainbowCircle() {
  for (int lap = 0; lap < 8; lap++) {
    int speed = 5 + 20 * lap;  
    for (int i = 0; i < NUM_LEDS; i++) {
      leds[i] = CHSV(i * 256.0 / NUM_LEDS + lap * 256.0 / 8, 255, 255); 
      FastLED.show();
      FastLED.clear();
      delay(1000 / speed); 
    }
  }
}
