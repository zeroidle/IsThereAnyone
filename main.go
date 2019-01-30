package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os/exec"
)

type Device struct {
	code       string
	name       string
	macaddress string
}

func ConnectRedis() redis.Client {
	client := redis.NewClient(&redis.Options{
		//Addr: RedisAddr,
		//Password: RedisPassword, // no password set
		Addr:     RedisAddr,
		Password: RedisPassword, // no password set

		DB: 0, // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	return *client
}

func GetDataFromRedis(client redis.Client) {
	err := client.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	result, err := client.LRange("isthereanyone", 0, -1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	for i, key := range result {
		fmt.Println(i, key)
		var dev map[string]interface{}
		err := json.Unmarshal([]byte(key), &dev)
		if err != nil {
			panic(err)
		}
		fmt.Println(dev["macaddress"])

	}
}

func main() {
	redisConn := ConnectRedis()
	GetDataFromRedis(redisConn)

	router := mux.NewRouter()
	router.HandleFunc("/test/{code}", GetData).Methods("GET")

	http.ListenAndServe(":9801", httpHandler(router))
}

func httpHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print(r.RemoteAddr, " ", r.Proto, " ", r.Method, " ", r.URL)
		handler.ServeHTTP(w, r)
	})
}

var devices []Device

type event struct {
	code        string
	Title       string
	Description string
}

func GetData(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)
	for _, i := range devices {
		if i.code == p["code"] {
			json.NewEncoder(w).Encode(L2ping(i.macaddress))

			return
		}
	}
	json.NewEncoder(w).Encode(&event{})
}

func L2ping(mac string) bool {
	log.Println("Checking ", mac)
	cmd := exec.Command("l2ping", "-c", "1", mac)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}
