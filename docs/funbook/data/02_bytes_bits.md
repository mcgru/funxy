# 02. Bytes и Bits

## Задача
Работать с бинарными данными: протоколы, файлы, сетевые пакеты.

---

## Bits — произвольные последовательности битов

Уникальная фича: длина не обязана быть кратна 8! Идеально для протоколов с bit fields.

### Литералы

```rust
// Бинарные литералы (1 бит на символ)
b1 = #b"10101010"    // 8 бит
b2 = #b"101"         // 3 бита (не byte-aligned!)
b3 = #b""            // пусто

// Hex литералы (4 бита на символ)
b4 = #x"FF"          // 8 бит: 11111111
b5 = #x"A5"          // 8 бит: 10100101

// Octal литералы (3 бита на символ)
b6 = #o"7"           // 3 бита: 111
b7 = #o"377"         // 9 бит: 011111111

print(b1)
print(b4)
```

### Создание и конвертация

```rust
import "lib/bits" (bitsFromBytes, bitsToBinary, bitsFromBinary, bitsToHex, bitsFromHex, bitsToBytes)

// Из строк
match bitsFromBinary("10101010") {
    Ok(b) -> print(bitsToBinary(b))  // "10101010"
    Fail(e) -> print(e)
}

match bitsFromHex("DEADBEEF") {
    Ok(b) -> print(bitsToHex(b))  // "deadbeef"
    Fail(e) -> print(e)
}
```

### Операции

```rust
import "lib/bits" (bitsGet, bitsSlice, bitsSet, bitsPadLeft, bitsPadRight)

b = #b"10101010"

// Длина
print(len(b))  // 8

// Доступ по индексу
match bitsGet(b, 0) {
    Some(bit) -> print(bit)  // 1
    Zero -> print("out of bounds")
}

// Slice [start, end)
part = bitsSlice(b, 0, 4)
print(part)  // #b"1010"

// Конкатенация
joined = #b"1111" ++ #b"0000"
print(joined)  // #b"11110000"

// Паддинг
padL = bitsPadLeft(#b"101", 8)
print(padL)  // #b"00000101"

padR = bitsPadRight(#b"101", 8)
print(padR)  // #b"10100000"
```

### Добавление числовых значений

```rust
import "lib/bits" (bitsNew, bitsAddInt, bitsAddFloat)

b = bitsNew()

// Целые числа (value, size in bits, endianness)
b = bitsAddInt(b, 255, 8)              // big endian (default)
b = bitsAddInt(b, 255, 8, "big")       // явно big endian
b = bitsAddInt(b, 1, 16, "little")     // little endian

// Float (IEEE 754)
bf = bitsAddFloat(bitsNew(), 3.14, 32)  // 32 bits
print(bf)
```

### Pattern Matching — парсинг бинарных протоколов

```rust
import "lib/bits" (bitsExtract, bitsInt)
import "lib/map" (mapGet)

// Пакет с bit fields:
// - version: 8 бит
// - flags: 4 бита
// - reserved: 4 бита
// - length: 16 бит
packet = #b"00000001010100000000000100000000"

// Спецификация полей
specs = [
    bitsInt("version", 8, "big"),
    bitsInt("flags", 4, "big"),
    bitsInt("reserved", 4, "big"),
    bitsInt("length", 16, "big")
]

// Извлечение!
match bitsExtract(packet, specs) {
    Ok(fields) -> {
        version = mapGet(fields, "version")
        flags = mapGet(fields, "flags")
        length = mapGet(fields, "length")
        print("Version: " ++ show(version))
        print("Flags: " ++ show(flags))
        print("Length: " ++ show(length))
    }
    Fail(err) -> print("Parse error: " ++ err)
}
```

### Spec-функции для парсинга

```rust
import "lib/bits" (bitsInt, bitsBytes, bitsRest)

// Целое число: bitsInt(name, size_in_bits, endianness)
spec1 = bitsInt("count", 16, "big")
spec2 = bitsInt("offset", 32, "little")

// Байты: bitsBytes(name, size_in_bytes)
spec3 = bitsBytes("payload", 16)

// Фиксированный размер
spec4 = bitsBytes("data", 64)

// Остаток битов
spec5 = bitsRest("tail")

```

### Практический пример: проверка PNG

