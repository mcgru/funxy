# Строки

## Представление строк

Строки в нашем языке представлены как `List<Char>` — список символов. Это позволяет использовать все операции со списками для работы со строками.

```rust
s = "hello"        // Тип: String (алиас для List<Char>)
print(s[0])        // 'h' - доступ к символу
print(len(s))      // 5
```

## Типы строковых литералов

### Обычные строки

```rust
s = "Hello, World!"
```

### Raw строки (многострочные)

Используют обратные кавычки и сохраняют форматирование:

```rust
text = `This is
a multi-line
string`
```

### Интерполированные строки

Позволяют встраивать выражения прямо в строку с помощью `${...}`:

```rust
name = "Alice"
age = 30

// Простая интерполяция
print("Hello, ${name}!")  // Hello, Alice!

// Выражения в интерполяции
x = 5
y = 3
print("${x} + ${y} = ${x + y}")  // 5 + 3 = 8

// Доступ к полям
person = { name: "Bob", age: 25 }
print("${person.name} is ${person.age}")  // Bob is 25

// Вызов функций
fun double(n) { n * 2 }
print("Double: ${double(10)}")  // Double: 20
```

### Что можно использовать в `${...}`

- Переменные: `${name}`
- Арифметические выражения: `${a + b * 2}`
- Доступ к полям: `${person.name}`
- Индексация: `${list[0]}`
- Вызовы функций: `${func(arg)}`
- Вложенные строки: `${"inner"}`

## Конкатенация строк

Оператор `++` (предпочтительно):

```rust
greeting = "Hello"
name = "World"
message = greeting ++ ", " ++ name ++ "!"
print(message)  // Hello, World!
```

Но интерполяция обычно удобнее и читабельнее.

## Unicode

Строки корректно работают с Unicode:

```rust
import "lib/list" (reverse)

s = "Привет"
print(len(s))      // 6 (символов, не байтов)
print(reverse(s))  // "тевирП"

for c in "日本語" {
    print(c)       // Печатает каждый символ
}
```

## lib/list функции на строках

Поскольку `String = List<Char>`, все функции списков работают:

```rust
import "lib/list" (head, last, tail, init, take, drop, reverse, filter, foldl, find)

s = "hello"

// Длина
print(len(s))                     // 5

// Доступ
print(head(s))                    // Some('h')
print(last(s))                    // Some('o')
print(tail(s))                    // "ello"
print(init(s))                    // "hell"

// Срезы
print(take(s, 3))                 // "hel"
print(drop(s, 2))                 // "llo"

// Поиск
print(find(fun(c) -> c == 'e', s)) // Some('e')

// Трансформация
print(reverse(s))                 // "olleh"
print(filter(fun(c) -> c != 'l', s)) // "heo"

// Свёртка
print(foldl(fun(acc, c) -> acc + 1, 0, s))  // 5
```

## Модуль lib/string

Модуль `lib/string` предоставляет специализированные строковые операции:

```rust
import "lib/string" (*)
```

### Split и Join

```rust
import "lib/string" (stringSplit, stringJoin, stringLines, stringWords)

// stringSplit: (String, String) -> List<String>
print(stringSplit("a,b,c", ","))           // ["a", "b", "c"]
print(stringSplit("one::two", "::"))       // ["one", "two"]

// stringJoin: (List<String>, String) -> String
print(stringJoin(["a", "b", "c"], ","))    // "a,b,c"
print(stringJoin(["a", "b", "c"], ""))     // "abc"

// stringWords: (String) -> List<String>
// Разбивает по пробелам
print(stringWords("hello   world  test"))  // ["hello", "world", "test"]
```

### Обрезка пробелов

```rust
// stringTrim: (String) -> String
stringTrim("  hello  ")             // "hello"

// stringTrimStart: (String) -> String
stringTrimStart("  hello")          // "hello"

// stringTrimEnd: (String) -> String
stringTrimEnd("hello  ")            // "hello"
```

### Регистр

```rust
// stringToUpper: (String) -> String
stringToUpper("hello")              // "HELLO"
stringToUpper("Привет")             // "ПРИВЕТ"

// stringToLower: (String) -> String
stringToLower("HELLO")              // "hello"

// stringCapitalize: (String) -> String
stringCapitalize("hello")           // "Hello"
```

### Поиск и замена

