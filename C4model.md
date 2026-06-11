# **1\. Diagrama de Contexto (Nível 1)**

O Diagrama de Contexto mostra o sistema de seguros como um todo, identificando os atores que interagem com ele e os sistemas externos com os quais se integra. Neste nível os detalhes de implementação são omitidos - o foco é em quem usa o sistema e qual valor ele entrega.

## **1.1 Atores internos**

Os seguintes perfis de usuário interagem diretamente com o sistema:

| **Ator**                | **Responsabilidade**                                                                                                       |
| ----------------------- | -------------------------------------------------------------------------------------------------------------------------- |
| **Gestor do shopping**  | Administra apólices do empreendimento, aprova contratos, acessa dashboards gerenciais e relatórios de cobertura de riscos. |
| **Corretor de seguros** | Gerencia contratos de seguro dos lojistas, propõe apólices, acompanha renovações e sinistros.                              |
| **Lojista / locatário** | Consulta e acompanha sua própria apólice, recebe notificações de vencimento e visualiza coberturas contratadas.            |

## **1.2 Sistemas externos**

O sistema depende dos seguintes serviços e plataformas externas:

| **Sistema externo**         | **Papel na arquitetura**                                                                                                 |
| --------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| **Seguradora parceira**     | Emite, valida e cancela apólices via API. O sistema de seguros atua como intermediário entre os usuários e a seguradora. |
| **Gateway de pagamento**    | Processa os prêmios de seguro (pagamentos periódicos). Exemplos: Stripe, Pagar.me.                                       |
| **Serviço de notificações** | Envia alertas de vencimento, renovação e sinistro via e-mail e SMS. Exemplos: SendGrid, Twilio.                          |
| **ERP do shopping**         | Fornece dados de contratos de locação que servem como base para emissão e vinculação de apólices.                        |

# **2\. Diagrama de Contêiner (Nível 2)**

O Diagrama de Contêiner aprofunda o nível de contexto mostrando os aplicativos (contêineres de software) que compõem o sistema. Cada contêiner é uma unidade implantável ou executável de forma independente.

## **2.1 Contêineres internos**

| **Contêiner**               | **Tecnologia**                         | **Responsabilidade**                                                                                                                                           |
| --------------------------- | -------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Frontend (SPA)**          | _React, Vite, TypeScript, TailwindCSS_ | Interface do usuário acessada pelo navegador. Comunica-se com a API do backend via JSON/HTTPS e redireciona fluxos de autenticação para o serviço de auth.     |
| **API Backend**             | _Go (Golang), RESTful_                 | Contém as regras de negócio: gestão de apólices, sinistros, relatórios e integrações. Orquestra chamadas para sistemas externos e persiste dados no banco.     |
| **Serviço de autenticação** | _JWT, OAuth2_                          | Gerencia autenticação e autorização dos diferentes perfis de usuário (gestor, corretor, lojista). Emite e valida tokens de acesso.                             |
| **Fila de mensagens**       | _RabbitMQ / Redis Streams_             | Desacopla operações assíncronas: emissão de apólices, processamento de sinistros e disparo de notificações. Evita bloqueio da API em chamadas externas lentas. |
| **Worker de notificações**  | _Go (worker process)_                  | Consome eventos da fila e dispara notificações (vencimento de apólice, confirmação de sinistro, renovação pendente) via serviço externo de e-mail/SMS.         |
| **Banco de dados**          | _PostgreSQL_                           | Persistência relacional de usuários, apólices, sinistros, contratos de locação vinculados, logs e auditoria.                                                   |

## **2.2 Comunicação entre contêineres**

| **Origem**            | **Destino**               | **Protocolo / Detalhe**             |
| --------------------- | ------------------------- | ----------------------------------- |
| **Usuário (browser)** | _Frontend (SPA)_          | HTTPS                               |
| **Frontend**          | _Serviço de autenticação_ | JSON/HTTPS - login, troca de token  |
| **Frontend**          | _API Backend_             | JSON/HTTPS - operações de negócio   |
| **API Backend**       | _Banco de dados_          | SQL/TCP (driver nativo Go)          |
| **API Backend**       | _Fila de mensagens_       | AMQP ou Redis protocol              |
| **Fila de mensagens** | _Worker de notificações_  | Consumo assíncrono de eventos       |
| **API Backend**       | _Seguradora parceira_     | JSON/HTTPS - API REST externa       |
| **API Backend**       | _Gateway de pagamento_    | JSON/HTTPS - API REST externa       |
| **Worker**            | _Serviço de notificações_ | JSON/HTTPS - SendGrid / Twilio API  |
| **API Backend**       | _ERP do shopping_         | JSON/HTTPS ou SOAP (depende do ERP) |

# **3\. Resumo das tecnologias**

| **Camada**         | **Tecnologia**                            | **Observação**                                          |
| ------------------ | ----------------------------------------- | ------------------------------------------------------- |
| **Frontend**       | _React + Vite + TypeScript + TailwindCSS_ | SPA com tipagem estática e build rápido.                |
| **Backend**        | _Go (Golang)_                             | API RESTful com alta performance e concorrência nativa. |
| **Autenticação**   | _JWT / OAuth2_                            | Controle de acesso por perfil de usuário.               |
| **Mensageria**     | _RabbitMQ ou Redis Streams_               | Operações assíncronas desacopladas.                     |
| **Banco de dados** | _PostgreSQL_                              | Persistência relacional com suporte a transações ACID.  |
| **Infraestrutura** | _Docker Compose_                          | Orquestração local de todos os serviços.                |
| **Notificações**   | _SendGrid / Twilio_                       | E-mail e SMS transacionais.                             |
| **Pagamentos**     | _Stripe / Pagar.me_                       | Processamento de prêmios de seguro.                     |
