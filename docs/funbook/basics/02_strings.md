# 02. Работа со строками

## Задача
Манипулировать текстовыми данными: разбивать, объединять, искать, форматировать.

## Основы

```rust
// Строки — это List<Char>
s = "Hello, World!"

// Длина
print(len(s))  // 13

// Доступ по индексу
print(s[0])    // 'H'
print(s[7])    // 'W'

// Конкатенация
full = "Hello" ++ ", " ++ "World!"
```

## Интерполяция

```rust
name = "Alice"
age = 30

// ${expression} внутри строки
print("Name: ${name}, Age: ${age}")
// Name: Alice, Age: 30

// Вычисления внутри
print("In 10 years: ${age + 10}")
// In 10 years: 40
```

## Многострочные строки

```rust
poem = "Roses are red,
Violets are blue,
Funxy is awesome,
And so are you!"

print(poem)
```

## Escape-последовательности

```rust
print("Line 1\nLine 2")     // перенос строки
print("Tab\there")          // табуляция
print("Quote: \"Hi\"")      // кавычки
print("Backslash: \\")      // обратный слэш
print("Dollar: \$")         // доллар (для интерполяции)
print("Null: \0")           // null byte
```

## lib/string — мощные функции

```rust
import "lib/string" (*)

s = "  Hello, World!  "

// Trimming
print(stringTrim(s))       // "Hello, World!"
print(stringTrimStart(s))  // "Hello, World!  "
print(stringTrimEnd(s))    // "  Hello, World!"

// Регистр
print(stringToUpper("hello"))      // "HELLO"
print(stringToLower("HELLO"))      // "hello"
print(stringCapitalize("hello"))   // "Hello"

// Split/Join
parts = stringSplit("a,b,c", ",")  // ["a", "b", "c"]
joined = stringJoin(parts, "-")    // "a-b-c"

// Удобные split
lines = stringLines("a\nb\nc")     // ["a", "b", "c"]
words = stringWords("hello world") // ["hello", "world"]

// Поиск
print(stringStartsWith("hello", "he"))   // true
print(stringEndsWith("hello", "lo"))     // true
print(stringIndexOf("hello", "ll"))      // Some(2)

// Замена
print(stringReplace("hello", "l", "L"))     // "heLlo" (первое)
print(stringReplaceAll("hello", "l", "L"))  // "heLLo" (все)

// Повтор и паддинг
print(stringRepeat("ab", 3))              // "ababab"
print(stringPadLeft("42", 5, '0'))        // "00042"
print(stringPadRight("hi", 5, ' '))       // "hi   "
```

## lib/char — работа с символами

```rust
import "lib/char" (*)

// Код символа
print(charToCode('A'))  // 65
print(charFromCode(65)) // 'A'

// Проверки
print(charIsUpper('A'))  // true
print(charIsLower('a'))  // true

// Преобразование
print(charToUpper('a'))  // 'A'
print(charToLower('Z'))  // 'z'
```

## Конвертация типов

```rust
// В строку
print(show(42))       // "42"
print(show(3.14))     // "3.14"
print(show(true))     // "true"
print(show([1,2,3]))  // "[1, 2, 3]"

// Из строки (возвращает Option)
print(read("42", Int))        // Some(42)
print(read("3.14", Float))    // Some(3.14)
print(read("abc", Int))       // Zero
```

## Строки как списки

```rust
import "lib/list" (map, filter, foldl, reverse, take, drop, contains)
import "lib/char" (charToUpper, charToCode)

// Строка — это List<Char>, поэтому работают все функции списков!

s = "hello"

// map
upper = map(charToUpper, s)
print(upper)  // "HELLO"

// filter
noVowels = filter(fun(c) -> !contains("aeiou", c), s)
print(noVowels)

// foldl
codes = foldl(fun(acc, c) -> acc + charToCode(c), 0, s)
print(codes)

// reverse
rev = reverse(s)
print(rev)  // "olleh"

// take/drop
print(take(s, 3))  // "hel"
print(drop(s, 2))  // "llo"
```

## Форматирование чисел

```rust
import "lib/string" (stringPadLeft)

fun formatPrice(cents: Int) -> String {
    dollars = cents / 100
    remainder = cents % 100
    "$" ++ show(dollars) ++ "." ++ stringPadLeft(show(remainder), 2, '0')
}

print(formatPrice(1299))  // "$12.99"
print(formatPrice(500))   // "$5.00"
```

## Парсинг данных

```rust
import "lib/string" (stringSplit, stringTrim)
import "lib/list" (map)

// Парсинг CSV строки
csv = "Alice, 30, Engineer"
parts = stringSplit(csv, ",") |> map(fun(s) -> stringTrim(s))
print(parts)  // ["Alice", "30", "Engineer"]

// Парсинг key=value
fun parseKV(s: String) -> (String, String) {
    parts = stringSplit(s, "=")
    (parts[0], parts[1])
}

(key, value) = parseKV("name=Alice")
print(key)    // name
print(value)  // Alice
```

## Regex (lib/regex)

```rust
import "lib/regex" (*)

text = "Email: alice@example.com, Phone: 123-456-7890"

// Проверка паттерна
print(regexMatch("[a-z]+@[a-z]+\\.[a-z]+", text))  // true

// Найти первое совпадение
match regexFind("[0-9-]+", text) {
    Some phone -> print(phone)  // "123-456-7890"
    Zero -> print("Not found")
}

// Найти все совпадения
numbers = regexFindAll("[0-9]+", text)  // ["123", "456", "7890"]

// Capture groups
match regexCapture("([a-z]+)@([a-z]+)", text) {
    Some groups -> {
        print(groups[1])  // "alice"
        print(groups[2])  // "example"
    }
    Zero -> print("No match")
}

// Замена
censored = regexReplaceAll("[0-9]", "X", text)
print(censored)
// "Email: alice@example.com, Phone: XXX-XXX-XXXX"

// Split по паттерну
parts = regexSplit("[,;]\\s*", "a, b; c")  // ["a", "b", "c"]
```

## Практический пример: парсер URL query string

```rust
import "lib/string" (stringSplit)
import "lib/list" (foldl)
import "lib/map" (mapPut, mapGet)

fun parseQuery(query: String) -> Map<String, String> {
    pairs = stringSplit(query, "&")
    foldl(fun(acc, pair) {
        parts = stringSplit(pair, "=")
        key = parts[0]
        value = if len(parts) > 1 { parts[1] } else { "" }
        mapPut(acc, key, value)
    }, %{}, pairs)
}

params = parseQuery("name=Alice&age=30")
match mapGet(params, "name") {
    Some(n) -> print(n)  // "Alice"
    Zero -> print("not found")
}
```

## Практический пример: slug генератор

```rust
import "lib/string" (stringToLower, stringTrim)
import "lib/regex" (regexReplaceAll)

fun slugify(title: String) -> String {
    s1 = stringTrim(title)
    s2 = stringToLower(s1)
    s3 = regexReplaceAll("[^a-z0-9\\s-]", "", s2)  // убрать спецсимволы
    s4 = regexReplaceAll("\\s+", "-", s3)          // пробелы -> дефисы
    regexReplaceAll("-+", "-", s4)                  // убрать дубли дефисов
}

print(slugify("Hello, World! This is Funxy"))
// "hello-world-this-is-funxy"
```

