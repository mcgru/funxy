# Системные функции (lib/sys)

Модуль `lib/sys` предоставляет доступ к системным функциям: аргументы командной строки, переменные окружения, завершение программы и выполнение внешних команд.

```rust
import "lib/sys" (*)
```

## Функции

### sysArgs

```rust
sysArgs() -> List<String>
```

Возвращает аргументы командной строки (без имени программы).

```rust
import "lib/sys" (sysArgs)

// Запуск: ./funxy script.lang hello world
arguments = sysArgs()
// arguments = ["script.lang", "hello", "world"]

print(len(arguments))  // количество аргументов
```

### sysEnv

```rust
sysEnv(name: String) -> Option<String>
```

Возвращает значение переменной окружения или `Zero` если не задана.

```rust
import "lib/sys" (sysEnv)

// PATH существует на всех системах
pathOpt = sysEnv("PATH")
hasPath = match pathOpt {
    Some(_) -> true
    Zero -> false
}
print(hasPath)  // true

// Несуществующая переменная
noVar = sysEnv("NONEXISTENT_VAR")
match noVar {
    Some(val) -> print("Value: " ++ val)
    Zero -> print("Not set")
}
// Not set
```

### sysExit

```rust
sysExit(code: Int) -> Nil
```

Завершает программу с указанным кодом возврата.

- `0` — успешное завершение
- `1` и выше — ошибка

```rust
import "lib/sys" (sysExit)

// Успешное завершение
sysExit(0)

// Завершение с ошибкой
sysExit(1)
```

**Важно:** код после `sysExit()` не выполняется.

### sysExec

```rust
sysExec(cmd: String, args: List<String>) -> { code: Int, stdout: String, stderr: String }
```

Выполняет внешнюю команду и возвращает результат в виде записи с полями:
- `code` — код возврата (0 = успех, -1 = команда не найдена)
- `stdout` — стандартный вывод
- `stderr` — стандартный вывод ошибок

```rust
import "lib/sys" (sysExec)

// Простая команда
echoResult = sysExec("echo", ["hello"])
print(echoResult.code == 0)  // true
print(len(echoResult.stdout) > 0)  // true

// Команда с несколькими аргументами
lsResult = sysExec("ls", ["-la", "/"])
print(lsResult.code == 0)  // true на Unix системах

// Несуществующая команда
badResult = sysExec("nonexistent_command", [])
print(badResult.code == -1)  // true - команда не найдена

// Проверка типов полей
print(getType(echoResult.code))    // Int
print(getType(echoResult.stdout))  // String
print(getType(echoResult.stderr))  // String
```

## Практические примеры

### CLI с аргументами

```rust
import "lib/sys" (sysArgs)
import "lib/list" (length, headOr, tail)

arguments = sysArgs()

// Проверка количества аргументов
if length(arguments) < 2 {
    print("Usage: program <command>")
} else {
    // Пропускаем имя скрипта, берём команду
    cmd = headOr(tail(arguments), "help")
    match cmd {
        "help" -> print("Commands: help, version")
        "version" -> print("v1.0.0")
        _ -> print("Unknown: " ++ cmd)
    }
}
```

### Проверка переменной окружения

```rust
import "lib/sys" (sysEnv)

// Получить порт из окружения или использовать 8080
port = match sysEnv("PORT") {
    Some(p) -> p
    Zero -> "8080"
}
print("Port: " ++ port)

// Проверка debug режима
isDebug = match sysEnv("DEBUG") {
    Some(_) -> true
    Zero -> false
}
print(isDebug)
```

### Получение вывода команды

```rust
import "lib/sys" (sysExec)

result = sysExec("git", ["status", "--short"])
if result.code == 0 {
    print("Git status:")
    print(result.stdout)
} else {
    print("Git error: " ++ result.stderr)
}
```

### Пример: простой grep

```rust
import "lib/sys" (sysArgs)
import "lib/io" (fileRead)
import "lib/list" (length, headOr, tail)
import "lib/string" (stringLines)
import "lib/regex" (regexMatch)

arguments = sysArgs()

if length(arguments) < 3 {
    print("Usage: grep <pattern> <file>")
} else {
    // arguments[0] = скрипт, [1] = pattern, [2] = file
    pattern = headOr(tail(arguments), "")
    filename = headOr(tail(tail(arguments)), "")

    match fileRead(filename) {
        Ok(content) -> {
            allLines = stringLines(content)
            // Простой поиск с regex
            for line in allLines {
                if regexMatch(pattern, line) {
                    print(line)
                }
            }
        }
        Fail(err) -> print("Error: " ++ err)
    }
}
```

## Сводка

| Функция | Тип | Описание |
|---------|-----|----------|
| `sysArgs` | `() -> List<String>` | Аргументы командной строки |
| `sysEnv` | `(String) -> Option<String>` | Переменная окружения |
| `sysExit` | `(Int) -> Nil` | Завершить программу |
| `sysExec` | `(String, List<String>) -> { code: Int, stdout: String, stderr: String }` | Выполнить команду |
