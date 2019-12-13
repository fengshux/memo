# Memo is a memory cache for golang.


## Installation

```bash
    go get github.com/fengshux/memo
```

## Usage


```golang
    // init
    m := memo.Default()
    // set
    m.Set("key",val)
    // get
    res := m.Get("key")
    // setex
    m.SetEx("key", 5*time.Second, val)
    //delete
    m.Del("key")

```
