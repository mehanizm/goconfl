Go Confluence
================

Minimalistic confluence api wrapper for confluence

Quick start:

```go
w, err := NewWiki(confluenceURL, BasicAuth(confluenceLogin, confluencePass))
if err != nil {
	panic(err)
}
content, err := w.GetContentByID(pageID, []string{"body.storage"})
if err != nil {
    panic(err)
}
fmt.Println(content.Body.Storage.Value)
```