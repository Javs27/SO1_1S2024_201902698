# Usa la imagen de Golang como base
FROM golang:latest

# Establece el directorio de trabajo en /app
WORKDIR /app

# Copia el código fuente al contenedor
COPY . .

# Descarga las dependencias del proyecto
RUN go get -d -v ./...

# Compila el código Go dentro del contenedor
RUN go build -o main .

# Expone el puerto 8080 en el contenedor
EXPOSE 8080

# Comando para ejecutar la aplicación al iniciar el contenedor
CMD ["./main"]
