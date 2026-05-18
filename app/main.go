package main

import (
    "fmt"
    "servidorHTTP/app/handlers"
    "servidorHTTP/app/utils"
    "log"
    "net"
    "net/http"
)

func main() {
    
    utils.ConnectToDB()

    
    fileserver := http.FileServer(http.Dir("./static"))
    http.Handle("/", fileserver)

    
    http.HandleFunc("/criar-paciente", handlers.PacienteHandler)
    http.HandleFunc("/listar-pacientes", handlers.ListPacientesHandler)
    http.HandleFunc("/atualizar-paciente", handlers.AtualizarPacienteHandler)
    http.HandleFunc("/deletar-paciente", handlers.DeletarPacienteHandler)

    // Obter IP local
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        log.Fatal(err)
    }

    var localIP string
    for _, addr := range addrs {
        if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
            localIP = ipNet.IP.String()
            break
        }
    }

    port := "3000"
    if localIP == "" {
        localIP = "127.0.0.1"
    }

    fmt.Printf(" Servidor de Saúde rodando em: http://%s:%s/\n", localIP, port)

    if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
        log.Fatal(err)
    }
}
