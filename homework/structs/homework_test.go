package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Option func(*GamePerson)

func charTo7Bit(c byte) byte {
	if c >= 32 && c <= 126 {
		return c - 32
	}
	return 0
}

func bit7ToChar(v byte) byte {
	return v + 32
}

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		person.name = [38]byte{}
		person.name[0] = byte(len(name))

		var bitPos int

		for _, ch := range name {
			val := charTo7Bit(byte(ch))

			for i := 6; i >= 0; i-- {
				bit := (val >> i) & 1

				byteIdx := 1 + (bitPos / 8)
				bitIdx := 7 - (bitPos % 8)

				if bit == 1 {
					person.name[byteIdx] |= (1 << bitIdx)
				}
				bitPos++
			}
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (maskMana << shiftMana)) | (uint64(mana) << shiftMana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (maskHealth << shiftHealth)) | (uint64(health) << shiftHealth)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (maskRespect << shiftRespect)) | (uint64(respect) << shiftRespect)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (maskStr << shiftStr)) | (uint64(strength) << shiftStr)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (maskExp << shiftExp)) | (uint64(experience) << shiftExp)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (maskLevel << shiftLevel)) | (uint64(level) << shiftLevel)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats ^= (1 << shiftHouse)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats ^= (1 << shiftGun)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats ^= (1 << shiftFamily)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ (maskType << shiftType)) | (uint64(personType) << shiftType)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x, y, z int32
	gold    uint32
	stats   uint64
	name    [38]byte
}

const (
	maskLevel   = 0x0F
	maskExp     = 0x0F
	maskStr     = 0x0F
	maskRespect = 0x0F
	maskHealth  = 0x03FF
	maskMana    = 0x03FF
	maskType    = 0x03

	shiftLevel   = 0
	shiftExp     = 4
	shiftStr     = 8
	shiftRespect = 12
	shiftHealth  = 16
	shiftMana    = 26
	shiftFamily  = 36
	shiftGun     = 37
	shiftHouse   = 38
	shiftType    = 39
)

func NewGamePerson(options ...Option) GamePerson {
	gamePerson := GamePerson{}

	for _, option := range options {
		option(&gamePerson)
	}

	return gamePerson
}

func (p *GamePerson) Name() string {
	length := int(p.name[0])
	if length == 0 {
		return ""
	}

	result := make([]byte, length)
	var bitPos int
	for i := range length {
		var val byte

		for j := 6; j >= 0; j-- {
			byteIdx := 1 + (bitPos / 8)
			bitIdx := 7 - (bitPos % 8)

			bit := (p.name[byteIdx] >> bitIdx) & 1
			if bit == 1 {
				val |= (1 << j)
			}
			bitPos++
		}
		result[i] = bit7ToChar(val)
	}

	return string(result)
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int((p.stats >> shiftMana) & maskMana)
}

func (p *GamePerson) Health() int {
	return int((p.stats >> shiftHealth) & maskHealth)
}

func (p *GamePerson) Respect() int {
	return int((p.stats >> shiftRespect) & maskRespect)
}

func (p *GamePerson) Strength() int {
	return int((p.stats >> shiftStr) & maskStr)
}

func (p *GamePerson) Experience() int {
	return int((p.stats >> shiftExp) & maskExp)
}

func (p *GamePerson) Level() int {
	return int((p.stats >> shiftLevel) & maskLevel)
}

func (p *GamePerson) HasHouse() bool {
	return (p.stats & (1 << shiftHouse)) != 0
}

func (p *GamePerson) HasGun() bool {
	return (p.stats & (1 << shiftGun)) != 0
}

func (p *GamePerson) HasFamily() bool {
	return (p.stats & (1 << shiftFamily)) != 0
}

func (p *GamePerson) Type() int {
	return int((p.stats >> shiftType) & maskType)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())

}
