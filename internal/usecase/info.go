package usecase

import currency "p.solovev/internal/entity"

type InfoUsecase struct {
	appVersion string
	hostname   string
}

var _ InfoContract = (*InfoUsecase)(nil)

func NewInfoUsecase(appVersion string, hostname string) *InfoUsecase {
	return &InfoUsecase{appVersion: appVersion, hostname: hostname}
}

func (i InfoUsecase) GetInfoAboutService() currency.InfoResponse {
	return currency.InfoResponse{Version: i.appVersion, Service: "currency", Author: "p.solovev", Hostname: i.hostname}
}
