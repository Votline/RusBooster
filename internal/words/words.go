package words

import (
	"RusBooster/internal/utils"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"strings"
)

func createTable(log *zap.Logger) {
	db, _ := utils.GetDb(log)
	defer db.Close()

	wordsTable := `
		CREATE TABLE IF NOT EXISTS words (
			word TEXT NOT NULL,
			explanation TEXT NOT NULL,
			task_id INT NOT NULL
		)
	`
	if _, err := db.Exec(wordsTable); err != nil {
		log.Fatal("Ошибка при попытке создать wordsTable", zap.Error(err))
	}
	log.Info("Таблица words создана или уже существует")
}

func GetAll(log *zap.Logger, wordsForTask map[string]string, currentTask int) error {
	db, psql := utils.GetDb(log)
	defer db.Close()

	query, args, err := psql.Select("word", "explanation").
		From("words").
		Where(sq.Eq{"task_id": currentTask}).
		ToSql()
	if err != nil {
		log.Error("Ошибка при попытке создать запрос", zap.Error(err))
		return err
	}

	rows, errX := db.Queryx(query, args...)
	if errX != nil {
		log.Error("Ошибка при попытке выполнить запрос", zap.Error(errX))
		createTable(log)
		return errX
	}

	var word, explanation string
	for rows.Next() {
		if errScan := rows.Scan(&word, &explanation); errScan != nil {
			log.Error("Ошибка при попытке считывания результата")
			return errScan
		}
		wordsForTask[word] = explanation
	}
	return nil
}

func SetSomething(log *zap.Logger, word string, explanation string, taskId int) string {
	if _, err := FindWord(log, taskId, word); err == nil {
		return fmt.Sprintf("Слово '%s' уже существует в бд", word)
	}
	db, psql := utils.GetDb(log)
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

func DeleteWord(log *zap.Logger, taskId int, word string) string {
	db, psql := utils.GetDb(log)
	defer db.Close()

	req, args, err := psql.Delete("words").
		Where(sq.Eq{"task_id": taskId, "word": word}).
		ToSql()
	if err != nil {
		return fmt.Sprintf("Ошибка при попытке создать запрос \n%v", zap.Error(err))
	}
	result, errExec := db.Exec(req, args...)
	if errExec != nil {
		return fmt.Sprintf("Ошибка при попытке выполнить запрос \n%v", zap.Error(errExec))
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Sprintf("Слова %s нет в базе данных", word)
	}
	return fmt.Sprintf("Успешно удалил слово: %s", word)
}

func FindWord(log *zap.Logger, taskId int, word string) (string, error){
	db, psql := utils.GetDb(log)
	defer db.Close()

	query, args, err := psql.Select("word", "explanation").
		From("words").
		Where(sq.Eq{"task_id": taskId}).
		Where(sq.Eq{"word": word}).
		ToSql()
	if err != nil {
		return fmt.Sprintf("Ошибка при попытке создать запрос \n%v", zap.Error(err)), err
	}
	var neededWord, neededExplanation string
	if errRow := db.QueryRow(query, args...).Scan(&neededWord, &neededExplanation); errRow != nil {
		return fmt.Sprintf("Ошибка при попытке создать запрос \n%v", zap.Error(errRow)), errRow
	}
	return fmt.Sprintf("Нужное слово: %v\nПояснение: %v", neededWord, neededExplanation), nil
}

func ShowAllWords(log *zap.Logger, userId int64, taskId int, targetSlice *[]string, targetPage *int, filter bool, neededLetter string) string {
	allWords := make(map[string]string)
	if err := GetAll(log, allWords, taskId); err != nil {
		return fmt.Sprintf("Ошибка при попытке заполнить allWords \n%v", zap.Error(err))
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
