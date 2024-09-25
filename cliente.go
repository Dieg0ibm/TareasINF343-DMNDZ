package main

import (
	"bufio"
	"bytes"
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

// ////////////////////////////////////////////////////FUNCIONES//////////////////////////////////////////////////////////////
func crearReserva() {
	scanner := bufio.NewScanner(os.Stdin)
	var idSala, idUsuario int
	var fecha, descripcion string

	fmt.Print("Ingrese el ID del usuario que realiza la reserva: ")
	fmt.Scanln(&idUsuario)
	fmt.Print("Ingrese el ID de la sala que se quiere reservar: ")
	fmt.Scanln(&idSala)
	fmt.Print("Ingrese la fecha en que quiere reservar la sala (Formato YYYY/MM/DD): ")
	fmt.Scanln(&fecha)
	fmt.Print("Ingrese una breve descripción de la reserva: ")
	if scanner.Scan() {
		descripcion = scanner.Text()
	}

	nuevaReserva := Reserva{
		ID_Sala:     idSala,
		ID_Usuario:  idUsuario,
		Fecha:       fecha,
		Descripcion: descripcion,
	}

	jsonData, err := json.Marshal(nuevaReserva)
	if err != nil {
		fmt.Println("Error al convertir la reserva a JSON:", err)
		return
	}

	response, err := http.Post("http://127.0.0.1:8080/api/reserva", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error al hacer la solicitud: ", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error al leer la respuesta del servidor: ", err)
		return
	}

	var createdReserva Reserva
	json.Unmarshal(body, &createdReserva)
	fmt.Println("Reserva creada exitosamente")
}

func verReservasPorUsuario() {
	var idUsuario int
	fmt.Print("Ingrese el ID del usuario: ")
	fmt.Scanln(&idUsuario)

	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/api/reserva?id_usuario=%d", idUsuario))
	if err != nil {
		fmt.Println("Error al hacer la solicitud:", err)
		return // Salir de la función, pero no del programa
	}
	defer response.Body.Close()

	// Verifica el código de estado HTTP
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)

		// Estructura para deserializar el mensaje de error
		var errorResponse struct {
			Message string `json:"message"`
		}
		json.Unmarshal(body, &errorResponse)

		fmt.Println(errorResponse.Message) // Mostrar solo el mensaje de error
		return                             // Salir de la función, pero no del programa
	}

	body, err := io.ReadAll(response.Body)

	// Deserializa el cuerpo de la respuesta
	var reservas []Reserva
	err = json.Unmarshal(body, &reservas)
	if err != nil {
		fmt.Println("Error al deserializar la respuesta:", err)
		return // Salir de la función, pero no del programa
	}

	// Verifica si hay reservas
	if len(reservas) == 0 {
		fmt.Println("No se encontraron reservas para este usuario.")
		return // Salir de la función, pero no del programa
	}

	for _, reserva := range reservas {
		fmt.Printf("Sala ID: %d, Fecha: %s, Descripción: %s\n", reserva.ID_Sala, reserva.Fecha, reserva.Descripcion)
	}
}

func verReservasPorSala() {
	var idSala int
	fmt.Print("Ingrese el ID de la sala: ")
	fmt.Scanln(&idSala)

	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/api/reserva?id_sala=%d", idSala))
	if err != nil {
		fmt.Println("Error al hacer la solicitud:", err)
		return // Salir de la función, pero no del programa
	}
	defer response.Body.Close()

	// Verifica el código de estado HTTP
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)

		// Estructura para deserializar el mensaje de error
		var errorResponse struct {
			Message string `json:"message"`
		}
		json.Unmarshal(body, &errorResponse)

		fmt.Println(errorResponse.Message) // Mostrar solo el mensaje de error
		return                             // Salir de la función, pero no del programa
	}

	body, err := io.ReadAll(response.Body)

	// Deserializa el cuerpo de la respuesta
	var reservas []Reserva
	err = json.Unmarshal(body, &reservas)
	if err != nil {
		fmt.Println("Error al deserializar la respuesta:", err)
		return // Salir de la función, pero no del programa
	}

	// Verifica si hay reservas
	if len(reservas) == 0 {
		fmt.Println("No se encontraron reservas para esta sala.")
		return // Salir de la función, pero no del programa
	}

	for _, reserva := range reservas {
		fmt.Printf("Usuario ID: %d, Fecha: %s, Descripción: %s\n", reserva.ID_Usuario, reserva.Fecha, reserva.Descripcion)
	}
}

func consultarReservaPorFecha() {
	var idSala int
	var fecha string

	fmt.Print("Ingrese el ID de la sala: ")
	fmt.Scanln(&idSala)
	fmt.Print("Ingrese la fecha (Formato YYYY/MM/DD): ")
	fmt.Scanln(&fecha)

	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/api/reserva?id_sala=%d&fecha=%s", idSala, fecha))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Intentamos decodificar la respuesta como un mensaje (cuando no hay reservas)
	var mensajeRespuesta map[string]string
	err = json.Unmarshal(body, &mensajeRespuesta)
	if err == nil {
		if msg, existe := mensajeRespuesta["message"]; existe {
			fmt.Println(msg) // Muestra el mensaje cuando no hay reservas
			return
		}
	}

	// Si no era un mensaje, intentamos decodificarlo como una lista de reservas
	var reservas []struct {
		IdSala      int    `json:"id_sala"`
		IdUsuario   int    `json:"id_usuario"`
		Fecha       string `json:"fecha"`
		Descripcion string `json:"descripcion"`
	}
	err = json.Unmarshal(body, &reservas)
	if err != nil {
		log.Fatal("Error al decodificar la respuesta del servidor:", err)
	}

	if len(reservas) == 0 {
		fmt.Println("La sala se encuentra disponible en la fecha consultada")
	} else {
		fmt.Println("La sala tiene una reserva activa en la fecha consultada")
	}
}

