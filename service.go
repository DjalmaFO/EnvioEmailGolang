package main

import (
	"log"
)

func pegarRelatorio(site, empresa, tabela string) (string, error) {
	log.Println("Inicio da consulta")

	sql := `
	SELECT 
		NVL(A.COD_PROD, 0) CODIGO,
		NVL(A.DESCR, ' ') DESCRICAO,
		NVL(A.T_PROD, ' ') TIPO,
		NVL(A.QTD_LIVRE, 0) LIVRE,
		NVL(A.QTD_RESERV, 0) RESERVAD0,
		NVL(A.QTD_TOTAL, 0) TOTAL
	FROM 
		` + tabela + ` A 
	WHERE 
		A.SITE = :0
		AND A.CLIENTE = :1`

	rows, err := db.Query(sql, site, empresa)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	log.Println("Montando relatório")
	// err = sqltocsv.WriteFile("relatorio.csv", rows)
	// if err != nil {
	// 	panic(err)
	// }

	exc := "relatorio.xlsx"
	err = generateXLSXFromRows(rows, exc)
	if err != nil {
		log.Fatal(err)
	}

	frase := `Segue anexo, relatorio.csv `

	err = enviarEmailAnexo(frase)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return "Email enviado com sucesso", nil
}
