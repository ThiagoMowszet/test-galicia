package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/atotto/clipboard" // Librería para copiar texto
)

// Modelo del programa
type model struct {
	cursor     int  // Para seleccionar entre OCOR y CRAFT
	subCursor  int  // Para navegar entre las subopciones
	selected   bool // Para verificar si ya se seleccionó una opción
	option     string // Guardará "OCOR" o "CRAFT" cuando se seleccionen
	subOptions bool // Verifica si se está dentro de las subopciones (OCOR 1, CRAFT 1, etc.)
	showLorem  bool // Si es true, mostramos "lorem ipsum"
	copyStatus string // Estado del botón de copiar
}

func (m model) Init() tea.Cmd {
	// Manejar señal de Ctrl+C para salir
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Saliendo del programa...")
		os.Exit(0)
	}()
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Si hay una tecla presionada
	case tea.KeyMsg:
		switch msg.String() {
		// Seleccionar una opción del menú
		case "enter":
			if m.showLorem {
				// Simular copiar texto al portapapeles
				clipboard.WriteAll("plantilla personalizada") // Usa atotto/clipboard para copiar texto
				m.copyStatus = "Texto copiado al portapapeles"
			} else if m.subOptions {
				// Si estamos en las subopciones, mostramos "lorem ipsum"
				m.showLorem = true
			} else if m.selected {
				// Si ya se seleccionó una opción, vamos a las subopciones
				m.subOptions = true
				m.subCursor = 0 // Reiniciamos el cursor de las subopciones
			} else {
				// Seleccionamos "OCOR" o "CRAFT"
				if m.cursor == 0 {
					m.option = "OCOR"
				} else {
					m.option = "CRAFT"
				}
				m.selected = true
			}

		// Volver al menú principal si se presiona "Esc"
		case "esc":
			if m.showLorem {
				m.showLorem = false
				m.copyStatus = ""
			} else if m.subOptions {
				m.subOptions = false
			} else {
				m.selected = false
				m.option = ""
			}

		// Navegar entre opciones
		case "up", "k":
			if m.subOptions && !m.showLorem && m.subCursor > 0 {
				m.subCursor--
			} else if !m.subOptions && !m.selected && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.subOptions && !m.showLorem && m.subCursor < 4 {
				m.subCursor++
			} else if !m.subOptions && !m.selected && m.cursor < 1 {
				m.cursor++
			}

		// Salir si se presiona Ctrl+C
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	// Bordes para simular el cuadro
	borderTop := "╭─────────────────────────╮"
	borderBottom := "╰─────────────────────────╯"
	spacing := "│                         │"

	// Si estamos mostrando "lorem ipsum"
	if m.showLorem {
		content := "plantilla personalizada"
		copyText := fmt.Sprintf("[ %s ]", m.copyStatus) // Muestra el estado de copiado

		return fmt.Sprintf("%s\n%s\n│   %s   │\n│   %s   │\n%s\n%s\n", borderTop, spacing, centerText(content, 25), centerText(copyText, 25), spacing, borderBottom)
	}

	// Si ya se seleccionó una opción y estamos dentro de las subopciones
	if m.subOptions {
		options := []string{"1", "2", "3", "4", "5"}
		content := ""

		for i, option := range options {
			prefix := "  " // Espacio si no está seleccionado
			color := "\033[32m" // Verde para las subopciones
			reset := "\033[0m"
			if m.subCursor == i {
				prefix = "> " // Flecha si está seleccionado
				color = "\033[33m" // Amarillo si está seleccionado
			}
			content += fmt.Sprintf("%s%s%s %s%s\n", prefix, color, m.option, option, reset)
		}

		return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", borderTop, spacing, centerText(content, 25), spacing, borderBottom)
	}

	// Menú principal
	s := ""
	if m.cursor == 0 {
		s += "> \033[34mOCOR\033[0m\n" // Azul cuando está seleccionado
	} else {
		s += "  OCOR\n"
	}

	if m.cursor == 1 {
		s += "> \033[31mCRAFT\033[0m\n" // Rojo cuando está seleccionado
	} else {
		s += "  CRAFT\n"
	}

	content := s
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", borderTop, spacing, centerText(content, 25), spacing, borderBottom)
}

// Función para centrar el texto dentro del cuadro
func centerText(text string, width int) string {
	lines := strings.Split(text, "\n")
	var centeredLines []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		padding := (width - len(trimmedLine)) / 2
		centeredLines = append(centeredLines, strings.Repeat(" ", padding)+trimmedLine)
	}
	return strings.Join(centeredLines, "\n")
}

func main() {
	// Retraso para que la transición sea más fluida
	time.Sleep(100 * time.Millisecond)

	p := tea.NewProgram(model{})
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
