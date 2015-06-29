package gomemcached

import (
    "fmt"
    "strconv"
    "testing"
)

var ring *Ring

func init() {
    ring = new(Ring)
    ring.nodeMap = make(map[string]*Node)
}

func TestAddNode(t *testing.T) {
    for i := 1; i < 5; i++ {
        ring.AddNode("127.0.0.1:1121" + strconv.Itoa(i))
    }
}

func TestGetNode(t *testing.T) {
    nmap := make(map[string]int)
    
    for i := 0; i < 400; i++ {
        key := "k" + strconv.Itoa(i)
        node := ring.GetNode(key)
        nmap[node.Url()] += 1
    }

    for k, v := range nmap {
        fmt.Printf("%s: %d\n", k, v)
    }
}
