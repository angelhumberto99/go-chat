package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

func main() {
	var msgs []string
	var name string
	menu := "1) Ver mensajes\n" + 
			"2) Enviar mensaje\n" + 
			"3) Salir\n"
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
				fmt.Print("> ")
				msg.Scan()
				err = gob.NewEncoder(c).Encode(name +": "+ msg.Text())
				if err != nil {
					fmt.Println(err)
				}
			case "3": // terminar cliente
				fmt.Println("Terminando cliente")

			default:
				fmt.Println("Opción incorrecta")
		}

		if input.Text() == "3" {
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