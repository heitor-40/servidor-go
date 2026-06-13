# C4 Model â€” Sistema de Seguros (Grupo-4-Projeto-Integrador)

> Diagramas em sintaxe PlantUML / C4-PlantUML.
> Para renderizar: cole cada bloco em https://www.plantuml.com/plantuml/uml/ ou use a extensÃ£o PlantUML no VS Code.

---

## NÃ­vel 1: Contexto

```plantuml
@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Context.puml

title Sistema de Seguros â€” CRM Shopping (NÃ­vel 1: Contexto)

Person(admin, "Administrador", "Gerencia usuÃ¡rios, acessa audit log completo e tem permissÃ£o total sobre apÃ³lices.")
Person(gestor, "Gestor", "Cria, edita e renova apÃ³lices. Faz upload de documentos e atualiza responsÃ¡veis.")
Person(visualizador, "Visualizador / Lojista", "Consulta apÃ³lices, dashboards e notificaÃ§Ãµes. Apenas leitura.")

System(seguros, "Sistema de Seguros", "CRM para gestÃ£o de apÃ³lices de seguros do shopping. Permite gerenciar apÃ³lices, documentos, notificaÃ§Ãµes, audit log e KPIs.")

System_Ext(sinistros, "Sistema de Sinistros", "Plataforma externa para abertura, acompanhamento e regulaÃ§Ã£o de sinistros vinculados Ã s apÃ³lices.")
System_Ext(seguradora, "Seguradora Parceira", "Sistema da seguradora responsÃ¡vel pela emissÃ£o, validaÃ§Ã£o e cobertura das apÃ³lices.")
SystemDb_Ext(bancoExterno, "Banco de Dados Externo", "Base de dados de terceiros (ex: ERP do shopping) consultada para validaÃ§Ã£o de contratos de locaÃ§Ã£o.")

Rel(admin, seguros, "Usa", "HTTPS")
Rel(gestor, seguros, "Usa", "HTTPS")
Rel(visualizador, seguros, "Usa", "HTTPS")

Rel(seguros, sinistros, "Consulta e registra sinistros", "JSON/HTTPS")
Rel(seguros, seguradora, "Valida apÃ³lices e coberturas", "JSON/HTTPS")
Rel(seguros, bancoExterno, "Consulta dados de contratos", "JDBC/HTTPS")

SHOW_LEGEND()
@enduml
```

---

## Ni­vel 2: Container

```plantuml
@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Container.puml

title Sistema de Seguros â€” CRM Shopping (NÃ­vel 2: ContÃªiner)

Person(usuario, "UsuÃ¡rio autenticado", "Admin, Gestor ou Visualizador")

System_Boundary(sistema, "Sistema de Seguros") {
    Container(frontend, "AplicaÃ§Ã£o Web (SPA)", "React 19, Vite, TypeScript, TailwindCSS", "Interface do usuÃ¡rio. Rotas: /dashboard, /seguros, /historico, /lojistas, /relatorios, /novo-sinistro. Consome a API via JSON/HTTPS.")
    Container(backend, "API Backend", "Go (Golang), net/http stdlib", "Servidor HTTP na porta 8080. Registra e serve todas as rotas /api/*. Valida JWT em cada request protegido. Serve o build estÃ¡tico do frontend.")
    ContainerDb(db, "Banco de Dados", "PostgreSQL 16 (flamboyant_seguros)", "Persiste: usuarios, seguros, coberturas, documentos, historico_apolice, audit_logs, notificacoes.")
}

System_Ext(sinistros, "Sistema de Sinistros", "Plataforma externa para abertura e acompanhamento de sinistros.")
System_Ext(seguradora, "Seguradora Parceira", "Sistema da seguradora â€” emissÃ£o e validaÃ§Ã£o de apÃ³lices.")
SystemDb_Ext(bancoExterno, "Banco de Dados Externo", "Base de dados de contratos de locaÃ§Ã£o (ERP do shopping).")

Rel(usuario, frontend, "Acessa via browser", "HTTPS :5173 (dev) / :8080 (prod)")
Rel(frontend, backend, "Consome API REST", "JSON/HTTPS :8080/api/*")
Rel(backend, db, "LÃª e escreve dados", "SQL/TCP :5432 â€” driver database/sql")

Rel(backend, sinistros, "Consulta e registra sinistros", "JSON/HTTPS")
Rel(backend, seguradora, "Valida apÃ³lices e coberturas", "JSON/HTTPS")
Rel(backend, bancoExterno, "Consulta dados de contratos", "JDBC/HTTPS")

SHOW_LEGEND()
@enduml
```

