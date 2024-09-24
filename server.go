package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

///////////////////////////////////////////////////// BASE DE DATOS /////////////////////////////////////////////////////

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./ServicioDeReservas.db")
	if err != nil {
		log.Fatal(err)
	}
}

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

///////////////////////////////////////////////////// POST /////////////////////////////////////////////////////

// Crear Sala
func crearSala(c *gin.Context) {
	var newSala Sala
	if err := c.BindJSON(&newSala); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := db.Exec(`INSERT INTO salas (nombre, ubicacion) VALUES (?, ?)`, newSala.Nombre, newSala.Ubicacion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newSala.ID = int(id)
	c.IndentedJSON(http.StatusCreated, newSala)
}

// Crear Usuario
func crearUsuario(c *gin.Context) {
	var newUsuario Usuario
	if err := c.BindJSON(&newUsuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := db.Exec(`INSERT INTO usuarios (nombre, departamento, descripcion) VALUES (?, ?, ?)`, newUsuario.Nombre, newUsuario.Departamento, newUsuario.Descripcion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newUsuario.ID = int(id)
	c.IndentedJSON(http.StatusCreated, newUsuario)
}

// Crear Reserva
func crearReserva(c *gin.Context) {
	var newReserva Reserva
	if err := c.BindJSON(&newReserva); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := db.Exec(`INSERT INTO reservas (id_sala, id_usuario, fecha, descripcion) VALUES (?, ?, ?, ?)`,
		newReserva.ID_Sala, newReserva.ID_Usuario, newReserva.Fecha, newReserva.Descripcion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newReserva)
}

///////////////////////////////////////////////////// GET /////////////////////////////////////////////////////

// Obtener todas las salas
func obtenerSalas(c *gin.Context) {
	rows, err := db.Query(`SELECT id, nombre, ubicacion FROM salas`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var salas []Sala
	for rows.Next() {
		var sala Sala
		if err := rows.Scan(&sala.ID, &sala.Nombre, &sala.Ubicacion); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		salas = append(salas, sala)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, salas)
}

// Obtener todos los usuarios
func obtenerUsuarios(c *gin.Context) {
	rows, err := db.Query(`SELECT id, nombre, departamento, descripcion FROM usuarios`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var usuarios []Usuario
	for rows.Next() {
		var usuario Usuario
		if err := rows.Scan(&usuario.ID, &usuario.Nombre, &usuario.Departamento, &usuario.Descripcion); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		usuarios = append(usuarios, usuario)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, usuarios)
}

// Obtener reservas
func obtenerReservas(c *gin.Context) {
	idSala := c.Query("id_sala")
	idUsuario := c.Query("id_usuario")
	fecha := c.Query("fecha")

	// Validar que al menos uno de los parámetros esté presente
	if idSala == "" && idUsuario == "" && fecha == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere al menos uno de los parámetros: 'id_sala', 'id_usuario' o 'fecha'."})
		return
	}

	// Lógica modularizada según los parámetros
	if idSala != "" && idUsuario != "" && fecha != "" {
		obtenerReservaPorSalaUsuarioYFecha(c, idSala, idUsuario, fecha)
	} else if idSala != "" && fecha != "" {
		obtenerReservaPorSalaYFecha(c, idSala, fecha)
	} else if idSala != "" {
		obtenerReservaPorSala(c, idSala)
	} else if idUsuario != "" {
		obtenerReservaPorUsuario(c, idUsuario)
	}
}

// Función para obtener reservas por id_sala, id_usuario y fecha
func obtenerReservaPorSalaUsuarioYFecha(c *gin.Context, idSala string, idUsuario string, fecha string) {
	query := `SELECT id_sala, id_usuario, fecha, descripcion FROM reservas WHERE id_sala = ? AND id_usuario = ? AND fecha = ?`
	rows, err := db.Query(query, idSala, idUsuario, fecha)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var reservas []Reserva
	for rows.Next() {
		var reserva Reserva
		if err := rows.Scan(&reserva.ID_Sala, &reserva.ID_Usuario, &reserva.Fecha, &reserva.Descripcion); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		reservas = append(reservas, reserva)
	}

	if len(reservas) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No se encontraron reservas para los parámetros proporcionados"})
		return
	}

	c.IndentedJSON(http.StatusOK, reservas)
}

// Función para obtener reservas por id_sala y fecha
func obtenerReservaPorSalaYFecha(c *gin.Context, idSala string, fecha string) {
	query := `SELECT id_sala, id_usuario, fecha, descripcion FROM reservas WHERE id_sala = ? AND fecha = ?`
	rows, err := db.Query(query, idSala, fecha)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var reservas []Reserva
	for rows.Next() {
		var reserva Reserva
		if err := rows.Scan(&reserva.ID_Sala, &reserva.ID_Usuario, &reserva.Fecha, &reserva.Descripcion); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		reservas = append(reservas, reserva)
	}

	// Verificar errores durante el procesamiento de las filas
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(reservas) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No se encontraron reservas para los parámetros proporcionados"})
		return
	}

	c.IndentedJSON(http.StatusOK, reservas)
}

// Función para obtener reservas por id_sala
func obtenerReservaPorSala(c *gin.Context, idSala string) {
	query := `SELECT id_sala, id_usuario, fecha, descripcion FROM reservas WHERE id_sala = ?`
	rows, err := db.Query(query, idSala)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var reservas []Reserva
	for rows.Next() {
		var reserva Reserva
		if err := rows.Scan(&reserva.ID_Sala, &reserva.ID_Usuario, &reserva.Fecha, &reserva.Descripcion); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		reservas = append(reservas, reserva)
	}

	if len(reservas) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No se encontraron reservas para la sala especificada"})
		return
	}

	c.IndentedJSON(http.StatusOK, reservas)
}

// Función para obtener reservas por id_usuario
func obtenerReservaPorUsuario(c *gin.Context, idUsuario string) {
	query := `SELECT id_sala, id_usuario, fecha, descripcion FROM reservas WHERE id_usuario = ?`
	rows, err := db.Query(query, idUsuario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var reservas []Reserva
	for rows.Next() {
		var reserva Reserva
		if err := rows.Scan(&reserva.ID_Sala, &reserva.ID_Usuario, &reserva.Fecha, &reserva.Descripcion); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		reservas = append(reservas, reserva)
	}

	if len(reservas) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No se encontraron reservas para el usuario especificado"})
		return
	}

	c.IndentedJSON(http.StatusOK, reservas)
}

func main() {
	r := gin.Default()

	// Rutas para manejar salas
	r.POST("/api/sala", crearSala)
	r.GET("/api/sala", obtenerSalas)

	// Rutas para manejar usuarios
	r.POST("/api/usuario", crearUsuario)
	r.GET("/api/usuario", obtenerUsuarios)

	// Rutas para manejar reservas
	r.POST("/api/reserva", crearReserva)
	r.GET("/api/reserva", obtenerReservas)

	// Iniciar el servidor en el puerto 8080
	r.Run("127.0.0.1:8080")
}
