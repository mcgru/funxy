# Регулярные выражения (lib/regex)

Модуль `lib/regex` предоставляет функции для работы с регулярными выражениями, основанными на синтаксисе RE2 (Go).

```rust
import "lib/regex" (*)
```

## Важно: escape-последовательности

В строках этого языка backslash передаётся как есть:

```rust
import "lib/regex" (regexMatch)

// Правильно:
regexMatch("\d+", "abc123")  // \d = цифра

// Неправильно (двойной backslash):
regexMatch("\\d+", "abc123")  // \\d = литерал \d
```

## Функции

### matchRe

```rust
regexMatch(pattern: String, str: String) -> Bool
```

Проверяет, совпадает ли паттерн где-либо в строке.

```rust
import "lib/regex" (regexMatch)

print(regexMatch("\d+", "abc123"))     // true
print(regexMatch("\d+", "abc"))        // false
print(regexMatch("^hello", "hello!"))  // true (начинается с "hello")
print(regexMatch("world$", "hello world"))  // true (заканчивается на "world")
```

### findRe

```rust
regexFind(pattern: String, str: String) -> Option<String>
```

Находит первое совпадение паттерна.

```rust
import "lib/regex" (regexFind)

result = regexFind("\d+", "abc123def456")
print(result)  // Some("123")

noMatch = regexFind("\d+", "abc")
print(noMatch)  // Zero
```

### findAllRe

```rust
regexFindAll(pattern: String, str: String) -> List<String>
```

Находит все совпадения паттерна.

```rust
import "lib/regex" (regexFindAll)

matches = regexFindAll("\d+", "a1b22c333")
print(matches)  // ["1", "22", "333"]

noMatches = regexFindAll("\d+", "abc")
print(noMatches)  // []
```

### capture

```rust
regexCapture(pattern: String, str: String) -> Option<List<String>>
```

Извлекает группы захвата из первого совпадения.

- Индекс 0: полное совпадение
- Индекс 1+: группы захвата `(...)` в порядке появления

```rust
import "lib/regex" (regexCapture)

// Паттерн с группами: (год)-(месяц)-(день)
result = regexCapture("(\d{4})-(\d{2})-(\d{2})", "Date: 2024-03-15")

match result {
    Some(groups) -> {
        print(groups[0])  // "2024-03-15" (полное совпадение)
        print(groups[1])  // "2024" (год)
        print(groups[2])  // "03" (месяц)
        print(groups[3])  // "15" (день)
    }
    Zero -> print("Нет совпадения")
}
```

### replaceRe

```rust
regexReplace(pattern: String, replacement: String, str: String) -> String
```

Заменяет **первое** совпадение паттерна.

```rust
import "lib/regex" (regexReplace)

result = regexReplace("\d+", "X", "a1b2c3")
print(result)  // "aXb2c3"
```

### replaceAllRe

```rust
regexReplaceAll(pattern: String, replacement: String, str: String) -> String
```

Заменяет **все** совпадения паттерна.

```rust
import "lib/regex" (regexReplaceAll)

// Замена всех чисел
result = regexReplaceAll("\d+", "X", "a1b2c3")
print(result)  // "aXbXcX"

// Обратные ссылки на группы ($1, $2, ...)
wrapped = regexReplaceAll("(\w+)", "[$1]", "hello world")
print(wrapped)  // "[hello] [world]"
```

### splitRe

```rust
regexSplit(pattern: String, str: String) -> List<String>
```

Разбивает строку по паттерну.

```rust
import "lib/regex" (regexSplit)

// По запятой и пробелам
parts = regexSplit(",\s*", "a, b,c,  d")
print(parts)  // ["a", "b", "c", "d"]

// По пробелам
words = regexSplit("\s+", "hello   world   foo")
print(words)  // ["hello", "world", "foo"]
```

### validateRe

```rust
regexValidate(pattern: String) -> Result<Nil, String>
```

Проверяет синтаксис регулярного выражения.

```rust
import "lib/regex" (regexValidate)

match regexValidate("\d+") {
    Ok(_) -> print("Valid")
    Fail(err) -> print("Error: ${err}")
}  // "Valid"

match regexValidate("[invalid") {
    Ok(_) -> print("Valid")
    Fail(err) -> print("Error: ${err}")
}  // "Error: ..."
```

## Синтаксис паттернов

Поддерживается синтаксис RE2 (Go):

| Паттерн | Описание |
|---------|----------|
| `.` | Любой символ (кроме newline) |
| `\d` | Цифра [0-9] |
| `\D` | Не цифра |
| `\w` | Слово [a-zA-Z0-9_] |
| `\W` | Не слово |
| `\s` | Пробельный символ |
| `\S` | Не пробельный |
| `^` | Начало строки |
| `$` | Конец строки |
| `*` | 0 или более |
| `+` | 1 или более |
| `?` | 0 или 1 |
| `{n}` | Ровно n раз |
| `{n,m}` | От n до m раз |
| `[abc]` | Любой из символов |
| `[^abc]` | Любой кроме |
| `(...)` | Группа захвата |
| `(?:...)` | Группа без захвата |
| `a\|b` | a или b |

## Практические примеры

### Валидация email

```rust
import "lib/regex" (regexMatch)

fun isValidEmail(email: String) -> Bool {
    pattern = "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}"
    regexMatch(pattern, email)
}

print(isValidEmail("test@example.com"))  // true
print(isValidEmail("invalid"))           // false
```

### Извлечение данных

```rust
import "lib/regex" (regexCapture, regexFindAll)

// Извлечь URL из текста
urls = regexFindAll("https?://[^\s]+", "Visit https://example.com or http://test.org")
print(urls)  // ["https://example.com", "http://test.org"]

// Извлечь части URL
result = regexCapture("(https?)://([^/]+)(/.*)?", "https://example.com/path")
match result {
    Some(parts) -> {
        print(parts[1])  // "https" (протокол)
        print(parts[2])  // "example.com" (хост)
        print(parts[3])  // "/path" (путь)
    }
    Zero -> ()
}
```

### Очистка данных

```rust
import "lib/regex" (regexReplaceAll, regexSplit)

// Удалить HTML теги
clean = regexReplaceAll("<[^>]+>", "", "<p>Hello <b>world</b></p>")
print(clean)  // "Hello world"

// Нормализация пробелов
normalized = regexReplaceAll("\s+", " ", "hello    world\nfoo")
print(normalized)  // "hello world foo"

// Разбор CSV (простой случай)
values = regexSplit(",", "a,b,c,d")
print(values)  // ["a", "b", "c", "d"]
```

## Сводка

| Функция | Тип | Описание |
|---------|-----|----------|
| `matchRe` | `(String, String) -> Bool` | Проверка совпадения |
| `findRe` | `(String, String) -> Option<String>` | Первое совпадение |
| `findAllRe` | `(String, String) -> List<String>` | Все совпадения |
| `capture` | `(String, String) -> Option<List<String>>` | Группы захвата |
| `replaceRe` | `(String, String, String) -> String` | Замена первого |
| `replaceAllRe` | `(String, String, String) -> String` | Замена всех |
| `splitRe` | `(String, String) -> List<String>` | Разбиение |
| `validateRe` | `(String) -> Result<Nil, String>` | Валидация паттерна |

## Ограничения

- Не поддерживаются lookbehind assertions (ограничение RE2)
- Backreference в паттерне (`\1`) не поддерживается
- Backreference в replacement (`$1`) поддерживается

