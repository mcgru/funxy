# Большие числа (lib/bignum)

Модуль `lib/bignum` предоставляет функции для работы с числами произвольной точности.

```rust
import "lib/bignum" (*)
```

## Типы

- **BigInt** — целые числа произвольной точности (литерал с суффиксом `n`)
- **Rational** — рациональные числа (дроби) с произвольной точностью (литерал с суффиксом `r`)

## BigInt

### Литералы

```rust
a = 123456789012345678901234567890n
b = -42n
c = 0n
```

### Создание

```rust
import "lib/bignum" (*)

// bigIntNew: (String) -> BigInt
x = bigIntNew("999999999999999999999999")

// bigIntFromInt: (Int) -> BigInt
y = bigIntFromInt(42)
```

### Конвертация

```rust
import "lib/bignum" (*)

// bigIntToString: (BigInt) -> String
x = bigIntNew("999999999999999999999999")
s = bigIntToString(x)  // "999999999999999999999999"

// bigIntToInt: (BigInt) -> Option<Int>
// Возвращает Some если значение влезает в Int, иначе Zero
bigIntToInt(bigIntFromInt(100))  // Some(100)
bigIntToInt(bigIntNew("99999999999999999999"))  // Zero (слишком большое)
```

### Операторы

BigInt поддерживает все арифметические операторы через трейт `Numeric`:

```rust
a = 1000000000000000000n
b = 2000000000000000000n

a + b   // 3000000000000000000
b - a   // 1000000000000000000
a * b   // 2000000000000000000000000000000000000
b / a   // 2
b % a   // 0
2n ** 100n  // 1267650600228229401496703205376
```

Сравнение через трейт `Order`:

```rust
a = 1000000000000000000n
b = 2000000000000000000n

a < b   // true
a == a  // true
a > b   // false
```

## Rational

Рациональные числа хранят точные дроби без потери точности.

### Литералы

```rust
half = 0.5r
third = 1.0r / 3.0r  // точно 1/3, не приближение
```

### Создание

```rust
import "lib/bignum" (*)

// ratFromInt: (Int, Int) -> Rational
// Автоматически упрощает дробь
r1 = ratFromInt(1, 3)   // 1/3
r2 = ratFromInt(2, 4)   // 1/2 (упрощено)
r3 = ratFromInt(-6, 8)  // -3/4 (упрощено)

// ratNew: (BigInt, BigInt) -> Rational
num = bigIntNew("22")
denom = bigIntNew("7")
pi_approx = ratNew(num, denom)  // 22/7
```

### Доступ к числителю и знаменателю

```rust
import "lib/bignum" (*)

// ratNumer: (Rational) -> BigInt
// ratDenom: (Rational) -> BigInt
frac = ratFromInt(3, 4)
bigIntToString(ratNumer(frac))    // "3"
bigIntToString(ratDenom(frac))  // "4"
```

### Конвертация

```rust
// ratToFloat: (Rational) -> Option<Float>
// Возвращает Zero если результат не влезает в Float (infinity/NaN)
ratToFloat(ratFromInt(1, 4))  // Some(0.25)
ratToFloat(ratFromInt(1, 3))  // Some(0.333...)

// ratToString: (Rational) -> String
// Формат "числитель/знаменатель"
ratToString(ratFromInt(3, 4))  // "3/4"
ratToString(ratFromInt(7, 1))  // "7/1" или "7"
```

### Операторы

Rational поддерживает арифметику через трейт `Numeric`:

```rust
import "lib/bignum" (*)

half = ratFromInt(1, 2)
third = ratFromInt(1, 3)

ratToString(half + third)  // "5/6"
ratToString(half - third)  // "1/6"
ratToString(half * third)  // "1/6"
ratToString(half / third)  // "3/2"
```

Сравнение:

```rust
import "lib/bignum" (*)

half = ratFromInt(1, 2)
third = ratFromInt(1, 3)

half > third   // true
half == half   // true
```

## Практические примеры

### Точная арифметика

Floating point имеет ошибки округления:

```rust
import "lib/bignum" (*)

// В обычных float: 0.1 + 0.2 = 0.30000000000000004
// В Rational: точно 3/10
sum = ratFromInt(1, 10) + ratFromInt(2, 10)
ratToString(sum)   // "3/10"
ratToFloat(sum)    // Some(0.3)
```

### Факториал больших чисел

```rust
import "lib/bignum" (*)

fun factorial(n: BigInt) -> BigInt {
    if bigIntToInt(n) == Some(0) { 1n }
    else { n * factorial(n - 1n) }
}

factorial(20n)   // 2432902008176640000
factorial(100n)  // огромное число (~158 цифр)
```

### Числа Фибоначчи

```rust
import "lib/bignum" (*)
import "lib/list" (range)

fun fib(n: BigInt) -> BigInt {
    if bigIntToInt(n) == Some(0) { 0n }
    else if bigIntToInt(n) == Some(1) { 1n }
    else { fib(n - 1n) + fib(n - 2n) }
}

// Или итеративно для больших n
fun fibIter(n: Int) -> BigInt {
    a = 0n
    b = 1n
    for i in range(0, n) {
        temp = a + b
        a = b
        b = temp
    }
    a
}

fibIter(100)  // 354224848179261915075
```

### Точные финансовые вычисления

```rust
import "lib/bignum" (*)

// Деньги как центы в Rational
dollars = fun(d: Int, c: Int) -> Rational {
    ratFromInt(d * 100 + c, 100)
}

price = dollars(19, 99)    // $19.99
tax = ratFromInt(8, 100)  // 8%
total = price + price * tax

ratToString(total)  // точная сумма
```

## Сводка

### BigInt

| Функция | Тип | Описание |
|---------|-----|----------|
| `bigIntNew` | `(String) -> BigInt` | Парсинг из строки |
| `bigIntFromInt` | `(Int) -> BigInt` | Int → BigInt |
| `bigIntToString` | `(BigInt) -> String` | BigInt → String |
| `bigIntToInt` | `(BigInt) -> Option<Int>` | BigInt → Int (если влезает) |

### Rational

| Функция | Тип | Описание |
|---------|-----|----------|
| `ratFromInt` | `(Int, Int) -> Rational` | Создать из числителя/знаменателя |
| `ratNew` | `(BigInt, BigInt) -> Rational` | Создать из BigInt |
| `ratNumer` | `(Rational) -> BigInt` | Получить числитель |
| `ratDenom` | `(Rational) -> BigInt` | Получить знаменатель |
| `ratToFloat` | `(Rational) -> Option<Float>` | Rational → Float |
| `ratToString` | `(Rational) -> String` | Rational → "a/b" |

### Операторы (работают автоматически)

- Арифметика: `+`, `-`, `*`, `/`, `%`, `**`
- Сравнение: `<`, `<=`, `>`, `>=`, `==`, `!=`

