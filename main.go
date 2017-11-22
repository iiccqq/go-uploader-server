package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

func main() {
	fmt.Println("本地ip有以下，请在手机端设置对应的ip，确保电脑80端口开放及没有被占用")
	printLocalIP()
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)
	http.ListenAndServe(":80", nil)
}

func upload(w http.ResponseWriter, r *http.Request) {
	filePath := r.FormValue("filePath")

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()
	var sep string
	if os.IsPathSeparator('\\') {
		sep = "\\"
	} else {
		sep = "/"
	}
	wd, _ := os.Getwd()
	var dirName = wd + sep + filePath + sep
	fmt.Println("文件目录为" + dirName)
	_, err = os.Stat(dirName)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(dirName, os.ModePerm)
			fmt.Println("创建成功" + dirName)
		}
	}
	f1, err := os.OpenFile(dirName+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("打开文件失败")
	}

	defer f1.Close()
	io.Copy(f1, file)
	fmt.Fprintln(w, "upload ok!")
}

func index(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(tpl))
}

func printLocalIP() {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
			}
		}
	}
}

const tpl = `<html>
<head>
<title>上传文件</title>
</head>
<body>
<form enctype="multipart/form-data" action="/upload" method="post">
 <input type="file" name="uploadfile" />
 <input type="hidden" name="filePath" value="files"/>
 <input type="submit" value="upload" />
</form>
</body>
</html>`
