package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// go test -v homework_test.go

// логика для наглядности и усвоения материала
// ((number >> 24) & 0x000000FF) Сдвигаем число на 3 байта (24 бита) вправо, оставляем младший байт с помощью маски
// ((number >> 8) & 0x0000FF00) Сдвигаем число на 1 байт вправо, оставляем байт на 2й ппозиции
// ((number << 8) & 0x00FF0000) Сдвигаем число на 1 байт влево, оставляем байт на 3й позиции
// ((number << 24) & 0xFF000000) Cдвигаем число на 3 байта влево, оставляем старший байт
// побитовым ИЛИ склеиваем результат
func ToLittleEndian(number uint32) uint32 {
	return ((number >> 24) & 0x000000FF) | ((number >> 8) & 0x0000FF00) | ((number << 8) & 0x00FF0000) | ((number << 24) & 0xFF000000)
}

func TestСonversion(t *testing.T) {
	tests := map[string]struct {
		number uint32
		result uint32
	}{
		"test case #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #4": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}
