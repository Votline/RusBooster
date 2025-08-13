package stat

import (
	"log"
	
	"RusBooster/internal/utils"
)

func appendUser(userId int64) {
	db, psql := utils.GetDb()
	defer db.Close()

	ins, args, err := psql.Insert("users").
		Columns("id").
		Values(userId).
		ToSql()
	if err != nil {
		log.Fatalf("Ошибка при попытке добавить id пользователя в users: %v", err)
	}
	if _, errExec := db.Exec(ins, args...); errExec != nil {
		log.Fatalf("Ошибка при попытке выполнить запроса: %v", errExec)
	}
}

func getDayForm(days int) string {
	if days == 1 {
		return "день"
	} else if days >= 2 && days <= 4 {
		return "дня"
	} else {
		return "дней"
	}
}
