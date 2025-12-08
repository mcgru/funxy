# Время (lib/time)

Модуль `lib/time` предоставляет функции для работы со временем.

```rust
import "lib/time" (*)
```

## Функции

### time

```rust
timeNow() -> Int
```

Возвращает Unix timestamp в секундах (время с 1 января 1970).

```rust
import "lib/time" (timeNow)

now = timeNow()
print(now)  // 1701234567
```

### clockNs и clockMs

```rust
clockNs() -> Int   // наносекунды
clockMs() -> Int   // миллисекунды
```

Монотонные часы для измерения времени. Не зависят от изменения системного времени.

```rust
import "lib/time" (clockMs)

start = clockMs()
// ... код ...
end = clockMs()
elapsed = end - start
print("Elapsed: ${elapsed} ms")
```

### sleep и sleepMs

```rust
sleep(seconds: Int) -> Nil
sleepMs(ms: Int) -> Nil
```

Приостанавливает выполнение программы.

```rust
import "lib/time" (sleep, sleepMs)

sleepMs(100)  // пауза 100 мс
sleep(1)      // пауза 1 секунда
```

## Практические примеры

### Бенчмарк функции

```rust
import "lib/time" (clockMs)

fun benchmark(f: () -> Nil, iterations: Int) -> Int {
    start = clockMs()
    for i in range(0, iterations) {
        f()
    }
    end = clockMs()
    end - start
}

fun heavyWork() {
    // ... какая-то работа ...
}

elapsed = benchmark(heavyWork, 1000)
print("1000 iterations: ${elapsed} ms")
```

### Таймаут

```rust
import "lib/time" (clockMs, sleepMs)

fun waitFor(condition: () -> Bool, timeoutMs: Int) -> Bool {
    start = clockMs()
    for !condition() {
        if clockMs() - start > timeoutMs {
            break false
        }
        sleepMs(10)
        true
    }
}
```

### Простой профайлер

```rust
import "lib/time" (clockNs)

fun timed<T>(label: String, f: () -> T) -> T {
    start = clockNs()
    result = f()
    elapsed = clockNs() - start
    print("${label}: ${elapsed / 1000000} ms")
    result
}

// Использование
result = timed("computation", fun() -> {
    // ... тяжёлые вычисления ...
    42
})
```

## Разница между timeNow() и clockMs()/clockNs()

| Функция | Тип | Изменяется при смене системного времени |
|---------|-----|----------------------------------------|
| `timeNow()` | Wall clock | Да |
| `clockMs()` | Monotonic | Нет |
| `clockNs()` | Monotonic | Нет |

**Правило**: используй `clockMs()`/`clockNs()` для измерения интервалов, `timeNow()` для работы с датами.

## Сводка

| Функция | Тип | Описание |
|---------|-----|----------|
| `time` | `() -> Int` | Unix timestamp (секунды) |
| `clockNs` | `() -> Int` | Монотонные наносекунды |
| `clockMs` | `() -> Int` | Монотонные миллисекунды |
| `sleep` | `(Int) -> Nil` | Пауза в секундах |
| `sleepMs` | `(Int) -> Nil` | Пауза в миллисекундах |

