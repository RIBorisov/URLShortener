# URL Shortener service
[![Coverage Status](https://img.shields.io/badge/coverage-44.7%25-brightgreen)](https://github.com/RIBorisov/gophermart/coverage.html)
# Сервис URL-коротких ссылок

## Описание

Сервис предоставляет API для создания и управления короткими ссылками. Он использует роутер на основе библиотеки `go-chi` и включает в себя различные middleware для обработки запросов.

## Маршруты

### Основные маршруты

- **GET /{id}**: Получение короткой ссылки по идентификатору.
  -  Этот маршрут позволяет получить короткую ссылку по ее идентификатору.
  - **Пример**: `GET /123`

- **POST /**: Создание новой короткой ссылки.
  -  Этот маршрут позволяет создать новую короткую ссылку.
  - **Пример**: `POST /` с телом запроса, содержащим URL, который нужно сократить.

### API Маршруты

#### /api/shorten

- **POST /**: Создание новой короткой ссылки.
  -  Этот маршрут позволяет создать новую короткую ссылку через API.
  - **Пример**: `POST /api/shorten` с телом запроса, содержащим URL, который нужно сократить.

- **POST /batch**: Создание нескольких новых коротких ссылок.
  -  Этот маршрут позволяет создать несколько новых коротких ссылок за один запрос.
  - **Пример**: `POST /api/shorten/batch` с телом запроса, содержащим список URL, которые нужно сократить.

#### /api/user

- **GET /urls**: Получение списка коротких ссылок пользователя.
  -  Этот маршрут позволяет получить список всех коротких ссылок, созданных пользователем.
  - **Middleware**: `CheckAuth` - проверяет аутентификацию пользователя.
  - **Пример**: `GET /api/user/urls`

### Удаление ссылок

- **DELETE /api/user/urls**: Удаление всех ссылок пользователя.
  -  Этот маршрут позволяет удалить все короткие ссылки, созданные пользователем.
  - **Пример**: `DELETE /api/user/urls`

### Пинг

- **GET /ping**: Проверка доступности сервиса.
  -  Этот маршрут позволяет проверить доступность сервиса.
  - **Пример**: `GET /ping`

### Debug

- **/debug**: Профилирование сервиса.
  -  Этот маршрут позволяет профилировать сервис для отладки и оптимизации.
  - **Пример**: `GET /debug`

## Middleware

- **Recoverer**: Middleware для восстановления после паники.
  -  Этот middleware перехватывает панику и возвращает HTTP ответ с кодом 500.

- **Auth**: Middleware для аутентификации.
  -  Этот middleware проверяет аутентификацию пользователя перед обработкой запроса.

- **Gzip**: Middleware для сжатия ответов.
  -  Этот middleware сжимает ответы для уменьшения объема передаваемых данных.

- **Log**: Middleware для логирования запросов.
  -  Этот middleware логгирует запросы и ответы для мониторинга и отладки.

- **CheckAuth**: Middleware для проверки аутентификации пользователя.
  - Этот middleware проверяет аутентификацию пользователя перед обработкой запроса к маршруту `/api/user/urls`.

## Запуск тестов

Чтобы запустить тесты и проверить покрытие, выполните следующую команду из корня репозитория:
```bash
make tests
```