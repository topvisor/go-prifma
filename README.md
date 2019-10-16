# topvisor/prifma
Простой proxy сервер на go

# build
Для сборки требуется golang https://golang.org/dl/

```shell script
go get -u github/topvisor/prifma/cmd/prifma
```

# run

```shell script
$GOPATH/bin/prifma --config /path/to/config.conf
```

# config
Пример конфигурации: https://github.com/topvisor/prifma/example/config/config.conf

#### include
Загрузить файлы конфигурации в текущий контекст

*Syntax* | **include** *glob* ...;
*Default* | &ndash;     
*Context* | *                  

## server
Настройки сервера

*Syntax* | **server** { ... } 
*Default* | &ndash;     
*Context* | main

#### listen_ip
Слушать ip

*Syntax* | **listen_ip** *ip*;
*Default* | listen_ip 0.0.0.0;     
*Context* | server

#### listen_port
Слушать port

*Syntax* | **listen_port** *port*;
*Default* | listen_port 3128;     
*Context* | server

#### listen_schema
Тип сервера

*Syntax* | **listen_schema** http;
*Default* | listen_schema http;     
*Context* | server

#### error_log
Лог ошибок

*Syntax* | **error_log** *path*;
*Default* | &ndash;  
*Context* | server

#### read_timeout
Максимальное время чтения входящего запроса (включая тело запроса)

*Syntax* | **read_timeout** *time*;
*Default* | read_timeout 0s;  
*Context* | server

#### read_header_timeout
Максимальное время чтения заголовков входящего запроса

*Syntax* | **read_header_timeout** *time*;
*Default* | read_header_timeout 0s;  
*Context* | server

#### write_timeout
Максимальное время ответа на входящий запрос

*Syntax* | **write_timeout** *time*;
*Default* | write_timeout 0s;  
*Context* | server

#### idle_timeout
Максимальное время ожидания запроса до закрытия соединения

*Syntax* | **idle_timeout** *time*;
*Default* | idle_timeout 0s;  
*Context* | server