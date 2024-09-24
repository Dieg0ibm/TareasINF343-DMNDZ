package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

///////////////////////////////////////////////////// ESTRUCTURAS /////////////////////////////////////////////////////

type Sala struct {
	ID        int    `json:"id"`
	Nombre    string `json:"nombre"`
	Ubicacion string `json:"ubicacion"`
}

type Usuario struct {
	ID           int    `json:"id"`
	Nombre       string `json:"nombre"`
	Departamento string `json:"departamento"`
	Descripcion  string `json:"descripcion"`
}

type Reserva struct {
	ID_Sala     int    `json:"id_sala"`
	ID_Usuario  int    `json:"id_usuario"`
	Fecha       string `json:"fecha"`
	Descripcion string `json:"descripcion"`
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func main() {
	for {
		fmt.Println("Menu")
		fmt.Println("1. Crear un usuario")
		fmt.Println("2. Crear una sala")
		fmt.Println("3. Consultar usuarios existentes")
		fmt.Println("4. Consultar salas existentes")
		fmt.Println("5. Administrar reservas")
		fmt.Println("6. Salir")
		fmt.Print("\nElige una opción: ")

		var opcion int
		_, err := fmt.Scan(&opcion) // Cambiar a fmt.Scan
		if err != nil {
			fmt.Println("Error: opción inválida.")
			// Limpiar el buffer de entrada en caso de error
			var dummy string
			fmt.Scanln(&dummy)
			continue
		}

		switch opcion {
		case 1:
			fmt.Println("Opción 1 seleccionada: Crear un usuario")
			// Aquí podrías implementar la lógica para crear un usuario
		case 2:
			fmt.Println("Opción 2 seleccionada: Crear una sala")
			// Aquí podrías implementar la lógica para crear una sala
		case 3:
			var usuarios []Usuario
			response, err := http.Get("http://127.0.0.1:8080/api/usuario")
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}

			// Deserializar el JSON en la estructura
			err = json.Unmarshal(body, &usuarios)
			if err != nil {
				log.Fatal(err)
			}

			// Imprimir los nombres de todos los usuarios
			for _, usuario := range usuarios {
				fmt.Printf("%s, ID = %d\n", usuario.Nombre, usuario.ID) // Imprime en el formato deseado
			}

		case 4:
			fmt.Println("Opción 4 seleccionada: Consultar salas existentes")
			response, err := http.Get("http://127.0.0.1:8080/api/sala")
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			// Deserializar y mostrar salas (aquí necesitarías definir la estructura Sala)
			var salas []Sala
			err = json.Unmarshal(body, &salas)
			if err != nil {
				log.Fatal(err)
			}
			for _, sala := range salas {
				fmt.Printf("%s, ID = %d\n", sala.Nombre, sala.ID)
			}

		case 5:
			fmt.Println("Opción 5 seleccionada: Administrar reservas")
			// Aquí podrías implementar la lógica para administrar reservas
		case 6:
			fmt.Println("Saliendo del programa...")
			os.Exit(0) // Finaliza el programa
		default:
			fmt.Println("Opción inválida. Por favor, elige una opción entre 1 y 6.")
		}

		fmt.Println() // Espacio entre interacciones del menú
	}
}
