package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

func to_utf8(str string) (string, error) {
	body, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return "", err
	}

	var f []byte
	encodings := []string{"sjis", "utf-8"}
	for _, enc := range encodings {
		if enc != "" {
			ee, _ := charset.Lookup(enc)
			if ee == nil {
				continue
			}
			var buf bytes.Buffer
			ic := transform.NewWriter(&buf, ee.NewDecoder())
			_, err := ic.Write(body)
			if err != nil {
				continue
			}
			err = ic.Close()
			if err != nil {
				continue
			}
			f = buf.Bytes()
			break
		}
	}
	return string(f), nil
}

func sendBroadcast(msg string, port int) {
	BROADCAST_IPv4 := net.IPv4(255, 255, 255, 255)
	socket, _ := net.DialUDP("udp4", nil, &net.UDPAddr{IP: BROADCAST_IPv4, Port: port})
	socket.Write([]byte(msg))
}

func main() {
	msg := ""
	port := 12342

	flag.Parse()
	if flag.NArg() >= 1 {
		str, _ := to_utf8(flag.Arg(0))
		msg = strings.Replace(str, "__tab__", "\t", 1)
		if flag.NArg() >= 2 {
			num1, err := strconv.Atoi(flag.Arg(1))
			if err != nil {
				msg = msg + "\t" + flag.Arg(1)
				if flag.NArg() >= 3 {
					num2, _ := strconv.Atoi(flag.Arg(2))
					port = num2
				}
			} else {
				port = num1
			}
		}
	}

	fmt.Println("msg:", msg)
	fmt.Println("port:", port)
	sendBroadcast(msg, port)
}
