package amp

import "encoding/binary"
//import "log"
import "fmt"
import "errors"

const PREFIXLENGTH = 2

//var count = 0

func (c *Client) unpackMaps(buf []byte) []byte {    
    //bytes_used := len(buf)
    //start_count := count
    for {
        item, left, err := getNext(buf)
        //log.Println(len(left),stop)
        
        if item != nil { 
            //count++ 
            c.handleIncoming(item)
        }
        
        if err != nil { 
            //log.Println("unpacked",count,"items","left",len(left))
            //unpacked_count := count - start_count
            //bytes_used = bytes_used - len(left)
            //if unpacked_count == 0 { log.Println("wtf!!!!!!!!!!!!!!!!!!!",err) }
            //if unpacked_count != bytes_used / 47 { log.Println("unpacked_count",unpacked_count,"not right for bytes used",bytes_used,err) }
            return left 
        } else { 
            buf = left
        }
    }
}

func getNext(buf []byte) (*map[string]string, []byte, error) {
    item := *resourceMap()
    i := 0
    length := len(buf)
    for {        
        
        if i + PREFIXLENGTH > length { 
            recycleMap(&item) 
            return nil, buf, errors.New(fmt.Sprintf("stop1 %d",i))
        }

        prefixBytes := []byte{buf[i], buf[i+1]}
        i += PREFIXLENGTH        
        prefix := int(binary.BigEndian.Uint16(prefixBytes))        

        // indicates end of incoming message
        if prefix == 0 {             
            if length == i {
                return &item, buf[i:], errors.New("valid end of incoming message")
            } else {
                return &item, buf[i:], nil
            }
        }                            
        
        // message overflow
        if i + prefix > length { 
            recycleMap(&item)
            return nil, buf, errors.New(fmt.Sprintf("stop2 %d",i))
        }

        // handling key 
        key := string(buf[i:i+prefix])
        i += prefix                
        // message overflow
        if i + PREFIXLENGTH > length { 
            recycleMap(&item)
            return nil, buf, errors.New(fmt.Sprintf("stop3 %d",i))
        }

        // handling value
        prefixBytes = []byte{buf[i], buf[i+1]}
        i += PREFIXLENGTH       
        prefix = int(binary.BigEndian.Uint16(prefixBytes)) 
        
        // message overflow
        if i + prefix > length { 
            recycleMap(&item)
            return nil, buf, errors.New(fmt.Sprintf("stop4 %d",i))
        }
        
        // assigning value    
        value := string(buf[i:i+prefix])
        i += prefix
        item[string(key)] = string(value)
    
    }
}


func packMap(m *map[string]string) *[]byte {
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
