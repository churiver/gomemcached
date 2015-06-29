package gomemcached
import (
    "bytes"
    "errors"
    "strconv"
    "strings"
)

type MemcachedClient struct {
    urlMap map[string]int
    ring *Ring
}

var (
    ERR_CONNECT = errors.New("Failed to connect to any memcached")
    ERR_READ = errors.New("Failed to read from memcached")
    ERR_STORE = errors.New("Failed to store key")
    ERR_GET = errors.New("Failed to get key's value")
    ERR_DELETE = errors.New("Failed to delete key")
    ERR_SERVER_NOT_FOUND = errors.New("Cannot find the server")
    ERR_SERVER_EXIST = errors.New("The server exists already")
    ERR_SERVER_BOUND = errors.New("Only one server left, cannot remove it")
)

func NewMemcachedClient() (*MemcachedClient) {
    return &MemcachedClient{urlMap: make(map[string]int), ring: NewRing()}
}

func (mc *MemcachedClient) AddServer(url string) (error) {
    if mc.urlMap[url] != 0 {
        return ERR_SERVER_EXIST
    }

    err := mc.ring.AddNode(url)
    if err == nil {
        mc.urlMap[url] = 1
    }
    return err
}

func (mc *MemcachedClient) RemoveServer(url string) (error) {
    if len(mc.urlMap) == 1 {
        return ERR_SERVER_BOUND
    }

    if mc.urlMap[url] == 0 {
        return ERR_SERVER_NOT_FOUND
    }

    err := mc.ring.RemoveNode(url)
    if err == nil {
        delete(mc.urlMap, url)
    }
    return err
}

func (mc *MemcachedClient) GetServerNum() (int) {
    return len(mc.urlMap)
}

func (mc *MemcachedClient) store(command, key, value string, flags, exptime int) error {
    var err error

    var buf bytes.Buffer
    buf.WriteString(command + " " + key + " " + strconv.Itoa(flags) + " " +
        strconv.Itoa(exptime) + " " + strconv.Itoa(len(value)) + " " +
        "\r\n" + value + "\r\n")

    conn := mc.ring.GetConn(key)
    _, err = conn.Write(buf.Bytes())
    if err != nil {
        return err
    }

    response := make([]byte, 512)
    _, err = conn.Read(response)
    if strings.TrimSpace(string(response)) == "ERROR" {
        return ERR_STORE
    }
    return err
}

func (mc *MemcachedClient) Set(key, value string, flags, exptime int) error {
    return mc.store("set", key, value, flags, exptime)
}

func (mc *MemcachedClient) Add(key, value string, flags, exptime int) error {
    return mc.store("add", key, value, flags, exptime)
}

func (mc *MemcachedClient) Replace(key, value string, flags, exptime int) error {
    return mc.store("replace", key, value, flags, exptime)
}

func (mc *MemcachedClient) Append(key, value string, flags, exptime int) error {
    return mc.store("append", key, value, flags, exptime)
}

func (mc *MemcachedClient) Prepend(key, value string, flags, exptime int) error {
    return mc.store("prepend", key, value, flags, exptime)
}

func (mc *MemcachedClient) Get(key string) (string, int, error) {
    var value string
    var flags int
    var err error

    var buf bytes.Buffer
    buf.WriteString("get " + key + "\r\n")

    conn := mc.ring.GetConn(key)
    _, err = conn.Write(buf.Bytes())
    if err != nil {
        return value, flags, err
    }

    response := make([]byte, 256)
    _, err = conn.Read(response)
    if err != nil {
        return value, flags, err
    }

    lines := strings.Split(string(response), "\r\n")
    if lines[0] == "END" || len(lines) < 3 {
        err = ERR_GET
    } else {
        flags, _ = strconv.Atoi(strings.Split(lines[0], " ")[2])
        value = lines[1]
    }
    return value, flags, err
}


func (mc *MemcachedClient) Delete(key string) error {
    var err error

    var buf bytes.Buffer
    buf.WriteString("delete " + key + "\r\n")

    conn := mc.ring.GetConn(key)
    _, err = conn.Write(buf.Bytes())
    if err != nil {
        return err
    }

    response := make([]byte, 64)
    _, err = conn.Read(response)
    if err == nil && strings.TrimSpace(string(response)) != "DELETED" {
        err = ERR_DELETE
    }
    return err
}
