package hbparser

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// Downloads raw data from Haxball endpoint.
func GetData() ([]byte, error) {
	resp, err := http.Get("https://www.haxball.com/rs/api/list")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, err
}

// All fields are stored as bytes. Some fields (e.g. "name") may be stored
// as strings or rune[] IN THE FUTURE, when the appropriate encoding
// will be found.
type Server struct {
	Link       []byte
	unknown1   [2]byte
	Name       []byte
	Flag       []byte
	unknown2   []byte
	Private    byte
	PlayersMax byte
	PlayersNow byte
}

// Pretty prints server data.
func (s *Server) Print() {
	fmt.Printf("Link: %s\nName: %v\nFlag: %s\nPrivate: %v\nPlayers: %v/%v\n", s.Link, s.Name, s.Flag, s.Private, s.PlayersNow, s.PlayersMax)
}

// Parses raw bytes into Server struct, based on bytes indices and dividers.
func Parse(body []byte) []Server {
	var ServerList []Server
	startIndex := 0
	body = body[1:]
	for startIndex < len(body) {
		singleServer := Server{}
		singleServer.Link = body[startIndex+2 : startIndex+13]
		singleServer.unknown1[0] = body[startIndex+14]
		singleServer.unknown1[1] = body[startIndex+17]
		var nameLength int
		for i, b := range body[startIndex+18:] {
			if b == 2 {
				nameLength = i
				break
			}
			singleServer.Name = append(singleServer.Name, b)
		}
		singleServer.Flag = body[startIndex+nameLength+19 : startIndex+nameLength+21]
		singleServer.unknown2 = body[startIndex+nameLength+21 : startIndex+nameLength+29]
		singleServer.Private = body[startIndex+nameLength+29]
		singleServer.PlayersMax = body[startIndex+nameLength+30]
		singleServer.PlayersNow = body[startIndex+nameLength+31]
		ServerList = append(ServerList, singleServer)
		startIndex = startIndex + nameLength + 32
	}
	return ServerList
}
