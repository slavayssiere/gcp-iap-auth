package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/slavayssiere/gcp-iap-auth/jwt"
)

type userIdentity struct {
	Subject string `json:"sub,omitempty"`
	Email   string `json:"email,omitempty"`
}

func authHandler(res http.ResponseWriter, req *http.Request) {
	claims, err := jwt.RequestClaims(req, cfg)
	if err != nil {
		if claims == nil || len(claims.Email) == 0 {
			log.Printf("Failed to authenticate (%v)\n", err)
		} else {
			log.Printf("Failed to authenticate %q (%v)\n", claims.Email, err)
		}
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := &userIdentity{
		Subject: claims.Subject,
		Email:   claims.Email,
	}
	expiresAt := time.Unix(claims.ExpiresAt, 0).UTC()
	log.Printf("Authenticated %q (token expires at %v)\n", user.Email, expiresAt)
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(user)
}
