#include <FastLED.h>

#define LED_PIN     6
#define COLOR_ORDER GRB
#define CHIPSET     WS2812B
#define NUM_LEDS    30

#define BRIGHTNESS  150
#define FLASH_RATE  500 // Flash rate in milliseconds

CRGB leds[NUM_LEDS];

void setup() {
  delay(3000);
  FastLED.addLeds<CHIPSET, LED_PIN, COLOR_ORDER>(leds, NUM_LEDS).setCorrection(TypicalLEDStrip);
  FastLED.setBrightness(BRIGHTNESS);
}

void loop() {
  FlashEffect(CRGB::Red, CRGB::Green, FLASH_RATE);
  FastLED.show();
  FastLED.delay(1000 / FRAMES_PER_SECOND);
}

void FlashEffect(CRGB onColor, CRGB offColor, unsigned long flashRate) {
  static unsigned long previousTime = 0;
  static boolean ledState = false;

  if (millis() - previousTime >= flashRate) {
    previousTime = millis();
    ledState = !ledState;

    for (int i = 0; i < NUM_LEDS; i++) {
      if (ledState) {
        leds[i] = onColor;
      } else {
        leds[i] = offColor;
      }
    }
  }
}
