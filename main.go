package main

/*
Code snipet for convert UTF-8 from Shift-JIS.

Reference source:
http://qiita.com/nobuhito/items/ff782f64e32f7ed95e43
*/

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
)

func main() {
	msg := ""
	flag.Parse()
	if flag.NArg() >= 1 {
		msg, _ = to_utf8(flag.Arg(0))
	}

	fmt.Println("msg:", msg)
}

// Shift-JIS -> UTF-8
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
