
# gh-checker

## Обзор

`gh-checker` — это легковесный сервис для взаимодействия с GitHub API, предназначенный для управления и отслеживания звезд и подписчиков GitHub репозиториев. Этот сервис позволяет периодически проверять, поставил ли пользователь звезду на конкретный репозиторий, и получать или обновлять список подписчиков для заданного пользователя. Информация хранится локально в базе данных SQLite с возможностью обновления на основе конфигурируемого интервала.

## Основные возможности

- Получение подписчиков GitHub для указанного пользователя.
- Проверка, поставил ли пользователь звезду на репозиторий.
- Локальное кэширование с использованием базы данных SQLite.
- Настраиваемое логирование и интервалы обновлений.
- Оптимизация работы с кэшированными и актуальными данными из GitHub API.
- Обработка ошибок и механизм повторных попыток запросов к GitHub API.

## Содержание

1. [Установка](#установка)
2. [Конфигурация](#конфигурация)
3. [Использование](#использование)
4. [Структура базы данных](#структура-базы-данных)
5. [Логирование](#логирование)
6. [API Эндпоинты](#api-эндпоинты)
7. [Участие в проекте](#участие-в-проекте)
8. [Лицензия](#лицензия)

## Установка

Для использования `gh-checker` требуется Go, установленный на вашем компьютере. Клонируйте репозиторий и установите необходимые зависимости:

```bash
git clone https://github.com/wrongjunior/gh-checker.git
cd gh-checker
go mod download
```

Соберите проект:

```bash
go build -o gh-checker
```

## Конфигурация

Конфигурация управляется через файл `config.yaml`. Пример конфигурации:

```yaml
github:
  api_key: "your-github-api-key"

database:
  path: "./gh-checker.db"

follower_check_interval: "10m"

logging:
  file_level: "info"
  console_level: "debug"
  file_path: "./app.log"
```

- `api_key`: Ключ API GitHub, необходимый для аутентификации.
- `path`: Путь к базе данных SQLite.
- `follower_check_interval`: Интервал для проверки новых подписчиков.
- `logging`: Уровни логов для файла и консоли, а также путь до файла логов.

## Использование

### Получение подписчиков

Для получения подписчиков указанного пользователя, сервис сначала проверяет, нужно ли обновить данные на основе интервала. Если кэшированные данные ещё действительны, будут возвращены они. В противном случае данные будут запрашиваться с GitHub API.

Пример использования на Go:

```go
import "gh-checker/internal/services"

followers, updated, err := services.UpdateFollowers("username", time.Hour*24)
if err != nil {
    log.Fatalf("Ошибка получения подписчиков: %v", err)
}

fmt.Println(followers)
```

### Проверка звёзд на репозитории

Чтобы проверить, поставил ли пользователь звезду на репозиторий:

```go
import "gh-checker/internal/services"

hasStar, err := services.UpdateStars("username", "repository", time.Hour*24)
if err != nil {
    log.Fatalf("Ошибка проверки звёзд: %v", err)
}

if hasStar {
    fmt.Println("Пользователь поставил звезду на репозиторий.")
} else {
    fmt.Println("Пользователь не поставил звезду на репозиторий.")
}
```

## Структура базы данных

Локальная база данных SQLite содержит три основные таблицы:

- `followers`: Хранит подписчиков пользователей GitHub.
- `last_check`: Хранит временные метки последней проверки подписчиков и звёзд.
- `stars`: Хранит информацию о звёздах, поставленных пользователями на репозитории.

Пример схемы базы данных:

```sql
CREATE TABLE IF NOT EXISTS followers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    follower TEXT NOT NULL,
    last_updated TIMESTAMP,
    UNIQUE(username, follower)
);

CREATE TABLE IF NOT EXISTS last_check (
    username TEXT NOT NULL,
    repository TEXT NOT NULL,
    last_checked TIMESTAMP,
    UNIQUE(username, repository)
);

CREATE TABLE IF NOT EXISTS stars (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    repository TEXT NOT NULL,
    last_updated TIMESTAMP,
    UNIQUE(username, repository)
);
```

## Логирование

Сервис поддерживает логирование как в файл, так и в консоль. Настраиваемые уровни логов и путь до логов указываются в файле `config.yaml`.

Поддерживаются следующие уровни логов:
- `debug`
- `info`
- `warn`
- `error`

Пример настройки логирования:

```go
import "gh-checker/internal/lib/logger"

logger.InitializeLogger(logger.LogConfig{
    FileLevel: "info",
    FilePath: "./app.log",
})

logger.Info("Приложение запущено")
```

## API Эндпоинты

Доступны следующие API эндпоинты:

### `/check-star`

Проверка, поставил ли пользователь звезду на репозиторий.

**Запрос:**

```json
{
  "username": "someuser",
  "repository": "somerepo"
}
```

**Ответ:**

```json
{
  "hasStar": true
}
```

### `/check-followers`

Проверка, является ли один пользователь подписчиком другого.

**Запрос:**

```json
{
  "follower": "userA",
  "followed": "userB"
}
```

**Ответ:**

```json
{
  "isFollowing": true
}
```

## Участие в проекте

См. [CONTRIBUTING.md](./CONTRIBUTING.md) для получения подробной информации о структуре коммитов.

## Лицензия

Этот проект лицензирован на условиях MIT License. См. файл [LICENSE](./LICENSE) для получения дополнительной информации.
