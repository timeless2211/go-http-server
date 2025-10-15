package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/timeless2211/go-http-server/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	// 1. Validate API key
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid API key", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", nil)
		return
	}

	type WebhookData struct {
		UserID uuid.UUID `json:"user_id"`
	}

	type WebhookRequest struct {
		Event string      `json:"event"`
		Data  WebhookData `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := WebhookRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	// If the event is anything other than user.upgraded, respond with 204
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Upgrade the user to Chirpy Red
	_, err = cfg.db.UpgradeToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		// If the user can't be found, respond with 404
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	// Return 204 No Content on success
	w.WriteHeader(http.StatusNoContent)
}
