# Use the official Python image as base
FROM python:3.9-slim

# Set the working directory in the container
WORKDIR /app

# Copy the consumer source code into the container
COPY . .

# Install pika library
RUN pip install pika

# Wait rabbitmq to start
CMD [ "python", "-u", "consumer.py"]