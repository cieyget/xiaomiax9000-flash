package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println(`     _    __  _____   ___   ___   ___    _____ _        _    ____  _   _ `)
	fmt.Println(`    / \   \ \/ / _ \ / _ \ / _ \ / _ \  |  ___| |      / \  / ___|| | | |`)
	fmt.Println(`   / _ \   \  / (_) | | | | | | | | | | | |_  | |     / _ \ \___ \| |_| |`)
	fmt.Println(`  / ___ \  /  \\__, | |_| | |_| | |_| | |  _| | |___ / ___ \ ___) |  _  |`)
	fmt.Println(` /_/   \_\/_/\_\ /_/ \___/ \___/ \___/  |_|   |_____/_/   \_\____/|_| |_|`)
	fmt.Println(`               OpenWrt XiaoMi AX9000 Flash Tool by CieyGet               `)
	fmt.Println(`                                                                         `)

	fmt.Println(`   ____ _             ____      _   `)
	fmt.Println(`  / ___(_) ___ _   _ / ___| ___| |_ `)
	fmt.Println(` | |   | |/ _ \ | | | |  _ / _ \ __|`)
	fmt.Println(` | |___| |  __/ |_| | |_| |  __/ |_ `)
	fmt.Println(`  \____|_|\___|\__, |\____|\___|\__|`)
	fmt.Println(`               |___/                `)

	cmd := exec.Command("ping", "baidu.com")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error waiting for command to finish:", err)
		return
	}
}
