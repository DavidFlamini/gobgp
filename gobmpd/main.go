// Copyright (C) 2015 Nippon Telegraph and Telephone Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/osrg/gobgp/packet/bmp"
	"github.com/googollee/go-socket.io"
	"net"
	"net/http"
	"os"
	"strconv"
)

func connLoop(conn *net.TCPConn, so socketio.Socket) {
	addr := conn.RemoteAddr()
	scanner := bufio.NewScanner(bufio.NewReader(conn))
	scanner.Split(bmp.SplitBMP)

	for scanner.Scan() {
		msg, err := bmp.ParseBMPMessage(scanner.Bytes())
		if err != nil {
			log.Info(err)
			continue
		}
		j, _ := json.Marshal(msg)
		fmt.Print(string(j), "\n")

		so.Emit("data", string(j))

	}
	log.Info("conn was closed ", addr)
}


func main() {

	//BGP Monitoring Protocol Server
	service := ":" + strconv.Itoa(bmp.BMP_DEFAULT_PORT)
	addr, _ := net.ResolveTCPAddr("tcp", service)

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Info(err)
		os.Exit(1)
	}

	//Socket.io Server
	server, err := socketio.NewServer(nil)
	if err != nil {
			log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")

			//BGP Monitoring Protocol Server
			for {
				conn, err := l.AcceptTCP()
				if err != nil {
					log.Info(err)
					continue
				}
				log.Info("Accepted a new connection from ", conn.RemoteAddr())

				go connLoop(conn, so)
			}

			so.On("disconnection", func() {
					log.Println("on disconnect")
			})
	})
	server.On("error", func(so socketio.Socket, err error) {
			log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))

}
