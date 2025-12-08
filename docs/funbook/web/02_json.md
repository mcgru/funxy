# 02. JSON обработка

## Задача
Парсить, создавать и трансформировать JSON данные.

---

## Парсинг JSON строки

```rust
import "lib/json" (jsonDecode, jsonEncode)

// JSON строка с экранированием
jsonStr = "{\"name\": \"Alice\", \"age\": 30, \"active\": true}"

match jsonDecode(jsonStr) {
    Ok(data) -> {
        // Доступ к полям через точку
        print(data.name)    // Alice
        print(data.age)     // 30
        print(data.active)  // true
    }
    Fail(err) -> print("Parse error: " ++ err)
}
```

---

## Создание JSON из записей

```rust
import "lib/json" (jsonEncode)

// Записи автоматически сериализуются
user = { name: "Bob", age: 25 }
print(jsonEncode(user))
// {"name":"Bob","age":25}

// Вложенные записи
profile = {
    user: { name: "Alice", age: 30 },
    settings: { theme: "dark", notifications: true }
}
print(jsonEncode(profile))
```

---

## Создание JSON из Map

```rust
import "lib/json" (jsonEncode)
import "lib/map" (*)

// Map с динамическими ключами
scores = %{ "Alice" => 100, "Bob" => 85, "Carol" => 92 }
print(jsonEncode(scores))
// {"Alice":100,"Bob":85,"Carol":92}
```

---

## Работа с массивами

```rust
import "lib/json" (jsonEncode)
import "lib/list" (filter, map)

users = [
    { id: 1, name: "Alice", role: "admin" },
    { id: 2, name: "Bob", role: "user" },
    { id: 3, name: "Carol", role: "admin" }
]

// Фильтрация
admins = filter(fun(u) -> u.role == "admin", users)
print(jsonEncode(admins))
// [{"id":1,"name":"Alice","role":"admin"},{"id":3,"name":"Carol","role":"admin"}]

// Трансформация
names = map(fun(u) -> u.name, users)
print(jsonEncode(names))
// ["Alice","Bob","Carol"]
```

---

## Вложенные структуры

```rust
import "lib/json" (jsonEncode)

company = {
    name: "Acme Corp",
    departments: [
        {
            name: "Engineering",
            floor: 3,
            employees: [
                { name: "Alice", title: "Senior Dev" },
                { name: "Bob", title: "Junior Dev" }
            ]
        },
        {
            name: "Sales",
            floor: 1,
            employees: [
                { name: "Carol", title: "Manager" }
            ]
        }
    ]
}

print(jsonEncode(company))

// Доступ к вложенным данным
for dept in company.departments {
    print(dept.name ++ " (floor " ++ show(dept.floor) ++ "):")
    for emp in dept.employees {
        print("  - " ++ emp.name ++ ", " ++ emp.title)
    }
}
```

---

## Чтение JSON из файла

```rust
import "lib/io" (fileRead)
import "lib/json" (jsonDecode)

// fileRead возвращает Result
match fileRead("config.json") {
    Ok(content) -> match jsonDecode(content) {
        Ok(config) -> {
            print("Host: " ++ config.host)
            print("Port: " ++ show(config.port))
        }
        Fail(e) -> print("Invalid JSON: " ++ e)
    }
    Fail(err) -> print("Cannot read file: " ++ err)
}
```

---

## Запись JSON в файл

```rust
import "lib/io" (fileWrite)
import "lib/json" (jsonEncode)

config = {
    host: "localhost",
    port: 8080,
    debug: true,
    maxConnections: 100,
    allowedOrigins: ["http://localhost:3000", "https://example.com"]
}

match fileWrite("config.json", jsonEncode(config)) {
    Ok(_) -> print("Config saved successfully")
    Fail(err) -> print("Cannot write file: " ++ err)
}

```

---

## Трансформация данных

