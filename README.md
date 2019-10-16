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

* *Syntax*: **include** *glob* ...;
* *Default*: &ndash;     
* *Context*: *         

## server
Настройки сервера

* *Syntax*: **server** { ... } 
* *Default*: &ndash;     
* *Context*: main

#### listen_ip
Слушать ip

* *Syntax*: **listen_ip** *ip*;
* *Default*: listen_ip 0.0.0.0;     
* *Context*: server

#### listen_port
Слушать port

* *Syntax*: **listen_port** *port*;
* *Default*: listen_port 3128;     
* *Context*: server

#### listen_schema
Тип сервера

* *Syntax*: **listen_schema** http;
* *Default*: listen_schema http;     
* *Context*: server

#### error_log
Лог ошибок

* *Syntax*: **error_log** *path* | off;
* *Default*: error_log off;
* *Context*: server

#### read_timeout
Максимальное время чтения входящего запроса (включая тело запроса)

* *Syntax*: **read_timeout** *time*;
* *Default*: read_timeout 0s;  
* *Context*: server

#### read_header_timeout
Максимальное время чтения заголовков входящего запроса

* *Syntax*: **read_header_timeout** *time*;
* *Default*: read_header_timeout 0s;  
* *Context*: server

#### write_timeout
Максимальное время ответа на входящий запрос

* *Syntax*: **write_timeout** *time*;
* *Default*: write_timeout 0s;  
* *Context*: server

#### idle_timeout
Максимальное время ожидания запроса до закрытия соединения

* *Syntax*: **idle_timeout** *time*;
* *Default*: idle_timeout 0s;  
* *Context*: server

## main

#### access_log
Лог запросов

* *Syntax*: **access_log** *path* | off;
* *Default*: access_log off; 
* *Context*: main, condition

#### dump_log
Расширенный лог запросов (для отладки)

* *Syntax*: **dump_log** *path* | off;
* *Default*: dump_log off;  
* *Context*: main, condition

#### basic_auth
"Basic" HTTP Authentication. Для включения требуется указать путь к файлу `htpasswd`

* *Syntax*: **basic_auth** *path* | off;
* *Default*: basic_auth off;  
* *Context*: main, condition

#### outgoing_ip_v4
IpV4 используемые prifma для запросов (случайный ip из списка)

* *Syntax*: **outgoing_ip_v4** *ip*... | off;
* *Default*: outgoing_ip_v4 0.0.0.0;  
* *Context*: main, condition

#### outgoing_ip_v6
IpV6 используемые prifma для запросов (случайный ip из списка)

* *Syntax*: **outgoing_ip_v6** *ip*... | off;
* *Default*: outgoing_ip_v6 ::0;  
* *Context*: main, condition

#### use_ip_header
Установить ip адрес для запроса исходя из переданных заголовков `Proxy-Use-IpV4` и `Proxy-Use-IpV6`

* *Syntax*: **use_ip_header** on | off;
* *Default*: use_ip_header off;  
* *Context*: main, condition

#### block_requests
Заблокировать входящие запросы (`423 Locked`)

* *Syntax*: **block_requests** on | off;
* *Default*: block_requests off;  
* *Context*: main, condition

## proxy_requests
Отправить исходящие запросы через прокси

* *Syntax*: **proxy_requests** *url* { ... } | *url*; | off;
* *Default*: proxy_requests off;  
* *Context*: main, condition

#### proxy_header
Установить заголовок, при отправке запроса через прокси указанный в `proxy_requests`

* *Syntax*: **proxy_header** *key* *val*;
* *Default*: &ndash; 
* *Context*: proxy_requests

## condition
Применить директивы при выполнении условия.

* *Syntax*: **condition** *key* *type* *val* { ... }
* *Default*: &ndash; 
* *Context*: proxy_requests

##### key
* `src_ip` - ip клиента
* `dst_domain` - домен, к которому будет выполнен исходящий запрос
* `header_*` - заголовок входящего запроса (например `header_user_agent`, `header_cookie`)

##### type 
* `=` - равенство
* `~` - регулярное выражение
* `cidr` - маска в формате cidr