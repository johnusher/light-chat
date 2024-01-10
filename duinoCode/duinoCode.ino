#include <FastLED.h>

#define LED_PIN     6  
#define COLOR_ORDER GRB  
#define CHIPSET     WS2812B
#define NUM_LEDS    30

#define BRIGHTNESS  150
#define FRAMES_PER_SECOND 60  

CRGB leds[NUM_LEDS];
uint16_t pos = 0;
uint8_t speed = 0;  // controls the speed of the move
uint8_t color_loop = 0;  // for random color

void setup() {
  delay(1000);
  FastLED.addLeds<CHIPSET, LED_PIN, COLOR_ORDER>(leds, NUM_LEDS);
  FastLED.setBrightness(BRIGHTNESS);
}

void loop() {
  MoveDot();
  FinalFlash();
}

void MoveDot() {
  fill_solid(leds, NUM_LEDS, CRGB::Black);
  leds[pos] = (pos >= NUM_LEDS / 2 ? CRGB::Yellow : CRGB::Blue);
  FastLED.show();
  FastLED.delay(1000 / (FRAMES_PER_SECOND + speed));

  // Handle turning around
  if (pos >= NUM_LEDS) {
    pos = NUM_LEDS - 1;
    for (int i = 0; i < 5; i++) {  // flashing yellow dot
      leds[pos] = CRGB::Black;
      FastLED.show();
      delay(100);
      leds[pos] = CRGB::Yellow;
      FastLED.show();
      delay(100);
    }
    pos--;
    speed -= 5;
  } else {
    pos++;
    speed += 2;
  }

   // Changing color to random at every loop
   color_loop++;
   if(color_loop > NUM_LEDS) {
     color_loop = 0;
     fill_solid(leds, NUM_LEDS, CHSV(random8(), 255, 192));
     FastLED.show();
     delay(500);
   }

}

void FinalFlash(){
  if(pos == 0){
    for(int i = 0; i < 4 ; i++){
      fill_solid(leds, NUM_LEDS, CRGB::White);
      FastLED.show();
      delay(100);
      fill_solid(leds, NUM_LEDS, CRGB::Black);
      FastLED.show();
      delay(100);
    }
    pos = 0;
    speed = 0; 
    color_loop = 0;  // reseting color_loop for next cycle
  }
}
