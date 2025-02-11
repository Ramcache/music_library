# 🎶 Онлайн-библиотека песен (Music Library)

Этот проект представляет собой RESTful API для управления библиотекой песен. Он позволяет добавлять, удалять, изменять и получать информацию о песнях, а также получать текст песни с пагинацией по куплетам. Проект использует PostgreSQL для хранения данных и интегрируется с внешним API для получения дополнительной информации о песнях.

---

## 📋 Функциональность

- **Добавление песни**: Добавление новой песни в формате JSON с автоматическим обогащением данными из внешнего API.
- **Получение списка песен**: Получение списка песен с фильтрацией по всем полям и пагинацией.
- **Получение текста песни**: Получение текста песни с пагинацией по куплетам.
- **Удаление песни**: Удаление песни по её ID.
- **Изменение данных песни**: Обновление информации о песне по её ID.

---

## 🛠️ Технологии

- **Язык программирования**: Go
- **Фреймворк**: Gin
- **База данных**: PostgreSQL
- **Миграции**: Goose
- **Документация API**: Swagger (используется `swag` для генерации)
- **Логирование**: Logrus
- **Конфигурация**: .env файл

---

## 🚀 Установка и запуск

### 1. Клонирование репозитория
```bash
git clone https://github.com/Ramcache/music_library.git
cd music_library
```

### 2. Настройка базы данных
1. Установите PostgreSQL и создайте базу данных:
   ```bash
   createdb songs
   ```
2. Настройте подключение к базе данных в файле `.env`:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER="user"
   DB_PASSWORD="password"
   DB_NAME="dbname"
   ```

### 3. Применение миграций
Установите `goose` и примените миграции:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir ./migrations postgres "user=postgres password=postgres dbname=songs sslmode=disable" up
```

### 4. Установка зависимостей
```bash
go mod tidy
```

### 5. Запуск сервера
```bash
go run cmd/server/main.go
```

Сервер будет доступен по адресу: `http://localhost:8080`.

---

## 📖 Документация API

Документация API доступна через Swagger UI:
- Откройте в браузере: `http://localhost:8080/swagger/index.html`.

Для генерации документации используйте `swag`:
```bash
swag init -g cmd/server/main.go
```

---

## 🛠️ Примеры запросов

### Добавление песни
```bash
curl -X POST -H "Content-Type: application/json" -d '{"group":"Muse","song":"Supermassive Black Hole"}' http://localhost:8080/songs
```

### Получение списка песен
```bash
curl http://localhost:8080/songs
```

### Получение текста песни
```bash
curl http://localhost:8080/songs/1/text?page=1&limit=2
```

### Удаление песни
```bash
curl -X DELETE http://localhost:8080/songs/1
```

### Обновление песни
```bash
curl -X PUT -H "Content-Type: application/json" -d '{"group":"Muse","song":"Starlight"}' http://localhost:8080/songs/1
```

---

## 🧑‍💻 Разработка

### Структура проекта
```
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── models/
│   └── service/
├── migrations/
├── docs/
├── .env
├── go.mod
├── go.sum
└── README.md
```

### Логирование
Логирование настроено с использованием `logrus`. Логи выводятся в консоль с уровнями `debug` и `info`.


## 🤝 Контакты

Если у вас есть вопросы или предложения, свяжитесь со мной:
- **Email**: ramcache@yandex.ru
- **GitHub**: [Ramcache](https://github.com/Ramcache)

---
