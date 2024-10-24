# Tarea 2 INF343-DMNDZ

Integrantes:
  - Diego Bahamondes, ROL: 202173052-4
  - Maximiliano Johnson, ROL: 202173097-4
  - Claudio Varela, ROL: 202141087-2

Instrucciones de uso:

- IP de la maquina virtual utilizada: 10.10.29.228
- Ir a la carpeta TareasINF343-DMNDZ/Tarea2
- Ejecutar cada uno de los siguientes comandos en terminales distintas (desde la shell utilice el comando 'screen' para tener acceso a cuatro terminales)
- **Ejecutar el servicio de gestión de inventario:** go run path_3/inventario.go
- **Ejecutar el servicio de facturación:** go run path_3/facturacion.go
- **Ejecutar el servicio de compras:** python3 path_2/compras.py 
- **Ejecutar el cliente leyendo el archivo vehicles.json:** ./path vehicles.json



Consideraciones:
- Versión de Go: 1.23.2
- El programa asume que los vehículos siempre tendrán stock disponible.
- Se ha detectado un error ocasional en la máquina virtual donde el canal de gRPC se cierra inesperadamente. Este problema es intermitente y no ocurre en todos los entornos, pero si se realizan **varios intentos** el programa funciona correctamente, actualizando la base de datos.
Se adjunta el mensaje de error: *rpc error: code Unknown desc = Exception calling application: Channel is closed*.


