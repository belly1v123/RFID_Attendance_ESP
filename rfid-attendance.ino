#include <WiFi.h>
#include <HTTPClient.h>
#include <SPI.h>
#include <MFRC522.h>
#include <ArduinoJson.h>  // For JSON parsing (install via Library Manager)

#define SS_PIN 5
#define RST_PIN 27

// RGB LED pins (PWM capable)
#define RED_PIN 13
#define GREEN_PIN 12
#define BLUE_PIN 14
#define BUZZER_PIN 27

const char* ssid = "Rons_pc";
const char* password = "rons-pc-pw0";

const char* serverURL = "http://192.168.137.1:3000/api/scan";

MFRC522 rfid(SS_PIN, RST_PIN);

void setup() {
  Serial.begin(115200);
  SPI.begin();
  rfid.PCD_Init();

  pinMode(RED_PIN, OUTPUT);
  pinMode(GREEN_PIN, OUTPUT);
  pinMode(BLUE_PIN, OUTPUT);
  pinMode(BUZZER_PIN, OUTPUT);

  connectWiFi();
}

void loop() {
  if (!rfid.PICC_IsNewCardPresent() || !rfid.PICC_ReadCardSerial()) {
    delay(100);
    return;
  }

  String uid = "";
  for (byte i = 0; i < rfid.uid.size; i++) {
    if (rfid.uid.uidByte[i] < 0x10) uid += "0";
    uid += String(rfid.uid.uidByte[i], HEX);
  }
  uid.toUpperCase();

  Serial.println("Scanned UID: " + uid);

  String response = sendUIDToServer(uid);
  if (response.length() > 0) {
    parseServerResponse(response);
  } else {
    indicateStatus("error");
  }

  rfid.PICC_HaltA();
  delay(2000); // Debounce delay before next read
}

void connectWiFi() {
  Serial.print("Connecting to WiFi");
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("\nConnected to WiFi!");
}

String sendUIDToServer(const String& uid) {
  if (WiFi.status() != WL_CONNECTED) {
    Serial.println("WiFi not connected!");
    return "";
  }

  HTTPClient http;
  http.begin(serverURL);
  http.addHeader("Content-Type", "application/json");

  String payload = "{\"uid\":\"" + uid + "\"}";
  int httpResponseCode = http.POST(payload);

  if (httpResponseCode > 0) {
    String res = http.getString();
    Serial.println("Server response: " + res);
    http.end();
    return res;
  } else {
    Serial.print("HTTP POST failed, error: ");
    Serial.println(httpResponseCode);
    http.end();
    return "";
  }
}

void parseServerResponse(const String& res) {
  StaticJsonDocument<200> doc;
  DeserializationError error = deserializeJson(doc, res);
  if (error) {
    Serial.print("JSON parse error: ");
    Serial.println(error.c_str());
    indicateStatus("error");
    return;
  }

  const char* status = doc["status"];
  if (status == nullptr) {
    Serial.println("No status in response");
    indicateStatus("error");
    return;
  }

  String stat = String(status);
  Serial.println("Status: " + stat);

  if (stat == "active") {
    indicateStatus("registered");
  } else if (stat == "disabled") {
    indicateStatus("disabled");
  } else if (stat == "unregistered") {
    indicateStatus("unregistered");
  } else {
    indicateStatus("error");
  }
}

// RGB + buzzer feedback
void indicateStatus(const String& status) {
  if (status == "registered") {
    setLED(0, 255, 0);    // Green
    beep(1, 100);
  } else if (status == "unregistered") {
    setLED(255, 0, 0);    // Red
    beep(3, 100);
  } else if (status == "disabled") {
    setLED(0, 0, 255);    // Blue
    beep(1, 500);
  } else {
    // Error - Yellow blink for example
    setLED(255, 255, 0);
    beep(2, 150);
  }
  delay(500);
  setLED(0, 0, 0); // Turn off LED
}

void beep(int times, int duration) {
  for (int i = 0; i < times; i++) {
    digitalWrite(BUZZER_PIN, HIGH);
    delay(duration);
    digitalWrite(BUZZER_PIN, LOW);
    delay(100);
  }
}

void setLED(int r, int g, int b) {
  analogWrite(RED_PIN, r);
  analogWrite(GREEN_PIN, g);
  analogWrite(BLUE_PIN, b);
}