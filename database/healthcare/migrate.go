package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"column:id;primaryKey"`
	Name      string `gorm:"column:name"`
	Password  string `gorm:"column:password"`
	Token     string `gorm:"column:token"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

type Contact struct {
	ID        string `gorm:"column:id;primaryKey"`
	FirstName string `gorm:"column:first_name"`
	LastName  string `gorm:"column:last_name"`
	Email     string `gorm:"column:email"`
	Phone     string `gorm:"column:phone"`
	UserID    string `gorm:"column:user_id"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

type Address struct {
	ID         string `gorm:"column:id;primaryKey"`
	ContactID  string `gorm:"column:contact_id"`
	Street     string `gorm:"column:street"`
	City       string `gorm:"column:city"`
	Province   string `gorm:"column:province"`
	PostalCode string `gorm:"column:postal_code"`
	Country    string `gorm:"column:country"`
	CreatedAt  int64  `gorm:"column:created_at"`
	UpdatedAt  int64  `gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "users"
}

func (Contact) TableName() string {
	return "contacts"
}

func (Address) TableName() string {
	return "addresses"
}

func main() {
	// Database connection
	dsn := "host=localhost user=postgres password=password dbname=healthcare port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to database successfully!")

	// Run migrations
	log.Println("Running migrations...")
	err = db.AutoMigrate(&User{}, &Contact{}, &Address{})
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	log.Println("Migrations completed successfully!")

	// Seed data
	log.Println("Seeding database with test data...")
	seedData(db)
	log.Println("Database seeding completed!")
}

func seedData(db *gorm.DB) {
	now := time.Now().Unix()

	// Hash password for test users
	hashedPassword1, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashedPassword2, _ := bcrypt.GenerateFromPassword([]byte("SecurePass456!"), bcrypt.DefaultCost)

	// Create test users
	users := []User{
		{
			ID:        "john.doe",
			Name:      "John Doe",
			Password:  string(hashedPassword1),
			Token:     "",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "jane.smith",
			Name:      "Jane Smith",
			Password:  string(hashedPassword2),
			Token:     "",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "bob.wilson",
			Name:      "Bob Wilson",
			Password:  string(hashedPassword1),
			Token:     "",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for _, user := range users {
		result := db.FirstOrCreate(&user, User{ID: user.ID})
		if result.Error != nil {
			log.Printf("Error creating user %s: %v", user.ID, result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("Created user: %s", user.ID)
		} else {
			log.Printf("User already exists: %s", user.ID)
		}
	}

	// Create test contacts for john.doe
	contacts := []Contact{
		{
			ID:        "550e8400-e29b-41d4-a716-446655440001",
			FirstName: "Alice",
			LastName:  "Johnson",
			Email:     "alice.johnson@example.com",
			Phone:     "+1-555-0101",
			UserID:    "john.doe",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "550e8400-e29b-41d4-a716-446655440002",
			FirstName: "Bob",
			LastName:  "Williams",
			Email:     "bob.williams@example.com",
			Phone:     "+1-555-0102",
			UserID:    "john.doe",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "550e8400-e29b-41d4-a716-446655440003",
			FirstName: "Charlie",
			LastName:  "Brown",
			Email:     "charlie.brown@example.com",
			Phone:     "+1-555-0103",
			UserID:    "john.doe",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Create test contacts for jane.smith
	contactsJane := []Contact{
		{
			ID:        "550e8400-e29b-41d4-a716-446655440004",
			FirstName: "David",
			LastName:  "Miller",
			Email:     "david.miller@example.com",
			Phone:     "+1-555-0201",
			UserID:    "jane.smith",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "550e8400-e29b-41d4-a716-446655440005",
			FirstName: "Emma",
			LastName:  "Davis",
			Email:     "emma.davis@example.com",
			Phone:     "+1-555-0202",
			UserID:    "jane.smith",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	allContacts := append(contacts, contactsJane...)
	for _, contact := range allContacts {
		result := db.FirstOrCreate(&contact, Contact{ID: contact.ID})
		if result.Error != nil {
			log.Printf("Error creating contact %s: %v", contact.FirstName, result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("Created contact: %s %s", contact.FirstName, contact.LastName)
		} else {
			log.Printf("Contact already exists: %s %s", contact.FirstName, contact.LastName)
		}
	}

	// Create test addresses
	addresses := []Address{
		{
			ID:         "660e8400-e29b-41d4-a716-446655440001",
			ContactID:  "550e8400-e29b-41d4-a716-446655440001",
			Street:     "123 Main Street",
			City:       "New York",
			Province:   "NY",
			PostalCode: "10001",
			Country:    "USA",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "660e8400-e29b-41d4-a716-446655440002",
			ContactID:  "550e8400-e29b-41d4-a716-446655440001",
			Street:     "456 Oak Avenue",
			City:       "Brooklyn",
			Province:   "NY",
			PostalCode: "11201",
			Country:    "USA",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "660e8400-e29b-41d4-a716-446655440003",
			ContactID:  "550e8400-e29b-41d4-a716-446655440002",
			Street:     "789 Pine Road",
			City:       "Boston",
			Province:   "MA",
			PostalCode: "02101",
			Country:    "USA",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "660e8400-e29b-41d4-a716-446655440004",
			ContactID:  "550e8400-e29b-41d4-a716-446655440003",
			Street:     "321 Elm Street",
			City:       "Chicago",
			Province:   "IL",
			PostalCode: "60601",
			Country:    "USA",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "660e8400-e29b-41d4-a716-446655440005",
			ContactID:  "550e8400-e29b-41d4-a716-446655440004",
			Street:     "555 Beach Boulevard",
			City:       "Miami",
			Province:   "FL",
			PostalCode: "33101",
			Country:    "USA",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "660e8400-e29b-41d4-a716-446655440006",
			ContactID:  "550e8400-e29b-41d4-a716-446655440005",
			Street:     "777 Valley Drive",
			City:       "Los Angeles",
			Province:   "CA",
			PostalCode: "90001",
			Country:    "USA",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	for _, address := range addresses {
		result := db.FirstOrCreate(&address, Address{ID: address.ID})
		if result.Error != nil {
			log.Printf("Error creating address %s: %v", address.Street, result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("Created address: %s, %s", address.Street, address.City)
		} else {
			log.Printf("Address already exists: %s, %s", address.Street, address.City)
		}
	}

	fmt.Println("\n=== Seed Data Summary ===")
	fmt.Println("Users created: 3")
	fmt.Println("  - john.doe (password: password123)")
	fmt.Println("  - jane.smith (password: SecurePass456!)")
	fmt.Println("  - bob.wilson (password: password123)")
	fmt.Println("\nContacts created: 5")
	fmt.Println("  - john.doe has 3 contacts")
	fmt.Println("  - jane.smith has 2 contacts")
	fmt.Println("\nAddresses created: 6")
	fmt.Println("  - Multiple addresses distributed across contacts")
	fmt.Println("\n=== Ready to test with Postman! ===")
}
