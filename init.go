package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type NRICWallet struct {
	NRIC         string `json:"nric" binding:"required"`
	WalletAddress string `json:"wallet_address" binding:"required"`
}

var db *sql.DB

func main() {
	var err error

	// Connect to the database
	db, err = sql.Open("postgres", "host=db port=5432 user=admin password=admin dbname=mydb sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Test the database connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Create a new Gin router
	r := gin.Default()

	// Define the POST route for collecting NRIC and wallet address
	r.POST("/nric-wallet", func(c *gin.Context) {
		var nw NRICWallet

		// Bind the request body to the NRICWallet struct
		if err := c.ShouldBindJSON(&nw); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the NRIC already exists in the database
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM nric_wallet WHERE nric = $1", nw.NRIC).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "NRIC already exists"})
			return
		}

		// Check if the wallet address already exists in the database
		err = db.QueryRow("SELECT COUNT(*) FROM nric_wallet WHERE wallet_address = $1", nw.WalletAddress).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Wallet address already exists"})
			return
		}

		// Insert the new NRIC and wallet address into the database
		hasher := md5.New()
		//hashing is kind of one-way fn. compared to encryption, since with correct key we can decrypt.
		//However, in context of API resp. tix. it's okay to use HASH since, it cannot to manipulated by third-party or faked.
		//combination of encryption and hashing may be the best choice for achieving the desired security properties.
		hasher.Write([]byte(nw.NRIC + nw.WalletAddress))
		receipt := hex.EncodeToString(hasher.Sum(nil))
		_, err = db.Exec("INSERT INTO nric_wallet (nric, wallet_address, receipt) VALUES ($1, $2, $3)", nw.NRIC, nw.WalletAddress, receipt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Return the receipt
		c.JSON(http.StatusCreated, gin.H{"receipt": receipt})
	})

	// Start the server
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
