package trainings

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
var errorUnknownActivity = errors.New("неизвестный тип тренировки")

// Структура реализует интерфейс
var _ actioninfo.DataParser = (*Training)(nil)

// Информация о тренировке
type Training struct {
	Steps                 int           // Количество шагов, проделанных за тренировку
	TrainingType          string        // Тип тренировки (бег или ходьба)
	Duration              time.Duration // Длительность тренировки
	personaldata.Personal               // Встроенная структура Personal из пакета personaldata (информация о пользователе)
}

// Парсит строку в поля структуры Training
// Принимает на вход строку (datastring) в формате: <кол-во шагов>,<вид активности>,<длительность активности>
// Возвращает:
//   - ошибка, если что-то пошло не так. (error)
func (t *Training) Parse(datastring string) (err error) {
	parseData := strings.Split(datastring, ",")
	if len(parseData) != 3 {
		return fmt.Errorf("%w: three comma-separated values are expected", errorFormateData)
	}

	steps, err := strconv.Atoi(parseData[0])
	if err != nil {
		return fmt.Errorf("%w. %w", errorFormateData, err)
	}
	if steps <= 0 {
		return fmt.Errorf("%w: first value in data must be > 0", errorFormateData)
	}

	if parseData[1] == "" {
		return fmt.Errorf("%w: second value in data is not defined", errorFormateData)
	}

	duration, err := time.ParseDuration(parseData[2])
	if err != nil {
		return fmt.Errorf("%w: %w", errorFormateData, err)
	}
	if duration <= 0 {
		return fmt.Errorf("%w: third value in the data must be > 0", errorFormateData)
	}
	// Все теперь можно записать в структуру
	t.Duration = duration
	t.Steps = steps
	t.TrainingType = parseData[1]

	return nil
}

// Возвращает сводную информацию об активности или ошибку
func (t Training) ActionInfo() (string, error) {

	distance := spentenergy.Distance(t.Steps, t.Height)               // Дистанция
	meanSpeed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration) // Ср.скорость

	var calories float64 = 0
	var err error

	switch strings.ToUpper(t.TrainingType) {
	case "ХОДЬБА":
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case "БЕГ":
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		err = fmt.Errorf("%w. %w", errorFormateData, errorUnknownActivity)
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType, t.Duration.Hours(), distance, meanSpeed, calories), nil
}
