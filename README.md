# A test assignment for Involta backend position

## The assignment

- [x] 1 Развернуть Reindexer в докере

- [x] 2 Написать на Go микросервис с постоянным подключением к реиндексеру посредством пакета от разработчика
- [x] 2.1 сделать вариативную конфигурацию
- [x] 2.1.1 конфигурация через локальный YAML файл
- [x] 2.1.2* конфигурация через Environment

- [x] 3 Проверки при запуске
- [x] 3.1 Проверка подключение к Reindexer 
- [x] 3.2 Проверка наличия необходимых коллекций
- [x] 3.2.1* В случае отсутствия необходимых коллекций их необходимо создать

- [x] 4 Микросервис должен создавать, редактировать, выводить информацию о списке имеющихся документов или заданного документа (CRUD)
- [x] 4.1* Для вывода списка предусмотреть параметры для пагинации и кол-ва выводимых документов

- [x] 5 Выдача одного документа должна кешироваться
- [x] 5.1* Кеш устаревает раз в 15 минут

- [x] 6 Документ содержит 2 уровня вложенности каждый из которых массив документов

- [x] 7* Массив документов первого уровня вложенности должен сортироваться по полю sort (целочисленное) (обратная сортировка)

- [ ] 8* Каждый документ перед отправкой ответа обрабатывается в отдельном потоке при этом общая сортировка не должна пострадать (обработка подразумевает исключение одного или нескольких полей из каждого уровня документа) 

\* - "Продвинутый" уровень

## Run
```
```
