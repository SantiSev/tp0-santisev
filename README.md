# TP0 | Parte 2: Repaso de Comunicaciones | Ejercicio 7

Esta documentación sirve como referencia sobre el funcionamiento del código y las decisiones tomadas para resolver los ejercicios.

# Decisiones Tomadas

Para esta seccion, se modifca el cliente y el server para poder procesar la loteria y anunciar los ganadores
Una vez que todas las apuestas estan almacenadas exitosamente de todas las agencias, se notifica a cada agencia su ganador y el server anuncia los ganadores de todas las loterias

> Todo el código está redactado en inglés, con excepción de algunos logs específicos que permanecen en español para garantizar la compatibilidad con los tests proporcionados.

# Cambios en script de `generar-compose.sh`

Se modifica la sección del servidor para agregar la cantidad de agencias que debe procesar antes de anunciar los resultados de la lotería:

```bash
cat > "$YAML_FILE" << EOF
name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - AGENCIES_AMOUNT=$AMOUNT_CLIENTS
    networks:
      - testing_net
    volumes:
      - ./server/config.ini:/config.ini
EOF
```

La cantidad de agencias es igual a la cantidad de clientes que se crearán, haciendo el sistema dinámico para loterias con cualquier número de clientes.

# Cambios en la Arquitectura del Servidor

## Server

Se modifica la lógica del servidor para procesar múltiples clientes secuencialmente y anunciar ganadores una vez completadas todas las agencias.

### Cambios Implementados

- **Contador de agencias**: Utiliza `processed_agencies` para rastrear la cantidad de clientes procesados
- **Control de límite**: Lee la variable de entorno `AGENCIES_AMOUNT` para determinar cuántos clientes esperar
- **Procesamiento secuencial**: Acepta y procesa un cliente a la vez hasta alcanzar el límite configurado
- **Anuncio de ganadores**: Una vez procesadas todas las agencias, ejecuta el sorteo y anuncia resultados

### Flujo del Servidor

```
1. Inicialización → Server(agencies_amount=N)
2. Bucle principal:
   a. Aceptar conexión → ConnectionManager.accept_connection()
   b. Incrementar contador → processed_agencies += 1
   c. Crear sesión → ClientManager.add_client(connection)
   d. Procesar cliente → ClientSession.begin()
   e. Verificar condición → processed_agencies < agencies_amount
3. Condición de salida → processed_agencies == agencies_amount
4. Anuncio de resultados → _tally_results()
5. Shutdown → _shutdown()
```

## Protocol

Se modifica el protocolo de comunicación para incluir el envío de resultados de ganadores desde el servidor hacia cada cliente.

### Cambios Implementados

- **Nueva función `send_winners()`**: Envía un mensaje al cliente anunciando los ganadores específicos de su agencia
- **Protocolo bidireccional**: Ahora el servidor puede iniciar comunicación con el cliente para notificar resultados
- **Mensaje de ganadores**: Incluye información específica de los ganadores de la agencia correspondiente
- **Comunicación asíncrona**: El envío de ganadores ocurre después de procesar todas las agencias

### Estructura del Protocolo de Ganadores

**Protocolo de envío de ganadores:**

```bash
[WINNER_HEADER][DATA_LENGTH][WINNER_DATA]
```

- `WINNER_HEADER`: Identificador del mensaje de ganadores
- `DATA_LENGTH`: Longitud de los datos de ganadores
- `WINNER_DATA`: Información de los ganadores de la agencia (DNI, número ganador, etc.)

### Flujo del Protocolo Actualizado

1. Cliente procesa y envía todas las apuestas → Protocolo normal de batches
2. Cliente envía DONE → Finalización de envío de apuestas
3. Servidor espera a que todas las agencias terminen
4. Servidor ejecuta sorteo → lottery_service.announce_winners()
5. Servidor envía ganadores individuales → send_winners() por cada cliente
6. Cliente recibe y procesa ganadores de su agencia → Logging de resultados

## Session

