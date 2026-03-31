package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Reset  = "\033[0m"
)

type Question struct {
	Text    string
	Options []string
	Answer  int
}

type GameState struct {
	Name      string
	Points    int
	Questions []Question
}

func (g *GameState) Init() {
	fmt.Printf("Seja bem vindo(a) ao quiz %sEscreva seu nome:%s\n", Yellow, Reset)

	reader := bufio.NewReader(os.Stdin)

	name, err := reader.ReadString('\n')
	if err != nil {
		panic("Erro ao ler a string")
	}

	g.Name = name

	fmt.Printf("Vamos ao jogo %s\n", g.Name)
}

func (g *GameState) ProcessCSV() {
	f, err := os.Open("quiz-go.csv")
	if err != nil {
		panic("Erro ao ler o arquivo")
	}

	// "defer" em Go é usado para adiar a execução de uma função até o fim da função atual é registrado no ponto onde aparece, mas só executa quando a função termina
	// "f.Close()" fecha o arquivo que foi aberto com "os.Open()"
	defer f.Close()

	reader := csv.NewReader(f)

	records, err := reader.ReadAll()
	if err != nil {
		panic("Erro ao ler o csv")
	}

	for index, record := range records {
		if index > 0 {
			correctAnswer, _ := toInt(record[5])
			question := Question{
				Text:    record[0],
				Options: record[1:5],
				Answer:  correctAnswer,
			}

			g.Questions = append(g.Questions, question)
		}
	}
}

func (g *GameState) Run() {
	// Exibir pergunta para o usuário
	for i, question := range g.Questions {
		fmt.Printf("%s %d. %s %s\n", Yellow, i+1, question.Text, Reset)

		// Vamos iterar sobre as opções que temos no game state
		// e exibir no terminal para o usuário
		for j, option := range question.Options {
			fmt.Printf("[%d] %s\n", j+1, option)
		}

		fmt.Println("Digite uma alternativa:")

		// Vamos coletar a entrada do usuário
		// Validar o caractere que foi inserido
		// Se for errado o usuário precisa inserir novamente
		var answer int
		var err error

		// loop infinito em Go (while), "for" sem expressão de parada
		for {
			reader := bufio.NewReader(os.Stdin)
			read, _ := reader.ReadString('\n')

			// "read[:len(read)-1]" operação de splice em um slice de string
			// remove o "\n" da string de leitura do terminal
			answer, err = toInt(read[:len(read)-1])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			break
		}

		// Validar a resposta
		// Exibir mensagem se correta ou não
		// Calcular a pontuação
		if answer == question.Answer {
			len, _ := fmt.Printf("%sParabéns você acertou!!!%s", Blue, Reset)
			fmt.Println()
			fmt.Println(strings.Repeat("-", len-1))
			fmt.Println()
			g.Points += 10
		} else {
			len, _ := fmt.Printf("%sOps! Errou!%s", Red, Reset)
			fmt.Println()
			fmt.Println(strings.Repeat("-", len-1))
			fmt.Println()
		}
	}
}

func main() {
	game := &GameState{}
	// Chama uma Goroutine para ler o arquivo enquanto o jogo processa outras partes do programa, deixando os valores disponíveis para leitura antecipadamente, acelerando o processamento geral do jogo
	go game.ProcessCSV()
	game.Init()
	game.Run()

	len, _ := fmt.Printf("%sFim de jogo, você fez %d pontos\n%s", Green, game.Points, Reset)
	fmt.Println(strings.Repeat("-", len-1))
}

func toInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("Não é permitido caractere diferente de número, por favor insira um número")
	}

	return i, nil
}
