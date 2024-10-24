// facturacion.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// //////////////////////////////////////// Estructuras //////////////////////////////////////////

type Vehicle struct {
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
	Price        int    `json:"price"`
}

type Customer struct {
	Id       string `json:"_id"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type Order struct {
	Id        string    `json:"_id"`
	OrderDate string    `json:"order_date"`
	Vehicles  []Vehicle `json:"vehicles"`
	Customer  Customer  `json:"customer"`
}

type Factura struct { // Se usa id generada automaticamente por Mongo
	Order_id     string `json:"order_id"`
	Amount       int    `json:"amount"`
	Invoice_date string `json:"invoice_date"`
}

// //////////////////////////////////////// Agregar factura a BD //////////////////////////////////////////
func CrearFactura(ctx context.Context, collection *mongo.Collection, order Order) error {
	total := 0
	for _, v := range order.Vehicles {
		total += int(v.Price)
	}

	factura := Factura{
		Order_id:     order.Id,
		Amount:       total,
		Invoice_date: order.OrderDate,
	}

	// Insertar en BD
	_, err := collection.InsertOne(ctx, factura)
	if err != nil {
		return fmt.Errorf("error al insertar la orden en MongoDB: %v", err)
	}

	return nil
}

func main() {

	//////////////// Conexión a MongoDB ////////////////
	ctx := context.Background()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Conectado a MongoDB")

	//////////////// Conexión a RabbitMQ ////////////////
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()
	ch, err := conn.Channel()
	defer ch.Close()

	//////////////// Definición de la cola ////////////////
	q, err := ch.QueueDeclare(
		"cola_facturacion", // Nombre de la cola
		true,               // Durable (persistente)
		false,              // Delete when unused
		false,              // Exclusive
		false,              // No-wait
		nil,                // Arguments
	)

	//////////////// Definición de mensajes ////////////////
	msgs, err := ch.Consume(
		q.Name, // Nombre de la cola
		"",     // Consumer tag (dejar vacío para auto-generar uno)
		true,   // Auto-ack (confirmación automática)
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Arguments
	)

	forever := make(chan bool) // Canal se mantiene corriendo

	coleccionesInvoices := client.Database("local").Collection("invoices")

	// Procesa mensajes
	go func() {
		for d := range msgs {
			log.Printf("Mensaje recibido: %s", d.Body)

			//  JSON a estructura de factura
			var order Order
			err := json.Unmarshal(d.Body, &order)
			if err != nil {
				log.Printf("Error: %s", err)
				continue
			}

			// Insertar la factura en DB
			err = CrearFactura(ctx, coleccionesInvoices, order)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println("Factura creada exitosamente")
		}
	}()

	log.Printf("Esperando mensajes... Para salir presione CTRL+C")
	<-forever
}
