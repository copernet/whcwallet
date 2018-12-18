package logic

import (
	"strconv"

	"github.com/copernet/whcwallet/model/view"
	"github.com/gin-gonic/gin"
)

type ErrNulldataLength struct {
	Msg string
}

func (en *ErrNulldataLength) Error() string {
	return en.Msg
}

type CheckParamsAndGetPayload interface {
	Check(*gin.Context) error
	CreatePayload(c *gin.Context) (string, error)
}

func FactoryForTxCheckAndCreate(txType int) CheckParamsAndGetPayload {
	switch txType {
	case 0:
		return &SimpleSendTransaction{}
	case 1:
		return &ParticipateCrowdSaleTransaction{}
	case 3:
		return &SendStoTransaction{}
	case 4:
		return &SendAllTransaction{}
	case 50:
		return &FixedIssuanceTransaction{}
	case 51:
		return &CrowdSaleIssuanceTransaction{}
	case 53:
		return &CloseCrowdSaleTransaction{}
	case 54:
		return &ManagedIssuanceTransaction{}
	case 55:
		return &SendGrantTransaction{}
	case 56:
		return &RevokeTransaction{}
	case 68:
		return &GetWhcTransaction{}
	case 70:
		return &ChangeIssuerTransaction{}
	case 185:
		return &FreezeTokenTransaction{}
	case 186:
		return &UnFreezeTokenTransaction{}
	default:
		return nil
	}
}

type SimpleSendTransaction struct {
	params view.SimpleSendTx
}

func (simpleSend *SimpleSendTransaction) Check(c *gin.Context) error {
	req := view.SimpleSendTx{}
	err := c.ShouldBind(&req)
	if err != nil {
		return err
	}

	simpleSend.params = req
	return nil
}

func (simpleSend *SimpleSendTransaction) CreatePayload(c *gin.Context) (string, error) {
	client := view.GetRPCIns()

	payload, err := client.WhcCreatePayloadSimpleSend(simpleSend.params.CID,
		c.PostForm("amount_to_transfer"))

	if err != nil {
		return "", err
	}
	return payload, nil
}

type ParticipateCrowdSaleTransaction struct {
	params view.ParticipateCrowdSale
}

func (partic *ParticipateCrowdSaleTransaction) Check(c *gin.Context) error {
	req := view.ParticipateCrowdSale{}
	err := c.ShouldBind(&req)
	if err != nil {
		return err
	}

	partic.params = req
	return nil
}

func (partic *ParticipateCrowdSaleTransaction) CreatePayload(c *gin.Context) (string, error) {
	client := view.GetRPCIns()
	payload, err := client.WhcCreatePayloadPartiCrowdSale(c.PostForm("amount_to_transfer"))
	if err != nil {
		return "", err
	}

	return payload, nil
}

type SendStoTransaction struct {
	params view.SendSto
}

func (sendSto *SendStoTransaction) Check(c *gin.Context) error {
	req := view.SendSto{}
	err := c.ShouldBind(&req)
	if err != nil {
		return err
	}

	sendSto.params = req
	return nil
}

func (sendSto *SendStoTransaction) CreatePayload(c *gin.Context) (string, error) {
	client := view.GetRPCIns()
	payload, err := client.WhcCreatePayloadSto(sendSto.params.PID,
		c.PostForm("amount"), sendSto.params.DisPID)
	if err != nil {
		return "", err
	}

	return payload, nil
}

type SendAllTransaction struct {
	params view.SendAll
}

func (sendAll *SendAllTransaction) Check(c *gin.Context) error {
	req := view.SendAll{}
	err := c.ShouldBind(&req)
	if err != nil {
		return err
	}

	sendAll.params = req
	return nil
}

func (sendAll *SendAllTransaction) CreatePayload(c *gin.Context) (string, error) {
	client := view.GetRPCIns()
	payload, err := client.WhcCreatePayloadSendAll(sendAll.params.Eco)
	if err != nil {
		return "", err
	}

	return payload, nil
}

type FixedIssuanceTransaction struct {
	params view.FixedIssuance
}

func (fixedIssuance *FixedIssuanceTransaction) Check(c *gin.Context) error {
	req := view.FixedIssuance{}
	err := c.ShouldBind(&req)
	if err != nil {
		return err
	}

	length := calPartialNulldataLength(req.Category, req.SubCategory, req.Name, req.Url, req.Data)
	// core: 220 - 19 for all data volume.
	// 5 field, every field have max 5 bytes for length indicator
	// So limit: 220 - 19 - 5*5 = 176
	limit := 177
	if length > limit {
		return &ErrNulldataLength{
			Msg: "the sum length of name/url/category/subcategory/data over the limit: " + strconv.Itoa(limit),
		}
	}

	fixedIssuance.params = req
	return nil
}

