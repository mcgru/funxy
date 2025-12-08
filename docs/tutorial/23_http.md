# HTTP клиент (lib/http)

Модуль `lib/http` предоставляет функции для выполнения HTTP-запросов.

```rust
import "lib/http" (*)
```

## Тип ответа

Все функции возвращают `Result<HttpResponse, String>`, где:

```rust
type HttpResponse = {
    status: Int,              // HTTP код статуса (200, 404, 500, ...)
    body: String,             // Тело ответа
    headers: List<(String, String)>  // Заголовки ответа
}
```

## Функции

### httpGet

```rust
httpGet(url: String) -> Result<HttpResponse, String>
```

Выполняет GET-запрос.

```rust
import "lib/http" (httpGet)

result = httpGet("https://api.example.com/users")

match result {
    Ok(resp) -> {
        print("Status: ${resp.status}")
        print("Body: ${resp.body}")
    }
    Fail(err) -> print("Error: ${err}")
}
```

### httpPost

```rust
httpPost(url: String, body: String) -> Result<HttpResponse, String>
```

Выполняет POST-запрос со строковым телом.

```rust
import "lib/http" (httpPost)

result = httpPost("https://api.example.com/users", "name=John&age=30")

match result {
    Ok(resp) -> print("Created: ${resp.status}")
    Fail(err) -> print("Error: ${err}")
}
```

### httpPostJson

```rust
httpPostJson(url: String, data: A) -> Result<HttpResponse, String>
```

Выполняет POST-запрос с автоматической сериализацией данных в JSON. Добавляет заголовок `Content-Type: application/json`.

```rust
import "lib/http" (httpPostJson)

user = { name: "John", age: 30, active: true }
result = httpPostJson("https://api.example.com/users", user)

match result {
    Ok(resp) -> print("Created: ${resp.status}")
    Fail(err) -> print("Error: ${err}")
}

// Списки тоже работают
items = [1, 2, 3]
httpPostJson("https://api.example.com/items", items)
```

### httpPut

```rust
httpPut(url: String, body: String) -> Result<HttpResponse, String>
```

Выполняет PUT-запрос для обновления ресурса.

```rust
import "lib/http" (httpPut)

result = httpPut("https://api.example.com/users/1", "{\"name\": \"Jane\"}")

match result {
    Ok(resp) -> print("Updated: ${resp.status}")
    Fail(err) -> print("Error: ${err}")
}
```

### httpDelete

```rust
httpDelete(url: String) -> Result<HttpResponse, String>
```

Выполняет DELETE-запрос.

```rust
import "lib/http" (httpDelete)

result = httpDelete("https://api.example.com/users/1")

match result {
    Ok(resp) -> print("Deleted: ${resp.status}")
    Fail(err) -> print("Error: ${err}")
}
```

### httpRequest

```
httpRequest(method: String, url: String, headers: List<(String, String)>, body: String = "", timeout: Int = 0) -> Result<HttpResponse, String>
```

Полный контроль над HTTP-запросом: метод, заголовки, тело, таймаут.

**Параметры с умолчаниями:**
- **body** — тело запроса (по умолчанию `""`)
- **timeout** — таймаут в миллисекундах. Если `0` — используется глобальный таймаут (по умолчанию 30000 мс)

```rust
import "lib/http" (httpRequest)

// Минимальный вызов - только method, url, headers
// body="" и timeout=0 (глобальный) подставятся автоматически
result = httpRequest("GET", "https://api.example.com/users", [])

// С телом, но глобальным таймаутом
headers = [("Content-Type", "application/json")]
result = httpRequest("POST", "https://api.example.com/users", headers, "{\"name\":\"John\"}")

// Полный контроль: все параметры явно
headers = [
    ("Authorization", "Bearer token123"),
    ("Content-Type", "application/json"),
    ("Accept", "application/json")
]
body = "{\"query\": \"search term\"}"
result = httpRequest("POST", "https://api.example.com/search", headers, body, 5000)

match result {
    Ok(resp) -> {
        print("Status: ${resp.status}")
        // Найти заголовок в ответе
        for header in resp.headers {
            (key, value) = header
            if key == "Content-Type" {
                print("Content-Type: ${value}")
            }
        }
    }
    Fail(err) -> print("Error: ${err}")
}
```

