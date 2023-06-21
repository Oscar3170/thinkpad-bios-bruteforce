// #include "HID-Project.h"

// const int buttonPin = 4;        // input pin for pushbutton
// int previousButtonState = HIGH; // for checking the state of a pushButton
// int counter = 0;                // button push counter

// int defaultDelay = 20; // default delay between key actions

// void keyboardButtonSetup() {
//     // make the pushButton pin an input:
//     pinMode(buttonPin, INPUT);
//     // generates seed from unconnected pin
//     randomSeed(analogRead(0));
//     // initialize control over the keyboard:
//     BootKeyboard.begin();
// }

// void keyboardButtonLoop() {
//     // read the pushbutton:
//     int buttonState = digitalRead(buttonPin);
//     // if the button state has changed,
//     if ((buttonState != previousButtonState)
//         // and it's currently pressed:
//         && (buttonState == HIGH)) {
//         // increment the button counter
//         counter++;

//         if (counter == 1) {
//             BootKeyboard.press(KEY_ESC);
//         }
//         else if (counter % 2 == 0) {
//             BootKeyboard.print("Pressed " + String(counter) + " times.");
//         }
//         else {
//             for (int i = 0; i < 20; i++) {
//                 BootKeyboard.press(KEY_BACKSPACE);
//             }
//         }
//     }
//     // save the current button state for comparison next time:
//     previousButtonState = buttonState;
// }