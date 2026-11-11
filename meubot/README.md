# Projeto_ED2 - ETAPA 1

**Autores: Arthur Rangel Silveira Golfette Lopes (GRR20251916) e Rodrigo Yuji Nonaka Rodrigues (GRR20251178)**\
DS143 Estruturas de Dados II, UFPR, 2026/2, ENTREGA DA ETAPA 1 

## Como compilar e rodar
### Linux:

    go build -o meubot .
    ../bin/servidor -cenario medio -etapa 1 -semente 1 ./meubot ../bin/indexado

### Windows: 

    go build -o meubot.exe .
    ..\bin\servidor.exe -cenario medio -etapa 1 -semente 1 .\meubot.exe ..\bin\indexado.exe

## Organização dos arquivos

- `main.go`: variáveis globais do bot e a função `main`, que lê o início da partida.
- `protocolo.go`: leitura das mensagens do servidor, com leitura de palavras e números.
- `chamados.go`: estrutura do chamado e a função `fechar`, marca um chamado como fechado e remove da tabela hash.
- `grade.go`: grade de baldes que divide o mapa em quadrados, uma para os chamados comuns e outra para os urgentes, e a busca do melhor chamado em anéis em volta da ambulância.
- `mapa.go`: posição dos hospitais, distância de Manhattan e busca do hospital mais próximo de um ponto.
- `ambulancias.go`: estrutura da ambulância, pontos por turno de viagem, a escolha da direção do movimento e a jogada de cada turno, com a reserva do chamado-alvo.
- `hash.go`: tabela hash própria, guarda o identificador do chamado e a posição dele no vetor.