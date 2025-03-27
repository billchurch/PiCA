package api

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/billchurch/pica/internal/ca"
	"github.com/billchurch/pica/internal/ca/commands"
	"github.com/billchurch/pica/internal/crypto"
	"github.com/billchurch/pica/internal/yubikey"
)

// Server represents the API server
type Server struct {
	CA          *ca.CA
	YubiKeySlot yubikey.PIVSlot
	CertDir     string
	CSRDir      string
}

// NewServer creates a new API server
func NewServer(ca *ca.CA, slot yubikey.PIVSlot, certDir, csrDir string) *Server {
	return &Server{
		CA:          ca,
		YubiKeySlot: slot,
		CertDir:     certDir,
		CSRDir:      csrDir,
	}
}

// StartServer starts the API server
func (s *Server) StartServer(addr string) error {
	// Ensure directories exist
	if err := os.MkdirAll(s.CertDir, 0755); err != nil {
		return fmt.Errorf("error creating certificate directory: %w", err)
	}
	if err := os.MkdirAll(s.CSRDir, 0755); err != nil {
		return fmt.Errorf("error creating CSR directory: %w", err)
	}

	// Set up routes
	http.HandleFunc("/api/health", s.handleHealth)
	http.HandleFunc("/api/submit-csr", s.handleSubmitCSR)
	http.HandleFunc("/api/certificates", s.handleListCertificates)
	http.HandleFunc("/api/certificate/", s.handleGetCertificate)
	http.HandleFunc("/api/revoke", s.handleRevokeCertificate)

	// Start the server
	log.Printf("Starting API server on %s", addr)
	return http.ListenAndServe(addr, nil)
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// CSRRequest represents a CSR submission
type CSRRequest struct {
	CSR     string `json:"csr"`
	Profile string `json:"profile"`
}

// handleSubmitCSR handles CSR submission
func (s *Server) handleSubmitCSR(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req CSRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing request: %s", err), http.StatusBadRequest)
		return
	}

	// Validate CSR
	block, _ := pem.Decode([]byte(req.CSR))
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		http.Error(w, "Invalid CSR PEM data", http.StatusBadRequest)
		return
	}

	// Parse CSR to get subject
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing CSR: %s", err), http.StatusBadRequest)
		return
	}

	// Verify CSR signature
	if err := csr.CheckSignature(); err != nil {
		http.Error(w, fmt.Sprintf("CSR signature verification failed: %s", err), http.StatusBadRequest)
		return
	}

	// Save CSR to file
	csrFilename := fmt.Sprintf("%s.csr", csr.Subject.CommonName)
	csrPath := filepath.Join(s.CSRDir, csrFilename)
	if err := os.WriteFile(csrPath, []byte(req.CSR), 0644); err != nil {
		http.Error(w, fmt.Sprintf("Error saving CSR: %s", err), http.StatusInternalServerError)
		return
	}

	// Generate certificate path
	certFilename := fmt.Sprintf("%s.crt", csr.Subject.CommonName)
	certPath := filepath.Join(s.CertDir, certFilename)

	// Sign the certificate
	cmd := commands.NewSignCommand(
		s.CA,
		csrPath,
		certPath,
		req.Profile,
		crypto.FromYubiKeySlot(s.YubiKeySlot),
	)

	if err := cmd.Execute(); err != nil {
		http.Error(w, fmt.Sprintf("Error signing certificate: %s", err), http.StatusInternalServerError)
		return
	}

	// Return the certificate
	certData, err := os.ReadFile(certPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading certificate: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"certificate": string(certData),
	})
}

// CertificateInfo represents certificate information
type CertificateInfo struct {
	Subject      string `json:"subject"`
	SerialNumber string `json:"serialNumber"`
	NotBefore    string `json:"notBefore"`
	NotAfter     string `json:"notAfter"`
	Status       string `json:"status"`
}

