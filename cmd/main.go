package main

import (
	"http/pkg/server"
	"log"
	"net"
	"os"
	"strconv"
)

func main() {
	host := "0.0.0.0"
	port := "9999"

	err := execute(host, port)
	if err != nil {
		os.Exit(1)
	}
}

//func execute(host string, port string) (err error) {
//	listener, err := net.Listen("tcp", net.JoinHostPort(host, port))
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	defer func() {
//		if err2 := listener.Close(); err2 != nil {
//			if err == nil {
//				err = err2
//				return
//			}
//			log.Println(err2)
//		}
//	}()
//	// to do server code
//
//	for {
//		conn, err := listener.Accept()
//		if err != nil {
//			log.Println(err)
//			continue
//		}
//
//		err = handle(conn)
//		if err != nil {
//			log.Println(err)
//			continue
//		}
//	}
//
//	return
//}

//func handle(conn net.Conn) (err error) {
//	defer func() {
//		if err2 := conn.Close(); err2 != nil {
//			if err == nil {
//				err = err2
//				return
//			}
//			log.Println(err)
//		}
//	}()
//
//	// TO DO handle connection
//
//	buf := make([]byte, 4096)
//	for {
//		read, err := conn.Read(buf)
//		if err == io.EOF {
//			log.Printf("%s", buf[:read])
//		return nil
//		}
//		if err != nil {
//			return err
//		}
//		log.Printf("%s", buf[:read])
//	}
//
//	return
//}

func execute(host string, port string) (err error) {
	const rn = "\r\n"
	srv := server.NewServer(net.JoinHostPort(host, port))
	srv.Register("/", func(conn net.Conn) {
		body := "Welcome to our web-site"
		_, err = conn.Write([]byte(
			"HTTP/1.1 200 OK" + rn +
				"Content-Length: " + strconv.Itoa(len(body)) + rn +
				"Content-Type: text/html" + rn +
				"Connection: close" + rn +
				rn +
				body,
		))
		if err != nil {
			log.Println(err)
		}

	})

	srv.Register("/about", func(conn net.Conn) {
		body := "About Golang Academy"
		_, err = conn.Write([]byte(
			"HTTP/1.1 200 OK" + rn +
				"Content-Length: " + strconv.Itoa(len(body)) + rn +
				"Content-Type: text/html" + rn +
				"Connection: close" + rn +
				rn +
				body,

		))
		if err != nil {
			log.Println(err)
		}
	})

	return srv.Start()
}
