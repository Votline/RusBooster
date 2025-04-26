package utils

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"math/rand"
	"sort"
	"strconv"
	"time"
	"unicode"
)

func ContainsRune(text string, item string) bool {
	for _, r := range text {
		if string(r) == item {
			return true
		}
	}
	return false
}

func Contains(slice []string, item string) bool {
	for _, value := range slice {
		if value == item {
			return true
		}
	}
	return false
}

func FindUpperCase(word string) rune {
	for _, r := range word {
		if unicode.IsUpper(r) {
			return r
		}
	}
	return 0
}

func AllEqual(letters []rune) bool {
	if len(letters) == 0 {
		return false
	}
	for _, r := range letters {
		if r != letters[0] {
			return false
		}
	}
	return true
}

func getDigitHash(s string) string {
	runes := []rune(s)
	sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
	return string(runes)
}

func CheckEquals(correct, verifiable string) bool {
	correct = getDigitHash(correct)
	verifiable = getDigitHash(verifiable)
	return correct == verifiable
}

func SumDigits(number int) int {
	sum := 0
	for number > 0 {
		digit := number % 10
		sum += digit
		number /= 10
	}
	return sum
}

func ActionAfter(log *zap.Logger, action func() error, delay int, errMessage string) {
	go func() {
		time.Sleep(time.Duration(delay) * time.Second)
		if err := action(); err != nil {
			log.Error(errMessage, zap.Error(err))
		}
	}()
}

func HandleError(log *zap.Logger, action func() error, errMsg string) string {
	if err := action(); err != nil {
		log.Error(errMsg, zap.Error(err))
		return GetReturnText(false)
	}
	return ""
}

func GetDb(log *zap.Logger) (*sqlx.DB, sq.StatementBuilderType) {
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres dbname=rsdb sslmode=disable")
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных rsdb\n", zap.Error(err))
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return db, psql
}

func GetTimeZoneForm(timeZone int) string {
	if timeZone > 0 {
		return fmt.Sprintf("+%d от МСК", timeZone)
	} else if timeZone < 0 {
		return fmt.Sprintf("-%d от МСК", timeZone)
	} else {
		return "МСК"
	}
}

func GetReturnText(newUser bool) string {
	rand.Seed(time.Now().UnixNano())
	var newUserPhrases = []string{
		"👋 Привет! Похоже, ты у нас впервые. Создаём профиль... Нажми кнопку ещё раз 🚀",
		"Добро пожаловать! Сейчас всё настроим. Просто ещё раз нажми кнопку 👇",
		"Профиль почти готов! Жми кнопку ещё раз, чтобы запустить 🚀",
	}
	var errorPhrases = []string{
		"Что-то пошло не так... Попробуй ещё раз 🔁",
		"Упс! Произошла ошибка. Повтори попытку 👇",
		"Мы столкнулись с технической проблемой 😢 Попробуй ещё раз",
		"Кажется, всё не так гладко 🧩 Попробуй снова — это должно помочь!",
		"Ошибка при обработке данных. Но ничего! Просто попробуй ещё раз 🤝",
	}
	if newUser {
		return newUserPhrases[rand.Intn(len(newUserPhrases))]
	} else {
		return errorPhrases[rand.Intn(len(errorPhrases))]
	}
}

func ToInt(val interface{}) (int, error) {
	switch value := val.(type) {
	case int:
		return value, nil
	case int64:
		return int(value), nil
	case float64:
		return int(value), nil
	case []uint8:
		if value, err := strconv.Atoi(string(value)); err != nil {
			return 0, err
		} else {
			return value, nil
		}
	case string:
		if value, err := strconv.Atoi(value); err != nil {
			return 0, err
		} else {
			return value, nil
		}
	default:
		return 0, nil
	}
}

func GetDayForm(num interface{}) string {
	var number int

	switch v := num.(type) {
	case int:
		number = v
	case int64:
		number = int(v)
	default:
		return "дней"
	}
	lastDigit := number % 10
	lastTwoDigits := number % 100

	if number == 0 {
		return "дней"
	} else if lastDigit == 1 && lastTwoDigits != 11 {
		return "день"
	} else if lastDigit == 2 && lastTwoDigits != 4 || !(lastDigit >= 12 && lastTwoDigits <= 14) {
		return "дня"
	} else {
		return "дней"
	}
}
