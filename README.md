# TP0 | Parte 2: Repaso de Comunicaciones | Ejercicio 5

Esta documentación sirve como referencia sobre el funcionamiento del código y las decisiones tomadas para resolver los ejercicios.

# Decisiones Tomadas

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado Lotería Nacional. Para la resolución de las mismas se utiliza como base el código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.

> Todo el código está redactado en inglés, con excepción de algunos logs específicos que permanecen en español para garantizar la compatibilidad con los tests proporcionados.

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

**ConnectionInterface**: Proporciona una abstracción de los servicios de sockets, permitiendo el uso de `send()` y `recv()` de información sin necesidad de manipular sockets directamente.

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

  **Limitación Actual:** Segun los requisitos de este ejercicio, por el momento solo procesa un solo bet de un solo cliente, pero esta arquitectura se va a escalar en los proximos ejercicios para

## Session

Gestiona las sesiones de usuario o conexión. Mantiene el estado de las interacciones, autenticación y el contexto de cada cliente conectado.

## Utils

Contiene utilidades auxiliares y funciones helper proporcionadas por la catedra para poder leer / escribir el archivo donde se almacenan los **Bets** (apuestas).

Contiene únicamente el archivo `utils.py` proporcionado por la cátedra, el cual no puede ser modificado según las especificaciones del enunciado.

# Arquitectura del Cliente

# Como Ejectuar

1. Crear el archivo de docker-compose mediante el uso del script `generar-compose.sh`

2. Correr el commando `make docker-compose-up` para iniciar los containers

   > **Opcional:** Podes ver el estado de los containers mediante el comando `docker ps`

3. Ejecutar `make docker-compose-logs` para ver los logs del servidor (se recomienda configurar el cliente para enviar mensajes durante un tiempo prolongado).

4. En otra terminal, correr `make docker-compose-down` para finalizar los procesos de forma controlada.

5. Volver a la terminal de logs para verificar que ambos procesos finalizaron con código de salida 0, indicando una terminación exitosa.)