```rust
import "lib/json" (jsonEncode)
import "lib/list" (map)

// Входные данные
people = [
    { firstName: "Alice", lastName: "Smith", birthYear: 1990 },
    { firstName: "Bob", lastName: "Jones", birthYear: 1985 },
    { firstName: "Carol", lastName: "Wilson", birthYear: 1992 }
]

// Трансформация: создать новую структуру
transformed = map(fun(p) -> {
    fullName: p.firstName ++ " " ++ p.lastName,
    initial: p.firstName[0],
    age: 2024 - p.birthYear
}, people)

print(jsonEncode(transformed))
// [{"fullName":"Alice Smith","initial":"A","age":34},...]
```

---

## Фильтрация и агрегация

```rust
import "lib/json" (jsonEncode)
import "lib/list" (filter, map, foldl)

orders = [
    { id: 1, customer: "Alice", total: 150.0, status: "completed" },
    { id: 2, customer: "Bob", total: 89.50, status: "pending" },
    { id: 3, customer: "Alice", total: 220.0, status: "completed" },
    { id: 4, customer: "Carol", total: 45.0, status: "completed" }
]

// Только завершённые заказы
completed = filter(fun(o) -> o.status == "completed", orders)

// Общая сумма
totalRevenue = foldl(fun(acc, o) -> acc + o.total, 0.0, completed)
print("Total revenue: $" ++ show(totalRevenue))  // $415.0

// Заказы по клиенту
aliceOrders = filter(fun(o) -> o.customer == "Alice", orders)
aliceTotal = foldl(fun(acc, o) -> acc + o.total, 0.0, aliceOrders)
print("Alice total: $" ++ show(aliceTotal))  // $370.0
```

---

## API Response обработка

```rust
import "lib/json" (jsonEncode)
import "lib/list" (map, filter, foldl)

// Симуляция ответа API
apiResponse = {
    status: "success",
    data: {
        users: [
            { id: 1, name: "Alice", email: "alice@example.com", active: true },
            { id: 2, name: "Bob", email: "bob@example.com", active: false },
            { id: 3, name: "Carol", email: "carol@example.com", active: true }
        ],
        pagination: { page: 1, totalPages: 5, perPage: 10 }
    }
}

// Извлечение только активных пользователей
activeUsers = filter(fun(u) -> u.active, apiResponse.data.users)

// Создание ответа для фронтенда
frontendResponse = {
    users: map(fun(u) -> { id: u.id, name: u.name }, activeUsers),
    meta: {
        count: len(activeUsers),
        page: apiResponse.data.pagination.page
    }
}

print(jsonEncode(frontendResponse))
```

---

## Сравнение: императивный vs функциональный

```rust
import "lib/list" (filter, map, foldl)

products = [
    { name: "Laptop", price: 999.0, category: "electronics" },
    { name: "Book", price: 15.0, category: "books" },
    { name: "Phone", price: 699.0, category: "electronics" },
    { name: "Desk", price: 200.0, category: "furniture" }
]

// Задача: найти среднюю цену электроники

// Императивный стиль
electronics1 = []
for p in products {
    if p.category == "electronics" {
        electronics1 = electronics1 ++ [p]
    }
}
sum1 = 0.0
for e in electronics1 {
    sum1 += e.price
}
avg1 = sum1 / intToFloat(len(electronics1))
print("Imperative avg: " ++ show(avg1))

// Функциональный стиль (одна цепочка!)
electronics2 = filter(fun(p) -> p.category == "electronics", products)
sum2 = foldl(fun(acc, p) -> acc + p.price, 0.0, electronics2)
avg2 = sum2 / intToFloat(len(electronics2))
print("Functional avg: " ++ show(avg2))
```

---

## Построение JSON для API

```rust
import "lib/json" (jsonEncode)

// Функция для создания стандартного API ответа
fun apiSuccess(data) {
    {
        status: "success",
        data: data,
        error: Nil
    }
}

fun apiError(message: String, code: Int) {
    {
        status: "error",
        data: Nil,
        error: { message: message, code: code }
    }
}

// Использование
print(jsonEncode(apiSuccess({ user: { id: 1, name: "Alice" } })))
// {"status":"success","data":{"user":{"id":1,"name":"Alice"}},"error":null}

print(jsonEncode(apiError("Not found", 404)))
// {"status":"error","data":null,"error":{"message":"Not found","code":404}}
```
