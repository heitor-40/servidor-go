package utils

import "fmt"

func DeletarPaciente(id int) error {
    _, err := DB.Exec("DELETE FROM pacientes WHERE id = $1", id)
    if err != nil {
        return fmt.Errorf("erro ao deletar paciente: %v", err)
    }
    return nil
}
