package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type ClaimRequest struct {
	NRIC          string `json:"nric" binding:"required"`
	WalletAddress string `json:"wallet_address" binding:"required"`
}

type ClaimResponse struct {
	Receipt string `json:"receipt"`
}

func main() {
	db, err := sql.Open("postgres", "postgres://user:password@postgres:5432/mydb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.POST("/claim", func(c *gin.Context) {
		var req ClaimRequest
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if NRIC is unique
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM claims WHERE nric = $1", req.NRIC).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "NRIC is already associated with a wallet address"})
			return
		}

		// Check if wallet address is unique
		err = db.QueryRow("SELECT COUNT(*) FROM claims WHERE wallet_address = $1", req.WalletAddress).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
			return
		}
		if count > 0 {
			//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet address is already associated with an NRIC"})
			return
		}

		// Insert claim into database
		_, err = db.Exec("INSERT INTO claims (nric, wallet_address) VALUES ($1, $2)", req.NRIC, req.WalletAddress)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert claim into database"})
			return
		}

		// Generate receipt hash
		hasher := md5.New()
		hasher.Write([]byte(fmt.Sprintf("%v", req)))
		receipt := hex.EncodeToString(hasher.Sum(nil))

		c.JSON(http.StatusOK, ClaimResponse{Receipt: receipt})
	})

	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

//alternate way to connect DB, declare const() method which holds this var.
func connectToDB() (*sql.DB, error) {
	dbHost := os.Getenv(dbHostKey)
	dbPort := os.Getenv(dbPortKey)
	dbUser := os.Getenv(dbUserKey)
	dbPassword := os.Getenv(dbPasswordKey)
	dbName := os.Getenv(dbNameKey)

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		return nil, errors.New("database configuration not provided")
	}

}
