package state

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type UserState struct {
	IsChoosing             bool     `json:"is_choosing"`
	IsChecking             bool     `json:"is_checking"`
	IsSetting              bool     `json:"is_setting"`
	IsActive               bool     `json:"is_active"`
	CurrentTask            int      `json:"current_task"`
	Answer                 int      `json:"answer"`
	OutputTask             string   `json:"output_task"`
	OutputAnswer           string   `json:"output_answer"`
	LastQuestion           string   `json:"last_question"`
	Explanations           string   `json:"explanations"`
	PartsOfAllWords        []string `json:"parts_of_all_words"`
	CurrentPageOfAllWords  int      `json:"current_page_of_all_words"`
	PartsOfFindWords       []string `json:"parts_of_find_words"`
	CurrentPageOfFindWords int      `json:"current_page_of_find_words"`
	PartsOfGuide           []string `json:"parts_of_guide"`
	CurrentPageOfGuide     int      `json:"current_page_of_guide"`
}

var (
	rdb *redis.Client
)

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func GetUserState(log *zap.Logger, userId int64) (*UserState, error) {
	if rdb == nil {
		log.Error("Ошибка, Redis не был инициализирован")
		InitRedis()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	key := fmt.Sprintf("user:%d", userId)
	data, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			log.Error("Такой записи в Redis не существует", zap.Error(err))
			newState := &UserState{}
			if errState := SetUserState(log, newState, userId); errState != nil {
				log.Error("Ошибка при создании новой записи", zap.Error(errState))
				return nil, err
			}
			return newState, nil
		}
		log.Error("Ошибка при получении данных из Redis", zap.Error(err))
		return nil, err
	}
	var state UserState
	if errUn := json.Unmarshal(data, &state); errUn != nil {
		log.Error("Ошибка при анмаршлинге данных из Redis", zap.Error(errUn))
		return nil, errUn
	}
	return &state, nil
}

func SetUserState(log *zap.Logger, state *UserState, userId int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if data, err := json.Marshal(state); err != nil {
		log.Error("Произошла ошибка при работе с redis: ", zap.Error(err))
		return err
	} else {
		key := fmt.Sprintf("user:%d", userId)
		if errSet := rdb.Set(ctx, key, data, 61*time.Minute).Err(); errSet != nil {
			log.Error("Произошла ошибка при установке значения в Redis", zap.Error(errSet))
			return errSet
		}
	}
	return nil
}
