package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	vaultAPI "github.com/hashicorp/vault/api"
)

type vaultBlob struct {
	LeaseDuration time.Duration `json:"lease_duration"`
	LeaseID       string        `json:"lease_id"`
	Data          struct {
		AccessKey     string `json:"access_key"`
		SecretKey     string `json:"secret_key"`
		SecurityToken string `json:"security_token"`
		Expiration    time.Time
	}
}

type awsBlob struct {
	Version         int    `json:"Version"`
	AccessKeyID     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"SessionToken"`
	Expiration      string `json:"Expiration"`
}

func main() {
	in, ok := os.LookupEnv("VACH_VAULT_BLOB")
	if !ok {
		log.Println("VACH_VAULT_BLOB must be set!")
		os.Exit(2)
	}

	inF, err := os.Open(in)
	if err != nil {
		log.Printf("Error opening input file: %s", err)
		os.Exit(2)
	}
	defer inF.Close()

	vaultBlob := vaultBlob{}

	dec := json.NewDecoder(inF)
	if err := dec.Decode(&vaultBlob); err != nil {
		log.Printf("Error decoding vault blob: %s", err)
		os.Exit(2)
	}

	vault, err := vaultAPI.NewClient(vaultAPI.DefaultConfig())
	if err != nil {
		log.Printf("Error establishing vault client: %s", err)
		os.Exit(2)
	}

	info, err := vault.Sys().Lookup(vaultBlob.LeaseID)
	if err != nil {
		log.Printf("Error obtaining information from vault: %s", err)
		os.Exit(2)
	}

	// We do this unchecked here because the explit type
	// referenced from Vault is a string.
	expiry := info.Data["expire_time"].(string)
	t, err := time.Parse(time.RFC3339Nano, expiry)
	if err != nil {
		log.Printf("Called vault, had a bad time: %s", err)
		os.Exit(2)
	}
	vaultBlob.Data.Expiration = t

	// Now we marshal between the two types
	awsBlob := awsBlob{
		Version:         1,
		AccessKeyID:     vaultBlob.Data.AccessKey,
		SecretAccessKey: vaultBlob.Data.SecretKey,
		SessionToken:    vaultBlob.Data.SecurityToken,
		Expiration:      vaultBlob.Data.Expiration.Format(time.RFC3339),
	}

	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(awsBlob); err != nil {
		log.Printf("Error encoding output: %s", err)
		os.Exit(2)
	}
	os.Exit(0)
}
