@startuml
title Sistema de Gestão de Apólices - Modelo C4 (versão compatível)

' ===== Diagrama de Contexto =====
actor "Usuário" as Usuario
actor "Administrador" as Administrador

package "Sistema de Gestão de Apólices" {
	rectangle "Frontend (React + Vite)" as Frontend_Context
	rectangle "Backend API (Go)" as API_Context
	database "PostgreSQL" as DB_Context
	folder "Armazenamento de arquivos" as Storage_Context
}

Usuario --> Frontend_Context : Usa / interage com
Frontend_Context --> API_Context : HTTP/JSON
Administrador --> API_Context : Administra via API
API_Context --> DB_Context : Lê/Grava dados (SQL)
API_Context --> Storage_Context : Armazena/Recupera arquivos

' ===== Diagrama de Contêineres =====
package "Contêineres" {
	rectangle "Frontend (React + Vite)\n(servido por Nginx)" as Frontend
	rectangle "Backend API (Go)\n(módulos: apólice, auth, audit, notificacao, uploads)" as API
	database "PostgreSQL" as DB
	folder "Armazenamento de arquivos" as Storage
	rectangle "Gerador de PDFs (Python)\n(generate_pdfs.py)" as PDF
}

Frontend --> API : Chama endpoints REST
API --> DB : Persiste e consulta (SQL)
API --> Storage : Grava/recupera arquivos
API --> PDF : Solicita geração de PDF

' ===== Diagrama de Componentes (Backend) =====
package "Componentes (Backend API)" {
	rectangle "Handlers / Routes\n(endpoints: /apolices, /auth, /audit, /notificacao, /uploads)" as Handlers
	rectangle "Serviços de Negócio" as Services
	rectangle "Repositório (Postgres)" as Repository
	rectangle "Auth / Middleware" as Auth
	rectangle "Audit Service" as Audit
	rectangle "Notificação Service" as Notificacao
	rectangle "Upload Handler" as Upload
}

Frontend --> Handlers : HTTP/JSON
Handlers --> Services : chama
Services --> Repository : lê/grava dados
Services --> Notificacao : solicita envio
Services --> Audit : registra ações
Handlers --> Auth : verifica autenticação
Upload --> Storage : salva arquivos

@enduml
