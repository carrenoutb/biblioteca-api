package db

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

// Connect abre la conexión a MySQL, reintentando unas cuantas veces
// mientras el contenedor de la base de datos termina de arrancar.
func Connect(dsn string) *sqlx.DB {
	var conn *sqlx.DB
	var err error

	maxRetries := 10
	for i := 1; i <= maxRetries; i++ {
		// sqlx.Connect ya hace Open + Ping, así que err == nil implica conexión verificada.
		conn, err = sqlx.Connect("mysql", dsn)
		if err == nil {
			break
		}

		log.Printf("Intento %d/%d: no se pudo conectar a MySQL (%v)", i, maxRetries, err)
		if i < maxRetries {
			time.Sleep(3 * time.Second)
		}
	}

	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos después de %d intentos: %v", maxRetries, err)
	}

	conn.SetMaxOpenConns(20)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(3 * time.Minute)

	log.Println("Conexión a MySQL establecida correctamente")
	return conn
}
