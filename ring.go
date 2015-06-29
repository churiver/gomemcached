package gomemcached
import (
    "crypto/md5"
    "fmt"
    "net"
    "sort"
    "strconv"
    "strings"
)

const REPLIC_NUM int = 3

type Node struct {
    url string
    conn net.Conn
}

func (node *Node) Url() string {
    return node.url
}

func (node *Node) Conn() net.Conn {
    return node.conn
}

type Ring struct {
    nodeMap map[string]*Node
    orderedKeys []string
}

func NewRing() *Ring {
    return &Ring{nodeMap: make(map[string]*Node)}
}

func (ring *Ring) PrintReplics() {
    for k, v := range ring.nodeMap {
        fmt.Printf("Replics. key = %s, url = %s\n", k, v.url)
    }
}

func (ring *Ring) AddNode(url string) (err error) {
    if strings.TrimSpace(url) == "" {
        return nil
    }

    var node *Node = new(Node)
    node.url = url
    node.conn, err = net.Dial("tcp", url)
    if err != nil {
        fmt.Print("Connect failed\n")
        return err
    }

    for i:= 1; i <= REPLIC_NUM; i++ {
        repUrl := url + ":" + strconv.Itoa(i)
        hashVal := fmt.Sprintf("%x", md5.Sum([]byte(repUrl)))
        ring.nodeMap[hashVal] = node
    }

    ring.orderedKeys = nil
    for key, _ := range ring.nodeMap {
    	ring.orderedKeys = append(ring.orderedKeys, key)
    }
    sort.Strings(ring.orderedKeys)

    return nil
}

func (ring *Ring) RemoveNode(url string) (err error) {
    if strings.TrimSpace(url) == "" {
        return nil
    }

    for i:= 1; i <= REPLIC_NUM; i++ {
        repUrl := url + ":" + strconv.Itoa(i)
        hashVal := fmt.Sprintf("%x", md5.Sum([]byte(repUrl)))
        node := ring.nodeMap[hashVal]
        if node != nil {
            err = node.conn.Close()
            if err != nil {
                fmt.Print("Close connection to %s failed\n", url)
            }
        }
        delete (ring.nodeMap, hashVal)
    }

    ring.orderedKeys = nil
    for key := range ring.nodeMap {
        ring.orderedKeys = append(ring.orderedKeys, key)
    }
    sort.Strings(ring.orderedKeys)

    return nil
}

func (ring *Ring) GetNode(key string) (*Node) {
    hashVal := fmt.Sprintf("%x", md5.Sum([]byte(key)))
    pos := 0
    for ; pos < len(ring.orderedKeys) && ring.orderedKeys[pos] < hashVal; pos++ {
    }

    if pos == len(ring.orderedKeys) {
        pos = 0
    }

    hashVal = ring.orderedKeys[pos]
    return ring.nodeMap[hashVal]
}

func (ring *Ring) GetConn(key string) (net.Conn) {
    return ring.GetNode(key).Conn()
}