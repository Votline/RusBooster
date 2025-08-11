# RusBooster

RusBooster — это Telegram-бот для подготовки к ЕГЭ по русскому языку, переписанный на Go. Проект вдохновлён приложением Duolingo, но ориентирован на подготовку к части ЕГЭ.

## Описание

Бот предлагает тренировочные задания по русскому языку в формате ЕГЭ, используя базу данных с перечнем слов и пояснений.

## Как начать использовать:
RusBooster уже доступен в Telegram! Чтобы начать:
1. Откройте Telegram
2. Найдите бота по username: [@RusBooster_bot]
3. Нажмите "Start" или напишите любое сообщение

### Основные функции:
- **Генерация заданий**: Пользователь выбирает тип задания, бот генерирует упражнение
- **Проверка ответов**: 
  - Верный ответ: +1 к статистике задания
  - Неверный ответ: -1 к статистике
- **Умные рекомендации**: Бот предлагает задания с наихудшими результатами при выборе задания
- **Статистика**: Реализована система streak (ежедневных цепочек) и напоминания

## Технологии
- **Go** (1.24.1) - основной язык
- **[PQ](https://github.com/lib/pq)** - драйвер для работы с PostgreSQL
- **[SQLX](https://github.com/jmoiron/sqlx)** - работа с базой данных 
- **[GO-Redis](https://github.com/redis/go-redis)** - кэширование
- **[Squirrel](https://github.com/Masterminds/squirrel)** - построение SQL-запросов
- **[Telebot](https://gopkg.in/telebot.v3)** - библиотека для работы с Telegram API

## Лицензия
Проект распространяется под лицензией [GNU AGPL v3](LICENSE).

Полные тексты лицензий доступны в [директории licenses](licenses/)

# RusBooster

RusBooster is a Telegram bot for preparing for the Russian Unified State Exam (USE), rewritten in Go. Inspired by Duolingo but specifically tailored for the test portion of the USE.

## Description

The bot provides Russian language practice exercises in USE format, using a database of words and explanations.

## Getting Started:
RusBooster is already available on Telegram! To begin:
1. Open Telegram
2. Search for the bot by username: [@RusBooster_bot]
3. Click "Start" or send any message

### Key Features:
- **Exercise Generation**: Users select exercise types, the bot generates practice questions
- **Answer Verification**:
  - Correct answer: +1 to exercise statistics
  - Incorrect answer: -1 to statistics
- **Smart Recommendations**: The bot suggests exercises with the weakest results when selecting tasks
- **Statistics**: Implemented streak system (daily chains) and reminders

## Technologies
- **Go** (1.24.1) - primary language
- **[PQ](https://github.com/lib/pq )** - a driver for working with PostgreSQL
- **[SQLX](https://github.com/jmoiron/sqlx)** - database operations 
- **[GO-Redis](https://github.com/redis/go-redis)** - caching
- **[Squirrel](https://github.com/Masterminds/squirrel)** - SQL query builder
- **[Telebot](https://gopkg.in/telebot.v3)** - Telegram API library

## License
This project is licensed under [GNU AGPL v3](LICENSE).

The full license texts are available in the [licenses directory](licenses/)

