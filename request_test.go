package main

import "testing"

func TestScanHost(t *testing.T) {
	scan := ScanHost("https://yahoo.com")

	if scan.GetHeaders("X-Frame-Options")[0] != "DENY" {
		t.Error("Failed to correctly scan yahoo")
	}

	scan2 := ScanHost("httpsss://fake.com")
	if scan2.Error == "" {
		t.Error("Failed to update error on failure")
	}
}
