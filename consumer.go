package main

import (
  "fmt"
  "os"
  "os/signal"
  "syscall"
  "github.com/confluentinc/confluent-kafka-go/kafka"
  "encoding/json"
)

func consumeUsers() {
  config := &kafka.ConfigMap{
    "metadata.broker.list": "velomobile-01.srvs.cloudkafka.com:9094,velomobile-02.srvs.cloudkafka.com:9094,velomobile-03.srvs.cloudkafka.com:9094",
    "security.protocol":    "SASL_SSL",
    "sasl.mechanisms":      "SCRAM-SHA-256",
    "sasl.username":        "xz6befqu",
    "sasl.password":        "ZnTqLiR0WxwLHX_jdiGChcbi4W-H9Mzd",
    "group.id":             "users-consumer1",
    "go.events.channel.enable":        true,
    "go.application.rebalance.enable": true,
    "default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"},
  }
  topic := "xz6befqu-users"

  sigchan := make(chan os.Signal, 1)
  signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
  c, err := kafka.NewConsumer(config)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
    os.Exit(1)
  }
  fmt.Printf("Created Consumer %v\n", c)
  err = c.Subscribe(topic, nil)
  run := true
  counter := 0
  commitAfter := 1000
  for run == true {
    select {
    case sig := <-sigchan:
      fmt.Printf("Caught signal %v: terminating\n", sig)
      run = false
    case ev := <-c.Events():
      switch e := ev.(type) {
      case kafka.AssignedPartitions:
        c.Assign(e.Partitions)
      case kafka.RevokedPartitions:
        c.Unassign()
      case *kafka.Message:
        fmt.Printf("%% Message on %s: %s\n", e.TopicPartition, string(e.Value))
        counter++
        if counter > commitAfter {
          c.Commit()
          counter = 0
        }
        createGopay(e.Value)

      case kafka.PartitionEOF:
        fmt.Printf("%% Reached %v\n", e)
      case kafka.Error:
        fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
        run = false
      }
    }
  }
  fmt.Printf("Closing consumer\n")
  c.Close()
}

type User struct {
  id  string
  email string
  phone string
  first_name string
  last_name string
  password_digest string
  gopay string
  created_at string
  updated_at string
}

func createGopay(v []byte){
  // db := connectDB()
  var user User
  err := json.Unmarshal(v, &user)
  if err != nil {
    fmt.Println("userFromJson:", err)
    os.Exit(1)
  }
  fmt.Println("email " + user.email)
}