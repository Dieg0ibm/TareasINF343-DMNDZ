// servicios.go
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

// //////////////////////////////////////// Estructuras //////////////////////////////////////////
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

// //////////////////////////////////////// Actualizar stock //////////////////////////////////////////
func ActualizarStock(ctx context.Context, collection *mongo.Collection, vehicle Vehicle) error {
	filtro := bson.D{
		{Key: "manufacturer", Value: vehicle.Manufacturer},
		{Key: "model", Value: vehicle.Model},
		{Key: "year", Value: vehicle.Year},
		{Key: "price", Value: vehicle.Price},
	}

	update := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "stock", Value: -1},
		}},
	}

	// Disminuir Stock
	result, err := collection.UpdateOne(ctx, filtro, update)
	if err != nil {
		return fmt.Errorf("Error actualizando el stock de %s %s (%d): %w", vehicle.Manufacturer, vehicle.Model, vehicle.Year, err)
	}

	// Comprobar
	if result.ModifiedCount == 0 {
		return fmt.Errorf("No se encontró %s %s (%d) o no se pudo actualizar el stock", vehicle.Manufacturer, vehicle.Model, vehicle.Year)
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
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Conectado a MongoDB!")

	//////////////// Conexión a RabbitMQ ////////////////
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()

	ch, err := conn.Channel()
	defer ch.Close()

	//////////////// Definición de la cola ////////////////
	q, err := ch.QueueDeclare(
		"cola_inventario", // Nombre de la cola
		true,              // Durable (persistente)
		false,             // Delete when unused
		false,             // Exclusive
		false,             // No-wait
		nil,               // Arguments
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
	forever := make(chan bool)

	coleccionAutos := client.Database("local").Collection("vehicles")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			//  JSON a Order
			var order Order
			err := json.Unmarshal(d.Body, &order)
			if err != nil {
				log.Printf("Error deserializing message: %s", err)
				continue
			}

			for _, v := range order.Vehicles {
				fmt.Printf("- %s %s (%d) - Precio: %.2f\n", v.Manufacturer, v.Model, v.Year, v.Price)

				// Actualizar el stock del vehículo
				err := ActualizarStock(ctx, coleccionAutos, v)
				if err != nil {
					log.Println(err)
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
