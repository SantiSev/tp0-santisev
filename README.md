# TP0 | Parte 2: Repaso de Comunicaciones | Ejercicio 6

Esta documentación sirve como referencia sobre el funcionamiento del código y las decisiones tomadas para resolver los ejercicios.

# Decisiones Tomadas

Para esta seccion, se modifca el cliente y el server para poder transmitir desde el client al server multiples bets a la vez mediante la lectura de archivos .csv y el envio de chunks desde el client hacia el server

Una vez que todas las apuestas estan almacenadas exitosamente, envia un mensaje de confirmacion al cliente

> Todo el código está redactado en inglés, con excepción de algunos logs específicos que permanecen en español para garantizar la compatibilidad con los tests proporcionados.

# Cambios en script de `generar-compose.sh`

Para soportar el procesamiento de múltiples apuestas desde archivos CSV, se modificó el script de generación de `docker-compose` para asignar a cada cliente un archivo de agencia específico.

```bash
for i in $(seq 1 "$AMOUNT_CLIENTS"); do
    cat >> docker-compose-dev.yaml << EOF

  client$i:
    container_name: client$i
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=$i
      - CLI_AGENCY_FILEPATH=/data/agency.csv  <-- [ ruta de agency csv filepath ]
    networks:
      - testing_net
    volumes:
      - ./client/config.yaml:/config.yaml
      - ./.data/agency-$i.csv:/data/agency.csv  <-- [ Mapeo específico por cliente ]
    depends_on:
      - server
EOF
done
```

También se agrega la variable de entorno `CLI_AGENCY_FILEPATH=/data/agency.csv` que especifica la ruta donde el contenedor debe buscar el archivo CSV de la agencia. **Importante:** Esta ruta debe coincidir exactamente con el punto de montaje del volumen, de lo contrario la aplicación fallará al no encontrar el archivo de datos.

El mapeo de volúmenes `./.data/agency-$i.csv:/data/agency.csv` permite que cada cliente acceda a su archivo de agencia específico (agency-1.csv, agency-2.csv, etc.) mientras mantiene una ruta estándar (`/data/agency.csv`) en el código. Esto proporciona aislamiento de datos entre clientes.

- **Limitaciones:** Si agregas mas clientes que archivos disponibles de agency, este script falla, lo mismo si los arvhivos csv no siguien la nomenclatura agency-{i}.csv

# Cambios en la Arquitectura del Servidor

## Business

## Config

## Protocol

## Server

## Session

# Cambios en la Arquitectura del Cliente

## Business

Se modifica el servicio `AgencyService` para procesar múltiples apuestas desde archivos CSV:

**Funcionalidad principal:**

- **Inicialización**: Recibe la ruta del archivo CSV de la agencia durante la construcción del servicio
- **Lectura de datos**: Utiliza `bufio.Scanner` de la librería estándar de Go para leer línea por línea el archivo de apuestas
- **Gestión de recursos**: Implementa el método `Close()` para liberar correctamente los recursos del scanner al finalizar

**Métodos implementados:**

- `NewAgencyService(filePath string)`: Constructor que inicializa el servicio con la ruta del archivo CSV
- `ReadBets()`: Lee y valida todas las apuestas del archivo, retornando una lista de apuestas estructuradas
- `Close()`: Cierra el scanner y libera recursos asociados

**Decisión de diseño:** Se optó por utilizar `bufio.Scanner` para manejar archivos de gran tamaño de forma eficiente, procesando línea por línea sin cargar todo el archivo en memoria simultáneamente.

## Config

Ahora en lugar de procesar una apuesta individual desde variables de entorno, se configura la instancia de `ClientConfig` con un nuevo parámetro llamado `AgencyFilePath`. Esta ruta de archivo se utiliza posteriormente en `AgencyService` para leer múltiples apuestas desde archivos CSV.

## Client

Se modifica la funcionalidad principal del cliente para procesar archivos CSV de apuestas en lotes y enviar una señal de finalización al servidor.

### Cambios Implementados

- **Integración con AgencyService**: El cliente ahora utiliza `AgencyService` para leer apuestas desde archivos CSV en lugar de variables de entorno
- **Procesamiento por batches**: Lee y envía múltiples apuestas por transacción (batches configurables)

### Flujo del Cliente

1. Inicialización → AgencyService.NewAgencyService(filePath, batchSize)
2. Conexión → ConnectionManager.Connect(serverAddress)
3. Bucle de procesamiento:
   - Leer batch → AgencyService.ReadBets(maxBatchAmount)
   - Enviar al servidor → AgencyHandler.SendBets(batch)
   - Recibir confirmación → AgencyHandler.RecvConfirmation()
   - Verificar datos restantes → AgencyService.HasData()
4. Señal de finalización → AgencyHandler.SendDone()
5. Cierre de recursos → AgencyService.Close() + ConnectionInterface.Close()

## Protocol

El protocolo mantiene la misma estructura base, pero incorpora mejoras para el procesamiento de múltiples apuestas:

**Cambios principales:**

- **Envío masivo**: Se transmiten lotes de apuestas (batches) en lugar de apuestas individuales
- **Confirmación simplificada**: La función `RecvConfirmation()` ahora recibe únicamente un header de estado en lugar de información detallada de cada batch enviada

### Estructura del Protocolo

**Protocolo de envío de batches:**

```bash
[HEADER_BYTE] [LENGTH_BYTE] [BET_DATA]
```

**Protocolo de recepción de batches:**

```bash
[SUCCESS_HEADER]
```

puede ser tando un valor de SUCCESS como uno de FAIL

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
