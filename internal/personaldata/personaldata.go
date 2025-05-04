package personaldata

import "fmt"

// Информация о пользователе
type Personal struct {
	Name   string  // Имя
	Weight float64 // Вес
	Height float64 // Рост
}

func (p Personal) Print() {
	fmt.Printf("Имя: %s\nВес: %.2f кг.\nРост: %.2f м.\n\n", p.Name, p.Weight, p.Height)
}
