package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/seth16888/wxbusiness/internal/data/entities"
	"github.com/seth16888/wxbusiness/internal/model"
	"github.com/seth16888/wxbusiness/internal/model/request"
	v1 "github.com/seth16888/wxproxy/api/v1"
	"go.uber.org/zap"
)

type MPMemberRepo interface {
	Find(c context.Context, appId string,
		params *request.MPMemberQuery) (*model.PageResult[*entities.MPMember], error)
	FindById(c context.Context, id string) (*entities.MPMember, error)
	UpdateRemark(c context.Context, id, remark string) error
	Save(c context.Context, members []*entities.MPMember) error // 存在则更新，不存在则创建
	BatchTagging(c context.Context, appId string, ids []string, tagId int64) error
	BatchUnTagging(c context.Context, appId string, ids []string, tagId int64) error
}

type MPBlackListRepo interface {
	Block(c context.Context, appId string, openids []string) error
	Unblock(c context.Context, appId string, openids []string) error
	Query(c context.Context, appId string) ([]*entities.MPMember, error)
}

type MPMemberUsecase struct {
	repo          MPMemberRepo
	blackListRepo MPBlackListRepo
	log           *zap.Logger
	apiProxy      *APIProxyUsecase
}

func (m *MPMemberUsecase) PullBlackList(c context.Context, appId string) error {
	mpIdVar := c.Value("MP_ID")
	if mpIdVar == nil {
		m.log.Error("get mp id error")
		return fmt.Errorf("get mp id error")
	}
	mpId := mpIdVar.(string)

	token, err := m.apiProxy.GetAccessToken(c, appId, mpId)
	if err != nil {
		m.log.Error("get access token error", zap.Error(err))
		return fmt.Errorf("get access token error")
	}

  // 拉取微信黑名单
	// 数量超过 1000 时，可通过填写 next_openid 的值
	fetchOpenIdFn := func(nextOpenid string) ([]string, string, error) {
		openids := make([]string, 0, 1000)
		req := v1.GetBlacklistReq{
			AccessToken: token,
			NextOpenid:  nextOpenid, // 下一个拉取的openid，不填默认从头开始拉取
		}
		res, err := m.apiProxy.cli.GetBlacklist(c, &req)
		if err != nil {
			m.log.Error("get black list error", zap.Error(err))
			return nil, "", fmt.Errorf("get black list error")
		}

		m.log.Debug("get black list", zap.Int64("total", res.Total), zap.Int64("count", res.Count))
		if res.Count > 0 {
			if res.OpenIDs != nil {
				for _, openid := range res.OpenIDs {
					openids = append(openids, openid)
				}
			}
		} else {
			m.log.Debug("no member")
		}

		if res.Count < 1000 { // 少于1000条，说明已经拉取完了
			return openids, "", nil
		}
		return openids, res.NextOpenid, nil
	}

  memberIds := make([]string, 0, 1000)
	nextOpenid := ""
	for {
		openids, next, err := fetchOpenIdFn(nextOpenid)
		if err != nil {
			m.log.Error("fetch blacklist error", zap.Error(err))
			return fmt.Errorf("fetch blacklist error")
		}
		if len(openids) == 0 {
			break
		}
		memberIds = append(memberIds, openids...)
		if next == "" {
			break
		}
		nextOpenid = next
	}
	m.log.Debug("get blacklist ids", zap.Int("count", len(memberIds)))

  // 拉黑
	if err := m.blackListRepo.Block(c, appId, memberIds); err != nil {
		m.log.Error("save blacklist error", zap.Error(err))
		return fmt.Errorf("save blacklist error")
	}

	return nil
}

