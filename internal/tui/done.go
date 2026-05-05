package tui

import (
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/tuna-os/bootc-installer-tui/internal/config"
)

type doneModel struct {
	cfg    *config.InstallConfig
	form   *huh.Form
	reboot bool
}

func newDoneModel(cfg *config.InstallConfig) *doneModel {
	m := &doneModel{cfg: cfg}
	m.reboot = true

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Reboot now?").
				Description("Remove the installation media before rebooting.").
				Affirmative("Reboot").
				Negative("Stay in installer").
				Value(&m.reboot),
		),
	)
	return m
}

func (m *doneModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *doneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}
	if m.form.State == huh.StateCompleted {
		if m.reboot {
			exec.Command("reboot").Run()
		}
		return m, tea.Quit
	}
	return m, cmd
}

func (m *doneModel) View() string {
	var sb strings.Builder
	sb.WriteString("\n\n")
	sb.WriteString(successStyle.Render("  ✓ Installation Complete!"))
	sb.WriteString("\n\n")
	sb.WriteString(accentStyle.Render("  The system has been installed successfully."))
	sb.WriteString("\n")
	sb.WriteString(mutedStyle.Render("  You can now reboot into your new system."))
	sb.WriteString("\n\n")
	sb.WriteString(m.form.View())
	return sb.String()
}
