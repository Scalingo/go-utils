package gomockgenerator

import (
	"strings"
	"testing"
)

func TestInterfaceSignatureResolvesStandardLibraryPackage(t *testing.T) {
	sig, err := interfaceSignature(t.Context(), "net", "Conn")
	if err != nil {
		t.Fatalf("interfaceSignature returned error: %v", err)
	}

	if sig == "" {
		t.Fatal("interfaceSignature returned an empty signature")
	}

	for _, method := range []string{"Read", "Write", "Close", "SetDeadline"} {
		if !strings.Contains(sig, method) {
			t.Fatalf("signature %q does not contain method %q", sig, method)
		}
	}
}
