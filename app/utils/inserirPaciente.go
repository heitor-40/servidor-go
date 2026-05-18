package utils

import "fmt"

func InserirPaciente(nome string, idade int, dataNascimento string, tipoSanguineo string) error {
    _, err := DB.Exec(
        "INSERT INTO pacientes (nome, idade, data_nascimento, tipo_sanguineo) VALUES ($1, $2, $3, $4)",
        nome, idade, dataNascimento, tipoSanguineo,
    )
    if err != nil {
        return fmt.Errorf("erro ao inserir paciente: %v", err)
    }
    return nil
}
