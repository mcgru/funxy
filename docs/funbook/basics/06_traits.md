# Трейты и операторы

## Базовый трейт

```rust
// Объявление трейта
trait Formatter<T> {
    fun format(val: T) -> String
}

// ADT тип
type UserId = MkUserId(Int)

instance Formatter UserId {
    fun format(val: UserId) -> String {
        match val { MkUserId(n) -> "User#" ++ show(n) }
    }
}

uid = MkUserId(42)
print(format(uid))  // User#42
```

---

## Супер трейты (наследование)

Трейт может наследоваться от других трейтов:

```rust
type Ordering = Lt | Eq | Gt

// Базовый трейт
trait MyEq<T> {
    fun myEq(a: T, b: T) -> Bool
}

// MyOrd наследует MyEq
trait MyOrd<T> : MyEq<T> {
    fun myCmp(a: T, b: T) -> Ordering
}

// ADT для денег
type Money = MkMoney(Int)

fun getAmount(m: Money) -> Int {
    match m { MkMoney(n) -> n }
}

// Сначала MyEq
instance MyEq Money {
    fun myEq(a: Money, b: Money) -> Bool {
        getAmount(a) == getAmount(b)
    }
}

// Потом MyOrd (требует MyEq)
instance MyOrd Money {
    fun myCmp(a: Money, b: Money) -> Ordering {
        if getAmount(a) < getAmount(b) { Lt }
        else if getAmount(a) > getAmount(b) { Gt }
        else { Eq }
    }
}

m1 = MkMoney(100)
m2 = MkMoney(200)
print(myEq(m1, m1))   // true
print(myCmp(m1, m2))  // Lt
```

---

## Перегрузка операторов

```rust
// ADT для 2D вектора
type Vec2 = MkVec2((Float, Float))

fun getX(v: Vec2) -> Float { match v { MkVec2((x, _)) -> x } }
fun getY(v: Vec2) -> Float { match v { MkVec2((_, y)) -> y } }

// Оператор сложения
trait Addable<T> {
    operator (+)(a: T, b: T) -> T
}

instance Addable Vec2 {
    operator (+)(a: Vec2, b: Vec2) -> Vec2 {
        MkVec2((getX(a) + getX(b), getY(a) + getY(b)))
    }
}

// Оператор равенства  
trait Equalable<T> {
    operator (==)(a: T, b: T) -> Bool
    operator (!=)(a: T, b: T) -> Bool {
        !(a == b)
    }
}

instance Equalable Vec2 {
    operator (==)(a: Vec2, b: Vec2) -> Bool {
        getX(a) == getX(b) && getY(a) == getY(b)
    }
}

v1 = MkVec2((1.0, 2.0))
v2 = MkVec2((3.0, 4.0))
v3 = v1 + v2
print(v3)           // MkVec2((4.0, 6.0))
print(v1 == v1)     // true
print(v1 != v2)     // true
```

---

## Оператор $ (применение)

Низкоприоритетное применение функции:

```rust
fun double(x: Int) -> Int { x * 2 }
fun inc(x: Int) -> Int { x + 1 }

// f $ x = f(x)
print(double $ 21)  // 42

// Правоассоциативный: f $ g $ x = f(g(x))
print(inc $ double $ 5)  // 11 = inc(double(5))

// Удобно для избежания скобок
print $ double $ inc $ 10  // 22
```

---

## Операторы как функции

Любой оператор можно использовать как функцию:

```rust
import "lib/list" (foldl)

// Оператор в переменной
add = (+)
print(add(1, 2))  // 3

// В higher-order функциях
sum = foldl((+), 0, [1, 2, 3, 4, 5])
print(sum)  // 15

product = foldl((*), 1, [1, 2, 3, 4])
print(product)  // 24
```

---

## Ограничения (constraints)

```rust
trait Showable<T> {
    fun render(val: T) -> String
}

instance Showable Int {
    fun render(val: Int) -> String { "Int:" ++ show(val) }
}

// Функция требует Showable
fun displayValue<T: Showable>(x: T) -> String {
    render(x)
}

print(displayValue(42))  // Int:42
```

---

## Реализации по умолчанию

