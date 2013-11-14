package wsloadtest

import (
    "code.google.com/p/go.net/websocket"
    "fmt"
    "time"
    "log"
)

type Client struct {
    profile *ClientProfile
    ws *websocket.Conn
}

type ClientProfile struct {
    Id int
    NumberOfRequestsPerMinute int
}

const (
    defaultOrigin = "http://localhost"
)

func NewClient(testUrl string, profile *ClientProfile) (cli *Client, err error) {
    ws, err := websocket.Dial(testUrl, "", defaultOrigin)

    if err != nil {
       goto Error
    }

    cli = &Client{profile: profile, ws: ws}
    go cli.reader()
    go cli.writer()
    return

Error:
    return nil, err
}

func (cli *Client) Close() {
    cli.log("shutting down...")
    cli.ws.Close()
}

// TODO: Make this a bit more sophisticated than
// having fixed intervals between requests.
func (cli *Client) sleepDuration() time.Duration {
    sleepAsFloat := float64(60) / float64(cli.profile.NumberOfRequestsPerMinute)
    sleepAsString := fmt.Sprintf("%vs", sleepAsFloat)
    duration, _ := time.ParseDuration(sleepAsString)

    log.Printf("%v", sleepAsString)

    return duration
}

func (cli *Client) writer() {
    sleepDuration := cli.sleepDuration()

    for {
        time.Sleep(sleepDuration)
        cli.log("sending 'hello' to the remote server.")
        err := websocket.Message.Send(cli.ws, "hello")
        if err != nil {
            cli.log("could not send message (%v).", err)
            break
        }
    }

    cli.Close()
}

func (cli *Client) reader() {
    for {
        var message = ""

        err := websocket.Message.Receive(cli.ws, &message)
        if err != nil {
            cli.log("could not receive message (%v)", err)
            break
        }
        cli.log("received '%v' from remote server.", message)
    }

    cli.Close()
}


func (cli *Client) log(format string, v ...interface{}) {
    prefix := fmt.Sprintf("Client %v: ", cli.profile.Id)
    log.Printf(prefix + format, v...)
}