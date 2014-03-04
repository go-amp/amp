package amp

import "net"
import "log"
import "fmt"
import "sync"

func (prot *AMP) connectionListener(netListen *net.TCPListener, service string) {    
    defer netListen.Close()
    log.Println("Waiting for clients") 
    for {
        conn, err := netListen.AcceptTCP()
        if err != nil {
            log.Println("Client error: ", err)
            break
        } else {            
            name := fmt.Sprintf("<%s<-%s>", conn.LocalAddr().String(), conn.RemoteAddr().String())
            log.Println("AMP.connectionListener accepted",name)            
            
            clientCreator(&name, conn, prot)
            
        }
    }
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
        
    newClient := clientCreator(&name, conn, prot)     
    return newClient, nil
}

func (prot *AMP) RegisterResponder(name string, responder chan *AskBox) {
    prot.commands_mutex.Lock()
    prot.commands[name] = responder
    prot.commands_mutex.Unlock()
}

func Init() *AMP { 
    prot := &AMP{make(map[string]chan *AskBox), make(map[string]*CallBox), &sync.Mutex{}, &sync.Mutex{}}
    return prot
}
