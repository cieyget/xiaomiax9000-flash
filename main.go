package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

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

// 实现自定义的WriteCloser接口，用于显示下载进度
type progressWriter struct {
	io.Writer
	total   int64 // 总字节数
	written int64 // 已写入字节数
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.Writer.Write(p)
	pw.written += int64(n)
	pw.printProgress()
	return n, err
}

func (pw *progressWriter) printProgress() {
	progress := float64(pw.written) / float64(pw.total) * 100
	fmt.Printf("\rDownloading：%.2f%%", progress)
}

func downloadFile(url string, filepath string) error {

	// 发起HTTP GET请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建目标文件
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 获取响应内容的长度
	contentLength := resp.Header.Get("Content-Length")
	total := int64(0)
	if contentLength != "" {
		total = stringsToInt64(contentLength)
	}

	// 创建自定义的progressWriter
	pw := &progressWriter{
		Writer: out,
		total:  total,
	}

	// 将HTTP响应的内容写入文件，并显示进度条
	_, err = io.Copy(pw, resp.Body)
	if err != nil {
		return err
	}
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Println("", green("Done！"))
	return nil
}

func calculateMD5(filepath string) (string, error) {
	// 打开文件
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 计算MD5校验和
	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	// 将校验和转换为字符串格式
	md5sum := fmt.Sprintf("%x", hash.Sum(nil))
	return md5sum, nil
}

func executeCommand(command string) (string, error) {
	// 执行命令
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func readFile(filepath string) (string, error) {
	// 读取文件内容
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func stringsToInt64(s string) int64 {
	var result int64
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return 0
	}
	return result
}

func getMD5FromAPI(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
func checkRootfsAndSwitch() {
	// 读取文件
	content, err := readFile("cmdline")
	if err != nil {
		fmt.Println("Readfile error:", err)
		return
	}
	// fmt.Println("Cmdline:", content)

	prefix := "ubi.mtd="
	rootfsStr := ""
	startIndex := strings.Index(content, prefix)
	if startIndex != -1 {
		startIndex += len(prefix)
		endIndex := strings.Index(content[startIndex:], " ")
		if endIndex == -1 {
			endIndex = len(content)
		} else {
			endIndex += startIndex
		}
		rootfsStr = strings.TrimSpace(content[startIndex:endIndex])
		// fmt.Printf("[%s]\n", rootfsStr)
	} else {
		fmt.Println("Cmdline parameter not obtained")
	}
	if rootfsStr != "rootfs" {
		fmt.Println("Check that the current boot system is not rootfs and automatically switch to rootfs")

		commands := []string{
			"nvram set flag_last_success=0",
			"nvram set flag_boot_rootfs=0",
			"nvram set flag_try_sys1_failed=0",
			"nvram set flag_try_sys2_failed=0",
			"nvram commit",
		}

		for _, cmd := range commands {
			err := exec.Command("sh", "-c", cmd).Run()
			if err != nil {
				fmt.Printf("Error executing command: %s, error: %v\n", cmd, err)
			} else {
				fmt.Printf("Command executed successfully: %s\n", cmd)
			}
		}
		fmt.Println("Reboot after 5 seconds")
		// 等待 5 秒
		time.Sleep(5 * time.Second)

		// 执行 reboot 命令
		// err = exec.Command("reboot").Run()
		// if err != nil {
		// 	fmt.Printf("Error executing reboot command: %v\n", err)
		// } else {
		// 	fmt.Println("Reboot command executed successfully")
		// }

		return
	}
	fmt.Printf("Boot: %s \n", rootfsStr)
}
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	return false
}
func main() {
	filename := "/tmp/temp"
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

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
	// 检查看是否需要切换rootfs
	fmt.Println("Check boot partition....")
	checkRootfsAndSwitch()

	if fileExists(filename) {
		err := os.Remove(filename)
		if err != nil {
			fmt.Printf("删除文件 %s 时发生错误: %v\n", filename, err)
		} else {
			fmt.Printf("CleanFile%s\n", filename)
		}
	}
	err := downloadFile("https://syun-1251974457.cos.ap-chengdu.myqcloud.com/xiaomi/Redmi_AX6000_RB06/firmware/xwrt/x-wrt-23.04-b202305152359-mediatek-filogic-xiaomi_redmi-router-ax6000-stock-initramfs-factory.ubi", filename)
	if err != nil {
		fmt.Printf("%s %v \n", red("Download Failed！"), err)
		return
	}

	// 获取API返回的MD5校验和
	apiMD5, err := getMD5FromAPI("https://syun-1251974457.cos.ap-chengdu.myqcloud.com/xiaomi/Redmi_AX6000_RB06/firmware/xwrt/x-wrt-23.04-b202305152359-mediatek-filogic-xiaomi_redmi-router-ax6000-stock-initramfs-factory.ubi.md5")
	// 字符串第一列
	parts := strings.Split(apiMD5, " ")
	firstColumn := parts[0]

	if err != nil {
		fmt.Printf("%s %v \n", red("Get MD5 Failed！"), err)
		return
	}

	// 计算本地文件的MD5校验和
	localMD5, err := calculateMD5(filename)
	if err != nil {
		fmt.Printf("%s %v \n", red("Error calculating file md5:"), err)
		return
	}

	// 对比MD5校验和
	fmt.Printf("Verify MD5 \n")
	if firstColumn != localMD5 {
		// 校验文件MD5如果失败则尝试重新下载
		fmt.Printf("%s %s %s %s\n", red("Failed to verify md5"), red(firstColumn), red("≠"), red(localMD5))
		return
	}

	fmt.Printf("%s %s\n", green("Verify MD5 successfully! :"), green(firstColumn))

	// 执行命令
	output, err := executeCommand("top")
	if err != nil {
		fmt.Println("执行命令出错:", err)
		return
	}
	fmt.Println("命令输出:", output)

}
