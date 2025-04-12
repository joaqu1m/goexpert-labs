package main

import "fmt"

type Endereco struct {
	Logradouro string
	Bairro     string
	Cidade     string
}

type Instituicao struct {
	Nome  string
	Email string
	CNAE  string
	Endereco
}

// sim, eh possivel encurtar a definição de um atributo caso o tipo tenha o exato mesmo nome que a chave. igualzinho no javascript S2

// isso é um receiver (ou receptor), ele é a maneira que o go encontrou de criar métodos em structs

// Esta função não altera o valor original pois não usa ponteiro
func (i Instituicao) SetNomeNonPointer(nome string) {
	i.Nome = nome
	println(i.Nome)
}

// por padrão, o nome não será alterado pois o receiver não foi definido como um ponteiro da struct, então é como se ele estivesse alterando uma nova cópia da struct

// Esta função altera o valor original pois usa ponteiro
func (i *Instituicao) SetNome(nome string) {
	i.Nome = nome
	println(i.Nome)
}

func (i *Instituicao) SetCNAE(cnae string) {
	i.CNAE = cnae
	println(i.CNAE)
}

type Pessoa struct {
	Nome  string
	Idade int
	Endereco
}

func (i *Pessoa) SetNome(nome string) {
	i.Nome = nome
	println(i.Nome)
}

type Nomeavel interface {
	SetNome(string)
}

func usarNovoNome(f Nomeavel) {
	f.SetNome("Novo Nome")
}

func usarNovoNomeFuncao(f func(string)) {
	f("Novo Nome")
}

func main() {
	// Exemplo com receiver sem ponteiro (não muda o original)
	itau1 := Instituicao{
		Nome:     "Itau",
		Email:    "email@example.com",
		CNAE:     "12345678",
		Endereco: Endereco{},
	}

	itau1.SetNomeNonPointer("Novo Itau")
	fmt.Println("Nome após SetNomeNonPointer:", itau1.Nome) // Não deve mudar

	// Exemplo com receiver com ponteiro (muda o original)
	itau2 := Instituicao{
		Nome:     "Itau",
		Email:    "email@example.com",
		CNAE:     "12345678",
		Endereco: Endereco{},
	}
	itau2.SetNome("Novo Itau")
	fmt.Println("Nome após SetNome:", itau2.Nome) // Deve mudar

	// passando métodos por parâmetro:
	itau3 := Instituicao{
		Nome:     "Itau Original",
		Email:    "email@example.com",
		CNAE:     "12345678",
		Endereco: Endereco{},
	}

	usarNovoNomeFuncao(itau3.SetNome)
	fmt.Println("Nome após usarNovoNomeFuncao:", itau3.Nome)

	// interfaces em go são automaticamente implementadas, ou seja, não é necessário declarar que uma struct implementa uma interface,
	// basta que ela tenha os métodos definidos na interface (no estilo left join, nem todos os métodos da struct precisam estar na interface,
	// mas todos os métodos da interface precisam estar na struct)
	itau4 := Instituicao{
		Nome:     "Itau Interface",
		Email:    "email@example.com",
		CNAE:     "12345678",
		Endereco: Endereco{},
	}

	usarNovoNome(&itau4)
	fmt.Println("Nome após usarNovoNome com interface:", itau4.Nome)

	// Demonstrando a interface com outro tipo que também implementa SetNome
	pessoa := Pessoa{
		Nome:  "João",
		Idade: 30,
	}

	usarNovoNome(&pessoa)
	fmt.Println("Nome da pessoa após usarNovoNome:", pessoa.Nome)
}