func (m *MPMemberUsecase) BatchUnTagging(c context.Context, appId string, ids []string, id int64) error {
	mpIdVar := c.Value("MP_ID")
	if mpIdVar == nil {
		m.log.Error("get mp id error")
		return fmt.Errorf("get mp id error")
	}
	mpId := mpIdVar.(string)

	token, err := m.apiProxy.GetAccessToken(c, appId, mpId)
	if err != nil {
		m.log.Error("get access token error", zap.Error(err))
		return fmt.Errorf("get access token error")
	}

	params := v1.BatchUnTaggingMembersRequest{
		AccessToken: token,
		Id:          id,
		OpenidList:  ids,
	}
	wxErr, err := m.apiProxy.cli.BatchUnTaggingMembers(c, &params)
	if err != nil {
		return err
	}
	if wxErr != nil && wxErr.Errcode != 0 {
		m.log.Error("api error", zap.Any("return", wxErr))
		return fmt.Errorf("call api error")
	}

	if err := m.repo.BatchUnTagging(c, appId, ids, id); err != nil {
		m.log.Error("batch UnTagging error", zap.Error(err))
		return fmt.Errorf("batch UnTagging error")
	}

	return nil
}

func (m *MPMemberUsecase) BatchTagging(c context.Context, appId string, ids []string, id int64) error {
	mpIdVar := c.Value("MP_ID")
	if mpIdVar == nil {
		m.log.Error("get mp id error")
		return fmt.Errorf("get mp id error")
	}
	mpId := mpIdVar.(string)

	token, err := m.apiProxy.GetAccessToken(c, appId, mpId)
	if err != nil {
		m.log.Error("get access token error", zap.Error(err))
		return fmt.Errorf("get access token error")
	}

	params := v1.BatchTaggingMembersRequest{
		AccessToken: token,
		Id:          id,
		OpenidList:  ids,
	}
	wxErr, err := m.apiProxy.cli.BatchTaggingMembers(c, &params)
	if err != nil {
		return err
	}
	if wxErr != nil && wxErr.Errcode != 0 {
		m.log.Error("api error", zap.Any("return", wxErr))
		return fmt.Errorf("call api error")
	}

	if err := m.repo.BatchTagging(c, appId, ids, id); err != nil {
		m.log.Error("batch Tagging error", zap.Error(err))
		return fmt.Errorf("batch Tagging error")
	}

	return nil
}

func (m *MPMemberUsecase) BatchUnblock(c context.Context, appId string, openids []string) error {
	mpIdVar := c.Value("MP_ID")
	if mpIdVar == nil {
		m.log.Error("get mp id error")
		return fmt.Errorf("get mp id error")
	}
	mpId := mpIdVar.(string)

	token, err := m.apiProxy.GetAccessToken(c, appId, mpId)
	if err != nil {
		m.log.Error("get access token error", zap.Error(err))
		return fmt.Errorf("get access token error")
	}

	params := v1.BlockMemberReq{
		AccessToken: token,
		OpenIds:     openids,
	}
	wxErr, err := m.apiProxy.cli.UnBlockMember(c, &params)
	if err != nil {
		return err
	}
	if wxErr != nil && wxErr.Errcode != 0 {
		m.log.Error("api error", zap.Any("return", wxErr))
		return fmt.Errorf("call api error")
	}

	if err := m.blackListRepo.Unblock(c, appId, openids); err != nil {
		m.log.Error("unblock member error", zap.Error(err))
		return fmt.Errorf("unblock member error")
	}
	return nil
}

func (m *MPMemberUsecase) BatchBlock(c context.Context, appId string, openids []string) error {
	mpIdVar := c.Value("MP_ID")
	if mpIdVar == nil {
		m.log.Error("get mp id error")
		return fmt.Errorf("get mp id error")
	}
	mpId := mpIdVar.(string)

	token, err := m.apiProxy.GetAccessToken(c, appId, mpId)
	if err != nil {
		m.log.Error("get access token error", zap.Error(err))
		return fmt.Errorf("get access token error")
	}

	params := v1.BlockMemberReq{
		AccessToken: token,
		OpenIds:     openids,
	}
	wxErr, err := m.apiProxy.cli.BlockMember(c, &params)
	if err != nil {
		return err
	}
	if wxErr != nil && wxErr.Errcode != 0 {
		m.log.Error("api error", zap.Any("return", wxErr))
		return fmt.Errorf("call api error")
	}

	if err := m.blackListRepo.Block(c, appId, openids); err != nil {
		m.log.Error("block member error", zap.Error(err))
		return fmt.Errorf("block member error")
	}
	return nil
}

