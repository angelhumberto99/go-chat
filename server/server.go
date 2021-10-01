package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
)

func listenClients(msgs *[]string, clients *[]net.Conn, c net.Conn) {
	var msg string
	for {
		err := gob.NewDecoder(c).Decode(&msg)
		if err != nil {
			fmt.Println(err)
		}
		if msg == "/quit" {
			for i, v := range *clients {
				if v == c {
					*clients = append((*clients)[:i], (*clients)[i+1:]...)
					break
				}
			}
			return
		}
		*msgs = append(*msgs, msg)
	}
}

func checkConnection(s net.Listener, clients *[]net.Conn, msgs *[]string) {
	var exists bool
	for {
		// peticiones del cliente
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
		}
		
		exists = false
		for _,v := range(*clients) {
			if v == c {
				exists = true
				break
			}
		}
		
		if !exists {
			*clients = append(*clients, c)
			go listenClients(msgs, clients, c)
		}
	}
}

func main() {
	var clients []net.Conn
	var msgs []string
	menu := "1) Ver mensajes\n" + 
			"2) Enviar mensaje\n" + 
			"3) Salir\n"
	input := bufio.NewScanner(os.Stdin)


	// se crea el servidor
	s, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer s.Close()
	go checkConnection(s, &clients, &msgs)

	for {
		fmt.Print(menu)
		input.Scan()
		switch input.Text() {
			case "1": // mostrar mensajes
				fmt.Println("Mensajes")
				for _,v := range msgs {
					fmt.Println(v)
				}
			case "2": // enviar mensaje
				for _,v := range clients {
					io.WriteString(v, "hola")
				}
			case "3": // terminar cliente
				fmt.Println("Terminando Servidor")
			default:
				fmt.Println("Opción incorrecta")
		}
		if input.Text() == "3" {
			break
		}
	}
}