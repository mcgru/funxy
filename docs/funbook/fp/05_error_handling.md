# 05. Error Handling

## Задача
Обрабатывать ошибки и отсутствующие значения типобезопасно, без exceptions.

---

## Три подхода к ошибкам

| Тип | Когда использовать | Пример |
|-----|-------------------|--------|
| `Option<T>` | Значение может отсутствовать | `find(predicate, list)` |
| `T?` (Nullable) | Быстрая проверка на null | `user.email?` |
| `Result<E, A>` | Нужна информация об ошибке | `fileRead(path)` |

---

## Option<T>: есть или нет

Option — встроенный тип: `Some T | Zero`.

```rust
fun safeDivide(a: Int, b: Int) -> Option<Int> {
    if b == 0 { Zero } else { Some(a / b) }
}

print(safeDivide(10, 2))  // Some(5)
print(safeDivide(10, 0))  // Zero
```

### Pattern Matching с Option

```rust
fun showResult(opt: Option<Int>) -> String {
    match opt {
        Some(x) -> "Got: " ++ show(x)
        Zero -> "Nothing"
    }
}

print(showResult(Some(42)))  // Got: 42
print(showResult(Zero))      // Nothing
```

### Функции для работы с Option

```rust
// map — применить функцию к значению внутри
fun mapOption(opt, f) {
    match opt {
        Some(x) -> Some(f(x))
        Zero -> Zero
    }
}

// flatMap — для цепочек операций
fun flatMapOption(opt, f) {
    match opt {
        Some(x) -> f(x)
        Zero -> Zero
    }
}

// getOrElse — значение по умолчанию
fun getOrElse(opt, default) {
    match opt {
        Some(x) -> x
        Zero -> default
    }
}

// orElse — альтернативный Option
fun orElse(opt, alternative) {
    match opt {
        Some(x) -> Some(x)
        Zero -> alternative
    }
}

// Примеры использования
x = Some(10)
doubled = mapOption(x, fun(n) -> n * 2)
print(doubled)  // Some(20)

value = getOrElse(Zero, 42)
print(value)  // 42

backup = orElse(Zero, Some(100))
print(backup)  // Some(100)
```

---

## T? (Nullable): быстрый null-check

`T?` — это синтаксический сахар для `T | Nil`.

```rust
type User = { name: String, id: Int }

// Функция может вернуть Nil
fun findUser(id: Int) -> User? {
    if id > 0 { { name: "Alice", id: id } } else { Nil }
}

user = findUser(1)

// Pattern matching с типами
match user {
    _: Nil -> print("Not found")
    u: User -> print("Found: " ++ u.name)
}
```

### Nullable поля в записях

```rust
type Profile = {
    name: String,
    email: String?,
    phone: String?,
    age: Int
}

fun showProfile(p: Profile) {
    print("Name: " ++ p.name)
    
    match p.email {
        _: Nil -> print("Email: not provided")
        e: String -> print("Email: " ++ e)
    }
    
    match p.phone {
        _: Nil -> print("Phone: not provided")
        ph: String -> print("Phone: " ++ ph)
    }
}

profile = { name: "Alice", email: "alice@example.com", phone: Nil, age: 30 }
showProfile(profile)
```

---

## Result<E, A>: успех или ошибка с информацией

Result — встроенный тип: `Ok A | Fail E`.

```rust
fun parseNumber(s: String) -> Result<String, Int> {
    match read(s, Int) {
        Some(n) -> Ok(n)
        Zero -> Fail("Cannot parse '" ++ s ++ "' as number")
    }
}

print(parseNumber("42"))     // Ok(42)
print(parseNumber("hello"))  // Fail("Cannot parse 'hello' as number")
```

### Pattern Matching с Result

```rust
fun handleResult(r: Result<String, Int>) -> String {
    match r {
        Ok(value) -> "Success: " ++ show(value)
        Fail(error) -> "Error: " ++ error
    }
}

print(handleResult(Ok(100)))        // Success: 100
print(handleResult(Fail("oops")))   // Error: oops
```

### Функции для работы с Result

