package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var ids []int

func main() {
	rand.Seed(time.Now().UnixNano())

	fileUser, err := os.Create("user_data.sql")
	if err != nil {
		panic(err)
	}
	defer fileUser.Close()

	for i := 0; i <= 1000; i++ {
		id := i
		balance := rand.Float64()*400 + 100

		fmt.Fprintf(fileUser, "INSERT INTO users (id, balance) VALUES ('%d', %.2f);\n", id, balance)
		ids = append(ids, id)
	}

	fmt.Println("SQL данные записаны в файл seed_data_ticket.sql")

	// записываем ids в JSON
	fileJSON, err := os.Create("user_ids.json")
	if err != nil {
		panic(err)
	}
	defer fileJSON.Close()

	encoder := json.NewEncoder(fileJSON)
	encoder.SetIndent("", "  ")
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
