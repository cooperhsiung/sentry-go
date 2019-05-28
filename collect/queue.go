package collect

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
)

func GetTaoBao() {
	r := gorequest.New()
	url := "http://mqadmin:-=@-:15672/api/queues/%2f/-"
	_, body, errs := r.Get(url).Type("json").End()

	if errs != nil {
		log.Fatal(errs)
	}

	fmt.Println(body)
}

func GetOperator() {

	r := gorequest.New()
	url := "http://guest:guest@-:15672/api/queues/%2f/-"
	_, body, errs := r.Get(url).Type("json").End()

	if errs != nil {
		log.Fatal(errs)
	}

	fmt.Println(body)

}
