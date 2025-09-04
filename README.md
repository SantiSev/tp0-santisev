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

El mapeo de volúmenes `./.data/agency-$i.csv:/data/agency.csv` permite que cada cliente acceda a su archivo de agencia específico (agency-1.csv, agency-2.csv, etc.) mientras mantiene una ruta estándar (`/data/agency.csv`) en el código.  Esto proporciona aislamiento de datos entre clientes.

- **Limitaciones:** Si agregas mas clientes que archivos disponibles de agency, este script falla, lo mismo si los arvhivos csv no siguien la nomenclatura agency-{i}.csv




# Cambios en la Arquitectura del Servidor

## Business

Se agrega un nuevo método al servicio `LotteryService` para consultar apuestas procesadas:

`def get_bets_by_agency(self, agency_id: int) -> list[Bet]:` que lee el archivo `bets.csv` y retorna todas las apuestas correspondientes a una agencia específica, permitiendo contar cuántas apuestas fueron procesadas por cada agencia.

**Decisión de diseño:** Se optó por almacenar todas las apuestas en archivo y realizar consultas directas desde el sistema de archivos para evitar mantener grandes volúmenes de datos en memoria, optimizando el uso de recursos del servidor.

## Config

Ahora en lugar de procesar una apuesta individual desde variables de entorno, se configura la instancia de `ClientConfig` con un nuevo parámetro llamado `AgencyFilePath`. Esta ruta de archivo se utiliza posteriormente en `AgencyService` para leer múltiples apuestas desde archivos CSV.

## Protocol

Define el protocolo de comunicación específico de la aplicación. Especifica el formato de mensajes, serialización/deserialización y las reglas de intercambio de datos entre cliente y servidor.

Este módulo contiene 2 clases fundamentales del servidor:

- **AgencyHandler**: Gestiona el protocolo de comunicación con clientes utilizando `ConnectionInterface`. Se encarga de:

  - Procesar headers de mensajes entrantes
  - Coordinar la recepción de apuestas mediante `BetParser`
  - Enviar confirmaciones de éxito/fallo al cliente
  - Manejar errores de comunicación durante el intercambio de datos

- **BetParser**: Responsable del procesamiento y transformación de datos. Funciones principales:
  - Parsear datos CSV recibidos como strings
  - Convertir información cruda en objetos `Bet` estructurados
  - Validar integridad de los datos antes de la conversión
  - Manejar casos de error en el parsing de apuestas

### Estructura del Protocolo

**Mensaje de Envío (Cliente → Servidor):**

```bash
[BET_HEADER][DATA_LENGTH][BET_DATA]
```

**Mensaje de Confirmación (Servidor → Cliente):**

```bash
[STATUS_HEADER][RESPONSE_LENGTH][CONFIRMATION_DATA]
```

## Server

Implementa la funcionalidad del servidor, incluyendo el manejo de conexiones entrantes, procesamiento de requests y la lógica específica del lado servidor.

Este módulo contiene 2 clases fundamentales:

- **ServerConfig**: Clase de configuración que encapsula todos los parámetros necesarios para la inicialización del servidor. Almacena información como el puerto de escucha, el número máximo de conexiones pendientes (backlog) y el nivel de logging. Actúa como un objeto de transferencia de datos que centraliza la configuración del sistema, eliminando la necesidad de pasar múltiples parámetros individuales entre componentes.

- **Server**: Clase principal que actúa como el núcleo orquestador del sistema servidor. Sus responsabilidades incluyen:

  **Inicialización del Sistema:**

  - Configura el `ConnectionManager` para gestionar conexiones TCP entrantes
  - Instancia `LotteryService` para manejar la lógica de negocio de apuestas
  - Establece `ClientManager` para administrar sesiones activas de clientes
  - Registra manejadores de señales (`SIGTERM`, `SIGINT`) para garantizar terminación controlada

  > **Nota:** El manejo de `SIGINT` no es requerido por el enunciado, pero se agregó para facilitar el debugging local y verificar el correcto funcionamiento del shutdown graceful durante el desarrollo.

  **Ciclo de Vida del Servidor:**

  1. **Inicio**: Activa la escucha en el puerto configurado
  2. **Aceptación**: Permanece bloqueado esperando conexiones entrantes
  3. **Sesión**: Crea una instancia `ClientSession` para cada cliente conectado
  4. **Delegación**: Transfiere el control al cliente para procesamiento de apuestas
  5. **Finalización**: Ejecuta shutdown graceful liberando todos los recursos

  **Limitación Actual:** Según los requisitos de este ejercicio, la implementación procesa únicamente un cliente de forma secuencial. Esta arquitectura será escalada en ejercicios posteriores.

