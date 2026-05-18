package handlers

import (
    "servidorHTTP/app/utils"
    "net/http"
    "strconv"
    "fmt"
)

func DeletarPacienteHandler(response http.ResponseWriter, request *http.Request) {
    if request.Method != http.MethodGet {
        http.Error(response, "Método não suportado", http.StatusMethodNotAllowed)
        return
    }

    idStr := request.URL.Query().Get("id")
    if idStr == "" {
        http.Error(response, "ID do paciente não informado", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(response, "ID inválido", http.StatusBadRequest)
        return
    }

    err = utils.DeletarPaciente(id)
    if err != nil {
        http.Error(response, fmt.Sprintf("Erro ao deletar paciente: %v", err), http.StatusInternalServerError)
        return
    }

    http.Redirect(response, request, "/listar-pacientes", http.StatusSeeOther)
