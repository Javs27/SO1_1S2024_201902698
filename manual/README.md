# Manual Técnico 
# Plataforma de Monitoreo y Señales a Procesos

**Introducción**

El presente documento describe los aspectos técnicos informáticos del sistema de información. El documento introducirá al personal técnico especializado encargado de las actividades de mantenimiento, revisión, solución de problemas, instalación y configuración del sistema.

**Requerimientos**

	 Software:
	 - Cualquier distribución derivada de debian

 
	 Hardware:
	 - Equipo con al menos 2 GB RAM   
	 - Equipo con al menos 4 GB disponible en
	   el disco duro
   

## **Instalación de aplicaciones**

Para implementar de manera correcta la plataforma de monitoreo en un sistema es necesario instalar ciertos programas, los cuales podrian necesitar de una conexión a internet de forma obligatoria.

| Software        | Descripción                                                                            | Comando de instalación en Ubuntu                         |
|-----------------|----------------------------------------------------------------------------------------|---------------------------------------------------------|
| Docker          | Plataforma de contenedores que permite empaquetar aplicaciones en contenedores         | `sudo apt-get install docker.io`                        |
| Docker Compose  | Herramienta para definir y ejecutar aplicaciones Docker multi-contenedor               | `sudo apt-get install docker-compose`                   |
| Gcc (GNU Compiler Collection) | Colección de compiladores para varios lenguajes de programación, incluyendo C, C++, Objective-C, Fortran, Ada, Go, y D | `sudo apt-get install build-essential`                  
|

## **Componentes**


**Backend**

Go como lenguaje principal, con los endpoints.

| Ruta                     | Descripción                                                | Detalles                                                   |
|--------------------------|------------------------------------------------------------|------------------------------------------------------------|
| `/ram`                   | Sirve datos en tiempo real de la RAM                       | Obtiene información actual sobre el uso de la RAM          |
| `/cpu`                   | Sirve datos en tiempo real de la CPU                       | Obtiene información actual sobre el uso de la CPU          |
| `/historical-ram`        | Sirve datos históricos de la RAM                           | Proporciona información sobre el uso de la RAM a lo largo del tiempo |
| `/historical-cpu`        | Sirve datos históricos de la CPU                           | Proporciona información sobre el uso de la CPU a lo largo del tiempo |
| `/processes/parents`     | Obtiene los procesos padres                                | Lista los procesos padres en el sistema                    |
| `/processes/details`     | Obtiene los procesos con detalles (hijos)                  | Lista los procesos junto con sus detalles, incluidos los procesos hijos |
| `/start`                 | Inicia la simulación de cambio de estado de un proceso     | Cambia el estado de un proceso a "iniciado"                |
| `/stop`                  | Detiene la simulación de cambio de estado de un proceso    | Cambia el estado de un proceso a "detenido"                |
| `/resume`                | Reanuda la simulación de cambio de estado de un proceso    | Cambia el estado de un proceso a "reanudado"               |
| `/kill`                  | Finaliza la simulación de cambio de estado de un proceso   | Cambia el estado de un proceso a "finalizado"              |
| `/getprocesses`          | Obtiene la lista de procesos                               | Lista todos los procesos actualmente en ejecución         |
| `/getprocesshistory`     | Obtiene el historial de procesos                           | Proporciona un historial de todos los procesos y sus cambios de estado |

**Base de datos**

  
Para la configuración de la base de datos, se utilizó MySQL alojado en un contenedor Docker, como se especifica en la configuración proporcionada. Esta configuración destaca el uso de la imagen `mysql:8.0` de Docker para instanciar un contenedor MySQL versión 8.0, la cual es una de las bases de datos de gestión relacional más populares y robustas, ampliamente utilizada para el almacenamiento y gestión de datos en aplicaciones web.

El contenedor se configura con varias opciones para asegurar su correcto funcionamiento y seguridad:

-   **Environment**: Se establecen variables de entorno dentro del contenedor. `MYSQL_ROOT_PASSWORD` se configura con el valor `admin` para definir la contraseña del usuario root de MySQL. `MYSQL_DATABASE` crea automáticamente una base de datos llamada `p1semi` al iniciar el contenedor, facilitando la separación y organización de los datos desde el inicio.
    
-   **Ports**: Se mapea el puerto `3306` del contenedor al puerto `3306` del host, permitiendo la conexión al servidor MySQL desde el host y otros contenedores o servicios que necesiten interactuar con la base de datos.
    
