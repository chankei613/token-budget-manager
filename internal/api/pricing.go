package api

import (
	"encoding/json"
	"net/http"

	"github.com/chankei613/token-budget-manager/internal/db"
)

func (s *Server) ListPricing() ([]db.ModelPricing, error) {
	var rows []db.ModelPricing
	err := s.DB.Order("model_id asc").Find(&rows).Error
	return rows, err
}

type SetPricingInput struct {
	ModelID          string  `json:"model_id"`
	InputPricePer1M  float64 `json:"input_price_per_1m"`
	OutputPricePer1M float64 `json:"output_price_per_1m"`
}

// SetPricing はモデル価格を作成/更新する（upsert）。
func (s *Server) SetPricing(in SetPricingInput) (db.ModelPricing, error) {
	if in.ModelID == "" {
		return db.ModelPricing{}, &apiError{"model_id is required"}
	}
	row := db.ModelPricing{
		ModelID:          in.ModelID,
		InputPricePer1M:  in.InputPricePer1M,
		OutputPricePer1M: in.OutputPricePer1M,
	}
	if err := s.DB.Save(&row).Error; err != nil {
		return db.ModelPricing{}, err
	}
	return row, nil
}

func (s *Server) lookupPricing(modelID string) (db.ModelPricing, bool) {
	var row db.ModelPricing
	if err := s.DB.First(&row, "model_id = ?", modelID).Error; err != nil {
		return db.ModelPricing{}, false
	}
	return row, true
}

func computeCost(inputTokens, outputTokens int64, price db.ModelPricing) float64 {
	return float64(inputTokens)/1_000_000*price.InputPricePer1M +
		float64(outputTokens)/1_000_000*price.OutputPricePer1M
}

func (s *Server) httpListPricing(w http.ResponseWriter, r *http.Request) {
	rows, err := s.ListPricing()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) httpSetPricing(w http.ResponseWriter, r *http.Request) {
	var body SetPricingInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	row, err := s.SetPricing(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, row)
}
