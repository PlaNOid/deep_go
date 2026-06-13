package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

const (
	NameMaxLength  = 42
	BitsPerChar    = 7
	BitsPerByte    = 8
	NameDataOffset = 1

	// calculate array lengh at compilation
	NameArraySize = NameDataOffset + (NameMaxLength*BitsPerChar+BitsPerByte-1)/BitsPerByte

	MinPrintableASCII = 32  // space
	MaxPrintableASCII = 126 // ~
)

type Option func(*GamePerson)

func charTo7Bit(c byte) byte {
	if c >= MinPrintableASCII && c <= MaxPrintableASCII {
		return c - MinPrintableASCII
	}
	return 0
}

func bit7ToChar(v byte) byte {
	return v + MinPrintableASCII
}

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		person.name = [NameArraySize]byte{}
		person.name[0] = byte(len(name))

		var bitPos int

		for _, ch := range name {
			val := charTo7Bit(byte(ch))

			// write symbol's bit, from head to tail
			for i := BitsPerChar - 1; i >= 0; i-- {
				bit := (val >> i) & 1
				byteIdx := NameDataOffset + (bitPos / BitsPerByte)
				bitIdx := (BitsPerByte - 1) - (bitPos % BitsPerByte)

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
	x, y, z int32               // 4B + 4B + 4B = 12B
	gold    uint32              // 4B
	stats   uint64              // 10b (mana) + 10b (health) + 4b (respect) + 4b (str) + 4b (exp) + 4b (lvl) + 3b (house, family, gun) + 2b (type) = 41b ~ 8B
	name    [NameArraySize]byte // encoded 95 ASCII symbols (26(A-Z), 26(a-z), 10(0-9), 33(! _ " etc)) 7b for one (2^7 = 128),  7b * 42(max lenght) = 294b, (294b / 8 ~ 37B) + 1B(lenght) = 38B
	// 12B + 4B + 8B + 38B = 62B align to 64B
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

		// read symbol's bit, from head to tail
		for j := BitsPerChar - 1; j >= 0; j-- {
			byteIdx := NameDataOffset + (bitPos / BitsPerByte)
			bitIdx := (BitsPerByte - 1) - (bitPos % BitsPerByte)

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
