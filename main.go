package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type ChoiceObject struct {
	Choice string `json:"choice"`
}

type ResultResponse struct {
	Status string `json:"status"`
	Result string `json:"result"`
}

func main() {
	var rock ChoiceObject
	rock.Choice = "Rock"
	var scissors ChoiceObject
	scissors.Choice = "Scissors"
	var paper ChoiceObject
	paper.Choice = "Paper"

	router := mux.NewRouter()
	router.HandleFunc("/setChoice", setChoice).Methods("POST")
	// router.HandleFunc("/setResult", setResult).Methods("POST")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Access-Control-Allow-Origin"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Println("starting server...")
	log.Fatal(http.ListenAndServe(":8844", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

func setChoice(w http.ResponseWriter, r *http.Request) {
	var selectedChoice ChoiceObject
	if err := json.NewDecoder(r.Body).Decode(&selectedChoice); err != nil {
		log.Printf("failed to unmarshal message: %s", err)
		http.Error(w, "Bad request", http.StatusTeapot)
		return
	}

	defer r.Body.Close()
	log.Println("received:", selectedChoice)
	judge(&selectedChoice, w, r)
}

func judge(choice *ChoiceObject, w http.ResponseWriter, r *http.Request) {
	yourChoice := *choice

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	serverChoice := r1.Intn(3) + 1
	log.Println(serverChoice)
	userChoice := 0

	switch yourChoice.Choice {
	case "Rock":
		userChoice = 1
		break
	case "Scissors":
		userChoice = 2
		break
	case "Paper":
		userChoice = 3
		break
	default:
		break
	}

	response := &ResultResponse{
		Status: "Success",
		Result: "",
	}

	if (userChoice)%3+1 == serverChoice {
		response.Result = "Win"
		log.Println("Win")
	} else if (serverChoice)%3+1 == userChoice {
		response.Result = "Lose"
		log.Println("Lose")
	} else {
		response.Result = "Draw"
		log.Println("Draw")
	}

	res, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error: %T", err)
	}

	responseString := string(res)
	fmt.Fprint(w, responseString)
}
