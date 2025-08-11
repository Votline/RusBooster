package state

import (
	"fmt"
	"log"
	"time"
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

type UserState struct {
	IsChoosing              bool
	IsChecking              bool
	IsSetting               bool
	IsActive                bool
	Answer                  int
	OutputAnswer            string
	OutputTask              string
	Explanations            string
	CurrentPageOfAllWords   int
	CurrentPageOfFindWords  int
	CurrentPageOfGuide      int
	PartsOfAllWords         []string
	PartsOfFindWords        []string
	PartsOfGuide            []string
}

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "rusbooster_redis:6379",
		Password: "",
		DB:       0,
	})
}

func GetUserState(userId int64) (*UserState, error) {
	if rdb == nil {
		log.Printf("Ошибка, Redis не был инициализирован")
		InitRedis()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	key := fmt.Sprintf("user:%d", userId)
	data, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Такой записи в Redis не существует: %v", err)
			newState := &UserState{}
			if errState := SetUserState(newState, userId); errState != nil {
				log.Printf("Ошибка при создании новой записи: %v", errState)
				return nil, err
			}
			return newState, nil
		}
		log.Printf("Ошибка при получении данных из Redis: %v", err)
		return nil, err
	}
	var state UserState
	if errUn := json.Unmarshal(data, &state); errUn != nil {
		log.Printf("Ошибка при анмаршлинге данных из Redis: %v", errUn)
		return nil, errUn
	}
	return &state, nil
}

func SetUserState(state *UserState, userId int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if data, err := json.Marshal(state); err != nil {
		log.Printf("Произошла ошибка при работе с redis: %v", err)
		return err
	} else {
		key := fmt.Sprintf("user:%d", userId)
		if errSet := rdb.Set(ctx, key, data, 61*time.Minute).Err(); errSet != nil {
			log.Printf("Произошла ошибка при установке значения в Redis: %v", errSet)
			return errSet
		}
	}
	return nil
}