---

## NÃ­vel 3: Componente

```plantuml
@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml

title API Backend â€” Go (NÃ­vel 3: Componente)

Container_Boundary(backend, "API Backend â€” Go") {

    Component(main, "cmd/api/main.go", "Go â€” entrypoint", "Ponto de entrada da aplicaÃ§Ã£o. Chama app.Run().")

    Component(app, "internal/app/app.go", "Go â€” compositor", "Carrega config, abre conexÃ£o com DB, instancia e conecta todos os pacotes internos. Registra o middleware chain.")

    Component(cfg, "pkg/config", "Go â€” configuraÃ§Ã£o", "LÃª variÃ¡veis de ambiente: PORT, PG_*, JWT_SECRET, JWT_EXPIRATION_HOURS, FRONTEND_DIR. Valida JWT_SECRET (mÃ­nimo 32 chars).")

    Component(mw, "internal/middleware", "Go â€” middleware HTTP", "Chain de middlewares: RequestID (gera UUID por request), Logger (mÃ©todo, path, status, duration), Recover (panic handler), CORS (headers Access-Control-*).")

    Component(auth, "internal/auth", "Go â€” autenticaÃ§Ã£o e usuÃ¡rios", "Handler: POST /api/auth/login, GET /api/auth/me. CRUD de usuÃ¡rios (admin only): GET/POST/PATCH /api/usuarios. AuthMiddleware valida JWT Bearer. RequireRole valida role (admin/gestor/visualizador).")

    Component(apolice, "internal/apolice", "Go â€” domÃ­nio principal", "CRUD completo de apÃ³lices (/api/apolices). Endpoints de KPIs: /kpis/history, /kpis/expiring-by-week, /kpis/coverage-history, /kpis/risk-by-segment, /kpis/health-score. Fila de aÃ§Ã£o (/api/fila-de-acao). Upload/download de documentos. ExportaÃ§Ã£o CSV. Busca full-text. Map layout e lojas.")

    Component(notificacao, "internal/notificacao", "Go â€” notificaÃ§Ãµes", "GET /api/notificacoes â€” lista notificaÃ§Ãµes do usuÃ¡rio logado. PATCH /api/notificacoes/marcar-lidas. DELETE /api/notificacoes/arquivadas e /{id}. Modelo: tipo vencida ou a_vencer, com JOIN em seguros.")

    Component(audit, "internal/audit", "Go â€” auditoria", "GET /api/admin/audit â€” lista logs de auditoria. POST /api/admin/audit â€” registra aÃ§Ã£o. Recebe eventos de: auth (login), apolice (criar, editar, renovar, excluir, exportar, upload_documento). Tabela audit_logs com payload anterior e novo em JSONB.")

    Component(database, "internal/database/postgres.go", "Go â€” conexÃ£o DB", "Abre e valida conexÃ£o com PostgreSQL via database/sql. Recebe PostgresConfig (host, port, user, password, dbname, sslmode).")

    Component(response, "pkg/response", "Go â€” helpers HTTP", "Padroniza responses JSON: Success(w, status, data, requestID) e Fail(w, status, message, requestID, errors).")
}

ContainerDb_Ext(db, "PostgreSQL", "flamboyant_seguros", "Tabelas: usuarios, seguros, coberturas, documentos, historico_apolice, audit_logs, notificacoes")
System_Ext(sinistros, "Sistema de Sinistros", "Plataforma externa de sinistros")
System_Ext(seguradora, "Seguradora Parceira", "EmissÃ£o e validaÃ§Ã£o de apÃ³lices")

Rel(main, app, "Chama")
Rel(app, cfg, "Carrega configuraÃ§Ã£o")
Rel(app, database, "Inicializa conexÃ£o")
Rel(app, mw, "Aplica chain de middlewares")
Rel(app, auth, "Registra rotas e injeta auditSvc + jwtSecret")
Rel(app, apolice, "Registra rotas e injeta auditSvc + jwtSecret")
Rel(app, notificacao, "Registra rotas e injeta jwtSecret")
Rel(app, audit, "Registra rotas, retorna Service compartilhado")
Rel(auth, audit, "Injeta AuditService â€” registra login e aÃ§Ãµes de usuÃ¡rio")
Rel(apolice, audit, "Injeta AuditService â€” registra CRUD e uploads")
Rel(auth, response, "Usa para serializar responses")
Rel(apolice, response, "Usa para serializar responses")
Rel(notificacao, response, "Usa para serializar responses")
Rel(audit, response, "Usa para serializar responses")
Rel(database, db, "SQL/TCP")
Rel(apolice, database, "Queries em seguros, coberturas, documentos, historico_apolice")
Rel(auth, database, "Queries em usuarios")
Rel(notificacao, database, "Queries em notificacoes JOIN seguros")
Rel(audit, database, "Queries em audit_logs")
Rel(apolice, sinistros, "Consulta e registra sinistros", "JSON/HTTPS")
Rel(apolice, seguradora, "Valida apÃ³lices e coberturas", "JSON/HTTPS")

SHOW_LEGEND()
@enduml
```

