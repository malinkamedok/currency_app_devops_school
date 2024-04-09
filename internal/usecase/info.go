package usecase

import currency "p.solovev/internal/entity"

type InfoUsecase struct {
	appVersion string
}

var _ InfoContract = (*InfoUsecase)(nil)

func NewInfoUsecase(appVersion string) *InfoUsecase {
	return &InfoUsecase{appVersion: appVersion}
}

func (i InfoUsecase) GetInfoAboutService() currency.InfoResponse {
	return currency.InfoResponse{Version: i.appVersion, Service: "currency", Author: "p.solovev"}
}
