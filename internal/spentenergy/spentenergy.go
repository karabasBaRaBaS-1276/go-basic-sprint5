package spentenergy

import (
	"errors"
	"fmt"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе.
)

// Неправильный формат данных
var errorFormateData = errors.New("wrong data format")

// Возвращает количество калорий, потраченных при ходьбе или ошибку.
// Принимает на вход:
//
//	steps    - количество шагов
//	weight   - вес пользователя
//	height   - рост пользователя
//	duration - продолжительность активности
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	runningSpentCalories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}
	return runningSpentCalories * walkingCaloriesCoefficient, nil
}

// Возвращает количество калорий, потраченных при беге или ошибку.
// Принимает на вход:
//
//	steps    - количество шагов
//	weight   - вес пользователя
//	height   - рост пользователя
//	duration - продолжительность активности
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("%w: steps must be > 0", errorFormateData)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("%w: weight must be > 0", errorFormateData)
	}
	if height <= 0 {
		return 0, fmt.Errorf("%w: height must be > 0", errorFormateData)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("%w: duration must be > 0", errorFormateData)
	}

	meanSpeed := MeanSpeed(steps, height, duration)
	if meanSpeed <= 0 {
		return 0, fmt.Errorf("%w: error in calculating average speed", errorFormateData)
	}

	durationInMinutes := duration.Minutes()
	return (weight * meanSpeed * durationInMinutes) / minInH, nil
}

// Возвращает среднюю скорость в км/ч
//
// Принимает на вход:
//
//	steps    - количество шагов
//	height   - рост пользователя
//	duration - продолжительность активности
func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distance := Distance(steps, height)
	return distance / duration.Hours()
}

// Возвращает дистанцию в километрах
//
// Принимает на вход:
//
//	steps  - количество шагов
//	height - рост пользователя
func Distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient
	return (stepLen * float64(steps)) / mInKm
}
