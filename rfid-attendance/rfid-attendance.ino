#include <WiFi.h>
#include <HTTPClient.h>
#include <SPI.h>
#include <MFRC522.h>
#include <ArduinoJson.h>

#define SS_PIN 5
#define RST_PIN 27

// RGB LED pins
#define RED_PIN 13
#define GREEN_PIN 12
#define BLUE_PIN 14
#define BUZZER_PIN 22 // Changed from 27 to 22 to avoid conflict

const char *ssid = "Pranjal_2.4";
const char *password = "PK@98400";

const char *serverURL = "http://192.168.137.1:3000/api/scan";

MFRC522 rfid(SS_PIN, RST_PIN);

void setup()
{
  Serial.begin(115200);
  SPI.begin();
  rfid.PCD_Init();

  pinMode(RED_PIN, OUTPUT);
  pinMode(GREEN_PIN, OUTPUT);
  pinMode(BLUE_PIN, OUTPUT);
  pinMode(BUZZER_PIN, OUTPUT);

  connectWiFi();
}

void loop()
{
  ensureWiFiConnected();

  if (!rfid.PICC_IsNewCardPresent() || !rfid.PICC_ReadCardSerial())
  {
    delay(100);
    return;
  }

  String uid = "";
  for (byte i = 0; i < rfid.uid.size; i++)
  {
    if (rfid.uid.uidByte[i] < 0x10)
      uid += "0";
    uid += String(rfid.uid.uidByte[i], HEX);
  }
  uid.toUpperCase();
  Serial.println("Scanned UID: " + uid);

  String response = sendUIDToServer(uid);
  if (response.length() > 0)
  {
    parseServerResponse(response);
  }
  else
  {
    indicateStatus("error");
  }

  rfid.PICC_HaltA();
  rfid.PCD_StopCrypto1(); // resets crypto communication
  delay(2000);            // Debounce delay
}

void connectWiFi()
{
  Serial.print("Connecting to WiFi");

  WiFi.disconnect(true);
  delay(100);
  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED)
  {
    Serial.print(".");
    blinkLED(0, 0, 255, 200); // Blink Blue while connecting
    delay(300);
  }

  Serial.println("\nConnected to WiFi!");
  setLED(0, 255, 0); // Solid Green when connected
  delay(500);
  setLED(0, 0, 0); // Turn off LED
}

void ensureWiFiConnected()
{
  if (WiFi.status() != WL_CONNECTED)
  {
    Serial.println("WiFi disconnected, reconnecting...");
    connectWiFi();
  }
}

String sendUIDToServer(const String &uid)
{
  if (WiFi.status() != WL_CONNECTED)
  {
    Serial.println("WiFi not connected!");
    return "";
  }

  HTTPClient http;
  http.begin(serverURL);
  http.setTimeout(5000); // 5 seconds timeout
  http.addHeader("Content-Type", "application/json");

  String payload = "{\"uid\":\"" + uid + "\"}";
  int httpResponseCode = http.POST(payload);

  String res = "";
  if (httpResponseCode > 0)
  {
    res = http.getString();
    Serial.println("Server response: " + res);
  }
  else
  {
    Serial.print("HTTP POST failed, error: ");
    Serial.println(httpResponseCode);
  }
  http.end(); // Always end to free resources
  return res;
}

void parseServerResponse(const String &res)
{
  StaticJsonDocument<200> doc;
  DeserializationError error = deserializeJson(doc, res);
  if (error)
  {
    Serial.print("JSON parse error: ");
    Serial.println(error.c_str());
    indicateStatus("error");
    return;
  }

  const char *status = doc["status"];
  if (status == nullptr)
  {
    Serial.println("No status in response");
    indicateStatus("error");
    return;
  }

  String stat = String(status);
  Serial.println("Status: " + stat);

  if (stat == "checked_in")
  {
    indicateStatus("registered");
  }
  else if (stat == "checkout")
  {
    indicateStatus("checked_out");
  }
  else if (stat == "unrecognized")
  {
    indicateStatus("unregistered");
  }
  else if (stat == "duplicate")
  {
    indicateStatus("duplicate");
  }
  else
  {
    indicateStatus("error");
  }
}

void indicateStatus(const String &status)
{
  if (status == "registered")
  {
    setLED(0, 255, 0); // Green
    beep(1, 100);
  }
  else if (status == "checked_out")
  {
    setLED(0, 0, 255); // Blue
    beep(1, 200);
  }
  else if (status == "duplicate")
  {
    setLED(255, 255, 0); // Yellow
    beep(2, 100);
  }
  else if (status == "unregistered")
  {
    setLED(255, 0, 0); // Red
    beep(3, 100);
  }
  else
  {
    setLED(255, 0, 255); // Purple for error
    beep(4, 80);
  }
  delay(500);
  setLED(0, 0, 0); // Turn off LED
}

void beep(int times, int duration)
{
  for (int i = 0; i < times; i++)
  {
    digitalWrite(BUZZER_PIN, HIGH);
    delay(duration);
    digitalWrite(BUZZER_PIN, LOW);
    delay(100);
  }
}

void setLED(int r, int g, int b)
{
  analogWrite(RED_PIN, r);
  analogWrite(GREEN_PIN, g);
  analogWrite(BLUE_PIN, b);
}

void blinkLED(int r, int g, int b, int duration)
{
  setLED(r, g, b);
  delay(duration);
  setLED(0, 0, 0);
  delay(100);
}
