// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package plugin

import (
	"fmt"
	"github.com/ysicing/sonarapi"
)

const (
	scanuser = "admin"
	scanpass = "admin"
)

func Api(sonarurl, sonaruser, sonarpass, sonarkey string) (*CIScan, error) {
	if sonaruser == "" || sonarpass == "" {
		sonaruser = scanuser
		sonarpass = scanpass
	}
	c, err := sonarapi.NewClient(sonarurl, sonaruser, sonarpass)
	if err != nil {
		return nil, err
	}
	return &CIScan{
		Key:  sonarkey,
		client: c,
	}, nil
}

type CIScan struct {
	Key string
	TokenPrefix string
	client *sonarapi.Client
}

func (ci *CIScan) tokenprefix() string  {
	if len(ci.TokenPrefix) == 0 {
		ci.TokenPrefix = getToday()
	}
	return ci.TokenPrefix
}

func (ci *CIScan) CheckProject() (bool, error) {
	s := sonarapi.ProjectsSearchOption{
		Q:                 ci.Key,
	}
	v, _, err := ci.client.Projects.Search(&s)
	if err != nil {
		return false, err
	}
	if len(v.Components) == 0 {
		return false, nil
	}
	for _, k := range v.Components {
		if k.Key == s.Q {
			return true, nil
		}
	}
	return false, nil
}

func (ci *CIScan) CreateProject() error {
	s := sonarapi.ProjectsCreateOption{
		Name:       ci.Key,
		Project:    ci.Key,
		Visibility: "public",
	}
	_, resp, err := ci.client.Projects.Create(&s)
	if err != nil {
		return err
	}
	if resp.StatusCode < 400 {
		return nil
	}
	return fmt.Errorf("resp code >= 400")
}

func (ci *CIScan) GenerateToken() (string, error) {
	if ci.SearchToken() {
		if err := ci.RevokeToken(); err != nil {
			return "", err
		}
	}

	s := sonarapi.UserTokensGenerateOption{
		Name: fmt.Sprintf("ci-%v-%v", ci.tokenprefix(), ci.Key),
	}
	v, _, err := ci.client.UserTokens.Generate(&s)
	if err != nil {
		return "", err
	}
	if len(v.Token) == 0 {
		return "", fmt.Errorf("token length is 0")
	}
	return v.Token, nil
}

func (ci *CIScan) RevokeToken() error {
	s := sonarapi.UserTokensRevokeOption{
		Name: fmt.Sprintf("ci-%v-%v", ci.tokenprefix(), ci.Key),
	}
	_, err := ci.client.UserTokens.Revoke(&s)
	if err != nil {
		return err
	}
	return nil
}

func (ci *CIScan) SearchToken() bool {
	s := sonarapi.UserTokensSearchOption{}
	v, _, err := ci.client.UserTokens.Search(&s)
	if err != nil {
		return false
	}
	// Name: fmt.Sprintf("ci-%v-%v", ci.tokenprefix(), ci.Key),
	for _, v := range v.UserTokens {
		if v.Name == fmt.Sprintf("ci-%v-%v", ci.tokenprefix(), ci.Key) {
			return true
		}
	}
	return false
}

func (ci *CIScan) Health() interface{} {
	v, _, err := ci.client.System.Status()
	if err != nil {
		return nil
	}
	if v == nil {
		return nil
	}
	return v
}