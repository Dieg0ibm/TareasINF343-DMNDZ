package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"Tarea2/orden" // Asegúrate de importar la ruta correcta de tu paquete

	"google.golang.org/grpc"
)

func main() {
	// Conectar al servidor gRPC
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close()

	client := orden.NewCompraServiceClient(conn)

	// Leer el archivo JSON
	file, err := os.Open("a.json")
	if err != nil {
		log.Fatalf("Error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Crear una estructura auxiliar para deserializar el JSON
	type CompraInput struct {
		Vehiculos []orden.Vehicle `json:"vehicles"` // Este es el campo en el JSON
		Cliente   orden.Customer  `json:"customer"` // Este es el campo en el JSON
	}

	var input CompraInput
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&input); err != nil {
		log.Fatalf("Error al decodificar JSON: %v", err)
	}

	// Crear un slice de punteros a vehículos
	var vehiculos []*orden.Vehicle
	for _, v := range input.Vehiculos {
		vehiculos = append(vehiculos, &v) // Añadir puntero a cada vehículo
	}

	// Asignar los datos a la variable compra
	compra := orden.Compra{
		Vehiculos: vehiculos,      // Slice de punteros a vehículos
		Cliente:   &input.Cliente, // Puntero al cliente
	}

	// Hacer la llamada al servicio gRPC
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.RealizarCompra(ctx, &compra)
	if err != nil {
		log.Fatalf("Error al realizar la compra: %v", err)
	}

	log.Printf("Respuesta del servidor: %s", response.Message)
}
