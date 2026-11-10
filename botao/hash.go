package main

// Tabela é uma tabela de dispersão (hash) com endereçamento aberto e
// sondagem linear, como na Aula 10 e na Atividade Avaliativa 3.
// Guarda pares chave -> valor: a chave é o id do chamado.
type Tabela struct {
	chaves  []int
	valores []int
	ocupada []bool
	n       int // quantos pares estão guardados
}

// NovaTabela cria uma tabela vazia com m posições (use um primo)
func NovaTabela(m int) *Tabela {
	return &Tabela{
		chaves:  make([]int, m),
		valores: make([]int, m),
		ocupada: make([]bool, m),
	}
}

// m devolve o tamanho do vetor
func (t *Tabela) m() int {
	return len(t.chaves)
}

// hash devolve a posição de origem da chave: método da divisão, k mod m
func (t *Tabela) hash(k int) int {
	return k % t.m()
}

// Busca devolve o valor da chave e true, ou 0 e false se ela não está na tabela
func (t *Tabela) Busca(k int) (int, bool) {
	p := t.hash(k)
	for t.ocupada[p] {
		if t.chaves[p] == k {
			return t.valores[p], true
		}
		p = (p + 1) % t.m() // sondagem linear: tenta a próxima, dando a volta no fim
	}
	return 0, false
}

// Insere guarda o par (k, v). Se a chave já existe, só troca o valor
func (t *Tabela) Insere(k, v int) {
	p := t.hash(k)
	for t.ocupada[p] {
		if t.chaves[p] == k {
			t.valores[p] = v
			return
		}
		p = (p + 1) % t.m()
	}
	t.chaves[p] = k
	t.valores[p] = v
	t.ocupada[p] = true
	t.n++

	if 10*t.n > 7*t.m() { // fator de carga passou de 0,7
		t.redimensiona(proximoPrimo(2 * t.m()))
	}
}

// redimensiona troca o vetor por um de m posições e reinsere todos os pares
func (t *Tabela) redimensiona(m int) {
	nova := NovaTabela(m)
	for p := 0; p < t.m(); p++ {
		if t.ocupada[p] {
			nova.Insere(t.chaves[p], t.valores[p])
		}
	}
	*t = *nova
}

// proximoPrimo devolve o menor primo maior ou igual a k
func proximoPrimo(k int) int {
	for {
		primo := k >= 2
		for d := 2; d*d <= k; d++ {
			if k%d == 0 {
				primo = false
				break
			}
		}
		if primo {
			return k
		}
		k++
	}
}
