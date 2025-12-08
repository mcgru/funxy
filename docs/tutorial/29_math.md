# Математические функции (lib/math)

Модуль `lib/math` предоставляет математические функции и константы.

```rust
import "lib/math" (*)
```

## Базовые операции

```rust
// abs: (Float) -> Float — модуль числа
abs(-5.5)         // 5.5
abs(3.14)         // 3.14

// absInt: (Int) -> Int — модуль целого
absInt(-10)       // 10

// sign: (Float) -> Int — знак числа (-1, 0, 1)
sign(-3.14)       // -1
sign(0.0)         // 0
sign(42.0)        // 1

// min, max: (Float, Float) -> Float
min(3.0, 7.0)     // 3
max(3.0, 7.0)     // 7

// minInt, maxInt: (Int, Int) -> Int
minInt(3, 7)      // 3
maxInt(3, 7)      // 7

// clamp: (Float, Float, Float) -> Float — ограничить диапазоном
clamp(5.0, 0.0, 10.0)   // 5 (в диапазоне)
clamp(-5.0, 0.0, 10.0)  // 0 (меньше min)
clamp(15.0, 0.0, 10.0)  // 10 (больше max)
```

## Округление

```rust
// floor: (Float) -> Int — округление вниз
floor(3.7)        // 3
floor(-3.7)       // -4

// ceil: (Float) -> Int — округление вверх
ceil(3.2)         // 4
ceil(-3.2)        // -3

// round: (Float) -> Int — округление к ближайшему
round(3.4)        // 3
round(3.5)        // 4

// trunc: (Float) -> Int — отбрасывание дробной части
trunc(3.7)        // 3
trunc(-3.7)       // -3
```

## Степени и корни

```rust
// sqrt: (Float) -> Float — квадратный корень
sqrt(16.0)        // 4
sqrt(2.0)         // 1.414...

// cbrt: (Float) -> Float — кубический корень
cbrt(27.0)        // 3
cbrt(-8.0)        // -2

// pow: (Float, Float) -> Float — возведение в степень
pow(2.0, 10.0)    // 1024
pow(4.0, 0.5)     // 2 (квадратный корень)

// exp: (Float) -> Float — e^x
exp(0.0)          // 1
exp(1.0)          // e (~2.718)
```

## Логарифмы

```rust
// log: (Float) -> Float — натуральный логарифм (ln)
log(1.0)          // 0
log(e())          // 1

// log10: (Float) -> Float — десятичный логарифм
log10(100.0)      // 2
log10(1000.0)     // 3

// log2: (Float) -> Float — двоичный логарифм
log2(8.0)         // 3
log2(1024.0)      // 10
```

## Тригонометрия

Все функции работают с радианами.

```rust
// sin, cos, tan: (Float) -> Float
sin(0.0)          // 0
cos(0.0)          // 1
tan(0.0)          // 0

sin(pi() / 2.0)   // 1
cos(pi())         // -1

// asin, acos, atan: (Float) -> Float — обратные функции
asin(0.0)         // 0
acos(1.0)         // 0
atan(1.0)         // pi/4 (~0.785)

// atan2: (Float, Float) -> Float — угол точки (y, x)
atan2(1.0, 1.0)   // pi/4
```

## Гиперболические функции

```rust
sinh(0.0)         // 0
cosh(0.0)         // 1
tanh(0.0)         // 0
```

## Константы

```rust
// pi: () -> Float
pi()              // 3.141592653589793

// e: () -> Float — число Эйлера
e()               // 2.718281828459045
```

## Конвертация

```rust
// Функции intToFloat и floatToInt теперь в prelude (не требуют импорта lib/math)

// intToFloat: (Int) -> Float
intToFloat(42)       // 42.0

// floatToInt: (Float) -> Int — отбрасывает дробную часть
floatToInt(3.99)       // 3
floatToInt(-3.99)      // -3
```

## Практические примеры

### Расстояние между точками

