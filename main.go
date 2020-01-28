package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

func readViaFile(path string) (string, string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	array := strings.Split(string(file), "\n")

	return array[0], array[1]
}

func main() {
	// Check for inputs
	ViaFile := flag.String("file", "", "Path to file containing username and password")
	Server := flag.String("server", "", "Hostname/IP of the radius server.")
	Port := flag.String("port", "1812", "Port of the radius server.")
	Secret := flag.String("secret", "", "Radius server shared secret")
	UserName := flag.String("username", "", "Username to authenticate with.")
	Password := flag.String("password", "", "Password of the user we are authenticating with")
	flag.Parse()

	// Script can accept username and passwords via path to file with username and password in ( see https://openvpn.net/community-resources/reference-manual-for-openvpn-2-2/ via-file) or via -username -password flags
	if *ViaFile != "" {
		if _, err := os.Stat(*ViaFile); os.IsNotExist(err) {
			log.Fatal(err)
		}
		*UserName, *Password = readViaFile(*ViaFile)
	} else if *UserName == "" {
		log.Fatal("-username not provided")
	} else if *Password == "" {
		log.Fatal("-password not provided")
	} else if *Secret == "" {
		log.Fatal("-secret not provided")
	}

	// Create new request to our radius server
	packet := radius.New(radius.CodeAccessRequest, []byte(*Secret))
	rfc2865.UserName_SetString(packet, *UserName)
	rfc2865.UserPassword_SetString(packet, *Password)
	response, err := radius.Exchange(context.Background(), packet, *Server+":"+*Port)
	if err != nil {
		log.Fatal(err)
	}

	// Do stuff with the response, either return success or error
	log.Println("User:", *UserName, " Code:", response.Code)

	if response.Code != 2 {
		log.Fatal("error, wanted Access-Accept, got: " + response.Code.String())
	}
}
