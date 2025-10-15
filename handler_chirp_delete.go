package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/timeless2211/go-http-server/internal/auth"
)

func (cfg *apiConfig) handlerChirpDelete(w http.ResponseWriter, r *http.Request) {
	// 1. Get and validate the bearer token
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid authorization header", err)
		return
	}

	// 2. Validate JWT and extract user ID
	userID, err := auth.ValidateJWT(tokenString, cfg.authSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	// 3. Get chirp ID from URL path
	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	// 4. Get the chirp from the database to check if it exists and who owns it
	chirp, err := cfg.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	// 5. Check if the user is the author of the chirp
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You are not authorized to delete this chirp", nil)
		return
	}

	// 6. Delete the chirp
	err = cfg.db.DeleteChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp", err)
		return
	}

	// 7. Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
