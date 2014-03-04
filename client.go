package amp

import "net"
import "log"
//import "time"

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
    //m := *data
    //if commandName, ok := m[COMMAND]; !ok {
        //msg := fmt.Sprintf("Incoming Ask data structure not valid, `%s` not found",COMMAND)
        //return errors.New(msg)
    //} else { 
        //if command,ok := c.prot.commands[commandName]; !ok {    
            //msg := fmt.Sprintf("Incoming Ask command `%s` does not exist",commandName)
            //return errors.New(msg)
        //} else {
            //ask := resourceAskBox()   
            //ask.Arguments = data
            //response := *resourceMap()
            //response[ANSWER] = m[ASK]
            //ask.Response = &response
            //ask.ReplyChannel = c.reply_handler   
            ////select {          
                ////case command.Responder <- ask:
                ////default:
            ////}
            ////log.Println("buffer size",len(command.Responder))
            //log.Println("sending to responder")
            //command.Responder <- ask
        //}
    //}
    return nil
}

func (c *Client) incomingAnswer(data *map[string]string) error {
    return nil
}

func (c *Client) handleIncoming(data *map[string]string) {
    m := *data
    if _,ok := m[ASK]; ok {
        c.incomingAsk(data)        
    } else if _,ok := m[ANSWER]; ok {
        c.incomingAnswer(data)        
    } else {
        // XXX handle error packets
    }
}

func (c *Client) Dispatch(box *CallBox) error {
    return nil
}


