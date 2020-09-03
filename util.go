package main

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/tealeg/xlsx"

	logrus "github.com/sirupsen/logrus"
	gomail "gopkg.in/mail.v2"
)

func enviarEmailAnexo(frase string) error {
	log.Println("Formatando e-mail")
	usuario := os.Getenv("MAIL_USER")
	senha := os.Getenv("MAIL_PASSWORD")
	servidor := os.Getenv("MAIL_SERVER")
	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		log.Println("Não foi possivel converter port para inteiro" + err.Error())
		return err
	}
	remetente := os.Getenv("MAIL_FROM")
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", remetente)

	para := os.Getenv("GRUPO_EMAIL")
	// Set E-Mail receivers
	m.SetHeader("To", para)
	// m.SetHeader("Cc", "anotherguy@example.com")
	// m.SetHeader("Bcc", "office@example.com", "anotheroffice@example.com")

	// Set E-Mail subject
	m.SetHeader("Subject", "Estoque Kit")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", frase)

	// Attach some file
	m.Attach("relatorio.xlsx")

	// Settings for SMTP server
	d := gomail.NewDialer(servidor, port, usuario, senha)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Save E-Mail in mymail.txt file

	// Get directory where binary is started
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	// Write contents of E-Mail into mymail.txt.
	// This is useful for debuging.
	var b bytes.Buffer
	m.WriteTo(&b)
	err = ioutil.WriteFile(dir+`mymail.txt`, b.Bytes(), 0777)
	if err != nil {
		panic(err)
	}

	log.Println("Enviando email")

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil
}

func iniciarProcesso() {
	log.Println("Configurando consulta")

	site, empresa, tabela, resp := configConsulta()
	if resp != 0 {
		log.Println("Arquivo .env inexistente ou mau formado")
	}

	r, err := pegarRelatorio(site, empresa, tabela)
	if err != nil {
		logrus.Fatal(err.Error(), map[string]string{"msg": "Erro ao realizar a requisição de e-mail"})
	} else {
		log.Println(r)
	}

	log.Println("Encerrando o serviço")
}

func configConsulta() (string, string, string, int) {
	erro := 0
	site, ok := os.LookupEnv("SITE_CONSULTA")
	if !ok {
		log.Fatal("Variável 'SITE_CONSULTA' não foi definida!")
		erro++
	}

	tabela, ok := os.LookupEnv("TABELA_CONSULTA")
	if !ok {
		log.Fatal("Variável 'TABELA_CONSULTA' não foi definida!")
		erro++
	}

	cliente, ok := os.LookupEnv("CLIENTE_CONSULTA")
	if !ok {
		log.Fatal("Variável 'CLIENTE_CONSULTA' não foi definida!")
		erro++
	}

	return site, cliente, tabela, erro
}

func generateXLSXFromRows(rows *sql.Rows, outf string) error {

	var err error

	// Get column names from query result
	colNames, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("error fetching column names, %s\n", err)
	}
	length := len(colNames)

	// Create a interface slice filled with pointers to interface{}'s
	pointers := make([]interface{}, length)
	container := make([]interface{}, length)
	for i := range pointers {
		pointers[i] = &container[i]
	}

	// Create output xlsx workbook
	xfile := xlsx.NewFile()
	xsheet, err := xfile.AddSheet("Sheet1")
	if err != nil {
		return fmt.Errorf("error adding sheet to xlsx file, %s\n", err)
	}

	// Write Headers to 1st row
	xrow := xsheet.AddRow()
	xrow.WriteSlice(&colNames, -1)

	// Process sql rows
	for rows.Next() {

		// Scan the sql rows into the interface{} slice
		err = rows.Scan(pointers...)
		if err != nil {
			return fmt.Errorf("error scanning sql row, %s\n", err)
		}

		xrow = xsheet.AddRow()

		// Here we range over our container and look at each column
		// and set some different options depending on the column type.
		for _, v := range container {
			xcell := xrow.AddCell()
			switch v := v.(type) {
			case string:
				xcell.SetString(v)
			case []byte:
				xcell.SetString(string(v))
			case int64:
				xcell.SetInt64(v)
			case float64:
				xcell.SetFloat(v)
			case bool:
				xcell.SetBool(v)
			case time.Time:
				xcell.SetDateTime(v)
			default:
				xcell.SetValue(v)
			}

		}

	}

	// Save the excel file to the provided output file
	err = xfile.Save(outf)
	if err != nil {
		return fmt.Errorf("error writing to output file %s, %s\n", outf, err)
	}

	return nil
}
