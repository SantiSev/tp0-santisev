# TP0 | Parte 1: Introducción a Docker | Ejercicio 4

Esta documentación sirve como referencia sobre el funcionamiento del código y las decisiones tomadas para resolver los ejercicios.

# Decisiones Tomadas

Este ejercio consiste en modificar servidor y cliente para que ambos sistemas terminen de forma graceful al recibir la signal SIGTERM.

## Server Implementacion

En el archivo [`server.py`](https://github.com/SantiSev/tp0-santisev/blob/e8065e4c7eb929eee945424698e593b9a7902405/server/common/server.py) agrege lo siguiente en la seccion de inicalizacion del server:

```py
class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

        signal.signal(signal.SIGTERM, self.__handle_shutdown)
```

y luego agrege este metodo del server:

```py
    def __handle_shutdown(self, signal_number, frame):
        """
        Handle server shutdown

        This function is called to gracefully shutdown the server,
        closing all active connections and freeing resources.
        """
        logging.debug(f"Signal received at frame: {frame}")

        if signal_number == signal.SIGTERM:
            logging.info('action: shutdown | result: in_progress')
            self._server_socket.close()
            logging.info('action: shutdown | result: success')
            exit(0)
```

Esta implementación establece un mecanismo de graceful shutdown para el servidor. Durante la inicialización, se registra un manejador de señales que asocia la señal SIGTERM con el método `__handle_shutdown`. 

Cuando el servidor recibe esta señal, intercepta la terminación abrupta, cierra el socket servidor de manera controlada y termina el proceso limpiamente.

## Client Implementacion

En el archivo [`client.py`](https://github.com/SantiSev/tp0-santisev/blob/d3386781721e1823b58647269191beb9235d39c5/client/common/client.go) agrege lo siguiente en la seccion del processo principal del cliente:


```go
func (c *Client) StartClientLoop() {

	// This is how i handle SIGTERM signals

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		<-sigChannel
		c.HandleShutdown()
		done <- true
	}()

	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {

		select {
		case <-done:
			log.Infof("action: exit | result: success | client_id: %v", c.config.ID)
			return
		default:
		}
    // . . . REST OF CODE
```

Esta implementación establece un mecanismo de **graceful shutdown** para el cliente utilizando goroutines y canales. 

Se crean dos canales: 
- **sigChannel:** para capturar señales del sistema operativo
- **done:** para coordinar la terminación. Una goroutine separada queda en espera de SIGTERM, y cuando la recibe, ejecuta el método de limpieza y señaliza la terminación a través del canal done.
Select statement en Go

El `select` statement en Go es una construcción de control que permite a una **goroutine esperar en múltiples operaciones de canal simultáneamente**. 

Funciona de manera similar a un switch, pero específicamente para **operaciones de canal**. En este caso, el **select verifica en cada iteración del loop si se ha recibido una señal de terminación a través del canal done.**

La estructura select evalúa todos los casos disponibles:

- `case <-done`: se ejecuta si hay un valor disponible en el canal done (señal de terminación)

- `default`: se ejecuta si ningún canal tiene datos disponibles, permitiendo que el loop continúe normalmente

Esta implementación permite que el cliente responda inmediatamente a señales de terminación sin tener que esperar a que termine la iteración actual del loop principal, proporcionando un shutdown más responsivo y controlado.


# Como Ejectuar

1. Crear el archivo de docker-compose mediante el uso del script `generar-compose.sh`

2. Correr el commando `make docker-compose-up` para iniciar los containers

    > **Opcional:** Podes ver el estado de los containers mediante el comando `docker ps`

3. Ejecutar `make docker-compose-logs` para ver los logs del servidor (se recomienda configurar el cliente para enviar mensajes durante un tiempo prolongado).

4. En otra terminal, correr `make docker-compose-down` para finalizar los procesos de forma controlada.

5. Volver a la terminal de logs para verificar que ambos procesos finalizaron con código de salida 0, indicando una terminación exitosa.)