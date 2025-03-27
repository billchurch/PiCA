# PiCA Architecture Diagrams

This document contains the architecture diagrams for the PiCA Certificate Authority system using Mermaid syntax.

## System Overview

```mermaid
flowchart TD
    RootCA[Root CA - Raspberry Pi</br>Offline] -- "Signs</br>(manual process)" --> SubCA
    YK1[YubiKey</br>Root CA Key] --- RootCA
    
    SubCA[Sub CA - Raspberry Pi</br>Online] <--> Client[Client Workstation</br>CSR Generation</br>Certificate Requests]
    YK2[YubiKey</br>Sub CA Key] --- SubCA
    SubCA <--> Services[Services</br>Certificate Usage]
    
    subgraph Root["Root CA Environment (Air-gapped)"]
        RootCA
        YK1
    end
    
    subgraph Trust["Trust Infrastructure"]
        SubCA
        YK2
    end
    
    subgraph Users["End Users"]
        Client
        Services
    end
    
    style Root fill:#f9f9f9,stroke:#333,stroke-width:2px
    style Trust fill:#f0f0ff,stroke:#333,stroke-width:2px
    style Users fill:#f0fff0,stroke:#333,stroke-width:2px
```

## Certificate Hierarchy

```mermaid
flowchart TD
    Root[Root CA] --> Sub[Sub CA]
    Sub --> Server[Server Certificates]
    Sub --> Client[Client Certificates]
    Sub --> Code[Code Signing Certificates]
    Sub --> Email[Email Certificates]
    
    style Root fill:#f9f,stroke:#333,stroke-width:2px,color:#000
    style Sub fill:#bbf,stroke:#333,stroke-width:2px,color:#000
    style Server fill:#dfd,stroke:#333,stroke-width:1px,color:#000
    style Client fill:#dfd,stroke:#333,stroke-width:1px,color:#000
    style Code fill:#dfd,stroke:#333,stroke-width:1px,color:#000
    style Email fill:#dfd,stroke:#333,stroke-width:1px,color:#000
```
## Software Architecture

```mermaid
flowchart TD
    subgraph Apps["Applications"]
        CLI[PiCA CLI]
        Web[PiCA Web]
    end
    
    subgraph Libs["Libraries"]
        CFSSL[CFSSL]
        Charm[Charm/Bubble]
        YKLib[YubiKey SDK]
        Crypto[Go Crypto]
    end
    
    subgraph Sys["System Layer"]
        Go[Go Lang]
        PKCS[PCSC/PKCS]
    end
    
    subgraph OS["Operating System"]
        RPiOS[Raspberry Pi OS]
    end
    
    subgraph HW["Hardware"]
        RPi[Raspberry Pi]
        YK[YubiKey]
    end
    
    CLI --> CFSSL
    CLI --> Charm
    CLI --> YKLib
    Web --> CFSSL
    Web --> YKLib
    Web --> Crypto
    
    CFSSL --> Go
    Charm --> Go
    YKLib --> Go
    YKLib --> PKCS
    Crypto --> Go
    
    Go --> RPiOS
    PKCS --> RPiOS
    
    RPiOS --> RPi
    RPiOS --> YK
    
    style Apps fill:#f9f9f9,stroke:#333,stroke-width:2px
    style Libs fill:#f0f0ff,stroke:#333,stroke-width:2px
    style Sys fill:#f0fff0,stroke:#333,stroke-width:2px
    style OS fill:#fff0f0,stroke:#333,stroke-width:2px
    style HW fill:#f0f8ff,stroke:#333,stroke-width:2px
```

## Certificate Issuance Process

```mermaid
sequenceDiagram
    participant Client
    participant WebUI as Web Interface
    participant SubCA as Sub CA
    participant YubiKey
    
    Client->>Client: Generate key pair and CSR
    Client->>WebUI: Submit CSR
    WebUI->>SubCA: Forward CSR
    SubCA->>SubCA: Validate CSR
    SubCA->>YubiKey: Request signing operation
    YubiKey->>YubiKey: Sign certificate with private key
    YubiKey->>SubCA: Return signature
    SubCA->>SubCA: Create certificate with signature
    SubCA->>WebUI: Return signed certificate
    WebUI->>Client: Deliver certificate
    Client->>Client: Install certificate
```

