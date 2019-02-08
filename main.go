package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type Device struct {
	code       int
	name       string
	macaddress string
	result     bool
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
		fmt.Println(dev["code"])

		devices = append(devices, Device{
			code:       int(dev["code"].(float64)),
			name:       dev["name"].(string),
			macaddress: dev["macaddress"].(string),
			result:     false,
		})
	}
}

func main() {
	runtime.GOMAXPROCS(4) // Use multi-cores

	redisConn := ConnectRedis()
	GetDataFromRedis(redisConn)

	for i := 0; i < len(devices); i++ {
		go Gathering(i)
	}
	router := mux.NewRouter()
	router.HandleFunc("/check/{code}", GetData).Methods("GET")
	router.HandleFunc("/view", ViewPage).Methods("GET")
	router.HandleFunc("/play", PlaySong).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":9801", httpHandler(router))
}

func httpHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print(r.RemoteAddr, " ", r.Proto, " ", r.Method, " ", r.URL)
		handler.ServeHTTP(w, r)
	})
}

var devices []Device

func ViewPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/view_page.html")
	t.Execute(w, nil)
}

func PlaySong(w http.ResponseWriter, r *http.Request) {
	fmt.Println("play")
}

type event struct {
	code        string
	Title       string
	Description string
}

func Gathering(pos int) {
	device := devices[pos]
	fmt.Println("-----------------------", device.name)
	for ok := true; ok; ok = true {
		result := L2ping(device.macaddress)

		devices[pos].result = result
		fmt.Println("scan bluetooth device ", device.name, " ", device.macaddress, " ", devices[pos].result)
		if result == true {
			time.Sleep(time.Second * 10)
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}

	fmt.Println("aa")
}
func GetData(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)
	for _, i := range devices {
		if strconv.Itoa(i.code) == p["code"] {
			rtnData := make(map[string]interface{})
			rtnData["code"] = i.code
			rtnData["result"] = i.result
			fmt.Println(i)

			json.NewEncoder(w).Encode(rtnData)

			return
		}
	}
	fmt.Println("nothing")
	json.NewEncoder(w).Encode(&event{})
}

func L2ping(mac string) bool {
	log.Println("Checking ", mac)
	ctx, cancel := context.WithTimeout(context.Background(), 8500*time.Millisecond)
	defer cancel()
	cmd := exec.CommandContext(ctx, "./l2ping", "-c", "1", mac)
	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("Command timed out")
		return false
	}
	if err != nil {
		log.Println(string(output), " ", err)
		return false
	}
	fmt.Println(string(output))
	return true

}
