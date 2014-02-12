package amp

/*
 * implements amp_diagram.svg
 * */
 
import "net"
 

type AnswerBox struct {
    Response *map[string]string    
    Error error
    Command *Command
    Callback chan *AnswerBox
}

type AskBox struct {
    Arguments *map[string]string
    Response *map[string]string
    Client *Connection
    Command *Command
}

type Command struct {
    Name string
    Responder chan *AskBox
    Arguments []string
    Response []string
}

type AMP struct {
    ConnList map[string]*Connection
    Commands map[string]*Command
    BoxCounter int
    Callbacks map[string]*AnswerBox
    GetBoxCounter chan chan int
    //ListenTCP chan string
    //ConnectTCP chan string
}

type ClientCreator struct {
    Name string
    Service string
}

type Connection struct {
    Name string
    Conn net.Conn
    Protocol *AMP
    Quit chan bool
    Closed bool
}
