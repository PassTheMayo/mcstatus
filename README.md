# MCStatus
A Go library for retrieving the status of a Minecraft server.

## Installation

```bash
go get -u github.com/PassTheMayo/mcstatus
```

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