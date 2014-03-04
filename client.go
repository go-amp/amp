package amp

import "net"
import "log"
//import "time"

const READ_BUFFER_SIZE int = 65535

var bytes_received = 0

func (c *Client) Reader() {    
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
        left = UnpackMaps(append(overflow[:len(left)], buf[:n]...))
        copy(overflow[:len(left)], left[:])
        
        //log.Println("amount left is",left)
        //time.Sleep(100 * time.Millisecond)         
        
        //log.Println("bytes_received",bytes_received)               
    }
}

func ClientCreator(name *string, conn *net.TCPConn) *Client {
    client := &Client{name, conn} 
    go client.Reader()
    return client
}


type Client struct {
    Name *string
    Conn *net.TCPConn
}