## Session

Gestiona las sesiones de usuario o conexión. Mantiene el estado de las interacciones, autenticación y el contexto de cada cliente conectado.

Este módulo contiene 2 clases fundamentales:

- **ClientSession**: Representa una sesión individual de cliente y maneja todo el ciclo de vida de la comunicación con una agencia. Sus responsabilidades incluyen:

  **Inicialización:**

  - Almacena la referencia de conexión (`ConnectionInterface`) para comunicación directa
  - Mantiene un ID único de agencia para identificación
  - Configura el `AgencyHandler` para manejar el protocolo de comunicación
  - Establece la referencia al `LotteryService` para procesamiento de apuestas

  **Procesamiento Principal (`begin`):**

  1. **Recepción**: Utiliza `AgencyHandler` para recibir apuestas del cliente
  2. **Almacenamiento**: Delega al `LotteryService` para persistir las apuestas
  3. **Confirmación**: Envía confirmación de éxito al cliente
  4. **Manejo de errores**: Captura excepciones y envía confirmación de fallo

  **Finalización (`finish`):**

  - Cierra la conexión de red de forma controlada
  - Libera recursos asociados a la sesión

- **ClientManager**: Actúa como un registro centralizado y coordinador de todas las sesiones activas. Sus funciones principales son:

  **Gestión de Sesiones:**

  - Mantiene una lista de todas las sesiones de cliente activas
  - Asigna IDs únicos secuenciales a cada nueva sesión
  - Proporciona una interfaz unificada para administrar múltiples clientes

  **Ciclo de Vida de Clientes:**

  - **`add_client()`**: Crea nuevas instancias de `ClientSession` para conexiones entrantes
  - **`remove_client()`**: Finaliza sesiones específicas y las elimina del registro
  - **`shutdown()`**: Termina todas las sesiones activas durante el cierre del servidor

  **Coordinación:**

  - Comparte la misma instancia de `LotteryService` entre todos los clientes
  - Garantiza que cada cliente tenga acceso a la lógica de negocio común
  - Facilita la gestión centralizada de recursos

  **Limitación Actual:** Para la escala de este ejercicio, la implementación de `ClientManager` no era estrictamente necesaria, sin embargo, proporciona una base sólida y escalable para resolver ejercicios posteriores que requerirán el manejo de múltiples clientes.

## Utils

Contiene utilidades auxiliares y funciones helper proporcionadas por la catedra para poder leer / escribir el archivo donde se almacenan los **Bets** (apuestas).

Contiene únicamente el archivo `utils.py` proporcionado por la cátedra, el cual no puede ser modificado según las especificaciones del enunciado.

# Cambios en la Arquitectura del Cliente

## Business

Contiene la lógica de negocio específica del cliente. Maneja las reglas y procesos relacionados con la generación, validación y preparación de apuestas para su envío al servidor.

Este módulo contiene el archivo `agency_service.go` que implementa la clase `AgencyService` con las siguientes responsabilidades:

- **Validación de apuestas**: Verifica que el formato de las apuestas sea correcto (6 campos separados por comas)
- **Lectura de datos**: Proporciona una interfaz para obtener las apuestas validadas
- **Gestión de agencia**: Mantiene el ID único de la agencia para identificación

**Limitación Actual:** La implementación actual procesa únicamente apuestas individuales obtenidas desde la configuración. En ejercicios posteriores, esta arquitectura se expandirá para leer múltiples apuestas desde archivos de agencias, manteniendo la misma estructura modular.

## Client

## Config

