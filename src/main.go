package main

import (
	"encryption-tools/src/services"
	"encryption-tools/src/utils"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	var mode string
	var input string
	var keyStr string
	var password string
	var inputFile string
	var outputFile string
	var keyLabel string
	var keyID string

	flag.StringVar(&mode, "mode", "encrypt", "Mode: encrypt/decrypt/encrypt-file/decrypt-file/encrypt-password/decrypt-password/list-keys")
	flag.StringVar(&input, "input", "", "Text to process")
	flag.StringVar(&keyStr, "key", "", "Hex-encoded 32-byte key")
	flag.StringVar(&password, "password", "", "Password for encryption/decryption")
	flag.StringVar(&inputFile, "input-file", "", "Input file path")
	flag.StringVar(&outputFile, "output-file", "", "Output file path")
	flag.StringVar(&keyLabel, "key-label", "", "Label for key storage")
	flag.StringVar(&keyID, "key-id", "", "Key ID from key storage")
	flag.Parse()

	keystore, err := utils.NewKeyStore(filepath.Join("keys", "keystore.json"))
	if err != nil {
		fmt.Printf("Error initializing keystore: %v\n", err)
		os.Exit(1)
	}

	var key []byte

	if keyID != "" {
		key, err = keystore.GetKey(keyID)
		if err != nil {
			fmt.Printf("Error retrieving key: %v\n", err)
			os.Exit(1)
		}
	} else if keyStr != "" {
		key, err = hex.DecodeString(keyStr)
		if err != nil || len(key) != 32 {
			fmt.Println("Invalid key. Must be 32 bytes hex-encoded")
			os.Exit(1)
		}
	}

	switch mode {
	case "encrypt", "decrypt":
		if key == nil {
			key = utils.GenerateKey()
			keyID, err = keystore.AddKey(key, keyLabel)
			if err != nil {
				fmt.Printf("Error storing key: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Generated key ID (save this): %s\n", keyID)
			fmt.Printf("Generated key (save this): %x\n", key)
		}

		service := services.NewEncryptionService(key)

		if mode == "encrypt" {
			result, err := service.Encrypt([]byte(input))
			if err != nil {
				fmt.Printf("Encryption error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Encrypted: %s\n", result)
		} else {
			result, err := service.Decrypt(input)
			if err != nil {
				fmt.Printf("Decryption error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Decrypted: %s\n", result)
		}

	case "encrypt-file", "decrypt-file":
		if key == nil {
			key = utils.GenerateKey()
			keyID, err = keystore.AddKey(key, keyLabel)
			if err != nil {
				fmt.Printf("Error storing key: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Generated key ID (save this): %s\n", keyID)
			fmt.Printf("Generated key (save this): %x\n", key)
		}

		service := services.NewFileEncryptionService(key)

		if mode == "encrypt-file" {
			err = service.EncryptFile(inputFile, outputFile)
			if err != nil {
				fmt.Printf("File encryption error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("File encrypted successfully")
		} else {
			err = service.DecryptFile(inputFile, outputFile)
			if err != nil {
				fmt.Printf("File decryption error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("File decrypted successfully")
		}

	case "encrypt-password", "decrypt-password":
		if password == "" {
			fmt.Println("Password is required for password-based encryption")
			os.Exit(1)
		}

		service := services.NewPasswordEncryptionService()

		if mode == "encrypt-password" {
			result, err := service.Encrypt(input, password)
			if err != nil {
				fmt.Printf("Password encryption error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Encrypted: %s\n", result)
		} else {
			result, err := service.Decrypt(input, password)
			if err != nil {
				fmt.Printf("Password decryption error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Decrypted: %s\n", result)
		}

	case "list-keys":
		keys := keystore.ListKeys()
		if len(keys) == 0 {
			fmt.Println("No keys stored")
			return
		}
		fmt.Println("Stored keys:")
		for _, k := range keys {
			fmt.Printf("ID: %s, Created: %s, Label: %s\n", k.KeyID, k.CreatedAt.Format("2006-01-02 15:04:05"), k.Label)
		}

	default:
		fmt.Println("Invalid mode. Use 'encrypt', 'decrypt', 'encrypt-file', 'decrypt-file', 'encrypt-password', 'decrypt-password', or 'list-keys'")
		os.Exit(1)
	}
}
