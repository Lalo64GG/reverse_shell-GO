package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
)

// Llave de cifrado AES (16, 24 o 32 bytes)
var encryptionKey = []byte("supersecretkey123") // Cambia esto por una llave más segura

func generateRandomPort(min, max int64)(int64, error){
	rangeVal := max - min +1

	// Generar un Puerto aleatorio en el rango [0, rangeVal]
	n, err := rand.Int(rand.Reader, big.NewInt(rangeVal))
	if err != nil {
		return 0, err
	}

	return n.Int64() + min, nil
}

func main() {
	min := int64(1024)
	max := int64(49151)

	port, err := generateRandomPort(min, max)
	if err != nil {
		fmt.Println("Error generando puerto aleatorio: ", err)
		return
	}

	for {
		conn := connectToServer("192.168.0.14:" + strconv.Itoa(int(port)))
		if conn != nil {
			defer conn.Close()

			// Manejar la comunicación con el servidor
			go handleInput(conn)
			go handleOutput(conn)
			go listenForCommands(conn)

			// Información del sistema operativo
			switch runtime.GOOS {
			case "linux":
				fmt.Println("Ejecutando en Linux")
			case "darwin":
				fmt.Println("Ejecutando en macOS")
			case "windows":
				fmt.Println("Ejecutando en Windows")
			default:
				fmt.Println("Sistema operativo no reconocido")
			}
		}
		time.Sleep(5 * time.Second) // Espera 5 segundos antes de intentar reconectar
	}
}

// Función para conectar al servidor con reintento en caso de fallo
func connectToServer(address string) net.Conn {
	for {
		conn, err := net.Dial("tcp", address)
		if err == nil {
			return conn
		}
		fmt.Println("Error al conectar, reintentando...")
		time.Sleep(5 * time.Second) // Reintenta cada 5 segundos
	}
}

// Manejar la entrada del usuario (stdin) y enviar al servidor con cifrado
func handleInput(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Println("Error de lectura de entrada:", err)
			return
		}

		encryptedData, err := encrypt(buf[:n])
		if err != nil {
			fmt.Println("Error en cifrado:", err)
			return
		}

		_, err = conn.Write(encryptedData)
		if err != nil {
			fmt.Println("Error al enviar datos:", err)
			return
		}
	}
}

// Manejar la salida del servidor (stdout) y desencriptar
func handleOutput(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error de lectura de salida:", err)
			return
		}

		decryptedData, err := decrypt(buf[:n])
		if err != nil {
			fmt.Println("Error en desencriptado:", err)
			return
		}

		_, err = os.Stdout.Write(decryptedData)
		if err != nil {
			fmt.Println("Error al escribir en salida estándar:", err)
			return
		}
	}
}

// Escuchar y ejecutar comandos específicos del servidor
func listenForCommands(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error de lectura de comandos:", err)
			return
		}
		command := string(buf[:n])

		// Comandos de ejemplo
		switch command {
		case "exit":
			fmt.Println("Cerrando conexión")
			conn.Close()
			os.Exit(0)
		default:
			fmt.Println("Comando desconocido:", command)
		}
	}
}

// Cifrado AES de datos
func encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// Desencriptar datos AES
func decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("datos incorrectos")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
