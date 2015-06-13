package gomemcachedclient

import (
    "os"
    "testing"
)

var mc *MemcachedClient

func TestMain(m *testing.M) {
    var err error
    mc, err = NewMemcachedClient("localhost:11211")
    if err != nil {
        os.Exit(1)
    }
    code := m.Run()
    os.Exit(code)
}

func TestSet(t *testing.T) {
    err := mc.Set("k1", "i am k1", 0, 60)
    if err != nil && err != ERR_STORE {
        t.Error("set failed")
    }
}

func TestAdd(t *testing.T) {
    err := mc.Add("k1", "i am k1", 0, 60)
    if err != nil && err != ERR_STORE {
        t.Error("add failed")
    }
}

func TestReplace(t *testing.T) {
    err := mc.Replace("k1", "i am k1", 0, 60)
    if err != nil && err != ERR_STORE {
        t.Error("replace failed")
    }
}

func TestAppend(t *testing.T) {
    err := mc.Append("k1", "i am k1", 0, 60)
    if err != nil && err != ERR_STORE {
        t.Error("append failed")
    }
}

func TestPrepend(t *testing.T) {
    err := mc.Prepend("k1", "i am k1", 0, 60)
    if err != nil && err != ERR_STORE {
        t.Error("prepend failed")
    }
}

func TestGet(t *testing.T) {
    _, _, err := mc.Get("k1")
    if err != nil {
        t.Error("get failed")
    }
}

func TestDelete(t *testing.T) {
    err := mc.Delete("k1")
    if err != nil && err != ERR_DELETE {
        t.Error("delete failed")
    }
}