func (fixedIssuance *FixedIssuanceTransaction) CreatePayload(c *gin.Context) (string, error) {
	client := view.GetRPCIns()
	payload, err := client.WhcCreatePayloadIssuanceFixed(fixedIssuance.params.Eco,
		fixedIssuance.params.Precision, fixedIssuance.params.PrevPID, fixedIssuance.params.Category,
		fixedIssuance.params.SubCategory, fixedIssuance.params.Name, fixedIssuance.params.Url,
		fixedIssuance.params.Data, fixedIssuance.params.TotalNumber)

	if err != nil {
		return "", err
	}

	return payload, nil
}

type CrowdSaleIssuanceTransaction struct {
	params view.CrowdSaleIssuance
}

func (crowSale *CrowdSaleIssuanceTransaction) Check(c *gin.Context) error {
	req := view.CrowdSaleIssuance{}
	err := c.ShouldBind(&req)
	if err != nil {
		return err
	}

	length := calPartialNulldataLength(req.Category, req.SubCategory, req.Name, req.Url, req.Data)
	// core: 220 - 36 for all data volume.
	// 5 field, every field have max 5 bytes for length indicator
	// So limit: 220 - 36 - 5*5 = 159
	limit := 159
	if length > limit {
		return &ErrNulldataLength{
			Msg: "the sum length of name/url/category/subcategory/data over the limit: " + strconv.Itoa(limit),
		}
	}

	crowSale.params = req
	return nil
}

func (crowSale *CrowdSaleIssuanceTransaction) CreatePayload(c *gin.Context) (string, error) {
	client := view.GetRPCIns()
	payload, err := client.WhcCreatePayloadIssuanceCrowdSale(crowSale.params.Eco, crowSale.params.Precision,
		crowSale.params.PrevPID, crowSale.params.DesiredPID, crowSale.params.DeadLine, crowSale.params.EarlyBird,
		0, crowSale.params.Category, crowSale.params.SubCategory, crowSale.params.Name,
		crowSale.params.Url, crowSale.params.Data, crowSale.params.Exchange, crowSale.params.TotalNumber)

	if err != nil {
		return "", err
	}

	return payload, nil
}

type CloseCrowdSaleTransaction struct {
	params view.CloseCrowdSale
}

func (closeCrowdSale *CloseCrowdSaleTransaction) Check(c *gin.Context) error {
	req := view.CloseCrowdSale{}
	err := c.ShouldBind(&req)
	if err != nil {
		return err
	}

	closeCrowdSale.params = req
	return nil
}

func (closeCrowdSale *CloseCrowdSaleTransaction) CreatePayload(c *gin.Context) (string, error) {
	client := view.GetRPCIns()
	payload, err := client.WhcCreatePayloadCloseCrowdSale(closeCrowdSale.params.CID)
	if err != nil {
		return "", err
	}

	return payload, nil
}

type ManagedIssuanceTransaction struct {
	params view.ManagedIssuance
}

func (managedIssuance *ManagedIssuanceTransaction) Check(c *gin.Context) error {
	req := view.ManagedIssuance{}
	err := c.ShouldBind(&req)
	if err != nil {
		return err
	}

	length := calPartialNulldataLength(req.Category, req.SubCategory, req.Name, req.Url, req.Data)
	// core: 220 - 11 for all data volume.
	// 5 field, every field have max 5 bytes for length indicator
	// So limit: 220 - 11 - 5*5 = 184
	limit := 184
	if length > limit {
		return &ErrNulldataLength{
			Msg: "the sum length of name/url/category/subcategory/data over the limit: " + strconv.Itoa(limit),
		}
	}

	managedIssuance.params = req
	return nil
}

func (managedIssuance *ManagedIssuanceTransaction) CreatePayload(c *gin.Context) (string, error) {
	client := view.GetRPCIns()
	payload, err := client.WhcCreatePayloadIssuanceManaged(managedIssuance.params.Eco,
		managedIssuance.params.Precision, managedIssuance.params.PrevPID,
		managedIssuance.params.Category, managedIssuance.params.SubCategory,
		managedIssuance.params.Name, managedIssuance.params.Url, managedIssuance.params.Data)

	if err != nil {
		return "", err
	}

	return payload, err
}

