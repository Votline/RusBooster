package guide

import (
	"fmt"
	"log"
	"strings"

	sq "github.com/Masterminds/squirrel"

	"RusBooster/internal/utils"
)

func createTable() {
	db, _ := utils.GetDb()
	defer db.Close()

	guidesTable := `
		CREATE TABLE IF NOT EXISTS guides (
			task_id INT PRIMARY KEY,
			guide TEXT
		);
	`
	if _, errExec := db.Exec(guidesTable); errExec != nil {
		log.Fatalf("Ошибка при попытке создать guidesTable: %v", errExec)
	}
}

func AppendGuide(taskId int, guideBody string) string {
	db, psql := utils.GetDb()
	defer db.Close()

	ins, args, err := psql.Insert("guides").
		Columns("task_id", "guide").
		Values(taskId, guideBody).
		ToSql()
	if err != nil {
		return fmt.Sprintf("Ошибка при попытке создать запрос: %v", err)
	}
	if _, errExec := db.Exec(ins, args...); errExec != nil {
		createTable()
		return fmt.Sprintf("Ошибка при попытке выполнить запрос: %v", errExec)
	}
	return fmt.Sprintf("Успешно! Гайд к заданию №%d добавлен!", taskId)
}

func DeleteGuide(taskId int) string {
	db, psql := utils.GetDb()
	defer db.Close()

	req, args, err := psql.Delete("guides").
		Where(sq.Eq{"task_id": taskId}).
		ToSql()
	if err != nil {
		return fmt.Sprintf("Ошибка при попытке создать запрос: %v", err)
	}
	if _, errExec := db.Exec(req, args...); errExec != nil {
		return fmt.Sprintf("Ошибка при выполнить запрос: %v", errExec)
	}
	return fmt.Sprintf("Успешно удалил гайд для задания №%d", taskId)
}

func ShowGuide(taskId int, targetPage *int, targetSlice *[]string) string {
	db, psql := utils.GetDb()
	defer db.Close()

	query, args, err := psql.Select("guide").
		From("guides").
		Where(sq.Eq{"task_id": taskId}).
		ToSql()
	if err != nil {
		log.Printf("Ошибка при попытке создать запрос: %v", err)
		return utils.GetReturnText(false)
	}

	var guideText string
	if errRow := db.QueryRowx(query, args...).Scan(&guideText); errRow != nil {
		log.Printf("Ошибка при попытке выполнить запрос или записать его: %v", errRow)
		return utils.GetReturnText(false)
	}
	
	*targetPage = 0
	*targetSlice = strings.Split(guideText, "@")
	
	return (*targetSlice)[0]
}
