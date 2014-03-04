package amp

import "net"
import "log"
import "time"

const READ_BUFFER_SIZE int = 65535

func (c *Client) Reader() {    
    buffer := make([]byte, READ_BUFFER_SIZE)
    for {
        //log.Println("ready for new read..")
        n, err := c.Conn.Read(buffer) 
        log.Println("received bytes",n)
        if err != nil {
            log.Println("connection reader error!!",err)        
            c.Conn.Close() 
            break    
        }        
        time.Sleep(100 * time.Millisecond)                        
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
