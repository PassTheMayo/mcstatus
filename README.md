# MCStatus
A Go library for retrieving the status of a Minecraft server.

## Installation

```bash
go get -u github.com/PassTheMayo/mcstatus
```

## Documentation

https://pkg.go.dev/github.com/PassTheMayo/mcstatus

## Usage

### Status

```go
import "github.com/PassTheMayo/mcstatus"

func main() {
    response, err := mcstatus.Status("play.hypixel.net", 25565)

    if err != nil {
        panic(err)
    }

    fmt.Println(response)
}
```

### Bedrock Status

```go
import "github.com/PassTheMayo/mcstatus"

func main() {
    response, err := mcstatus.StatusBedrock("127.0.0.1", 19132)

    if err != nil {
        panic(err)
    }

    fmt.Println(response)
}
```

### Basic Query

```go
import "github.com/PassTheMayo/mcstatus"

func main() {
    response, err := mcstatus.BasicQuery("play.hypixel.net", 25565)

    if err != nil {
        panic(err)
    }

    fmt.Println(response)
}
```

### Full Query

```go
import "github.com/PassTheMayo/mcstatus"

func main() {
    response, err := mcstatus.FullQuery("play.hypixel.net", 25565)

    if err != nil {
        panic(err)
    }

    fmt.Println(response)
}
```

### RCON

```go
import "github.com/PassTheMayo/mcstatus"

func main() {
    client := mcstatus.NewRCON()

    if err := client.Dial("127.0.0.1", 25575); err != nil {
        panic(err)
    }

    if err := client.Login("mypassword"); err != nil {
        panic(err)
    }

    if err := client.Run("say Hello, world!"); err != nil {
        panic(err)
    }

    fmt.Println(<- client.Messages)

    if err := client.Close(); err != nil {
        panic(err)
    }
}
```

## Send Vote

```go
import "github.com/PassTheMayo/mcstatus"

func main() {
    err := mcstatus.SendVote("127.0.0.1", 8192, mcstatus.VoteOptions{
		ServiceName: "my-service",
		Username:    "PassTheMayo",
		Token:       "abc123", // server's Votifier token
		UUID:        "",       // recommended but not required, UUID with dashes
		Timestamp:   time.Now(),
		Timeout:     time.Second * 5,
	})

    if err != nil {
        panic(err)
    }
}
```