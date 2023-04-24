package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// UserData : data structure for storing NRIC and wallet address
type UserData struct {
	NRIC         string `json:"nric"`
	WalletAddr   string `json:"walletAddr"`
	ReceiptHash  string `json:"receiptHash"`
}

//global database variable
// postgres running on 3306
var db *sql.DB

// db initialize
func initDB() error {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS userdata (
		id SERIAL PRIMARY KEY,
		nric TEXT UNIQUE,
		walletaddr TEXT UNIQUE,
		receipthash TEXT NOT NULL
	)`)
	return err
}

//enerate a receipt hash
func generateReceiptHash(userdata *UserData) (string, error) {
	if userdata.NRIC == "" || userdata.WalletAddr == "" {
		return "", errors.New("NRIC and walletAddr required")
	}
	hash := sha256.New()
	_, err := hash.Write([]byte(userdata.NRIC + userdata.WalletAddr))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

//handler for the POST /api/userdata endpoint
func handlePostUserData(c *gin.Context) {
	var userdata UserData
	if err := c.ShouldBindJSON(&userdata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Generate the receipt hash
	receiptHash, err := generateReceiptHash(&userdata)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userdata.ReceiptHash = receiptHash
	// Insert the userdata into the database
	_, err = db.Exec("INSERT INTO userdata (nric, walletaddr, receipthash) VALUES ($1, $2, $3)", userdata.NRIC, userdata.WalletAddr, userdata.ReceiptHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"receiptHash": userdata.ReceiptHash})
}

func main() {
	// Initialize the database
	if err := initDB(); err != nil {
		log.Fatal(err)
	}
	// Create a new Gin router
	r := gin.Default()
	// Define the POST /api/userdata endpoint
	r.POST("/api/userdata", handlePostUserData)
	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
