package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func getData() ([]byte, error) {
	resp, err := http.Get("https://www.haxball.com/rs/api/list")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, err
}

type server struct {
	link       []byte
	unknown1   [2]byte
	name       []byte
	flag       []byte
	unknown2   []byte
	private    byte
	playersMax byte
	playersNow byte
}

func (s *server) print() {
	fmt.Printf("Link: %s\nName: %s\nFlag: %s\nPrivate: %v\nPlayers: %v/%v\nUnknown1: %v\nUnknown2: %v\n", s.link, s.name, s.flag, s.private, s.playersNow, s.playersMax, s.unknown1, s.unknown2)
}

func parseServers(body []byte) []server {
	var serverList []server
	startIndex := 0
	body = body[1:]
	for startIndex < len(body) {
		singleServer := server{}
		singleServer.link = body[startIndex+2 : startIndex+13]
		singleServer.unknown1[0] = body[startIndex+14]
		singleServer.unknown1[1] = body[startIndex+17]
		var nameLength int
		for i, b := range body[startIndex+18:] {
			if b == 2 {
				nameLength = i
				break
			}
			singleServer.name = append(singleServer.name, b)
		}
		singleServer.flag = body[startIndex+nameLength+19 : startIndex+nameLength+21]
		singleServer.unknown2 = body[startIndex+nameLength+21 : startIndex+nameLength+29]
		singleServer.private = body[startIndex+nameLength+29]
		singleServer.playersMax = body[startIndex+nameLength+30]
		singleServer.playersNow = body[startIndex+nameLength+31]
		serverList = append(serverList, singleServer)
		startIndex = startIndex + nameLength + 32
	}
	return serverList
}

func main() {
	body, err := getData()
	if err != nil {
		log.Fatalln(err)
	}
	serverList := parseServers(body)
	for _, server := range serverList {
		server.print()
                fmt.Println()
	}
}
