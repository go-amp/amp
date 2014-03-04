package amp

import "net"
import "sync"

type Client struct {
    Name *string
    Conn *net.TCPConn
    prot *AMP
}

type AskBox struct {
    data *map[string]string
    client *Client
}

type CallBox struct {
    Args map[string]string
    Response *map[string]string
    Command *string    
    Callback chan *CallBox
    CallbackArgs *interface{}
}

type AMP struct {
    commands map[string]chan *AskBox
    callbacks map[string]*CallBox
    commands_mutex *sync.Mutex
    callbacks_mutex *sync.Mutex
}

