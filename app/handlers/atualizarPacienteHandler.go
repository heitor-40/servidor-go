package handlers

import (
    "servidorHTTP/app/utils"
    "net/http"
    "strconv"
    "fmt"
)

func AtualizarPacienteHandler(response http.ResponseWriter, request *http.Request) {
    if request.Method != http.MethodPost {
        http.Error(response, "Método não suportado", http.StatusMethodNotAllowed)
        return
    }

    idStr := request.FormValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(response, "ID inválido", http.StatusBadRequest)
        return
    }

    nome := request.FormValue("nome")
    idade := request.FormValue("idade")
    dataNascimento := request.FormValue("data_nascimento")
    tipoSanguineo := request.FormValue("tipo_sanguineo")

    if nome == "" || tipoSanguineo == "" {
        http.Error(response, "Nome e tipo sanguíneo são obrigatórios", http.StatusBadRequest)
        return
    }

    idadeInt, err := strconv.Atoi(idade)
    if err != nil {
        http.Error(response, "Idade deve ser um número válido", http.StatusBadRequest)
        return
    }

    err = utils.AtualizarPaciente(id, nome, idadeInt, dataNascimento, tipoSanguineo)
    if err != nil {
        http.Error(response, fmt.Sprintf("Erro ao atualizar paciente: %v", err), http.StatusInternalServerError)
        return
    }

    http.Redirect(response, request, "/listar-pacientes", http.StatusSeeOther)
}
