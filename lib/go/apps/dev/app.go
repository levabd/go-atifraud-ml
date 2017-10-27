package main

import (
	"log"
	"github.com/valyala/fasthttp"
	"os/exec"
	"fmt"
	"os"
	"strings"
	"github.com/levabd/go-atifraud-ml/lib/go/helpers"
)

var (
	connection                       = &fasthttp.Client{}
	req                              = fasthttp.AcquireRequest()
	resp                             = fasthttp.AcquireResponse()
)
func init() {
	err := helpers.LoadEnv()
	if err != nil {
		log.Fatalln(err)
	}
}
func main() {

	APP_ROOT_DIR:=os.Getenv("APP_ROOT_DIR")

	cmd := exec.Command(APP_ROOT_DIR+"/lib/python/train")
	cmd.Wait()
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error python exec: ", err)
		os.Exit(-1)
	}
	if strings.Contains(string(out), "Education finished") {
		log.Println("Education finished")
	} else {
		log.Println("Problem detected while educate model")
	}
}

