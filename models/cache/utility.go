package cache

import (
	"strings"
	"time"
)

func StoreUtility(meals []string, timeslots string,periods string, cachetime time.Duration) error {
	contents := ""

	// store the meals to the redis cache
	for _, meal := range meals {
		contents += (meal+"/")
	}
	err := Save(timeslots+"/"+periods, contents, cachetime)
	if err != nil {
		return err
	}

	return nil
}
func GetUtility(timeslots string, periods string) ([]string, error) {
	// get the meals from the redis cache
	value, err := Get(timeslots + "/" + periods)
	if err != nil {
		return nil, err
	}
	// meals
	meals := strings.Split(value, "/")
	meals = meals[:len(meals)-1]
	return meals, nil
}
