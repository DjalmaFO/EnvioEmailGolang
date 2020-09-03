package main

import (
	"database/sql"
	"log"

	_ "github.com/godror/godror"
	"github.com/labstack/echo"
)

const envFilename = ".env"

const driverOracle = "godror"

var db *sql.DB
var e *echo.Echo

var serverAddress string

func main() {
	configurar()

	defer db.Close()

	iniciarProcesso()

	log.Println("Serviço encerrado")
}
