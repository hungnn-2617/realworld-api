package scripts

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/realworld-api/internal/config"
	"github.com/realworld-api/internal/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config.Load()

	token, err := utils.GenerateToken(1, "testuser", "test@example.com")
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	fmt.Println("Generated Token:")
	fmt.Println(token)
	fmt.Println()

	claims, err := utils.ValidateToken(token)
	if err != nil {
		log.Fatalf("Failed to validate token: %v", err)
	}

	fmt.Println("Token Claims:")
	fmt.Printf("  UserID: %d\n", claims.UserID)
	fmt.Printf("  Username: %s\n", claims.Username)
	fmt.Printf("  Email: %s\n", claims.Email)
	fmt.Printf("  ExpiresAt: %v\n", claims.ExpiresAt)

	fmt.Println()
	fmt.Println("Password Hashing Test:")
	password := "secretpassword123"
	hashed, err := utils.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	fmt.Printf("  Original: %s\n", password)
	fmt.Printf("  Hashed: %s\n", hashed)
	fmt.Printf("  Match: %v\n", utils.CheckPassword(hashed, password))
	fmt.Printf("  Wrong password match: %v\n", utils.CheckPassword(hashed, "wrongpassword"))
}
