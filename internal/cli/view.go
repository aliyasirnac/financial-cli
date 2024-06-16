package cli

import "fmt"

func (m cliModel) listView() string {
	return docStyle.Render(m.list.View())
}

func (m cliModel) addView() string {
	return fmt.Sprintf(
		"Add a New Product\n\n%s\n%s\n%s\n\nPress Enter to Add, Tab to Navigate, Esc to Cancel\n",
		m.inputs[0].View(),
		m.inputs[1].View(),
		m.inputs[2].View(),
	)
}

func (m cliModel) updateView() string {
	return fmt.Sprintf(
		"Update Product\n\n%s\n%s\n%s\n\nPress Enter to Update, Tab to Navigate, Esc to Cancel\n",
		m.inputs[0].View(),
		m.inputs[1].View(),
		m.inputs[2].View(),
	)
}

func (m cliModel) detailsView() string {
	totalPrice := m.selectedItem.product.Price * uint16(m.selectedItem.product.Count)
	return fmt.Sprintf(
		"Ürün Detayları\n\nİsmi: %s\nFiyatı: %d\nTotal Harcama: %d TL\n\nESC tuşuna basarak geri dön.\n",
		m.selectedItem.product.Name,
		m.selectedItem.product.Price,
		totalPrice,
	)
}

func (m cliModel) helpView() string {
	return "Help Menu\n\n" +
		"a - Yeni ürün ekle\n" +
		"u - Seçili ürünü güncelle\n" +
		"i - Satın almayı arttır\n" +
		"d - Seçili ürünü sil\n" +
		"enter - View product details\n" +
		"tab - Navigate inputs\n" +
		"esc- Cancel or go back\n\n" +
		"Press Esc to Go Back"
}

func (m cliModel) alertView() string {
	return fmt.Sprintf(
		"Alert\n\n%s\n\nPress Esc to Go Back",
		m.alertMessage,
	)
}
