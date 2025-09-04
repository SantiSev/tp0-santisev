# TP0 | Parte 2: Repaso de Comunicaciones | Ejercicio 5

Esta documentación sirve como referencia sobre el funcionamiento del código y las decisiones tomadas para resolver los ejercicios.

# Decisiones Tomadas

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado Lotería Nacional. Para la resolución de las mismas se utiliza como base el código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.

> Todo el código está redactado en inglés, con excepción de algunos logs específicos que permanecen en español para garantizar la compatibilidad con los tests proporcionados.

# Cambios en script de `generar-compose.sh`

Para evitar hardcodear un bet, se lo agrega como variable de entorno, que por el momento esta hardcodead en el script pero es modificable una vez tenida el archivo de docker-compose

```bash
for i in $(seq 1 "$AMOUNT_CLIENTS"); do
    cat >> "$YAML_FILE" << EOF

  client$i:
    container_name: client$i
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=$i
      - BET_NUMBER=67890
      - CLIENT_DOCUMENT=1
      - CLIENT_FIRST_NAME=Santi
      - CLIENT_LAST_NAME=Sev
      - CLIENT_BIRTHDATE=2000-08-10
    networks:
      - testing_net
    volumes:
      - ./client/config.yaml:/config.yaml
    depends_on:
      - server
EOF
done
```


# Arquitectura del Servidor

## Estructura de Directorios

El servidor se encuentra en la carpeta [server](https://github.com/SantiSev/tp0-santisev/blob/8a1db90ec3edaf039c5237e04bcac69bc3339ed7/server) en la raíz del repositorio que contiene estos componentes:

```bash
    server
    ├── common
    ├── config.ini
    ├── Dockerfile
    ├── main.py
    └── tests
```

- `common`: Módulo base que agrupa la lógica de negocio, protocolos de comunicación, configuración y utilidades esenciales del servidor.

- `config.ini`: Archivo de configuración.

- `Dockerfile`: Archivo usado para construir la imagen de Docker para el servidor _(provisto por la cátedra, no se modifica)._

- `main`: Entrypoint del Server

- `tests`: Carpeta que contiene pruebas automatizadas para verificar que el código funciona correctamente. _(provisto por la cátedra, no se modifica)._

## Directorio Common

Dentro de la seccion common tenemos los siguientes modulos:

```bash
server
├── common
│   ├── business
│   ├── config
│   ├── network
│   ├── protocol
│   ├── server
│   ├── session
│   └── utils
```

## Business

Contiene la lógica de negocio principal de la aplicación. Define las reglas y procesos específicos del dominio, separando la lógica empresarial de los detalles de implementación técnica.

Esta carpeta solo tiene el archivo `lottery_service.py` Que tiene como objetivo gestionar la lógica de negocio de las loterías.

En el alcance actual del ejercicio, se encarga exclusivamente de almacenar las apuestas utilizando la función `store_bets()` del archivo `utils.py`.

## Config

Maneja la configuración del sistema, incluyendo la lectura de archivos de configuración, variables de entorno, parámetros de inicialización e inicializacion de logs. Centraliza toda la gestión de configuración.

Esta carpeta solo tiene el archivo `config.py` que tiene 2 funcionciones:

`def initialize_config() -> ServerConfig`: Lee el archivo de configuración y las variables de entorno, creando una instancia de `ServerConfig` con todos los parámetros necesarios para el servidor. Esta instancia se utiliza para inicializar correctamente el servidor.

`def initialize_log(logging_level)`: Inicializa el sistema de logs del servidor según el nivel de logging especificado:

- **INFO**: Muestra mensajes de tipo INFO, ERROR y CRITICAL.
- **DEBUG**: Además de los anteriores, incluye mensajes de tipo DEBUG.

## Network

Proporciona las abstracciones y funcionalidades de red de bajo nivel. Maneja conexiones TCP, sockets y operaciones de comunicación básicas entre procesos.

Este modulo contiene 2 clases fundamnetales al server:

**ConnectionInterface**: Proporciona una abstracción de los servicios de sockets, permitiendo el uso de `send()`, `recv()` y `close()` sin la necesidad de manipular sockets directamente.

**ConnectionManager**: Implementa el patrón Acceptor del sistema. Permanece a la espera de conexiones entrantes y, una vez que un cliente se conecta al servidor, devuelve una instancia de `ConnectionInterface` para gestionar correctamente el envío y recepción de mensajes entre servidor y cliente.

### Manejo de Short Read/Write

La clase ConnectionInterface implementan mecanismos robustos para manejar lecturas y escrituras parciales:

- `receive function`: Utiliza el método `_receive_all()` para garantizar la recepción completa de datos mediante un bucle que continúa hasta obtener exactamente la cantidad de bytes solicitada, evitando problemas de short reads.

- `send function`: Utiliza `sendall()` para asegurar el envío completo de datos, manejando automáticamente las escrituras parciales que pueden ocurrir en redes congestionadas.

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

# Arquitectura del Cliente

## Estructura de Directorios

El servidor se encuentra en la carpeta [client](https://github.com/SantiSev/tp0-santisev/blob/050d2ac88f9682d9b8b60ad64b0c1c16bc8196da/client) en la raíz del repositorio que contiene estos componentes:

```bash
  client/
  ├── common
  ├── config.yaml
  ├── Dockerfile
  └── main.go
```

- `common`: Módulo base que agrupa la lógica de negocio, protocolos de comunicación, configuración y utilidades esenciales del cliente.

- `config.ini`: Archivo de configuración.

- `Dockerfile`: Archivo usado para construir la imagen de Docker para el cliente _(provisto por la cátedra, no se modifica)._

- `main`: Entrypoint del Client

- `tests`: Carpeta que contiene pruebas automatizadas para verificar que el código funciona correctamente. _(provisto por la cátedra, no se modifica)._

## Directorio Common

```bash
client
├── common
│   ├── business
│   ├── client
│   ├── config
│   ├── network
│   └── protocol
```

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