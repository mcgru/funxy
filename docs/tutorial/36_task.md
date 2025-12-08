# Асинхронные вычисления (lib/task)

`lib/task` предоставляет Task — абстракцию для асинхронных вычислений (Futures/Promises).

## Импорт

```rust
import "lib/task" (*)
```

## Тип Task

`Task<T>` представляет асинхронное вычисление, которое:
- Может завершиться успешно со значением типа T
- Может завершиться с ошибкой (String)
- Может быть отменено

## Создание задач

### async — запуск асинхронного вычисления

```rust
import "lib/task" (async, await)

// Запускает функцию в отдельной горутине
task = async(fun() -> Int {
    // тяжёлые вычисления
    42 * 42
})

result = await(task)
print(result)  // Ok(1764)
```

### taskResolve / taskReject — готовые задачи

```rust
import "lib/task" (taskResolve, taskReject, await)

// Уже завершённая задача со значением
resolved = taskResolve(42)
print(await(resolved))  // Ok(42)

// Уже завершённая задача с ошибкой
fun getRejected() -> Result<String, Int> {
    rejected = taskReject("something went wrong")
    await(rejected)
}
print(getRejected())  // Fail("something went wrong")
```

## Ожидание результата

### await — базовое ожидание

```rust
import "lib/task" (async, await)

task = async(fun() -> String { "fetched data" })

match await(task) {
    Ok(data) -> print("Got: ${data}")
    Fail(err) -> print("Error: ${err}")
}
```

### awaitTimeout — ожидание с таймаутом

```rust
import "lib/task" (async, awaitTimeout)

task = async(fun() -> Int { 42 })

// Таймаут в миллисекундах
match awaitTimeout(task, 5000) {
    Ok(result) -> print("Got result: ${result}")
    Fail("timeout") -> print("Task took too long!")
    Fail(err) -> print("Error: ${err}")
}
```

## Множественные задачи

### awaitAll — ждать все

```rust
import "lib/task" (async, awaitAll)

tasks = [
    async(fun() -> Int { 1 }),
    async(fun() -> Int { 2 }),
    async(fun() -> Int { 3 })
]

match awaitAll(tasks) {
    Ok(results) -> {
        // results = [1, 2, 3]
        print(results)
    }
    Fail(err) -> print("One task failed: ${err}")
}
```

### awaitAny — первый успешный

```rust
import "lib/task" (async, awaitAny)

// Возвращает первый успешный результат, игнорируя ошибки
tasks = [
    async(fun() -> String { "server1" }),
    async(fun() -> String { "server2" }),
    async(fun() -> String { "server3" })
]

match awaitAny(tasks) {
    Ok(data) -> print("Got data: ${data}")
    Fail(_) -> print("All servers failed")
}
```

### awaitFirst — первый завершённый

```rust
import "lib/task" (async, awaitFirst)

tasks = [
    async(fun() -> Int { 1 }),
    async(fun() -> Int { 2 })
]

// Возвращает первый завершённый результат (успех или ошибка)
match awaitFirst(tasks) {
    Ok(v) -> print("First completed with: ${v}")
    Fail(e) -> print("First completed with error: ${e}")
}
```

### Варианты с таймаутом

```rust
import "lib/task" (async, awaitAllTimeout, awaitAnyTimeout, awaitFirstTimeout)

tasks = [
    async(fun() -> Int { 1 }),
    async(fun() -> Int { 2 })
]

// awaitAllTimeout — все с таймаутом
print(awaitAllTimeout(tasks, 10000))

// awaitAnyTimeout — первый успешный с таймаутом
// print(awaitAnyTimeout(tasks, 5000))

// awaitFirstTimeout — первый завершённый с таймаутом
// print(awaitFirstTimeout(tasks, 5000))
```

## Управление задачами

### taskIsDone — проверка завершения

```rust
import "lib/task" (async, taskIsDone)

task = async(fun() -> Int { 42 })

// Неблокирующая проверка
if taskIsDone(task) { print("Task finished") }
else { print("Still running...") }
```

### taskCancel — отмена

```rust
import "lib/task" (async, taskCancel, taskIsCancelled)

task = async(fun() -> Int {
    // долгая операция
    42
})

// Отмена (если задача ещё не началась)
taskCancel(task)

if taskIsCancelled(task) { print("Task was cancelled") }
```

## Пул горутин

По умолчанию максимум 1000 параллельных задач. Можно настроить:

