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

// ////////////////////////////////// Estructuras ////////////////////////////////////

type CompraInput struct {
	Vehiculos []orden.Vehicle `json:"vehicles"`
	Cliente   orden.Customer  `json:"customer"`
}

func main() {
	//////////////////////////////////// Conectar al servidor gRPC ////////////////////////////////////

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close()

	client := orden.NewCompraServiceClient(conn)

	//////////////////////////////////// Leer el archivo JSON ////////////////////////////////////

	// Comprobar paso de argumento
	if len(os.Args) < 2 {
		log.Fatalf("%s <archivo.json>", os.Args[0])
	}

	// Abrir JSON
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

	var vehiculos []*orden.Vehicle
	for _, v := range input.Vehiculos {
		vehiculos = append(vehiculos, &v)
	}

	//////////////////////////////////// Asignar los datos a la compra ////////////////////////////////////

	compra := orden.Compra{
		Vehiculos: vehiculos,
		Cliente:   &input.Cliente,
	}

	//////////////////////////////////// LLamada gRPC ////////////////////////////////////
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := client.RealizarCompra(ctx, &compra)
	if err != nil {
		log.Fatalf("Error al realizar la compra: %v", err)
	}

	log.Printf("Respuesta del servidor: %s", response.Message)
}
