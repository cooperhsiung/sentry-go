package notice

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
)

func SendWx(message string) {

	r := gorequest.New()

	data := map[string]interface{}{
		"type":    "text",
		"touser":  "--",
		"toparty": "",
		"agentid": "1000003",
		"content": message,
	}

	_, body, errs := r.Post("https://--/api/wx_v2.php").Type("form").Send(data).End()
	if errs != nil {
		fmt.Println(errs)
	}

	fmt.Println(body)
}
