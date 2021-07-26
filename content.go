package goconfl

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type Content struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Title  string `json:"title"`
	Body   struct {
		Storage struct {
			Value          string `json:"value"`
			Representation string `json:"representation"`
		} `json:"storage"`
	} `json:"body"`
	Version struct {
		Number int `json:"number"`
	} `json:"version"`
}

type SearchResults struct {
	Results []*Content `json:"results"`
}

func (w *Wiki) contentEndpoint(pageID string) (*url.URL, error) {
	return url.ParseRequestURI(w.endPoint.String() + "/content/" + pageID)
}

func (w *Wiki) DeleteContentByID(pageID string) error {
	contentEndPoint, err := w.contentEndpoint(pageID)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", contentEndPoint.String(), nil)
	if err != nil {
		return err
	}

	_, err = w.sendRequest(req)
	if err != nil {
		return err
	}
	return nil
}

func (w *Wiki) GetContentByID(pageID string, expand []string) (*Content, error) {
	contentEndPoint, err := w.contentEndpoint(pageID)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set("expand", strings.Join(expand, ","))
	contentEndPoint.RawQuery = data.Encode()

	req, err := http.NewRequest("GET", contentEndPoint.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := w.sendRequest(req)
	if err != nil {
		return nil, err
	}
	// // Save a copy of this request for debugging.
	// requestDump, err := httputil.DumpRequest(req, true)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(requestDump))

	content := new(Content)
	err = json.Unmarshal(res, content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (w *Wiki) UpdateContentByID(content *Content) (*Content, error) {
	jsonbody, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	contentEndPoint, err := w.contentEndpoint(content.Id)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", contentEndPoint.String(), strings.NewReader(string(jsonbody)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := w.sendRequest(req)
	if err != nil {
		return nil, err
	}

	var newContent Content
	err = json.Unmarshal(res, &newContent)
	if err != nil {
		return nil, err
	}

	return &newContent, nil
}

func (w *Wiki) GetChildrenByID(pageID string) ([]string, error) {
	contentEndPoint, err := w.contentEndpoint(pageID)
	if err != nil {
		return nil, err
	}
	childrenEndPoint, err := url.ParseRequestURI(contentEndPoint.String() + "/child/page")
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set("limit", "250")
	childrenEndPoint.RawQuery = data.Encode()

	req, err := http.NewRequest("GET", childrenEndPoint.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := w.sendRequest(req)
	if err != nil {
		return nil, err
	}

	children := new(SearchResults)
	err = json.Unmarshal(res, children)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(children.Results))

	for childNum, child := range children.Results {
		result[childNum] = child.Id
	}

	return result, nil
}
