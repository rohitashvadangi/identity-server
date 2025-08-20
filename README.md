# identity-server

A modern **Identity Provider (IdP)** written in Go — implementing **OAuth 2.1, OpenID Connect (OIDC), SAML 2.0, and WebAuthn** with a **security-first design**.

This project is built to demonstrate **cloud architecture, application security, and identity engineering expertise**.  
It’s my personal project to showcase real-world design for **secure authentication & authorization at scale**.

---

## 🚀 Features (Work in Progress)

### ✅ Core Protocols
- [ ] OAuth 2.1: Authorization Code + PKCE, Client Credentials
- [ ] OpenID Connect Core: ID Token, UserInfo
- [ ] Refresh Tokens, Rotation & Reuse Detection
- [ ] Dynamic Client Registration (RFC 7591)
- [ ] Discovery: `/.well-known/openid-configuration`
- [ ] JWKS endpoint (rotating signing keys)

### 🔐 Security-first Enhancements
- [ ] Pushed Authorization Requests (PAR)
- [ ] JWT-secured Authorization Response (JARM)
- [ ] Sender-constrained tokens: DPoP + mTLS
- [ ] JWT Access Tokens (RFC 9068)
- [ ] Token Exchange (RFC 8693)
- [ ] FAPI2 & non-repudiation hooks

### 🧑‍💻 User Experience
- [ ] Modern MFA: WebAuthn / Passkeys
- [ ] TOTP (Authenticator App), Recovery Codes
- [ ] Step-up Authentication (per-scope ACR/AMR)
- [ ] Consent Screen with fine-grained scopes

### 🌍 Federation
- [ ] OIDC Broker (Google, GitHub, etc.)
- [ ] SAML 2.0 IdP (minimal implementation)
- [ ] SAML SP adapter

### 🔑 Key Management
- [ ] Pluggable KMS/HSM (AWS KMS, CloudHSM)
- [ ] Automated Key Rotation
- [ ] KID pinning, health checks

### 🏢 Multi-tenancy
- [ ] Tenant-scoped clients, keys, and branding
- [ ] Per-tenant discovery endpoints

### 📊 Operations & Observability
- [ ] Structured logs + audit logs
- [ ] Prometheus metrics
- [ ] OpenTelemetry tracing

### 🛠️ Developer Experience
- [ ] Go SDK for RPs & Resource Servers
- [ ] Java Servlet filter
- [ ] Example apps: Go, Java, Next.js
- [ ] Terraform module for AWS deployment

---

## 📐 Architecture (planned)

```mermaid
flowchart LR
  subgraph Client
    App[Web/Mobile App]
  end

  subgraph IdentityServer
    Auth[Auth Service]
    Keys[Key Service]
    Storage[(Postgres)]
    Cache[(Redis)]
  end

  subgraph External
    AWSKMS[AWS KMS/HSM]
    OIDC[OIDC Providers]
    SAML[SAML SPs]
  end

  App -->|OAuth/OIDC| Auth
  Auth --> Storage
  Auth --> Cache
  Auth --> Keys
  Keys --> AWSKMS
  Auth --> OIDC
  Auth --> SAML


⚠️ Security Note:
Private key is committed for demo purposes only.
In a real-world setup, private keys must be stored in a secure vault
(AWS KMS, Azure Key Vault, GCP KMS, or HashiCorp Vault).

