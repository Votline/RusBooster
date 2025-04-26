package stat

import (
	"RusBooster/internal/utils"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func createTable(log *zap.Logger) {
	db, _ := utils.GetDb(log)
	defer db.Close()

	usersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT NOT NULL UNIQUE PRIMARY KEY,
			current_task INTEGER DEFAULT 1,
			current_score INTEGER DEFAULT 0,
			worst_task_result INTEGER DEFAULT 0,
			worst_task_score INTEGER DEFAULT 0,
			best_task_result INTEGER DEFAULT 0,
			best_task_score INTEGER DEFAULT 0,
			streak INTEGER DEFAULT 0,
			last_active_date INTEGER DEFAULT 0,
			time_zone INTEGER DEFAULT 0,
			streak_freeze INTEGER DEFAULT 0
		);`
	if _, err := db.Exec(usersTable); err != nil {
		log.Fatal("Ошибка при попытке создать usersTable", zap.Error(err))
	}
	log.Info("Таблица users создана или уже существует")
}

func appendUser(log *zap.Logger, userId int64) {
	db, psql := utils.GetDb(log)
	defer db.Close()

	ins, args, err := psql.Insert("users").
		Columns("id").
		Values(userId).
		ToSql()
	if err != nil {
		log.Fatal("Ошибка при попытке добавить id пользователя в users", zap.Error(err))
	}
	if _, errExec := db.Exec(ins, args...); errExec != nil {
		log.Fatal("Ошибка при попытке выполнить запроса", zap.Error(errExec))
	}
}

func GetStatistic(log *zap.Logger, userName string, userId int64) string {
	var message string

	db, psql := utils.GetDb(log)
	defer db.Close()

	query, args, err := psql.Select("*").
		From("users").
		Where(sq.Eq{"id": userId}).
		ToSql()
	if err != nil {
		log.Error("Ошибка при попытке получить данные пользователя\n", zap.Error(err))
		return utils.GetReturnText(false)
	}

	rows, errX := db.Queryx(query, args...)
	if errX != nil {
		log.Error("Ошибка при попытке выполнить запроса\n", zap.Error(errX))
		createTable(log)
		return utils.GetReturnText(false)
	}

	if rows.Next() {
		result := make(map[string]interface{})
		if errMap := rows.MapScan(result); errMap != nil {
			log.Fatal("Ошибка при попытке считывания результата\n", zap.Error(errMap))
			return utils.GetReturnText(false)
		}
		currentTask, _ := utils.ToInt(result["current_task"])
		bestTaskResult, _ := utils.ToInt(result["best_task_result"])
		bestTaskScore, _ := utils.ToInt(result["best_task_score"])
		worstTaskResult, _ := utils.ToInt(result["worst_task_result"])
		worstTaskScore, _ := utils.ToInt(result["worst_task_score"])
		streak, _ := utils.ToInt(result["streak"])
		timeZone, _ := utils.ToInt(result["time_zone"])
		message = fmt.Sprintf("	👋Привет, %s!\nТекущее задание: №%d\nНаилучшая успеваимость: №%d, %d\nНаихудшая успеваимость: №%d, %d\nТы занимаешся уже: %d %s подряд!👏\nТекущий часовой пояс: %s\n", userName, currentTask, bestTaskResult, bestTaskScore, worstTaskResult, worstTaskScore, streak, utils.GetDayForm(streak), utils.GetTimeZoneForm(timeZone))
		return message
	} else if errRows := rows.Err(); errRows != nil {
		log.Error("Ошибка во время чтения rows\n", zap.Error(errRows))
	}
	log.Error("Ошибка при попытке чтения результата")
	appendUser(log, userId)
	return utils.GetReturnText(true)
}

func SetSomething(log *zap.Logger, userId int64, value int, request string) error {
	db, psql := utils.GetDb(log)
	defer db.Close()

	upd, args, err := psql.Update("users").
		Set(request, value).
		Where(sq.Eq{"id": userId}).
		ToSql()
	if err != nil {
		log.Error("Ошибка при попытке создать запрос", zap.Error(err))
		return err
	}
	if _, errExec := db.Exec(upd, args...); errExec != nil {
		log.Error("Ошибка при попытке выполнить запрос", zap.Error(err))
		return errExec
	}
	return nil
}

func GetSomething(log *zap.Logger, userId int64, request string) (int, error) {
	db, psql := utils.GetDb(log)
	defer db.Close()

	query, args, err := psql.Select(request).
		From("users").
		Where(sq.Eq{"id": userId}).
		ToSql()
	if err != nil {
		log.Error("Ошибка попытке создать запрос", zap.Error(err))
		return 0, err
	}
	var something int
	if errRow := db.QueryRow(query, args...).Scan(&something); errRow != nil {
		log.Error("Ошибка при попытке выполнить запрос", zap.Error(errRow))
		return 0, err
	}
	return something, nil
}

func GetAllChatIDs(log *zap.Logger) []int64 {
	db, psql := utils.GetDb(log)
	defer db.Close()
	ids := []int64{}

	query, _, err := psql.Select("id").
		From("users").
		ToSql()
	if err != nil {
		log.Warn("Ошибка при попытке создать запрос")
		return nil
	}
	rows, errX := db.Queryx(query)
	if errX != nil {
		log.Warn("Ошибка при попытке выполнить запрос")
		return nil
	}
	for rows.Next() {
		result := make(map[string]interface{})
		if errScan := rows.MapScan(result); errScan != nil {
			log.Warn("Ошибка при попытке сканировать результат")
			continue
		}
		id, _ := result["id"].(int64)
		ids = append(ids, id)
	}
	return ids
}
