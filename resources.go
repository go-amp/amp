package amp

var map_resource chan *map[string]string = make(chan *map[string]string, 100)

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
