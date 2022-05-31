package main

import (
//        "bufio"
        "fmt"
        "net"
        "os"
		"bytes"
		"io"
//        "strings"
//        "time"
)

type Serv_ses struct {
	Domain string
	EHLO bool
	State int
}

func main() {

	Serv := Serv_ses {
		Domain: "localhost",
		EHLO: true,
		State: 0,
	}
	l, err := net.Listen("tcp",":1025")
	if err != nil {
		fmt.Println("error creating listener: ",err)
        os.Exit(1)
	}
	defer l.Close()

	fmt.Println("****************************************************")
	fmt.Println("email server ready to accept connection on port 1025")
	fmt.Println("****************************************************")

	c, err := l.Accept()
    if err != nil {
		fmt.Println("error accepting con: ", err)
		os.Exit(1)
	}

	fmt.Println("email server accepted connection: ready for msgs!")
	fmt.Println("****************************************************")

	cl_reply := make([]byte, 2048)
// max line size is 2000

	_, err = c.Write([]byte("220 testserver.com ready \n"))
	if err != nil {
		fmt.Println("error in writing greet msg: ", err)
		os.Exit(1)
	}

	state :=0

	for i:=0; i<20; i++ {

		n, err := c.Read(cl_reply)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error in reading reply msg state ",i, ": ", err)
			os.Exit(1)
		}

// need to check for /lf/cr
		idx :=  bytes.Index(cl_reply[:n],[]byte("\r\n"))
		if idx < 0 {
			fmt.Println("no crlf ending! bad transmission!")
			_, err := c.Write([]byte("500 bad command no crlf\n\r"))
// normally need to terminate session
			if err != nil {
				fmt.Println("error writing nocrlf msg to client!")
			}
			os.Exit(1)
		}
		fmt.Println("CL->Srv: ", string(cl_reply[:idx]))

		switch Serv.State {
		case 0:
			bres := bytes.Equal(cl_reply[:4],[]byte("EHLO"))
			if !bres {
				fmt.Println("error did not receive EHLO in state: ",state)
				os.Exit(1)
			}

			cl_domain := cl_reply[4:idx]
			fmt.Println("client domain: ", string(cl_domain))
// need to sub domain names
			_, err := c.Write([]byte("250- server greets client\r\n"))
			if err != nil {
				fmt.Println("error state 0 writing msg greet! error: ", err)
				os.Exit(1)
			}
			fmt.Println("Srv->CL: 250- server greets client")

			_, err = c.Write([]byte("250- SIZE\r\n"))
			if err != nil {
				fmt.Println("error state 0 writing msg SIZE! error: ", err)
				os.Exit(1)
			}
			fmt.Println("Srv->CL: 250- SIZE")

			_, err = c.Write([]byte("250- DSN\r\n"))
			if err != nil {
				fmt.Println("error state 0 writing msg DSN! error: ", err)
				os.Exit(1)
			}
			fmt.Println("Srv->CL: 250- DSN")

			_, err = c.Write([]byte("250- 8BITMIME\r\n"))
			if err != nil {
				fmt.Println("error state 0 writing msg 8BITMIME! error: ", err)
				os.Exit(1)
			}
			fmt.Println("Srv->CL: 250- 8BITMIME")

			_, err = c.Write([]byte("250 HELP\r\n"))
			if err != nil {
				fmt.Println("error state 0 writing msg HELP! error: ", err)
				os.Exit(1)
			}
			fmt.Println("Srv->CL: 250 HELP")
			fmt.Println("completed EHLO Seq!")
			state +=1

		case 1:
			fmt.Println("****************************************************")

	    default:
		}
	}

	fmt.Println("****************************************************")
	fmt.Println("end server!")
}
