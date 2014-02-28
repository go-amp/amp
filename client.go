package amp

import "log"
import "bytes"
import "errors"
import "fmt"
import "net"
import "runtime"
//import "encoding/binary"

const READ_BUFFER_SIZE int = 81920

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
    tag := c.Protocol.GetBoxCounter()
    m := *call.Arguments    
    m[ASK] = tag
    m[COMMAND] = call.Command.Name
    send := PackMap(call.Arguments)    
    _, err := c.Conn.Write(*send)    
    if err != nil { RecycleCallBox(call); log.Println("error sending:",err); return "", err 
        }  else { c.Protocol.AssignCall(tag, call) }    
    return tag, nil
}



func (c *Client) IncomingAnswer(data *map[string]string) error {    
    m := *data        
    //log.Println("Incoming answer..",m)
    tag := m[ANSWER]
    if answer, ok := c.Protocol.GetCall(tag); !ok {
        msg := fmt.Sprintf("callback for incoming answer `%x` not found!!", tag)
        return errors.New(msg)
    } else {                
        answer.Response = data   
        outer:
        for i := 0; i < 10; i++ { 
            select {    
                case answer.Callback <- answer:
                    break outer
                default:                    
                    if i == 9 { log.Println("callback is not responding!!!", answer,"try",i,"buffer length",len(answer.Callback)) }
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
                        if i == 9 { log.Println("command's responder not responding!!",command,ask,"try",i,"buffer length",len(command.Responder)) }                                                
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

func (c *Client) OldReader() {
    buffer := make([]byte, READ_BUFFER_SIZE)
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

func (c *Client) Reader() {
    //buffer := make([][]byte, 1)
    buffer := make([]byte, READ_BUFFER_SIZE)
    //over := make([]byte, READ_BUFFER_SIZE)
    //b_over := make([]byte, READ_BUFFER_SIZE)
    //b_over_length := 0
    //buf_index := 0
    //var index_buf, index_over, l_message, l_overflow int
    var overflow int = 0    
    for {
        readBytes, error := c.Conn.Read(buffer) 
        if error != nil {
            c.Close()         
            break
        }
        
                        
        
        overflow = UnpackMaps(&buffer, readBytes, c.incoming_handler)
        log.Println("received",readBytes,"bytes",error, overflow) 
        if overflow > 0 {
            overflowed := READ_BUFFER_SIZE - overflow
            log.Fatal("overflow ",overflowed)
        }
        //if l_overflow > 0 {
            //copy(over[index_over:index_over+length_remaining], buf[:length_remaining])
            //index_buf = length_remaining
            //UnpackMaps(&over, index_over+length_remaining, c.incoming_handler)
        //} else { index_buf = 0 }
                
        
               
        
        for k := 0; k <= readBytes; k++ {
            buffer[k] = 0x00
        }
    }
}