```rust
import "lib/math" (sqrt)

fun distance(x1: Float, y1: Float, x2: Float, y2: Float) -> Float {
    dx = x2 - x1
    dy = y2 - y1
    sqrt(dx * dx + dy * dy)
}

distance(0.0, 0.0, 3.0, 4.0)  // 5
```

### Градусы в радианы

```rust
import "lib/math" (pi, sin, cos)

fun degToRad(deg: Float) -> Float {
    deg * pi() / 180.0
}

fun radToDeg(rad: Float) -> Float {
    rad * 180.0 / pi()
}

sin(degToRad(90.0))   // 1
cos(degToRad(180.0))  // -1
```

### Площадь круга

```rust
import "lib/math" (pi)

fun circleArea(r: Float) -> Float {
    pi() * r * r
}

circleArea(1.0)   // pi
circleArea(2.0)   // 4*pi
```

### Нормализация значения

```rust
import "lib/math" (clamp)

fun normalize(value: Float, min: Float, max: Float) -> Float {
    (clamp(value, min, max) - min) / (max - min)
}

normalize(5.0, 0.0, 10.0)   // 0.5
normalize(-5.0, 0.0, 10.0)  // 0
normalize(15.0, 0.0, 10.0)  // 1
```

### Квадратное уравнение

```rust
import "lib/math" (sqrt, abs)

fun solveQuadratic(a: Float, b: Float, c: Float) -> Option<(Float, Float)> {
    discriminant = b * b - 4.0 * a * c
    if discriminant < 0.0 {
        Zero
    } else {
        sqrtD = sqrt(discriminant)
        x1 = (-b + sqrtD) / (2.0 * a)
        x2 = (-b - sqrtD) / (2.0 * a)
        Some((x1, x2))
    }
}

// x² - 5x + 6 = 0 → x = 2, x = 3
solveQuadratic(1.0, -5.0, 6.0)  // Some((3, 2))
```

## Сводка

| Функция | Тип | Описание |
|---------|-----|----------|
| `abs` | `(Float) -> Float` | Модуль |
| `absInt` | `(Int) -> Int` | Модуль целого |
| `sign` | `(Float) -> Int` | Знак (-1, 0, 1) |
| `min` | `(Float, Float) -> Float` | Минимум |
| `max` | `(Float, Float) -> Float` | Максимум |
| `minInt` | `(Int, Int) -> Int` | Минимум целых |
| `maxInt` | `(Int, Int) -> Int` | Максимум целых |
| `clamp` | `(Float, Float, Float) -> Float` | Ограничить диапазоном |
| `floor` | `(Float) -> Int` | Округление вниз |
| `ceil` | `(Float) -> Int` | Округление вверх |
| `round` | `(Float) -> Int` | К ближайшему |
| `trunc` | `(Float) -> Int` | Отбросить дробь |
| `sqrt` | `(Float) -> Float` | Квадратный корень |
| `cbrt` | `(Float) -> Float` | Кубический корень |
| `pow` | `(Float, Float) -> Float` | Степень |
| `exp` | `(Float) -> Float` | e^x |
| `log` | `(Float) -> Float` | Натуральный логарифм |
| `log10` | `(Float) -> Float` | Десятичный логарифм |
| `log2` | `(Float) -> Float` | Двоичный логарифм |
| `sin` | `(Float) -> Float` | Синус |
| `cos` | `(Float) -> Float` | Косинус |
| `tan` | `(Float) -> Float` | Тангенс |
| `asin` | `(Float) -> Float` | Арксинус |
| `acos` | `(Float) -> Float` | Арккосинус |
| `atan` | `(Float) -> Float` | Арктангенс |
| `atan2` | `(Float, Float) -> Float` | Угол точки (y, x) |
| `sinh` | `(Float) -> Float` | Гиперболический синус |
| `cosh` | `(Float) -> Float` | Гиперболический косинус |
| `tanh` | `(Float) -> Float` | Гиперболический тангенс |
| `pi` | `() -> Float` | Число π |
| `e` | `() -> Float` | Число e |

**Примечание:** Функции `intToFloat` и `floatToInt` перемещены в prelude (доступны без импорта).

