package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/timeless2211/go-http-server/internal/database"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	// Get the optional author_id query parameter
	authorIDStr := r.URL.Query().Get("author_id")

	var chirps []database.Chirp
	var err error

	// If author_id is provided, filter by author
	if authorIDStr != "" {
		authorID, parseErr := uuid.Parse(authorIDStr)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id", parseErr)
			return
		}
		chirps, err = cfg.db.GetChirpsByAuthor(r.Context(), authorID)
	} else {
		// Otherwise, get all chirps
		chirps, err = cfg.db.GetChirps(r.Context())
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get chirps", err)
		return
	}

	apiChirps := make([]Chirp, len(chirps))
	for i, dbChirp := range chirps {
		apiChirps[i] = Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
	}

	// Get the optional sort query parameter (default to "asc")
	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "" {
		sortOrder = "asc"
	}

	// Sort chirps by created_at
	sort.Slice(apiChirps, func(i, j int) bool {
		if sortOrder == "desc" {
			return apiChirps[i].CreatedAt.After(apiChirps[j].CreatedAt)
		}
		// Default to ascending order
		return apiChirps[i].CreatedAt.Before(apiChirps[j].CreatedAt)
	})

	respondWithJSON(w, http.StatusOK, apiChirps)
}
