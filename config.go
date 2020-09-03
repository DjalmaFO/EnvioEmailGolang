package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func configurar() {

	err := godotenv.Load(envFilename)
	if err != nil {
		log.Fatal(err)
	}

	// Serviço
	srvaddr, ok := os.LookupEnv("ESTOQUE_SKF_ADDRESS")
	if !ok {
		log.Fatal("Variável 'ESTOQUE_SKF_ADDRESS' não foi definida!")
	}
	srvport, ok := os.LookupEnv("ESTOQUE_SKF_PORT")
	if !ok {
		log.Fatal("Variável 'ESTOQUE_SKF_PORT' não foi definida!")
	}
	serverAddress = srvaddr + ":" + srvport

	// DB MOracle
	dbuser, ok := os.LookupEnv("ORACLE_DB_USER")
	if !ok {
		log.Fatal("Variável 'ORACLE_DB_USER' não foi definida!")
	}
	dbpwd, ok := os.LookupEnv("ORACLE_DB_PASSWD")
	if !ok {
		log.Fatal("Variável 'ORACLE_DB_PASSWD não foi definida!")
	}
	dbsrv, ok := os.LookupEnv("ORACLE_DB_SERVER")
	if !ok {
		log.Fatal("Variável 'ORACLE_DB_SERVER' não foi definida!")
	}
	dbport, ok := os.LookupEnv("ORACLE_DB_PORT")
	if !ok {
		log.Fatal("Variável 'ORACLE_DB_PORT' não foi definida!")
	}
	dbname, ok := os.LookupEnv("ORACLE_DB_SID")
	if !ok {
		log.Fatal("Variável 'ORACLE_DB_SID' não foi definida!")
	}

	dsn := dbuser + "/" + dbpwd + "@" + dbsrv + ":" + dbport + "/" + dbname

	log.Println("Tentando estabelecer conexão com o SERVIÇO")
	//log.Println(dsn)
	db, err = sql.Open(driverOracle, dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Banco de dados conectado!")

}
