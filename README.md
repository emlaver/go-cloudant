# go-cloudant
A Cloudant library for Golang.

[![Build Status](https://travis.ibm.com/cloudant/go-cloudant.svg?token=wExLyHssVGfwxehwRf9S&branch=master)](https://travis.ibm.com/cloudant/go-cloudant)

_The API is not fully baked at this time and may change._

## Description
A [Cloudant](https://cloudant.com/) library for Golang.

## Installation
```bash
go get github.ibm.com/cloudant/go-cloudant
```

## Supported Features
- Session authentication
- Keep-Alive & Connection Pooling
- Configurable request retrying
- Hard limit on request concurrency
- Stream `/_all_docs` & `/_changes`
- Manage `_bulk_docs` uploads

## Getting Started

### Creating a Cloudant client:
```go
// create a Cloudant client (max. request concurrency 5) with default retry configuration:
//   - maximum retries per request:     3
//   - random retry delay minimum:      5  seconds
//   - random retry delay maximum:      30 seconds
client1, err1 := cloudant.CreateClient("user123", "pa55w0rd01", "https://user123.cloudant.com", 5)

// create a Cloudant client (max. request concurrency 20) with _custom_ retry configuration:
//   - maximum retries per request:     5
//   - random retry delay minimum:      10  seconds
//   - random retry delay maximum:      60 seconds
client2, err2 := cloudant.CreateClientWithRetry("user123", "pa55w0rd01", "https://user123.cloudant.com", 20, 5, 10, 60)
```

### `Get` a document:
```go
// create a Cloudant client (max. request concurrency 5)
client, err := cloudant.CreateClient("user123", "pa55w0rd01", "https://user123.cloudant.com", 5)
db, err := client.GetOrCreate("my_database")

type Doc struct {
    Id     string    `json:"_id"`
    Rev    string    `json:"_rev"`
    Foo    string    `json:"foo"`
}

doc = new(Doc)
err = db.Get("my_doc", doc)

fmt.Println(doc.Foo)  // prints 'foo' key
```

### `Set` a document:
```go
// create a Cloudant client (max. request concurrency 5)
client, err := cloudant.CreateClient("user123", "pa55w0rd01", "https://user123.cloudant.com", 5)
db, err := client.GetOrCreate("my_database")

myDoc := &Doc{
        Id:     "my_doc_id",
        Rev:    "2-xxxxxxx",
        Foo:    "bar",
}

newRev, err := db.Set(myDoc)

fmt.Println(newRev)  // prints '_rev' of new document revision
```

### `Delete` a document:
```go
// create a Cloudant client (max. request concurrency 5)
client, err := cloudant.CreateClient("user123", "pa55w0rd01", "https://user123.cloudant.com", 5)
db, err := client.GetOrCreate("my_database")

err := db.Delete("my_doc_id", "2-xxxxxxx")
```

### Using `_bulk_docs`:
```go
// create a Cloudant client (max. request concurrency 5)
client, err := cloudant.CreateClient("user123", "pa55w0rd01", "https://user123.cloudant.com", 5)
db, err := client.GetOrCreate("my_database")

myDoc1 := Doc{
        Id:     "doc1",
        Rev:    "1-xxxxxxx",
        Foo:    "bar",
}

myDoc2 := Doc{
        Id:     "doc2",
        Rev:    "2-xxxxxxx",
        Foo:    "bar",
}

myDoc3 := Doc{
        Id:     "doc3",
        Rev:    "3-xxxxxxx",
        Foo:    "bar",
}

uploader := db.Bulk(50) // new uploader using batch size 50

uploader.FireAndForget(myDoc1)

upload.Flush() // uploads all received documents

r2 := uploader.UploadNow(myDoc2) // uploaded as soon as it's received by a worker

r2.Wait()
if r2.Error != nil {
    fmt.Println(r2.Response.Id, r2.Response.Rev) // prints new document '_id' and 'rev'
}

r3 := uploader.Upload(myDoc3) // queues until the worker creates a full batch of 50 documents

upload.Stop() // uploads any queued documents before stopping

r3.Wait()
if r3.Error != nil {
    fmt.Println(r3.Response.Id, r3.Response.Rev) // prints new document '_id' and 'rev'
}
```

### Using `/_all_docs`:
```go
// create a Cloudant client (max. request concurrency 5)
client, err := cloudant.CreateClient("user123", "pa55w0rd01", "https://user123.cloudant.com", 5)
db, err := client.GetOrCreate("my_database")

docs, err := db.All()

// OR include some query options...
//
// q := &cloudant.AllQuery{
//     Limit:	    123,
//     StartKey:    "bar",
//     EndKey:      "foo",
// }
// docs, err := db.AllQ(q)

for{
    doc, more := <-docs
	if more {
	    fmt.Println(doc.Id, doc.Rev)  // prints document 'id' and 'rev'
	} else {
	    break
	}
}
```

### Using `/_changes`:
```go
// create a Cloudant client (max. request concurrency 5)
client, err := cloudant.CreateClient("user123", "pa55w0rd01", "https://user123.cloudant.com", 5)
db, err := client.GetOrCreate("my_database")

changes, err := db.Changes()

for{
    change, more := <-changes
	if more {
	    fmt.Println(doc.Seq, doc.Id, doc.Rev)  // prints change 'seq', 'id' and 'rev'
	} else {
	    break
	}
}
```