## Certificate Revocation Process

```mermaid
sequenceDiagram
    participant Admin
    participant WebUI as Web Interface
    participant SubCA as Sub CA
    participant YubiKey
    participant CRL as CRL Distribution Point
    
    Admin->>WebUI: Submit revocation request
    WebUI->>SubCA: Forward request
    SubCA->>SubCA: Validate request
    SubCA->>SubCA: Update revocation list
    SubCA->>YubiKey: Request CRL signing
    YubiKey->>YubiKey: Sign CRL with private key
    YubiKey->>SubCA: Return signature
    SubCA->>SubCA: Create signed CRL
    SubCA->>CRL: Publish updated CRL
    SubCA->>WebUI: Confirm revocation
    WebUI->>Admin: Display confirmation
```

## Deployment Architecture

```mermaid
flowchart TD
    subgraph Offline["Offline Environment"]
        RootCA[Root CA Raspberry Pi]
        RootYK[Root CA YubiKey]
        RootCA --- RootYK
    end
    
    subgraph Online["Online Environment"]
        SubCA[Sub CA Raspberry Pi]
        SubYK[Sub CA YubiKey]
        SubCA --- SubYK
        
        Firewall --- SubCA
        Firewall --- Web[Web Server]
        
        Client1[Client Systems] --- Firewall
        Client2[Client Systems] --- Firewall
    end
    
    RootCA -.-> |"Manual</br>Certificate</br>Transfer"| SubCA
    
    style Offline fill:#fff0f0,stroke:#333,stroke-width:2px,stroke-dasharray: 5 5
    style Online fill:#f0f0ff,stroke:#333,stroke-width:2px
```

## Project Structure

```mermaid
flowchart TD
    Root[PiCA Project] --> CMD[cmd/]
    Root --> Internal[internal/]
    Root --> Web[web/]
    Root --> RPI[rpi-images/]
    Root --> Configs[configs/]
    Root --> Scripts[scripts/]
    Root --> Docs[docs/]
    
    CMD --> CLI[pica/]
    CMD --> WebApp[pica-web/]
    
    Internal --> CA[ca/]
    Internal --> UI[ui/]
    Internal --> YK[yubikey/]
    
    CA --> Commands[commands/]
    UI --> Pages[pages/]
    
    Web --> API[api/]
    Web --> HTML[html/]
    
    RPI --> RootImg[root-ca/]
    RPI --> SubImg[sub-ca/]
    
    Configs --> CFSSL[cfssl/]
    
    style Root fill:#f9f,stroke:#333,stroke-width:2px,color:#000
    style CMD fill:#bbf,stroke:#333,stroke-width:1px,color:#000
    style Internal fill:#bbf,stroke:#333,stroke-width:1px,color:#000
    style Web fill:#bbf,stroke:#333,stroke-width:1px,color:#000
    style RPI fill:#bbf,stroke:#333,stroke-width:1px,color:#000
    style Configs fill:#bbf,stroke:#333,stroke-width:1px,color:#000
    style Scripts fill:#bbf,stroke:#333,stroke-width:1px,color:#000
    style Docs fill:#bbf,stroke:#333,stroke-width:1px,color:#000
```
## YubiKey Integration

```mermaid
flowchart LR
    subgraph App["PiCA Application"]
        CLI[CLI Module]
        Web[Web Module]
        CA[CA Module]
    end
    
    subgraph YK["YubiKey Integration"]
        YKAPI[YubiKey API]
        PIV[PIV Module]
        KeyOps[Key Operations]
    end
    
    subgraph HW["Hardware"]
        YubiKey[YubiKey Device]
        Slot1[Slot 9A 82</br>Root CA Key]
        Slot2[Slot 9B 83</br>Sub CA Key]
    end
    
    CA --> YKAPI
    CLI --> CA
    Web --> CA
    
    YKAPI --> PIV
    PIV --> KeyOps
    KeyOps --> YubiKey
    
    YubiKey --> Slot1
    YubiKey --> Slot2
    
    style App fill:#f0f0ff,stroke:#333,stroke-width:2px
    style YK fill:#fff0f0,stroke:#333,stroke-width:2px
    style HW fill:#f0fff0,stroke:#333,stroke-width:2px
```
