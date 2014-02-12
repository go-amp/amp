package amp

/*
 * implements amp_diagram.svg
 * */
 
import "go/cmn"
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


func (ask *AskBox) Reply() error {    
    var err error
    err = CheckArgs(&ask.Command.Response, ask.Response)
    if err != nil {
        cmn.Log("reply failed!",err)
        // XXX Need to send error box here
        return err
    }    
    send := PackMap(ask.Response)    
    //UnpackMap(send,len(*send))
    _, err = ask.Client.Conn.Write(*send)    
    // XXX need to handle cleanup of connection for error here
    if err != nil {
        cmn.Log("reply failed!",err)
        return err
    }
    return nil
}

func (prot *AMP) connectionListener(netListen net.Listener, service string) {
    clientNum := 0
    defer netListen.Close()
    cmn.Log("Waiting for clients") 
    for {
        conn, err := netListen.Accept()
        if err != nil {
            cmn.Log("Client error: ", err)
            break
        } else {
            clientNum += 1
            name := fmt.Sprintf("<%s<-%s>", conn.LocalAddr().String(), conn.RemoteAddr().String())
            cmn.Log("AMP.connectionListener accepted",name)
            quitChannel := make(chan bool)
            cmn.Log("name is",name)
            newClient := &Connection{name, conn, prot, quitChannel, false} 
            cmn.Log("Connection created",newClient)
            go newClient.Reader()
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
        cmn.Log("Error: Could not resolve address")
        return err
    } else {
        cmn.Log("ListenTCP",*tcpAddr)
        netListen, err := net.Listen(tcpAddr.Network(), tcpAddr.String())
        if err != nil {
            cmn.Log("Error: could not listen")
            return err
        } else {
            go prot.connectionListener(netListen, service)
       }
    }
    return nil
}

func (prot *AMP) ConnectTCP(service string) (*Connection, error) {
    cmn.Log("ConnectTCP",service)
    conn, err := net.Dial("tcp", service)
    if err != nil {
        cmn.Log("error!",err)
        return nil, err
    }
    name := fmt.Sprintf("<%s->%s>", conn.LocalAddr().String(), conn.RemoteAddr().String())    
    cmn.Log("AMP.ConnectTCP connected",name)
    quitChannel := make(chan bool)    
    newClient := &Connection{name, conn, prot, quitChannel, false} 
    go newClient.Reader()
    cmn.Log("name is",newClient)
    return newClient, nil
}

func Init(commands *map[string]*Command) *AMP {
    connList := make(map[string]*Connection)
    boxCounter := 0
    callbacks := make(map[string]*AnswerBox)    
    prot := &AMP{connList, *commands, boxCounter, callbacks, make(chan chan int)} 
    go prot.BoxCounterIncrementer()
    cmn.Log("AMP initialized.")   
    return prot
}


