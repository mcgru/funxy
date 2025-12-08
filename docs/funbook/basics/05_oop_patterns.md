# 05. ООП паттерны в Funxy

## Для тех, кто привык к классам и объектам

Funxy не имеет классов в традиционном смысле, но все привычные паттерны легко реализуются — часто проще и безопаснее.

---

## Records = Data Classes

```rust
// Java/Kotlin: data class User(val name: String, val age: Int)
// Python: @dataclass class User

// Funxy:
type User = { name: String, age: Int, email: String }

// Создание "объекта"
user = { name: "Alice", age: 30, email: "alice@example.com" }

// Доступ к полям
print(user.name)  // Alice

// Иммутабельное обновление (как copy() в Kotlin)
older = { ...user, age: user.age + 1 }
print(older.age)  // 31
```

---

## Traits = Интерфейсы

```rust
// Java: interface Formatter { String format(); }
// Funxy:
trait Formatter<T> {
    fun format(self: T) -> String
}

// "Класс" User
type User = { name: String, age: Int }

// Реализация интерфейса
instance Formatter User {
    fun format(self: User) -> String {
        "User(" ++ self.name ++ ", " ++ show(self.age) ++ ")"
    }
}

// "Класс" Product
type Product = { name: String, price: Float }

instance Formatter Product {
    fun format(self: Product) -> String {
        self.name ++ ": $" ++ show(self.price)
    }
}

// Полиморфизм!
fun printItem<T: Formatter>(item: T) -> Nil { print(format(item)) }

u: User = { name: "Alice", age: 30 }
p: Product = { name: "Book", price: 19.99 }
printItem(u)  // User(Alice, 30)
printItem(p)  // Book: $19.99
```

---

## Методы на типах

```rust
// Java: class User { String greet() { return "Hi, " + name; } }

// Funxy — функции, принимающие тип первым аргументом
type User = { name: String, age: Int }

fun greet(user: User) -> String { "Hi, I'm " ++ user.name }

fun isAdult(user: User) -> Bool { user.age >= 18 }

fun haveBirthday(user: User) -> User { { ...user, age: user.age + 1 } }

// Использование (как методы через pipe)
alice = { name: "Alice", age: 30 }

print(alice |> greet)       // Hi, I'm Alice
print(alice |> isAdult)     // true

older = alice |> haveBirthday
print(older.age)            // 31
```

---

## Инкапсуляция через модули

Funxy использует модули для инкапсуляции. Экспортируются только нужные символы.

**counter.lang:**
```rust
// Модуль counter
module counter

// Приватный тип (не экспортируется)
type CounterState = { value: Int }

// Публичный "конструктор"
export fun newCounter(initial: Int) -> CounterState { { value: initial } }

// Публичные "методы"  
export fun increment(c: CounterState) -> CounterState { { value: c.value + 1 } }
export fun decrement(c: CounterState) -> CounterState { { value: c.value - 1 } }
export fun getValue(c: CounterState) -> Int { c.value }
// ...
```

**main.lang:**
```rust
import "counter" (newCounter, increment, getValue)

c = newCounter(0)
c2 = c |> increment |> increment |> increment
print(getValue(c2))  // 3
// ...
```

---

## ADT = Sealed Classes / Enums с данными

```rust
// Kotlin: sealed class Shape
// Java 17+: sealed interface Shape permits Circle, Rectangle

// Funxy:
type Shape = Circle(Float) | Rectangle((Float, Float))

// Паттерн "visitor" встроен в язык!
fun area(shape: Shape) -> Float {
    match shape {
        Circle(r) -> 3.14159 * r * r
        Rectangle((w, h)) -> w * h
    }
}

fun describe(shape: Shape) -> String {
    match shape {
        Circle(r) -> "Circle with radius " ++ show(r)
        Rectangle((w, h)) -> "Rectangle " ++ show(w) ++ "x" ++ show(h)
    }
}

// Использование
shapes = [Circle(5.0), Rectangle((4.0, 3.0))]

for s in shapes {
    print(describe(s) ++ " has area " ++ show(area(s)))
}
// Circle with radius 5 has area 78.53975
// Rectangle 4x3 has area 12
```

---

## Builder паттерн

```rust
// Java: new UserBuilder().name("Alice").age(30).build()

// Funxy — просто обновление записей
type User = { name: String, age: Int, email: String, role: String }

// "Builder" — просто дефолты + обновление
defaultUser :- { name: "", age: 0, email: "", role: "user" }

fun withName(u: User, name: String) -> User { { ...u, name: name } }
fun withAge(u: User, age: Int) -> User { { ...u, age: age } }
fun withEmail(u: User, email: String) -> User { { ...u, email: email } }
fun withRole(u: User, role: String) -> User { { ...u, role: role } }

// Fluent API через pipe
admin = defaultUser
    |> fun(u) -> withName(u, "Alice")
    |> fun(u) -> withAge(u, 30)
    |> fun(u) -> withEmail(u, "alice@example.com")
    |> fun(u) -> withRole(u, "admin")

print(admin.name)   // Alice
print(admin.role)   // admin
```

