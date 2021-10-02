package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
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

		// separación de cadenas
		if strings.Contains(msg, "/file") {
			info := strings.Split(msg, "?")[0]
			data := strings.Split(info, ">")

			fileName := (data[0])[len("/file<"):]
			user := (data[1])[1:len(data[1])-1]
			fileInfo := msg[len(info)+1:]

			if user != name {
				dest, err := os.Create("./files/"+fileName)
				if err != nil {
					fmt.Println(err)
				} 
				dest.Write([]byte(fileInfo))
				dest.Close()
				*msgs = append(*msgs, user + " envío: " + fileName)
			} else {
				*msgs = append(*msgs, "Enviaste: " + fileName)
			}

		} else if strings.Contains(msg, name) {
			// si el mensaje recibido lleva el nombre del cliente
			// entonces se reemplaza por la palabra "Tú"
			*msgs = append(*msgs, "Tú:" + msg[len(name)+1:])
		} else {
			*msgs = append(*msgs, msg)
		}
	}
}

func listDirectory() string{
	dir, err := ioutil.ReadDir("./files")
	var files []string
	input := bufio.NewScanner(os.Stdin)

    if err != nil {
        fmt.Println(err)
    }
	for _, v := range dir {
		files = append(files, v.Name())
	}

	fmt.Println("Archivos")
	if len(files) == 0 {
		return "-1"
	}
	// menú
	for i, v := range files {
		fmt.Printf("%d) %s\n", i, v)
	}
	fmt.Printf("%d) Salir\n", len(files))
	input.Scan()
	resp,_ := strconv.Atoi(input.Text())
	if len(files) == resp {
		return "-1"
	}
	return files[resp]
}

func sendFile(fileName, userName string, c net.Conn) {
	origin, err := os.ReadFile("./files/"+fileName)
	if err != nil {
		fmt.Println(err)
	}
	// información para enviar
	fileStr := "/file<"+fileName+">("+userName+")?"+string(origin) 
	err = gob.NewEncoder(c).Encode(fileStr)
	if err != nil {
		fmt.Println(err)
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
				fileName := listDirectory()
				if fileName != "-1" {
					sendFile(fileName, name, c)
				}
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