### httpSetTimeout

```rust
httpSetTimeout(milliseconds: Int) -> Nil
```

Устанавливает таймаут для всех последующих запросов (по умолчанию 30000 мс = 30 секунд).

```rust
import "lib/http" (httpSetTimeout, httpGet)

// Установить таймаут в 5 секунд
httpSetTimeout(5000)

// Теперь запросы будут прерываться через 5 секунд
result = httpGet("https://slow-api.example.com/data")
```

## Практические примеры

### Получение и парсинг JSON

```rust
import "lib/http" (httpGet)
import "lib/json" (jsonDecode)

type User = { name: String, email: String }

fun fetchUser(id: Int) -> Result<String, User> {
    match httpGet("https://api.example.com/users/${id}") {
        Ok(resp) -> {
            if resp.status == 200 {
                jsonDecode(resp.body)
            } else {
                Fail("HTTP ${resp.status}")
            }
        }
        Fail(err) -> Fail(err)
    }
}

match fetchUser(1) {
    Ok(user) -> print("User: ${user.name}")
    Fail(err) -> print("Error: ${err}")
}
```

### REST API клиент

```rust
import "lib/http" (httpGet, httpPostJson, httpDelete)
import "lib/json" (jsonDecode)

type User = { id: Int, name: String, email: String }

baseUrl = "https://api.example.com"

fun getUsers() -> Result<String, List<User>> {
    match httpGet("${baseUrl}/users") {
        Ok(resp) -> {
            if resp.status == 200 { jsonDecode(resp.body) }
            else { Fail("HTTP ${resp.status}") }
        }
        Fail(err) -> Fail(err)
    }
}

fun createUser(name: String, email: String) -> Result<String, User> {
    data = { name: name, email: email }
    match httpPostJson("${baseUrl}/users", data) {
        Ok(resp) -> {
            if resp.status == 201 { jsonDecode(resp.body) }
            else { Fail("HTTP ${resp.status}") }
        }
        Fail(err) -> Fail(err)
    }
}

fun deleteUser(userId: Int) -> Result<String, Nil> {
    match httpDelete("${baseUrl}/users/${userId}") {
        Ok(resp) -> {
            if resp.status == 204 { Ok(Nil) }
            else { Fail("HTTP ${resp.status}") }
        }
        Fail(err) -> Fail(err)
    }
}
```

### Работа с заголовками

```rust
import "lib/http" (httpRequest)

fun authenticatedGet(url: String, token: String) -> Result<String, String> {
    headers = [
        ("Authorization", "Bearer ${token}"),
        ("Accept", "application/json")
    ]
    
    match httpRequest("GET", url, headers, "") {
        Ok(resp) -> {
            if resp.status == 200 { Ok(resp.body) }
            else if resp.status == 401 { Fail("Unauthorized") }
            else { Fail("HTTP ${resp.status}") }
        }
        Fail(err) -> Fail(err)
    }
}
```

### Повторные попытки

```rust
import "lib/http" (httpGet)
import "lib/time" (sleepMs)

fun fetchWithRetry(url: String, maxRetries: Int) -> Result<String, String> {
    fun attempt(n: Int) -> Result<String, String> {
        if n > maxRetries {
            Fail("Max retries exceeded")
        } else {
            match httpGet(url) {
                Ok(resp) -> {
                    if resp.status == 200 { Ok(resp.body) }
                    else if resp.status >= 500 {
                        // Server error - retry
                        sleepMs(1000 * n)  // Exponential backoff
                        attempt(n + 1)
                    }
                    else { Fail("HTTP ${resp.status}") }
                }
                Fail(_) -> {
                    sleepMs(1000 * n)
                    attempt(n + 1)
                }
            }
        }
    }
    attempt(1)
}
```

## Сводка