---

## Factory паттерн

```rust
type Shape = Circle(Float) | Rectangle((Float, Float))

// Factory функции
fun createCircle(radius: Float) -> Shape { Circle(radius) }
fun createSquare(side: Float) -> Shape { Rectangle((side, side)) }
fun createRectangle(w: Float, h: Float) -> Shape { Rectangle((w, h)) }

// Factory с валидацией
fun createValidCircle(radius: Float) -> Result<String, Shape> {
    if radius <= 0.0 { Fail("Radius must be positive") }
    else { Ok(Circle(radius)) }
}

print(createValidCircle(5.0))   // Ok(Circle(5))
print(createValidCircle(-1.0))  // Fail("Radius must be positive")
```

---

## Strategy паттерн

```rust
// Java: interface Strategy { int execute(int a, int b); }

// Funxy — просто функции!
type Strategy = (Int, Int) -> Int

fun add(a: Int, b: Int) -> Int { a + b }
fun multiply(a: Int, b: Int) -> Int { a * b }
fun power(a: Int, b: Int) -> Int { a ** b }

// Context
fun executeStrategy(strategy: Strategy, a: Int, b: Int) -> Int {
    strategy(a, b)
}

// Использование
print(executeStrategy(add, 5, 3))       // 8
print(executeStrategy(multiply, 5, 3))  // 15
print(executeStrategy(power, 5, 3))     // 125

// Или даже проще — передать лямбду
print(executeStrategy(fun(a, b) -> a - b, 5, 3))  // 2
```

---

## Observer паттерн

```rust
import "lib/list" (forEach)

// Type alias для функционального типа
type Observer = (Int) -> Nil
type Subject = { observers: List<Observer>, value: Int }

fun createSubject(initial: Int) -> Subject { 
    { observers: [], value: initial }
}

fun subscribe(subject: Subject, observer: Observer) -> Subject {
    { ...subject, observers: subject.observers ++ [observer] }
}

fun notify(subject: Subject) -> Nil {
    forEach(fun(obs) -> obs(subject.value), subject.observers)
}

fun setValue(subject: Subject, value: Int) -> Subject {
    updated = { ...subject, value: value }
    notify(updated)
    updated
}

// Использование
subject = createSubject(0)
s1 = subscribe(subject, fun(v) -> print("Observer 1: " ++ show(v)))
s2 = subscribe(s1, fun(v) -> print("Observer 2: " ++ show(v)))

_ = setValue(s2, 42)
// Observer 1: 42
// Observer 2: 42
```

---

## State через замыкания

```rust
// Инкапсулированное мутабельное состояние
fun createCounter() {
    count = 0
    {
        get: fun() -> count,
        inc: fun() { count += 1 },
        dec: fun() { count -= 1 }
    }
}

counter = createCounter()
counter.inc()
counter.inc()
counter.inc()
print(counter.get())  // 3
counter.dec()
print(counter.get())  // 2
```

---

## Наследование? Композиция!

```rust
import "lib/list" (contains)

// ООП: class Admin extends User
// Funxy: композиция вместо наследования

type User = { name: String, email: String }
type Admin = { user: User, permissions: List<String> }

// Функции для User
fun greetUser(u: User) -> String { "Hello, " ++ u.name }

// Admin может использовать функции User
fun greetAdmin(a: Admin) -> String { greetUser(a.user) ++ " (Admin)" }

fun hasPermission(a: Admin, perm: String) -> Bool {
    contains(a.permissions, perm)
}

// Создание
admin = {
    user: { name: "Alice", email: "alice@example.com" },
    permissions: ["read", "write", "delete"]
}

print(greetAdmin(admin))                    // Hello, Alice (Admin)
print(hasPermission(admin, "delete"))       // true
```

---

## Сравнение парадигм

| ООП концепт | Funxy эквивалент |
|-------------|-----------------|
| Class | `type` (record) |
| Interface | `trait` |
| Method | Функция с типом как первый аргумент |
| Inheritance | Композиция записей |
| Private fields | Модули (export только нужное) |
| Sealed class | ADT (sum types) |
| Factory | Функция-конструктор |
| Builder | Pipe + функции обновления |
| Strategy | Функции первого класса |
| Singleton | Константа в модуле |

---

## Преимущества подхода Funxy

1. Иммутабельность по умолчанию — нет "случайных" мутаций
2. Exhaustive matching — компилятор проверит все случаи ADT
3. Нет null pointer exceptions — Option/Result типы
4. Простота — меньше boilerplate, больше сути
5. Композиция > наследование — гибче и понятнее
