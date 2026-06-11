# pipelineTest

Демонстрационное приложение для библиотеки [pipelineLibrary](https://github.com/mr-tortilla/pipelineLibrary).

Рекурсивно обходит заданную директорию, вычисляет MD5-хеш для каждого файла
и выводит результат в stdout. Ошибки выводятся в stderr. Вычисление хешей параллелизировано.

## Требования

- Go 1.25+

## Установка

```bash
git clone https://github.com/mr-tortilla/pipelineTest
cd pipelineTest
go mod tidy
```

## Запуск

```bash
go run ./cmd/md5walk/ <directory> [parallelism]
```

**Аргументы:**
- `<directory>` - путь к директории для обхода (обязательный)
- `[parallelism]` - количество параллельных воркеров для вычисления MD5 (по умолчанию: 10)

## Примеры

```bash
# обход текущей директории с параллелизмом по умолчанию (10)
go run ./cmd/md5walk/ .

# обход C:\Users с 5 воркерами
go run ./cmd/md5walk/ C:\Users 5
```

**Пример вывода:**

```bash
a3f5c2d1e4b6789012345678abcdef01  ./cmd/md5walk/main.go
b7d1e4f2a3c5678901234567bcdef012  ./cmd/md5walk/walk_node.go
ERR: open ./locked.txt: permission denied
```

## Остановка

Для остановки нажмите `Ctrl+C` — пайплайн корректно завершит работу.

## Архитектура

Приложение построено на четырёх нодах:

- **WalkNode** - рекурсивно обходит директорию и пишет пути файлов в канал
- **HashNode** - читает путь из канала, вычисляет MD5-хеш, ошибки пишет в канал ошибок. Запускается в N экземплярах параллельно
- **PrintNode** - читает результаты и выводит их в stdout
- **ErrNode** - читает ошибки и выводит их в stderr

```
WalkNode -> [paths] -> HashNode x N -> [results] -> PrintNode -> stdout -> [errs] -> ErrNode -> stderr
```

Параллелизм HashNode реализован через один общий канал `paths` -
Go автоматически распределяет файлы между свободными воркерами.
Библиотека сама закрывает каналы когда все писатели завершились.

## Последовательность завершения

1. WalkNode обошла все файлы -> завершается
2. Библиотека закрывает `paths`
3. HashNode читает `ok=false` из `paths` -> все воркеры завершаются
4. Библиотека закрывает `results` и `errs`
5. PrintNode и ErrNode читают `ok=false` -> завершаются
6. `p.Wait()` разблокируется -> программа завершается