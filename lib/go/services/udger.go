package services

import (
	"github.com/udger/udger"
	"path/filepath"
	"os"
	"github.com/levabd/go-atifraud-ml/lib/go/helpers"
)

func init() {
	err := helpers.LoadEnv()
	if err != nil {
		Logger.Fatalln(err)
	}
}

var instantiated *udger.Udger = nil

func GetUdgerInstance() (*udger.Udger, error){
	if instantiated == nil {
		path_to_udger_db :=filepath.Join(os.Getenv("APP_ROOT_DIR"),"data", "db", "udgerdb_v3.dat")
		println("path_to_udger_db: ", path_to_udger_db)

		u, err := udger.New(path_to_udger_db)
		if err != nil {
			Logger.Fatalln(err)
			return nil, err;
		}
		instantiated = u;
	}
	return instantiated, nil;
}

//func IsCrawler(client_ip string, client_ua string) bool {
//	cmd := exec.Command("python3", filepath.Join(os.Getenv("APP_ROOT_DIR"),"lib", "python", "isCrawler.py"), client_ip, client_ua)
//	out, err := cmd.CombinedOutput()
//
//	if err != nil {
//		fmt.Println("error python exec: ", err)
//		os.Exit(-1)
//	}
//
//	return strings.TrimRight(string(out) , "\n")  == "True"
//}