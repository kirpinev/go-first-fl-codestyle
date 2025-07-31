// Package main реализует текстовую RPG игру с системой классов персонажей.
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// CharacterClass представляет тип класса персонажа
type CharacterClass string

// Константы для классов персонажей
const (
	WarriorClass CharacterClass = "warrior"
	MageClass    CharacterClass = "mage"
	HealerClass  CharacterClass = "healer"
)

// Базовые характеристики
const (
	BaseAttack  = 5
	BaseDefense = 10
	BaseStamina = 80
)

// Stats представляет характеристики персонажа
type Stats struct {
	Attack  int
	Defense int
	Stamina int
}

// Character представляет игрового персонажа
type Character struct {
	Name  string
	Class CharacterClass
	Stats Stats
}

// Action представляет действие, которое может выполнить персонаж
type Action interface {
	Execute(character *Character) string
	GetName() string
}

// AttackAction реализует действие атаки
type AttackAction struct{}

func (a AttackAction) GetName() string {
	return "attack"
}

func (a AttackAction) Execute(character *Character) string {
	damage := character.calculateAttackDamage()
	return fmt.Sprintf("%s нанес урон противнику равный %d.", character.Name, damage)
}

// DefenseAction реализует действие защиты
type DefenseAction struct{}

func (d DefenseAction) GetName() string {
	return "defense"
}

func (d DefenseAction) Execute(character *Character) string {
	defense := character.calculateDefenseValue()
	return fmt.Sprintf("%s блокировал %d урона.", character.Name, defense)
}

// SpecialAction реализует специальное действие
type SpecialAction struct{}

func (s SpecialAction) GetName() string {
	return "special"
}

func (s SpecialAction) Execute(character *Character) string {
	return character.useSpecialAbility()
}

// Game представляет игровую сессию
type Game struct {
	reader  *bufio.Scanner
	actions map[string]Action
}

// NewGame создает новую игру
func NewGame() *Game {
	game := &Game{
		reader: bufio.NewScanner(os.Stdin),
		actions: make(map[string]Action),
	}
	
	// Регистрируем доступные действия
	game.registerAction(AttackAction{})
	game.registerAction(DefenseAction{})
	game.registerAction(SpecialAction{})
	
	return game
}

// registerAction регистрирует новое действие в игре
func (g *Game) registerAction(action Action) {
	g.actions[action.GetName()] = action
}

// readInput читает ввод пользователя с обработкой ошибок
func (g *Game) readInput(prompt string) (string, error) {
	fmt.Print(prompt)
	if !g.reader.Scan() {
		return "", fmt.Errorf("ошибка чтения ввода")
	}
	return strings.TrimSpace(g.reader.Text()), nil
}

// createCharacter создает нового персонажа
func (g *Game) createCharacter() (*Character, error) {
	name, err := g.readInput("...назови себя: ")
	if err != nil {
		return nil, err
	}
	
	if name == "" {
		return nil, fmt.Errorf("имя не может быть пустым")
	}

	fmt.Printf("Здравствуй, %s\n", name)
	fmt.Printf("Сейчас твоя выносливость — %d, атака — %d и защита — %d.\n", 
		BaseStamina, BaseAttack, BaseDefense)
	fmt.Println("Ты можешь выбрать один из трёх путей силы:")
	fmt.Println("Воитель, Маг, Лекарь")

	class, err := g.chooseCharacterClass()
	if err != nil {
		return nil, err
	}

	character := &Character{
		Name:  name,
		Class: class,
		Stats: Stats{
			Attack:  BaseAttack,
			Defense: BaseDefense,
			Stamina: BaseStamina,
		},
	}

	return character, nil
}

// chooseCharacterClass позволяет игроку выбрать класс персонажа
func (g *Game) chooseCharacterClass() (CharacterClass, error) {
	validClasses := map[string]CharacterClass{
		"warrior": WarriorClass,
		"mage":    MageClass,
		"healer":  HealerClass,
	}

	classDescriptions := map[CharacterClass]string{
		WarriorClass: "Воитель — дерзкий воин ближнего боя. Сильный, выносливый и отважный.",
		MageClass:    "Маг — находчивый воин дальнего боя. Обладает высоким интеллектом.",
		HealerClass:  "Лекарь — могущественный заклинатель. Черпает силы из природы, веры и духов.",
	}

	for {
		input, err := g.readInput("Введи название персонажа: Воитель — warrior, Маг — mage, Лекарь — healer: ")
		if err != nil {
			return "", err
		}

		class, exists := validClasses[strings.ToLower(input)]
		if !exists {
			fmt.Println("Неизвестный класс персонажа. Попробуйте еще раз.")
			continue
		}

		fmt.Println(classDescriptions[class])

		confirm, err := g.readInput("Нажми (Y), чтобы подтвердить выбор, или любую другую кнопку, чтобы выбрать другого персонажа: ")
		if err != nil {
			return "", err
		}

		if strings.ToLower(confirm) == "y" {
			return class, nil
		}
	}
}

