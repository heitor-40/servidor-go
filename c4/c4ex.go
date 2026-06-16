```mermaid
flowchart TB

Admin[Administrador]
Gestor[Gestor]
Logista[Visualizador e Logista]

Sistema[Sistema de Seguros<br>CRM para gestão de apólices,<br>documentos, notificações,<br>auditoria e KPIs]

Seguradora[Seguradora Parceira]
Sinistros[Sistema de Sinistros]
Banco[Banco de Dados Externo]

Admin -->|Gerencia usuários e auditorias| Sistema
Gestor -->|Gerencia apólices e documentos| Sistema
Logista -->|Consulta informações| Sistema

Sistema -->|Valida apólices e coberturas| Seguradora

Sistema -->|Registra e consulta sinistros| Sinistros
Sinistros -->|Retorna andamento e status| Sistema

Sistema -->|Consulta contratos de locação| Banco
```
