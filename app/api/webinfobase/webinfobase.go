package webinfobase

import (
	"context"
	"encoding/xml"
)

type (
	Config struct {
		Address string `yaml:"address"`
		Port    int64  `yaml:"port"`
	}

	Service interface {
		WSDL() []byte
		CheckInfoBases(ctx context.Context, clientID string) (CheckInfoBasesResponse, error)
		GetInfoBases(ctx context.Context, clientID string) (GetInfoBasesResponse, error)
	}

	CheckInfoBasesResponse struct {
		URL     string
		Changed bool
	}

	GetInfoBasesResponse struct {
		ClientID string
		Text     string
	}

	CheckInfoBasesRequest struct {
		XMLName xml.Name `xml:"CheckInfoBases"`
		ID      string   `xml:"ID" validate:"required,uuid"`
		Code    string   `xml:"Code" validate:"required,uuid"`
	}

	CheckInfoBasesRequestBody struct {
		XMLName               xml.Name              `xml:"Body"`
		CheckInfoBasesRequest CheckInfoBasesRequest `xml:"https://titovcode.com/WebCommonInfoBases CheckInfoBases"`
	}

	CheckInfoBasesRequestWrapper struct {
		XMLName xml.Name                  `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
		Body    CheckInfoBasesRequestBody `xml:"Body"`
	}

	GetInfoBasesRequest struct {
		XMLName xml.Name `xml:"GetInfoBases"`
		ID      string   `xml:"ID" validate:"required,uuid"`
	}

	GetInfoBasesRequestBody struct {
		XMLName             xml.Name            `xml:"Body"`
		GetInfoBasesRequest GetInfoBasesRequest `xml:"https://titovcode.com/WebCommonInfoBases GetInfoBases"`
	}

	GetInfoBasesRequestWrapper struct {
		XMLName xml.Name                `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
		Body    GetInfoBasesRequestBody `xml:"Body"`
	}
)