// startTraining запускает тренировочный режим
func (g *Game) startTraining(character *Character) error {
	character.showClassDescription()
	g.showInstructions()

	for {
		input, err := g.readInput("Введи команду: ")
		if err != nil {
			return err
		}

		if input == "skip" {
			fmt.Println("тренировка окончена")
			return nil
		}

		action, exists := g.actions[input]
		if !exists {
			fmt.Println("Неизвестная команда. Попробуйте: attack, defense, special или skip")
			continue
		}

		result := action.Execute(character)
		fmt.Println(result)
	}
}

// showInstructions показывает инструкции игроку
func (g *Game) showInstructions() {
	fmt.Println("Потренируйся управлять своими навыками.")
	fmt.Println("Введи одну из команд:")
	fmt.Println("  attack — чтобы атаковать противника")
	fmt.Println("  defense — чтобы блокировать атаку противника")
	fmt.Println("  special — чтобы использовать свою суперсилу")
	fmt.Println("  skip — чтобы закончить тренировку")
}

// Run запускает игру
func (g *Game) Run() error {
	fmt.Println("Приветствую тебя, искатель приключений!")
	fmt.Println("Прежде чем начать игру...")

	character, err := g.createCharacter()
	if err != nil {
		return fmt.Errorf("ошибка создания персонажа: %w", err)
	}

	return g.startTraining(character)
}

// showClassDescription показывает описание класса персонажа
func (c *Character) showClassDescription() {
	descriptions := map[CharacterClass]string{
		WarriorClass: "%s, ты Воитель - отличный боец ближнего боя.",
		MageClass:    "%s, ты Маг - превосходный укротитель стихий.",
		HealerClass:  "%s, ты Лекарь - чародей, способный исцелять раны.",
	}

	if desc, exists := descriptions[c.Class]; exists {
		fmt.Printf(desc+"\n", c.Name)
	}
}

// calculateAttackDamage вычисляет урон атаки в зависимости от класса
func (c *Character) calculateAttackDamage() int {
	damageRanges := map[CharacterClass][2]int{
		WarriorClass: {3, 5},
		MageClass:    {5, 10},
		HealerClass:  {-3, -1},
	}

	if dmgRange, exists := damageRanges[c.Class]; exists {
		return c.Stats.Attack + randRange(dmgRange[0], dmgRange[1])
	}
	return c.Stats.Attack
}

// calculateDefenseValue вычисляет значение защиты в зависимости от класса
func (c *Character) calculateDefenseValue() int {
	defenseRanges := map[CharacterClass][2]int{
		WarriorClass: {5, 10},
		MageClass:    {-2, 2},
		HealerClass:  {2, 5},
	}

	if defRange, exists := defenseRanges[c.Class]; exists {
		return c.Stats.Defense + randRange(defRange[0], defRange[1])
	}
	return c.Stats.Defense
}

// useSpecialAbility использует специальную способность класса
func (c *Character) useSpecialAbility() string {
	abilities := map[CharacterClass]struct {
		name  string
		value int
	}{
		WarriorClass: {"Выносливость", c.Stats.Stamina + 25},
		MageClass:    {"Атака", c.Stats.Attack + 40},
		HealerClass:  {"Защита", c.Stats.Defense + 30},
	}

	if ability, exists := abilities[c.Class]; exists {
		return fmt.Sprintf("%s применил специальное умение `%s %d`", 
			c.Name, ability.name, ability.value)
	}
	return "неизвестный класс персонажа"
}

// randRange возвращает случайное число в заданном диапазоне (включительно)
func randRange(min, max int) int {
	if min > max {
		min, max = max, min
	}
	return rand.Intn(max-min+1) + min
}

// initRandom инициализирует генератор случайных чисел
func initRandom() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	initRandom()
	
	game := NewGame()
	if err := game.Run(); err != nil {
		fmt.Printf("Ошибка игры: %v\n", err)
		os.Exit(1)
	}
}