func (m *MPMemberUsecase) GetBlackList(c context.Context, appId string) ([]*entities.MPMember, error) {
	docs, err := m.blackListRepo.Query(c, appId)
	if err != nil {
		m.log.Error("query black list error", zap.Error(err))
		return nil, fmt.Errorf("query black list error")
	}
	return docs, nil
}

func (m *MPMemberUsecase) UpdateRemark(c context.Context, appId, id, openId, remark string) error {
	app, err := GetAppInfoFromCtx(c)
	if err != nil {
		m.log.Error("get app info error", zap.Error(err))
		return fmt.Errorf("get app info error")
	}
	mpId := app.MpId
	accessToken, err := m.apiProxy.GetAccessToken(c, appId, mpId)
	if err != nil {
		m.log.Error("get access token error", zap.Error(err))
		return fmt.Errorf("get access token error")
	}

	req := v1.UpdateMemberRemarkRequest{
		AccessToken: accessToken,
		Openid:      openId,
		Remark:      remark,
	}
	if wxErr, err := m.apiProxy.cli.UpdateMemberRemark(c, &req); err != nil || wxErr.Errcode != 0 {
		m.log.Error("update remark error", zap.Error(err), zap.Any("wxErr", wxErr))
		return fmt.Errorf("update remark error")
	}

	_, err = m.repo.FindById(c, id)
	if err != nil {
		m.log.Error("update member remark error", zap.Error(err))
		return fmt.Errorf("data not found")
	}

	if err := m.repo.UpdateRemark(c, id, remark); err != nil {
		m.log.Error("update member remark error", zap.Error(err))
		return fmt.Errorf("update member remark error")
	}
	return nil
}

func (m *MPMemberUsecase) GetMemberInfo(c context.Context,
	appId string, id string,
) (*entities.MPMember, error) {
	return m.repo.FindById(c, id)
}

func (m *MPMemberUsecase) Query(c context.Context, appId string,
	params *request.MPMemberQuery,
) (*model.PageResult[*entities.MPMember], error) {
	docs, err := m.repo.Find(c, appId, params)
	if err != nil {
		m.log.Error("query member error", zap.Error(err))
		return nil, fmt.Errorf("query member error")
	}
	return docs, nil
}

