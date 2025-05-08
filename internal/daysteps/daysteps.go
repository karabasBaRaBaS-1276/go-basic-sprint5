package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/actioninfo"
	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// Неправильный формат данных
var errorFormateData = errors.New("wrong data format")

// Структура реализует интерфейс
var _ actioninfo.DataParser = (*DaySteps)(nil)

// Информация о прогулках
type DaySteps struct {
	Steps                 int           // Количество шагов
	Duration              time.Duration // Длительность прогулки
	personaldata.Personal               // Встроенная структура Personal из пакета personaldata (информация о пользователе)
}

// Парсит строку в поля структуры DaySteps
// Принимает на вход строку (datastring) в формате: <кол-во шагов>,<длительность активности>
// Возвращает:
//   - ошибка, если что-то пошло не так. (error)
func (ds *DaySteps) Parse(datastring string) (err error) {
	parseData := strings.Split(datastring, ",")
	if len(parseData) != 2 {
		return fmt.Errorf("%w: two comma-separated values are expected", errorFormateData)
	}

	steps, err := strconv.Atoi(parseData[0])
	if err != nil {
		return fmt.Errorf("%w: %v", errorFormateData, err)
	}
	if steps <= 0 {
		return fmt.Errorf("%w: first value in data must be > 0", errorFormateData)
	}

	duration, err := time.ParseDuration(parseData[1])
	if err != nil {
		return fmt.Errorf("%w: %v", errorFormateData, err)
	}
	if duration <= 0 {
		return fmt.Errorf("%w: second value in data must be > 0", errorFormateData)
	}
	// Все теперь можно записать в структуру
	ds.Steps = steps
	ds.Duration = duration

	return nil
}

// Возвращает сводную информацию об активности или ошибку
func (ds DaySteps) ActionInfo() (string, error) {
	if ds.Steps <= 0 {
		return "", fmt.Errorf("%w", errorFormateData)
	}
	distance := spentenergy.Distance(ds.Steps, ds.Height)

	caloriesDay, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Weight, ds.Height, ds.Duration)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", ds.Steps, distance, caloriesDay), nil
}
