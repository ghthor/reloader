package main

import (
    "fmt"
    "flag"
    "net"
    "os"
    "exec"
    "path"
    "log"
)

var exeName string

func init() {
    exeName = path.Base(os.Args[0])
}

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

func openListener() Cmd {
    log.Println("Opening UnixSocket /tmp/" + exeName)
    l, err := net.Listen("unix", "/tmp/" + exeName)
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

    m := Cmd(msg[:n][0])

    log.Println("Received: ", m)

    return m
}

type loopFunc func() loopFunc

func runServer() loopFunc {
    msg := openListener()
    switch (msg) {
    case QUIT, ERROR:
        os.Exit(int(msg))
    case REBUILD:
        if !rebuild() {
            return runServer
        }
        fallthrough
    case RELOAD:
        reloadServer()
        os.Exit(int(msg))
        return nil
    }
    return nil
}

func rebuild() bool {
    cmd := exec.Command("make")

    log.Println("Rebuilding")
    err := cmd.Run()
    if err != nil {
        fmt.Println("Error Running Make:", err)
        return false
    }
    log.Println("Rebuild Success")
    return true
}

func reloadServer() {
    log.Println("Reloading Executable")
    err := os.Exec(os.Args[0], os.Args, os.Environ())
    if err != nil {
        fmt.Println("Error During Exec:", err)
    }
}

func sendCmd(cmd Cmd) {
    c, err := net.Dial("unix", "/tmp/" + exeName)
    defer c.Close()
    if err != nil {
        fmt.Println("Error Dialing:", "/tmp/" + exeName, err)
        return
    }

    c.Write([]byte{byte(cmd)})
}

func main() {
    cmd := flag.String("c", "rebuild", "List o possible Commands")
    flag.Parse()

    switch (*cmd) {
        case "quit":
            sendCmd(QUIT)
        case "rebuild":
            sendCmd(REBUILD)
        case "reload":
            sendCmd(RELOAD)
        case "server":
            f := runServer()
            for f != nil {
                f = f()
            }
        default:
    }
}
