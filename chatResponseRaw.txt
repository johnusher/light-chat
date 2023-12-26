#include <FastLED.h>

#define LED_PIN     6
#define COLOR_ORDER GRB
#define CHIPSET     WS2812B
#define NUM_LEDS    30

#define BRIGHTNESS  150
#define FRAMES_PER_SECOND 30

CRGB leds[NUM_LEDS];

void setup() {
  delay(1000);
  FastLED.addLeds<CHIPSET, LED_PIN, COLOR_ORDER>(leds, NUM_LEDS).setCorrection(TypicalLEDStrip);
  FastLED.setBrightness(BRIGHTNESS);
}

void loop() {
  randomColorMovingLight(); // Modified to change color randomly

  FastLED.show();
  FastLED.delay(1000 / FRAMES_PER_SECOND);
}

void randomColorMovingLight() {
  static int pos = 0;
  static bool direction = true;
  static CRGB currentColor = CRGB::Blue;
  
  // Turn off previous LED
  leds[pos] = CRGB::Black;
  
  // Move the position
  if (direction) {
    pos++;
    if (pos == NUM_LEDS) {
      pos = NUM_LEDS - 1;
      direction = false;
      currentColor = randomColor();
    }
  } else {
    pos--;
    if (pos == -1) {
      pos = 0;
      direction = true;
      currentColor = randomColor();
    }
  }
  
  // Set the current LED to the random color
  leds[pos] = currentColor;
}

CRGB randomColor() {
  int r = random(256);
  int g = random(256);
  int b = random(256);
  return CRGB(r, g, b);
}