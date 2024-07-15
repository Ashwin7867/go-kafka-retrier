# Go Kafka Retrirer

The go-kafka-retrier package is a Go-based library designed to enhance the reliability of Kafka message processing. It implements a configurable retry mechanism and Dead Letter Queue (DLQ) handling to ensure unprocessable messages are handled appropriately.

# Key Features

- Configurable Retry Mechanism: Allows setting the maximum number of retries and delay between retries.
- DLQ Handling: Redirects messages to a Dead Letter Queue after exceeding the retry limit.
- Modular Design: Separates configuration and retry logic into distinct packages.
- Comprehensive Logging: Provides detailed logs for monitoring and debugging.
- Sample Application: Includes a demonstration application for easy integration and usage.

## Installation

1. Clone the repository:

```bash
git clone https://github.com/Ashwin7867/go-kafka-retrier.git
cd go-kafka-retrier
```

2. Install dependencies:

Ensure you have Go installed, then run:

```bash
go mod tidy
```

## Usage

### 1. Configuration
-Create a .env file in your project directory with the following content:

```env
KAFKA_BROKERS=localhost:9092
CONSUMER_GROUP=my-consumer-group
TOPIC=test-topic
DLQ_TOPIC=test-dlq-topic
MAX_RETRIES=3
RETRY_DELAY=5s
```

### 2. Integration
- Import the package and use it in your Go application:

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/Ashwin7867/go-kafka-retrier/config"
    "github.com/Ashwin7867/go-kafka-retrier/retrier"
)

func main() {
    cfg := config.LoadConfig()

    retrier, err := retrier.NewRetrier(cfg)
    if err != nil {
        log.Fatalf("Failed to create retrier: %v", err)
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt, syscall.SIGTERM)
        <-c
        cancel()
    }()

    if err := retrier.ProcessMessages(ctx); err != nil {
        log.Fatalf("Failed to process messages: %v", err)
    }
}

```

### 3. Running the applciation
- Build and run your application:
 
```bash
go run main.go
```



## Testing the DLQ Mechanism

### 1. Simulate Message Failures:

- Modify the HandleMessage function in retrier/retrier.go to simulate failures:

```go
func (r *Retrier) HandleMessage(msg *sarama.ConsumerMessage) error {
    if string(msg.Value) == "simulate-failure" {
        return fmt.Errorf("simulated error")
    }
    return nil
}

```

### 2. Send Test Messages:

- Use Kafka console producer to send test messages:

```bash
kafka-console-producer.sh --broker-list localhost:9092 --topic test-topic
> simulate-failure

```
### 3. Check DLQ Messages:

- Use Kafka console consumer to read messages from the DLQ topic:

```bash
kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test-dlq-topic --from-beginning

```


## Code Structure:

 1. config/config.go: Handles loading and parsing configuration from the environment.
 2. retrier/retrier.go: Implements the retry mechanism and message processing logic.
 3. retrier/handler.go: Manages message consumption and retry handling.




## Summary:
The kafka-retrier package improves the stability and efficiency of Kafka message handling in Go applications by providing a robust retry mechanism and DLQ handling. By integrating this package, you can ensure that unprocessable messages are managed effectively, contributing to the overall reliability of your message processing system.



### Contributing

Feel free to contribute by opening issues or submitting pull requests.