Administra la configuración del cliente, incluyendo la lectura del archivo `config.yaml`, parámetros de conexión y inicialización del sistema de logging.

Este módulo contiene el archivo `config.go` con dos funciones principales:

**`InitConfig()`**: Inicializa la configuración del cliente mediante la lectura de archivos de config y variables de entorno

> **Nota:** Para el alcance de este ejercicio, los atributos de la apuesta enviada al servidor se almacenan en variables de entorno. En ejercicios posteriores, esta implementación será reemplazada por la lectura de archivos de apuestas para mayor escalabilidad.

**`InitLogger()`**: Configura el sistema de logging con niveles de verbosidad (INFO / DEBUG)

## Network

Proporciona las abstracciones de red para el lado cliente. Implementa las funcionalidades de conexión TCP, envío y recepción de datos, y manejo de la comunicación de bajo nivel con el servidor.

Este módulo contiene 2 clases fundamentales:

- **ConnectionManager**: Gestiona la establecimiento de conexiones TCP hacia el servidor con retry automático (máximo 3 intentos con intervalos de 100ms). Cuando logra conectarse al server, devuelve una instancia de `ConnectionInterface`

- **ConnectionInterface**: Abstrae las operaciones de socket TCP proporcionando métodos `Connect()`, `SendData()`, `ReceiveData()` y `Close()` para comunicación confiable con el servidor.

### Manejo de Short Read/Write

- `SendData()` utiliza un loop que continúa escribiendo hasta enviar todos los bytes

- **`ReceiveData()`**: Utiliza la función estándar de Go `io.ReadFull()` que garantiza lectura completa del buffer

## Protocol

Define e implementa el protocolo de comunicación desde la perspectiva del cliente. Maneja la serialización de datos de apuestas, el formato de mensajes enviados al servidor y el procesamiento de las respuestas de confirmación recibidas.

Este módulo contiene la clase `AgencyHandler` que gestiona el intercambio de mensajes con el servidor:

### AgencyHandler

Implementa el protocolo de comunicación cliente-servidor para el envío de apuestas y recepción de confirmaciones. Utiliza la instancia de ConnectionInterface para manejar envio y recepcion de datos

**Métodos principales:**

#### `SendBets(bet string, connSock *ConnectionInterface)`

Envía apuestas al servidor siguiendo el protocolo definido:

1. **Header**: Envía el byte identificador del tipo de mensaje (`HEADER`)
2. **Longitud**: Envía un byte indicando la longitud de los datos de la apuesta
3. **Datos**: Envía los datos de la apuesta en formato string

**Protocolo de envío:**

```bash
[HEADER_BYTE] [LENGTH_BYTE] [BET_DATA]
```

#### `RecvConfirmation(connSock *ConnectionInterface)`

Recibe y procesa la confirmación del servidor:

1. **Verificación de header**: Lee el header de respuesta (`SUCCESS_HEADER`)
2. **Longitud del mensaje**: Obtiene la longitud de la respuesta
3. **Datos de confirmación**: Lee los datos de confirmación (DNI y número de apuesta)
4. **Logging**: Registra el resultado de la operación

**Protocolo de recepción:**

```
[SUCCESS_HEADER][RESPONSE_LENGTH][CONFIRMATION_DATA]
```

# Cómo Ejecutar

1. **Limpieza inicial**: Ejecutar `make docker-compose-down` para asegurar un inicio limpio
2. **Inicio de contenedores**: Ejecutar `make docker-compose-up` para iniciar los contenedores de servidor y cliente
3. **Visualización de logs**: Ejecutar `make docker-compose-logs` para ver los resultados y outputs del servidor y clientes
4. **Verificación de estado**: Ejecutar `docker ps -a` para confirmar que los contenedores finalizaron con exit status 0

## Script de Automatización

> **Alternativa conveniente:** Se incluye el script `run_local_test.sh` que automatiza los primeros 3 comandos y genera un archivo `logs.txt` con el output de `make docker-compose-logs` para visualización offline.

### Uso del script:

```bash
./run_local_test.sh
```

Este script ejecuta automáticamente toda la secuencia de testing y guarda los logs en un archivo para análisis posterior.
