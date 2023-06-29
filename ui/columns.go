package ui

type Columns struct {
	sumWidth int
	columns  []Column
}
type Column struct {
	Name  string
	Width int
}

func (c Columns) GetSize(name string, maxX int) int {
	for _, col := range c.columns {
		if col.Name == name {
			if col.Width != -1 {
				return col.Width
			}

			return maxX - c.sumWidth
		}
	}

	return -2
}

func (c Columns) Add(name string, width int) Columns {
	c.sumWidth = c.sumWidth + width
	c.columns = append(c.columns, Column{
		Name:  name,
		Width: width,
	})
	return c
}
