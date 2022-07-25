## Running Direction:

* Prerequisite: [GNU make](https://www.gnu.org/software/make/) && installed [Docker](https://www.docker.com/products/docker-desktop/) at your local machine.

* Build up and start all the services listed in the ```docker-compose.yml```
```
make up_build
```
* Start the front-end connection:
```
make start
```
At this moment, you can hit the url: http://localhost:8081 to browse the website.
* Stop the front-end connection:
```
make stop
```
* Make docker-compose down:
```
make down
```

* All ```make``` instructions:
```make
EJChen@EJsMacBook GoMicro_project % make help
 Choose a command:
  up                    starts all containers in the background without forcing build
  up_build              stops docker-compose(if running), builds all project and starts docker compose
  down                  stop docker compose
  compile_broker        builds the broker binary as a linux executable
  compile_auth          builds the auth binary as a linux executable
 compile_logger         builds the logger service binary
 compile_listener       builds the listener service binary
 compile_mail           build mail service binary
  compile_front         builds the front-end binary
  compile_front_linux   builds the front-end linux executable
  start                 starts the front-end
  stop                  stop the front-end
```

