package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/exp/rand"
)

type model struct {
	Helth     int
	Trees     int
	Gold      int
	KillKrips int
	damage    int
	armor     int
	turn      int
	gameOver  bool
	message   string
}

const (
	maxHealth     = 100
	tangoHealth   = 50
	creepsGoldMin = 152
	creepsGoldMax = 192
	creepsHealth  = 40
	pudgeGold     = 3000
	pudgeHealth   = 50
)

func initialModel() model {
	return model{
		Helth:     maxHealth,
		Trees:     0,
		Gold:      0,
		KillKrips: 0,
		armor:     5,
		damage:    5,
		turn:      0,
		gameOver:  false,
		message:   "Добро пожаловать в игру!",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) eatTango() {
	if m.Helth < maxHealth {
		if m.Trees >= 10 && m.armor >= 1 {
			m.message = "Похоже что вы съели слишком много деревьев, из-за этого ваше тело настолько увеличилось что оно порвало всю вашу броню"
			m.armor = 0
		}
		m.Helth += tangoHealth
		if m.Helth > maxHealth {
			m.Helth = maxHealth
		}
		m.message = "Вы нашли танго и использовали его, съев дерево"
		m.Trees += 1
	} else {
		m.message = "Вы нашли танго, но у вас был полный запас здоровья"
	}
}

func (m *model) attackCreeps() {
	if m.Helth >= 60 {
		var NewGold = rand.Intn(creepsGoldMax-creepsGoldMin+1) + creepsGoldMin
		var NewHelth = creepsHealth
		m.message = fmt.Sprintf("Вы смогли их убить, хотя и они вас немного ранили. Золота заработано %d, Очков здоровья потерянно %d", NewGold, NewHelth)
		m.Gold += NewGold
		m.Helth -= NewHelth
		m.KillKrips += 4
	} else {
		m.message = fmt.Sprintf("У вас не хватило здоровья и вы погибли. Заработано за игру %d, Сьедено деревьев за игру %d Убито крипов %d", m.Gold, m.Trees, m.KillKrips)
		m.Helth = 0
		m.gameOver = true
	}
}

func (m *model) fightPudge() {
	if m.Helth >= 60 && m.Gold >= 1000 && m.armor >= 5 {
		var NewGold = pudgeGold
		var StealHelth = pudgeHealth
		m.message = fmt.Sprintf("Вы смогли их убить, но это было слишком тяжело. Золота заработано %d, Очков здоровья потерянно %d", NewGold, StealHelth)
		m.Gold += NewGold
		m.Helth -= StealHelth
	} else {
		m.message = fmt.Sprintf("У вас не хватило здоровья и вы погибли. Заработано за игру %d, Сьедено деревьев за игру %d Убито крипов %d", m.Gold, m.Trees, m.KillKrips)
		m.Helth = 0
		m.gameOver = true
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.gameOver {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "w":
			Clear()
			m.turn += 1
			m.eatTango()
		case "a":
			Clear()
			m.turn += 1
			m.attackCreeps()
		case "d":
			Clear()
			m.turn += 1
			m.fightPudge()
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.gameOver {
		return fmt.Sprintf(
			"Игра окончена!\n"+
				"Заработано золота: %d\n"+
				"Съедено деревьев: %d\n"+
				"Убито крипов: %d\n",
			m.Gold, m.Trees, m.KillKrips,
		)
	}
	return fmt.Sprintf(
		"Ход: %d\n"+
			"Колличество жизней: %d\n"+
			"Колличество золота: %d\n"+
			"Колличество брони: %d\n"+
			"Сообщение: %s\n"+
			"Нажмите W что бы идти вперед, A для поворота налево, D для поворота направо и q для выхода.\n",
		m.turn, m.Helth, m.Gold, m.armor, m.message,
	)
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		panic(err)
	}
}

func Clear() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
