# Certificate Lifecycle in PiCA

This document visualizes the complete lifecycle of certificates in the PiCA system.

## Certificate Lifecycle Overview

```mermaid
stateDiagram-v2
    [*] --> Requested: CSR Submitted
    
    Requested --> Pending: CSR Validation
    Pending --> Rejected: Validation Failed
    Pending --> Issued: Certificate Signed
    
    Issued --> Active: Certificate Deployed
    
    Active --> Expired: Time Passed
    Active --> Revoked: Requested by Owner/Admin
    
    Rejected --> [*]
    Expired --> [*]
    Revoked --> [*]
    
    note right of Requested
        CSR submitted via 
        Web UI or API
    end note
    
    note right of Pending
        - CSR format validation
        - Subject validation
        - Policy compliance check
    end note
    
    note right of Issued
        Certificate signed by 
        Sub CA YubiKey
    end note
    
    note right of Active
        Certificate in use by
        services or clients
    end note
    
    note right of Revoked
        - Added to CRL
        - OCSP status updated
    end note
```

## Root CA Certificate Creation Process

```mermaid
sequenceDiagram
    participant Admin
    participant RootCA as PiCA Root CA
    participant YubiKey as Root YubiKey
    
    Admin->>RootCA: Initialize Root CA
    RootCA->>YubiKey: Generate key in PIV slot 9A
    YubiKey-->>RootCA: Return public key
    RootCA->>RootCA: Create Root CA CSR
    RootCA->>YubiKey: Sign Root CA certificate
    YubiKey-->>RootCA: Return signed certificate
    RootCA->>RootCA: Store Root CA certificate
    RootCA-->>Admin: Confirm Root CA creation
    
    Note over Admin,YubiKey: Root private key never leaves YubiKey
```

## Sub CA Certificate Creation Process

```mermaid
sequenceDiagram
    participant Admin
    participant SubCA as PiCA Sub CA
    participant SubYK as Sub YubiKey
    participant RootCA as PiCA Root CA
    participant RootYK as Root YubiKey
    
    Admin->>SubCA: Initialize Sub CA
    SubCA->>SubYK: Generate key in PIV slot 9B
    SubYK-->>SubCA: Return public key
    SubCA->>SubCA: Create Sub CA CSR
    
    Admin->>Admin: Transfer CSR to Root CA (offline)
    Admin->>RootCA: Submit Sub CA CSR
    RootCA->>RootCA: Validate CSR
    RootCA->>RootYK: Sign Sub CA certificate
    RootYK-->>RootCA: Return signed certificate
    RootCA-->>Admin: Provide signed certificate
    
    Admin->>Admin: Transfer certificate to Sub CA
    Admin->>SubCA: Import Sub CA certificate
    SubCA->>SubCA: Store certificate
    SubCA-->>Admin: Confirm Sub CA ready
    
    Note over Admin,RootYK: Manual, air-gapped process
```

## End-Entity Certificate Issuance

```mermaid
sequenceDiagram
    participant User
    participant WebUI as Web Interface
    participant SubCA as PiCA Sub CA
    participant SubYK as Sub YubiKey
    
    User->>User: Generate key pair and CSR
    User->>WebUI: Submit CSR via web UI
    WebUI->>SubCA: Forward CSR
    
    SubCA->>SubCA: Validate CSR format
    SubCA->>SubCA: Check subject information
    SubCA->>SubCA: Verify CSR signature
    SubCA->>SubCA: Apply certificate policy
    
    SubCA->>SubYK: Sign certificate
    SubYK-->>SubCA: Return signature
    
    SubCA->>SubCA: Create certificate
    SubCA->>SubCA: Record in database
    SubCA->>WebUI: Return signed certificate
    WebUI->>User: Deliver certificate
    
    Note over SubCA,SubYK: All signing operations occur on YubiKey
```

## Certificate Revocation Process

```mermaid
sequenceDiagram
    participant User
    participant WebUI as Web Interface
    participant SubCA as PiCA Sub CA
    participant SubYK as Sub YubiKey
    participant CRL as CRL Distribution
    
    User->>WebUI: Request certificate revocation
    WebUI->>SubCA: Forward revocation request
    
    SubCA->>SubCA: Validate request
    SubCA->>SubCA: Update revocation database
    SubCA->>SubCA: Generate new CRL
    
    SubCA->>SubYK: Sign CRL
    SubYK-->>SubCA: Return signature
    
    SubCA->>SubCA: Complete CRL
    SubCA->>CRL: Publish CRL to distribution points
    SubCA->>WebUI: Confirm revocation
    WebUI->>User: Display confirmation
```

## Certificate Renewal Process

```mermaid
sequenceDiagram
    participant User
    participant WebUI as Web Interface
    participant SubCA as PiCA Sub CA
    participant SubYK as Sub YubiKey
    
    User->>User: Generate new CSR using existing key or new key
    User->>WebUI: Submit renewal request with CSR
    WebUI->>SubCA: Forward renewal request
    
    SubCA->>SubCA: Validate CSR
    SubCA->>SubCA: Check existing certificate
    SubCA->>SubCA: Verify renewal eligibility
    
    SubCA->>SubYK: Sign new certificate
    SubYK-->>SubCA: Return signature
    
    SubCA->>SubCA: Create new certificate
    SubCA->>SubCA: Update certificate database
    SubCA->>WebUI: Return renewed certificate
    WebUI->>User: Deliver new certificate
    
    Note over User,SubCA: Original certificate remains valid until expiration
```

## Certificate Status Checking

```mermaid
sequenceDiagram
    participant Client
    participant OCSP as OCSP Responder
    participant SubCA as PiCA Sub CA
    participant CRL as CRL Repository
    
    Note over Client,CRL: Option 1: OCSP Checking
    Client->>OCSP: OCSP Request
    OCSP->>SubCA: Query certificate status
    SubCA-->>OCSP: Return status
    OCSP-->>Client: OCSP Response
    
    Note over Client,CRL: Option 2: CRL Checking
    Client->>CRL: Download CRL
    CRL-->>Client: Return CRL
    Client->>Client: Check certificate against CRL
```

## Root CA Renewal Process (When Needed)

```mermaid
sequenceDiagram
    participant Admin
    participant RootCA as PiCA Root CA
    participant NewRootYK as New Root YubiKey
    participant OldRootYK as Old Root YubiKey
    participant SubCA as PiCA Sub CA
    
    Admin->>RootCA: Initiate Root CA renewal
    RootCA->>NewRootYK: Generate new key in PIV slot
    NewRootYK-->>RootCA: Return public key
    RootCA->>RootCA: Create new Root CA CSR
    
    RootCA->>OldRootYK: Sign new Root CA cert with old key
    OldRootYK-->>RootCA: Return cross-signed certificate
    
    RootCA->>RootCA: Store new Root CA certificate
    RootCA-->>Admin: Provide new Root CA certificate
    
    Admin->>Admin: Distribute new Root CA certificate
    Admin->>SubCA: Update Root CA certificate
    SubCA->>SubCA: Trust new Root CA
    
    Note over Admin,SubCA: Carefully planned transition process
```

This document visualizes all the critical processes in the PiCA Certificate Authority system, showing how certificates are managed throughout their lifecycle.
