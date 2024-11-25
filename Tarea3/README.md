# Preguntas

### ¿Qué ocurre cuando n= 1 en ambos sistemas?

R: En ambos sistemas, todo se procesa de manera secuencial. No hay beneficios de concurrencia, ya que solo una goroutine procesa las solicitudes.

### ¿Qué ocurre cuando n= número de solicitudes en ambos sistemas?

R:
- *Sistema A*: Cada goroutine intenta acceder a la estructura compartida, pero debido al mutex global, solo una puede modificarla a la vez, lo que limita la concurrencia efectiva.
- *Sistema B*: Las goroutines pueden trabajar simultáneamente en diferentes paralelos, ya que cada paralelo tiene su propio mutex. Esto mejora la eficiencia si las solicitudes se distribuyen entre múltiples paralelos.

### ¿Como podría mejorar el sistema A para que sea más eficiente?
R:
- Dividir el mutex global en varios mutexes locales asociados a cada paralelo, como en el Sistema B. Esto permitiría concurrencia sobre diferentes recursos.


### ¿Como podría mejorar el sistema B para que sea más eficiente?
R:
- Colocar solicitudes en orden para que los procesos pueden acceder a ellas de mejor manera y reducir la espera.
- Utilizar técnicas de balanceo de carga para distribuir las solicitudes de manera más equitativa entre los paralelos.

### Escriba una analogía de los sistemas A y B con un sistema de inscripción presencial, donde:

- *Sistema A*:
  - Hay una sola persona que controla todos los cupos, y cada estudiante debe esperar a que esta termine antes de procesar otra solicitud.
- *Sistema B*:
  - Se asigna un encargado a cada paralelo que controla exclusivamente los cupos de este, permitiendo que múltiples estudiantes sean procesados simultáneamente en diferentes paralelos.
- *Diferencia*: El Sistema B aprovecha mejor los recursos humanos (encargados por paralelo) y permite mayor concurrencia.

# Diagramas de flujo
## Sistema A
![Imagen 1: Sistema A](Tarea3/sistema_A1.drawio.png)
## Sistema B
![Imagen 2: Sistema B](Tarea3/sistema_B1.drawio.png)
