package main

import (
	"context"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/videocoin/marketplace/internal/auth"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/model"
	"log"
	"os"
	"strings"
)

func main() {
	ctx := context.Background()

	dbURI := os.Getenv("DBURI")
	if dbURI == "" {
		dbURI = "host=127.0.0.1 port=5432 dbname=marketplace sslmode=disable"
	}

	ds, err := datastore.NewDatastore(ctx, dbURI)
	if err != nil {
		log.Fatalf("failed to init datastore: %s", err)
	}

	authSecret := os.Getenv("AUTH_SECRET")
	if authSecret == "" {
		authSecret = "secret"
	}

	address := os.Getenv("ACCOUNT_ADDRESS")
	isNew := strings.ToLower(os.Getenv("NEW_ACCOUNT"))

	var account *model.Account

	if !ethcommon.IsHexAddress(address) {
		log.Fatal("invalid account address")
	}

	if isNew == "y" {
		account = &model.Account{Address: strings.ToLower(address)}
		err := ds.Accounts.Create(ctx, account)
		if err != nil {
			log.Fatalf("failed to get account by address: %s", err)
		}
	} else {
		account, err = ds.Accounts.GetByAddress(ctx, address)
		if err != nil {
			log.Fatalf("failed to get account by address: %s", err)
		}
	}

	if account != nil {
		token, err := auth.CreateAuthToken(ctx, authSecret, account)
		if err != nil {
			log.Fatalf("failed to generate auth token: %s", err)
		}

		fmt.Printf("Token: %s\n", token)
	}
}
