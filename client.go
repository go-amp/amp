package amp

import "log"
import "bytes"
import "errors"
import "fmt"
import "net"
import "runtime"

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

func ClientCreator(name *string, conn *net.Conn, prot *AMP) *Client {
    quitChannel := make(chan bool)
    incoming_handler := make(chan *map[string]string)
    reply_handler := make(chan *Ask)
    client := &Client{*name, *conn, prot, quitChannel, false, incoming_handler, reply_handler} 
    go client.Reader()
    go client.IncomingHandler()
    go client.ReplyHandler()
    return client
}

func (c *Client) CallRemote(call *CallBox) (string, error) { 
    fetch := make(chan int)
    c.Protocol.GetBoxCounter <- fetch
    counter := <- fetch
    close(fetch)
    tag := fmt.Sprintf("%x", counter)
    m := *call.Arguments    
    m[ASK] = tag
    m[COMMAND] = call.Command.Name
    send := PackMap(call.Arguments)    
    _, err := c.Conn.Write(*send)    
    if err != nil { RecycleCallBox(call); log.Println("error sending:",err); return "", err }    
    c.Protocol.Callbacks[tag] = call    
    //log.Println("callremote",call)
    return tag, nil
}

func (c *Client) IncomingAnswer(data *map[string]string) error {    
    m := *data        
    //log.Println("Incoming answer..",m)
    tag := m[ANSWER]
    if answer, ok := c.Protocol.Callbacks[tag]; !ok {
        msg := fmt.Sprintf("callback for incoming answer `%x` not found", tag)
        return errors.New(msg)
    } else {                
        answer.Response = data   
        outer:
        for i := 0; i < 10; i++ { 
            select {    
                case answer.Callback <- answer:
                    break outer
                default:
                    log.Println("callback is not responding!!!", answer,"try",i,"buffer length",len(answer.Callback))
                    runtime.Gosched()     
            }
        }        
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
            //log.Println("incoming ask",ask)
            outer:
            for i := 0; i < 10; i++ {
                select {                                      
                    case command.Responder <- ask:
                        //log.Println("sent command ask",ask)
                        break outer
                    default:
                        log.Println("command's responder not responding!!",command,ask,"try",i,"buffer length",len(command.Responder))                             
                        runtime.Gosched()               
                }        
            }
        }
    }
    return nil
}


func (c *Client) ReplyHandler() {
    for {
        ask := <- c.reply_handler   
        send := PackMap(ask.Response)                
        _, err := c.Conn.Write(*send)    
        // XXX should probably close the client if not already if it's an error to send
        if err != nil {
            log.Println("reply failed!",err)            
        }    
        recycleAskBox(ask)
    }
}

func (c *Client) IncomingHandler() {
    for {
        data := <- c.incoming_handler
        m := *data
        if _,ok := m[ASK]; ok {
            err := c.IncomingAsk(data)
            if err != nil { log.Println("error:",err) }
        } else if _,ok := m[ANSWER]; ok {
            err := c.IncomingAnswer(data)
            if err != nil { log.Println("error:",err) }
        } else {
            // XXX handle error packets
        }
    }
}

func (c *Client) Reader() {
    buffer := make([]byte, 81920)
    for {        
        bytesRead, error := c.Conn.Read(buffer)        
        if error != nil {
            c.Close()         
            break
        }                              
        //log.Println("received",bytesRead,"bytes")
        UnpackMaps(&buffer, bytesRead, c.incoming_handler)        
        for i := 0; i <= bytesRead; i++ {
            buffer[i] = 0x00
        }            
    }     
}

