package template

import "html/template"

// FuncMap содержит все вспомогательные функции для шаблонов
var FuncMap = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"subtract": func(a, b int) int {
		return a - b
	},
	"sequence": func(start, end int) []int {
		var sequence []int
		for i := start; i <= end; i++ {
			sequence = append(sequence, i)
		}
		return sequence
	},
}
