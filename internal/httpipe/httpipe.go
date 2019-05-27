package httpipe

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Result []byte
	Err    error
}

func NewPipe(ctx context.Context, url string,
	headers map[string]string) <-chan Response {

	pipe := make(chan Response, 2)

	go func() {
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		for {
			select {
			case <-ctx.Done():
				close(pipe)
				return
			default:
				req, _ := http.NewRequest("GET", url, nil)
				for k, v := range headers {
					req.Header.Set(k, v)
				}
				resp, err := client.Do(req)
				if resp == nil || err != nil {
					log.Printf("Nil or error response: %v", err)
					pipe <- Response{
						Result: nil,
						Err:    err,
					}
					if resp != nil {
						defer resp.Body.Close()
					}
					break
				}

				defer resp.Body.Close()

				bytes, err := ioutil.ReadAll(resp.Body)
				log.Printf("Found response: %v", string(bytes))
				pipe <- Response{
					Result: bytes,
					Err:    err,
				}
			}
		}
	}()

	return pipe
}
