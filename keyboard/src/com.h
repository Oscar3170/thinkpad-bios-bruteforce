#pragma once

#include <Arduino.h>
#include <Wire.h>
#include "debug.h"
// #include "HID-Project.h"

#define BUFFER_SIZE 512


typedef struct buffer_t {
    char   data[BUFFER_SIZE];      // !< Array to buffer incoming bytes
    size_t len;            // !< How many bytes are currently in the buffer
} buffer_t;

namespace com
{

    void begin();
    void begin(int address, int ledPin);
    void onReceive(int len);
    void onRequest();
    void setLED(int len);
    void readMessage(int len);

    void update();
    bool hasData();
    const buffer_t& getBuffer();
    void sendDone();

    void process(const char* str, size_t len);
};
