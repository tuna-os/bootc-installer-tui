package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/tuna-os/bootc-installer-tui/internal/config"
)

type confirmModel struct {
	cfg     *config.InstallConfig
	form    *huh.Form
	proceed bool
}

func newConfirmModel(cfg *config.InstallConfig) *confirmModel {
	m := &confirmModel{cfg: cfg}
	m.proceed = false

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Start installation?").
				Description("Review the summary above. This will ERASE the target disk.").
				Affirmative("Install now").
				Negative("Go back").
				Value(&m.proceed),
		),
	)
	return m
}

func (m *confirmModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		if m.proceed {
			return m, func() tea.Msg { return stepMsg{next: stepProgress} }
		}
		return m, func() tea.Msg { return stepMsg{next: stepUser} }
	}
	return m, cmd
}

func (m *confirmModel) View() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(titleStyle.Render("  Step 8 of 9: Confirm Installation"))
	sb.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(colorMuted).Width(18)
	valueStyle := lipgloss.NewStyle().Foreground(colorText).Bold(true)

	row := func(label, value string) string {
		return "  " + labelStyle.Render(label+":") + " " + valueStyle.Render(value) + "\n"
	}

	encStr := "Disabled"
	if m.cfg.EncryptionEnabled {
		encStr = "LUKS (passphrase)"
	}

	keysStr := "None"
	if len(m.cfg.SSHKeys) > 0 {
		keysStr = fmt.Sprintf("%d key(s) imported", len(m.cfg.SSHKeys))
	}

	sb.WriteString(boxStyle.Render(
		titleStyle.Render("Installation Summary") + "\n\n" +
			row("Image", m.cfg.Image) +
			row("Disk", m.cfg.DiskDevice) +
			row("Filesystem", m.cfg.Filesystem) +
			row("Hostname", m.cfg.Hostname) +
			row("Encryption", encStr) +
			row("Username", m.cfg.Username) +
			row("Full name", m.cfg.FullName) +
			row("SSH keys", keysStr),
	))

	sb.WriteString("\n\n")
	sb.WriteString(warningStyle.Render("  ⚠  ALL DATA ON " + m.cfg.DiskDevice + " WILL BE DESTROYED"))
	sb.WriteString("\n\n")
	sb.WriteString(m.form.View())
	return sb.String()
}