-   **Volumes**: Se utiliza un volumen llamado `db-data` para persistir los datos almacenados en `/var/lib/mysql` dentro del contenedor. Esto asegura que los datos no se pierdan cuando el contenedor se detiene o se elimina, proporcionando durabilidad y estabilidad a la base de datos.
    
-   **Healthcheck**: Define cómo verificar si el contenedor está funcionando correctamente. En este caso, se usa `mysqladmin ping -h localhost` para verificar la disponibilidad del servidor MySQL. El chequeo se realiza cada 10 segundos (`interval`), con un tiempo de espera (`timeout`) de 5 segundos. Si la prueba falla 5 veces consecutivas (`retries`), Docker considerará el contenedor como no saludable.

## Inicialización de Tablas
Para automatizar la creación de tablas necesarias para nuestra aplicación, incluimos un script SQL en el directorio mapeado a /docker-entrypoint-initdb.d dentro del contenedor de MySQL.

**Descripción de las Tablas de la Base de Datos** 

A continuación se presenta una descripción de cada tabla creada en la base de datos, incluyendo su propósito y los datos principales que almacena:

| Nombre de la Tabla    | Descripción                                                                                                                                                                                                                           |
|-----------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `datos_cpu`           | Almacena registros de uso de CPU, incluyendo el total y el porcentaje de uso, junto con la fecha de registro.                                                                                                                        |
| `datos_ram`           | Guarda información sobre el uso de la RAM, incluyendo total, en uso, libre, porcentaje en uso y la fecha de registro.                                                                                                               |
| `processes`           | Contiene datos sobre los procesos ejecutándose, como el ID del proceso, nombre, ID del usuario, estado y uso de RAM.                                                                                                                 |
| `process_log`         | Registra eventos o estados de procesos específicos, incluyendo un ID de mensaje, ID del proceso, estado, mensaje descriptivo y la fecha del registro.                                                                                |
| `process_relations`   | Define las relaciones entre procesos, especificando relaciones padre-hijo entre procesos a través de sus ID. Incluye restricciones de clave foránea para asegurar la integridad referencial con la tabla `processes`.                 |
| `process_summary`     | Proporciona un resumen de los estados de los procesos en un momento dado, incluyendo conteos de procesos en ejecución, durmiendo, zombis, detenidos y el total, junto con la fecha de registro.                                      |

Estas tablas son fundamentales para el almacenamiento y manejo de datos dentro de la aplicación, permitiendo un seguimiento detallado del rendimiento del sistema y la gestión de procesos.




## Frontend

En el proyecto, se empleó React para el desarrollo del frontend, creando una interfaz dinámica y moderna. La configuración del entorno se basa en Docker, utilizando Node.js para el entorno de ejecución de React y Nginx como servidor web en producción, optimizado para eficiencia con la versión `nginx:alpine`. La configuración de Nginx se personaliza a través de volúmenes, y se configura para servir la aplicación en el puerto 80, asegurando una integración fluida entre el frontend, el backend, y servicios dependientes. Esta estructura ofrece una solución robusta y escalable para despliegues modernos.

## Módulos
  
En el desarrollo del proyecto, se crearon módulos de kernel específicos para mejorar la monitorización y gestión de recursos del sistema. Estos módulos se centran en dos áreas críticas: la **RAM y la CPU.**

-   **Módulo de RAM**: Este módulo, cargado mediante el comando `sudo insmod ram_so1_1s2024.ko`, permite a los usuarios ver el porcentaje de uso de la memoria RAM en tiempo real. La funcionalidad de este módulo es esencial para optimizar el rendimiento del sistema, permitiendo a los desarrolladores y administradores de sistemas monitorear el consumo de memoria y tomar decisiones informadas sobre la gestión de recursos.
    
-   **Módulo de CPU**: A través del comando `sudo insmod cpu_so1_1s2024.ko`, se carga un módulo diseñado para monitorear el uso de la CPU. Este módulo no solo proporciona información sobre el porcentaje de uso de la CPU, sino que también entrega detalles sobre los procesos en ejecución, incluyendo identificadores de proceso (PIDs) y otros datos relevantes. Esta información es vital para la administración eficiente del procesador, permitiendo la identificación de procesos que consumen una cantidad desproporcionada de recursos y la toma de medidas para optimizar el rendimiento general del sistema.