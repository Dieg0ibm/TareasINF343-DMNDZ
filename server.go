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
	db, err = sql.Open("sqlite3", "./Reservas.db")
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
	var nuevaSala Sala
	if err := c.BindJSON(&nuevaSala); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := db.Exec(`INSERT INTO salas (nombre, ubicacion) VALUES (?, ?)`, nuevaSala.Nombre, nuevaSala.Ubicacion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	nuevaSala.ID = int(id)
	c.IndentedJSON(http.StatusCreated, nuevaSala)
}

// Crear Usuario
func crearUsuario(c *gin.Context) {
	var nuevoUsuario Usuario
	if err := c.BindJSON(&nuevoUsuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := db.Exec(`INSERT INTO usuarios (nombre, departamento, descripcion) VALUES (?, ?, ?)`, nuevoUsuario.Nombre, nuevoUsuario.Departamento, nuevoUsuario.Descripcion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	nuevoUsuario.ID = int(id)
	c.IndentedJSON(http.StatusCreated, nuevoUsuario)
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

func obtenerReservas(c *gin.Context) {
	idSala := c.Query("id_sala")
	idUsuario := c.Query("id_usuario")
	fecha := c.Query("fecha")

	if idSala == "" && idUsuario == "" && fecha == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere al menos uno de los parámetros: 'id_sala', 'id_usuario' o 'fecha'."})
		return
	}

	// Endpoint depende de los datos
	if idSala != "" && fecha != "" {
		obtenerReservaPorSalaYFecha(c, idSala, fecha)
	} else if idSala != "" {
		obtenerReservaPorSala(c, idSala)
	} else if idUsuario != "" {
		obtenerReservaPorUsuario(c, idUsuario)
	}
}

// Obtener reservas por id_sala y fecha
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

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusNotFound, gin.H{"message": "La sala se encuentra disponible en la fecha consultada"})
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
		c.JSON(http.StatusNotFound, gin.H{"message": "El usuario no tiene reservas"})
		return
	}

	c.IndentedJSON(http.StatusOK, reservas)
}

/////////////////////////////////////////// DELETE ///////////////////////////////////////////

func eliminarReserva(c *gin.Context) {
	idSala := c.Query("id_sala")
	idUsuario := c.Query("id_usuario")
	fecha := c.Query("fecha")

	// Verificar si la reserva existe
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM reservas WHERE id_sala = ? AND id_usuario = ? AND fecha = ?)`
	err := db.QueryRow(checkQuery, idSala, idUsuario, fecha).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "No se encontró la reserva"})
		return
	}

	// Eliminar la reserva si existe
	deleteQuery := `DELETE FROM reservas WHERE id_sala = ? AND id_usuario = ? AND fecha = ?`
	_, err = db.Exec(deleteQuery, idSala, idUsuario, fecha)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reserva eliminada exitosamente"})
}

func main() {
	r := gin.Default()

	// Rutas salas
	r.POST("/api/sala", crearSala)
	r.GET("/api/sala", obtenerSalas)

	// Rutas usuarios
	r.POST("/api/usuario", crearUsuario)
	r.GET("/api/usuario", obtenerUsuarios)

	// Rutas reservas
	r.POST("/api/reserva", crearReserva)
	r.GET("/api/reserva", obtenerReservas)
	r.DELETE("/api/reserva", eliminarReserva)

	// Rutas auxiliares
	r.GET("/api/getNombreUser", getNombreUsuario)
	r.GET("/api/getNombreSala", getNombreSala)

	// Iniciar el servidor en el puerto 8080
	r.Run("127.0.0.1:8080")
}

///////////////////////////////////////////////////// GET AUXILIARES (Para prints) /////////////////////////////////////////////////////

func getNombreUsuario(c *gin.Context) {
	id := c.Query("id")
	row := db.QueryRow(`SELECT nombre FROM usuarios WHERE id = ?`, id)

	var nombre string
	if err := row.Scan(&nombre); err != nil {
		return
	}

	c.String(http.StatusOK, nombre)
}

func getNombreSala(c *gin.Context) {
	id := c.Query("id")
	row := db.QueryRow(`SELECT nombre FROM salas WHERE id = ?`, id)

	var nombre string
	if err := row.Scan(&nombre); err != nil {
		return
	}

	c.String(http.StatusOK, nombre)
}