// Pull 拉取微信公众号粉丝
func (m *MPMemberUsecase) Pull(c context.Context, appId string) error {
	mpIdVar := c.Value("MP_ID")
	if mpIdVar == nil {
		m.log.Error("get mp id error")
		return fmt.Errorf("get mp id error")
	}
	mpId := mpIdVar.(string)

	token, err := m.apiProxy.GetAccessToken(c, appId, mpId)
	if err != nil {
		m.log.Error("get access token error", zap.Error(err))
		return fmt.Errorf("get access token error")
	}

	// 拉取微信公众号粉丝Openid
	// 微信接口调用，每次最多拉取10000条，需要多次调用
	fetchOpenIdFn := func(nextOpenid string) ([]string, string, error) {
		openids := make([]string, 0, 10000)
		req := v1.GetMemberListRequest{
			AccessToken: token,
			NextOpenid:  nextOpenid, // 下一个拉取的openid，不填默认从头开始拉取
		}
		res, err := m.apiProxy.cli.GetMemberList(c, &req)
		if err != nil {
			m.log.Error("get member list error", zap.Error(err))
			return nil, "", fmt.Errorf("get member list error")
		}

		m.log.Debug("get member list", zap.Int64("total", res.Total), zap.Int64("count", res.Count))
		if res.Count > 0 {
			if res.Data != nil && res.Data.Openid != nil {
				for _, openid := range res.Data.Openid {
					openids = append(openids, openid.Openid)
				}
			}
		} else {
			m.log.Debug("no member")
		}

		if res.Count < 10000 { // 少于10000条，说明已经拉取完了
			return openids, "", nil
		}
		return openids, res.NextOpenid, nil
	}

	memberIds := make([]string, 0, 10000)
	nextOpenid := ""
	for {
		openids, next, err := fetchOpenIdFn(nextOpenid)
		if err != nil {
			m.log.Error("fetch member error", zap.Error(err))
			return fmt.Errorf("fetch member error")
		}
		if len(openids) == 0 {
			break
		}
		memberIds = append(memberIds, openids...)
		if next == "" {
			break
		}
		nextOpenid = next
	}
	m.log.Debug("get member ids", zap.Int("count", len(memberIds)))

	// 获取粉丝信息
	// 微信接口调用，每次最多拉取100条，需要多次调用
	fetchMemberInfoFn := func(openids []string) ([]*entities.MPMember, error) {
		req := v1.BatchGetMemberInfoRequest{
			AccessToken: token,
			UserList:    []*v1.BatchGetMemberInfoRequest_OpenIdList{},
		}
		for _, openid := range openids {
			req.UserList = append(req.UserList, &v1.BatchGetMemberInfoRequest_OpenIdList{
				Openid: openid,
			})
		}

		res, err := m.apiProxy.cli.BatchGetMemberInfo(c, &req)
		if err != nil {
			m.log.Error("batch get member info error", zap.Error(err))
			return nil, fmt.Errorf("batch get member info error")
		}
		m.log.Debug("get member info", zap.Int("total", len(res.GetUserListInfo())))

		members := make([]*entities.MPMember, 0, len(res.GetUserListInfo()))
		for _, info := range res.GetUserListInfo() {
			if info.Subscribe == 0 { // 未关注, 跳过
				continue
			}
			now := time.Now().Unix()
			fans := &entities.MPMember{
				AppId:          appId,
				MpId:           mpId,
				Subscribe:      1,
				OpenId:         info.Openid,
				NickName:       "",
				Sex:            0,
				Language:       info.Language,
				City:           "",
				Province:       "",
				Country:        "",
				SubscribeTime:  info.SubscribeTime,
				UnionId:        info.Unionid,
				Remark:         info.Remark,
				GroupId:        info.Groupid,
				Tags:           []*entities.MemberTag{},
				SubscribeScene: info.SubscribeScene,
				QrScene:        info.QrScene,
				QrSceneStr:     info.QrSceneStr,
				MessageCount:   0,
				CommentCount:   0,
				StarComment:    0,
				PraiseCount:    0,
				PraiseAmounts:  0,
				CreatedAt:      now,
				UpdatedAt:      now,
				Blocked:        false,
			}
			// tags
			if info.TagidList != nil {
				for _, tagId := range info.TagidList {
					fans.Tags = append(fans.Tags, &entities.MemberTag{
						MpId:  mpId,
						TagId: tagId,
						AppId: appId,
					})
				}
			}
			members = append(members, fans)
		}

		return members, nil
	}

	// 分批拉取粉丝信息
	batchSize := 100
	for i := 0; i < len(memberIds); i += batchSize {
		end := i + batchSize
		if end > len(memberIds) {
			end = len(memberIds)
		}
		members, err := fetchMemberInfoFn(memberIds[i:end])
		if err != nil {
			m.log.Error("fetch member info error", zap.Error(err))
			return fmt.Errorf("fetch member info error")
		}
		if len(members) == 0 {
			continue
		}
		// 保存粉丝信息
		if err := m.repo.Save(c, members); err != nil {
			m.log.Error("save member error", zap.Error(err))
			return fmt.Errorf("save member error")
		}
	}

	return nil
}

func NewMPMemberUsecase(log *zap.Logger, repo MPMemberRepo,
	apiProxy *APIProxyUsecase, block MPBlackListRepo,
) *MPMemberUsecase {
	return &MPMemberUsecase{repo: repo, log: log, apiProxy: apiProxy, blackListRepo: block}
}
