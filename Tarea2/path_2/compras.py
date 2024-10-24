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

        ###################### Conexión MongoDB ######################
        self.client = MongoClient('mongodb://localhost:27017/')
        self.db = self.client['local']
        self.inventory = self.db['vehicles']  
        self.orders = self.db['orders']  

        ###################### Conexión RabbitMQ ######################
        self.rabbitmq_connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
        self.rabbitmq_channel = self.rabbitmq_connection.channel()
        self.rabbitmq_channel.queue_declare(queue='order_queue', durable=True)
    

    def enviar_orden_a_rabbitmq(self, order_data):
        ############ Cola facturacion ############
        self.rabbitmq_channel.basic_publish(
            exchange='',
            routing_key='cola_facturacion',
            body=json.dumps(order_data),
            properties=pika.BasicProperties(
                delivery_mode=2,  #Persistente
            )
        )

        ############ Cola inventario ############
        self.rabbitmq_channel.basic_publish(
            exchange='',
            routing_key='cola_inventario',
            body=json.dumps(order_data),
            properties=pika.BasicProperties(
                delivery_mode=2,
            )
        )

    def RealizarCompra(self, request, context):
        customer = request.cliente
        vehicles = request.vehiculos
        
        vehiculos = []
        precios = []

        # Obtener el precio
        for auto in vehicles:
            resultado = self.inventory.find_one({"manufacturer": auto.manufacturer, "model": auto.model, "year": auto.year})
            if resultado:
                precio = resultado.get('price', 0)
                precios.append(precio) 
                vehiculos.append(f"{auto.manufacturer} {auto.model} ({auto.year}) - Precio: {precio}")
            else:
                precios.append(0)  # Si no se encuentra agrega 0
                vehiculos.append(f"{auto.manufacturer} {auto.model} ({auto.year}) - Precio no disponible")

        mensaje_respuesta = f"Compra realizada por {customer.name} {customer.lastname} para los vehículos: {', '.join(vehiculos)}"

        # Crear orden en BD
        order_data = {
            "order_date": datetime.now().strftime("%Y/%m/%d"),
            "vehicles": [
                {
                    "manufacturer": v.manufacturer,
                    "model": v.model,
                    "year": v.year,
                    "price": precios[i]
                } for i, v in enumerate(vehicles)
            ],
            "customer": {
                "name": customer.name,
                "lastname": customer.lastname,
                "email": customer.email,
                "phone": customer.phone
            }
        }

        agregar_orden = self.orders.insert_one(order_data)
        print()
        order_data['_id'] = str(agregar_orden.inserted_id)
        self.enviar_orden_a_rabbitmq(order_data)
        return orden_pb2.CompraResponse(message=mensaje_respuesta)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    orden_pb2_grpc.add_CompraServiceServicer_to_server(CompraService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Servidor gRPC iniciado en el puerto 50051")
    server.wait_for_termination()

if __name__ == "__main__":
    serve()
