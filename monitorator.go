package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)


const monitoramento = 3
const delay = 3


func main() {
	for {
		exibeMenu()

		comando := pegaInputDoUsuario()

		switch comando {
		case 1:
			fmt.Println("Iniciando monitoramento dos sites...")
			logs := monitorandoSites()
			registraLog(logs)
		case 2:
			fmt.Println("Exibindo logs do sistema...")
			imprimeLogs()
		case 0:
			fmt.Println("Saindo...")
			os.Exit(0)
		default:
			fmt.Println("Não conheço esse comando.")
			os.Exit(-1)
		}
	}

}


func pegaInputDoUsuario() int {
	var comando int
	fmt.Scan(&comando)
	endOfLine()
	return comando
}


func exibeMenu() {
	fmt.Println("1 - Iniciar monitoramento.")
	fmt.Println("2 - Exibir logs.")
	fmt.Println("0 - Sair.")
	endOfLine()
}


func endOfLine() {
	fmt.Println("-----------")
}


func monitorandoSites() []string {
	sites := leSitesDoArquivo()
	var logsArr []string

	for i := 0; i < monitoramento; i++ {
		for index, site := range sites {
			fmt.Println("Testando site:", index+1, ":", site)
			arr := testaSite(site)
			
			for _, v := range arr {
				logsArr = append(logsArr, v)
			}
			
		}
		time.Sleep(delay * time.Second)
		endOfLine()
	}

	return logsArr
}


func testaSite(site string) []string {
	resp, _ := http.Get(site)
	var listOk []string
	var listErr []string

	if resp.StatusCode == 200 {
		statusCodeOk := strconv.Itoa(resp.StatusCode)
		listOk = append(listOk, site, statusCodeOk)
		return listOk
	} else {
		statusCodeErr := strconv.Itoa(resp.StatusCode)
		listErr = append(listErr, site, statusCodeErr)
		return listErr
	}
}


func leSitesDoArquivo() []string {
	var slices []string
	file, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println(err)
	} else {
		sites := bufio.NewReader(file)

		for {
			lines, err := sites.ReadString('\n')
			lines = strings.TrimSpace(lines)

			if err == io.EOF {
				break
			}

			slices = append(slices, lines)
		}

	}

	file.Close()

	return slices
}


func registraLog(sites []string) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) // 0666 refere-se ao tipo de permissão.

	if err != nil {
		fmt.Println(err)
		fmt.Println("Ocorreu um erro inesperado ao tentar criar ou ler/escrever no arquivo de log.")
		return
	}

	for i := 0; i < len(sites); i++ {
		file.WriteString(time.Now().Format("02/01/2006 15:04:05")+ " " + sites[i] + " " + "\n")
		continue
	}
	
	file.Close()
	return
}


func imprimeLogs() {
	file, err := os.Open("log.txt")

	if err != nil {
		fmt.Println(err)
	} else {
		logs := bufio.NewReader(file)

		for {
			if err != nil {
				fmt.Println(err)
			} else {
				lines, err := logs.ReadString('\n')
				lines = strings.TrimSpace(lines)

				if err == io.EOF {
					break
				}

				fmt.Println(lines)
			}
		}
	}
	
	endOfLine()
}

