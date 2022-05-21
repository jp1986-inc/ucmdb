package rest

import (
	"context"
	"encoding/json"
	"fmt"
)

var ctx = context.Background()

type ucmdb struct {
	client *Client
}

type Ucmdb interface {
	CreateDataModel(tql *TopologyData) (*DataModelChange, error)
	ExecuteQuery(tql TopologyQuery) (*TopologyData, error)
	DeleteConfigurationItem(ucmdbId string) (*DataModelChange, error)
	UpdateConfigurationItem(id string, data DataInConfigurationItem) (*DataModelChange, error)
	GetConfigurationItem(ucmdbId string) (*DataInConfigurationItem, error)
}

func (u *ucmdb) ExecuteQuery(tql TopologyQuery) (*TopologyData, error) {
	req, err := u.client.newRequest("POST", "topologyQuery", tql)

	if err != nil {
		return nil, err
	}
	res, err := u.client.sendRequest(req)
	if err != nil {
		return nil, err
	}

	td := &TopologyData{}
	err = json.Unmarshal(res.([]byte), td)
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (u *ucmdb) CreateDataModel(tdata *TopologyData) (*DataModelChange, error) {
	req, err := u.client.newRequest("POST", "dataModel", tdata)
	if err != nil {
		return nil, err
	}
	res, err := u.client.sendRequest(req)
	if err != nil {
		return nil, err
	}

	dm := &DataModelChange{}
	err = json.Unmarshal(res.([]byte), dm)
	if err != nil {
		return nil, err
	}
	return dm, nil
}

func (u *ucmdb) UpdateConfigurationItem(id string, data DataInConfigurationItem) (*DataModelChange, error) {
	req, err := u.client.newRequest("PUT", fmt.Sprintf("dataModel/ci/%s", id), data)
	if err != nil {
		return nil, err
	}
	res, err := u.client.sendRequest(req)
	if err != nil {
		return nil, err
	}

	dm := &DataModelChange{}
	err = json.Unmarshal(res.([]byte), dm)
	if err != nil {
		return nil, err
	}
	return dm, nil
}

func (u *ucmdb) DeleteConfigurationItem(ucmdbId string) (*DataModelChange, error) {
	req, err := u.client.newRequest("DELETE", fmt.Sprintf("dataModel/ci/%s", ucmdbId), nil)
	if err != nil {
		return nil, err
	}
	res, err := u.client.sendRequest(req)
	if err != nil {
		return nil, err
	}

	dm := &DataModelChange{}
	err = json.Unmarshal(res.([]byte), dm)
	if err != nil {
		return nil, err
	}
	return dm, nil
}

func (u *ucmdb) GetConfigurationItem(ucmdbId string) (*DataInConfigurationItem, error) {
	req, err := u.client.newRequest("GET", fmt.Sprintf("dataModel/ci/%s", ucmdbId), nil)
	if err != nil {
		return nil, err
	}
	res, err := u.client.sendRequest(req)
	if err != nil {
		return nil, err
	}

	dci := &DataInConfigurationItem{}
	err = json.Unmarshal(res.([]byte), dci)
	if err != nil {
		return nil, err
	}
	return dci, nil
}
