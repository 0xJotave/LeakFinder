package testes

import "fmt"

func main() {
	// API Keys
	var apiKey = "sk_test_4eC39HqLyjWDarjtT1zdp7dc" // Exemplo de chave de API do Stripe
	var dbPassword = "p@ssw0rd123"                  // Exemplo de senha de banco de dados
	var jwtSecret = "s3cr3tkey"                     // Exemplo de segredo JWT

	// Print some info (simulando um log)
	fmt.Println("Running application...")
	fmt.Println("API Key:", apiKey)
	fmt.Println("Database Password:", dbPassword)
	fmt.Println("JWT Secret:", jwtSecret)

	// Simulando uso de variáveis sensíveis
	connectToDatabase(dbPassword)
}

func connectToDatabase(password string) {
	fmt.Println("Connecting to database with password:", password)
}
