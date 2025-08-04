# TestBeeline

Сервис для обработки XML-данных пользователей с последующей отправкой преобразованных данных на удаленный сервер.

## Описание

TestBeeline - это HTTP-сервер, который:
- Принимает XML-данные с информацией о пользователях
- Преобразует возраст пользователей в возрастные группы
- Конвертирует данные в JSON формат
- Отправляет обработанные данные на удаленный сервер

## Возможности

- ✅ Асинхронная обработка пользователей с использованием горутин
- ✅ Структурированное логирование с помощью slog
- ✅ Обработка ошибок и отслеживание запросов через Request ID
- ✅ RESTful API на базе Echo framework
- ✅ Автоматическая категоризация пользователей по возрастным группам

## Требования

- Go 1.21+
- Docker и Docker Compose (для контейнерного запуска)
- Make (для использования Makefile)
- Echo web framework
- Стандартная библиотека Go (encoding/xml, encoding/json, net/http)

## Установка

```bash
git clone https://github.com/skywalker-jpg/Beeline.git
cd TestBeeline
go mod download
```

## Конфигурация

Перед запуском необходимо настроить:
- URL удаленного сервера (`serverURL`)
- Порт для запуска сервиса
- Уровень логирования

Пример конфигурации можно задать через переменные окружения или конфигурационный файл.

## API Endpoints

### POST /api/v1/process

Обрабатывает XML-данные пользователей и отправляет их на удаленный сервер.

**Request:**
- Content-Type: application/xml
- Body: XML с данными пользователей

**Пример входных данных:**
```xml
<users>
    <user>
        <id>1</id>
        <name>Иван Иванов</name>
        <email>ivan@example.com</email>
        <age>28</age>
    </user>
    <user>
        <id>2</id>
        <name>Петр Петров</name>
        <email>petr@example.com</email>
        <age>40</age>
    </user>
</users>
```

**Response:**
```json
{
    "message": "Processing completed",
    "users_received": 2,
    "users_sent": 2,
    "errors": 0
}
```

## Структура данных

### Входной формат (XML)
```go
type User struct {
    ID    string
    Name  string
    Email string
    Age   int
}
```

### Выходной формат (JSON)
```go
type UserJSON struct {
    ID       string `json:"id"`
    FullName string `json:"fullName"`
    Email    string `json:"email"`
    AgeGroup string `json:"ageGroup"`
}
```

### Возрастные группы
- **до 25**: возраст < 25
- **от 25 до 35**: возраст >= 25 и <= 35
- **старше 35**: возраст > 35

## Сборка и запуск

### Локальная сборка

### Запуск через Docker Compose

```bash
# Сборка бинарного файла
make build

# Запуск всех сервисов
make up

# Остановка сервисов
make down

# Очистка сборочных артефактов
make clean
```

### Makefile команды

- `make build` - компиляция приложения в `./bin/beeline/`
- `make up` - запуск через docker-compose с автоматической пересборкой
- `make down` - остановка всех контейнеров
- `make clean` - удаление скомпилированных файлов

### Переменные окружения для сборки

- `GOOS` - целевая ОС (по умолчанию: linux)
- `GOARCH` - архитектура процессора (автоматически определяется)
- `CGO` - включение CGO (по умолчанию: 0 - выключено)

## Логирование

Сервис использует структурированное логирование с помощью `slog`. Каждый запрос отслеживается через уникальный Request ID.

Уровни логирования:
- **INFO**: Основные операции (получение запроса, отправка данных)
- **ERROR**: Ошибки обработки
- **DEBUG**: Детальная информация о обработке каждого пользователя

## Обработка ошибок

Сервис корректно обрабатывает следующие типы ошибок:
- Невалидный XML формат
- Ошибки чтения тела запроса
- Ошибки отправки данных на удаленный сервер
- Ошибки сетевого соединения

## Производительность

- Параллельная обработка пользователей с использованием горутин
- Использование каналов для сбора результатов
- Эффективная работа с памятью при больших объемах данных

## Примеры использования

### Отправка данных через curl:
```bash
curl -X POST http://localhost:8080/api/v1/process \
  -H "Content-Type: application/xml" \
  -d @users.xml
```

### Пример интеграции на Go:
```go
client := &http.Client{Timeout: 30 * time.Second}
xmlData := []byte(`<users>...</users>`)

req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/process", bytes.NewBuffer(xmlData))
req.Header.Set("Content-Type", "application/xml")

resp, err := client.Do(req)
// обработка ответа
```

## Contributing

1. Fork репозиторий
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Commit изменения (`git commit -m 'Add some amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Откройте Pull Request