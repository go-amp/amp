package amp
 
import "net"

type Command struct {
    Name string
    Responder chan *Ask
    Arguments []string
    Response []string
}

type Ask struct {
    Arguments *map[string]string
    Response *map[string]string
    ReplyChannel chan *Ask    
}

type CallBox struct {
    Arguments *map[string]string
    Response *map[string]string
    Command *Command    
    Callback chan *CallBox
    CallbackArgs *interface{}
}

type AMP struct {
    ConnList map[string]*Client
    Commands map[string]*Command
    BoxCounter int
    Callbacks map[string]*CallBox
}

type Client struct {
    Name string
    Conn net.TCPConn
    Protocol *AMP
    Quit chan bool
    Closed bool
    incoming_handler chan *map[string]string
    reply_handler chan *Ask
}
