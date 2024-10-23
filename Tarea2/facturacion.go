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

// Definición de la estructura de la factura
type Invoice struct {
	ID          string `json:"id"`
	OrderID     string `json:"order_id"`
	Amount      int64  `json:"amount"`
	InvoiceDate string `json:"invoice_date"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Guardar factura en MongoDB
func SaveInvoice(ctx context.Context, collection *mongo.Collection, invoice Invoice) error {
	_, err := collection.InsertOne(ctx, invoice)
	return err
}

func main() {
	// Conexión a MongoDB
	ctx := context.Background()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Conectado a MongoDB!")

	// Conexión a RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declarar la cola
	q, err := ch.QueueDeclare(
		"invoice_queue", // Nombre de la cola
		true,            // Durable (persistente)
		false,           // Delete when unused
		false,           // Exclusive
		false,           // No-wait
		nil,             // Arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Configurar el consumo de mensajes
	msgs, err := ch.Consume(
		q.Name, // Nombre de la cola
		"",     // Consumer tag (dejar vacío para auto-generar uno)
		true,   // Auto-ack (confirmación automática)
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Arguments
	)
	failOnError(err, "Failed to register a consumer")

	invoicesCollection := client.Database("local").Collection("invoices")

	// Canal para mantener el proceso corriendo
	forever := make(chan bool)

	// Goroutine para procesar los mensajes
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			// Deserializar el mensaje de JSON a la estructura de Invoice
			var invoice Invoice
			err := json.Unmarshal(d.Body, &invoice)
			if err != nil {
				log.Printf("Error deserializing message: %s", err)
				continue
			}

			// Guardar la factura en la base de datos
			err = SaveInvoice(ctx, invoicesCollection, invoice)
			if err != nil {
				log.Printf("Error saving invoice: %s", err)
				continue
			}
			fmt.Printf("Factura guardada: %+v\n", invoice)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
