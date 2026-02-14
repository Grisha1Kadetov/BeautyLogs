# BeautyLogs

Линтер для Go-проектов, проверяющий **строки логирования**  
(реализован на базе `golang.org/x/tools/go/analysis`).

## Поддержка любых логгеров

`BeautyLogs` **не привязан к конкретным библиотекам логирования**.  
Он может работать с **любыми логгерами**, если указать их в конфигурации.
(формат: `<пакет>:<функция1,функция2,...>`)

---

## Что проверяет

Для функций логирования проверяються правила:

- **`only-eng`** — только латинские буквы
- **`lowercase`** — первая буква должна быть строчной
- **`special-char`** — запрет спецсимволов и эмодзи
- **`sensitive`** — запрет чувствительных данных (учитываеться `camelCase`, например: `log.Print("apiKey")`)

### Значения по умолчанию

**Чувствительные ключевые слова:**
password, secret, token, key, credential, username, email, phone

**Разрешённые спецсимволы:**
:,.-_()=

**Поддерживаемые логгеры по умолчанию:**
- `fmt` : `Printf`, `Println`, `Print`
- `go.uber.org/zap` : `Debug`, `Info`, `Warn`, `Error`, и др.
- `log/slog` : `Debug`, `Info`, `Warn`, `Error`, `Log`, и др.

---

## Использование как самостоятельного линтера

Точка входа:  
`cmd/beautylogs/main`

### Сборка
```bash
go build -o beautylogs ./cmd/beautylogs
```

### Запуск
```bash
./beautylogs ./...
```

### Флаги
```bash
./beautylogs \
  -only-eng=true \
  -lowercase=true \
  -special-char=true \
  -sensitive=true \
  -ignore-special-chars=":,.-_()=" \
  -sensitive-keys="password,secret,token,key,credential,username,email,phone" \
  -logger="fmt:Printf,Println,Print" \
  -logger="log/slog:Info,Error,Warn,Debug" \
  ./...
```

Примечания:
- sensitive-keys — список через запятую с ключивыми словами для правила **`sensitive`**
- ignore-special-chars — строка из разрешённых рун для **`special-char`**
- logger можно указывать несколько раз формат: `<пакет>:<функция1,функция2,...>`

## Интеграция с golangci-lint (module plugin)

### Сборка кастомного golangci-lint

Для использования линтера как плагина `golangci-lint` требуется собрать
кастомный бинарник.

Создай файл `.custom-gcl.yml`:

```yaml
version: v2.9.0
plugins:
  - module: "github.com/Grisha1Kadetov/BeautyLogs"
    import: "github.com/Grisha1Kadetov/BeautyLogs/pkg/plugin"
    version: v0.1.1
```

Собери бинарник:

```bash
golangci-lint custom -v
```

В результате будет создан файл `./custom-gcl`.

---

### Использование

Запуск линтера:

```bash
./custom-gcl run ./...
```

---

### Пример конфигурации `.golangci.yml`

```yaml
version: "2"

linters:
  default: none
  enable:
    - beautylogs
  settings:
    custom:
      beautylogs:
        type: "module"
        description: "BeautyLogs: проверка строк логирования"
        settings:
          only-eng: true
          lowercase: true
          special-char: true
          sensitive: true

          sensitive-keys:
            - password
            - secret
            - token
            - key
            - credential
            - username
            - email
            - phone

          ignore-special-chars: ":,.-_()="

          logger:
            - "fmt:Printf,Println,Print"
            - "log/slog:Info,Error,Warn,Debug"
```
