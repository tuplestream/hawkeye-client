package hawkeye

import (
	"bufio"
	"crypto/tls"
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

	if u.Port() == "" {
		var presumedPort string
		if u.Scheme == "http" {
			presumedPort = "80"
		} else if u.Scheme == "https" {
			presumedPort = "443"
		} else {
			log.Panic("Unsupported scheme for target " + u.Scheme)
		}

		u, e = url.Parse(fmt.Sprintf("%s://%s:%s/%s", u.Scheme, u.Hostname(), presumedPort, u.Path))
	}
	return u
}

func InitiateConnection(filename string, auth string) (net.Conn, *bufio.Writer) {
	hawkeyeTarget := GetHawkeyeTarget()
	host := hawkeyeTarget.Hostname() + ":" + hawkeyeTarget.Port()
	var conn net.Conn
	var err error

	if hawkeyeTarget.Scheme == "https" {
		conn, err = tls.Dial("tcp", host, nil)
	} else {
		conn, err = net.Dial("tcp", host)
	}

	if err != nil {
		log.Panic("error connecting to " + host)
		return nil, nil
	}

	req, err := http.NewRequest("GET", GetHawkeyeTarget().String(), nil)
	handleErr(err)

	req.Header.Add("Connection", "Upgrade")
	req.Header.Add("Upgrade", "hawkeye/1.0.0alpha1")
	req.Header.Add("User-Agent", "hawkeye/client-go1.0.0alpha1")

	if auth != "" {
		req.Header.Add("Authorization", "Bearer "+auth)
	}

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
