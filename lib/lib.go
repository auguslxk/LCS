package lib

import (
	"os"
	"unicode"
)

type matrix struct {
	data [][]int
}

func initMatrix(width, height int, m *matrix) {
	m.data = make([][]int, height+1)
	for i := 0; i <= height; i++ {
		m.data[i] = make([]int, width+1)
	}

	for i := 0; i <= height; i++ {
		m.data[i][0] = 0
	}

	for i := 0; i <= width; i++ {
		m.data[0][i] = 0
	}
}

func GetLCS(strA, strB []rune) ([]rune, []int, []int) {
	stra := processRune(strA)
	strb := processRune(strB)

	if len(stra) > len(strb) {
		stra, strb = strb, stra
	}

	height := len(stra)
	width := len(strb)

	m := &matrix{}
	initMatrix(width, height, m)

	for i := 1; i <= height; i++ {
		for j := 1; j <= width; j++ {
			if stra[i-1] == strb[j-1] {
				m.data[i][j] = m.data[i-1][j-1] + 1
			} else {
				m.data[i][j] = Max(m.data[i-1][j], m.data[i][j-1]).(int)
			}
		}
	}

	indexA := make([]int, m.data[height][width])
	indexB := make([]int, m.data[height][width])
	lcsIndex := m.data[height][width] - 1
	var lcs = []rune{}
	for height > 0 && width > 0 {
		if m.data[height][width] == m.data[height-1][width] {
			height--
		} else if m.data[height][width] == m.data[height][width-1] {
			width--
		} else {
			height--
			width--
			indexA[lcsIndex] = width
			indexB[lcsIndex] = height
			lcsIndex--
			lcs = append(stra[height:height+1], lcs...)
		}
	}
	return lcs, indexA, indexB
}

func Max(a, b interface{}) interface{} {
	switch a.(type) {
	case int:
		if a.(int) > b.(int) {
			return a
		}
		return b
	case float32:
		if a.(float32) > b.(float32) {
			return a
		}
		return b
	case float64:
		if a.(float64) > b.(float64) {
			return a
		}
		return b
	}
	return nil
}

func processRune(word []rune) []rune {
	product := []rune{}
	for _, char := range word {
		if unicode.IsUpper(char) {
			product = append(product, unicode.ToLower(char))
		} else {
			product = append(product, char)
		}
	}
	return product
}

func ReadFile(filename string) (string, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return "", err
	}

	buf := make([]byte, 1024)
	content := []byte{}
	for {
		size, err := file.Read(buf)
		if err != nil {
			break
		}
		if size == 0 {
			break
		}
		content = append(content, buf[0:size]...)
	}
	return string(content), nil
}