---

## NÃ­vel 4: CÃ³digo

### Diagrama de classes â€” Pacote `internal/apolice`

```plantuml
@startuml
title internal/apolice (NÃ­vel 4: CÃ³digo â€” Diagrama de Classes)

skinparam classAttributeIconSize 0

interface Repository {
    +List() []Apolice, error
    +Get(luc string) Apolice, error
    +Create(model Apolice) Apolice, error
    +Update(luc string, model Apolice) Apolice, error
    +Delete(luc string) error
    +SearchApolices(query string) []Apolice, error
    +GetCoberturas(luc string) []Cobertura, error
    +GetHistorico(luc string) []HistoricoApolice, error
    +GetHistoricoGlobal(limit int) []HistoricoApolice, error
    +GetAtividadesRecentes(limit int) []AtividadeRecente, error
    +UpdateObservacoes(luc string, obs string) error
    +Renovar(luc string, novoVencimento Time, novoValor float64, ator string, descricao string) error
    +UpdateResponsavel(luc string, responsavelID int64, ator string) error
    +GetLojas() []LojaInfo, error
    +GetUsuarios() []Usuario, error
    +GetDocumentoByID(id string) Documento, error
    +GetDocumentosByApolice(luc string) []Documento, error
    +CreateDocumento(doc Documento) Documento, error
    +DeleteDocumento(id string) error
}

class Apolice {
    +Luc string
    +Loja string
    +Segmento string
    +Seguradora string
    +Vigencia time.Time
    +Vencimento time.Time
    +Status string
    +Cobertura float64
    +DiasRestantes int
    +Responsavel string
    +ResponsavelID *int64
    +Observacoes string
    +CNPJ string
    +NumeroApolice string
}

class Cobertura {
    +ID int
    +ApoliceLuc string
    +Nome string
    +Descricao string
    +Valor float64
}

class HistoricoApolice {
    +ID int
    +ApoliceLuc string
    +Data time.Time
    +Descricao string
    +Ator string
}

class Documento {
    +ID int
    +ApoliceLuc string
    +Nome string
    +ArquivoPath string
    +DataAdicao time.Time
    +DeletedAt *time.Time
}

class AtividadeRecente {
    +ID string
    +Luc string
    +NomeLoja string
    +Acao string
    +Responsavel string
    +Timestamp time.Time
}

class LojaInfo {
    +Luc string
    +Nome string
    +Segmento string
}

class Payload {
    +Luc string
    +Loja string
    +Segmento string
    +Seguradora string
    +Vigencia string
    +Vencimento string
    +Cobertura float64
    +Responsavel string
    +Observacoes string
}

class Response {
    +Luc string
    +Loja string
    +Segmento string
    +Seguradora string
    +Vigencia string
    +Vencimento string
    +Status string
    +Cobertura float64
    +DiasRestantes int
    +Responsavel string
    +ResponsavelID *int64
    +Observacoes string
    +CNPJ string
    +NumeroApolice string
}

class RenovacaoPayload {
    +NovaVigencia string
    +NovoValor float64
}

class UpdateResponsavelPayload {
    +ResponsavelID int64
}

class ValidationError {
    +Message string
    +Error() string
}

class ErrNotFound <<sentinel>> {
    apÃ³lice nÃ£o encontrada
}

class PostgresRepository {
    -db *sql.DB
    +NewRepository(db *sql.DB) *PostgresRepository
    +List() []Apolice, error
    +Get(luc string) Apolice, error
    +Create(model Apolice) Apolice, error
    +Update(luc string, model Apolice) Apolice, error
    +Delete(luc string) error
    +SearchApolices(query string) []Apolice, error
    +GetCoberturas(luc string) []Cobertura, error
    +GetHistorico(luc string) []HistoricoApolice, error
    +GetHistoricoGlobal(limit int) []HistoricoApolice, error
    +GetAtividadesRecentes(limit int) []AtividadeRecente, error
    +UpdateObservacoes(luc string, obs string) error
    +Renovar(...) error
    +UpdateResponsavel(...) error
    +GetLojas() []LojaInfo, error
    +GetUsuarios() []Usuario, error
    +GetDocumentoByID(id string) Documento, error
    +GetDocumentosByApolice(luc string) []Documento, error
    +CreateDocumento(doc Documento) Documento, error
    +DeleteDocumento(id string) error
    -insertHistoricoTx(tx *sql.Tx, luc, descricao, ator string) error
    -scanApolice(scanner) Apolice, error
}

class Service {
    -repo Repository
    +NewService(repo Repository) *Service
    +List() []Apolice, error
    +SearchApolices(query string) []Apolice, error
    +Get(luc string) Apolice, error
    +Create(payload Payload) Apolice, error
    +Update(luc string, payload Payload) Apolice, error
    +Delete(luc string) error
    +GetFilaDeAcao() []Apolice, error
    +GetCoberturas(luc string) []Cobertura, error
    +GetHistorico(luc string) []HistoricoApolice, error
    +GetHistoricoGlobal(limit int) []HistoricoApolice, error
    +GetAtividadesRecentes(limit int) []AtividadeRecente, error
    +UpdateObservacoes(luc string, obs string) error
    +Renovar(luc string, novaVigencia string, novoValor float64, ator string) error
    +GetLojas() []LojaInfo, error
    +GetDocumentoByID(id string) Documento, error
    +GetDocumentosByApolice(luc string) []Documento, error
    +CreateDocumento(doc Documento) Documento, error
    +DeleteDocumento(id string) error
    -buildModel(payload Payload) Apolice, error
}

class Handler {
    -service *Service
    -auditSvc *audit.Service
    +NewHandler(service *Service, auditSvc *audit.Service) *Handler
    +Collection(w, r) void
    +Item(routePrefix string) http.HandlerFunc
    +FilaDeAcao(w, r) void
    +SearchApolices(w, r) void
    +GetCoberturas(w, r) void
    +GetHistorico(w, r) void
    +GetAtividadesRecentes(w, r) void
    +GetKPIHistory(w, r) void
    +GetExpiringByWeek(w, r) void
    +GetCoverageHistory(w, r) void
    +GetRiskBySegment(w, r) void
    +GetHealthScore(w, r) void
    +UpdateObservacoes(w, r) void
    +UpdateApoliceResponsavel(w, r) void
    +RenovarApolice(w, r) void
    +GetMapLayout(w, r) void
    +GetLojas(w, r) void
    +GetDocumentos(w, r) void
    +UploadDocumento(w, r) void
    +DownloadDocumento(w, r) void
    +DeleteDocumento(w, r) void
    +Exportar(w, r) void
    -logAudit(r, acao, entidade, entidadeID, anterior, novo) void
    -writeError(w, requestID, err) void
}

class DomainFunctions <<module>> {
    +calculateDaysRemaining(vencimento time.Time) int
    +calculatePolicyStatus(vencimento time.Time) string
    +calculateHealthScore(items []Apolice, at time.Time) int
    +buildKPIHistoryPoints(items []Apolice, weeks int, metric string) []KPIHistoryPoint
    +countConformesAt(items []Apolice, at time.Time) int
    +countVencidasAt(items []Apolice, at time.Time) int
    +ParseDate(value string) time.Time, error
    +ToResponse(model Apolice) Response
    +toJSON(v any) string
}

PostgresRepository ..|> Repository : implementa
Service --> Repository : depende da interface
Service ..> Apolice : produz e consome
Service ..> Payload : recebe como entrada
Service ..> DomainFunctions : utiliza calculateDaysRemaining e calculatePolicyStatus
Handler --> Service : chama
Handler ..> Payload : decodifica do body HTTP
Handler ..> Response : retorna via ToResponse()
Handler ..> RenovacaoPayload : decodifica do body HTTP
Handler ..> UpdateResponsavelPayload : decodifica do body HTTP
Handler ..> ValidationError : trata via writeError()
Handler ..> ErrNotFound : trata via writeError()

@enduml
```

