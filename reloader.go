package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"syscall"
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
	log.Println("Action listen unix /tmp/" + exeName)
	l, err := net.Listen("unix", "/tmp/"+exeName)
	defer l.Close()
	if err != nil {
		log.Println("Error", err)
		return ERROR
	}

	c, err := l.Accept()
	defer c.Close()
	if err != nil {
		log.Println("Error accept unix", err)
		return ERROR
	}

	msg := make([]byte, 0, 1024)
	n, err := c.Read(msg[:1])
	if err != nil && err != io.EOF {
		log.Println("Error ", err)
		return ERROR
	}

	m := Cmd(msg[:n][0])

	log.Println("Received", m)

	return m
}

type loopFunc func() loopFunc

func runServer() loopFunc {
	msg := openListener()
	switch msg {
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

	log.Println("Action rebuild")
	err := cmd.Run()
	if err != nil {
		log.Println("Error execute make", err)
		return false
	}
	log.Println("Status rebuilt")
	return true
}

func reloadServer() {
	log.Println("Action reload")
	err := syscall.Exec(os.Args[0], os.Args, os.Environ())
	if err != nil {
		log.Println("Error", err)
	}
}

func sendCmd(cmd Cmd) {
	log.Printf("Action send cmd %s", cmd.String())
	c, err := net.Dial("unix", "/tmp/"+exeName)
	defer c.Close()
	if err != nil {
		log.Println("Error", err)
		return
	}

	c.Write([]byte{byte(cmd)})
}

func main() {
	cmd := flag.String("c", "rebuild", "List o possible Commands")
	flag.Parse()

	switch *cmd {
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
