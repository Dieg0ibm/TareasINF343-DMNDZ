package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"Tarea2/orden"

	"google.golang.org/grpc"
)

// ///////////////////////// Estructura ///////////////////////////
type CompraInput struct {
	Vehiculos []orden.Vehicle `json:"vehicles"`
	Cliente   orden.Customer  `json:"customer"`
}

// ///////////////////////// Conexión gRPC ///////////////////////////
func conectarServer() (*grpc.ClientConn, orden.CompraServiceClient, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	client := orden.NewCompraServiceClient(conn)
	return conn, client, nil
}

// ///////////////////////// Crear compra y enviar  ///////////////////////////
func realizarCompra(client orden.CompraServiceClient, input CompraInput) error {
	var vehiculos []*orden.Vehicle
	for _, v := range input.Vehiculos {
		vehiculos = append(vehiculos, &v)
	}

	compra := orden.Compra{
		Vehiculos: vehiculos,
		Cliente:   &input.Cliente,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := client.RealizarCompra(ctx, &compra)
	return err
}

func main() {
	////////////////////// Extraer JSON //////////////////////

	// Comprueba que existan 2 argumentos
	if len(os.Args) < 2 {
		log.Fatalf("%s <archivo.json>", os.Args[0])
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Error al abrir el archivo: %v", err)
	}
	defer file.Close()

	var input CompraInput
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&input); err != nil {
		log.Fatalf("Error al deserializar JSON: %v", err)
	}

	var conn *grpc.ClientConn
	var client orden.CompraServiceClient

	for {
		conn, client, err = conectarServer()
		if err != nil {
			time.Sleep(time.Second) // Espera antes de reintentar
			continue
		}

		if err = realizarCompra(client, input); err == nil {
			log.Printf("Compra realizada exitosamente.")
			break
		} else {
			conn.Close()            // Cerrar la conexión antes de reintentar
			time.Sleep(time.Second) // Esperar antes de reintentar
		}
	}

	if conn != nil {
		defer conn.Close() // Cerrar la conexión al final
	}
}
