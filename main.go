package main

import (
    "fmt"
    "flag"
    "net"
    "os"
)

type Cmd byte

func (c Cmd) String() string {
    r := "Cmd["
    switch c {
        case QUIT:
            return r + "quit" + "]"
        case REBUILD:
            return r + "rebuild" + "]"
        case RELOAD:
            return r + "reload" + "]"
        case ERROR:
            fallthrough
        default:
    }
    return r + "error" + "]"
}

const (
    QUIT Cmd = iota
    REBUILD
    RELOAD
    ERROR
)

func openListener() int {
    l, err := net.Listen("unix", "/tmp/reloader-test")
    defer l.Close()
    if err != nil {
        fmt.Println("Error ListenUnix:", err)
        return ERROR
    }

    c, err := l.Accept()
    defer c.Close()
    if err != nil {
        fmt.Println("Error Accepting Conn:", err)
        return ERROR
    }

    msg := make([]byte, 0, 1024)
    n, err := c.Read(msg[:1])
    if err != nil && err != os.EOF {
        fmt.Println("Error Reading from Conn:", err)
        return ERROR
    }

    fmt.Println("Read:", n, "bytes")
    fmt.Println(msg[:n])

    m := msg[:n][0];

    return int(m)
}

func sendMsg(msg int) {
    c, err := net.Dial("unix", "/tmp/reloader-test")
    defer c.Close()
    if err != nil {
        fmt.Println("Error Dialing:", "/tmp/reloader-test", err)
        return
    }

    c.Write([]byte{byte(msg)})
}

func main() {
    cmd := flag.String("c", "server", "List o possible Commands")

    flag.Parse()
    fmt.Println("Command: ", *cmd)

    switch (*cmd) {
        case "server":
            msg := openListener()
            switch (msg) {
                case QUIT, ERROR:
                    os.Exit(msg)
                case RELOAD:
                    err := os.Exec("./reloader-test.app", []string{"./reloader-test.app"}, os.Environ())
                    if err != nil {
                        fmt.Println("Error During Exec:", err)
                    }
                    os.Exit(msg)
            }
        case "quit":
            sendMsg(QUIT)
        case "reload":
            sendMsg(RELOAD)
        default:
    }
}
