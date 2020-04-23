package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"
)

func main() {

	workdir := "app-arm64-v8a-debug"

	cmd := exec.Command("rm", "-rf", workdir)
	fmt.Printf("%s", "删除旧版本........")
	runComd(cmd)

	cmd = exec.Command("java", "-jar", "apktool_2.4.1.jar", "d", "/Users/zhui/AndroidStudioProjects/ClashForAndroid/app/build/outputs/apk/debug/"+workdir+".apk")
	fmt.Printf("%s", "正在反编译.........")
	runComd(cmd)

	b, _ := ioutil.ReadFile("./" + workdir + "/AndroidManifest.xml")
	cmd = exec.Command("rm", "-rf", "./"+workdir+"/AndroidManifest.xml")
	fmt.Printf("%s", "修改AndroidManifest.xml........")
	runComd(cmd)
	all := strings.ReplaceAll(string(b), " android:extractNativeLibs=\"false\"", "")
	_ = ioutil.WriteFile("./"+workdir+"/AndroidManifest.xml", []byte(all), 0644)

	cmd = exec.Command("rm", "-rf", "./"+workdir+"/lib")
	fmt.Printf("%s", "正在替换内核........")
	runComd(cmd)

	cmd = exec.Command("cp", "-r", "./lib", "./"+workdir+"/lib")
	fmt.Printf("%s", "替换内核成功........")
	runComd(cmd)

	cmd = exec.Command("java", "-jar", "apktool_2.4.1.jar", "b", workdir, "-o", "new.apk")
	fmt.Printf("%s", "正在重新打包........")
	runComd(cmd)

	format := time.Now().Format("20060102-150405")
	cmd = exec.Command("./apksigner", "sign", "--ks", "./clashr.jks", "--ks-key-alias", "key0", "--ks-pass", "pass:123456", "--out", "./clashR"+format+".apk", "./new.apk")
	fmt.Printf("%s", "重签名成功........")
	runComd(cmd)

	cmd = exec.Command("rm", "-rf", "./app-debug")
	runComd(cmd)
	cmd = exec.Command("rm", "-rf", "./new.apk")
	fmt.Printf("%s", "清理目录........")
	runComd(cmd)

}

func runComd(cmd *exec.Cmd) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(opBytes))
}
