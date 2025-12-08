# 01. Async/Await

## Задача
Выполнять асинхронные операции (I/O, HTTP запросы) без блокировки.

---

## Базовый синтаксис

```rust
import "lib/task" (async, await)

// Создать асинхронную задачу
task = async(fun() -> {
    // длительная операция
    42
})

// Дождаться результата
result = await(task)
print(result)
```

---

## Простой пример

```rust
import "lib/task" (async, await)
import "lib/time" (sleepMs)

fun slowAdd(a: Int, b: Int) -> Int {
    sleepMs(1000)  // 1 секунда
    a + b
}

task = async(fun() -> slowAdd(2, 3))
print("Task started...")
result = await(task)
print("Result: " ++ show(result))  // Result: 5
```

---

## Параллельное выполнение

```rust
import "lib/task" (async, await)
import "lib/time" (sleepMs)

fun fetchUser(id: Int) {
    sleepMs(500)
    { id: id, name: "User " ++ show(id) }
}

// Запускаем три запроса ПАРАЛЛЕЛЬНО
task1 = async(fun() -> fetchUser(1))
task2 = async(fun() -> fetchUser(2))
task3 = async(fun() -> fetchUser(3))

// Ждём все результаты
user1 = await(task1)
user2 = await(task2)
user3 = await(task3)

// Общее время: ~500ms вместо ~1500ms!
print(user1)
print(user2)
print(user3)
```

---

## HTTP запросы параллельно

```rust
import "lib/task" (async, await)
import "lib/http" (httpGet)
import "lib/json" (jsonDecode)

fun fetchData(url: String) {
    match httpGet(url) {
        Ok(resp) -> match jsonDecode(resp.body) {
            Ok(data) -> data
            Fail(e) -> { error: e }
        }
        Fail(e) -> { error: e }
    }
}

// Параллельные запросы к разным API
usersTask = async(fun() -> fetchData("https://api.example.com/users"))
postsTask = async(fun() -> fetchData("https://api.example.com/posts"))

users = await(usersTask)
posts = await(postsTask)

print("Got users and posts")
```

---

## Пул задач

```rust
import "lib/task" (async, await)
import "lib/list" (map)

fun processInParallel(items, process) {
    tasks = map(fun(item) -> async(fun() -> process(item)), items)
    map(fun(t) -> await(t), tasks)
}

urls = [
    "https://api.example.com/1",
    "https://api.example.com/2",
    "https://api.example.com/3"
]

results = processInParallel(urls, fun(url) -> { url: url, status: "fetched" })
print(results)
```

---

## Последовательность vs Параллельность

```rust
import "lib/task" (async, await)
import "lib/time" (sleepMs, timeNow)

fun slowTask(ms: Int) -> Int {
    sleepMs(ms)
    ms
}

// Последовательно: 300ms
start1 = timeNow()
slowTask(100)
slowTask(100)
slowTask(100)
elapsed1 = timeNow() - start1
print("Sequential: " ++ show(elapsed1) ++ "ms")

// Параллельно: ~100ms
start2 = timeNow()
t1 = async(fun() -> slowTask(100))
t2 = async(fun() -> slowTask(100))
t3 = async(fun() -> slowTask(100))
results = [await(t1), await(t2), await(t3)]
elapsed2 = timeNow() - start2
print("Parallel: " ++ show(elapsed2) ++ "ms")
```

---

## Практический пример: Web Scraper

```rust
import "lib/task" (async, await)
import "lib/http" (httpGet)
import "lib/list" (map)

fun scrapeUrl(url: String) {
    match httpGet(url) {
        Ok(resp) -> {
            url: url,
            status: resp.status,
            size: len(resp.body)
        }
        Fail(e) -> {
            url: url,
            status: 0,
            size: 0
        }
    }
}

urls = [
    "https://example.com",
    "https://github.com"
]

// Скрейпим все URL параллельно
tasks = map(fun(url) -> async(fun() -> scrapeUrl(url)), urls)
results = map(fun(t) -> await(t), tasks)

for result in results {
    match result {
        Ok(r) -> print(r.url ++ " - status " ++ show(r.status) ++ " - " ++ show(r.size) ++ " bytes")
        Fail(e) -> print("Error: " ++ e)
    }
}
```

---

## Когда использовать async

- I/O операции (файлы, сеть)
- HTTP запросы
- Базы данных
- Любые "ждущие" операции

Async позволяет не блокировать выполнение программы пока ждём результат.
