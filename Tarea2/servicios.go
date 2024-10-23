package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Definición de una estructura para representar los vehículos
type Vehicle struct {
	Manufacturer string  `json:"manufacturer"`
	Model        string  `json:"model"`
	Year         int     `json:"year"`
	Price        float64 `json:"price"`
}

type Order struct {
	OrderDate string    `json:"order_date"`
	Vehicles  []Vehicle `json:"vehicles"`
	Customer  struct {
		Name     string `json:"name"`
		Lastname string `json:"lastname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	} `json:"customer"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Actualizar Stock
func ActualizarStock(ctx context.Context, collection *mongo.Collection, vehicle Vehicle) error {
	// Filtro para encontrar el vehículo específico en la base de datos.
	filter := bson.D{
		{Key: "manufacturer", Value: vehicle.Manufacturer},
		{Key: "model", Value: vehicle.Model},
		{Key: "year", Value: vehicle.Year},
		{Key: "price", Value: vehicle.Price},
	}

	// Definir la actualización para disminuir el stock.
	update := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "stock", Value: -1}, // Decrementar el stock en 1
		}},
	}

	// Realizar la actualización.
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating stock for %s %s (%d): %w", vehicle.Manufacturer, vehicle.Model, vehicle.Year, err)
	}

	// Comprobar el resultado de la actualización.
	if result.ModifiedCount == 0 {
		return fmt.Errorf("no se encontró el vehículo %s %s (%d) o no se pudo actualizar el stock", vehicle.Manufacturer, vehicle.Model, vehicle.Year)
	}

	return nil
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
		"order_queue", // Nombre de la cola
		true,          // Durable (persistente)
		false,         // Delete when unused
		false,         // Exclusive
		false,         // No-wait
		nil,           // Arguments
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

	vehiclesCollection := client.Database("local").Collection("vehicles")

	// Canal para mantener el proceso corriendo
	forever := make(chan bool)

	// Goroutine para procesar los mensajes
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			// Deserializar el mensaje de JSON a la estructura de Order
			var order Order
			err := json.Unmarshal(d.Body, &order)
			if err != nil {
				log.Printf("Error deserializing message: %s", err)
				continue
			}

			// Procesar la orden
			fmt.Println("Nueva orden recibida:")
			fmt.Printf("Fecha: %s\n", order.OrderDate)
			fmt.Printf("Cliente: %s %s\n", order.Customer.Name, order.Customer.Lastname)
			fmt.Printf("Vehículos:\n")
			for _, v := range order.Vehicles {
				fmt.Printf("- %s %s (%d) - Precio: %.2f\n", v.Manufacturer, v.Model, v.Year, v.Price)

				// Actualizar el stock del vehículo
				err := ActualizarStock(ctx, vehiclesCollection, v)
				if err != nil {
					log.Println(err) // Muestra el error si no se puede actualizar el stock
					continue
				}
				fmt.Printf("Stock actualizado para %s %s (%d).\n", v.Manufacturer, v.Model, v.Year)
			}
			fmt.Printf("Email: %s, Teléfono: %s\n", order.Customer.Email, order.Customer.Phone)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
