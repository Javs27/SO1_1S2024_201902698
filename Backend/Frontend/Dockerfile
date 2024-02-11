# Establecer la imagen base
FROM node:14-alpine AS build

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el archivo package.json y el archivo package-lock.json (si existe)
COPY package*.json ./

# Instalar las dependencias
RUN npm install

# Copiar los archivos del frontend
COPY . .

# Compilar la aplicación de React
RUN npm run build

# ----------------------

# Establecer una nueva imagen base para el servidor web
FROM nginx:alpine

# Copiar los archivos de construcción del frontend desde la etapa de compilación
COPY --from=build /app/build /usr/share/nginx/html

# Exponer el puerto 80
EXPOSE 80

# Comando para iniciar el servidor web
CMD ["nginx", "-g", "daemon off;"]
