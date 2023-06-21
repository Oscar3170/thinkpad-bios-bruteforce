#include "com.h"

#define REQ_SOT 0x01     // !< Start of transmission
#define REQ_EOT 0x04     // !< End of transmission


namespace com
{
    buffer_t receive_buf;
    buffer_t data_buf;

    bool start_parser         = false;
    bool ongoing_transmission = false;

    int ledPin = 7;
    int address = 0x8;

    void begin() {
        pinMode(ledPin, OUTPUT);
        digitalWrite(ledPin, HIGH);
        Wire.begin(address);
        Wire.onReceive(onReceive);
        Wire.onRequest(onRequest);

        receive_buf.len = 0;
    }

    void begin(int address, int ledPin) {
        com::address = address;
        com::ledPin = ledPin;
        begin();
    }

    void onRequest() {
        bool done = receive_buf.len + data_buf.len == 0;
        Serial.println(String{""} + "Requested status, done: " + done);
        Wire.write((uint8_t*)&done, sizeof(bool));
    }

    void onReceive(int len) {
        if (receive_buf.len + (unsigned int)len <= BUFFER_SIZE) {
            Wire.readBytes(&receive_buf.data[receive_buf.len], len);
            receive_buf.len += len;
        }
    }

    void process(const char* str, size_t len) {
        Serial.println(str);
    }

    bool hasData() {
        return data_buf.len > 0 && start_parser;
    }

    const buffer_t& getBuffer() {
        return data_buf;
    }

    void update() {
        if (!start_parser && (receive_buf.len > 0) && (data_buf.len < BUFFER_SIZE)) {
            unsigned int i = 0;

            debugs("RECEIVED ");

            // ! Skip bytes until start of transmission
            while (i < receive_buf.len && !ongoing_transmission) {
                if (receive_buf.data[i] == REQ_SOT) {
                    ongoing_transmission = true;
                    debugs("[SOT] ");
                }
                ++i;
            }

            debugs("'");

            while (i < receive_buf.len && ongoing_transmission) {
                char c = receive_buf.data[i];

                if (c == REQ_EOT) {
                    start_parser         = true;
                    ongoing_transmission = false;
                } else {
                    debug(c);

                    data_buf.data[data_buf.len] = c;
                    ++data_buf.len;
                }

                if (data_buf.len == BUFFER_SIZE) {
                    start_parser         = true;
                    ongoing_transmission = false;
                }

                ++i;
            }

            debugs("' ");

            if (start_parser && !ongoing_transmission) {
                debugs("[EOT]");
            } else if (!start_parser && ongoing_transmission) {
                debugs("...");
            } else if (!start_parser && !ongoing_transmission) {
                debugs("DROPPED");
            }

            debugln();

            receive_buf.len = 0;
        }
    }

    void sendDone() {
        data_buf.len = 0;
        start_parser = false;
    }
}
