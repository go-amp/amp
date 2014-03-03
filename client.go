package amp

import "log"
import "bytes"
import "errors"
import "fmt"
import "net"
//import "time"
//import "runtime"
//import "encoding/binary"

const READ_BUFFER_SIZE int = 65535

func (c *Client) Close() {
    //log.Println("Closing connection")
    // XXX need to close the channels
    c.Conn.Close()    
}

func (c *Client) Equal(other *Client) bool {
    if bytes.Equal([]byte(c.Name), []byte(other.Name)) {
        if c.Conn == other.Conn {
            return true
        }
    }
    return false
}

// removing for simplicity for now
//func CheckArgs(args *[]string, data *map[string]string) error {
    //var arg string
    //m := *data
    //for _,arg = range *args {
        //if _,ok := m[arg]; !ok {                            
            //msg := fmt.Sprintf("Found missing argument `%s`",arg)
            //return errors.New(msg)
        //}
    //}
    //return nil
//}

func ClientCreator(name *string, conn *net.TCPConn, prot *AMP) *Client {
    quitChannel := make(chan bool)
    incoming_handler := make(chan *map[string]string)
    reply_handler := make(chan *Ask)
    writer := make(chan *[]byte)
    client := &Client{*name, *conn, prot, quitChannel, false, incoming_handler, reply_handler, writer} 
    
    go client.Reader()
    go client.IncomingHandler()
    go client.ReplyHandler()
    go client.IOHandler()
    return client
}

func (c *Client) CallRemote(call *CallBox) (string, error) {     
    tag := <- boxcounter_producer 
    m := *call.Arguments    
    m[ASK] = tag
    m[COMMAND] = call.Command.Name
    call_mutex.Lock()    
    c.Protocol.Callbacks[tag] = call    
    call_mutex.Unlock()    
    send := PackMap(call.Arguments)    
    log.Println("sending to write",send)
    c.writer <- send
    //_, err := c.Conn.Write(*send)        
    
    //if err != nil {         
        //neterr, ok := err.(net.Error)
        //if ok && neterr.Timeout() {
            //log.Panic("error callremote",neterr)             
        //} else { log.Panic(err) }
        //call_mutex.Lock()    
        //delete(c.Protocol.Callbacks, tag)
        //call_mutex.Unlock()
        //RecycleCallBox(call)       
        //return "", err  
    //} 
    return tag, nil
}

func (c *Client) IOHandler() {
    for {
        send := <- c.writer
        //c.Conn.SetWriteDeadline(time.Now().Add(1e9))
        //test := []byte("hello")
        log.Println("writing",send)
        _, err := c.Conn.Write(*send)     
        if err != nil {         }
            //neterr, ok := err.(net.Error)
            //if ok && neterr.Timeout() {
                //log.Panic("error IOHandler: ",neterr)             
            //} else { log.Panic(err) }
        //}
    }
}



func (c *Client) IncomingAnswer(data *map[string]string) error {    
    m := *data            
    tag := m[ANSWER]    
    if answer, ok := c.Protocol.GetCall(tag); !ok {
        msg := fmt.Sprintf("callback for incoming answer `%s` not found!!", tag)        
        return errors.New(msg)
    } else {                
        answer.Response = data  
        //select { 
            //case answer.Callback <- answer:
            //default:
        //}
        //answer.Callback <- answer
    }
    return nil
}

func (c *Client) IncomingAsk(data *map[string]string) error {
    m := *data
    if commandName, ok := m[COMMAND]; !ok {
        msg := fmt.Sprintf("Incoming Ask data structure not valid, `%s` not found",COMMAND)
        return errors.New(msg)
    } else { 
        if command,ok := c.Protocol.Commands[commandName]; !ok {    
            msg := fmt.Sprintf("Incoming Ask command `%s` does not exist",commandName)
            return errors.New(msg)
        } else {
            ask := resourceAskBox()   
            ask.Arguments = data
            response := *resourceMap()
            response[ANSWER] = m[ASK]
            ask.Response = &response
            ask.ReplyChannel = c.reply_handler   
            //select {          
                //case command.Responder <- ask:
                //default:
            //}
            //log.Println("buffer size",len(command.Responder))
            log.Println("sending to responder")
            command.Responder <- ask
        }
    }
    return nil
}


func (c *Client) ReplyHandler() {
    for {
        ask := <- c.reply_handler   
        log.Println("in replyhandler")
        send := PackMap(ask.Response)          
        c.writer <- send
        //c.Conn.SetWriteDeadline(time.Now().Add(1e9))      
        //_, err := c.Conn.Write(*send)    
        // XXX should probably close the client if not already if it's an error to send
        //if err != nil {
            
            //neterr, ok := err.(net.Error)
            //if ok && neterr.Timeout() {
                //log.Panic("error callremote",err)             
            //} else { log.Panic("reply failed!",err) }
        //}    
        recycleAskBox(ask)
    }
}

func (c *Client) IncomingHandler() {
    for {
        data := <- c.incoming_handler
        log.Println("received data in IcomingHandler")
        m := *data
        if _,ok := m[ASK]; ok {
            err := c.IncomingAsk(data)
            if err != nil { log.Println("error:",err,m) }
        } else if _,ok := m[ANSWER]; ok {
            err := c.IncomingAnswer(data)
            if err != nil { log.Println("error:",err,m) }            
        } else {
            // XXX handle error packets
        }
    }
}

func (c *Client) Reader() {    
    buffer := make([]byte, READ_BUFFER_SIZE)
    packet_slice := make([]byte, 0)
    overflow_slice := make([]byte, 0)
    var overflow int = 0    
    //var readBytes int
    //defer func() {
        //if r := recover(); r != nil {
            //log.Fatal("Recovered in f", r, b, i, readBytes, message_start)
        //}
    //}()
    for {
        log.Println("ready for new read..")
        readBytes, error := c.Conn.Read(buffer) 
        log.Println("received bytes",readBytes)
        if error != nil {
            log.Println("connection reader error!!",error)
            c.Close()                     
            break
        }
        //log.Println("received",readBytes,error)
        // this is probably slow as fuck but here we go
        time.Sleep(.1 * time.Second)
        //packet_slice = append(overflow_slice, buffer[:readBytes]...)        
        ////overflow = UnpackMaps(&packet_slice, len(packet_slice), c.incoming_handler)        
            
        //if overflow > 0 {            
            //overflow_slice = packet_slice[overflow:]            
        //} else if len(overflow_slice) > 0 {
            //overflow_slice = packet_slice[0:0]
        //}
        
                           
    }
}