```rust
import "lib/bits" (bitsExtract, bitsInt, bitsFromBytes)
import "lib/map" (mapGet)

fun isPNG(data) -> Bool {
    specs = [
        bitsInt("magic1", 32, "big"),
        bitsInt("magic2", 32, "big")
    ]
    match bitsExtract(data, specs) {
        Ok(fields) -> {
            m1 = mapGet(fields, "magic1")
            m2 = mapGet(fields, "magic2")
            match (m1, m2) {
                (Some(v1), Some(v2)) -> v1 == 0x89504E47 && v2 == 0x0D0A1A0A
                _ -> false
            }
        }
        Fail(_) -> false
    }
}

// PNG magic bytes: 89 50 4E 47 0D 0A 1A 0A
print(isPNG(#x"89504E470D0A1A0A"))  // true
print(isPNG(#x"FFD8FFE0"))          // false (это JPEG)
```

---

## Bytes — последовательности байтов

### Создание

```rust
import "lib/bytes" (bytesFromString, bytesFromList, bytesFromHex, bytesToString)

// Из строки
b = bytesFromString("Hello")

// Из списка байтов
b2 = bytesFromList([0x48, 0x65, 0x6C, 0x6C, 0x6F])

// Из hex строки
match bytesFromHex("48656C6C6F") {
    Ok(b) -> {
        match bytesToString(b) {
            Ok(s) -> print(s)  // "Hello"
            Fail(e) -> print(e)
        }
    }
    Fail(e) -> print(e)
}
```

### Конвертация

```rust
import "lib/bytes" (bytesFromString, bytesToString, bytesToList, bytesToHex)

b = bytesFromString("Hello")

// В строку
match bytesToString(b) {
    Ok(s) -> print(s)  // "Hello"
    Fail(e) -> print("Not valid UTF-8")
}

// В список байтов
list = bytesToList(b)
print(list)  // [72, 101, 108, 108, 111]

// В hex
hex = bytesToHex(b)
print(hex)  // "48656c6c6f"
```

### Операции

```rust
import "lib/bytes" (bytesFromString, bytesConcat, bytesSlice, bytesContains, bytesStartsWith, bytesEndsWith, bytesIndexOf, bytesSplit, bytesJoin)

b1 = bytesFromString("Hello")
b2 = bytesFromString(" World")

// Конкатенация
joined = bytesConcat(b1, b2)

// Slice
part = bytesSlice(b1, 0, 3)  // "Hel"

// Поиск
print(bytesContains(b1, bytesFromString("ell")))     // true
print(bytesStartsWith(b1, bytesFromString("He")))    // true
print(bytesEndsWith(b1, bytesFromString("lo")))      // true
print(bytesIndexOf(b1, bytesFromString("ll")))       // Some(2)

// Split/Join
parts = bytesSplit(joined, bytesFromString(" "))
back = bytesJoin(parts, bytesFromString("-"))
```

### Числовое кодирование

```rust
import "lib/bytes" (bytesEncodeInt, bytesDecodeInt, bytesEncodeFloat, bytesDecodeFloat)

// Int -> Bytes
b = bytesEncodeInt(256, 2)            // 2 bytes, big endian (default)
b2 = bytesEncodeInt(256, 2, "little") // little endian

// Bytes -> Int
n = bytesDecodeInt(b)
print(n)  // 256

// Float
bf = bytesEncodeFloat(3.14, 4)    // 32-bit float
f = bytesDecodeFloat(bf, 4)
print(f)  // 3.14...
```

---

## Сравнение: Bits vs Bytes

| Аспект | Bits | Bytes |
|--------|------|-------|
| Минимальная единица | 1 бит | 8 бит |
| Литералы | `#b"101"`, `#x"FF"`, `#o"7"` | `bytesFromString("...")` |
| Не кратно 8 | Да | Нет |
| Pattern matching | bitsExtract | — |
| Работа с протоколами | Отлично | Хорошо |
| Работа с файлами | Можно | Отлично |

---

## Use Cases

Bits:
- TCP/IP headers с bit flags
- Сжатие (Huffman codes)
- Видео кодеки (H.264 NAL units)
- Криптография
- Embedded протоколы

Bytes:
- Файловый I/O
- HTTP bodies
- JSON/XML данные
- Изображения
