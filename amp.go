package amp

/*
 * implements amp_diagram.svg
 * */
 
import "log"
import "net"
import "fmt"
import "sync"
 
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


var boxcounter_mutex = &sync.Mutex{}
var call_mutex = &sync.Mutex{}
var boxcounter_producer = make(chan string)
var call_deregistration = make(chan string)
var call_registration = make(chan *CallBox)
var call_dispatch = make(chan *map[string]string)

func (prot *AMP) connectionListener(netListen *net.TCPListener, service string) {
    clientNum := 0
    defer netListen.Close()
    log.Println("Waiting for clients") 
    for {
        conn, err := netListen.AcceptTCP()
        if err != nil {
            log.Println("Client error: ", err)
            break
        } else {
            clientNum += 1
            name := fmt.Sprintf("<%s<-%s>", conn.LocalAddr().String(), conn.RemoteAddr().String())
            log.Println("AMP.connectionListener accepted",name)            
            //log.Println("name is",name)
            ClientCreator(&name, conn, prot)
            //log.Println("Client created",newClient)            
        }
    }
}

//func (prot *AMP) CallRegistrar() {
    //for {
        //select {
            //case call := <- call_registration:
                //m := *call.Arguments    
                //tag := m[ASK]
                //prot.Callbacks[tag] = call
            //case data := <- call_dispatch:
                //m := *data
                //tag := m[ANSWER]
                //answer, ok := prot.Callbacks[tag]
                //if !ok { 
                    //log.Println(fmt.Sprintf("callback for incoming answer `%s` not found!!", tag)) 
                //} else {
                    //answer.Response = data   
                    //answer.Callback <- answer
                    //delete(prot.Callbacks, tag)
                //}
                
            //case tag := <- call_deregistration:
                //delete(prot.Callbacks, tag)
        //}
    //}
//}

func (prot *AMP) TagProduction() {
    for {
        prot.BoxCounter += 1
        tag := fmt.Sprintf("%x", prot.BoxCounter)  
        boxcounter_producer <- tag
    }
}

func (prot *AMP) GetCall(tag string) (*CallBox, bool) {
    call_mutex.Lock()    
    call, ok := prot.Callbacks[tag]
    delete(prot.Callbacks, tag)    
    call_mutex.Unlock()
    return call, ok    
}

func (prot *AMP) ListenTCP(service string) error {
    tcpAddr, err := net.ResolveTCPAddr("tcp", service) 
    if err != nil {
        log.Println("Error: Could not resolve address")
        return err
    } else {
        log.Println("ListenTCP",*tcpAddr)
        netListen, err := net.ListenTCP(tcpAddr.Network(), tcpAddr)
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
    //conn, err := net.Dial("tcp", service)
    serverAddr, err := net.ResolveTCPAddr("tcp", service)
    if err != nil {
        log.Println("error!",err)
        return nil, err
    }
    conn, err := net.DialTCP("tcp", nil, serverAddr)
    if err != nil {
        log.Println("error!",err)
        return nil, err
    }
    name := fmt.Sprintf("<%s->%s>", conn.LocalAddr().String(), conn.RemoteAddr().String())    
    log.Println("AMP.ConnectTCP connected",name)
        
    newClient := ClientCreator(&name, conn, prot)     
    return newClient, nil
}


func Init(commands *map[string]*Command) *AMP {
    connList := make(map[string]*Client)
    boxCounter := 0
    callbacks := make(map[string]*CallBox)    
    prot := &AMP{connList, *commands, boxCounter, callbacks}     
    //go prot.CallRegistrar()
    go prot.TagProduction()    
    log.Println("AMP initialized.")   
    return prot
}


