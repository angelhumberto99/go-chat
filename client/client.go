package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
)

func listenServer(msgs *[]string, c net.Conn, name string) {
	var msg string
	for {
		err := gob.NewDecoder(c).Decode(&msg)
		if err != nil {
			fmt.Println(err)
			// terminamos el programa
			os.Exit(1)
		}
		// si el mensaje recibido lleva el nombre del cliente
		// entonces se reemplaza por la palabra "Tú"
		if strings.Contains(msg, name) {
			*msgs = append(*msgs, "Tú:" + msg[len(name)+1:])
		} else {
			*msgs = append(*msgs, msg)
		}
	}
}

func main() {
	var msgs []string
	var name string
	menu := "1) Mostrar mensajes/archivos\n" + 
			"2) Enviar mensaje\n" +
			"3) Enviar archivo\n" + 
			"4) Salir\n"
	input := bufio.NewScanner(os.Stdin)
	
	
	// pedimos el nombre al usuario
	fmt.Print("Ingrese su nombre: ")
	input.Scan()
	name = input.Text()

	// conectamos el cliente al servidor
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
	}

	// escuchara todas las respuestas del servidor
	go listenServer(&msgs, c, name)

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
				msg := bufio.NewScanner(os.Stdin)
				fmt.Print(">>> ")
				msg.Scan()
				err = gob.NewEncoder(c).Encode(name +": "+ msg.Text())
				if err != nil {
					fmt.Println(err)
				}
			case "3": // enviar archivo
				// todo
			case "4": // terminar cliente
				fmt.Println("Terminando cliente")

			default:
				fmt.Println("Opción incorrecta")
		}

		if input.Text() == "4" {
			// se envia un mensaje al servidor 
			// para eliminar la conexion al cliente
			err = gob.NewEncoder(c).Encode("/quit")
			if err != nil {
				fmt.Println(err)
			}
			break
		}

	}
	// terminamos la llamada
	c.Close()
}