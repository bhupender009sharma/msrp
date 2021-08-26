package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var mutex sync.Mutex

func sendMessage(c net.Conn) {

	var text string = ""

	fmt.Println("Want send Message  by: \n 1)Automatic\n 2)Manual")
	var choice int
	fmt.Scan(&choice)
	if choice == 2 {
		fmt.Println("Enter your Message here")
		fmt.Println("--------------------")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {

			if scanner.Text() == "exit;" {
				//your actions on exit...
				//...
				fmt.Println("Exiting........\n")
				break
			}

			text += scanner.Text() + "\n"
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

		text = text + "\n\r"
		c.Write([]byte(text))

		message, _ := bufio.NewReader(c).ReadString('\r')

		fmt.Print("\n Reply:\n ", message)

	} else {
		for {
			text = `MSRP cf147547 SEND
		To-Path: msrp://[2405:200:610:1587:21::10]:62222/z3ogutjj9q5ob;tcp
		From-Path: msrp://[2409:4040:e98:9906:bda6:a64a:ecf:f71f]:9/3eb80645b9;tcp
		Message-ID: dfc5eb13
		Success-Report: no
		Failure-Report: yes
		Byte-Range: 1-286/286
		Content-Type: message/cpim

		From: <sip:anonymous@anonymous.invalid>
		To: <sip:anonymous@anonymous.invalid>
		DateTime: 2021-05-04T19:33:18+05:30
		NS: imdn <urn:ietf:params:imdn>
		imdn.Disposition-Notification: positive-delivery, display
		imdn.Message-ID: 1088

		Content-Type: text/plain
		Content-Length: 25

		Hello
		-------cf147547$`

			text = text + "\n\r"
			c.Write([]byte(text))

			message, _ := bufio.NewReader(c).ReadString('\r')

			fmt.Print("\n Reply:\n ", message)

			time.Sleep(1 * time.Second)
		}
	}

}

func main() {

	CONNECT := "127.0.0.1:5050"

	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err) // if there is an error in connecting, then print it
		return
	}

	for {
		go func() {

			buf := make([]byte, 1024)

			n, e := c.Read(buf)
			if e != nil {
				fmt.Println(err)
				return
			}

			if n > 0 {

				var netData string = string(buf[:n])
				fmt.Println("Incoming message: \n ", netData)

				res1 := strings.Split(netData, "\n") //slipts the recieved data where new line is found
				tpath := strings.Split(res1[1], "-Path:")
				fpath := strings.Split(res1[2], "-Path:")
				msrp := strings.Split(res1[0], "SEND")
				var l int = len(res1)
				text := msrp[0] + "200 OK\nTo-Path:" + fpath[1] + "\nFrom-Path:" + tpath[1] + "\n" + res1[l-3] + "\n\r"

				c.Write([]byte(text)) // Write takes []byte type of data input only
				//fmt.Fprintf(c, text)

			}

		}()

		sendMessage(c)

		time.Sleep(500 * time.Millisecond)
	}
}
