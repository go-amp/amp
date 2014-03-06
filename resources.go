package amp

//var map_resource chan *map[string]string = make(chan *map[string]string, 100)
//var callbox_resource chan *CallBox = make(chan *CallBox, 100)
var askbox_resource chan *AskBox = make(chan *AskBox, 100)

//func resourceMap() *map[string]string {
    //var m *map[string]string    
    //select {
        //case m = <- map_resource:     
            ////log.Println("reusing map",m)
        //default:        
            //r := make(map[string]string)            
            //m = &r
            ////log.Println("creating new map",m)
    //}
    //return m
//}

//func recycleMap(m *map[string]string) { 
    ////log.Println("recycling map",m)
    //for k, _ := range *m {
        //delete(*m, k)
    //}
    //select {
        //case map_resource <- m:
            //m = nil
        //default:
    //}    
//}

//func ResourceCallBox() *CallBox {
    //var callbox *CallBox
    //select {
        //case callbox = <- callbox_resource:
            ////callbox.Args = resourceMap()
            ////log.Println("reusing callbox",callbox)
            //return callbox
        //default:
            //callbox = &CallBox{make(map[string][]byte), make(map[string][]byte), nil, nil}            
            //return callbox
            ////log.Println("creating new callbox",callbox)
    //}
    
//}

//func RecycleCallBox(callbox *CallBox) {
    
    //for k, _ := range callbox.Args {
        //delete(callbox.Args, k)
    //}
    //for k, _ := range callbox.Response {
        //delete(callbox.Response, k)
    //}
   
    //callbox.Callback = nil
    //callbox.CallbackArgs = nil
    //select {
        //case callbox_resource <- callbox:
            //callbox = nil
        //default:
    //}      
//}

func resourceAskBox() *AskBox {
    var ask *AskBox
    select {
        case ask = <- askbox_resource:
            return ask
        default:
            ask = &AskBox{make(map[string][]byte), make(map[string][]byte), nil}
            return ask
    }
}

func recycleAskBox(ask *AskBox) {    
    for k, _ := range ask.Args {
        delete(ask.Args, k)
    }
    for k, _ := range ask.Response {
        delete(ask.Response, k)
    }
    
    ask.client = nil    
    select {
        case askbox_resource <- ask:
            ask = nil
        default:
    }    
}
