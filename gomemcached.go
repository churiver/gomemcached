package gomemcachedclient
import (
    "bytes"
    "errors"
    "fmt"
    "net"
    "strconv"
    "strings"
)

type MemcachedClient struct {
    Urls []string
    Conns []net.Conn
}

var (
    ERR_CONNECT = errors.New("Failed to connect to any memcached")
    ERR_READ = errors.New("Failed to read from memcached")
    ERR_STORE = errors.New("Failed to store key")
    ERR_GET = errors.New("Failed to get key's value")
    ERR_DELETE = errors.New("Failed to delete key")
)

func NewMemcachedClient(serverUrls string) (*MemcachedClient, error) {
    var err error

    var client *MemcachedClient = new(MemcachedClient)
    client.Urls = strings.Split(serverUrls, ",")

    for i := 0; i < len(client.Urls); i++ {
        //fmt.Printf("%s\n", client.Urls[i])
        conn, err := net.Dial("tcp", client.Urls[i])
        if err != nil {
            fmt.Printf("Failed to connect to %s\n", client.Urls[i])
            continue
        }
        fmt.Printf("Connected to %s\n", client.Urls[i])
        client.Conns = append(client.Conns, conn)
    }
    if len(client.Conns) == 0 {
        err = ERR_CONNECT
        //fmt.Printf("No memcached found. Exiting...\n")
        //os.Exit(1)
    }
    return client, err
}

/* 
    TODO
    1. Add support for noreply in args
*/
func (client *MemcachedClient) store(command, key, value string, flags, exptime int) error {
    var err error

    var buf bytes.Buffer
    buf.WriteString(command + " " + key + " " + strconv.Itoa(flags) + " " +
        strconv.Itoa(exptime) + " " + strconv.Itoa(len(value)) + " " +
        "\r\n" + value + "\r\n")
    //fmt.Printf("Debug. store: %s\n", buf.String())

    _, err = client.Conns[0].Write(buf.Bytes())
    if err != nil {
        //fmt.Print("Set: write to server failed\n")
        return err
    }

    response := make([]byte, 512)
    _, err = client.Conns[0].Read(response)
    //fmt.Printf("store: response from server: %s\n", string(response))
    if strings.TrimSpace(string(response)) == "ERROR" {
        err = ERR_STORE
    }
    return err
}

func (client *MemcachedClient) Set(key, value string, flags, exptime int) error {
    return client.store("set", key, value, flags, exptime)
}

func (client *MemcachedClient) Add(key, value string, flags, exptime int) error {
    return client.store("add", key, value, flags, exptime)
}

func (client *MemcachedClient) Replace(key, value string, flags, exptime int) error {
    return client.store("replace", key, value, flags, exptime)
}

func (client *MemcachedClient) Append(key, value string, flags, exptime int) error {
    return client.store("append", key, value, flags, exptime)
}

func (client *MemcachedClient) Prepend(key, value string, flags, exptime int) error {
    return client.store("prepend", key, value, flags, exptime)
}

func (client *MemcachedClient) Get(key string) (string, int, error) {
    var value string
    var flags int
    var err error

    var buf bytes.Buffer
    buf.WriteString("get " + key + "\r\n")
    //fmt.Printf("Debug. Get: %s\n", buf.String())

    _, err = client.Conns[0].Write(buf.Bytes())
    if err != nil {
        //fmt.Printf("Get: Write to %s failed\n", client.Urls[0])
        return value, flags, err
    }

    response := make([]byte, 256)
    _, err = client.Conns[0].Read(response)
    if err != nil {
        //fmt.Printf("Get: Read from %s failed\n", client.Urls[0])
        return value, flags, err
    }

    //fmt.Printf("Debug. Get: value = %s\n", string(response))
    lines := strings.Split(string(response), "\r\n")
    if lines[0] == "END" || len(lines) < 3 {
        err = ERR_GET
    } else {
        flags, _ = strconv.Atoi(strings.Split(lines[0], " ")[2])
        value = lines[1]
    }
    return value, flags, err
}

/* TODO func GetMulti */

func (client *MemcachedClient) Delete (key string) error {
    var err error

    var buf bytes.Buffer
    buf.WriteString("delete " + key + "\r\n")

    _, err = client.Conns[0].Write(buf.Bytes())
    if err != nil {
        return err
    }

    response := make([]byte, 64)
    _, err = client.Conns[0].Read(response)
    if err == nil && strings.TrimSpace(string(response)) != "DELETED" {
        err = ERR_DELETE
    }
    return err
}
