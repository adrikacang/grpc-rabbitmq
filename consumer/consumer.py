import pika

# Callback function to process incoming messages
def callback(ch, method, properties, body):
    print(" [x] Received %r" % body)

# Connect to RabbitMQ server
credentials = pika.PlainCredentials('user', 'password')
connection = pika.BlockingConnection(pika.ConnectionParameters('localhost', credentials=credentials))
channel = connection.channel()

# Declare a queue
channel.queue_declare(queue='bookings')

# Consume messages from the queue
channel.basic_consume(queue='bookings',
                      auto_ack=True,
                      on_message_callback=callback)

print(' [*] Waiting for messages. To exit press CTRL+C')
channel.start_consuming()