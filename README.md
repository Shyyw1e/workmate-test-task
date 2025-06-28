# workmate-test-task

Task Manager — это простой HTTP-сервис на Go, позволяющий создавать, отслеживать и отменять длительные задачи (3–5 минут выполнения). Все данные хранятся **в памяти**, без внешних баз или очередей.

Проект написан с учётом инженерных практик и масштабируемости: логика разделена по слоям, есть поддержка контекста и цветное логирование с помощью `slog`.

---

## 🚀 Запуск

### ⚙️ Требования

- Go 1.21+
- Порт `:8080` свободен

### 🔧 Сборка и запуск

```bash
git clone https://github.com/yourname/task-manager
cd task-manager
go run cmd/main.go
```

Вот готовый фрагмент для вставки в твой `README.md`, оформленный в полноценном Markdown:

---

## 🌐 API Эндпоинты

### 🔹 POST `/tasks`

Создание новой задачи:

```bash
curl -X POST http://localhost:8080/tasks
```

**Ответ:**

```json
{
  "id": "3d1b6bc5-19d6-4e17-9c1e-7503f301f5d1",
  "status": "created",
  "created_at": "2025-06-25T12:00:00Z"
}
```

---

### 🔹 GET `/tasks/{id}`

Получить статус задачи:

```bash
curl http://localhost:8080/tasks/3d1b6bc5-19d6-4e17-9c1e-7503f301f5d1
```

**Ответ:**

```json
{
  "id": "3d1b6bc5-19d6-4e17-9c1e-7503f301f5d1",
  "status": "running",
  "created_at": "2025-06-25T12:00:00Z",
  "started_at": "2025-06-25T12:00:01Z"
}
```

---

### 🔹 GET `/tasks?status=running`

Получить список всех задач (опционально с фильтрацией по статусу):

```bash
curl http://localhost:8080/tasks
curl http://localhost:8080/tasks?status=running
```

Возможные значения параметра `status`:

* `created`
* `running`
* `done`
* `failed`

---

### 🔹 DELETE `/tasks/{id}`

Удалить (отменить) задачу:

```bash
curl -X DELETE http://localhost:8080/tasks/3d1b6bc5-19d6-4e17-9c1e-7503f301f5d1
```

**Ответ:**

* `204 No Content` — успешно удалено
* `404 Not Found` — задача не найдена

---

## 🧠 Архитектура

```
task-manager/
├── cmd/               # Точка входа (main.go)
├── internal/
│   └── task/          # Модель, сервис, HTTP-хендлеры
├── pkg/
│   └── logger/        # Кастомный цветной slog-логгер
├── go.mod / go.sum
└── README.md
```

---

## 📦 Задачи (`Task`)

```json
{
  "id": "string",
  "status": "created | running | done | failed",
  "created_at": "timestamp",
  "started_at": "timestamp",
  "finished_at": "timestamp",
  "duration": "number of seconds",
  "result": "string"
}
```

---

## 🎨 Цветной логгер

Используется `log/slog` с кастомным `TextHandler`, окрашивающим логи:

* 🟦 `INFO` — синий
* 🟨 `DEBUG` — жёлтый
* 🟥 `ERROR`, `WARN` — красный
* 🟩 `SUCCESS` (на базе INFO + атрибут) — зелёный

---

## 🧪 Тестирование

### Запуск всех тестов:

```bash
go test ./internal/task -v
```

### Покрытие:

* ✅ `TaskService` — создание, удаление, отмена
* ✅ HTTP-хендлеры (`POST`, `GET`, `DELETE`, `GET list`)
* ✅ Обработка ошибок и статусов

---

## 🔐 Особенности

* Все данные хранятся **в памяти** — перезапуск сервиса очищает список задач.
* Асинхронность реализована через `goroutine + context.WithCancel`.
* Возможность **отменить задачу до завершения**.
* Поддержка **фильтрации по статусу** (`?status=...`).

---
