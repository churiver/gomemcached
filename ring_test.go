package gomemcached

import (
    "fmt"
    //"os"
    "strconv"
    "testing"
)

var ring *Ring

func TestMain(m *testing.M) {
    ring = new(Ring)
    ring.nodes = make(map[string]*Node)
}

func TestAddNode(t *testing.T) {
    for i := 1; i < 5; i++ {
        ring.AddNode("127.0.0.1:1121" + strconv.Itoa(i))
    }
}

func TestGetNode(t *testing.T) {
    for i := 0; i < 10; i++ {
        node, _ := ring.GetNode("k" + strconv.Itoa(i))
        fmt.Printf(node.url)
    }
}
