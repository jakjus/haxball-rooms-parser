package cmd

import (
	"bytes"
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
	Link        []byte
	unknown     []byte
	encodedName []byte
	Name        string
	Flag        []byte
	Private     byte
	PlayersMax  byte
	PlayersNow  byte
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

// Parses raw bytes into Server struct, based on bytes indices and dividers.
func Parse(body []byte) []Server {
	var ServerList []Server
	body = body[1:]
	var link, rest []byte
	for len(body) > 0 {
		singleServer := Server{}
		link, body, _ = bytes.Cut(body, []byte{9, 0})
		singleServer.Link = link[2 : len(link)-2]
		indBody := len(body)
		for i, char := range body {
			if i < 2 {
				continue
			}
			if body[i-2] == 0 && body[i-1] == 11 && char > body[i-1] {
				indBody = i - 2
				break
			}
		}
		rest = body[:indBody]
		body = body[indBody:]
		rest = rest[1:]
		indc2 := 0
		for i, char := range rest {
			if (char == 0) || (char == 2) || (char == 3) {
				indc2 = i
				break
			}
		}
		if indc2 == 0 {
			continue
		}
		singleServer.encodedName = rest[:indc2]
		flagPart := rest[indc2:]
		singleServer.Name = decodeName(singleServer.encodedName)
                if (len(flagPart) < 5) {
                  continue
                }
		singleServer.Flag = flagPart[1:3]
		singleServer.unknown = flagPart[3 : len(flagPart)-3]
		singleServer.Private = flagPart[len(flagPart)-3]
		singleServer.PlayersMax = flagPart[len(flagPart)-2]
		singleServer.PlayersNow = flagPart[len(flagPart)-1]
		ServerList = append(ServerList, singleServer)
	}
	return ServerList
}