```rust
trait MathOps<T> {
    fun mathAdd(a: T, b: T) -> T
    fun mathMul(a: T, b: T) -> T
    
    // Реализация по умолчанию
    fun mathSquare(x: T) -> T {
        mathMul(x, x)
    }
    
    fun mathDouble(x: T) -> T {
        mathAdd(x, x)
    }
}

type MyNum = MkMyNum(Int)

fun getVal(n: MyNum) -> Int { match n { MkMyNum(v) -> v } }

instance MathOps MyNum {
    fun mathAdd(a: MyNum, b: MyNum) -> MyNum {
        MkMyNum(getVal(a) + getVal(b))
    }
    fun mathMul(a: MyNum, b: MyNum) -> MyNum {
        MkMyNum(getVal(a) * getVal(b))
    }
    // mathSquare и mathDouble работают автоматически!
}

n = MkMyNum(5)
print(mathSquare(n))  // MkMyNum(25)
print(mathDouble(n))  // MkMyNum(10)
```

---

## Пользовательские операторы (UserOp)

Доступные слоты для пользовательских операторов:

| Оператор | Трейт | Ассоциативность | Типичное использование |
|----------|-------|-----------------|------------------------|
| `<>` | `Semigroup` | Правая | Объединение (Semigroup) |
| `<\|>` | `UserOpChoose` | Левая | Альтернатива |
| `<:>` | `UserOpCons` | Правая | Cons-подобный prepend |
| `<~>` | `UserOpSwap` | Левая | Обмен/Swap |
| `<$>` | `UserOpMap` | Левая | Functor map |
| `=>` | `UserOpImply` | Правая | Импликация |
| `<\|` | `UserOpPipeLeft` | Правая | Правый pipe |
| `$` | (встроенный) | Правая | Применение функции |

### Semigroup (<>)

```rust
type Text = MkText(String)

fun getText(t: Text) -> String { match t { MkText(s) -> s } }

instance Semigroup Text {
    operator (<>)(a: Text, b: Text) -> Text {
        match (a, b) { (MkText(x), MkText(y)) -> MkText(x ++ y) }
    }
}

t1 = MkText("Hello")
t2 = MkText(" ")
t3 = MkText("World")
result = t1 <> t2 <> t3  // правоассоциативный
print(getText(result))   // Hello World
```

### UserOpChoose (<|>) — Альтернатива

```rust
type Maybe = MkJust(Int) | MkNothing

fun getMaybe(m: Maybe) -> Int { 
    match m { 
        MkJust(x) -> x 
        MkNothing -> -1 
    } 
}

instance UserOpChoose Maybe {
    operator (<|>)(a: Maybe, b: Maybe) -> Maybe {
        match a {
            MkJust(_) -> a
            MkNothing -> b
        }
    }
}

m1 = MkNothing
m2 = MkJust(42)
print(getMaybe(m1 <|> m2))  // 42 (первый non-Nothing)
```

### UserOpImply (=>) — Импликация

```rust
type Logic = MkLogic(Bool)

fun getLogic(l: Logic) -> Bool { match l { MkLogic(b) -> b } }

instance UserOpImply Logic {
    operator (=>)(a: Logic, b: Logic) -> Logic {
        // a => b эквивалентно !a || b
        match (a, b) {
            (MkLogic(x), MkLogic(y)) -> MkLogic(!x || y)
        }
    }
}

lt = MkLogic(true)
lf = MkLogic(false)
print(getLogic(lt => lt))  // true  (true => true)
print(getLogic(lt => lf))  // false (true => false)
print(getLogic(lf => lt))  // true  (false => anything)
```

### Операторы как функции

```rust
type Text = MkText(String)
fun getText(t: Text) -> String { match t { MkText(s) -> s } }

instance Semigroup Text {
    operator (<>)(a: Text, b: Text) -> Text {
        match (a, b) { (MkText(x), MkText(y)) -> MkText(x ++ y) }
    }
}

// Оператор в переменной
combine = (<>)
t4 = combine(MkText("A"), MkText("B"))
print(getText(t4))  // AB
```

---

## Поддерживаемые операторы

| Оператор | Описание |
|----------|----------|
| `+`, `-`, `*`, `/`, `%`, `**` | Арифметика |
| `==`, `!=` | Равенство |
| `<`, `>`, `<=`, `>=` | Сравнение |
| `&`, `\|`, `^`, `<<`, `>>` | Битовые |
| `++` | Конкатенация |
| `::` | Cons (prepend) |
| `\|>` | Pipe |
| `$` | Применение функции |
| `,,` | Композиция |
| `<>` | Semigroup combine |
| `<\|>` | Alternative choice |
| `<:>` | Custom cons |
| `<~>` | Swap |
| `<$>` | Functor map |
| `=>` | Implication |
