package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Token struct {
	LN []int
}

type Paralelo struct {
	Nombre string
	Cupos  int
	mu     sync.Mutex // Mutex específico para cada paralelo
}

type Proceso struct {
	ID         int
	RN         []int
	Token      *Token
	TieneToken bool
	Solicitud  chan int
}

var (
	procesos     []*Proceso
	paralelos    map[string]*Paralelo
	mutexArchivo sync.Mutex // Mutex para proteger escrituras en archivos
	wg           sync.WaitGroup
)

func cargarParalelos(nombreArchivo string) (map[string]*Paralelo, error) {
	datos, err := os.ReadFile(nombreArchivo)
	if err != nil {
		return nil, err
	}
	lineas := strings.Split(string(datos), "\n")
	paralelos := make(map[string]*Paralelo)
	for _, linea := range lineas {
		if strings.TrimSpace(linea) == "" {
			continue
		}
		partes := strings.Fields(linea)
		if len(partes) != 2 {
			continue
		}
		nombre := strings.TrimSpace(partes[0])
		cupos, _ := strconv.Atoi(strings.TrimSpace(partes[1]))
		paralelos[nombre] = &Paralelo{Nombre: nombre, Cupos: cupos}
	}
	return paralelos, nil
}

// Extraer solicitud
func extraerSolicitud() (string, []string, error) {
	mutexArchivo.Lock() // Sección crítica
	defer mutexArchivo.Unlock()

	datos, err := os.ReadFile("solicitudes.txt")
	if err != nil {
		return "", nil, err
	}
	lineas := strings.Split(string(datos), "\n")

	for i, linea := range lineas {
		if strings.TrimSpace(linea) != "" {
			partes := strings.Fields(linea)
			if len(partes) < 2 {
				continue
			}
			estudiante := strings.TrimSpace(partes[0])
			preferencias := partes[1:]

			lineas = append(lineas[:i], lineas[i+1:]...)
			err = os.WriteFile("solicitudes.txt", []byte(strings.Join(lineas, "\n")), 0644)
			if err != nil {
				return "", nil, err
			}
			return estudiante, preferencias, nil
		}
	}
	return "", nil, fmt.Errorf("No hay solicitudes disponibles")
}

// Inscribir estudiante en un paralelo
func inscribirEstudiante(estudiante string, preferencias []string) (string, error) {
	for _, paralelo := range preferencias {
		if p, ok := paralelos[paralelo]; ok {
			p.mu.Lock()
			if p.Cupos > 0 {
				p.Cupos--
				p.mu.Unlock()
				fmt.Printf("Estudiante %s inscrito en paralelo %s. Cupos restantes: %d\n", estudiante, paralelo, p.Cupos)
				return fmt.Sprintf("%s %s", estudiante, paralelo), nil
			}
			p.mu.Unlock()
			fmt.Printf("Paralelo %s no tiene cupos disponibles.\n", paralelo)
		} else {
			fmt.Printf("Paralelo %s no existe.\n", paralelo)
		}
	}
	return fmt.Sprintf("%s no pudo ser inscrito (sin cupos disponibles)", estudiante), nil
}

// Registrar resultado en inscritos.txt
func registrarLog(idProceso int, resultado string) error {
	mutexArchivo.Lock()
	defer mutexArchivo.Unlock()

	archivo, err := os.OpenFile("inscritos.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer archivo.Close()

	entrada := fmt.Sprintf("P%d: %s\n", idProceso+1, resultado)
	_, err = archivo.WriteString(entrada)
	return err
}

// Manejar solicitudes para cada proceso
func (p *Proceso) ManejarSolicitudes() {
	defer wg.Done()

	for {
		estudiante, preferencias, err := extraerSolicitud()
		if err != nil {
			fmt.Printf("P%d: %v\n", p.ID+1, err)
			break
		}

		// Esperar el token
		for !p.TieneToken {
		}

		resultado, err := inscribirEstudiante(estudiante, preferencias)
		if err != nil {
			fmt.Printf("P%d error al inscribir: %v\n", p.ID+1, err)
			return
		}

		if err := registrarLog(p.ID, resultado); err != nil {
			fmt.Printf("P%d error al registrar log: %v\n", p.ID+1, err)
			return
		}

		// Liberar token
		p.Token.LN[p.ID] = p.RN[p.ID]
		tokenTransferido := false

		// Transferir token al siguiente proceso disponible
		for i := 1; i < len(procesos); i++ {
			siguienteID := (p.ID + i) % len(procesos)
			if procesos[siguienteID].TieneToken == false {
				procesos[siguienteID].Token = p.Token
				procesos[siguienteID].TieneToken = true
				fmt.Printf("P%d transfiere el token a P%d\n", p.ID+1, siguienteID+1)
				tokenTransferido = true
				break
			}
		}

		if !tokenTransferido {
			p.TieneToken = true
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Uso: go run main.go <número de procesos>")
		return
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil || n <= 0 {
		fmt.Println("El número de procesos debe ser un entero positivo.")
		return
	}

	paralelos, err = cargarParalelos("paralelos.txt")
	if err != nil {
		fmt.Printf("Error al cargar paralelos: %v\n", err)
		return
	}

	procesos = make([]*Proceso, n)
	for i := 0; i < n; i++ {
		procesos[i] = &Proceso{
			ID:         i,
			RN:         make([]int, n),
			Token:      nil,
			TieneToken: i == 0,
			Solicitud:  make(chan int, n),
		}
		if i == 0 {
			procesos[i].Token = &Token{LN: make([]int, n)}
		}
	}

	for _, proc := range procesos {
		wg.Add(1)
		go proc.ManejarSolicitudes()
	}

	wg.Wait()
	fmt.Println("Ejecución completada.")
}
