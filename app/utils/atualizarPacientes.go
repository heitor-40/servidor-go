package utils

import "fmt"

func AtualizarPaciente(id int, nome string, idade int, dataNascimento string, tipoSanguineo string) error {
    _, err := DB.Exec(
        "UPDATE pacientes SET nome=$1, idade=$2, data_nascimento=$3, tipo_sanguineo=$4 WHERE id=$5",
        nome, idade, dataNascimento, tipoSanguineo, id,
    )
    if err != nil {
        return fmt.Errorf("erro ao atualizar paciente: %v", err)
    }
    return nil
}
