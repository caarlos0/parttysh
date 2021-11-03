package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "embed"

	"github.com/caarlos0/env/v6"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
	"github.com/muesli/termenv"
)

//go:embed frames/0.txt
var f0 string

//go:embed frames/1.txt
var f1 string

//go:embed frames/2.txt
var f2 string

//go:embed frames/3.txt
var f3 string

//go:embed frames/4.txt
var f4 string

//go:embed frames/5.txt
var f5 string

//go:embed frames/6.txt
var f6 string

//go:embed frames/7.txt
var f7 string

//go:embed frames/8.txt
var f8 string

//go:embed frames/9.txt
var f9 string

type Config struct {
	Host string `env:"PARTTYSH_HOST" envDefault:"127.0.0.1"`
	Port int    `env:"PARTTYSH_PORT" envDefault:"2222"`
}

func main() {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln(err)
	}

	// force colors as we might start it from systemd which has no interactive term and no colors
	lipgloss.SetColorProfile(termenv.ANSI256)

	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)),
		wish.WithHostKeyPath(".ssh/parttysh"),
		wish.WithMiddleware(
			bm.Middleware(teaHandler(cfg)),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Starting SSH server on %s:%d", cfg.Host, cfg.Port)
	err = s.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

func teaHandler(cfg Config) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		_, _, active := s.Pty()
		if !active {
			fmt.Println("no active terminal, skipping")
			return nil, nil
		}

		return newModel(), []tea.ProgramOption{tea.WithAltScreen()}
	}
}

func newModel() model {
	spin := spinner.NewModel()
	spin.Spinner = spinner.Spinner{
		Frames: []string{f0, f1, f2, f3, f4, f5, f6, f7, f8, f9},
		FPS:    time.Second / 15,
	}
	var colors []lipgloss.Style
	for _, color := range []string{
		"#FF0000", // red
		"#FFFF00", // yellow
		"#00FF00", // green
		"#0247FE", // blue
		"#FF00FF", // magenta
		"#00FFFF", // cyan
		"#FFFFFF", // white
	} {
		colors = append(colors, lipgloss.NewStyle().Foreground(lipgloss.Color(color)))
	}
	return model{
		spin:   spin,
		colors: colors,
	}
}

type model struct {
	spin   spinner.Model
	colors []lipgloss.Style
}

// Init initializes the confetti after a small delay
func (m model) Init() tea.Cmd {
	return spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	m.spin.Style = m.colors[rand.Intn(len(m.colors))]
	s, cmd := m.spin.Update(msg)
	m.spin = s
	return m, cmd
}

// View displays all the particles on the screen
func (m model) View() string {
	return m.spin.View()
}
