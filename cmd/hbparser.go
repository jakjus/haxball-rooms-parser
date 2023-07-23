package cmd

import (
	"fmt"
	"io"
	"log"
        "bytes"
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
	unknown1   []byte
	encodedName[]byte
	Name       string
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
                // If not invalid surrogate in utf-8 interpretation
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
				break
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
	//fmt.Printf("Link: %s\nName: %v\n%x\nFlag: %s\nPrivate: %v\nPlayers: %v/%v\n", s.Link, s.Name, s.Name, s.Flag, s.Private, s.PlayersNow, s.PlayersMax)
        //fmt.Printf("%b\n", s.Name)
	//fmt.Printf("%v\n", utf16.IsSurrogate(r))
	//fmt.Printf("%v\n", utf16.Encode(r))
	//fmt.Printf("%s\n%b\n%v\n%v\n----\n", s.Name, s.Name, r, size)
	fmt.Printf("%s\n----\n", s.Name)
}

// Parses raw bytes into Server struct, based on bytes indices and dividers.
func Parse(body []byte) []Server {
	var ServerList []Server
	//startIndex := 0
	body = body[1:]
        //fmt.Printf("\n%s\n%v", body[:128], body[:128])
	for len(body) > 0 {
		singleServer := Server{}
                _, body, _ = bytes.Cut(body, []byte{0b0})
                singleServer.Link, body, _ = bytes.Cut(body, []byte{0b0})
                if len(singleServer.Link) < 5 {
                  continue
                }
                singleServer.Link = singleServer.Link[1:]
                singleServer.unknown1, body, _ = bytes.Cut(body, []byte{0b0})
                body = body[1:]

		var nameLength int
		for i, b := range body {
			if b == 0b10 || b == 0b11 {
				nameLength = i
				break
			}
			singleServer.encodedName = append(singleServer.encodedName, b)
		}
                singleServer.Name = decodeName(singleServer.encodedName)
                //startIndex += nameLength
                body = body[nameLength:]
		singleServer.Flag = body[1:3]
		singleServer.unknown2 = body[3:11]
		singleServer.Private = body[11]
		singleServer.PlayersMax = body[12]
		singleServer.PlayersNow = body[13]
                body = body[14:]
               //// fmt.Printf("%v\n", body[startIndex:startIndex+18])
               //// fmt.Printf("%s\n", body[startIndex:startIndex+18])
		ServerList = append(ServerList, singleServer)
		//startIndex = startIndex + 32
	}
	return ServerList
}
