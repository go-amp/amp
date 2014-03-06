package amp

import "net"
import "log"
import "time"
import "fmt"
import "errors"
import "bufio"

var ASK = "_ask"
var ANSWER = "_answer"
var COMMAND = "_command"
/* not using these yet, as they are used uncomment 
var ERROR = "_error"
var ERROR_CODE = "_error_code"
var ERROR_DESCRIPTION = "_error_description"
var UNKNOWN_ERROR_CODE = "UNKNOWN"
var UNHANDLED_ERROR_CODE = "UNHANDLED"

var MAX_KEY_LENGTH = 0xff
var MAX_VALUE_LENGTH = 0xffff
*/

var bytes_received = 0

func (c *Client) incoming() {        
    var err error
    for {    
        m := make(map[string][]byte)
        err = get(c.reader, m)
        if err != nil { log.Println(err); break }
        
        // handle m
        if _,ok := m[ASK]; ok {
            err = c.incomingAsk(m)        
            if err != nil { log.Println("error: ",err) }
        } else if _,ok := m[ANSWER]; ok {
            err = c.incomingAnswer(m)        
            if err != nil { log.Println("error: ",err) }
        } else {
            // XXX handle error packets
        }
    } 
    
}

func clientCreator(name *string, conn *net.TCPConn, prot *AMP) *Client {    
    writer := bufio.NewWriter(conn)
    reader := bufio.NewReader(conn)
    client := &Client{name, conn, prot, writer, reader} 
    go client.incoming()
    return client
}

func (c *Client) incomingAsk(m map[string][]byte) error {
    
    if commandName, ok := m[COMMAND]; !ok {
        msg := fmt.Sprintf("Incoming Ask data structure not valid, `%s` not found",COMMAND)
        return errors.New(msg)
    } else { 
        if command_responder, ok := c.prot.getCommandResponder(string(commandName)); !ok {    
        
            msg := fmt.Sprintf("Incoming Ask command `%s` does not exist",commandName)
            return errors.New(msg)
        } else {            
            ask := resourceAskBox()   
            ask.Args = m            
            ask.client = c
            ask.Response[ANSWER] = m[ASK]                    
            command_responder <- ask
        }
    }
    return nil
}

func (c *Client) incomingAnswer(m map[string][]byte) error {
         
    tag := string(m[ANSWER])
    if box, ok := c.prot.getCallback(tag); !ok {
        
        msg := fmt.Sprintf("callback for incoming answer `%s` not found!!", tag)        
        return errors.New(msg)
    } else {                
        box.Response = m
        box.Callback <- box
    }
    return nil
}


func (c *Client) CallRemote(commandName string, box *CallBox) error {
    tag := <- c.prot.tagger    
    box.Args[ASK] = []byte(tag)
    box.Args[COMMAND] = []byte(commandName)
    c.prot.registerCallback(box, tag)
    buf := *pack(box.Args)     
    _, err := c.writer.Write(buf)    
    c.writer.Write.Flush()
    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}

func (ask *AskBox) Reply() error {
    buf := *pack(ask.Response) 
    //ask.client.Conn.SetWriteDeadline(time.Now().Add(1e9)) 
    //_, err := ask.client.Conn.Write(*send)
    _, err := ask.client.writer.Write(buf) 
    ask.client.writer.Write.Flush()   
    recycleAskBox(ask)
    if err != nil {
        log.Println(err) 
        return err
    }
    
    
    return nil
}