### Diagrama de sequÃªncia â€” Fluxo de renovaÃ§Ã£o de apÃ³lice

```plantuml
@startuml
title Fluxo de RenovaÃ§Ã£o de ApÃ³lice (NÃ­vel 4: CÃ³digo â€” Diagrama de SequÃªncia)

actor Gestor
participant "RequireRole\n(Middleware)" as MW
participant "Handler" as H
participant "Service" as S
participant "PostgresRepository" as R
participant "audit.Service" as A
database "PostgreSQL" as DB

Gestor -> MW : POST /api/apolices/{id}/renovar\nAuthorization: Bearer <jwt>
MW -> MW : Valida JWT (HS256)\nVerifica role: admin | gestor

alt Token invÃ¡lido ou role insuficiente
    MW --> Gestor : 401 / 403 Forbidden
end

MW -> H : http.HandlerFunc chamado

H -> H : Decodifica RenovacaoPayload\n{nova_vigencia, novo_valor}

alt JSON invÃ¡lido
    H --> Gestor : 400 Bad Request
end

H -> S : Renovar(luc, novaVigencia, novoValor, ator)
S -> S : ParseDate(novaVigencia)\nValida: layouts "2006-01-02" ou "02/01/2006"

alt Data invÃ¡lida
    S --> H : ValidationError
    H --> Gestor : 400 Bad Request
end

S -> R : Renovar(luc, novoVencimento, novoValor, ator, "RenovaÃ§Ã£o realizada")
R -> DB : BEGIN TRANSACTION
R -> DB : UPDATE seguros SET vencimento=$1, cobertura=$2 WHERE luc=$3
DB --> R : rowsAffected

alt rowsAffected == 0
    R -> DB : ROLLBACK
    R --> S : ErrNotFound
    S --> H : ErrNotFound
    H --> Gestor : 404 Not Found
end

R -> DB : INSERT INTO historico_apolice (apolice_luc, data, descricao, ator)
R -> DB : COMMIT
DB --> R : OK
R --> S : nil (sucesso)
S --> H : nil

H -> A : logAudit(r, "renovar", "apolice", luc, nil, detalhesJSON)
note over A, DB
  goroutine assÃ­ncrona â€” nÃ£o bloqueia a resposta
end note
A --> DB : INSERT INTO audit_logs (user_id, acao, entidade, entidade_id, payload_novo, ip, user_agent, timestamp)

H --> Gestor : 200 OK {"message": "ApÃ³lice renovada com sucesso"}

@enduml
```
        Component(apolice, "internal/apolice", "Go â€” domÃ­nio principal", "CRUD completo de apÃ³lices (/api/apolices). Endpoints de KPIs: /kpis/history, /kpis/expiring-by-week, /kpis/coverage-history, /kpis/risk-by-segment, /kpis/health-score. Fila de aÃ§Ã£o (/api/fila-de-acao). Upload/download de documentos. ExportaÃ§Ã£o CSV. Busca full-text. Map layout e lojas.")

        Component(notificacao, "internal/notificacao", "Go â€” notificaÃ§Ãµes", "GET /api/notificacoes â€” lista notificaÃ§Ãµes do usuÃ¡rio logado. PATCH /api/notificacoes/marcar-lidas. DELETE /api/notificacoes/arquivadas e /{id}. Modelo: tipo vencida ou a_vencer, com JOIN em seguros.")

        Component(audit, "internal/audit", "Go â€” auditoria", "GET /api/admin/audit â€” lista logs de auditoria. POST /api/admin/audit â€” registra aÃ§Ã£o. Recebe eventos de: auth (login), apolice (criar, editar, renovar, excluir, exportar, upload_documento). Tabela audit_logs com payload anterior e novo em JSONB.")

        Component(database, "internal/database/postgres.go", "Go â€” conexÃ£o DB", "Abre e valida conexÃ£o com PostgreSQL via database/sql. Recebe PostgresConfig (host, port, user, password, dbname, sslmode).")

        Component(response, "pkg/response", "Go â€” helpers HTTP", "Padroniza responses JSON: Success(w, status, data, requestID) e Fail(w, status, message, requestID, errors).")
    }

    ContainerDb(db, "PostgreSQL", "flamboyant_seguros", "Tabelas: usuarios, seguros, coberturas, documentos, historico_apolice, audit_logs, notificacoes")

    Rel(main, app, "Chama")
    Rel(app, cfg, "Carrega configuraÃ§Ã£o")
    Rel(app, database, "Inicializa conexÃ£o")
    Rel(app, mw, "Aplica chain de middlewares")
    Rel(app, auth, "Registra rotas e injeta auditSvc + jwtSecret")
    Rel(app, apolice, "Registra rotas e injeta auditSvc + jwtSecret")
    Rel(app, notificacao, "Registra rotas e injeta jwtSecret")
    Rel(app, audit, "Registra rotas, retorna Service compartilhado")
    Rel(auth, audit, "Injeta AuditService â€” registra login e aÃ§Ãµes de usuÃ¡rio")
    Rel(apolice, audit, "Injeta AuditService â€” registra CRUD e uploads")
    Rel(auth, response, "Usa para serializar responses")
    Rel(apolice, response, "Usa para serializar responses")
    Rel(notificacao, response, "Usa para serializar responses")
    Rel(audit, response, "Usa para serializar responses")
    Rel(database, db, "SQL/TCP")
    Rel(apolice, database, "Queries em seguros, coberturas, documentos, historico_apolice")
    Rel(auth, database, "Queries em usuarios")
    Rel(notificacao, database, "Queries em notificacoes JOIN seguros")
    Rel(audit, database, "Queries em audit_logs")