Se modifica el manejo de sesiones para procesar múltiples lotes de apuestas y enviar resultados de ganadores a cada cliente al finalizar.

### Cambios Implementados

- **Envío de resultados**: Nueva función `send_results()` para notificar ganadores específicos de cada agencia. La sesión utiliza `AgencyHandler` para enviar los resultados de la lotería a cada cliente

# Cambios en la Arquitectura del Cliente

## Business

Se modifica el servicio `AgencyService` para parsear los resultados de la loteria agregnado la funcion nueva de `ShowResults()` que consiste en parsear el string obtenido del `AgencyHandler` e imprimir los ganadores de la loteria y cuantos ganadores hubo


## Client

Se modifica el cliente para mantener la conexión activa hasta recibir los resultados de la lotería desde el servidor.

### Flujo del Cliente

1. Inicialización → AgencyService.NewAgencyService(filePath, batchSize)
2. Conexión → ConnectionManager.Connect(serverAddress)
3. Bucle de procesamiento:
   - Leer batch → AgencyService.ReadBets(maxBatchAmount)
   - Enviar al servidor → AgencyHandler.SendBets(batch)
   - Recibir confirmación → AgencyHandler.RecvConfirmation()
   - Verificar datos restantes → AgencyService.HasData()
4. Señal de finalización → AgencyHandler.SendDone()
5. Esperar los resultados de la loteria → AgencyHandler.GetResults()
5. Cierre de recursos → AgencyService.Close() + ConnectionInterface.Close()

## Protocol

El protocolo mantiene la misma estructura base para el envío de batches de apuestas, pero incorpora una nueva funcionalidad para la recepción de resultados de la lotería.

### Nueva Funcionalidad: GetResults

**Función `GetResults()`**: Permite al cliente recibir los ganadores específicos de su agencia una vez finalizado el sorteo.

#### Protocolo de Recepción de Ganadores

**Estructura del mensaje:**
```
[WINNERS_HEADER][WINNER_COUNT][WINNER_DATA]
```

- **`WINNERS_HEADER`**: Identificador del mensaje de ganadores (1 byte)
- **`WINNER_COUNT`**: Longitud de los datos de ganadores (2 bytes, BigEndian)
- **`WINNER_DATA`**: Información de los ganadores en formato string

#### Características Principales

- **Espera bloqueante**: El cliente permanece en `ReceiveData()` hasta recibir el mensaje del servidor
- **Validación de header**: Verifica que el mensaje sea del tipo `WINNERS_HEADER` antes de procesar
- **Manejo de longitud variable**: Utiliza `WINNER_COUNT` para determinar cuántos bytes leer

#### Integración en el Flujo del Cliente

```
1. Cliente envía todas las apuestas (batches)
2. Cliente envía DONE → Finalización
3. Cliente ejecuta GetResults() → Espera bloqueante
4. Servidor procesa todas las agencias
5. Servidor envía ganadores → WINNERS_HEADER + datos
6. Cliente recibe y procesa resultados
```

# Cómo Ejecutar

1. generar un archivo .yaml de docker-compose mediante la funcion

```bash
./generar-compose.sh docker-compose-dev.yaml 2
```

2. **Limpieza inicial**: Ejecutar `make docker-compose-down` para asegurar un inicio limpio
3. **Inicio de contenedores**: Ejecutar `make docker-compose-up` para iniciar los contenedores de servidor y cliente
4. **Visualización de logs**: Ejecutar `make docker-compose-logs` para ver los resultados y outputs del servidor y clientes
5. **Verificación de estado**: Ejecutar `docker ps -a` para confirmar que los contenedores finalizaron con exit status 0

## Script de Automatización

> **Alternativa conveniente:** Se incluye el script `run_local_test.sh` que automatiza los primeros 3 comandos y genera un archivo `logs.txt` con el output de `make docker-compose-logs` para visualización offline.

### Uso del script:

```bash
./run_local_test.sh
```

Este script ejecuta automáticamente toda la secuencia de testing y guarda los logs en un archivo para análisis posterior.
