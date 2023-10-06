package wrapper

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/golang/glog"
)

type Dapr struct {
	Host string
	Port string
}

type Message struct {
	ComponentName string
	Topic         string
	Data          []byte
}

func New(host, port string) *Dapr {
	if host == "" {
		host = "http://localhost"
	}
	if port == "" {
		port = "3500"
	}
	return &Dapr{Host: host, Port: port}
}

func (d *Dapr) BaseUrl() string {
	return d.Host + ":" + d.Port
}

// Get secret from a local secret store

func (d *Dapr) GetSecretFromFile(secretstore, key string) (string, error) {
	//getResponse, err := http.Get(DAPR_HOST + ":" + DAPR_HTTP_PORT + "/v1.0/secrets/" + DAPR_SECRET_STORE + "/" + SECRET_NAME)
	url := d.BaseUrl() + "/v1.0/secrets/" + secretstore + "/" + key
	getResponse, err := http.Get(url)

	if err != nil {
		return "", nil
	}
	result, err := io.ReadAll(getResponse.Body)
	if err != nil {
		return "", nil
	}

	var mapData map[string]string
	if err = json.Unmarshal(result, &mapData); err != nil {
		return "", nil
	} else {

	}

	return mapData[key], nil
}

func (d *Dapr) Publish(pubsubComponentName, topic string, message []byte) error {
	client := http.Client{}
	url := d.BaseUrl() + "/v1.0/publish/" + pubsubComponentName + "/" + topic
	req, err := http.NewRequest("POST", url, bytes.NewReader(message))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Publish an event using Dapr pub/sub
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (d *Dapr) PublichCh(chMessage chan Message) {
	client := http.Client{}
	for msg := range chMessage {
		url := d.BaseUrl() + "/v1.0/publish/" + msg.ComponentName + "/" + msg.Topic
		reader := bytes.NewReader(msg.Data)
		req, err := http.NewRequest("POST", url, reader)
		if err != nil {
			//return err
			glog.Errorln(err)
		}

		req.Header.Set("Content-Type", "application/json")

		// Publish an event using Dapr pub/sub
		res, err := client.Do(req)
		if err != nil {
			glog.Errorln(err)
		}
		res.Body.Close()
	}

}
