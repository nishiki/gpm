package gpm

import "testing"

func TestCreateEmptyEntry(t *testing.T) {
  var entry Entry
  err := entry.Verify()
  if err == nil {
    t.Error("an entry Without an ID must return an error")
  }
}

func TestCreateEntryWithoutName(t *testing.T) {
  var entry Entry
  entry.GenerateID()
  if entry.ID == "" {
    t.Error("generateID can't be generate a void ID")
  }

  err := entry.Verify()
  if err == nil {
    t.Error("an entry without a name must return an error")
  }
}

func TestCreateEntryWithName(t *testing.T) {
  entry := Entry{ Name: "test" }
  entry.GenerateID()
  err := entry.Verify()
  if err != nil {
    t.Errorf("an entry with a name mustn't return an error: %s", err)
  }
}

func TestCreateEntryWithBadURI(t *testing.T) {
  entry := Entry{ Name: "test", URI: "url/bad:" }
  entry.GenerateID()
  err := entry.Verify()
  if err == nil {
    t.Error("an entry with a bad URI must return an error")
  }
}

func TestCreateEntryWithGoodURI(t *testing.T) {
  entry := Entry{ Name: "test", URI: "http://localhost:8081" }
  entry.GenerateID()
  err := entry.Verify()
  if err != nil {
    t.Errorf("an entry with a good URI mustn't return an error: %s", err)
  }
}

func TestGenerateOTPCode(t *testing.T) {
  entry := Entry{ OTP: "JBSWY3DPEHPK3PXP" }
  code, time, err := entry.OTPCode()
  if err != nil {
    t.Errorf("must generate an OTP code without error: %s", err)
  }
  if len(code) != 6 {
    t.Errorf("must generate an OTP code with 6 chars: %s", code)
  }
  if time < 0 || time > 30 {
    t.Errorf("time must be between 0 and 30: %d", time)
  }
}
