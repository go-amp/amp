package amp

import "go/cmn"
import "container/list"
import "encoding/binary"

var PREFIXLENGTH = 2

func PrintList(l *list.List) {
    cmn.Log("PrintList..")
    for e := l.Front(); e != nil; e = e.Next() {
        cmn.Log(e.Value)
    }
}

func UnpackMaps(buffer *[]byte, length int) *list.List {
    /*
     * Unpacks N number of maps from a []byte.  Maps are separate by a key length of 0.
     * */
    //cmn.Log("UnpackMap",length)
    b := *buffer
    var i int = 0  
    retList := list.New()
        
    outer: 
        for {
            ret := make(map[string]string)
            for {                
                /* key
                 * */
                prefixBytes := []byte{b[i], b[i+1]}
                i += PREFIXLENGTH        
                prefix := int(binary.BigEndian.Uint16(prefixBytes))        
                if i >= length { retList.PushBack(&ret); break outer }
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
                //cmn.Log("unpacked -",key,":",value)
                ret[string(key)] = string(value)
            }
            //cmn.Log("breaking early",ret)
            retList.PushBack(&ret)                    
        }
    return retList
}

func PackMap(m *map[string]string) *[]byte {
    cmn.Log("packing - ",*m)                       
    length := 0
    for k, v := range *m {         
        length += len(k)
        length += PREFIXLENGTH
        length += len(v)
        length += PREFIXLENGTH
        // 2 is prefixLength
    }
    //cmn.Log("length is",length)
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
        //cmn.Log(buf)
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
        //cmn.Log(buf)
        stop = start + PREFIXLENGTH
        copy(array[start:stop], prefix)
        start = stop
        stop = start + length
        copy(array[start:stop], v)
        start = stop        
    }
    //cmn.Log(array)
    return &array
    
}

func UnpackList(b string) *list.List {
    //cmn.Log("unpacking - ",b)                    
    var i int = 0    
    ret := list.New()
    // list is a linked list
    for {                
        prefixBytes := []byte{b[i], b[i+1]}
        i += PREFIXLENGTH        
        prefix := int(binary.BigEndian.Uint16(prefixBytes))        
        if prefix == 0 { break }
        value := b[i:i+prefix]
        //cmn.Log("string",value)
        i += prefix        
        ret.PushBack(value)
    }
    
    return ret
}

func PackList(l *list.List) *[]byte {
    //cmn.Log("Packing..")
    length := 0
    for e := l.Front(); e != nil; e = e.Next() {
        val := e.Value        
        length += len(val.(string))
        length += PREFIXLENGTH
        // 2 is prefixLength
    }
    var array = make([]byte, length + PREFIXLENGTH)
    start := 0
    stop := 0
    var prefix = make([]byte, PREFIXLENGTH)
    for e := l.Front(); e != nil; e = e.Next() {        
        val := e.Value        
        length = len(val.(string))
        binary.BigEndian.PutUint16(prefix, uint16(length))
        //cmn.Log(buf)
        stop = start + PREFIXLENGTH
        copy(array[start:stop], prefix)
        start = stop
        stop = start + length
        copy(array[start:stop], val.(string))
        start = stop        
    }
    //cmn.Log(array)
    return &array
}
