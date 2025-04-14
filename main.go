package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/promwish"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/accesscontrol"
	"github.com/charmbracelet/wish/activeterm"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
)

//go:embed frames/*.txt
var fsys embed.FS

var (
	port        = flag.Int("port", 2222, "port to listen on")
	metricsPort = flag.Int("metrics-port", 9222, "port to listen on")
)

func main() {
	flag.Parse()
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("0.0.0.0:%d", *port)),
		wish.WithHostKeyPath(".ssh/parttysh"),
		wish.WithMiddleware(
			bm.Middleware(teaHandler()),
			lm.Middleware(),
			promwish.Middleware(fmt.Sprintf("0.0.0.0:%d", *metricsPort), "parttysh"),
			accesscontrol.Middleware(),
			activeterm.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Starting SSH server on 0.0.0.0:%d", *port)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err = s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()
	<-done
	if err := s.Close(); err != nil {
		log.Fatalln(err)
	}
}

func teaHandler() func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	var frames []string
	if err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		bts, err := fsys.ReadFile(path)
		frames = append(frames, string(bts))
		return err
	}); err != nil {
		log.Fatalln(err)
	}
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		return newModel(frames), []tea.ProgramOption{tea.WithAltScreen()}
	}
}

func newModel(frames []string) model {
	spin := spinner.New()
	spin.Spinner = spinner.Spinner{
		Frames: frames,
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
	return m.spin.Tick
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
