package s3pp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type Config struct {
	AWSCredentials AWSCredentials
	Bucket         string
	Region         string
	Expires        time.Duration
	Key            Condition
	Conditions     []Condition
}

type AWSCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
}

type Form struct {
	Action string `json:"action"`

	// name -> value, generated from the Name and Value of the provided Conditions
	Fields map[string]string `json:"fields"`
}

// New encodes and signs the POST Policy using AWS Signature v4 and returns the required
// fields to create a form for the POST request.
func New(c Config) (*Form, error) {
	date := time.Now().UTC()
	signing := signingKey(c.AWSCredentials.SecretAccessKey, date, c.Region, "s3")
	scope := scope(c.AWSCredentials.AccessKeyID, date, c.Region, "s3")
	conditions := append(
		c.Conditions,
		c.Key,
		Match("bucket", c.Bucket),
		Match("x-amz-algorithm", "AWS4-HMAC-SHA256"),
		Match("x-amz-credential", scope),
		Match("x-amz-date", date.Format("20060102T150405Z")),
	)

	b, err := json.Marshal(map[string]interface{}{
		"expiration": date.Add(c.Expires),
		"conditions": conditions,
	})
	if err != nil {
		return nil, err
	}

	policy := base64.StdEncoding.EncodeToString(b)
	fields := make(map[string]string)
	for _, c := range conditions {
		fields[c.Name()] = c.Value()
	}
	fields["policy"] = policy
	fields["x-amz-signature"] = signPolicy(policy, signing)

	return &Form{Action: bucketURL(c.Bucket), Fields: fields}, nil
}

func bucketURL(bucket string) string {
	return fmt.Sprintf("https://%s.s3.amazonaws.com/", bucket)
}

func scope(accessKeyID string, date time.Time, region, service string) string {
	return fmt.Sprintf(
		"%s/%s/%s/%s/aws4_request",
		accessKeyID,
		date.Format("20060102"),
		region,
		service,
	)
}

// signingKey derives a key from your AWS secret access key and the given scope.
// see: http://docs.aws.amazon.com/general/latest/gr/sigv4-calculate-signature.html
func signingKey(secret string, date time.Time, region, service string) []byte {
	kDate := sign(date.Format("20060102"), []byte("AWS4"+secret))
	kRegion := sign(region, kDate)
	kService := sign(service, kRegion)
	return sign("aws4_request", kService)
}

func signPolicy(policy string, key []byte) string {
	return hex.EncodeToString(sign(policy, key))
}

func sign(msg string, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(msg))
	return mac.Sum(nil)
}
