package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"unicode/utf16"
	"unicode/utf8"
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

// All fields are stored as bytes. Some fields can be
// turned into string with regular bytes-to-string methods,
// e.g. "string(s.Name)"
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

func decodeName(b []byte) string {
	var newName string
	for len(b) > 0 {
		var r rune
		var size int
		r, _ = utf8.DecodeRune(b)
		if r != 65533 {
			r, size = utf8.DecodeRune(b)
			newName = newName + string(r)
			b = b[size:]
		} else {
			if len(b) < 6 {
				// It can happen when there is an invalid surrogate
				// at the end of room name.
				// It is shown as question mark in room list.
				invalidSurrogate, _ := utf8.DecodeRune(b)
				newName = newName + string(invalidSurrogate)
				b = b[3:]
				continue
			}
			const (
				maskx = 0b00111111
				mask2 = 0b00011111
				mask3 = 0b00001111
				mask4 = 0b00000111
			)

			b1 := rune(b[0]&mask3)<<12 | rune(b[1]&maskx)<<6 | rune(b[2]&maskx)
			b2 := rune(b[3]&mask3)<<12 | rune(b[4]&maskx)<<6 | rune(b[5]&maskx)
			r = utf16.DecodeRune(b1, b2)
			newName = newName + string(r)
			b = b[6:]
		}
	}
	return newName
}

// Pretty prints server data.
func (s *Server) Print() {
	//fmt.Printf("Link: %s\nName: %v\n%s\n%x\n%b\nFlag: %s\nPrivate: %v\nPlayers: %v/%v\n", s.Link, s.Name, s.Name, s.Name, s.Name, s.Flag, s.Private, s.PlayersNow, s.PlayersMax)
	n := decodeName(s.Name)
	//fmt.Printf("%v\n", utf16.IsSurrogate(r))
	//fmt.Printf("%v\n", utf16.Encode(r))
	//fmt.Printf("%s\n%b\n%v\n%v\n----\n", s.Name, s.Name, r, size)
	fmt.Printf("%s\n----\n", n)
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
