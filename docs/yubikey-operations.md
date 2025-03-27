# YubiKey Operations in PiCA

This document visualizes the YubiKey operations in the PiCA Certificate Authority system.

## YubiKey PIV Slot Allocation

```mermaid
graph TD
    YK[YubiKey] --- Slot9A[Slot 9A/82</br>Root CA]
    YK --- Slot9B[Slot 9B/83</br>Sub CA]
    YK --- Slot9C["Slot 9C/84</br>OCSP Signing</br>(Optional)"]
    YK --- Slot9D["Slot 9D/85</br>Audit Key</br>(Optional)"]
    
    style YK fill:#f9f9f9,stroke:#000,stroke-width:2px, color:#000
    style Slot9A fill:#ff8080,stroke:#000,stroke-width:1px, color:#000
    style Slot9B fill:#80ff80,stroke:#000,stroke-width:1px, color:#000
    style Slot9C fill:#8080ff,stroke:#000,stroke-width:1px, color:#000
    style Slot9D fill:#ffff66,stroke:#000,stroke-width:1px, color:#000
```

## YubiKey Setup Process

```mermaid
flowchart TD
    Start[Start Setup] --> ResetPIV[Reset PIV Application]
    ResetPIV --> ChangePIN[Change Default PIN]
    ChangePIN --> ChangePUK[Change Default PUK]
    ChangePUK --> ChangeMGMT[Change Management Key]
    ChangeMGMT --> GenKey[Generate Private Key in Slot]
    GenKey --> GenCert[Generate Self-Signed Certificate]
    GenCert --> Export[Export Certificate]
    Export --> End[Setup Complete]
    
    style Start fill:#f9f9f9,stroke:#333,stroke-width:1px, color: #000
    style End fill:#f9f9f9,stroke:#333,stroke-width:1px, color: #000
    style GenKey fill:#ddffdd,stroke:#333,stroke-width:2px, color: #000
    style GenCert fill:#ddffdd,stroke:#333,stroke-width:2px, color: #000
```

## YubiKey PIN Protection

```mermaid
stateDiagram-v2
    [*] --> Ready: YubiKey Inserted
    
    Ready --> WaitingForPIN: Sign Request
    WaitingForPIN --> Authenticated: Correct PIN
    WaitingForPIN --> Retry: Wrong PIN
    Authenticated --> Signing: Perform Operation
    Signing --> Ready: Operation Complete
    
    Retry --> WaitingForPIN: Attempts Remaining
    Retry --> Locked: Too Many Failed Attempts
    
    Locked --> WaitingForPUK: PUK Required
    WaitingForPUK --> PINReset: Correct PUK
    WaitingForPUK --> PUKRetry: Wrong PUK
    
    PINReset --> Ready: New PIN Set
    PUKRetry --> WaitingForPUK: Attempts Remaining
    PUKRetry --> Blocked: Too Many Failed PUK Attempts
    
    Blocked --> [*]: Factory Reset Required
    
    note right of WaitingForPIN
        Default: 3 attempts
    end note
    
    note right of WaitingForPUK
        Default: 8 attempts
    end note
```

## Certificate Signing with YubiKey

```mermaid
sequenceDiagram
    participant App as PiCA Application
    participant PCSC as PC/SC Daemon
    participant YK as YubiKey
    
    App->>App: Generate signing request
    App->>PCSC: Forward to YubiKey
    PCSC->>YK: Request PIN
    YK-->>PCSC: Prompt for PIN
    PCSC-->>App: Request PIN from user
    App->>App: Collect PIN
    App->>PCSC: Submit PIN
    PCSC->>YK: Forward PIN
    YK->>YK: Validate PIN
    
    alt PIN Valid
        YK->>YK: Unlock private key
        YK->>YK: Perform signing operation
        YK-->>PCSC: Return signature
        PCSC-->>App: Forward signature
        App->>App: Create certificate with signature
    else PIN Invalid
        YK-->>PCSC: Return error
        PCSC-->>App: Forward error
        App->>App: Display error to user
    end
```

## YubiKey States and Transitions

```mermaid
stateDiagram-v2
    [*] --> Uninitialized: New or Reset YubiKey
    
    Uninitialized --> Configured: Setup Process
    Configured --> Ready: YubiKey Ready for Use
    
    Ready --> InUse: Signing Operations
    InUse --> Ready: Operation Complete
    
    Ready --> Locked: Failed PIN Attempts
    Locked --> Ready: PUK Unlock
    
    Locked --> Blocked: Failed PUK Attempts
    Blocked --> Uninitialized: Factory Reset
    
    Configured --> Compromised: Security Incident
    Compromised --> Decommissioned: Remove from Service
    
    Ready --> Decommissioned: End of Life
    Decommissioned --> [*]
```

## YubiKey Management Workflow

