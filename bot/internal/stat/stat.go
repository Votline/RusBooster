package stat

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	sq "github.com/Masterminds/squirrel"

	"RusBooster/internal/utils"
)

func createTable() {
	db, _ := utils.GetDb()
	defer db.Close()

	usersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id                BIGINT NOT NULL UNIQUE PRIMARY KEY,
			current_task      INTEGER DEFAULT 1,
			current_score     INTEGER DEFAULT 0,
			worst_task_result INTEGER DEFAULT 0,
			worst_task_score  INTEGER DEFAULT 0,
			best_task_result  INTEGER DEFAULT 0,
			best_task_score   INTEGER DEFAULT 0,
			streak            INTEGER DEFAULT 0,
			last_active_date  INTEGER DEFAULT 0,
			time_zone         INTEGER DEFAULT 0,
			streak_freeze     INTEGER DEFAULT 0
		);`
	if _, err := db.Exec(usersTable); err != nil {
		log.Fatalf("Ошибка при попытке создать usersTable: %v", err)
	}
	log.Printf("Таблица users создана или уже существует")
}

func GetStatistic(userName string, userId int64) string {
	var message string

	db, psql := utils.GetDb()
	defer db.Close()

	query, args, err := psql.Select("*").
		From("users").
		Where(sq.Eq{"id": userId}).
		ToSql()
	if err != nil {
		log.Printf("Ошибка при попытке получить данные пользователя: %v", err)
		return utils.GetReturnText(false)
	}

	rows, errX := db.Queryx(query, args...)
	if errX != nil {
		log.Printf("Ошибка при попытке выполнить запроса: %v", errX)
		createTable()
		return utils.GetReturnText(false)
	}

	if rows.Next() {
		result := make(map[string]interface{})
		if errMap := rows.MapScan(result); errMap != nil {
			log.Fatalf("Ошибка при попытке считывания результата: %v", errMap)
			return utils.GetReturnText(false)
		}
		currentTask, _     := utils.ToInt(result["current_task"])
		bestTaskResult, _  := utils.ToInt(result["best_task_result"])
		bestTaskScore, _   := utils.ToInt(result["best_task_score"])
		worstTaskResult, _ := utils.ToInt(result["worst_task_result"])
		worstTaskScore, _  := utils.ToInt(result["worst_task_score"])
		streak, _          := utils.ToInt(result["streak"])
		timeZone, _        := utils.ToInt(result["time_zone"])
		message = fmt.Sprintf("	👋Привет, %s!\nТекущее задание: №%d\nНаилучшая успеваимость: №%d, %d\nНаихудшая успеваимость: №%d, %d\nТы занимаешся уже: %d %s подряд!👏\nТекущий часовой пояс: %s\n", userName, currentTask, bestTaskResult, bestTaskScore, worstTaskResult, worstTaskScore, streak, getDayForm(streak), utils.GetTimeZoneForm(timeZone))
		return message
	} else if errRows := rows.Err(); errRows != nil {
		log.Printf("Ошибка во время чтения rows: %v", errRows)
	}
	appendUser(userId)
	return utils.GetReturnText(true)
}

func SetSomething(userId int64, value int, request string) error {
	db, psql := utils.GetDb()
	defer db.Close()

	query, args, err := psql.Update("users").
		Set(request, value).
		Where(sq.Eq{"id": userId}).
		ToSql()
	if err != nil {
		log.Printf("Ошибка при попытке создать запрос: %v", err)
		return err
	}
	if _, err := db.Exec(query, args...); err != nil {
		log.Printf("Ошибка при попытке выполнить запрос: %v", err)
		return err
	}
	return nil
}
func SetExclude(userID int64, newValues map[string]interface{}, exclude ...string) error {
	db, psql := utils.GetDb()
	defer db.Close()

	excludeSet := make(map[string]struct{}, len(exclude))
	for _, col := range exclude {
		excludeSet[col] = struct{}{}
	}

	query := psql.Update("users").Where(sq.Eq{"id": userID})

	for col, val := range newValues {
		if _, skip := excludeSet[col]; skip {
			continue
		}
		query.Set(col, val)
	}

	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Printf("Ошибка при попытке создать запрос: %v", err)
		return err
	}

	if _, err := db.Exec(sqlStr, args...); err != nil {
		log.Printf("Ошибка при попытке выполнить запрос: %v", err)
		return err
	}

	return nil
}

func GetSomething(userId int64, request string) (int, error) {
	db, psql := utils.GetDb()
	defer db.Close()

	query, args, err := psql.Select(request).
		From("users").
		Where(sq.Eq{"id": userId}).
		ToSql()
	if err != nil {
		log.Printf("Ошибка при попытке создать запрос: %v", err)
		return 0, err
	}
	var result int
	if errRow := db.QueryRowx(query, args...).Scan(&result); errRow != nil {
		log.Printf("Ошибка при попытке выполнить запрос: %v", errRow)
		return 0, errRow
	}
	return result, nil
}
func GetExclude(userID int64, exclude ...string) (map[string]interface{}, error) {
	db, psql := utils.GetDb()
	defer db.Close()

	query, args, err := psql.Select("*").
		From("users").
		Where(sq.Eq{"id": userID}).
		ToSql()
	if err != nil {
		log.Printf("Ошибка при попытке создать запрос: %v", err)
		return nil, err
	}

	row := db.QueryRowx(query, args...)
	result := make(map[string]interface{})
	if err := row.MapScan(result); err != nil {
		log.Printf("Ошибка при попытке сканирования результата: %v", err)
		return nil, err
	}

	excludeSet := make(map[string]struct{}, len(exclude))
	for _, col := range exclude {
		excludeSet[col] = struct{}{}
	}
	for col := range excludeSet {
		delete(result, col)
	}
	return result, nil
}

func GetAllChatIDs() []int64 {
	db, psql := utils.GetDb()
	defer db.Close()

	query, args, err := psql.Select("id").
		From("users").
		ToSql()
	if err != nil {
		log.Printf("Ошибка при попытке создать запрос: %v", err)
		return []int64{}
	}

	rows, errX := db.Queryx(query, args...)
	if errX != nil {
		log.Printf("Ошибка при попытке выполнить запрос: %v", errX)
		return []int64{}
	}

	var userIds []int64
	for rows.Next() {
		var userId int64
		if errScan := rows.Scan(&userId); errScan != nil {
			log.Printf("Ошибка при попытке считывания результата: %v", errScan)
			continue
		}
		userIds = append(userIds, userId)
	}
	return userIds
}
