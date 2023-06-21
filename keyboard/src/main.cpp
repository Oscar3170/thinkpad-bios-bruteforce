#include "com.h"
#include "keyboard.h"


void setup() {
    Serial.begin(9600);
    com::begin();
}

void loop() {
    com::update();
    if (com::hasData()) {
        const buffer_t& buffer = com::getBuffer();

        keyboard::type(buffer.data, buffer.len);

        com::sendDone();
    }
}
