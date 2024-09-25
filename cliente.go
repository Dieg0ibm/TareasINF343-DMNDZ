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

///////////////////////////////////////////////////// FUNCIONES AUX //////////////////////////////////////////////////////////////

// Se usa para verificar que la sala no esté ocupada en ese día
func consultarReservaPorFecha(idSala int, fecha string) string {

	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/api/reserva?id_sala=%d&fecha=%s", idSala, fecha))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	var reservas []Reserva

	json.Unmarshal(body, &reservas)

	if len(reservas) == 0 {
		mensaje := "La sala se encuentra disponible en la fecha consultada"
		return mensaje
	} else {
		mensaje := "La sala ya tiene una reserva en la fecha consultada"
		return mensaje
	}
}

////////////////////////////////////////////////////// FUNCIONES MENU //////////////////////////////////////////////////////////////

func MenucrearReserva() {
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

	if consultarReservaPorFecha(idSala, fecha) != "La sala ya tiene una reserva en la fecha consultada" {
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

		body, _ := io.ReadAll(response.Body)

		var createdReserva Reserva
		json.Unmarshal(body, &createdReserva)
		fmt.Println("Reserva creada exitosamente")
	} else {
		fmt.Println("La sala ya está reservada para esta fecha")
	}
}

func verReservasPorUsuario() {
	var idUsuario int
	fmt.Print("Ingrese el ID del usuario: ")
	fmt.Scanln(&idUsuario)

	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/api/reserva?id_usuario=%d", idUsuario))
	if err != nil {
		fmt.Println("Error al hacer la solicitud:", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	var reservas []Reserva
	json.Unmarshal(body, &reservas)

	if len(reservas) == 0 {
		fmt.Println("No se encontraron reservas para este usuario.")
		return
	}

	for _, reserva := range reservas {
		responseUsuario, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/api/getNombreUser?id=%d", reserva.ID_Usuario))
		if err != nil {
			log.Fatalf("Error al obtener el nombre del usuario: %v", err)
		}
		defer responseUsuario.Body.Close()
		nombre_usuario, err := io.ReadAll(responseUsuario.Body)

		responseSala, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/api/getNombreSala?id=%d", reserva.ID_Sala))
		if err != nil {
			log.Fatalf("Error al obtener el nombre de la sala: %v", err)
		}
		defer responseSala.Body.Close()
		nombre_sala, err := io.ReadAll(responseSala.Body)

		fmt.Printf("%s tiene reserva en la sala %s en %s", nombre_usuario, nombre_sala, reserva.Fecha)
	}
}

func verReservasPorSala() {
	var idSala int
	fmt.Print("Ingrese el ID de la sala: ")
	fmt.Scanln(&idSala)

	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/api/reserva?id_sala=%d", idSala))
	if err != nil {
		fmt.Println("Error al hacer la solicitud:", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	var reservas []Reserva
	json.Unmarshal(body, &reservas)

	if len(reservas) == 0 {
		fmt.Println("No se encontraron reservas para este usuario.")
		return
	}

	for _, reserva := range reservas {
		responseUsuario, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/api/getNombreUser?id=%d", reserva.ID_Usuario))
		if err != nil {
			log.Fatalf("Error al obtener el nombre del usuario: %v", err)
		}
		defer responseUsuario.Body.Close()
		nombre_usuario, err := io.ReadAll(responseUsuario.Body)

		responseSala, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/api/getNombreSala?id=%d", reserva.ID_Sala))
		if err != nil {
			log.Fatalf("Error al obtener el nombre de la sala: %v", err)
		}
		defer responseSala.Body.Close()
		nombre_sala, err := io.ReadAll(responseSala.Body)
		fmt.Printf("La sala %s tiene una reserva de %s en %s", nombre_sala, nombre_usuario, reserva.Fecha)
	}
}

func MenuconsultarReservaPorFecha() string {
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

	var reservas []Reserva

	json.Unmarshal(body, &reservas)

	if len(reservas) == 0 {
		mensaje := "La sala se encuentra disponible en la fecha consultada"
		return mensaje
	} else {
		mensaje := "La sala ya tiene una reserva en la fecha consultada"
		return mensaje
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
		fmt.Println("No se encontraron reservas con los datos seleccionados")
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
			if scanner.Scan() {
				departamento = scanner.Text()
			}

			nuevoUsuario := Usuario{
				Nombre:       nombre,
				Departamento: departamento,
				Descripcion:  descripcion,
			}

			jsonData, _ := json.Marshal(nuevoUsuario)

			response, err := http.Post("http://127.0.0.1:8080/api/usuario", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Error al hacer la solicitud: ", err)
				return
			}

			defer response.Body.Close()

			body, _ := io.ReadAll(response.Body)

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

			jsonData, _ := json.Marshal(nuevaSala)

			response, err := http.Post("http://127.0.0.1:8080/api/sala", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Error al hacer la solicitud: ", err)
				return
			}
			defer response.Body.Close()

			body, _ := io.ReadAll(response.Body)

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

			// JSON a Usuario
			json.Unmarshal(body, &usuarios)

			for _, usuario := range usuarios {
				fmt.Printf("%s, ID = %d\n", usuario.Nombre, usuario.ID)
			}
			///////////////////////////////////////////////////////////////////////////////////////////////////////////////
		case 4:
			response, err := http.Get("http://127.0.0.1:8080/api/sala")
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)

			var salas []Sala
			json.Unmarshal(body, &salas)
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
				MenucrearReserva()
			case 2:
				verReservasPorUsuario()
			case 3:
				verReservasPorSala()
			case 4:
				fmt.Println(MenuconsultarReservaPorFecha())
			case 5:
				cancelarReserva()
			default:
				fmt.Println("Opción inválida.")
			}

			///////////////////////////////////////////////////////////////////////////////////////////////////////////////
		case 6:
			fmt.Println("Fin del programa!")
			os.Exit(0)
		default:
			fmt.Println("Opción inválida. Por favor, elige una opción entre 1 y 6.")
		}

		fmt.Println()
	}
}
