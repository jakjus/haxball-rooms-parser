package cmd

import (
	"fmt"
	"io"
	"log"
        "bytes"
        //"strings"
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

// Pretty prints server data.
func (s *Server) TableRow() {
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
	body = body[1:]
        var link, rest []byte
        for len(body) > 0 {
          singleServer := Server{}
          link, body, _ = bytes.Cut(body, []byte{9, 0})
          fmt.Println("link", link)
          singleServer.Link = link[2:len(link)-2]
          indBody := len(body)
          for i, char := range body {
            if (i<2) { continue }
            if (body[i-2] == 0 && body[i-1] == 11 && char > body[i-1]) {
              indBody = i-2
              break
            }
          }
          rest = body[:indBody]
          fmt.Println(rest)
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
          //encodedName, flagPart, _ := bytes.Cut(rest[:indc], []byte{2, 0, 3})
          singleServer.Name = decodeName(singleServer.encodedName)
          singleServer.Flag = flagPart[1:3]
          fmt.Printf("flag: %v\n", flagPart)
          singleServer.unknown = flagPart[3:len(flagPart)-3]
          //fmt.Printf("priv: %v\n", rest[indc:])
          singleServer.Private = flagPart[len(flagPart)-3]
          singleServer.PlayersMax = flagPart[len(flagPart)-2]
          singleServer.PlayersNow = flagPart[len(flagPart)-1]
          if singleServer.PlayersMax > 30 {
            fmt.Println(singleServer.Name)
            fmt.Printf("%v\n", singleServer.encodedName)
            fmt.Printf("flagpart: %v\n", flagPart)
            fmt.Println("rest")
            fmt.Println(rest)
            fmt.Println("body")
            fmt.Println(body[:100])
          }
          ServerList = append(ServerList, singleServer)
        }
        /*
        for _, s := range servers {
          linkEnd := bytes.IndexByte(s[2:], 0b0)
          //fmt.Println(s[2:13])
          link := s[2:linkEnd]
          fmt.Println(linkEnd, link)
        }
        if (body[i] == 9 && body[i+1] == 0) {
          fmt.Printf("\n")
        }
        if (body[i] == 0) {
          fmt.Printf("\n")
        }
        fmt.Printf("\n%s\n%v", body[:128], body[:128])
	for len(body) > 0 {
		singleServer := Server{}
                _, body, _ = bytes.Cut(body, []byte{0b0})
                singleServer.Link, body, _ = bytes.Cut(body, []byte{0b0})
                if len(singleServer.Link) < 5 {
                  fmt.Printf("error link is too short")
                  fmt.Printf("%v\n", singleServer.Link)
                  fmt.Printf("%v\n", body[:10])
                  continue
                }
                singleServer.Link = singleServer.Link[1:]
                singleServer.unknown1, body, _ = bytes.Cut(body, []byte{0b0})
                body = body[1:]

		var nameLength int
		for i, b := range body {
			if b == 0b10 || b == 0b11 || b == 0b0 {
				nameLength = i
				break
			}
			singleServer.encodedName = append(singleServer.encodedName, b)
		}
                singleServer.Name = decodeName(singleServer.encodedName)
                body = body[nameLength:]
		singleServer.Flag = body[1:3]
                //unknown2End := bytes.IndexByte(body[3:], 0b1000001)+3
                unknown2End := 11
                if body[11] != 0b0 && body[11] != 0b1 {
                  unknown2End = 12
                }
                //if bytes.Index(body[3:18], byte(65)) == -1 {
                //}
                if body[unknown2End+2] > body[unknown2End+1] {
                  fmt.Println(singleServer.Name)
                  fmt.Printf("%v\n%b\n\n", body[3:18], body[3:18])
                }
                //fmt.Println(unknown2End)
                singleServer.Geo = body[3:9]
		singleServer.Private = body[unknown2End]
		singleServer.PlayersMax = body[unknown2End+1]
		singleServer.PlayersNow = body[unknown2End+2]
                body = body[unknown2End+2:]
               //// fmt.Printf("%v\n", body[startIndex:startIndex+18])
               //// fmt.Printf("%s\n", body[startIndex:startIndex+18])
		ServerList = append(ServerList, singleServer)
		//startIndex = startIndex + 32
	}
        */
	return ServerList
}
