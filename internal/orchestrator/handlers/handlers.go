package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"finalProject/internal/orchestrator/logic"
)

type Request struct {
	Expression string `json:"expression"`
}


var Expressions []Expression

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	var request Request
	_ = json.Unmarshal(body, &request)
	logic.NewEx(request.Expression)
}
