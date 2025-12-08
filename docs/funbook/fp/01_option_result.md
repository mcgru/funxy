# 01. Option и Result

## Задача
Обрабатывать отсутствующие значения и ошибки без null/exceptions.

## Option: может быть значение, а может не быть

Option — встроенный тип: `Some T | Zero`.

```rust
fun safeDivide(a: Int, b: Int) -> Option<Int> {
    if b == 0 { Zero } else { Some(a / b) }
}

print(safeDivide(10, 2))  // Some(5)
print(safeDivide(10, 0))  // Zero
```

## Обработка Option

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

## Быстрая проверка: `T?` синтаксис

```rust
// T? эквивалентно T | Nil
type User = { name: String, id: Int }

fun findUser(id: Int) -> User? {
    if id > 0 { { name: "Alice", id: id } } else { Nil }
}

user = findUser(1)
match user {
    _: Nil -> print("Not found")
    u: User -> print("Found: " ++ u.name)
}
```

## Result: успех или ошибка с информацией

Result — встроенный тип: `Ok A | Fail E`.

```rust
fun parseNumber(s: String) -> Result<String, Int> {
    match read(s, Int) {
        Some(n) -> Ok(n)
        Zero -> Fail("Cannot parse: " ++ s)
    }
}

print(parseNumber("42"))     // Ok(42)
print(parseNumber("hello"))  // Fail("Cannot parse: hello")
```

## Обработка Result

```rust
fun handleResult(r: Result<String, Int>) -> String {
    match r {
        Ok(value) -> "Success: " ++ show(value)
        Fail(error) -> "Error: " ++ error
    }
}

print(handleResult(Ok(100)))          // Success: 100
print(handleResult(Fail("oops")))     // Error: oops
```

## Функции для работы с Option

```rust
// Применить функцию к значению внутри Option
fun mapOption(opt, f) {
    match opt {
        Some(x) -> Some(f(x))
        Zero -> Zero
    }
}

// Значение по умолчанию
fun getOrElse(opt, default) {
    match opt {
        Some(x) -> x
        Zero -> default
    }
}

// Примеры
x = Some(10)
doubled = mapOption(x, fun(n) -> n * 2)  // Some(20)
print(doubled)

value = getOrElse(Zero, 42)  // 42
print(value)
```

## Цепочка операций (flatMap)

```rust
// Цепочка операций с Option
fun flatMapOption(opt, f) {
    match opt {
        Some(x) -> f(x)
        Zero -> Zero
    }
}

// Пример: безопасное деление с цепочкой
fun safeDivide(a: Int, b: Int) -> Option<Int> {
    if b == 0 { Zero } else { Some(a / b) }
}

result = Some(100)
    |> fun(opt) -> flatMapOption(opt, fun(x) -> safeDivide(x, 2))
    |> fun(opt) -> flatMapOption(opt, fun(x) -> safeDivide(x, 5))

print(result)  // Some(10)
```

## Практический пример: валидация формы

```rust
import "lib/list" (contains)

type ValidationError = EmptyField(String)
                     | TooShort(String)
                     | InvalidEmail(String)

type User = { name: String, email: String }

fun validateName(name: String) -> Result<ValidationError, String> {
    if len(name) == 0 { Fail(EmptyField("name")) }
    else if len(name) < 2 { Fail(TooShort("name")) }
    else { Ok(name) }
}

fun validateEmail(email: String) -> Result<ValidationError, String> {
    if !contains(email, '@') { Fail(InvalidEmail(email)) }
    else { Ok(email) }
}

fun validateUser(name: String, email: String) -> Result<ValidationError, User> {
    match (validateName(name), validateEmail(email)) {
        (Ok(n), Ok(e)) -> Ok({ name: n, email: e })
        (Fail(err), _) -> Fail(err)
        (_, Fail(err)) -> Fail(err)
    }
}

// Использование
match validateUser("Alice", "alice@example.com") {
    Ok(user) -> print("Valid: " ++ user.name)
    Fail(EmptyField(f)) -> print("Field " ++ f ++ " is empty")
    Fail(TooShort(f)) -> print(f ++ " is too short")
    Fail(InvalidEmail(e)) -> print("Invalid email: " ++ e)
}
```

## Комбинирование нескольких Result

```rust
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

results = [Ok(1), Ok(2), Ok(3)]
print(sequence(results))  // Ok([1, 2, 3])

withError = [Ok(1), Fail("oops"), Ok(3)]
print(sequence(withError))  // Fail("oops")
```

## Когда использовать что?

| Ситуация | Тип |
|----------|-----|
| Значение может отсутствовать | `Option<T>` или `T?` |
| Нужна информация об ошибке | `Result<E, A>` |
| Быстрый null-check | `T?` (T \| Nil) |
| Сложная обработка ошибок | `Result<E, A>` |

## Преобразование между типами

```rust
fun optionToResult(opt, error) {
    match opt {
        Some(x) -> Ok(x)
        Zero -> Fail(error)
    }
}

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