type SendGrantTransaction struct {
	params view.SendGrant
}

func (grant *SendGrantTransaction) Check(c *gin.Context) error {
	req := view.SendGrant{}
	err := c.ShouldBind(&req)
	if err != nil {
		return err
	}

	if req.Note != nil {
		if len(*req.Note) > 200 {
			return &ErrNulldataLength{
				Msg: "the sum length of 'note' over the limit: " + strconv.Itoa(200),
			}
		}
	}

	grant.params = req
	return nil
}

func (grant *SendGrantTransaction) CreatePayload(c *gin.Context) (string, error) {
	client := view.GetRPCIns()
	payload, err := client.WhcCreatePayloadGrant(grant.params.CID,
		c.PostForm("amount"), grant.params.Note)
	if err != nil {
		return "", err
	}

	return payload, nil
}

type RevokeTransaction struct {
	params view.SendRevoke
}

func (revoke *RevokeTransaction) Check(c *gin.Context) error {
	req := view.SendRevoke{}
	err := c.ShouldBind(&req)
	if err != nil {
		return err
	}

	if req.Note != nil {
		if len(*req.Note) > 200 {
			return &ErrNulldataLength{
				Msg: "the sum length of 'note' over the limit: " + strconv.Itoa(200),
			}
		}
	}

	revoke.params = req
	return nil
}

func (revoke *RevokeTransaction) CreatePayload(c *gin.Context) (string, error) {
	client := view.GetRPCIns()
	payload, err := client.WhcCreatePayloadRevoke(revoke.params.CID,
		c.PostForm("amount"), revoke.params.Note)
	if err != nil {
		return "", err
	}

	return payload, nil
}

type GetWhcTransaction struct {
	params view.BurnBCH
}

func (burn *GetWhcTransaction) Check(c *gin.Context) error {
	req := view.BurnBCH{}
	err := c.ShouldBind(&req)
	if err != nil {
		return err
	}

	burn.params = req
	return nil
}

func (burn *GetWhcTransaction) CreatePayload(c *gin.Context) (string, error) {
	client := view.GetRPCIns()
	payload, err := client.WhcCreatePayloadBurnBCH()
	if err != nil {
		return "", err
	}

	return payload, nil
}

type ChangeIssuerTransaction struct {
	params view.ChangeIssuer
}

func (changeIssuer *ChangeIssuerTransaction) Check(c *gin.Context) error {
	req := view.ChangeIssuer{}
	err := c.ShouldBind(&req)
	if err != nil {
		return err
	}

	changeIssuer.params = req
	return nil
}

func (changeIssuer *ChangeIssuerTransaction) CreatePayload(c *gin.Context) (string, error) {
	client := view.GetRPCIns()
	payload, err := client.WhcCreatePayloadChangeIssuer(changeIssuer.params.CID)
	if err != nil {
		return "", err
	}

	return payload, nil
}

const (
	defaultDataLenght = 520
)

func calPartialNulldataLength(data ...string) int {
	volume := make([]byte, 0, defaultDataLenght)
	for _, item := range data {
		volume = append(volume, []byte(item)...)
	}

	return len(volume)
}

type FreezeTokenTransaction struct {
	params view.FreezeToken
}

func (simpleSend *FreezeTokenTransaction) Check(c *gin.Context) error {
	req := view.FreezeToken{}
	err := c.ShouldBind(&req)
	if err != nil {
		return err
	}

	simpleSend.params = req
	return nil
}

func (simpleSend *FreezeTokenTransaction) CreatePayload(c *gin.Context) (string, error) {
	client := view.GetRPCIns()

	payload, err := client.WhcCreatePayloadFreeze(simpleSend.params.FrozenAddress,
		simpleSend.params.PropertyId, simpleSend.params.Amount)

	if err != nil {
		return "", err
	}
	return payload, nil
}

type UnFreezeTokenTransaction struct {
	params view.FreezeToken
}

func (simpleSend *UnFreezeTokenTransaction) Check(c *gin.Context) error {
	req := view.FreezeToken{}
	err := c.ShouldBind(&req)
	if err != nil {
		return err
	}

	simpleSend.params = req
	return nil
}

func (simpleSend *UnFreezeTokenTransaction) CreatePayload(c *gin.Context) (string, error) {
	client := view.GetRPCIns()

	payload, err := client.WhcCreatePayloadUnFreeze(simpleSend.params.FrozenAddress,
		simpleSend.params.PropertyId, simpleSend.params.Amount)

	if err != nil {
		return "", err
	}
	return payload, nil
}
