package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
)

func GetHawkeyeTarget() *url.URL {
	rawURL := os.Getenv("TUPLESTREAM_HAWKEYE_TARGET")
	if rawURL == "" {
		rawURL = "https://incoming.tuplestream.com"
	}
	u, e := url.Parse(rawURL)
	handleErr(e)
	return u
}

func InitiateConnection(filename string) (net.Conn, *bufio.Writer) {
	hawkeyeTarget := GetHawkeyeTarget()
	conn, err := net.Dial("tcp", hawkeyeTarget.Host)

	if err != nil {
		log.Print("error connecting to " + hawkeyeTarget.Host)
		return nil, nil
	}

	req, err := http.NewRequest("GET", GetHawkeyeTarget().String(), nil)
	handleErr(err)

	req.Header.Add("Connection", "Upgrade")
	req.Header.Add("Upgrade", "hawkeye/1.0.0alpha1")
	req.Header.Add("User-Agent", "hawkeye/client-go1.0.0alpha1")

	handleErr(err)

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	err = req.Write(writer)
	handleErr(err)
	writer.Flush()

	resp, err := http.ReadResponse(reader, req)
	if resp.StatusCode != 101 {
		log.Fatal("Couldn't upgrade HTTP connection, closing. Got status: " + resp.Status)
	}
	handleErr(err)

	fmt.Println(resp.Status)

	controlMessage := make(map[string]string)
	encoder := json.NewEncoder(writer)

	controlMessage["__hawkeye_filename"] = filename

	err = encoder.Encode(controlMessage)
	handleErr(err)
	writer.Flush()

	ok, err := reader.ReadString('\n')
	handleErr(err)
	if ok == "OK\n" {
		log.Print("handshake successful")
	}
	return conn, writer
}

func handleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