```rust
// map — применить функцию к успешному значению
fun mapResult(r, f) {
    match r {
        Ok(a) -> Ok(f(a))
        Fail(e) -> Fail(e)
    }
}

// flatMap — для цепочек
fun flatMapResult(r, f) {
    match r {
        Ok(a) -> f(a)
        Fail(e) -> Fail(e)
    }
}

// mapError — преобразовать ошибку
fun mapError(r, f) {
    match r {
        Ok(a) -> Ok(a)
        Fail(e) -> Fail(f(e))
    }
}

// getOrElse
fun getOrElseResult(r, default) {
    match r {
        Ok(a) -> a
        Fail(_) -> default
    }
}

// Примеры
doubled = mapResult(Ok(5), fun(x) -> x * 2)
print(doubled)  // Ok(10)

withDefault = getOrElseResult(Fail("error"), 0)
print(withDefault)  // 0
```

---

## Реальные примеры из lib/*

### Чтение файла

```rust
import "lib/io" (fileRead, fileWrite)

fun loadConfig(path: String) {
    match fileRead(path) {
        Ok(content) -> {
            print("Loaded " ++ show(len(content)) ++ " bytes")
            content
        }
        Fail(err) -> {
            print("Error reading " ++ path ++ ": " ++ err)
            "{}"  // default empty config
        }
    }
}

config = loadConfig("config.json")

```

### HTTP запросы

```rust
import "lib/http" (httpGet)
import "lib/json" (jsonDecode)

fun fetchUser(id: Int) {
    url = "https://api.example.com/users/" ++ show(id)
    
    match httpGet(url) {
        Ok(response) -> {
            if response.status == 200 {
                Ok(jsonDecode(response.body))
            } else {
                Fail("HTTP " ++ show(response.status))
            }
        }
        Fail(err) -> Fail("Network error: " ++ err)
    }
}

```

---

## Цепочки операций

### Проблема: вложенные match

```rust
// Пример: mock функции
fun getUser(id: Int) -> Result<String, { id: Int, name: String }> {
    if id > 0 { Ok({ id: id, name: "Alice" }) } else { Fail("User not found") }
}

fun getEmail(user) -> Result<String, String> {
    Ok(user.name ++ "@example.com")
}

fun validateEmail(email: String) -> Result<String, String> {
    Ok(email)
}

// Глубокая вложенность - работает, но трудно читать
fun processUserBad(id: Int) {
    match getUser(id) {
        Fail(e) -> Fail(e)
        Ok(user) -> match getEmail(user) {
            Fail(e) -> Fail(e)
            Ok(email) -> match validateEmail(email) {
                Fail(e) -> Fail(e)
                Ok(valid) -> Ok({ user: user, email: valid })
            }
        }
    }
}

result = processUserBad(1)
print(result)
```

### Решение: flatMap chain

```rust
// Вспомогательные функции для Result
fun mapResult(r, f) {
    match r { Ok(v) -> Ok(f(v)), Fail(e) -> Fail(e) }
}

fun flatMapResult(r, f) {
    match r { Ok(v) -> f(v), Fail(e) -> Fail(e) }
}

// Mock функции
fun getUser(id: Int) { if id > 0 { Ok({ id: id, name: "Alice" }) } else { Fail("Not found") } }
fun getEmail(user) { Ok(user.name ++ "@example.com") }
fun validateEmail(email: String) { Ok(email) }

// Линейная цепочка - чище
fun processUserGood(id: Int) {
    result = getUser(id)
    result = flatMapResult(result, fun(user) -> mapResult(getEmail(user), fun(email) -> { user: user, email: email }))
    result = flatMapResult(result, fun(data) -> mapResult(validateEmail(data.email), fun(valid) -> { user: data.user, email: valid }))
    result
}

result = processUserGood(1)
print(result)
```

---

## Кастомные типы ошибок

