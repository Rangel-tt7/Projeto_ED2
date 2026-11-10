package main

import (
	"math/rand"
	"testing"
)

// confere verifica se a tabela tem exatamente as chaves esperadas, nas posições esperadas
func confere(t *testing.T, tab *TabelaHash, esperado map[int]int) {
	t.Helper()
	for p := 0; p < tab.tamanho(); p++ {
		k, deveria := esperado[p]
		if deveria && (!tab.ocupada[p] || tab.chaves[p] != k) {
			t.Errorf("posição %d: esperava a chave %d", p, k)
		}
		if !deveria && tab.ocupada[p] {
			t.Errorf("posição %d: esperava vazia, tem a chave %d", p, tab.chaves[p])
		}
	}
}

// Exemplo da Aula 10 / AA3: m = 11 e as chaves 22, 33, 5, 16, 27, 38 e 44.
func TestInsereBusca(t *testing.T) {
	tab := CriarTabela(11)
	for _, k := range []int{22, 33, 5, 16, 27, 38, 44} {
		tab.Inserir(k, k*10)
	}
	confere(t, tab, map[int]int{0: 22, 1: 33, 2: 44, 5: 5, 6: 16, 7: 27, 8: 38})

	for _, k := range []int{22, 33, 5, 16, 27, 38, 44} {
		v, ok := tab.Buscar(k)
		if !ok || v != k*10 {
			t.Errorf("Busca(%d) = %d, %v; esperava %d, true", k, v, ok, k*10)
		}
	}
	if _, ok := tab.Buscar(49); ok {
		t.Errorf("Busca(49) achou uma chave que não existe")
	}
}

// Primeiro exemplo da Remove na AA3: o 1 fica, o 10 desce.
func TestRemoveVoltaCircular(t *testing.T) {
	tab := CriarTabela(11)
	for _, k := range []int{9, 20, 31, 1, 10} {
		tab.Inserir(k, 0)
	}
	confere(t, tab, map[int]int{0: 31, 1: 1, 2: 10, 9: 9, 10: 20})

	if !tab.Remover(20) {
		t.Fatal("Remove(20) devolveu false")
	}
	confere(t, tab, map[int]int{0: 10, 1: 1, 9: 9, 10: 31})
}

// Segundo exemplo da AA3: o 21 e o 31 descem dando a volta no fim do vetor, o 2 fica.
func TestRemoveDistancia(t *testing.T) {
	tab := CriarTabela(11)
	for _, k := range []int{9, 20, 21, 31, 2} {
		tab.Inserir(k, 0)
	}
	confere(t, tab, map[int]int{0: 21, 1: 31, 2: 2, 9: 9, 10: 20})

	tab.Remover(20)
	confere(t, tab, map[int]int{0: 31, 2: 2, 9: 9, 10: 21})

	if tab.Remover(20) {
		t.Error("Remove(20) de novo devolveu true")
	}
}

// Oráculo cruzado: faz as mesmas operações aleatórias na Tabela e num map
// (permitido em teste) e confere se as duas concordam.
// Detecta remoção que deixa chave inalcançável ou conta n errado.
func TestOraculo(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	tab := CriarTabela(11)
	ref := map[int]int{}
	for i := 0; i < 100000; i++ {
		k := 1 + r.Intn(3000)
		if r.Intn(3) < 2 {
			tab.Inserir(k, i)
			ref[k] = i
		} else {
			_, estava := ref[k]
			if tab.Remover(k) != estava {
				t.Fatalf("Remove(%d) discordou do map", k)
			}
			delete(ref, k)
		}
	}
	if tab.quantidade != len(ref) {
		t.Fatalf("n = %d, o map tem %d", tab.quantidade, len(ref))
	}
	for k, v := range ref {
		if got, ok := tab.Buscar(k); !ok || got != v {
			t.Fatalf("Busca(%d) = %d, %v; esperava %d", k, got, ok, v)
		}
	}
}
