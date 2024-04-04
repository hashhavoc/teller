package utils

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/btcsuite/btcd/btcec/v2"
)

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

// StacksPrivateKey represents a private key with a flag indicating if it should generate compressed public keys.
type StacksPrivateKey struct {
	PrivateKey *btcec.PrivateKey
	Compressed bool
}

// makeRandomPrivKey generates a new random StacksPrivateKey.
func MakeRandomPrivKey() (*StacksPrivateKey, error) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, err
	}

	// Assuming you want to generate compressed public keys by default.
	// Adjust based on your needs.
	return &StacksPrivateKey{
		PrivateKey: privKey,
		Compressed: true,
	}, nil
}

// // Automatically detect and decode a hex string.
// // It returns the decoded string or an error if the decoding fails.
// func autoDetectAndDecodeHexString(hexInput string) (string, error) {
// 	// Check if the input is potentially a uint128 encoded value
// 	if len(hexInput) >= 6 && hexInput[:2] == "0x" {
// 		// Attempt to decode as uint128
// 		u, err := uint128.FromString(hexInput[6:])
// 		if err == nil {
// 			// Successfully decoded as uint128
// 			return u.String(), nil
// 		}
// 		// If decoding as uint128 fails, fall through to string decoding
// 	}

// 	// Default to assuming the input is a string encoded value
// 	if len(hexInput) < 14 {
// 		return "", errors.New("invalid hex string length for string decoding")
// 	}
// 	decoded, err := hex.DecodeString(hexInput[14:])
// 	if err != nil {
// 		return "", fmt.Errorf("failed to decode hex string: %v", err)
// 	}
// 	return string(decoded), nil
// }

// func main() {
// 	// Example hex inputs
// 	examples := []string{
// 		"0x070d00000009486972655669626573",
// 		"0x070d000000055649424553",
// 		"0x070100000000000000000000000000000008",
// 	}

// 	for _, example := range examples {
// 		decoded, err := autoDetectAndDecodeHexString(example)
// 		if err != nil {
// 			fmt.Printf("Error decoding: %v\n", err)
// 			continue
// 		}
// 		fmt.Printf("Decoded value: %s\n", decoded)
// 	}
// }