```rust
import "lib/task" (taskGetGlobalPool, taskSetGlobalPool)

// Получить текущий лимит
poolLimit = taskGetGlobalPool()
print("Current limit: ${poolLimit}")

// Установить новый лимит
taskSetGlobalPool(2000)
```

## Комбинаторы

### taskMap — трансформация результата

```rust
import "lib/task" (async, await, taskMap)

task = async(fun() -> Int { 10 })
doubled = taskMap(task, fun(x) -> Int { x * 2 })

match await(doubled) {
    Ok(v) -> print(v)  // 20
    Fail(_) -> print("failed")
}
```

### taskFlatMap — цепочка задач

```rust
import "lib/task" (async, await, taskFlatMap)

type User = { id: Int, name: String }
type Post = { title: String }

fun fetchUser(userId: Int) {
    async(fun() -> User { { id: userId, name: "Alice" } })
}

fun fetchPosts(user: User) {
    async(fun() -> List<Post> { [{ title: "Post 1" }] })
}

// Цепочка: сначала пользователь, потом его посты
result = taskFlatMap(fetchUser(1), fetchPosts)

match await(result) {
    Ok(posts) -> print(posts)
    Fail(err) -> print("Error: ${err}")
}
```

### taskCatch — обработка ошибок

```rust
import "lib/task" (async, await, taskCatch)

riskyTask = async(fun() -> Int {
    if true { panic("failure") } else { 42 }
})

// Восстановление при ошибке
safeTask = taskCatch(riskyTask, fun(err) -> Int {
    0  // значение по умолчанию
})

match await(safeTask) {
    Ok(v) -> print("Got: " ++ show(v))
    Fail(e) -> print("Handler failed: " ++ e)
}
```

**Важно:** Если функция-обработчик внутри `taskCatch` сама вызовет `panic`, результат будет `Fail`.

## Паттерны использования

### Параллельная обработка

```rust
import "lib/task" (async, awaitAll)
import "lib/list" (map)

// Параллельная обработка списка
items = [1, 2, 3]
tasks = map(fun(item) -> async(fun() -> Int { item * 2 }), items)
result = awaitAll(tasks)
print(result)  // Ok([2, 4, 6])
```

### Worker pool

```rust
import "lib/task" (async, awaitAll, taskSetGlobalPool)
import "lib/list" (map)

// Обработка очереди задач с ограничением параллелизма
queue = [1, 2, 3, 4, 5]
taskSetGlobalPool(2)  // Максимум 2 параллельных задачи

tasks = map(fun(job) -> async(fun() -> Int { job * 10 }), queue)
result = awaitAll(tasks)
print(result)  // Ok([10, 20, 30, 40, 50])
```

### Гонка с таймаутом

```rust
import "lib/task" (async, await, awaitTimeout)

// Fetch with fallback
fun fetchWithFallback() -> Result<String, String> {
    primary = async(fun() -> String { "primary data" })
    
    match awaitTimeout(primary, 1000) {
        Ok(data) -> Ok(data)
        Fail(_) -> {
            // Primary failed or timed out, try fallback
            fallback = async(fun() -> String { "fallback data" })
            await(fallback)
        }
    }
}

print(fetchWithFallback())  // Ok("primary data")
```

### Pipeline с комбинаторами

```rust
import "lib/task" (taskResolve, taskMap, taskFlatMap, taskCatch, async, await)

// Pipeline with combinators
t1 = taskResolve(1)
t2 = taskMap(t1, fun(x) -> Int { x + 1 })      // 2
t3 = taskMap(t2, fun(x) -> Int { x * 3 })      // 6
t4 = taskFlatMap(t3, fun(x) { async(fun() -> Int { x + 10 }) })  // 16
t5 = taskCatch(t4, fun(err) -> Int { 0 })       // fallback
result = await(t5)

match result {
    Ok(v) -> print("Result: ${v}")  // 16
    Fail(e) -> print("Error: ${e}")
}
```

## Важные замечания

1. **await возвращает Result** — всегда проверяйте результат через match или оператор ??

2. **panic ловится** — если функция внутри async вызывает panic, это будет Fail в Result

3. **taskCancel** — отмена работает только для задач, которые ещё не начали выполняться

4. **Глобальный пул** — защищает от создания слишком большого количества горутин (taskSetGlobalPool/taskGetGlobalPool)

5. **Комбинаторы не блокируют** — taskMap, taskFlatMap, taskCatch возвращают новый Task сразу

