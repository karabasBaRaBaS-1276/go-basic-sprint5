package actioninfo

import (
	"fmt"
	"log"
)

// Пакет реализует вывод общей информации обо всех тренировках и прогулках
type DataParser interface {
	Parse(datastring string) (err error) // Парсит строку с данными об активности
	ActionInfo() (string, error)         // Сводная информация об активности
}

func Info(dataset []string, dp DataParser) {
	for _, datastring := range dataset {
		err := dp.Parse(datastring)
		if err != nil {
			log.Println(err)
			continue
		}
		result, err := dp.ActionInfo()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(result)
	}
}
