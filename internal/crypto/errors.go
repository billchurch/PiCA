package crypto

import (
	"errors"
)

// Error definitions for the crypto package
var (
	// ErrProviderNotFound is returned when a requested provider type is not registered
	ErrProviderNotFound = errors.New("crypto provider not found")
	
	// ErrNotConnected is returned when an operation is attempted on a provider that is not connected
	ErrNotConnected = errors.New("crypto provider not connected")
	
	// ErrSlotNotFound is returned when the requested slot does not exist
	ErrSlotNotFound = errors.New("slot not found")
	
	// ErrKeyNotFound is returned when a key is not found in the requested slot
	ErrKeyNotFound = errors.New("key not found in slot")
	
	// ErrCertNotFound is returned when a certificate is not found in the requested slot
	ErrCertNotFound = errors.New("certificate not found in slot")
	
	// ErrInvalidKeyType is returned when an invalid key type is specified
	ErrInvalidKeyType = errors.New("invalid key type")
	
	// ErrInvalidCertificate is returned when an invalid certificate is provided
	ErrInvalidCertificate = errors.New("invalid certificate")
	
	// ErrOperationNotSupported is returned when an operation is not supported by the provider
	ErrOperationNotSupported = errors.New("operation not supported by this provider")
	
	// ErrInvalidAlgorithm is returned when an invalid algorithm is specified
	ErrInvalidAlgorithm = errors.New("invalid algorithm")
)
