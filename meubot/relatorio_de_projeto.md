# Relatório de projeto: Arthur Rangel Silveira Golfette Lopes e Rodrigo Yuji Nonaka Rodrigues

## Achar os chamados próximos de uma ambulância
Estrutura: A estrutura adotada foi a grade de baldes, que consiste em um mapa dividido em quadrados, cada quadrado guarda os chamados existentes nele. A ambulância olha o quadrado onde está e depois os vizinhos, cada vez mais longe, e para quando os chamados mais distantes não podem ter nota melhor que a já achada.
Custo por turno: Só olha os quadrados perto da ambulância, então quase não depende do total de chamados.
Alternativa descartada: Foi descartada a opção de ler todos os chamados, visto que seria lento e demorado.

## Descartar os chamados que expiraram
Estrutura: Foi adotado um vetor com uma lista de identificadores para cada turno. Quando um chamado chega, o identificador dele vai na lista do turno do prazo, dessa forma, ao passar vai fechando os chamados da lista.
Custo por turno: Só os chamados que expiram naquele turno.
Alternativa descartada: Olhar todos os chamados a cada turno procurando os vencidos, seria lento.

## Remover um chamado recolhido pela outra equipe, sabendo só o identificador
Estrutura: É usada uma tabela hash que leva o identificador até a posição do chamado no vetor.
Custo por turno: O custo gira em torno de O(1) em média por chamado recolhido.
Alternativa descartada: Foi descartada a ideia de procurar o identificador no vetor inteiro visto que com muitos chamados o processo seria lento.

## Evitar que duas ambulâncias sigam para o mesmo chamado
Estrutura: A struct de cada chamado tem um campo que armazena a ambulância que o atende.
Custo por turno: É constante por ambulância.
Alternativa descartada: Guardar as reservas numa tabela separada, porque o campo no próprio chamado dá a mesma resposta mais rápido.

## Achar o hospital mais próximo de uma ambulância ocupada
Estrutura: Lista de hospitais, função que calcula a distância de cada um e comparar qual o mais próximo.
Custo por turno: Número de hospitais vezes o número de (ambulâncias ocupadas + chamadas do novo turno).
Alternativa descartada: Calcular antes o hospital mais próximo de cada ponto do mapa, porque com poucos hospitais a lista já é rápida e o cálculo gastaria memória.
