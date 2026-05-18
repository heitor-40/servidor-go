package handlers

import (
    "servidorHTTP/app/utils"
    "net/http"
    "fmt"
)

func ListPacientesHandler(response http.ResponseWriter, request *http.Request) {
    // Use utils.DB direto, não ConnectToDB()
    rows, err := utils.DB.Query("SELECT id, nome, idade, data_nascimento, tipo_sanguineo FROM pacientes ORDER BY id")
    if err != nil {
        http.Error(response, "Erro ao buscar pacientes", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var html string
    for rows.Next() {
        var id, idade int
        var nome, dataNasc, tipoSang string
        
        err := rows.Scan(&id, &nome, &idade, &dataNasc, &tipoSang)
        if err != nil {
            continue
        }
        
        html += fmt.Sprintf(`
            <tr>
                <td>%d</td>
                <td>%s</td>
                <td>%d</td>
                <td>%s</td>
                <td>%s</td>
                <td>
                    <a href="/atualizar-paciente.html?id=%d">Editar</a>
                    <a href="/deletar-paciente?id=%d" onclick="return confirm('Tem certeza?')">Excluir</a>
                </td>
            </tr>`, id, nome, idade, dataNasc, tipoSang, id, id)
    }

    response.Header().Set("Content-Type", "text/html; charset=utf-8")
    response.Write([]byte(`
        <!DOCTYPE html>
        <html>
        <head>
            <title>Pacientes</title>
            <style>
                table { border-collapse: collapse; width: 100%; }
                th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
                th { background-color: #4CAF50; color: white; }
            </style>
        </head>
        <body>
            <h1>🏥 Pacientes Cadastrados</h1>
            <table>
                <tr>
                    <th>ID</th>
                    <th>Nome</th>
                    <th>Idade</th>
                    <th>Nascimento</th>
                    <th>Tipo Sanguíneo</th>
                    <th>Ações</th>
                </tr>
                ` + html + `
            </table>
            <br>
            <a href="/">← Voltar</a>
            <a href="/criar-paciente.html">+ Novo</a>
        </body>
        </html>
    `))
}
