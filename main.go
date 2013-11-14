package main

import (
  "log"
  "os"
  "os/signal"
  "syscall"
  "github.com/mikeys/go-websocket-redis-test1/wsloadtest"
)

const (
  defaultNumberOfClients = 249
  testUrl = "ws://localhost:8080"
)

var (
  clients []*wsloadtest.Client
)

func main() {
  log.Printf("Press Ctrl-C to exit.")
  createClients()
  waitForSignalAndShutdown()
}

func createClients() {
  clients = make([]*wsloadtest.Client, defaultNumberOfClients, defaultNumberOfClients)

  for i := 0; i < defaultNumberOfClients; i++ {
    createClient(i)
  }
}

func createClient(id int) {
  log.Printf("Creating client %v...", id)

  profile := wsloadtest.ClientProfile{Id: id, NumberOfRequestsPerMinute: 10}
  client, err := wsloadtest.NewClient(testUrl, &profile)
  
  if err != nil {
    log.Fatalf("Could not create client %v:\n%v\nshutting down...", id, err)
  }

  clients[id] = client

  log.Printf("Client %v created.", id)
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