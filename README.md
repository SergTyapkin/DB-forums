# API базы данных для форума
Данный репозиторий содержит реализацию API для курса СУБД в рамках ["Технопарка"](https://park.mail.ru/).

## Немного о коде
* Он реализует API, описанное на [app.swaggerhub.com](https://app.swaggerhub.com/apis-docs/SergTyapkin/DB-froums/0.1.0)
* Написан файл для формирования Docker-контейнера, который разворачивает реализованный сервис вместе с PostgreSQL 9.5.

## Состав Docker-контейнеров
Docker-контейнер организован по следующему приципу:

* Приложение:
    * использует и объявляет порт 5000 (http);
* PostgreSQL:
    * использует и объявляет порт 5342 (http);

## Как собрать и запустить контейнер
Для сборки контейнера монжо выполнить команду вида:
```bash
docker build -t s.tyapkin https://github.com/SergTyapkin/DB-forums.git
```
Или команды:
```bash
git clone https://github.com/SergTyapkin/DB-forums.git DB-forums
cd DB-forums/
docker build -t s.tyapkin .
```

После этого будет создан Docker-образ с именем `s.tyapkin` (опция `-t`).

Запустить ранее собранный контейнер можно командой вида:
```bash
docker run -p 5000:5000 --name s.tyapkin -t s.tyapkin
```
После этого можно получить доступ к запущенному в контейнере приложению по URL: ```http://localhost:5000/``` (базовая часть URL: ```/api```)

Получить список запущенных контейнеров можно командой:
```bash
docker ps
```

Остановить контейнер можно командой:
```bash
docker kill s.tyapkin
```
