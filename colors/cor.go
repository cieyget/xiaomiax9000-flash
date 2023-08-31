package main

import (
	"fmt"
	"math"

	"github.com/fatih/color"
)

func gradientText(text string, colors []color.Attribute, lineIndex int) string {
	length := len(text)
	numColors := len(colors)
	gradientText := ""

	for i, char := range text {
		colorIndex := int(math.Mod(float64((i+lineIndex)*numColors)/float64(length), float64(numColors)))
		colorCode := colors[colorIndex]
		coloredChar := color.New(colorCode).SprintFunc()(string(char))
		gradientText += coloredChar
	}

	return gradientText
}

func main() {
	lines := []string{
		`     _    __  _____   ___   ___   ___    _____ _        _    ____  _   _ `,
		`    / \   \ \/ / _ \ / _ \ / _ \ / _ \  |  ___| |      / \  / ___|| | | |`,
		`   / _ \   \  / (_) | | | | | | | | | | | |_  | |     / _ \ \___ \| |_| |`,
		`  / ___ \  /  \\__, | |_| | |_| | |_| | |  _| | |___ / ___ \ ___) |  _  |`,
		` /_/   \_\/_/\_\ /_/ \___/ \___/ \___/  |_|   |_____/_/   \_\____/|_| |_|`,
		`               OpenWrt XiaoMi AX9000 Flash Tool by CieyGet               `,
		`                                                                         `,
	}

	// 类似于 oh-my-zsh 的颜色数组（蓝色、洋红、青色、绿色、黄色）
	colors := []color.Attribute{color.FgBlue, color.FgMagenta, color.FgCyan, color.FgGreen, color.FgYellow}

	for index, line := range lines {
		fmt.Println(gradientText(line, colors, index))
	}
}
