# Тестовое задание на позицию стажера-бекендера

Цель задания – разработать чат-сервер, предоставляющий HTTP API для работы с чатами и сообщениями пользователя.

Детали реализации:

* Писать код можно на любом языке программирования
    * Сервис реализован на языке `Golang`
* В качестве хранилища данных можно использовать любую технологию
    * Для хранения данных используется `MySQL`
    * SQL описывающий БД находится в [файле](./db/install_db.sql)
    * Архитектура сервиса предоставляет возможность использовать др. хранилища но для этого требуется реализовать [интерфейс](./src/connector.go#L9)
* При перезапуске сервера добавленные данные должны сохраняться
    * данные `MySQL` хранятся в контейнере
* Сервер должен быть доступен на порту 9000
    * порт может быть задан переменной окружения `PORT`
* Визуализация данных в виде пользовательского интерфейса (веб-приложение, мобильное приложение) не требуется – достаточно только обозначенного ниже API, доступного из командной строки. Однако простор фантазии не ограничиваем, покуда соблюдаются основные требования
    * для просмотра записей можно использовать `Adminer` (порт 8080, логин root, пароль password)
    * контейнер использует сеть хоста, так что можно использовать инструменты IDE для просмотра БД
* Предоставить инструкцию по запуску приложения. В идеале (но не обязательно) – использовать контейнеризацию с возможностью запустить проект командой `docker-compose up`
    * файл [`docker-compose`](docker-compose.yml)
* Финальную версию нужно выложить на github.com

## Инструкция
Сервис реализован в соответствии с заданием.
Единственное расхождение id в `uint` а не в `string`.

В случае ошибки на стороне сервера возвещается код `500`

Если ошибка вызвана неверными данными то вернется код `400`
а в теле может присутствовать json обьясняющий что пошло не так, например:
~~~json
{
  "code": 0,
  "description": "Описание ошибки"
}
~~~  

Возможны следующие коды:
* 0 - сущность уже существует (при получении данных)
* 1 - сущность не существует (при создании)
* 2 - передан пустой параметр
## Основные сущности

Ниже перечислены основные сущности, которыми должен оперировать сервер.

### User

Пользователь приложения. Имеет следующие свойства:

* **id** - уникальный идентификатор пользователя
* **username** - уникальное имя пользователя
* **created_at** - время создания пользователя

### Chat

Отдельный чат. Имеет следующие свойства:

* **id** - уникальный идентификатор чата
* **name** - уникальное имя чата
* **users** - список пользователей в чате, отношение многие-ко-многим
* **created_at** - время создания

### Message

Сообщение в чате. Имеет следующие свойства:

* **id** - уникальный идентификатор сообщения
* **chat** - ссылка на идентификатор чата, в который было отправлено сообщение
* **author** - ссылка на идентификатор отправителя сообщения, отношение многие-к-одному
* **text** - текст отправленного сообщения
* **created_at** - время создания

## Основные API методы

Методы обрабатывают HTTP POST запросы c телом, содержащим все необходимые параметры в JSON.

### Добавить нового пользователя

Запрос:

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"username": "user_1"}' \
  http://localhost:9000/users/add
```

Ответ: `id` созданного пользователя или HTTP-код ошибки.

### Создать новый чат между пользователями

Запрос:

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name": "chat_1", "users": ["<USER_ID_1>", "<USER_ID_2>"]}' \
  http://localhost:9000/chats/add
```

Ответ: `id` созданного чата или HTTP-код ошибки.

Количество пользователей не ограничено.

### Отправить сообщение в чат от лица пользователя

Запрос:

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"chat": "<CHAT_ID>", "author": "<USER_ID>", "text": "hi"}' \
  http://localhost:9000/messages/add
```

Ответ: `id` созданного сообщения или HTTP-код ошибки.

### Получить список чатов конкретного пользователя

Запрос:

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"user": "<USER_ID>"}' \
  http://localhost:9000/chats/get
```

Ответ: cписок всех чатов со всеми полями, отсортированный по времени создания последнего сообщения в чате (от позднего к раннему). Или HTTP-код ошибки.

### Получить список сообщений в конкретном чате

Запрос:

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"chat": "<CHAT_ID>"}' \
  http://localhost:9000/messages/get
```

Ответ: список всех сообщений чата со всеми полями, отсортированный по времени создания сообщения (от раннего к позднему). Или HTTP-код ошибки.

## Фидбек

1) Время до запуска: 3 – ~25 минут (docker-compose network host, не прокидывался порт, коннект к mysql не использовал хост и порт)
2) Чистота кода: 3 - плохо написаны комменты, где-то их нет, в файле postgres лежит mysql, не очень симпатично структурировано, непонятный текст ошибок, неаккуратные имена методов с ошибками, не было go.sum, после вставки в бд читает из нее сразу же
3) Структура бд: 4 - E1_ префиксы, не использует IF EXISTS, но вообще - норм
4) Соответствие ТЗ: 4, не возвращаются id где-то, created_at пустой, но работает!
5) Описание в ридми: 1 – хорошее понятное описание есть детали реализации
6) Тесты: 0 - нет
7) Бонусы: 1 - adminer для просмотра данных

