package amp

import "net"
import "log"
import "time"
import "fmt"
import "errors"

var ASK = "_ask"
var ANSWER = "_answer"
var COMMAND = "_command"
/* not using these yet, as they are used uncomment */
//var ERROR = "_error"
//var ERROR_CODE = "_error_code"
//var ERROR_DESCRIPTION = "_error_description"
//var UNKNOWN_ERROR_CODE = "UNKNOWN"
//var UNHANDLED_ERROR_CODE = "UNHANDLED"

//var MAX_KEY_LENGTH = 0xff
//var MAX_VALUE_LENGTH = 0xffff

const READ_BUFFER_SIZE int = 65535

var bytes_received = 0

func (c *Client) reader() {    
    buf := make([]byte, READ_BUFFER_SIZE)
    overflow := make([]byte, READ_BUFFER_SIZE)
    left := buf[:0]
    for {
        //log.Println("ready for new read..")
        n, err := c.Conn.Read(buf) 
        //log.Println("received bytes",n)
        if err != nil {
            log.Println("connection reader error!!",err)        
            c.Conn.Close() 
            break    
        }       
        
        bytes_received += n
        
        //log.Println("pre amount left is",left)
        //if len(left) > 0 { log.Println("left...",len(left)) }
        left = c.unpackMaps(append(overflow[:len(left)], buf[:n]...))
        if len(left) > READ_BUFFER_SIZE { log.Fatal(fmt.Sprintf("Client.reader overflow problem with overflow bytes `%d` greater then overflow buffer size of `%d`", len(left), READ_BUFFER_SIZE)) }
        copy(overflow[:len(left)], left[:])
                
        //log.Println("amount left is",left)
        //time.Sleep(100 * time.Millisecond)         
        
        //log.Println("bytes_received",bytes_received)               
    }
}

func clientCreator(name *string, conn *net.TCPConn, prot *AMP) *Client {
    client := &Client{name, conn, prot} 
    go client.reader()
    return client
}

func (c *Client) incomingAsk(data *map[string]string) error {
    m := *data
    if commandName, ok := m[COMMAND]; !ok {
        msg := fmt.Sprintf("Incoming Ask data structure not valid, `%s` not found",COMMAND)
        return errors.New(msg)
    } else { 
        if command_responder, ok := c.prot.getCommandResponder(commandName); !ok {    
            msg := fmt.Sprintf("Incoming Ask command `%s` does not exist",commandName)
            return errors.New(msg)
        } else {            
            ask := resourceAskBox()   
            ask.Args = data            
            ask.client = c
            ask.Response[ANSWER] = m[ASK]                    
            command_responder <- ask
        }
    }
    return nil
}

func (c *Client) incomingAnswer(data *map[string]string) error {
    m := *data            
    tag := m[ANSWER]
    if box, ok := c.prot.getCallback(tag); !ok {
        msg := fmt.Sprintf("callback for incoming answer `%s` not found!!", tag)        
        return errors.New(msg)
    } else {                
        box.Response = data  
        box.Callback <- box
    }
    return nil
}

func (c *Client) handleIncoming(data *map[string]string) {
    m := *data
    if _,ok := m[ASK]; ok {
        err := c.incomingAsk(data)        
        if err != nil { log.Println("error: ",err) }
    } else if _,ok := m[ANSWER]; ok {
        err := c.incomingAnswer(data)        
        if err != nil { log.Println("error: ",err) }
    } else {
        // XXX handle error packets
    }
}

func (c *Client) CallRemote(commandName string, box *CallBox) error {
    tag := <- c.prot.tagger    
    box.Args[ASK] = tag
    box.Args[COMMAND] = commandName
    c.prot.registerCallback(box, tag)
    send := packMap(&box.Args) 
    c.Conn.SetWriteDeadline(time.Now().Add(1e9)) 
    _, err := c.Conn.Write(*send)
    if err != nil {
        neterr, ok := err.(net.Error)
        if ok && neterr.Timeout() {
            log.Println("error callremote",neterr)             
        } else { log.Println(err) }
        return err
    }
    return nil
}

func (ask *AskBox) Reply() error {
    send := packMap(&ask.Response) 
    ask.client.Conn.SetWriteDeadline(time.Now().Add(1e9)) 
    _, err := ask.client.Conn.Write(*send)
    recycleAskBox(ask)
    if err != nil {
        neterr, ok := err.(net.Error)
        if ok && neterr.Timeout() {
            log.Println("error callremote",neterr)             
        } else { log.Println(err) }
        return err
    }
    
    
    return nil
}


