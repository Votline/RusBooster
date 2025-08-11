package utils

import (
	"os"
	"fmt"
	"log"
	"sort"
	"time"
	"strconv"
	"unicode"
	"math/rand"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	sq "github.com/Masterminds/squirrel"
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

func ActionAfter(action func() error, delay int, errMessage string) {
	go func() {
		time.Sleep(time.Duration(delay) * time.Second)
		if err := action(); err != nil {
			log.Printf("%s: %v", errMessage, err)
		}
	}()
}

func HandleError(action func() error, errMsg string) string {
	if err := action(); err != nil {
		log.Printf("%s: %v", errMsg, err)
		return GetReturnText(false)
	}
	return ""
}

func GetDb() (*sqlx.DB, sq.StatementBuilderType) {
	dbUser := os.Getenv("POSTGRES_USER")
	dbName := os.Getenv("POSTGRES_DB")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbPass := os.Getenv("POSTGRES_PASSWORD")

	dsn := "host=postgres port=" + dbPort + " user=" + dbUser + " dbname=" + dbName + " password=" + dbPass + " sslmode=disable"

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных rsdb: %v", err)
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

func ToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("неподдерживаемый тип: %T", value)
	}
}

func GetDayForm(days int) string {
	if days == 1 {
		return "день"
	} else if days >= 2 && days <= 4 {
		return "дня"
	} else {
		return "дней"
	}
}