```

---

## Tabela de entidades do banco de dados

ExtraÃ­da diretamente de `migrations/init/01_schema.sql`:

| Tabela | Chave principal | DescriÃ§Ã£o |
|---|---|---|
| `usuarios` | `id SERIAL` | UsuÃ¡rios do sistema com role: admin, gestor, visualizador |
| `seguros` | `luc VARCHAR` | ApÃ³lices de seguro (entidade central) |
| `coberturas` | `id SERIAL` | Coberturas vinculadas a uma apÃ³lice (`apolice_luc`) |
| `documentos` | `id SERIAL` | Arquivos PDF/imagem vinculados a uma apÃ³lice (soft delete) |
| `historico_apolice` | `id SERIAL` | Timeline de eventos de cada apÃ³lice |
| `audit_logs` | `id BIGSERIAL` | Log imutÃ¡vel de todas as aÃ§Ãµes (payload JSONB antes/depois) |
| `notificacoes` | `id SERIAL` | Alertas de apÃ³lices vencidas ou a vencer por usuÃ¡rio |

---

## Self-Refine â€” VerificaÃ§Ã£o de fronteiras

| VerificaÃ§Ã£o | Resultado |
|---|---|
| Os nomes dos componentes correspondem aos pacotes Go reais? | âœ… Sim â€” `internal/apolice`, `internal/auth`, `internal/notificacao`, `internal/audit`, `internal/middleware`, `internal/database`, `pkg/config`, `pkg/response` |
| As rotas listadas existem no cÃ³digo? | âœ… Sim â€” verificadas em `*/routes.go` de cada pacote |
| Os roles (admin, gestor, visualizador) batem com o schema? | âœ… Sim â€” `CHECK constraint` em `migrations/init/01_schema.sql` |
| HÃ¡ sistemas externos que foram omitidos? | âœ… Nenhum â€” o cÃ³digo nÃ£o consome APIs externas reais |
| A dependÃªncia `audit â†’ outros pacotes` estÃ¡ representada corretamente? | âœ… Sim â€” `auditSvc` Ã© criado em `app.go` e injetado explicitamente |
| O frontend serve arquivos estÃ¡ticos via backend em produÃ§Ã£o? | âœ… Sim â€” `registerStaticFiles` em `app.go` + `build.outDir` em `vite.config.ts` |

---

## NÃ­vel 4 â€” Diagrama de CÃ³digo (Pacote: `internal/apolice`)

### RaciocÃ­nio tÃ©cnico

O NÃ­vel 4 foca no contÃªiner de maior complexidade de domÃ­nio â€” o pacote `internal/apolice`. Ele Ã© o Ãºnico que define uma **interface `Repository`** explÃ­cita (em `repository.go`), seguindo o padrÃ£o de inversÃ£o de dependÃªncia do Go. Isso significa que `Service` depende da abstraÃ§Ã£o, nÃ£o da implementaÃ§Ã£o concreta `PostgresRepository`. O `Handler` recebe `*Service` e `*audit.Service` opcionais por injeÃ§Ã£o direta. O `dto.go` separa os tipos de entrada (`Payload`, `RenovacaoPayload`, `UpdateResponsavelPayload`) dos tipos de domÃ­nio (`Apolice`, `Cobertura`, `Documento`) e dos tipos de saÃ­da (`Response`). O `errors.go` define os dois Ãºnicos erros de domÃ­nio: `ErrNotFound` (sentinel) e `ValidationError` (tipo customizado com interface `error`). Todos os mÃ©todos do `Handler` que escrevem dados chamam `logAudit()` â€” wrapper que invoca `audit.Service.LogFromRequest()` de forma segura (nil-safe).

```mermaid
classDiagram
    direction TB

    %% â”€â”€ Interfaces â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
    class Repository {
        <<interface>>
        +List() []Apolice, error
        +Get(luc string) Apolice, error
        +Create(model Apolice) Apolice, error
        +Update(luc string, model Apolice) Apolice, error
        +Delete(luc string) error
        +SearchApolices(query string) []Apolice, error
        +GetCoberturas(luc string) []Cobertura, error
        +GetHistorico(luc string) []HistoricoApolice, error
        +GetHistoricoGlobal(limit int) []HistoricoApolice, error
        +GetAtividadesRecentes(limit int) []AtividadeRecente, error
        +UpdateObservacoes(luc string, obs string) error
        +Renovar(luc string, novoVencimento Time, novoValor float64, ator string, descricao string) error
        +UpdateResponsavel(luc string, responsavelID int64, ator string) error
        +GetLojas() []LojaInfo, error
        +GetUsuarios() []Usuario, error
        +GetDocumentoByID(id string) Documento, error
        +GetDocumentosByApolice(luc string) []Documento, error
        +CreateDocumento(doc Documento) Documento, error
        +DeleteDocumento(id string) error
    }

    %% â”€â”€ Structs de domÃ­nio (model.go) â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
    class Apolice {
        +Luc string
        +Loja string
        +Segmento string
        +Seguradora string
        +Vigencia time.Time
        +Vencimento time.Time
        +Status string
        +Cobertura float64
        +DiasRestantes int
        +Responsavel string
        +ResponsavelID *int64
        +Observacoes string
        +CNPJ string
        +NumeroApolice string
    }

    class Cobertura {
        +ID int
        +ApoliceLuc string
        +Nome string
        +Descricao string
        +Valor float64
    }

    class HistoricoApolice {
        +ID int
        +ApoliceLuc string
        +Data time.Time
        +Descricao string
        +Ator string
    }

    class Documento {
        +ID int
        +ApoliceLuc string
        +Nome string
        +ArquivoPath string
        +DataAdicao time.Time
        +DeletedAt *time.Time
    }

    class AtividadeRecente {
        +ID string
        +Luc string
        +NomeLoja string
        +Acao string
        +Responsavel string
        +Timestamp time.Time
    }

    class LojaInfo {
        +Luc string
        +Nome string
        +Segmento string
    }

    %% â”€â”€ DTOs (dto.go) â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
    class Payload {
        +Luc string
        +Loja string
        +Segmento string
        +Seguradora string
        +Vigencia string
        +Vencimento string
        +Cobertura float64
        +Responsavel string
        +Observacoes string
    }

    class Response {
        +Luc string
        +Loja string
        +Segmento string
        +Seguradora string
        +Vigencia string
        +Vencimento string
        +Status string
        +Cobertura float64
        +DiasRestantes int
        +Responsavel string
        +ResponsavelID *int64
        +Observacoes string
        +CNPJ string
        +NumeroApolice string
    }

    class RenovacaoPayload {
        +NovaVigencia string
        +NovoValor float64
    }

    class UpdateResponsavelPayload {
        +ResponsavelID int64
    }

    %% â”€â”€ Erros (errors.go) â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
    class ValidationError {
        +Message string
        +Error() string
    }

    class ErrNotFound {
        <<sentinel>>
        apÃ³lice nÃ£o encontrada
    }

    %% â”€â”€ ImplementaÃ§Ã£o concreta (repository.go) â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
    class PostgresRepository {
        -db *sql.DB
        +NewRepository(db *sql.DB) *PostgresRepository
        +List() []Apolice, error
        +Get(luc string) Apolice, error
        +Create(model Apolice) Apolice, error
        +Update(luc string, model Apolice) Apolice, error
        +Delete(luc string) error
        +SearchApolices(query string) []Apolice, error
        +GetCoberturas(luc string) []Cobertura, error
        +GetHistorico(luc string) []HistoricoApolice, error
        +GetHistoricoGlobal(limit int) []HistoricoApolice, error
        +GetAtividadesRecentes(limit int) []AtividadeRecente, error
        +UpdateObservacoes(luc string, obs string) error
        +Renovar(...) error
        +UpdateResponsavel(...) error
        +GetLojas() []LojaInfo, error
        +GetUsuarios() []Usuario, error
        +GetDocumentoByID(id string) Documento, error
        +GetDocumentosByApolice(luc string) []Documento, error
        +CreateDocumento(doc Documento) Documento, error
        +DeleteDocumento(id string) error
        -insertHistoricoTx(tx *sql.Tx, luc, descricao, ator string) error
        -scanApolice(scanner) Apolice, error
    }

    %% â”€â”€ Camada de serviÃ§o (service.go) â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
    class Service {
        -repo Repository
        +NewService(repo Repository) *Service
        +List() []Apolice, error
        +SearchApolices(query string) []Apolice, error
        +Get(luc string) Apolice, error
        +Create(payload Payload) Apolice, error
        +Update(luc string, payload Payload) Apolice, error
        +Delete(luc string) error
        +GetFilaDeAcao() []Apolice, error
        +GetCoberturas(luc string) []Cobertura, error
        +GetHistorico(luc string) []HistoricoApolice, error
        +GetHistoricoGlobal(limit int) []HistoricoApolice, error
        +GetAtividadesRecentes(limit int) []AtividadeRecente, error
        +UpdateObservacoes(luc string, obs string) error
        +Renovar(luc string, novaVigencia string, novoValor float64, ator string) error
        +GetLojas() []LojaInfo, error
        +GetDocumentoByID(id string) Documento, error
        +GetDocumentosByApolice(luc string) []Documento, error
        +CreateDocumento(doc Documento) Documento, error
        +DeleteDocumento(id string) error
        -buildModel(payload Payload) Apolice, error
    }

    %% â”€â”€ Camada de apresentaÃ§Ã£o (handler.go) â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
    class Handler {
        -service *Service
        -auditSvc *audit.Service
        +NewHandler(service *Service, auditSvc *audit.Service) *Handler
        +Collection(w, r) void
        +Item(routePrefix string) http.HandlerFunc
        +FilaDeAcao(w, r) void
        +SearchApolices(w, r) void
        +GetCoberturas(w, r) void
        +GetHistorico(w, r) void
        +GetAtividadesRecentes(w, r) void
        +GetKPIHistory(w, r) void
        +GetExpiringByWeek(w, r) void
        +GetCoverageHistory(w, r) void
        +GetRiskBySegment(w, r) void
        +GetHealthScore(w, r) void
        +UpdateObservacoes(w, r) void
        +UpdateApoliceResponsavel(w, r) void
        +RenovarApolice(w, r) void
        +GetMapLayout(w, r) void
        +GetLojas(w, r) void
        +GetDocumentos(w, r) void
        +UploadDocumento(w, r) void
        +DownloadDocumento(w, r) void
        +DeleteDocumento(w, r) void
        +Exportar(w, r) void
        -logAudit(r, acao, entidade, entidadeID, anterior, novo) void
        -writeError(w, requestID, err) void
    }

    %% â”€â”€ FunÃ§Ãµes auxiliares de domÃ­nio (service.go) â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
    class DomainFunctions {
        <<module>>
        +calculateDaysRemaining(vencimento time.Time) int
        +calculatePolicyStatus(vencimento time.Time) string
        +calculateHealthScore(items []Apolice, at time.Time) int
        +buildKPIHistoryPoints(items []Apolice, weeks int, metric string) []KPIHistoryPoint
        +countConformesAt(items []Apolice, at time.Time) int
        +countVencidasAt(items []Apolice, at time.Time) int
        +ParseDate(value string) time.Time, error
        +ToResponse(model Apolice) Response
        +toJSON(v any) string
    }

    %% â”€â”€ Relacionamentos â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
    PostgresRepository ..|> Repository : implementa
    Service --> Repository : depende da interface
    Service ..> Apolice : produz e consome
    Service ..> Payload : recebe como entrada
    Service ..> DomainFunctions : utiliza calculateDaysRemaining e calculatePolicyStatus
    Handler --> Service : chama
    Handler ..> Payload : decodifica do body HTTP
    Handler ..> Response : retorna via ToResponse()
    Handler ..> RenovacaoPayload : decodifica do body HTTP
    Handler ..> UpdateResponsavelPayload : decodifica do body HTTP
    Handler ..> ValidationError : trata via writeError()
    Handler ..> ErrNotFound : trata via writeError()
```

---

## NÃ­vel 4 â€” Diagrama de SequÃªncia: Fluxo de RenovaÃ§Ã£o de ApÃ³lice

### RaciocÃ­nio tÃ©cni
