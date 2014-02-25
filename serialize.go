package amp

//import "log"
//import "container/list"
import "encoding/binary"

var PREFIXLENGTH = 2

func UnpackMaps(buffer *[]byte, length int, incoming_handler chan *map[string]string)  {
    /*
     * Unpacks N number of maps from a []byte.  Maps are separate by a key length of 0.
     * */
    //log.Println("UnpackMap",length)
    b := *buffer
    var i int = 0  
    //retList := list.New()
        
    outer: 
        for {
            //ret := make(map[string]string)
            ret := *resourceMap()
            for {                
                /* key
                 * */
                prefixBytes := []byte{b[i], b[i+1]}
                i += PREFIXLENGTH        
                prefix := int(binary.BigEndian.Uint16(prefixBytes))        
                if i >= length { incoming_handler <- &ret; break outer }
                if prefix == 0 { break }             
                key := string(b[i:i+prefix])
                i += prefix
                /* value
                 * */
                prefixBytes = []byte{b[i], b[i+1]}
                i += PREFIXLENGTH       
                prefix = int(binary.BigEndian.Uint16(prefixBytes))        
                if i >= length { break outer }
                value := string(b[i:i+prefix])
                i += prefix
                //log.Println("unpacked -",key,":",value)
                ret[string(key)] = string(value)
            }
            //log.Println("breaking early",ret)
            //retList.PushBack(&ret)                    
            incoming_handler <- &ret
        }
    //return retList
}

func PackMap(m *map[string]string) *[]byte {
    //log.Println("packing - ",m)                       
    length := 0
    for k, v := range *m {         
        length += len(k)
        length += PREFIXLENGTH
        length += len(v)
        length += PREFIXLENGTH
        // 2 is prefixLength
    }
    //log.Println("length is",length)
    var array = make([]byte, length + PREFIXLENGTH)
    /*
     * 2 null terminating bytes 
     * - A single NUL will separate every key, and a double NUL separates
      messages.  This provides some redundancy when debugging traffic dumps.
      * */
    start := 0
    stop := 0
    var prefix = make([]byte, PREFIXLENGTH)
    for k, v := range *m {                 
        /* for key
         * */
        length = len(k)
        binary.BigEndian.PutUint16(prefix, uint16(length))
        //log.Println(buf)
        stop = start + PREFIXLENGTH
        copy(array[start:stop], prefix)
        start = stop
        stop = start + length
        copy(array[start:stop], k)
        start = stop        
        /* now for value
         * */
        length = len(v)
        binary.BigEndian.PutUint16(prefix, uint16(length))
        //log.Println(buf)
        stop = start + PREFIXLENGTH
        copy(array[start:stop], prefix)
        start = stop
        stop = start + length
        copy(array[start:stop], v)
        start = stop        
    }
    //log.Println(array)
    return &array
    
}


