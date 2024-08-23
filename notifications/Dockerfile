FROM golang:1.21

# Copy the .env file from the host machine to the /app directory inside the container

# Copy your application binary (assuming it's named "app") to the /app directory
COPY ./bin/app /app/app

ENV IS_DOCKER=1
# Set the working directory to /app
WORKDIR /app

# Command to run your application
CMD ["./app"]
