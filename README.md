# TP0 | Parte 1: Introducción a Docker | Ejercicio 1

Esta documentación sirve como referencia sobre el funcionamiento del código y las decisiones tomadas para resolver los ejercicios.

# Decisiones Tomadas
El objetivo de este ejercicio generar-compose.sh que permita crear una definición de Docker Compose con una cantidad configurable de clientes.

El diseño del script se basó en el archivo **docker-compose-dev.yaml** incluido originalmente en el proyecto base al momento de realizar el fork.

Para el diseño de este script lo vamos a explicar en 4 partes:

### 1. Inicialización y Parámetros Requeridos

```bash
if [ "$#" -ne 2 ]; then
   echo "Usage: $0 <yaml_file> <amount_clients>"
   exit 1
fi

YAML_FILE="$1"
AMOUNT_CLIENTS="$2"
```
El script requiere exactamente 2 parámetros:

- **yaml_file:** Archivo YAML de entrada (aunque no se utiliza en el código actual)
- **amount_clients:** Número de clientes que se desean crear

Si no se proporcionan exactamente 2 argumentos, el script muestra el mensaje de uso y termina con código de error 1.

### 2. Creación del Servicio Servidor
```bash
cat > docker-compose-dev.yaml << EOF
name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
    networks:
      - testing_net
    volumes:
      - ./server/config.ini:/config.ini
EOF
```
Esta sección crea la estructura base del archivo Docker Compose con:

- **Nombre del proyecto:** tp0
- **Nombre del contenedor:** server
- **Imagen:** server:latest
- **Punto de entrada:** python3 /main.py
- **Variable de entorno** para mostrar output de Python sin buffer
- **Conectado a la red testing_net**
- **Monta el archivo de configuración:** ./server/config.ini


### 3. Creación de Clientes Dinámicos
```bash
bashfor i in $(seq 1 "$AMOUNT_CLIENTS"); do
    cat >> docker-compose-dev.yaml << EOF

  client$i:
    container_name: client$i
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=$i
    networks:
      - testing_net
    volumes:
      - ./client/config.yaml:/config.yaml
    depends_on:
      - server
EOF
done
```
La creacion de los clientes consiste en un bucle for que genera dinámicamente los servicios cliente:

- **Cantidad:** Según el parámetro $AMOUNT_CLIENTS
- **Nomenclatura:** __client1, client2, client3, etc.__

Configuración por cliente:

- **Imagen:** client:latest
- **ID único mediante variable de entorno:** CLI_ID
- **Conectado a la misma red que el servidor**
- **Monta configuración desde:** ./client/config.yaml
- **Dependencia:** Todos los clientes dependen del servidor (se inicia después)

### 4. Creacion de servicio red

Finalmente, se define la red testing_net con una subred específica 172.25.125.0/24.

> Esta sección fue extraída directamente del archivo **docker-compose-dev.yaml** incluido en el proyecto base.


# Como Ejectuar

El script se encuentra en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

Ejemplo de uso:
```bash
./generar-compose.sh [DOCKER-COMPOSE.YAML] [CANT_CLIENTES]
```

Ejemplo de uso:
```bash
./generar-compose.sh docker-compose-dev.yaml 5
```
