from concurrent import futures
import grpc
import orden_pb2_grpc
import orden_pb2
from pymongo import MongoClient
from datetime import datetime

class CompraService(orden_pb2_grpc.CompraServiceServicer):
    def __init__(self):
        # Conectar a la base de datos de MongoDB
        self.client = MongoClient('mongodb://localhost:27017/')
        self.db = self.client['local']  # Cambia 'local' por el nombre de tu base de datos
        self.inventory = self.db['vehicles']  # Cambia 'vehicles' por el nombre de tu colección de precios
        self.orders = self.db['orders']  # Cambia 'orders' por el nombre de tu colección de órdenes

    def RealizarCompra(self, request, context):
        # Acceder a los datos del mensaje
        customer = request.cliente
        vehicles = request.vehiculos
        
        # Inicializar descripciones de vehículos y precios
        vehiculo_descripciones = []
        precios = []

        # Obtener el precio de cada vehículo de la base de datos
        for v in vehicles:
            resultado = self.inventory.find_one({"manufacturer": v.manufacturer, "model": v.model, "year": v.year})
            if resultado:
                precio = resultado.get('price', 0)
                precios.append(precio)  # Guardar el precio en la lista de precios
                vehiculo_descripciones.append(f"{v.manufacturer} {v.model} ({v.year}) - Precio: {precio}")
            else:
                precios.append(0)  # Si no se encuentra el vehículo, agregar 0
                vehiculo_descripciones.append(f"{v.manufacturer} {v.model} ({v.year}) - Precio no disponible")

        # Construir el mensaje de respuesta
        mensaje = f"Compra realizada por {customer.name} {customer.lastname} para los vehículos: {', '.join(vehiculo_descripciones)}"

        # Crear la orden para persistirla en la base de datos
        order_data = {
            "order_date": datetime.now().strftime("%Y/%m/%d"),  # Fecha actual en formato YYYY/MM/DD
            "vehicles": [
                {
                    "manufacturer": v.manufacturer,
                    "model": v.model,
                    "year": v.year,
                    "price": precios[i]  # Usar el precio correspondiente de la lista
                } for i, v in enumerate(vehicles)
            ],
            "customer": {
                "name": customer.name,
                "lastname": customer.lastname,
                "email": customer.email,
                "phone": customer.phone
            }
        }

        # Insertar la orden en la colección 'orders'
        self.orders.insert_one(order_data)  # Corregido el acceso a la colección de órdenes

        return orden_pb2.CompraResponse(message=mensaje)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    orden_pb2_grpc.add_CompraServiceServicer_to_server(CompraService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Servidor gRPC iniciado en el puerto 50051")
    server.wait_for_termination()

if __name__ == "__main__":
    serve()