// handleListCertificates handles certificate listing
func (s *Server) handleListCertificates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read certificates from the certificate directory
	files, err := os.ReadDir(s.CertDir)
	if err != nil {
		log.Printf("Error reading certificate directory: %v", err)
		http.Error(w, "Failed to list certificates", http.StatusInternalServerError)
		return
	}

	certs := []CertificateInfo{}

	// Process each file in the directory
	for _, file := range files {
		fileName := file.Name()
		
		// Skip directories and non-certificate files
		if file.IsDir() || !(filepath.Ext(fileName) == ".pem" || filepath.Ext(fileName) == ".crt") {
			continue
		}

		// Read and parse certificate
		certPath := filepath.Join(s.CertDir, fileName)
		certData, err := os.ReadFile(certPath)
		if err != nil {
			log.Printf("Error reading certificate file %s: %v", fileName, err)
			continue
		}

		// Decode PEM block
		block, _ := pem.Decode(certData)
		if block == nil || block.Type != "CERTIFICATE" {
			log.Printf("Failed to decode PEM block from %s", fileName)
			continue
		}

		// Parse certificate
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			log.Printf("Failed to parse certificate from %s: %v", fileName, err)
			continue
		}

		// Create certificate info
		certInfo := CertificateInfo{
			Subject:      cert.Subject.CommonName,
			SerialNumber: fmt.Sprintf("%X", cert.SerialNumber),
			NotBefore:    cert.NotBefore.Format("2006-01-02"),
			NotAfter:     cert.NotAfter.Format("2006-01-02"),
			Status:       "Valid", // Default to valid as we don't check CRL in this example
		}

		certs = append(certs, certInfo)
	}

	log.Printf("Found %d certificates in %s", len(certs), s.CertDir)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"certificates": certs,
	})
}

// handleGetCertificate handles certificate retrieval
func (s *Server) handleGetCertificate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract serial number from path
	serialNumber := r.URL.Path[len("/api/certificate/"):]
	if serialNumber == "" {
		http.Error(w, "Serial number required", http.StatusBadRequest)
		return
	}

	// Read all certificates and find the one with the matching serial number
	files, err := os.ReadDir(s.CertDir)
	if err != nil {
		log.Printf("Error reading certificate directory: %v", err)
		http.Error(w, "Failed to access certificates", http.StatusInternalServerError)
		return
	}

	var matchingCertPath string
	var matchingCert *x509.Certificate

	// Look for a certificate with the matching serial number
	for _, file := range files {
		if file.IsDir() || !(filepath.Ext(file.Name()) == ".pem" || filepath.Ext(file.Name()) == ".crt") {
			continue
		}

		certPath := filepath.Join(s.CertDir, file.Name())
		certData, err := os.ReadFile(certPath)
		if err != nil {
			continue
		}

		block, _ := pem.Decode(certData)
		if block == nil || block.Type != "CERTIFICATE" {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			continue
		}

		// Compare serial number (ignore case and allow for different formats)
		certSerial := fmt.Sprintf("%X", cert.SerialNumber)
		if strings.EqualFold(certSerial, serialNumber) {
			matchingCertPath = certPath
			matchingCert = cert
			break
		}
	}

	if matchingCert == nil {
		http.Error(w, "Certificate not found", http.StatusNotFound)
		return
	}

	// Read the certificate file to return the PEM data
	certData, err := os.ReadFile(matchingCertPath)
	if err != nil {
		log.Printf("Error reading certificate file: %v", err)
		http.Error(w, "Failed to read certificate", http.StatusInternalServerError)
		return
	}

	// Create certificate info
	certInfo := CertificateInfo{
		Subject:      matchingCert.Subject.CommonName,
		SerialNumber: fmt.Sprintf("%X", matchingCert.SerialNumber),
		NotBefore:    matchingCert.NotBefore.Format("2006-01-02"),
		NotAfter:     matchingCert.NotAfter.Format("2006-01-02"),
		Status:       "Valid", // Default to valid as we don't check CRL in this example
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"certificate": string(certData),
		"info":        certInfo,
	})
}

// RevokeRequest represents a certificate revocation request
type RevokeRequest struct {
	SerialNumber string `json:"serialNumber"`
	Reason       string `json:"reason"`
}

// handleRevokeCertificate handles certificate revocation
func (s *Server) handleRevokeCertificate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req RevokeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing request: %s", err), http.StatusBadRequest)
		return
	}

	// Revoke the certificate
	cmd := commands.NewRevokeCommand(
		s.CA,
		req.SerialNumber,
		req.Reason,
		crypto.FromYubiKeySlot(s.YubiKeySlot),
	)

	if err := cmd.Execute(); err != nil {
		http.Error(w, fmt.Sprintf("Error revoking certificate: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": "Certificate revoked successfully",
	})
}