```rust
// stringReplace: (String, String, String) -> String
// Заменяет первое вхождение
stringReplace("hello", "l", "L")    // "heLlo"

// stringReplaceAll: (String, String, String) -> String
// Заменяет все вхождения
stringReplaceAll("hello", "l", "L") // "heLLo"
stringReplaceAll("banana", "a", "o") // "bonono"

// stringStartsWith: (String, String) -> Bool
stringStartsWith("hello", "hel")    // true
stringStartsWith("hello", "ell")    // false

// stringEndsWith: (String, String) -> Bool
stringEndsWith("hello", "lo")       // true
stringEndsWith("hello", "ll")       // false

// stringIndexOf: (String, String) -> Option<Int>
stringIndexOf("hello", "ll")      // Some(2)
stringIndexOf("hello", "x")       // Zero
```

### Повтор и выравнивание

```rust
import "lib/string" (stringRepeat, stringPadLeft, stringPadRight)

// stringRepeat: (String, Int) -> String
print(stringRepeat("ab", 3))               // "ababab"
print(stringRepeat("x", 5))                // "xxxxx"

// stringPadLeft: (String, Int, Char) -> String
print(stringPadLeft("42", 5, '0'))         // "00042"
print(stringPadLeft("hello", 3, '-'))      // "hello" (без изменений если >= длины)

// stringPadRight: (String, Int, Char) -> String
print(stringPadRight("42", 5, '-'))        // "42---"
```

## Pipe оператор

Строковые функции отлично работают с pipe:

```rust
import "lib/string" (*)

result = "  HELLO WORLD  " 
    |> stringTrim 
    |> stringToLower 
    |> stringCapitalize
print(result)  // "Hello world"

// Обработка CSV
csv = "a,b,c"
parts = stringSplit(csv, ",")
formatted = stringJoin(parts, " | ")
print(formatted)  // "a | b | c"
```

## Практические примеры

### Подсчёт слов

```rust
import "lib/string" (stringWords)

fun wordCount(text: String) -> Int {
    len(stringWords(text))
}

print(wordCount("hello world test"))  // 3
```

### Title Case

```rust
import "lib/string" (stringWords, stringJoin, stringCapitalize)
import "lib/list" (map)

fun titleCase(s: String) -> String {
    words = stringWords(s)
    capitalized = map(stringCapitalize, words)
    stringJoin(capitalized, " ")
}

print(titleCase("hello world"))  // "Hello World"
```

### Форматирование чисел

```rust
import "lib/string" (stringPadLeft)

fun formatNumber(n: Int, width: Int) -> String {
    stringPadLeft(show(n), width, '0')
}

print(formatNumber(42, 5))  // "00042"
```

### Парсинг CSV

```rust
import "lib/string" (stringSplit, stringTrim)
import "lib/list" (map)

parseCSVLine = fun(line: String) -> List<String> {
    stringSplit(line, ",") |> map(stringTrim)
}

parseCSVLine("  a , b , c  ") // ["a", "b", "c"]
```

### Отладочный вывод

```rust
fun debug(label, value) {
    print("[DEBUG] ${label} = ${value}")
}

x = 100
debug("x", x)  // [DEBUG] x = 100
```

## Сводка lib/string

| Функция | Тип | Описание |
|---------|-----|----------|
| `stringSplit` | `(String, String) -> List<String>` | Разбить по разделителю |
| `stringJoin` | `(List<String>, String) -> String` | Соединить с разделителем |
| `stringLines` | `(String) -> List<String>` | Разбить по строкам |
| `stringWords` | `(String) -> List<String>` | Разбить по пробелам |
| `stringTrim` | `(String) -> String` | Обрезать пробелы |
| `stringTrimStart` | `(String) -> String` | Обрезать слева |
| `stringTrimEnd` | `(String) -> String` | Обрезать справа |
| `stringToUpper` | `(String) -> String` | В верхний регистр |
| `stringToLower` | `(String) -> String` | В нижний регистр |
| `stringCapitalize` | `(String) -> String` | Первую букву заглавной |
| `stringReplace` | `(String, String, String) -> String` | Заменить первое |
| `stringReplaceAll` | `(String, String, String) -> String` | Заменить все |
| `stringStartsWith` | `(String, String) -> Bool` | Проверка префикса |
| `stringEndsWith` | `(String, String) -> Bool` | Проверка суффикса |
| `stringIndexOf` | `(String, String) -> Option<Int>` | Найти подстроку |
| `stringRepeat` | `(String, Int) -> String` | Повторить N раз |
| `stringPadLeft` | `(String, Int, Char) -> String` | Дополнить слева |
| `stringPadRight` | `(String, Int, Char) -> String` | Дополнить справа |
