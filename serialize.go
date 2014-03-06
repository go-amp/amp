package amp

import "encoding/binary"
//import "log"
import "fmt"
import "errors"

const PREFIXLENGTH = 2

//var count = 0

func scan(reader *bufio.Reader, v []byte) error {
    i := 0
    for {
        n, err := reader.Read(v[i:])
        if err != nil { return err }
        i += n
        if i == len(v) { return nil }
    }
    return nil
}

func get(reader *bufio.Reader, m map[string][]byte) error {
    prefix := make([]byte, 2)
    l := 0
    var err error
    for {                
        //k        
        err = scan(reader, prefix[:])        
        if err != nil { return err }
        l = int(binary.BigEndian.Uint16(prefix)) 
        // indicates end of message 
        if l == 0 { return nil }              
        
        k := make([]byte, l)
        err = scan(reader, k[:])                
        if err != nil { return err }
        
        //v        
        err = scan(reader, prefix[:])                
        if err != nil { return err }
        l = int(binary.BigEndian.Uint16(prefix))        
                
        v := make([]byte, l)
        err = scan(reader, v[:])                
        if err != nil { return err }
        
        // assign
        m[string(k)] = v        
    }
    return nil
}


func pack(m map[string][]byte) *[]byte {
    l := 0
    for k, v := range m {         
        l += len(k)
        l += PREFIXLENGTH
        l += len(v)
        l += PREFIXLENGTH        
    }
    
    var r = make([]byte, l + PREFIXLENGTH)
    i := 0
    for k, v := range m {
        //k        
        l = len(k)
        binary.BigEndian.PutUint16(r[i:i+PREFIXLENGTH], uint16(l))
        i += PREFIXLENGTH
        copy(r[i:i+l],k)
        i += l
        //v        
        l = len(v)
        binary.BigEndian.PutUint16(r[i:i+PREFIXLENGTH], uint16(l))
        i += PREFIXLENGTH
        copy(r[i:i+l],v)
        i += l
    }
    return &r
}
