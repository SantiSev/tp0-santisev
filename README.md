# TP0 | Parte 1: Introducción a Docker | Ejercicio 3

Esta documentación sirve como referencia sobre el funcionamiento del código y las decisiones tomadas para resolver los ejercicios.

# Decisiones Tomadas

El objetivo de este ejercicio es permitir que los cambios en los archivos de configuración sean efectivos sin necesidad de reconstruir las imágenes de Docker.

Este ejercicio se resolvio mediante un branch desde la rama ej2 y crear un script para verificar el correcto funcionamiento del servidor utilizando el comando netcat para interactuar con el mismo

Seccion Modificada:

```bash
#!/bin/bash
SERVER_PORT=$(grep 'SERVER_PORT' server/config.ini | cut -d'=' -f2)

PING_MESSAGE="hello darkness my old friend..."

SERVER_RESPONSE=$(echo "$PING_MESSAGE" | docker run --rm -i --network container:server busybox nc server $SERVER_PORT)

echo "Received response: '$SERVER_RESPONSE'"

RESULT="fail"
if [ "$SERVER_RESPONSE" = "$PING_MESSAGE" ]; then
    RESULT="success"
fi

echo "action: test_echo_server | result: $RESULT"
```

### Como funciona el script

## 1. Obtención del puerto del servidor

```bash
SERVER_PORT=$(grep 'SERVER_PORT' server/config.ini | cut -d'=' -f2)
```

- Busca la línea que contiene 'SERVER_PORT' en el archivo server/config.ini
- Extrae el valor después del signo '=' usando cut

> Por ejemplo, si el archivo contiene SERVER_PORT=8080, obtendrá 8080

## 2. Definición del mensaje de prueba
```bash
PING_MESSAGE="hello darkness my old friend..."
```

Define un mensaje específico que se enviará al servidor para la prueba

## 3. Envío del mensaje y captura de respuesta
   ```bash
   SERVER_RESPONSE=$(echo "$PING_MESSAGE" | docker run --rm -i --network container:server busybox nc server $SERVER_PORT)
   ```
   Esta es la parte más compleja:

- `echo "$PING_MESSAGE"` envía el mensaje a través de un pipe
- `docker run --rm -i --network container:server busybox` ejecuta un contenedor temporal de **busybox**, 

> **busybox** una imagen ligera de Linux que incluye utilidades básicas de Unix en un solo ejecutable. Esto permite ejecutar comandos como `nc` (netcat) sin necesidad de instalar herramientas adicionales en el sistema anfitrión.

- `--rm`: elimina el contenedor al finalizar
- `-i`: modo interactivo (permite recibir entrada por pipe)
- `--network container:server`: usa la misma red que el contenedor llamado "server"
- `nc server $SERVER_PORT:` usa **netcat** para conectarse al servidor en el puerto especificado

Finalmente, la respuesta del servidor se captura en SERVER_RESPONSE

## 4. Evaluación del resultado
   ```bash
   RESULT="fail"
   if [ "$SERVER_RESPONSE" = "$PING_MESSAGE" ]; then
   RESULT="success"
   fi
   ```

Por defecto asume que falló
Si la respuesta del servidor es exactamente igual al mensaje enviado, marca como exitoso
Esto confirma que el servidor está funcionando como un __"echo server"__ (devuelve lo que recibe)

## 5. Reporte final
   ```bash
   echo "action: test_echo_server | result: $RESULT"
   ```

Muestra el resultado final en un formato estructurado

# Como Ejectuar

El script se encuentra en la raíz del proyecto, primero se debe levantar el server con el commando:

```bash
make docker-compose-up
```

Luego, correr el script

```bash
./validar-echo-server.sh
```

