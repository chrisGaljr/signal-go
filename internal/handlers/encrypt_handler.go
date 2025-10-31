package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"signal/main/internal/utils"
)

type EncryptRequestBody struct {
	NewCipherKey         string `json:"newCipherKey"`
	OldDecryptedPassword string `json:"oldPassword"`
}

func EncryptionRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Error: invalid method", http.StatusBadRequest)
		return
	}

	var body EncryptRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Error decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	oldKey := []byte(os.Getenv("CIPHER_KEY"))
	oldPsw, err := base64.StdEncoding.DecodeString(body.OldDecryptedPassword)
	if err != nil {
		http.Error(w, "Error decoding old password: "+err.Error(), http.StatusBadRequest)
		return
	}

	password, err := utils.Decrypt(oldKey, oldPsw)
	if err != nil {
		http.Error(w, "Error decrypting old password: "+err.Error(), http.StatusBadRequest)
		return
	}

	newPassword, err := utils.HashText([]byte(body.NewCipherKey), password)
	if err != nil {
		http.Error(w, "Error getting new cipher: "+err.Error(), http.StatusBadRequest)
		return
	}

	http.ResponseWriter.Write(w, []byte(newPassword))
}
