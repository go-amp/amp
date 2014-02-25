package amp

/*
 * implements amp_diagram.svg
 * */
 
import "log"
import "net"
import "fmt"
 
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




func (prot *AMP) connectionListener(netListen net.Listener, service string) {
    clientNum := 0
    defer netListen.Close()
    log.Println("Waiting for clients") 
    for {
        conn, err := netListen.Accept()
        if err != nil {
            log.Println("Client error: ", err)
            break
        } else {
            clientNum += 1
            name := fmt.Sprintf("<%s<-%s>", conn.LocalAddr().String(), conn.RemoteAddr().String())
            log.Println("AMP.connectionListener accepted",name)            
            //log.Println("name is",name)
            ClientCreator(&name, &conn, prot)
            //log.Println("Client created",newClient)            
        }
    }
}

func (prot *AMP) BoxCounterIncrementer() {
    for {
        dispatch := <- prot.GetBoxCounter
        prot.BoxCounter += 1
        dispatch <- prot.BoxCounter
    }
}

func (prot *AMP) ListenTCP(service string) error {
    tcpAddr, err := net.ResolveTCPAddr("tcp", service) 
    if err != nil {
        log.Println("Error: Could not resolve address")
        return err
    } else {
        log.Println("ListenTCP",*tcpAddr)
        netListen, err := net.Listen(tcpAddr.Network(), tcpAddr.String())
        if err != nil {
            log.Println("Error: could not listen")
            return err
        } else {
            go prot.connectionListener(netListen, service)
       }
    }
    return nil
}

func (prot *AMP) ConnectTCP(service string) (*Client, error) {
    //log.Println("ConnectTCP",service)
    conn, err := net.Dial("tcp", service)
    if err != nil {
        log.Println("error!",err)
        return nil, err
    }
    name := fmt.Sprintf("<%s->%s>", conn.LocalAddr().String(), conn.RemoteAddr().String())    
    log.Println("AMP.ConnectTCP connected",name)
        
    newClient := ClientCreator(&name, &conn, prot) 
    //log.Println("name is",newClient)
    return newClient, nil
}

func init() {
}

func Init(commands *map[string]*Command) *AMP {
    connList := make(map[string]*Client)
    boxCounter := 0
    callbacks := make(map[string]*CallBox)    
    prot := &AMP{connList, *commands, boxCounter, callbacks, make(chan chan int)} 
    go prot.BoxCounterIncrementer()
    log.Println("AMP initialized.")   
    return prot
}


