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
- **Telebot** - библиотека для работы с Telegram API
- **SQLx** - работа с базой данных
- **Redis** - кэширование
- **Zap** - логирование
- **Squirrel** - построение SQL-запросов

## Установка и запуск
1. Клонировать репозиторий:
  - `git clone https://github.com/yourusername/RusBooster.git`
2. Установить зависимости:
  - `go mod download`
3. Создать URL и скопировать его:
  - `lt --port 8080`
4. Запустить бота с URL:
  - `WEBHOOK_URL="URL" go run main.go`

## Лицензия
Проект распространяется под лицензией [GNU AGPL v3](LICENSE).



# RusBooster

RusBooster is a Telegram bot for preparing for the Russian Unified State Exam (ЕГЭ), rewritten in Go. Inspired by Duolingo but specifically tailored for the test portion of the ЕГЭ.

## Description

The bot provides Russian language practice exercises in ЕГЭ format, using a database of words and explanations.

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
- **Telebot** - Telegram API library
- **SQLx** - database operations
- **Redis** - caching
- **Zap** - logging
- **Squirrel** - SQL query builder

## Installation and Setup
1. Clone the repository:
  - `git clone https://github.com/yourusername/RusBooster.git`
2. Install dependencies:
  - `go mod download`
3. Create and copy URL:
  - `lt --port 8080`
4. Run the bot with URL:
  - `WEBHOOK_URL="URL" go run main.go`

## License
This project is licensed under [GNU AGPL v3](LICENSE).