```mermaid
flowchart TD
    Inventory[Inventory Management] --> Setup
    Setup[YubiKey Setup] --> Deploy
    Deploy[Deploy to CA System] --> Use
    Use[Operational Use] --> Rotate
    Rotate[Key/PIN Rotation] --> Use
    
    Use --> Incident{Security Incident?}
    Incident -->|Yes| Compromise[Handle Compromise]
    Incident -->|No| Continue
    
    Continue --> EOL{End of Life?}
    EOL -->|Yes| Decom[Decommission]
    EOL -->|No| Use
    
    Compromise --> Decom
    Decom --> Archive[Archive Information]
    Archive --> NewYK[New YubiKey]
    NewYK --> Inventory
    
    style Incident fill:#ffdddd,stroke:#333,stroke-width:2px, color: #000
    style Compromise fill:#ffdddd,stroke:#333,stroke-width:2px, color: #000
    style Inventory fill:#ddffdd,stroke:#333,stroke-width:1px, color: #000
    style Setup fill:#ddffdd,stroke:#333,stroke-width:1px, color: #000
    style Archive fill:#ffffdd,stroke:#333,stroke-width:1px, color: #000
```

## YubiKey Touch Policy Options

```mermaid
graph TD
    YK[YubiKey] --- Never[Never</br>No Touch Required]
    YK --- Always[Always</br>Touch Required Each Time]
    YK --- Cached[Cached</br>Touch Required Once</br>Then Cached for 15 Seconds]
    
    style YK fill:#f9f9f9,stroke:#333,stroke-width:2px, color:#000
    style Never fill:#ffdddd,stroke:#333,stroke-width:1px, color:#000
    style Always fill:#ddffdd,stroke:#333,stroke-width:1px, color:#000
    style Cached fill:#ddddff,stroke:#333,stroke-width:1px, color:#000
    
    Never -.-> |Least Secure</br>No Physical Interaction| SecurityLow[Low Security]
    Always -.-> |Most Secure</br>Requires Physical Presence| SecurityHigh[High Security]
    Cached -.-> |Balance of Security</br>and Convenience| SecurityMed[Medium Security]
```

## YubiKey Backup and Recovery Strategy

```mermaid
flowchart TD
    Primary[Primary YubiKey] --> |Disaster</br>Recovery| Secondary[Secondary YubiKey]
    Primary --> |Regular</br>Operation| Normal[Normal CA Operations]
    
    Secondary --> |Activation</br>Required| Activate[Activate Backup YubiKey]
    Activate --> NewPrimary[New Primary Operations]
    
    subgraph Secure Storage
        Secondary
        RecoveryKit[Recovery Documentation]
        PINs[PINs and PUKs]
        MGMKeys[Management Keys]
    end
    
    Secondary -.-> RecoveryKit
    RecoveryKit -.-> PINs
    RecoveryKit -.-> MGMKeys
    
    style Primary fill:#ddffdd,stroke:#333,stroke-width:2px, color:#000
    style Secondary fill:#ffdddd,stroke:#333,stroke-width:2px, color:#000
    style Secure Storage fill:#ffffdd,stroke:#333,stroke-width:1px,stroke-dasharray: 5 5, color:#000
```

## PIN Protection for Management Key

```mermaid
sequenceDiagram
    participant Admin
    participant App as PiCA Application
    participant YK as YubiKey
    
    Admin->>App: Request management operation
    App->>YK: Request with protected management key
    YK->>YK: Check if PIN protection enabled
    YK-->>App: Request PIN
    App->>Admin: Prompt for PIN
    Admin->>App: Provide PIN
    App->>YK: Submit PIN
    
    alt PIN Valid
        YK->>YK: Use management key for operation
        YK->>YK: Perform administrative function
        YK-->>App: Return success
        App-->>Admin: Display success
    else PIN Invalid
        YK-->>App: Return error
        App-->>Admin: Display error
    end
    
    Note over Admin,YK: Management key never exposed, protected by PIN
```

## Hardware-Based Attestation Process

```mermaid
sequenceDiagram
    participant Admin
    participant App as PiCA Application
    participant YK as YubiKey
    participant Verify as Verification System
    
    Admin->>App: Request key generation with attestation
    App->>YK: Generate key with attestation option
    YK->>YK: Generate key in requested slot
    YK->>YK: Sign attestation certificate using device key
    YK-->>App: Return public key and attestation
    App-->>Admin: Provide attestation certificate
    
    Admin->>Verify: Submit for verification
    Verify->>Verify: Validate against YubiKey CA certs
    Verify->>Verify: Verify attestation signature
    Verify-->>Admin: Confirm key generated on hardware
    
    Note over YK,Verify: Proves key was generated on authentic YubiKey
```

## YubiKey in PiCA Software Architecture

```mermaid
flowchart TD
    subgraph PiCA
        UI[UI Component] --> CA[CA Module]
        WebAPI[Web API] --> CA
        CA --> YKInteg[YubiKey Integration]
    end
    
    YKInteg --> PCSC[PC/SC Layer]
    PCSC --> YKLib[YubiKey Libraries]
    YKLib --> Driver[USB Driver]
    Driver --> YK[YubiKey Hardware]
    
    style PiCA fill:#f0f0ff,stroke:#333,stroke-width:2px
    style YKInteg fill:#ddffdd,stroke:#333,stroke-width:2px, color:#000
    style YK fill:#ffdddd,stroke:#333,stroke-width:2px, color:#000
```

These diagrams provide comprehensive visualizations of how YubiKeys are integrated and used within the PiCA Certificate Authority system, from initial setup to operational use and security mechanisms.
