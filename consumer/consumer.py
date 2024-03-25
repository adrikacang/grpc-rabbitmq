import pika
import time

# Callback function to process incoming messages
def callback(ch, method, properties, body):
    print(" [x] Received %r" % body)

# Connect to RabbitMQ server
while True:
    try:
        credentials = pika.PlainCredentials('user', 'password')
        connection = pika.BlockingConnection(pika.ConnectionParameters('rabbitmq', credentials=credentials))
        channel = connection.channel()

        # Declare a queue
        channel.queue_declare(queue='bookings')

        # Consume messages from the queue
        channel.basic_consume(queue='bookings',
                            auto_ack=True,
                            on_message_callback=callback)
        break
    except Exception as err:
        time.sleep(1)

print(' [*] Waiting for messages.')
channel.start_consuming()