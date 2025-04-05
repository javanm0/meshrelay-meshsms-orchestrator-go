package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"
    "time"
)

var (
    port                 = "3040"
    messagesApiEndpoint  = os.Getenv("MESSAGES_API_ENDPOINT")
    smsApiEndpoint       = os.Getenv("SMS_API_ENDPOINT")
    processIntervalStr   = os.Getenv("PROCESS_INTERVAL")
    processInterval      time.Duration
)

func init() {
    var err error
    processInterval, err = time.ParseDuration(processIntervalStr)
    if err != nil {
        log.Fatalf("Invalid PROCESS_INTERVAL: %v\n", err)
    }
}

type Message struct {
    ID          string `json:"_id"`
    PhoneNumber string `json:"phoneNumber"`
    Message     string `json:"message"`
    MessageSent bool   `json:"messageSent"`
}

func fetchMessages() ([]Message, error) {
    resp, err := http.Get(messagesApiEndpoint)
    if err != nil {
        log.Printf("Error fetching messages: %v\n", err)
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %v\n", err)
        return nil, err
    }

    var messages []Message
    err = json.Unmarshal(body, &messages)
    if err != nil {
        log.Printf("Error unmarshalling messages: %v\n", err)
        return nil, err
    }

    return messages, nil
}

func sendSMS(to, body, messageID string) error {
    payload := fmt.Sprintf(`{"to":"%s","body":"%s"}`, to, body)
    resp, err := http.Post(smsApiEndpoint, "application/json", strings.NewReader(payload))
    if err != nil {
        log.Printf("Error sending SMS: %v\n", err)
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        updatePayload := fmt.Sprintf(`{"id":"%s"}`, messageID)
        req, err := http.NewRequest(http.MethodPut, messagesApiEndpoint, strings.NewReader(updatePayload))
        if err != nil {
            log.Printf("Error creating update request: %v\n", err)
            return err
        }
        req.Header.Set("Content-Type", "application/json")

        client := &http.Client{}
        updateResp, err := client.Do(req)
        if err != nil {
            log.Printf("Error updating message status: %v\n", err)
            return err
        }
        defer updateResp.Body.Close()

        if updateResp.StatusCode == http.StatusOK {
            log.Printf("Message with ID %s marked as sent\n", messageID)
        }
    }

    return nil
}

func processMessages() {
    messages, err := fetchMessages()
    if err != nil {
        log.Println("Error fetching messages, skipping processing")
        return
    }

    for _, message := range messages {
        if message.PhoneNumber != "" && !message.MessageSent {
            log.Printf("Processing message with ID: %s\n", message.ID)
            err := sendSMS(message.PhoneNumber, message.Message, message.ID)
            if err == nil {
                log.Printf("SMS sent to %s\n", message.PhoneNumber)
            }
        } else {
            log.Printf("Skipping message with ID: %s as it is already sent or missing phone number\n", message.ID)
        }
    }
}

func main() {
    log.Printf("Server is running on port %s\n", port)

    // Initial call
    go processMessages()

    // Call at the specified interval
    ticker := time.NewTicker(processInterval)
    defer ticker.Stop()

    for range ticker.C {
        processMessages()
    }
}