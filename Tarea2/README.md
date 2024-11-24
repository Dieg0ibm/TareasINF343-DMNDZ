# Tarea 2 INF343-DMNDZ

Integrantes:
  - Diego Bahamondes, ROL: 202173052-4
  - Maximiliano Johnson, ROL: 202173097-4
  - Claudio Varela, ROL: 202141087-2

Instrucciones de uso:


EJECUTAR EN EL SIGUIENTE ORDEN Y DESDE LA CARPETA TareasINF343-DMNDZ/Tarea2/

MV 3 (10.10.29.228) "Inventario y facturacion"

- go run path_3/inventario.go
- go run path_3/facturacion.go

MV 2 (10.10.29.233) "Servicio compras"

- python3 path_2/compras.py

MV 1 (10.10.29.227) "Cliente"

- ./Tarea2 vehicles.json

Puedes revisar la base de datos en la MV 2 con mongosh, ocupamos la db "local"



Consideraciones:
- Versión de Go: 1.23.2
- El programa asume que los vehículos siempre tendrán stock disponible.
- Se ha detectado un error ocasional en la máquina virtual donde el canal de gRPC se cierra inesperadamente. Este problema es intermitente y no ocurre en todos los entornos, pero si se realizan **varios intentos** el programa funciona correctamente, actualizando la base de datos.
Se adjunta el mensaje de error: *rpc error: code Unknown desc = Exception calling application: Channel is closed*.


