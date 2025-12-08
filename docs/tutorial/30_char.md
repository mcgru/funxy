# Символы (lib/char)

Модуль `lib/char` предоставляет функции для работы с отдельными символами (Char).

```rust
import "lib/char" (*)
```

## Конвертация

### charToCode

```rust
charToCode(c: Char) -> Int
```

Возвращает Unicode код символа.

```rust
charToCode('A')   // 65
charToCode('a')   // 97
charToCode('0')   // 48
charToCode(' ')   // 32
charToCode('\n')  // 10
charToCode('\t')  // 9
charToCode('\0')  // 0
charToCode('\\')  // 92
```

### charFromCode

```rust
charFromCode(code: Int) -> Char
```

Создаёт символ из Unicode кода.

```rust
charFromCode(65)     // 'A'
charFromCode(97)     // 'a'
charFromCode(1071)   // 'Я' (Cyrillic)
charFromCode(12354)  // 'あ' (Japanese hiragana)
```

## Классификация

### charIsUpper

```rust
charIsUpper(c: Char) -> Bool
```

Проверяет, является ли символ заглавной буквой.

```rust
charIsUpper('A')  // true
charIsUpper('Z')  // true
charIsUpper('a')  // false
charIsUpper('5')  // false

// Работает с Unicode
charIsUpper(charFromCode(1071))  // true ('Я')
```

### charIsLower

```rust
charIsLower(c: Char) -> Bool
```

Проверяет, является ли символ строчной буквой.

```rust
charIsLower('a')  // true
charIsLower('z')  // true
charIsLower('A')  // false
charIsLower('5')  // false

// Работает с Unicode
charIsLower(charFromCode(1103))  // true ('я')
```

## Преобразование регистра

### charToUpper

```rust
charToUpper(c: Char) -> Char
```

Преобразует символ в верхний регистр.

```rust
charToUpper('a')  // 'A'
charToUpper('z')  // 'Z'
charToUpper('A')  // 'A' (без изменений)
charToUpper('5')  // '5' (без изменений)

// Unicode
charToUpper(charFromCode(1103))  // 'я' -> 'Я'
```

### charToLower

```rust
charToLower(c: Char) -> Char
```

Преобразует символ в нижний регистр.

```rust
charToLower('A')  // 'a'
charToLower('Z')  // 'z'
charToLower('a')  // 'a' (без изменений)
charToLower('5')  // '5' (без изменений)

// Unicode
charToLower(charFromCode(1071))  // 'Я' -> 'я'
```

## Практические примеры

### Проверка первой буквы

```rust
import "lib/char" (charIsUpper)

fun startsWithUpper(s: String) -> Bool {
    if len(s) == 0 { false }
    else { charIsUpper(s[0]) }
}

print(startsWithUpper("Hello"))  // true
print(startsWithUpper("world"))  // false
```

### Capitalize

```rust
import "lib/char" (charToUpper)
import "lib/list" (tail)

fun capitalize(s: String) -> String {
    if len(s) == 0 { "" }
    else { charToUpper(s[0]) :: tail(s) }
}

print(capitalize("hello"))  // "Hello"
```

### Проверка ASCII

```rust
import "lib/char" (charToCode)

fun isAscii(c: Char) -> Bool {
    charToCode(c) < 128
}

fun isDigit(c: Char) -> Bool {
    code = charToCode(c)
    code >= 48 && code <= 57
}

fun isLetter(c: Char) -> Bool {
    code = charToCode(c)
    (code >= 65 && code <= 90) || (code >= 97 && code <= 122)
}

isDigit('5')   // true
isLetter('a')  // true
```

### ROT13

```rust
import "lib/char" (charToCode, charFromCode, charIsUpper, charIsLower)
import "lib/list" (map)

fun rot13char(c: Char) -> Char {
    code = charToCode(c)
    if charIsUpper(c) {
        charFromCode(65 + (code - 65 + 13) % 26)
    } else if charIsLower(c) {
        charFromCode(97 + (code - 97 + 13) % 26)
    } else {
        c
    }
}

fun rot13(s: String) -> String {
    map(rot13char, s)
}

print(rot13("Hello"))  // "Uryyb"
print(rot13("Uryyb"))  // "Hello"
```

## Работа со строками

Поскольку `String = List<Char>`, функции lib/char хорошо сочетаются с lib/list:

```rust
import "lib/char" (charToUpper, charIsLower)
import "lib/list" (map, filter, length)

s = "Hello World"

// Все в верхний регистр (без lib/string)
upper = map(charToUpper, s)  // "HELLO WORLD"

// Только строчные
lowers = filter(charIsLower, s)  // "elloorld"

// Подсчёт строчных
count = length(filter(charIsLower, s))  // 8
```

## Сводка

| Функция | Тип | Описание |
|---------|-----|----------|
| `charToCode` | `(Char) -> Int` | Символ → Unicode код |
| `charFromCode` | `(Int) -> Char` | Unicode код → символ |
| `charIsUpper` | `(Char) -> Bool` | Заглавная буква? |
| `charIsLower` | `(Char) -> Bool` | Строчная буква? |
| `charToUpper` | `(Char) -> Char` | В верхний регистр |
| `charToLower` | `(Char) -> Char` | В нижний регистр |

## Escape-последовательности

В char литералах поддерживаются:

| Escape | Значение |
|--------|----------|
| `'\n'` | Перенос строки (10) |
| `'\t'` | Табуляция (9) |
| `'\r'` | Возврат каретки (13) |
| `'\0'` | Null (0) |
| `'\\'` | Обратный слэш (92) |
| `'\''` | Одинарная кавычка (39) |

