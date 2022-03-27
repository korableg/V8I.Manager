package webinfobase

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/korableg/V8I.Manager/app/api/client"
	"github.com/korableg/V8I.Manager/app/api/onecdb"
	"text/template"
)

//go:embed webcommoninfobases.wsdl
var wsdlTmplt string

type (
	Service interface {
		WSDL() []byte
		CheckInfoBases(ctx context.Context, clientID string) (CheckInfoBasesResponse, error)
		GetInfoBases(ctx context.Context, clientID string) (GetInfoBasesResponse, error)
	}

	service struct {
		wsdl          []byte
		v8iBuilder    onecdb.V8IBuilder
		clientService client.Service
	}
)

func NewService(address string, port int, v8iBuilder onecdb.V8IBuilder, clientService client.Service) (*service, error) {
	if v8iBuilder == nil {
		return nil, errors.New("v8i builder is not defined")
	}

	if clientService == nil {
		return nil, errors.New("client service is not defined")
	}

	w, err := template.New("wsdl").Parse(wsdlTmplt)
	if err != nil {
		return nil, fmt.Errorf("parse wsdl: %w", err)
	}

	buf := bytes.NewBuffer(nil)
	baseURL := fmt.Sprintf("http://%s:%d", address, port)
	if err = w.Execute(buf, baseURL); err != nil {
		return nil, fmt.Errorf("execute wsdl template: %w", err)
	}

	s := &service{
		wsdl:          buf.Bytes(),
		v8iBuilder:    v8iBuilder,
		clientService: clientService,
	}

	return s, nil
}

func (s *service) WSDL() []byte {
	return s.wsdl
}

func (s *service) CheckInfoBases(_ context.Context, _ string) (CheckInfoBasesResponse, error) {
	return CheckInfoBasesResponse{
		URL:     "/WebCommonInfoBases",
		Changed: true,
	}, nil
}

func (s *service) GetInfoBases(ctx context.Context, clientID string) (GetInfoBasesResponse, error) {
	if clientID == "00000000-0000-0000-0000-000000000000" {
		var err error
		if clientID, err = s.clientService.NewClient(ctx); err != nil {
			return GetInfoBasesResponse{}, fmt.Errorf("couldn't create new client: %w", err)
		}
	}

	v8iData, err := s.v8iBuilder.BuildV8I(ctx)
	if err != nil {
		return GetInfoBasesResponse{}, err
	}

	return GetInfoBasesResponse{
		ClientID: clientID,
		Text:     v8iData,
	}, nil
}
