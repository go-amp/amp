package amp

//import "log" 

var map_resource chan *map[string]string = make(chan *map[string]string, 5)
var askbox_resource chan *Ask = make(chan *Ask, 5)
var callbox_resource chan *CallBox = make(chan *CallBox, 5)

func resourceMap() *map[string]string {
    var m *map[string]string    
    select {
        case m = <- map_resource:     
            //log.Println("reusing map",m)
        default:        
            r := make(map[string]string)            
            m = &r
            //log.Println("creating new map",m)
    }
    return m
}

func recycleMap(m *map[string]string) { 
    //log.Println("recycling map",m)
    for k, _ := range *m {
        delete(*m, k)
    }
    select {
        case map_resource <- m:
            m = nil
        default:
    }    
}

func resourceAskBox() *Ask {
    var ask *Ask
    select {
        case ask = <- askbox_resource:
            //log.Println("reusing askbox",ask)
        default:
            ask = &Ask{nil, nil, nil}
            //log.Println("creating new askbox",ask)
    }
    return ask
}

func recycleAskBox(ask *Ask) {
    //log.Println("recycling ask",ask)
    if ask.Arguments != nil {
        recycleMap(ask.Arguments)
        ask.Arguments = nil
    }
    if ask.Response != nil {
        recycleMap(ask.Response)
        ask.Response = nil
    }
    ask.ReplyChannel = nil
    select {
        case askbox_resource <- ask:
            ask = nil
        default:
    }    
}

func ResourceCallBox() *CallBox {
    var callbox *CallBox
    select {
        case callbox = <- callbox_resource:
            callbox.Arguments = resourceMap()
            //log.Println("reusing callbox",callbox)
        default:
            callbox = &CallBox{nil, nil, nil, nil, nil}
            callbox.Arguments = resourceMap()
            //log.Println("creating new callbox",callbox)
    }
    return callbox
}

func RecycleCallBox(callbox *CallBox) {
    //log.Println("recycling callbox",callbox)
    if callbox.Arguments != nil {
        recycleMap(callbox.Arguments)
        callbox.Arguments = nil
    }
    if callbox.Response != nil {
        recycleMap(callbox.Response)
        callbox.Response = nil
    }
    callbox.Command = nil
    callbox.Callback = nil
    callbox.CallbackArgs = nil
    select {
        case callbox_resource <- callbox:
            callbox = nil
        default:
    }      
}