func cancelarReserva() {
	var idSala, idUsuario int
	var fecha string

	fmt.Print("Ingrese el ID del usuario que hizo la reserva: ")
	fmt.Scanln(&idUsuario)
	fmt.Print("Ingrese el ID de la sala reservada: ")
	fmt.Scanln(&idSala)
	fmt.Print("Ingrese la fecha de la reserva (Formato YYYY/MM/DD): ")
	fmt.Scanln(&fecha)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://127.0.0.1:8080/api/reserva?id_sala=%d&id_usuario=%d&fecha=%s", idSala, idUsuario, fecha), nil)
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Reserva cancelada con éxito.")
	} else {
		fmt.Println("Error al cancelar la reserva.")
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func main() {
	scanner := bufio.NewScanner(os.Stdin)
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
		fmt.Scanln(&opcion)

		///////////////////////////////////////////////////////////////////////////////////////////////////////////////
		switch opcion {
		case 1:
			var nombre string
			var departamento string
			var descripcion string

			fmt.Print("Ingrese el nombre del usuario: ")
			if scanner.Scan() {
				nombre = scanner.Text()
			}

			fmt.Print("Ingrese el departamento al que pertenece: ")
			if scanner.Scan() {
				departamento = scanner.Text()
			}

			fmt.Print("Ingrese una breve descripcion del usuario: ")
			fmt.Scanln(&descripcion)
			if scanner.Scan() {
				departamento = scanner.Text()
			}

			nuevoUsuario := Usuario{
				Nombre:       nombre,
				Departamento: departamento,
				Descripcion:  descripcion,
			}

			jsonData, err := json.Marshal(nuevoUsuario)
			if err != nil {
				fmt.Println("Error al convertir el usuario a JSON:", err)
				return
			}

			response, err := http.Post("http://127.0.0.1:8080/api/usuario", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Error al hacer la solicitud: ", err)
				return
			}
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error al leer la respuesta del servidor: ", err)
				return
			}

			var createdUsuario Usuario
			json.Unmarshal(body, &createdUsuario)

			fmt.Printf("Usuario creado exitosamente, su ID de usuario es: %d\n", createdUsuario.ID)
			///////////////////////////////////////////////////////////////////////////////////////////////////////////////
		case 2:
			var nombre_sala string
			var ubicacion_sala string

			fmt.Print("Ingrese el nombre de la sala : ")
			if scanner.Scan() {
				nombre_sala = scanner.Text()
			}

			fmt.Print("Ingrese la ubicacion de la sala : ")
			if scanner.Scan() {
				ubicacion_sala = scanner.Text()
			}

			nuevaSala := Sala{
				Nombre:    nombre_sala,
				Ubicacion: ubicacion_sala,
			}

			jsonData, err := json.Marshal(nuevaSala)
			if err != nil {
				fmt.Println("Error al convertir el usuario a JSON:", err)
				return
			}

			response, err := http.Post("http://127.0.0.1:8080/api/sala", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Error al hacer la solicitud: ", err)
				return
			}
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error al leer la respuesta del servidor: ", err)
				return
			}

			var salaCreada Usuario
			json.Unmarshal(body, &salaCreada)

			fmt.Printf("Sala creada exitosamente, el ID de la sala es: %d\n", salaCreada.ID)
			///////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
			///////////////////////////////////////////////////////////////////////////////////////////////////////////////
		case 4:
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
			///////////////////////////////////////////////////////////////////////////////////////////////////////////////
		case 5:
			fmt.Println("1. Crear una reserva")
			fmt.Println("2. Ver reservas por usuario")
			fmt.Println("3. Ver reservas por sala")
			fmt.Println("4. Consultar si sala tiene reserva en cierta fecha")
			fmt.Println("5. Cancelar una reserva")

			var subopcion int
			fmt.Print("Elige una opción: ")
			fmt.Scanln(&subopcion)

			switch subopcion {
			case 1:
				crearReserva()
			case 2:
				verReservasPorUsuario()
			case 3:
				verReservasPorSala()
			case 4:
				consultarReservaPorFecha()
			case 5:
				cancelarReserva()
			default:
				fmt.Println("Subopción inválida.")
			}

			///////////////////////////////////////////////////////////////////////////////////////////////////////////////
		case 6:
			fmt.Println("Saliendo del programa...")
			os.Exit(0) // Finaliza el programa
		default:
			fmt.Println("Opción inválida. Por favor, elige una opción entre 1 y 6.")
		}

		fmt.Println() // Espacio entre interacciones del menú
	}
}
