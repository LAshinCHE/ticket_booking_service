package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var ticketIDs []int

func main() {
	rand.Seed(time.Now().UnixNano())

	fileTicket, err := os.Create("ticket_data.sql")
	if err != nil {
		panic(err)
	}
	defer fileTicket.Close()
	fileIds, err := os.Create("ids")
	if err != nil {
		panic(err)
	}
	defer fileIds.Close()
	for i := 0; i <= 3000; i++ {
		id := i
		price := rand.Float64()*100 + 100
		available := true

		fmt.Fprintf(fileTicket, "INSERT INTO tickets (id, price, available) VALUES ('%d', %.2f, %t);\n", id, price, available)
		ticketIDs = append(ticketIDs, id)
	}

	fmt.Println("SQL данные записаны в файл seed_data_ticket.sql")
	saveIDs(ticketIDs, "ticket_ids.json")
}

func saveIDs(ids []int, name string) {
	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	err = encoder.Encode(ids)
	if err != nil {
		panic(err)
	}

	fmt.Println("Данные успешно сгенерированы.")
}

func escape(s string) string {
	return replace(s, "'", "''")
}

func replace(s, old, new string) string {
	return string([]rune(s))
}
