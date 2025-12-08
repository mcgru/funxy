# 03. BigInt и Rational

## Задача
Работать с числами произвольной точности: большие целые и точные дроби.

---

## BigInt — произвольно большие целые

Обычный `Int` ограничен размером машинного слова. `BigInt` не имеет лимита.

### Создание

```rust
import "lib/bignum" (bigIntNew, bigIntFromInt, bigIntToString, bigIntToInt)

// Из строки (для очень больших чисел)
huge = bigIntNew("123456789012345678901234567890")

// Из обычного Int
small = bigIntFromInt(42)

// Конвертация обратно
print(bigIntToString(huge))  // "123456789012345678901234567890"

// В Int (может не поместиться!)
match bigIntToInt(small) {
    Some(n) -> print(n)       // 42
    Zero -> print("Too big!")
}
```

### Арифметика

```rust
import "lib/bignum" (bigIntNew, bigIntToString)

a = bigIntNew("1000000000000000000")
b = bigIntNew("2000000000000000000")

// BigInt поддерживает стандартные операторы
sum = a + b   // 3000000000000000000
diff = b - a  // 1000000000000000000
prod = a * b  // 2000000000000000000000000000000000000
quot = b / a  // 2

print(bigIntToString(prod))
```

### Практический пример: факториал

```rust
import "lib/bignum" (bigIntFromInt, bigIntToString)

fun factorial(n: Int) -> BigInt {
    fun go(i: Int, acc: BigInt) -> BigInt {
        if i <= 1 { acc } else { go(i - 1, acc * bigIntFromInt(i)) }
    }
    go(n, bigIntFromInt(1))
}

print(bigIntToString(factorial(50)))
// 30414093201713378043612608166064768844377641568960512000000000000
```

### Практический пример: числа Фибоначчи

```rust
import "lib/bignum" (bigIntFromInt, bigIntToString)

fun fib(n: Int) -> BigInt {
    fun go(i: Int, a: BigInt, b: BigInt) -> BigInt {
        if i == 0 { a }
        else { go(i - 1, b, a + b) }
    }
    go(n, bigIntFromInt(0), bigIntFromInt(1))
}

print(bigIntToString(fib(100)))
// 354224848179261915075
```

---

## Rational — точные дроби

`Rational` хранит числитель и знаменатель как `BigInt`. Никакой потери точности!

### Создание

```rust
import "lib/bignum" (ratFromInt, ratNew, bigIntNew, ratNumer, ratDenom, bigIntToString)

// Из двух Int
half = ratFromInt(1, 2)       // 1/2
third = ratFromInt(1, 3)      // 1/3
twoThirds = ratFromInt(2, 3)  // 2/3

// Из BigInt
huge = ratNew(bigIntNew("1"), bigIntNew("3"))

// Доступ к частям
print(bigIntToString(ratNumer(half)))  // "1"
print(bigIntToString(ratDenom(half)))  // "2"
```

### Арифметика

```rust
import "lib/bignum" (ratFromInt, ratToString)

a = ratFromInt(1, 2)   // 1/2
b = ratFromInt(1, 3)   // 1/3

// Стандартные операторы
sum = a + b    // 5/6
diff = a - b   // 1/6
prod = a * b   // 1/6
quot = a / b   // 3/2

print(ratToString(sum))   // "5/6"
print(ratToString(prod))  // "1/6"
print(ratToString(quot))  // "3/2"
```

### Compound операторы

```rust
import "lib/bignum" (ratFromInt, ratToString)

r = ratFromInt(1, 2)

r += ratFromInt(1, 4)   // 3/4
r *= ratFromInt(2, 1)   // 3/2 (= 6/4 упрощается)
r -= ratFromInt(1, 2)   // 1/1 = 1

print(ratToString(r))  // "1/1"
```

### Автоматическое упрощение

```rust
import "lib/bignum" (ratFromInt, ratToString)

// Дроби автоматически приводятся к несократимому виду
frac = ratFromInt(4, 8)  
print(ratToString(frac))  // "1/2" (не "4/8")

frac2 = ratFromInt(15, 25)
print(ratToString(frac2))  // "3/5"
```

### Конвертация в Float

```rust
import "lib/bignum" (ratFromInt, ratToFloat)

r = ratFromInt(1, 3)

match ratToFloat(r) {
    Some(f) -> print(f)  // 0.3333333333333333
    Zero -> print("Cannot convert")
}
```

### Практический пример: точные финансовые расчёты

```rust
import "lib/bignum" (ratFromInt, ratToString)

// Проблема Float: 0.1 + 0.2 != 0.3
print(0.1 + 0.2)  // 0.30000000000000004

// Решение: Rational
a = ratFromInt(1, 10)   // 0.1
b = ratFromInt(2, 10)   // 0.2
c = ratFromInt(3, 10)   // 0.3

print(a + b == c)  // true

// Расчёт скидки
fun applyDiscount(price: Rational, discountPercent: Int) -> Rational {
    discount = ratFromInt(discountPercent, 100)
    price * (ratFromInt(1, 1) - discount)
}

price = ratFromInt(9999, 100)  // $99.99
discounted = applyDiscount(price, 15)  // 15% скидка
print(ratToString(discounted))  // точный результат!
```

### Практический пример: сложные проценты

```rust
import "lib/bignum" (ratFromInt, ratToString)

// Сложные проценты: A = P * (1 + r/n)^(n*t)
// Точный расчёт без потери точности

fun compoundInterest(
    principal: Rational,
    rate: Rational,
    times: Int,
    years: Int
) -> Rational {
    base = ratFromInt(1, 1) + rate / ratFromInt(times, 1)
    periods = times * years
    
    // Возведение в степень
    fun pow(r: Rational, n: Int) -> Rational {
        if n == 0 { ratFromInt(1, 1) }
        else { r * pow(r, n - 1) }
    }
    
    principal * pow(base, periods)
}

// $1000 под 5% годовых, ежемесячно, 10 лет
result = compoundInterest(
    ratFromInt(1000, 1),
    ratFromInt(5, 100),
    12,
    10
)
print(ratToString(result))
```

---

## Сравнение типов

| Тип | Диапазон | Точность | Скорость |
|-----|----------|----------|----------|
| `Int` | ±2^63 | Точный | Быстро |
| `Float` | ±10^308 | Приблизительный | Быстро |
| `BigInt` | Бесконечный | Точный | Медленнее |
| `Rational` | Бесконечный | Точный | Медленнее |

---

## Когда использовать

- Int — счётчики, индексы, обычная арифметика
- Float — научные расчёты, графика, где небольшая погрешность допустима
- BigInt — криптография, факториалы, числа Фибоначчи
- Rational — финансы, точные пропорции, где 0.1 + 0.2 = 0.3 критично
