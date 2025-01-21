# Encryption Tools

A secure and efficient encryption tool built with Go.

## Features
- AES-GCM encryption for maximum security
- Multiple encryption methods:
  - Key-based encryption
  - Password-based encryption
  - File encryption
- Secure key management system
- Base64 encoded output
- Command-line interface

## Quick Start

1. Build the project:
```batch
run.bat build
```

2. Text Encryption:
```batch
run.bat run -mode encrypt -input "your secret text" -key-label "my-key"
```

3. Text Decryption:
```batch
run.bat run -mode decrypt -key-id YOUR_KEY_ID -input ENCRYPTED_TEXT
```

4. File Encryption:
```batch
run.bat run -mode encrypt-file -input-file "original.txt" -output-file "encrypted.bin" -key-label "file-key"
```

5. Password-based Encryption:
```batch
run.bat run -mode encrypt-password -input "secret text" -password "your-secure-password"
```

6. List Stored Keys:
```batch
run.bat run -mode list-keys
```

## Build Requirements
- Go 1.21 or higher

## Security Features
- AES-GCM for authenticated encryption
- PBKDF2 for password-based key derivation
- Secure random number generation
- Safe key storage system
