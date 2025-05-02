# TODO: 
- have all these fields in Implant: ? 

```json
type Implant struct {
  gorm.Model
  // 	LastSeen time.Time // Timestamp of last beacon or activity
  // PublicKey      []byte    // If using asymmetric encryption
  // NonceCounter   uint64    // To prevent replay attacks
  // SessionToken   string    // JWT or session ID
  PendingTasks   []Task // Pending tasks for this implant
CompletedTasks []Task // History or recent completed tasks
//	Labels         []string  // Optional: e.g., ["phishing", "test", "persistent"]
// UserAgent string // Optional: helps with fingerprinting
}
```
