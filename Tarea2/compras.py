from concurrent import futures
import grpc
import orden_pb2_grpc
import orden_pb2
from pymongo import MongoClient
from datetime import datetime

import pika
import json

class CompraService(orden_pb2_grpc.CompraServiceServicer):
    def __init__(self):
        # Conectar a la base de datos de MongoDB
        self.client = MongoClient('mongodb://localhost:27017/')
        self.db = self.client['local']
        self.inventory = self.db['vehicles']  
        self.orders = self.db['orders']  

         # Conectar a RabbitMQ
        self.rabbitmq_connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
        self.rabbitmq_channel = self.rabbitmq_connection.channel()
        # Declarar la cola donde se enviarán las órdenes
        self.rabbitmq_channel.queue_declare(queue='order_queue', durable=True)
    
    def enviar_orden_a_rabbitmq(self, order_data):
        """
        Método para enviar la orden a RabbitMQ.
        """
        # Enviar la orden a RabbitMQ en formato JSON
        self.rabbitmq_channel.basic_publish(
            exchange='',
            routing_key='order_queue',  # Cola donde se enviarán las órdenes
            body=json.dumps(order_data),  # Convertir los datos de la orden a JSON
            properties=pika.BasicProperties(
                delivery_mode=2,  # Hacer el mensaje persistente
            )
        )
    

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
        insert_result = self.orders.insert_one(order_data)

        # Convertir el ObjectId a string para hacerlo serializable
        order_data['_id'] = str(insert_result.inserted_id)

        # Enviar la orden a RabbitMQ
        self.enviar_orden_a_rabbitmq(order_data)

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
