# s3pp

A package to help you create [POST policies](http://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-HTTPPOSTConstructPolicy.html) to upload files directly to Amazon S3, see the [AWS docs](http://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-UsingHTTPPOST.html) for the all the available conditions and values.

## Example

A working example is available in [examples/form.go](examples/form.go), below are the relevant parts:
```go
form, err := s3pp.New(s3pp.Config{
        AWSCredentials: s3pp.AWSCredentials{
                AccessKeyID:     "key",
                SecretAccessKey: "secret",
        },
        Bucket:  "mybucket",
        Region:  "us-east-1",
        Expires: 10 * time.Minute,
        Key:     s3pp.Match("key", uuid.New()),
        Conditions: []s3pp.Condition{
                s3pp.Match("acl", "public-read"),
                s3pp.Match("success_action_status", "201"),
        },
})
```
`form.Fields` will contain all the fields generated from the conditions for the form and you can pass any additional conditions in `Conditions`. Available conditions are documented here: [Creating a POST Policy](http://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-HTTPPOSTConstructPolicy.html).

The form:
```html
<form action="{{.Action}}" method="POST" enctype="multipart/form-data">
  {{range $name, $value := .Fields}}
    <input type="text" name="{{$name}}" value="{{$value}}" readonly>
  {{end}}

  <input type="file" name="file">
  <input type="submit" value="Upload">
</form>
```
