package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"tg-api/pkg/e"
	"tg-api/pkg/types"
)

// TYPE

type Client struct {
	host     string
	basePath string
	client   http.Client
}

// CONSTRUCTOR

func New(host string, token string) Client {
	return Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

//CONST

const (
	getUpdatesMethod     = "getUpdates"
	sendMessageMethod    = "sendMessage"
	forwardMessageMethod = "forwardMessage"
	copyMessageMethod    = "copyMessage"
	sendPhotoMethod      = "sendPhoto"
)

// PUBLIC METHOD

func (c *Client) Updates(offset int, limit int) (data []types.Update, err error) {
	defer func() {
		err = e.WrapIfErr("can't get updates:", err)
	}()

	query := url.Values{}
	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))

	updates, err := c.doRequestGET(getUpdatesMethod, query)
	if err != nil {
		return nil, err
	}

	var res types.UpdatesResponse
	if err := json.Unmarshal(updates, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) (data []byte, err error) {
	defer func() {
		err = e.WrapIfErr("can't send message", err)
	}()

	query := url.Values{}
	query.Add("chat_id", strconv.Itoa(chatID))
	query.Add("text", text)

	data, err = c.doRequestGET(sendMessageMethod, query)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) ForwardMessage(chatID int, fromID int, messageID int) (data []byte, err error) {
	defer func() {
		err = e.WrapIfErr("can't send message", err)
	}()

	query := url.Values{}
	query.Add("chat_id", strconv.Itoa(chatID))
	query.Add("from_chat_id", strconv.Itoa(fromID))
	query.Add("message_id", strconv.Itoa(messageID))

	data, err = c.doRequestGET(forwardMessageMethod, query)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) CopyMessage(chatID int, fromID int, messageID int) (data []byte, err error) {
	defer func() {
		err = e.WrapIfErr("can't send message", err)
	}()

	query := url.Values{}
	query.Add("chat_id", strconv.Itoa(chatID))
	query.Add("from_chat_id", strconv.Itoa(fromID))
	query.Add("message_id", strconv.Itoa(messageID))

	data, err = c.doRequestGET(copyMessageMethod, query)
	if err != nil {
		return nil, err
	}

	return data, nil
}
func (c *Client) SendPhoto(chatID int, format string, file string) (data []byte, err error) {
	defer func() {
		err = e.WrapIfErr("can't send message", err)
	}()

	query := url.Values{}
	query.Add("chat_id", strconv.Itoa(chatID))

	data, err = c.doRequestPOST(sendPhotoMethod, query, "photo", path.Join("image", format), file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// PRIVATE METHOD

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) doRequestGET(method string, query url.Values) (data []byte, err error) {
	defer func() {
		err = e.WrapIfErr("can't do request:", err)
	}()

	u := url.URL{
		Scheme:   "https",
		Host:     c.host,
		Path:     path.Join(c.basePath, method),
		RawQuery: query.Encode(),
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
func (c *Client) doRequestPOST(method string, query url.Values, fieldname string, contenttype string, filename string) (data []byte, err error) {
	defer func() {
		err = e.WrapIfErr("can't do request:", err)
	}()

	u := url.URL{
		Scheme:   "https",
		Host:     c.host,
		Path:     path.Join(c.basePath, method),
		RawQuery: query.Encode(),
	}

	//
	var buf bytes.Buffer

	w := multipart.NewWriter(&buf)

	h := textproto.MIMEHeader{}
	h.Add("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
		fieldname, filepath.Base(filename)))
	h.Add("Content-Type", contenttype)

	fw, err := w.CreatePart(h)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if _, err := io.Copy(fw, file); err != nil {
		return nil, err
	}

	w.Close()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, u.String(), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
