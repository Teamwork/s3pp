package s3pp

import (
	"encoding/hex"
	"testing"
	"time"
)

// http://docs.aws.amazon.com/general/latest/gr/signature-v4-examples.html
func TestSigningKey(t *testing.T) {
	secret := "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY"
	date := time.Date(2012, 2, 15, 0, 0, 0, 0, time.UTC)
	region := "us-east-1"
	service := "iam"

	key := hex.EncodeToString(signingKey(secret, date, region, service))
	want := "f4780e2d9f65fa895f9c67b32ce1baf0b0d8a43505a000a1a9e090d414db404d"
	if key != want {
		t.Errorf("got %s, want %s", key, want)
	}
}

// http://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-post-example.html
func TestSigningPolicy(t *testing.T) {
	secret := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	date := time.Date(2013, 8, 6, 0, 0, 0, 0, time.UTC)
	region := "us-east-1"
	service := "s3"
	policy := "eyAiZXhwaXJhdGlvbiI6ICIyMDEzLTA4LTA3VDEyOjAwOjAwLjAwMFoiLA0KICAiY29uZGl0aW9ucyI6IFsNCiAgICB7ImJ1Y2tldCI6ICJleGFtcGxlYnVja2V0In0sDQogICAgWyJzdGFydHMtd2l0aCIsICIka2V5IiwgInVzZXIvdXNlcjEvIl0sDQogICAgeyJhY2wiOiAicHVibGljLXJlYWQifSwNCiAgICB7InN1Y2Nlc3NfYWN0aW9uX3JlZGlyZWN0IjogImh0dHA6Ly9leGFtcGxlYnVja2V0LnMzLmFtYXpvbmF3cy5jb20vc3VjY2Vzc2Z1bF91cGxvYWQuaHRtbCJ9LA0KICAgIFsic3RhcnRzLXdpdGgiLCAiJENvbnRlbnQtVHlwZSIsICJpbWFnZS8iXSwNCiAgICB7IngtYW16LW1ldGEtdXVpZCI6ICIxNDM2NTEyMzY1MTI3NCJ9LA0KICAgIFsic3RhcnRzLXdpdGgiLCAiJHgtYW16LW1ldGEtdGFnIiwgIiJdLA0KDQogICAgeyJ4LWFtei1jcmVkZW50aWFsIjogIkFLSUFJT1NGT0ROTjdFWEFNUExFLzIwMTMwODA2L3VzLWVhc3QtMS9zMy9hd3M0X3JlcXVlc3QifSwNCiAgICB7IngtYW16LWFsZ29yaXRobSI6ICJBV1M0LUhNQUMtU0hBMjU2In0sDQogICAgeyJ4LWFtei1kYXRlIjogIjIwMTMwODA2VDAwMDAwMFoiIH0NCiAgXQ0KfQ=="

	signed := signPolicy(policy, signingKey(secret, date, region, service))
	want := "21496b44de44ccb73d545f1a995c68214c9cb0d41c45a17a5daeec0b1a6db047"
	if signed != want {
		t.Errorf("got %s, want %s", signed, want)
	}
}

func TestScope(t *testing.T) {
	date := time.Date(2016, 1, 30, 0, 0, 0, 0, time.UTC)
	s := scope("ACCESSKEY", date, "us-east-1", "s3")
	want := "ACCESSKEY/20160130/us-east-1/s3/aws4_request"
	if s != want {
		t.Errorf("got %s, want %s", s, want)
	}
}

func TestConditions(t *testing.T) {
	cases := []struct {
		cond        Condition
		name, value string
	}{
		{Any("key"), "key", ""},
		{StartsWith("key", "attachments/"), "key", "attachments/"},
		{Match("key", "file.zip"), "key", "file.zip"},
		{ContentLengthRange(1, 1000), "content-length-range", ""},
	}
	for _, c := range cases {
		if c.cond.Name() != c.name {
			t.Errorf("Name() == %q, want %q", c.cond.Name(), c.name)
		}
		if c.cond.Value() != c.value {
			t.Errorf("Value() == %q, want %q", c.cond.Value(), c.value)
		}
	}
}

func TestGeneratedConditonsAddedToFields(t *testing.T) {
	form, err := New(Config{Key: Any("key")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"x-amz-algorithm", "x-amz-credential", "x-amz-date"}
	for _, field := range expected {
		if _, ok := form.Fields[field]; !ok {
			t.Errorf("Fields missing expected key %q", field)
		}
	}
}

func TestContentLengthRangeExcludedFromFields(t *testing.T) {
	form, err := New(Config{
		Key:        Any("key"),
		Conditions: []Condition{ContentLengthRange(0, 100)},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := form.Fields["content-length-range"]; ok {
		t.Errorf("content-length-range shouldn't be included in the form fields")
	}
}
