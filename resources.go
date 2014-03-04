package amp

var map_resource chan *map[string]string = make(chan *map[string]string, 100)
var callbox_resource chan *CallBox = make(chan *CallBox, 100)

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

func ResourceCallBox() *CallBox {
    var callbox *CallBox
    select {
        case callbox = <- callbox_resource:
            //callbox.Args = resourceMap()
            //log.Println("reusing callbox",callbox)
            return callbox
        default:
            callbox = &CallBox{make(map[string]string), nil, nil, nil, nil}            
            return callbox
            //log.Println("creating new callbox",callbox)
    }
    
}

func RecycleCallBox(callbox *CallBox) {
    //log.Println("recycling callbox",callbox)
    //if callbox.Arguments != nil {
        //recycleMap(callbox.Arguments)
        //callbox.Arguments = nil
    //}
    for k, _ := range callbox.Args {
        delete(callbox.Args, k)
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
