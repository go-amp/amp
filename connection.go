package amp

import "go/cmn"
import "bytes"
import "errors"
import "fmt"

func (c *Connection) Close() {
    cmn.Log("Closing connection")    
    c.Conn.Close()    
}

func (c *Connection) Equal(other *Connection) bool {
    if bytes.Equal([]byte(c.Name), []byte(other.Name)) {
        if c.Conn == other.Conn {
            return true
        }
    }
    return false
}

func CheckArgs(args *[]string, data *map[string]string) error {
    var arg string
    m := *data
    for _,arg = range *args {
        if _,ok := m[arg]; !ok {                            
            msg := fmt.Sprintf("Found missing argument `%s`",arg)
            return errors.New(msg)
        }
    }
    return nil
}

func (c *Connection) CallRemote(command *Command, args *map[string]string, callback chan *AnswerBox) (string, error) {    
    m := *args        
    fetch := make(chan int)
    c.Protocol.GetBoxCounter <- fetch
    counter := <- fetch
    tag := fmt.Sprintf("%x", counter)
    m[ASK] = tag
    m[COMMAND] = command.Name
    cmn.Log("CallRemote",m)
    answer := &AnswerBox{nil, nil, command, callback}
    // XXX need to add callback variables eventually
    c.Protocol.Callbacks[tag] = answer
    send := PackMap(args)    
    _, err := c.Conn.Write(*send)    
    if err != nil { close(callback); cmn.Log("error sending:",err); return "", err }    
    return tag, nil
}

func (c *Connection) IncomingAnswer(data *map[string]string) error {
    m := *data    
    cmn.Log("IncomingAnswer",m)
    tag := m[ANSWER]
    if answer, ok := c.Protocol.Callbacks[tag]; !ok {
        msg := fmt.Sprintf("callback for incoming answer `%x` not found", tag)
        return errors.New(msg)
    } else {
        err := CheckArgs(&answer.Command.Response, data)
        if err != nil {
            return err
        } else {
            //delete(m, ANSWER)
            answer.Response = data
            select {
                case answer.Callback <- answer:
                    //cmn.Log("sent ask to callback")                                
                default:
                    msg := fmt.Sprintf("Incoming Answer command `%s`'s callback is not responding.",answer.Command.Name)
                    return errors.New(msg)
            }                
        }
    }
    return nil
}

func (c *Connection) IncomingAsk(data *map[string]string) error {
    m := *data
    if commandName, ok := m[COMMAND]; !ok {
        msg := fmt.Sprintf("Incoming Ask data structure not valid, `%s` not found",COMMAND)
        return errors.New(msg)
    } else { 
        if command,ok := c.Protocol.Commands[commandName]; !ok {    
            msg := fmt.Sprintf("Incoming Ask command `%s` does not exist",commandName)
            return errors.New(msg)
        } else {
            err := CheckArgs(&command.Arguments, data)
            if err != nil {
                return err
            } else {
                answerData := make(map[string]string)
                answerData[ANSWER] = m[ASK]
                ask := &AskBox{data, &answerData, c, command}                
                select {
                    case command.Responder <- ask:
                        // do nothing
                    default:
                        msg := fmt.Sprintf("Responder for `%s`'s cis not responding.",command.Name)
                        return errors.New(msg)
                }
                        
            }
        }
    }
    return nil
}

func (c *Connection) Reader() {
    buffer := make([]byte, 81920)
    for {
        /* XXX need to make this multipart (what if the packet is bigger then 81k?)
         * Assuming we have full receipt for now, but we need to do multiple passes until the end of message is found
         * */
        bytesRead, error := c.Conn.Read(buffer)        
        if error != nil {
            c.Close()
            cmn.Log("c.Conn.Read error -",error)
            break
        }
        //cmn.Log("Read ", bytesRead, " bytes:",string(buffer[:bytesRead]))  
                              
        dataList := UnpackMaps(&buffer, bytesRead)
        //PrintList(dataList)
        // perhaps unpackmap should throw an error
        for i := 0; i <= bytesRead; i++ {
            buffer[i] = 0x00
        }    
        for e := dataList.Front(); e != nil; e = e.Next() {
            data := e.Value.(*map[string]string)
            m := *data
            
            if _,ok := m[ASK]; ok {
                err := c.IncomingAsk(data)
                if err != nil {
                    cmn.Log(err)
                }
            } else if _,ok := m[ANSWER]; ok {
                err := c.IncomingAnswer(data)
                if err != nil {
                    cmn.Log(err)
                }
            } else {
                cmn.Log("got packet that does not make sense",m)
            }
        }
    } 
    cmn.Log("ClientReader stopped for ", c.Name)
}

