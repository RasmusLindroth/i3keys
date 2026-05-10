package tui

import (
	"fmt"
	"log"
	"os"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/RasmusLindroth/i3keys/internal/i3parse"
	"github.com/RasmusLindroth/i3keys/internal/keyboard"
	"github.com/RasmusLindroth/i3keys/internal/xlib"
)

type keyMap struct {
	Enter     key.Binding
	Left      key.Binding
	Right     key.Binding
	Back      key.Binding
	Quit      key.Binding
	ForceQuit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right, k.Back, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right},
		{k.Back, k.Quit},
	}
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title string
	desc  string
}

func (i item) Title() string {
	return i.title
}
func (i item) Description() string {
	return i.desc
}
func (i item) FilterValue() string {
	return i.title
}

type Focus uint

const (
	ModeListFocus Focus = iota
	GroupListFocus
	KeyboardFocus
)

type model struct {
	theme     Theme
	layout    string
	modifiers map[string][]string

	currentFocus Focus

	list list.Model
	vp   viewport.Model
	help help.Model
	keys keyMap

	modeListIndex  int
	groupListIndex int

	modeItems []list.Item
	groups    []i3parse.ModifierGroup
	modes     []i3parse.Mode

	tmpKeyboard keyboard.Keyboard
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.ForceQuit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Quit):
			if m.currentFocus == KeyboardFocus {
				return m, tea.Quit
			}
		case key.Matches(msg, m.keys.Back):
			if m.currentFocus == GroupListFocus && m.list.FilterState() != list.FilterApplied {
				m.currentFocus = ModeListFocus
				m.list.Title = "Modes"
				m.list.SetItems(m.modeItems)
				m.list.Select(m.modeListIndex)
				return m, nil
			}
			if m.currentFocus == KeyboardFocus {
				m.currentFocus = GroupListFocus
				return m, nil
			}

		case key.Matches(msg, m.keys.Left):
			if m.currentFocus == KeyboardFocus {
				m.vp.ScrollLeft(5)
			}
		case key.Matches(msg, m.keys.Right):
			if m.currentFocus == KeyboardFocus {
				m.vp.ScrollRight(5)
			}
		case key.Matches(msg, m.keys.Enter):
			if m.currentFocus == ModeListFocus && m.list.FilterState() != list.Filtering {
				items := []list.Item{}
				m.groupListIndex = 0
				m.currentFocus = GroupListFocus
				gi := m.list.GlobalIndex()
				m.modeListIndex = gi
				if gi >= len(m.modes) && len(m.modes) > 0 {
					return m, nil
				}
				mode := m.modes[gi]
				groups := i3parse.GetModifierGroups(mode.Bindings)
				m.groups = groups
				for _, group := range groups {
					mods := strings.Join(group.Modifiers, "+")
					if mods == "" {
						mods = "default"
					}
					desc := fmt.Sprintf("bindings: %d", len(group.Bindings))
					items = append(items, item{
						title: mods,
						desc:  desc,
					})
				}
				m.list.Title = "Groups"
				m.list.SetItems(items)
				m.list.GoToStart()
				return m, nil
			}
			if m.currentFocus == GroupListFocus && m.list.FilterState() != list.Filtering {
				index := m.list.GlobalIndex()
				if index < len(m.groups) {
					kb, err := keyboard.MapKeyboard(m.layout, m.groups[index], m.modifiers)
					if err != nil {
						return m, nil
					}
					s := m.drawKeyboard(kb)
					m.vp.SetContent(s)
					m.vp.SetXOffset(0)
					m.currentFocus = KeyboardFocus
					return m, nil
				}

			}
		}
	case tea.WindowSizeMsg:
		horizontal, vertical := docStyle.GetFrameSize()

		m.list.SetSize(
			msg.Width-horizontal,
			msg.Height-vertical,
		)

		m.vp.SetWidth(msg.Width - horizontal)

		availableHeight := msg.Height - vertical
		helpHeight := lipgloss.Height(m.help.View(m.keys))

		m.vp.SetHeight(availableHeight - helpHeight)

	}

	var cmd tea.Cmd
	if m.currentFocus == ModeListFocus || m.currentFocus == GroupListFocus {
		m.list, cmd = m.list.Update(msg)
	}
	if m.currentFocus == KeyboardFocus {
		var cmds []tea.Cmd
		m.vp, cmd = m.vp.Update(msg)
		cmds = append(cmds, cmd)
		m.help, cmd = m.help.Update(msg)
		cmds = append(cmds, cmd)
		cmd = tea.Batch(cmds...)
	}

	return m, cmd
}

func (m model) View() tea.View {
	content := ""
	switch m.currentFocus {
	case ModeListFocus, GroupListFocus:
		content = m.list.View()
	case KeyboardFocus:
		helpView := m.help.View(m.keys)

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			m.vp.View(),
			helpView,
		)
	}
	v := tea.NewView(docStyle.Render(content))
	v.AltScreen = true
	return v
}

func (m model) drawKeyboard(kb keyboard.Keyboard) string {
	rows := []string{}
	inUse := lipgloss.NewStyle().Foreground(m.theme.KeyInUse)

	for _, keyRow := range kb.Keys {
		r := []string{}
		for _, key := range keyRow {
			k := key.Symbol
			if key.InUse {
				k = inUse.Render(k)
			}
			r = append(r, k)
		}
		rows = append(rows, strings.Join(r, ", "))
	}
	return lipgloss.JoinVertical(lipgloss.Top, rows...)
}

func Start(wm string, layout string) {
	modes, _, err := i3parse.ParseFromRunning(wm, true)

	if err != nil {
		log.Fatalln(err)
	}

	layout = strings.ToUpper(layout)

	modifiers := xlib.GetModifiers()

	items := []list.Item{}
	for _, mode := range modes {
		groups := i3parse.GetModifierGroups(mode.Bindings)

		modeName := mode.Name
		if modeName == "" {
			modeName = "default"
		}
		items = append(items, item{
			title: modeName,
			desc:  fmt.Sprintf("groups: %d", len(groups)),
		})
	}

	theme := ActiveTheme()
	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(theme.NormalTitle)
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.Foreground(theme.NormalDesc)
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(theme.SelectedTitle).BorderForeground(theme.SelectedBorder)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(theme.SelectedDesc).BorderForeground(theme.SelectedBorder)
	delegate.Styles.DimmedTitle = delegate.Styles.DimmedTitle.Foreground(theme.DimmedTitle)
	delegate.Styles.DimmedDesc = delegate.Styles.DimmedDesc.Foreground(theme.DimmedDesc)

	keys := keyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter", "select"),
			key.WithHelp("enter", "select"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		ForceQuit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "force quit"),
		),
	}
	m := model{
		keys:         keys,
		theme:        theme,
		list:         list.New(items, delegate, 0, 0),
		vp:           viewport.New(),
		help:         help.New(),
		modeItems:    items,
		layout:       layout,
		modifiers:    modifiers,
		modes:        modes,
		currentFocus: ModeListFocus,
	}
	m.vp.SetHorizontalStep(5)

	m.list.Title = "Modes"
	m.list.SetShowStatusBar(false)
	m.list.SetFilteringEnabled(true)
	m.list.KeyMap.Quit = m.keys.Quit

	m.list.Styles.Title = m.list.Styles.Title.Background(theme.TitleBG).Foreground(theme.TitleText)

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
