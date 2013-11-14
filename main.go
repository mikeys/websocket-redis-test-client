package main

import (
    "flag"
    "log"
    "time"
    "math/rand"
    "os"
    "os/signal"
    "syscall"
    "github.com/mikeys/websocket-redis-test-client/wsloadtest"
)

const (
    minInitialDelayInSeconds = 0
    maxInitialDelayInSeconds = 10
)

var (
    testUrl string
    numberOfClients int
    numberOfRequestsPerMinute int
    
    clients []*wsloadtest.Client
)

func init() {
    const (
        defaultTestUrl = "ws://localhost:8080"
        defaultNumberOfClients = 100
        defaultNumberOfRequestsPerMinute = 40

        testUrlDesc = "Websocket server URL"
        numberOfClientsDesc = "Number of websocket clients to load"
        numberOfRequestsPerMinuteDesc = "Number of requests each client will execute"
    )

    flag.StringVar(&testUrl, "test-url", defaultTestUrl, testUrlDesc)
    flag.StringVar(&testUrl, "u", defaultTestUrl, testUrlDesc + " (shorthand)")
    flag.IntVar(&numberOfClients, "number-of-clients", defaultNumberOfClients, numberOfClientsDesc)
    flag.IntVar(&numberOfClients, "c", defaultNumberOfClients, numberOfClientsDesc + " (shorthand)")
    flag.IntVar(&numberOfRequestsPerMinute, "requests-per-minute", defaultNumberOfRequestsPerMinute, numberOfRequestsPerMinuteDesc)
    flag.IntVar(&numberOfRequestsPerMinute, "r", defaultNumberOfRequestsPerMinute, numberOfRequestsPerMinuteDesc + " (shorthand)")
}

func main() {
    flag.Parse()
    printSettings()
    createClients()
    waitForSignalAndShutdown()
}

func printSettings() {
    log.Printf("Generating %v websocket clients", numberOfClients)
    log.Printf("Websocket server URL: %v", testUrl)
    log.Printf("Number of requests per minute: %v", numberOfRequestsPerMinute)
    log.Printf("Press Ctrl-C at any time to exit...")
    log.Printf("===============================\n")
    time.Sleep(time.Duration(5) * time.Second)
}

func createClients() {
    clients = make([]*wsloadtest.Client, numberOfClients, numberOfClients)

    for i := 0; i < numberOfClients; i++ {
        createClient(i)
    }
}

func createClient(id int) {
    log.Printf("Creating client %v...", id)

    profile := buildClientProfile(id)
    client, err := wsloadtest.NewClient(testUrl, &profile)
  
    if err != nil {
        log.Fatalf("Could not create client %v:\n%v\nshutting down...", id, err)
    }

    clients = append(clients, client)

    log.Printf("Client %v created.", id)
}

func buildClientProfile(id int) wsloadtest.ClientProfile {
    return wsloadtest.ClientProfile{Id: id,
        NumberOfRequestsPerMinute: numberOfRequestsPerMinute,
        InitialDelayInSeconds: generateRandomDelayInSeconds()}
}

func generateRandomDelayInSeconds() float32 {
    randomInt := rand.Intn(maxInitialDelayInSeconds - minInitialDelayInSeconds) + minInitialDelayInSeconds
    return float32(randomInt) + rand.Float32()
}

func waitForSignalAndShutdown() {
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
    <-sigc

    log.Printf("Shutting down...")
    shutdown()
}

func shutdown() {
    for _, client := range clients {
        if client != nil {
            client.Close()
        }
    }
}