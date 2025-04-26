package guide

import (
	"RusBooster/internal/utils"
	"fmt"
	"strings"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

func createTable(log *zap.Logger) {
	db, _ := utils.GetDb(log)
	defer db.Close()

	guidesTable := `
		CREATE TABLE IF NOT EXISTS guides (
			task_id INT PRIMARY KEY,
			guide TEXT
		);
	`
	if _, errExec := db.Exec(guidesTable); errExec != nil {
		log.Fatal("Ошибка при попытке создать guidesTable", zap.Error(errExec))
	}
}

func AppendGuide(log *zap.Logger, taskId int, guideBody string) string {
	db, psql := utils.GetDb(log)
	defer db.Close()

	ins, args, err := psql.Insert("guides").
		Columns("task_id", "guide").
		Values(taskId, guideBody).
		ToSql()
	if err != nil {
		return fmt.Sprintf("Ошибка при попытке создать запрос \n%v", zap.Error(err))
	}
	if _, errExec := db.Exec(ins, args...); errExec != nil {
		createTable(log)
		return fmt.Sprintf("Ошибка при попытке выполнить запрос \n%v", zap.Error(errExec))
	}
	return fmt.Sprintf("Успешно! Гайд к заданию №%d добавлен!", taskId)
}

func DeleteGuide(log *zap.Logger, taskId int) string {
	db, psql := utils.GetDb(log)
	defer db.Close()

	req, args, err := psql.Delete("guides").
		Where(sq.Eq{"task_id": taskId}).
		ToSql()
	if err != nil {
		return fmt.Sprintf("Ошибка при попытке создать запрос \n%v", zap.Error(err))
	}
	if _, errExec := db.Exec(req, args...); errExec != nil {
		return fmt.Sprintf("Ошибка при выполнить запрос \n%v", zap.Error(errExec))
	}
	return fmt.Sprintf("Успешно удалил гайд для задания №%d", taskId)
}

func ShowGuide(log *zap.Logger, taskId int, targetPage *int, targetSlice *[]string) string {
	db, psql := utils.GetDb(log)
	defer db.Close()

	query, args, err := psql.Select("guide").
		From("guides").
		Where(sq.Eq{"task_id": taskId}).
		ToSql()
	if err != nil {
		log.Error("Ошибка при попытке создать запрос: \n", zap.Error(err))
		return utils.GetReturnText(false)
	}

	var guide string
	if errRow := db.QueryRow(query, args...).Scan(&guide); errRow != nil {
		log.Error("Ошибка при попытке выполнить запрос или записать его \n", zap.Error(errRow))
		return utils.GetReturnText(false)
	}
	
	*targetPage = 0	
	*targetSlice = strings.Split(guide, "@")
	
	return (*targetSlice)[0] 
}
