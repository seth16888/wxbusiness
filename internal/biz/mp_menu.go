package biz

import (
	"context"

	"github.com/seth16888/wxbusiness/internal/data/entities"
	"github.com/seth16888/wxbusiness/internal/model/r"
	"github.com/seth16888/wxcommon/mp"
	v1 "github.com/seth16888/wxproxy/api/v1"
	"go.uber.org/zap"
)

type MenuRepo interface {
	GetMenuInfo(ctx context.Context, appId string) (*entities.MPMenu, error)
	SaveMenu(ctx context.Context, apiMenu *entities.MPMenu) error
  DeleteMenu(ctx context.Context, pId string) error
}

// MPMenuUsecase 微信公众号菜单业务逻辑处理类
//
// Pull 同步原菜单：如果是官方平台设置，转换为自定义菜单，如果是自定义菜单，则直接同步。
// 官方平台设置的菜单不保存，只转换为自定义菜单。
// 数据库中保存自定义菜单。
type MPMenuUsecase struct {
	repo     MenuRepo
	appRepo  AppRepo
	tokenUc  *AccessTokenUsecase
	apiProxy *APIProxyUsecase
	log      *zap.Logger
}

func NewMPMenuUsecase(repo AppRepo,
	tokenUc *AccessTokenUsecase,
	apiProxy *APIProxyUsecase,
	log *zap.Logger,
	menuRepo MenuRepo,
) *MPMenuUsecase {
	return &MPMenuUsecase{appRepo: repo, tokenUc: tokenUc, apiProxy: apiProxy, log: log, repo: menuRepo}
}

func (u *MPMenuUsecase) Pull(ctx context.Context, appId string) *r.R {
	app, err := GetAppInfoFromCtx(ctx)
	if err != nil {
		u.log.Error("get app info error", zap.Error(err))
		return r.Error(400, "get app info error")
	}
	mpId := app.MpId
	// accessToken
	tokenRes, err := u.tokenUc.FetchAccessToken(ctx, appId, mpId)
	if err != nil {
		return r.Error(403, "fetch access token error")
	}

	// 查询接口
	reply, err := u.apiProxy.cli.PullMenu(ctx, &v1.AccessTokenParam{AccessToken: tokenRes.AccessToken})
	if err != nil {
		u.log.Error("pull menu error", zap.Error(err))
		return r.Error(400, "pull menu error")
	}

	if reply.IsMenuOpen == 0 { // 菜单没有开启，不处理
		return r.SuccessData(nil)
	}

	// 转换为自定义菜单
	apiMenu := u.convertToApiMenu(reply)
	apiMenu.AppId = appId
	apiMenu.MpId = mpId

	// 保存自定义菜单
	err = u.repo.SaveMenu(ctx, apiMenu)
	if err != nil {
		u.log.Error("save menu error", zap.Error(err))
		return r.Error(400, "save menu error")
	}

	return r.SuccessData(nil)
}

func (u *MPMenuUsecase) convertToApiMenu(reply *v1.SelfMenuReply) *entities.MPMenu {
	typeFunc := func(d *entities.MenuButton, value string) {
		switch d.Type {
		case "text": // TODO: 自动回复：文本
			d.Type = "click"
			d.Key = value
		case "img", "voice":
			d.Type = "media_id"
			d.MediaID = value
		case "video":
			d.Type = "view"
			d.URL = value
		case "news":
			d.Type = "article_view_limited"
			d.MediaID = value
		}
	}
	menu := &entities.MPMenu{
		AppId:           "",
		MpId:            "",
		MenuID:          0,
		Button:          []*entities.MenuButton{},
		Conditionalmenu: []*entities.ConditionalMenuRes{},
		CreatedAt:       0,
		UpdatedAt:       0,
	}
	for _, btn := range reply.SelfmenuInfo.Button {
		doc := entities.MenuButton{
			Type:       btn.Type,
			Name:       btn.Name,
			Key:        btn.Key,
			URL:        btn.Url,
			MediaID:    btn.Value,
			AppID:      "",
			PagePath:   "",
			SubButtons: []*entities.MenuButton{},
		}
		typeFunc(&doc, btn.Value)
		for _, subBtn := range btn.SubButton.List {
			docSub := entities.MenuButton{
				Type:       subBtn.Type,
				Name:       subBtn.Name,
				Key:        subBtn.Key,
				URL:        subBtn.Url,
				MediaID:    subBtn.Value,
				AppID:      "",
				PagePath:   "",
				SubButtons: []*entities.MenuButton{},
			}
			typeFunc(&docSub, subBtn.Value)
			doc.SubButtons = append(doc.SubButtons, &docSub)
		}
		menu.Button = append(menu.Button, &doc)
	}

	return menu
}

func (u *MPMenuUsecase) Create(ctx context.Context, pId string, params *mp.CreateMenuReq) *r.R {
	app, err := GetAppInfoFromCtx(ctx)
	if err != nil {
		u.log.Error("get app info error", zap.Error(err))
		return r.Error(400, "get app info error")
	}
	mpId := app.MpId
	// accessToken
	akRes, err := u.tokenUc.FetchAccessToken(ctx, pId, mpId)
	if err != nil {
		return r.Error(403, "fetch access token error")
	}

	err = u.apiProxy.CreateMenu(ctx, akRes.AccessToken, params)
	if err != nil {
		return r.Error(400, "create menu error")
	}

  // 保存自定义菜单
  reply, err := u.apiProxy.cli.PullMenu(ctx, &v1.AccessTokenParam{AccessToken: akRes.AccessToken})
  if err != nil {
		u.log.Error("pull menu error", zap.Error(err))
		return r.Error(400, "pull menu error")
	}

	if reply.IsMenuOpen == 0 { // 菜单没有开启，不处理
		return r.SuccessData(nil)
	}

	// 转换为自定义菜单
	apiMenu := u.convertToApiMenu(reply)
	apiMenu.AppId = pId
	apiMenu.MpId = mpId

  err = u.repo.SaveMenu(ctx, apiMenu)
  if err != nil {
    u.log.Error("save menu error", zap.Error(err))
    return r.Error(400, "save menu error")
  }

	return r.Success()
}

// GetMenuInfo 获取菜单信息
func (u *MPMenuUsecase) GetMenuInfo(ctx context.Context, pId string) *r.R {
	menu, err := u.repo.GetMenuInfo(ctx, pId)
	if err != nil {
		return r.Error(404, "not found")
	}

	return r.SuccessData(menu)
}

// Delete 删除菜单
func (u *MPMenuUsecase) Delete(ctx context.Context, pId string) *r.R {
	app, err := GetAppInfoFromCtx(ctx)
	if err != nil {
		u.log.Error("get app info error", zap.Error(err))
		return r.Error(400, "get app info error")
	}
	mpId := app.MpId
	// accessToken
	akRes, err := u.tokenUc.FetchAccessToken(ctx, pId, mpId)
	if err != nil {
		return r.Error(403, "fetch access token error")
	}

	wxErr, err := u.apiProxy.cli.DeleteMenu(ctx, &v1.AccessTokenParam{AccessToken: akRes.AccessToken})
	if err != nil {
		return r.Error(400, "delete menu error")
	}
	if wxErr.Errcode != 0 {
		return r.Error(400, "delete menu error")
	}

	// 删除数据库中的菜单信息
	err = u.repo.DeleteMenu(ctx, pId)
	if err != nil {
		u.log.Error("delete db menu error", zap.Error(err))
		return r.Error(400, "delete db menu error")
	}

	return r.Success()
}