| Функция | Тип | Описание |
|---------|-----|----------|
| `httpGet` | `String -> Result<HttpResponse, String>` | GET запрос |
| `httpPost` | `(String, String) -> Result<HttpResponse, String>` | POST со строкой |
| `httpPostJson` | `(String, A) -> Result<HttpResponse, String>` | POST с JSON |
| `httpPut` | `(String, String) -> Result<HttpResponse, String>` | PUT запрос |
| `httpDelete` | `String -> Result<HttpResponse, String>` | DELETE запрос |
| `httpRequest` | `(String, String, List<(String,String)>, String, Int) -> Result<HttpResponse, String>` | Полный контроль (с таймаутом) |
| `httpSetTimeout` | `Int -> Nil` | Установить таймаут |

## HTTP Сервер

### httpServe

```rust
httpServe(port: Int, handler: (HttpRequest) -> HttpResponse) -> Result<Nil, String>
```

Запускает HTTP-сервер на указанном порту. Блокирует выполнение программы.

```rust
type HttpRequest = {
    method: String,              // "GET", "POST", etc.
    path: String,                // "/api/users"
    query: String,               // "id=1&name=test"
    headers: List<(String, String)>,
    body: String
}
```

```
import "lib/http" (httpServe)

fun handler(req: HttpRequest) -> HttpResponse {
    if req.path == "/" {
        { status: 200, body: "Hello, World!", headers: [] }
    } else if req.path == "/api/data" {
        { status: 200, body: "{\"value\": 42}", headers: [("Content-Type", "application/json")] }
    } else {
        { status: 404, body: "Not Found", headers: [] }
    }
}

// Start server on port 8080
print("Starting server on http://localhost:8080")
httpServe(8080, handler)
```

### Пример: JSON API сервер

```
import "lib/http" (httpServe)
import "lib/json" (jsonEncode, jsonDecode)
import "lib/list" (length)

type User = { id: Int, name: String }
type UserInput = { name: String }

users = [
    { id: 1, name: "Alice" },
    { id: 2, name: "Bob" }
]

fun handler(req: HttpRequest) -> HttpResponse {
    match (req.method, req.path) {
        ("GET", "/users") -> {
            { status: 200, body: jsonEncode(users), headers: [("Content-Type", "application/json")] }
        }
        ("POST", "/users") -> {
            match jsonDecode(req.body) : Result<String, UserInput> {
                Ok(data) -> {
                    newUser: User = { id: length(users) + 1, name: data.name }
                    { status: 201, body: jsonEncode(newUser), headers: [("Content-Type", "application/json")] }
                }
                Fail(err) -> {
                    { status: 400, body: "{\"error\": \"Invalid JSON\"}", headers: [("Content-Type", "application/json")] }
                }
            }
        }
        _ -> { status: 404, body: "Not Found", headers: [] }
    }
}

httpServe(8080, handler)
```

### httpServeAsync (неблокирующий сервер)

```
httpServeAsync(port: Int, handler: (HttpRequest) -> HttpResponse) -> Int
```

Запускает HTTP-сервер в фоновом режиме и возвращает ID сервера. Не блокирует выполнение программы.

```
import "lib/http" (httpServeAsync, httpServerStop, httpGet)

fun handler(req: HttpRequest) -> HttpResponse {
    { status: 200, body: "Hello!", headers: [] }
}

// Запуск сервера в фоне
serverId = httpServeAsync(8080, handler)
print("Server started with ID: ${serverId}")

// Теперь можно делать запросы к серверу
response = httpGet("http://localhost:8080/")
print(response)

// Остановка сервера
httpServerStop(serverId)
print("Server stopped")
```

### httpServerStop

```rust
httpServerStop(serverId: Int) -> Nil
```

Останавливает запущенный сервер по его ID.

```rust
serverId = httpServeAsync(8080, handler)

// ... работа с сервером ...

httpServerStop(serverId)  // Graceful shutdown
```

## Ограничения

### Клиент
- Только синхронные запросы
- Нет встроенной поддержки cookies (можно передать в заголовках)
- Нет автоматических редиректов (можно обработать вручную)
- Таймаут глобальный (применяется ко всем запросам)

### Сервер
- Нет поддержки HTTPS (используйте reverse proxy)
- Нет маршрутизации (реализуйте в handler)

