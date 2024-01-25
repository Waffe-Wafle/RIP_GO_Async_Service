package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ProbabilityResponse структура для ответа
type ProbabilityResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ProbabilityHandler обработчик запроса
func ProbabilityHandler(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	requestID, ok := requestData["id"].(string)
	if !ok || requestID == "" {
		http.Error(w, "No request id", http.StatusForbidden)
		return
	}

	// Ответ 200 и выполнение работы в асинхронном режиме
	go doWork(requestID)

	response := ProbabilityResponse{Message: "OK"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// doWork функция для выполнения работы
func doWork(requestID string) {
	time.Sleep(7 * time.Second)

	var result string
	// if rand.Float64() < 0.7 {
	// 	result = "paid"
	// } else {
	// 	result = "opened"
	// }
	result = "paid"

	payload := map[string]interface{}{
		"key":    "gfdswaqASFGHGFD",
		"id":     requestID,
		"status": result,
	}

	// Создание HTTP-запроса
	reqData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:8000/api/soft_loading_api/v1/asinc_pay_service/", bytes.NewBuffer(reqData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	// Отправка HTTP-запроса
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Чтение ответа
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	// Вывод ответа
	fmt.Println("Response:", response)
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/pay", ProbabilityHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	fmt.Println("Server is running on :8080")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
