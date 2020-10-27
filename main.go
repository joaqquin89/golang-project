package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type Person struct {
	ID        string   `json:"id,omitempty"`
	FirstName string   `json:"firstname,omitempty"`
	LastName  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Country string `json:"country,omitempty"`
}

var people []Person

func GetPeopleEndpoint(w http.ResponseWriter, req *http.Request) {

	f, err := os.Open("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		fmt.Println("------------------")
		dec := json.NewDecoder(strings.NewReader(scanner.Text()))
		for {
			var doc Person
			err := dec.Decode(&doc)
			if err == io.EOF {
				// all done
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%+v\n", doc)
			json.NewEncoder(w).Encode(doc)
		}
	}

}

func GetPersonEndpoint(w http.ResponseWriter, req *http.Request) {

	params := mux.Vars(req)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
		json.NewEncoder(w).Encode(&Person{})
	}
}

func CreatePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var person Person
	fmt.Println(params)
	_ = json.NewDecoder(req.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

func DeletePersonEndpoint(w http.ResponseWriter, req *http.Request) {

}

func main() {
	router := mux.NewRouter()
	people = append(people, Person{ID: "1", FirstName: "Ryan", LastName: "Ray", Address: &Address{City: "dubling", State: "california", Country: "USA"}})
	people = append(people, Person{ID: "2", FirstName: "joaquin", LastName: "jachura", Address: &Address{City: "birmingham", State: "alabama", Country: "USA"}})

	// verificar si existe  el file data.txt
	if _, err := os.Stat("data.txt"); err == nil {
		fmt.Printf("File exists we don't need to create another file\n")
		fmt.Printf("Reading the file\n")
	} else {
		//en caso de no existir el file , lo que hacemos es crear dicho file
		err := ioutil.WriteFile("data.txt", []byte("first line\n"), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	for i, s := range people {
		// convert the output in json obj
		out, err := json.Marshal(s)
		if err != nil {
			panic(err)
		}
		// read the whole file at once
		b, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		str := string(b)
		// //check whether s contains substring text
		fmt.Println(strings.Contains(str, string(out)))
		if strings.Contains(str, string(out)) == false {
			fmt.Println(i)
			file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err)
			}
			defer file.Close()

			file.WriteString(string(out))
			file.WriteString("\n")
		}

	}
	*fmt.Println("done")

	//endpoints
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))
}
