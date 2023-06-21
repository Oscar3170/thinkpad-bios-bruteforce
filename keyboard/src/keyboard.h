#include <Arduino.h>
#include "HID-Project.h"
#include "debug.h"

// const int buttonPin = 4;        // input pin for pushbutton
// int previousButtonState = HIGH; // for checking the state of a pushButton
// int counter = 0;                // button push counter

int defaultDelay = 20; // default delay between key actions

namespace keyboard {
    void begin() {BootKeyboard.begin();}


    String parseBuffer(const char* buffer, size_t len) {
        String str;
        for (size_t i = 0; i < len; i++) {
            str += buffer[i];
        }
        return str;
    }

    void type(const char* buffer, size_t len) {
        String str = parseBuffer(buffer, len);
        Serial.println(str);
        if (str == String("ESCAPE")) {
            Serial.println("Pressing escape");
            BootKeyboard.write(KEY_ESC);
        } else if (str == String("TAB")) {
            Serial.println("Pressing tab");
            BootKeyboard.write(KEY_TAB);
        } else if (str == String("SHIFT")) {
            Serial.println("Pressing shift");
            BootKeyboard.write(KEY_LEFT_SHIFT);
        } else {
            BootKeyboard.print(str);
        }
    }
}
