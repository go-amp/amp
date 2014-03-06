package amp

import "net"
import "sync"
import "bufio"

type Client struct {
    Name *string
    Conn *net.TCPConn
    prot *AMP
    writer *bufio.Writer
    reader *bufio.Reader
}

type AskBox struct {
    Args map[string][]byte
    Response map[string][]byte
    client *Client
}

//type CallBox struct {    
    //Callback chan *CallBox
    //CallbackArgs *interface{}
//}

type AMP struct {
    commands map[string]chan *AskBox
    callbacks map[string]chan map[string][]byte
    commands_mutex *sync.Mutex
    callbacks_mutex *sync.Mutex
    boxCounter int
    tagger chan string
}

