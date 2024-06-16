package cli

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aliyasirnac/financialManagement/internal/model"
	"github.com/aliyasirnac/financialManagement/internal/services"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type cliItem struct {
	product model.Product
}

func (i cliItem) Title() string       { return i.product.Name }
func (i cliItem) Description() string { return i.product.Description }
func (i cliItem) FilterValue() string { return i.product.Name }

type state int

const (
	listing state = iota
	adding
	updating
	details
	showHelp
	alert
)

type cliModel struct {
	list           list.Model
	inputs         []textinput.Model
	focusIndex     int
	productService *services.ProductService
	state          state
	selectedItem   cliItem
	alertMessage   string
}

func initialModel(productService *services.ProductService) cliModel {
	products, err := productService.GetProducts()
	if err != nil {
		log.Fatal(err)
	}

	var items []list.Item
	for _, p := range products {
		items = append(items, cliItem{product: p})
	}

	inputs := make([]textinput.Model, 3) // Count field removed from the inputs array
	for i := range inputs {
		inputs[i] = textinput.New()
		inputs[i].CharLimit = 32
		inputs[i].Width = 20
	}

	inputs[0].Placeholder = "Name"
	inputs[1].Placeholder = "Description"
	inputs[2].Placeholder = "Price"

	m := cliModel{
		list:           list.New(items, list.NewDefaultDelegate(), 0, 0),
		inputs:         inputs,
		productService: productService,
		state:          listing,
	}
	m.list.Title = "Products"
	return m
}

func InitBubbleTea() {
	productService, err := services.NewProductService()
	if err != nil {
		log.Fatal(err)
	}

	m := initialModel(productService)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m cliModel) Init() tea.Cmd {
	return nil
}

func (m cliModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "?":
			m.state = showHelp
			return m, nil
		case "a":
			if m.state == listing {
				m.state = adding
			}
			return m, nil
		case "u":
			if m.state == listing {
				if len(m.list.Items()) == 0 {
					m.alertMessage = "Cannot update. No products available."
					m.state = alert
					return m, nil
				}
				selectedItem := m.list.SelectedItem().(cliItem)
				m.selectedItem = selectedItem
				m.inputs[0].SetValue(selectedItem.product.Name)
				m.inputs[1].SetValue(selectedItem.product.Description)
				m.inputs[2].SetValue(fmt.Sprintf("%d", selectedItem.product.Price))
				m.state = updating
			}
			return m, nil
		case "enter":
			if m.state == listing {
				if len(m.list.Items()) == 0 {
					m.alertMessage = "Cannot show details. No products available."
					m.state = alert
					return m, nil
				}
				selectedItem := m.list.SelectedItem().(cliItem)
				m.selectedItem = selectedItem
				m.state = details
				return m, nil
			}
			if m.state == adding || m.state == updating {
				name := m.inputs[0].Value()
				description := m.inputs[1].Value()
				price, err := strconv.ParseUint(m.inputs[2].Value(), 10, 16)
				if err != nil {
					log.Println("Invalid price")
					return m, nil
				}
				if m.state == adding {
					newProduct := model.Product{Name: name, Description: description, Price: uint16(price)}
					err = m.productService.CreateProduct(&newProduct)
					if err != nil {
						log.Println("Failed to create product:", err)
						return m, nil
					}
					m.list.InsertItem(len(m.list.Items()), cliItem{product: newProduct})
				} else if m.state == updating {
					m.selectedItem.product.Name = name
					m.selectedItem.product.Description = description
					m.selectedItem.product.Price = uint16(price)
					err = m.productService.UpdateProduct(&m.selectedItem.product)
					if err != nil {
						log.Println("Failed to update product:", err)
						return m, nil
					}
					m.list.SetItem(m.list.Index(), m.selectedItem)
				}
				m.state = listing
				for i := range m.inputs {
					m.inputs[i].SetValue("")
				}
				return m, nil
			}
		case "tab":
			if m.state == adding || m.state == updating {
				m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
				cmds := make([]tea.Cmd, len(m.inputs))
				for i := 0; i < len(m.inputs); i++ {
					if i == m.focusIndex {
						cmds[i] = m.inputs[i].Focus()
						m.inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
						m.inputs[i].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
					} else {
						m.inputs[i].Blur()
						m.inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
						m.inputs[i].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
					}
				}
				return m, tea.Batch(cmds...)
			}
		case "i":
			if m.state == listing {
				if len(m.list.Items()) > 0 {
					selectedItem := m.list.SelectedItem().(cliItem)
					selectedItem.product.Count++
					err := m.productService.UpdateProduct(&selectedItem.product)
					if err != nil {
						log.Println("Failed to increment product count:", err)
					} else {
						m.list.SetItem(m.list.Index(), selectedItem)
					}
				} else {
					m.alertMessage = "Cannot increment. No products available."
					m.state = alert
				}
			}
			return m, nil
		case "d":
			if m.state == listing {
				if len(m.list.Items()) > 0 {
					selectedItem := m.list.SelectedItem().(cliItem)
					err := m.productService.DeleteProduct(selectedItem.product.ID)
					if err != nil {
						log.Println("Failed to delete product:", err)
					} else {
						m.list.RemoveItem(m.list.Index())
					}
				} else {
					m.alertMessage = "Cannot delete. No products available."
					m.state = alert
				}
			}
			return m, nil
		case "esc":
			if m.state == adding || m.state == updating {
				m.state = listing
				for i := range m.inputs {
					m.inputs[i].SetValue("")
				}
				return m, nil
			} else if m.state == details || m.state == showHelp || m.state == alert {
				m.state = listing
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	if m.state == adding || m.state == updating {
		cmd = m.updateInputs(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m cliModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m cliModel) View() string {
	switch m.state {
	case adding:
		return m.addView()
	case updating:
		return m.updateView()
	case details:
		return m.detailsView()
	case showHelp:
		return m.helpView()
	case alert:
		return m.alertView()
	default:
		return m.listView()
	}
}
