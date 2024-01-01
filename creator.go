package curlfuzz

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/google/uuid"
)

type Data struct {
	Raw      string
	Name     string
	Payload  *string
	TempName string
	Key      string
}
type Config struct {
	ProxyIP       string
	Payload       string
	ProxyPort     string
	Verbose       bool
	URI           []string
	AddQueryParam string
	Temp          *template.Template
	TypeAttack    string
	Curl          string
	Key           string
}

func (c *Config) create_template(new_raw http.Request, newDumpRequest []byte, index string) {
	flag.Parse()
	home, _ := os.UserHomeDir()
	fmt.Println("home, ", home)

	// temp := template.Must(template.ParseFiles("builder-template/template.yaml"))
	var temp *template.Template
	files, err := os.ReadDir(fmt.Sprintf("%s/builder-template", home))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, file := range files {
		fmt.Println(file.Name())
		temp = template.Must(template.ParseFiles(fmt.Sprintf("%s/builder-template/%s", home, file.Name())))
		var config Config
		config.Temp = temp

		nucleiPATH := fmt.Sprintf("fuzzing-%s-%v", file.Name(), uuid.New())
		_, err := os.Stat("/tmp/" + c.Key)
		if os.IsNotExist(err) {
			errDir := os.MkdirAll("/tmp/"+c.Key, os.ModePerm)
			if errDir != nil {
				log.Fatal(err)
			}
		}
		file, err := os.Create("/tmp/" + c.Key + "/" + nucleiPATH + ".yaml")
		if err != nil {
			log.Fatal(err)
			os.Exit(3)
		}
		data_temp := Data{TempName: nucleiPATH, Payload: nil, Raw: fmt.Sprintf("%s", newDumpRequest), Name: nucleiPATH, Key: c.Key}
		err = config.Temp.Execute(file, data_temp)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Creating template in ", nucleiPATH)
	}

}
