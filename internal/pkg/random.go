package pkg

import (
	"math/rand"
)

func Random(min, max int) int {
	return (rand.Intn(max-min+1) + min) // Генерация случайного числа в диапазоне [min, max]
}
