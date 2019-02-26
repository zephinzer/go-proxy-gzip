// this is an echoserver that can be used to verify requests sent by the
// main proxy-gzip server. this server listens on port 1339, which is also
// the default for the main application if you're using the docker-compose
// to spin it up
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// ListenInterface defines the interface the server will listen on
const ListenInterface = "0.0.0.0:1339"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("-----------------------------------------")
		log.Println("incoming request at", r.URL.EscapedPath())

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Panic(err)
		}
		defer r.Body.Close()

		response := map[string]interface{}{
			"source":  "test-server",
			"host":    r.Host,
			"method":  r.Method,
			"path":    r.URL.EscapedPath(),
			"headers": r.Header,
			"body":    string(body),
		}

		log.Println(response)
		responseData, err := json.Marshal(response)
		if err != nil {
			log.Println("error: ", err.Error())
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
		} else {
			log.Println("successful response made")
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(responseData))
		}
		log.Println("-----------------------------------------")
	})
	log.Printf("listening on 'https://%s'", ListenInterface)
	err := http.ListenAndServeTLS(ListenInterface, "echoserver.crt", "echoserver.key", nil)
	if err != nil {
		log.Panic(err)
	}
}