```rust
import "lib/list" (contains)

// ADT для ошибок валидации
type ValidationError = EmptyField(String)
                     | TooShort(String)
                     | InvalidFormat(String)

// Функции валидации
fun validateUsername(name: String) -> Result<ValidationError, String> {
    if len(name) == 0 { 
        Fail(EmptyField("username")) 
    } else if len(name) < 3 { 
        Fail(TooShort("username")) 
    } else { 
        Ok(name) 
    }
}

fun validateEmail(email: String) -> Result<ValidationError, String> {
    if len(email) == 0 {
        Fail(EmptyField("email"))
    } else if !contains(email, '@') {
        Fail(InvalidFormat("email"))
    } else {
        Ok(email)
    }
}

// Красивый вывод ошибок
fun showValidationError(e: ValidationError) -> String {
    match e {
        EmptyField(field) -> "Field '" ++ field ++ "' is required"
        TooShort(field) -> "Field '" ++ field ++ "' is too short"
        InvalidFormat(field) -> "Field '" ++ field ++ "' has invalid format"
    }
}

// Использование
match validateUsername("ab") {
    Ok(name) -> print("Valid username: " ++ name)
    Fail(err) -> print("Validation error: " ++ showValidationError(err))
}
// Validation error: Field 'username' is too short
```

---

## Комбинирование нескольких Result

```rust
// Если все Ok — вернуть список значений
// Если хоть один Fail — вернуть первую ошибку
fun sequence(results) {
    match results {
        [] -> Ok([])
        [Ok(x), rest...] -> match sequence(rest) {
            Ok(xs) -> Ok([x] ++ xs)
            Fail(e) -> Fail(e)
        }
        [Fail(e), _...] -> Fail(e)
    }
}

// Примеры
allOk = [Ok(1), Ok(2), Ok(3)]
print(sequence(allOk))  // Ok([1, 2, 3])

withError = [Ok(1), Fail("oops"), Ok(3)]
print(sequence(withError))  // Fail("oops")
```

### Собрать все ошибки (не только первую)

```rust
import "lib/list" (foldl)

// Собирает все Ok значения или все Fail ошибки
fun sequenceAll(results) {
    foldl(fun(acc, r) -> {
        match (acc, r) {
            (Ok(xs), Ok(x)) -> Ok(xs ++ [x])
            (Fail(es), Fail(e)) -> Fail(es ++ [e])
            (Fail(es), Ok(_)) -> Fail(es)
            (Ok(_), Fail(e)) -> Fail([e])
        }
    }, Ok([]), results)
}

// Пример
results = [Ok(1), Ok(2), Ok(3)]
match sequenceAll(results) {
    Ok(values) -> print("All valid: " ++ show(values))
    Fail(errors) -> print("Errors: " ++ show(errors))
}

resultsWithErrors = [Ok(1), Fail("error1"), Fail("error2")]
match sequenceAll(resultsWithErrors) {
    Ok(values) -> print("All valid: " ++ show(values))
    Fail(errors) -> print("Errors: " ++ show(errors))
}

```

---

## Преобразование между типами

```rust
// Option -> Result
fun optionToResult(opt, errorMsg) {
    match opt {
        Some(x) -> Ok(x)
        Zero -> Fail(errorMsg)
    }
}

// Result -> Option (теряем информацию об ошибке)
fun resultToOption(res) {
    match res {
        Ok(x) -> Some(x)
        Fail(_) -> Zero
    }
}

// Примеры
print(optionToResult(Some(42), "missing"))  // Ok(42)
print(optionToResult(Zero, "missing"))      // Fail("missing")

print(resultToOption(Ok(42)))       // Some(42)
print(resultToOption(Fail("err")))  // Zero
```

---

## Когда использовать что?

| Ситуация | Используйте | Пример |
|----------|------------|--------|
| Значение может отсутствовать, причина не важна | `Option<T>` | `find()`, `head()` |
| Быстрый null-check для полей | `T?` | `user.middleName?` |
| Операция может не удаться, нужна причина | `Result<String, T>` | `fileRead()`, `httpGet()` |
| Валидация с детальными ошибками | `Result<CustomError, T>` | Form validation |
| Может быть несколько ошибок | `Result<List<E>, T>` | Batch validation |
