package words

import (
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
	sq "github.com/Masterminds/squirrel"

	"RusBooster/internal/utils"
)

func createTable() {
	db, _ := utils.GetDb()
	defer db.Close()

	wordsTable := `
		CREATE TABLE IF NOT EXISTS words (
			word TEXT NOT NULL,
			explanation TEXT NOT NULL,
			task_id INT NOT NULL
		)
	`
	if _, err := db.Exec(wordsTable); err != nil {
		log.Fatalf("Ошибка при попытке создать wordsTable: %v", err)
	}
	log.Printf("Таблица words создана или уже существует")
}

func GetAll(wordsForTask map[string]string, currentTask int) error {
	db, psql := utils.GetDb()
	defer db.Close()

	query, args, err := psql.Select("word", "explanation").
		From("words").
		Where(sq.Eq{"task_id": currentTask}).
		ToSql()
	if err != nil {
		log.Printf("Ошибка при попытке создать запрос: %v", err)
		return err
	}

	rows, errX := db.Queryx(query, args...)
	if errX != nil {
		log.Printf("Ошибка при попытке выполнить запрос: %v", errX)
		createTable()
		return errX
	}

	var word, explanation string
	for rows.Next() {
		if errScan := rows.Scan(&word, &explanation); errScan != nil {
			log.Printf("Ошибка при попытке считывания результата: %v", errScan)
			return errScan
		}
		wordsForTask[word] = explanation
	}
	return nil
}

func SetSomething(word string, explanation string, taskId int) string {
	if _, err := FindWord(taskId, word); err == nil {
		return fmt.Sprintf("Слово '%s' уже существует в бд", word)
	}
	db, psql := utils.GetDb()
	defer db.Close()

	ins, args, err := psql.Insert("words").
		Columns("word", "explanation", "task_id").
		Values(word, explanation, taskId).
		ToSql()
	if err != nil {
		return "Ошибка при попытке создать запрос"
	}
	if _, errExec := db.Exec(ins, args...); errExec != nil {
		return "Ошибка при попытке выполнить запрос. Слово: " + word
	}
	return fmt.Sprintf("Успешно! Слово: %s\nи пояснение: %s\nдобавлены к заданию: %d", word, explanation, taskId)
}

func DeleteWord(taskId int, word string) string {
	db, psql := utils.GetDb()
	defer db.Close()

	req, args, err := psql.Delete("words").
		Where(sq.Eq{"task_id": taskId, "word": word}).
		ToSql()
	if err != nil {
		return fmt.Sprintf("Ошибка при попытке создать запрос: %v", err)
	}
	result, errExec := db.Exec(req, args...)
	if errExec != nil {
		return fmt.Sprintf("Ошибка при попытке выполнить запрос: %v", errExec)
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Sprintf("Слова %s нет в базе данных", word)
	}
	return fmt.Sprintf("Успешно удалил слово: %s", word)
}

func FindWord(taskId int, word string) (string, error){
	db, psql := utils.GetDb()
	defer db.Close()

	query, args, err := psql.Select("explanation").
		From("words").
		Where(sq.Eq{"task_id": taskId, "word": word}).
		ToSql()
	if err != nil {
		return fmt.Sprintf("Ошибка при попытке создать запрос: %v", err), err
	}
	var explanation string
	if errRow := db.QueryRowx(query, args...).Scan(&explanation); errRow != nil {
		return fmt.Sprintf("Ошибка при попытке создать запрос: %v", errRow), errRow
	}
	return explanation, nil
}

func ShowAllWords(userId int64, taskId int, targetSlice *[]string, targetPage *int, filter bool, neededLetter string) string {
	allWords := make(map[string]string)
	if err := GetAll(allWords, taskId); err != nil {
		return fmt.Sprintf("Ошибка при попытке заполнить allWords: %v", err)
	}

	var cnt int = 0
	var message string = ""
	for word, explanation := range allWords {
		if filter == true && string(utils.FindUpperCase(explanation)) != neededLetter {
			continue
		}
		if cnt < 5 {
			cnt++
			message += fmt.Sprintf("Слово: %s\nОбъяснение: %s\n\n", word, explanation)
		} else {
			cnt = 0
			message += "@"
		}
	}

	*targetSlice = strings.Split(message, "@")
	*targetPage = 0
	return (*targetSlice)[0]
}
