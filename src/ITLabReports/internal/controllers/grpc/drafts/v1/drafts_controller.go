package drafts

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/internal/common/errors"
	drafts "github.com/RTUITLab/ITLab-Reports/internal/domain/drafts/service"
	user "github.com/RTUITLab/ITLab-Reports/internal/domain/user/service"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/RTUITLab/ITLab/proto/reports/types"
	reportsgrpc "github.com/RTUITLab/ITLab/proto/reports/v1"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type DraftController struct {
	draftService drafts.DraftsService
	auth         user.UserService
	marshaller   DraftMarshaller

	reportsgrpc.UnimplementedDraftsServer

	AuthErrorsHandler
	ErrorFormatter
}

func NewDraftController(
	draftService drafts.DraftsService,
	auth user.UserService,
) *DraftController {
	return &DraftController{
		draftService: draftService,
		auth:         auth,
	}
}

func NewDraftControllerFrom(i *do.Injector) (*DraftController, error) {
	draftService := do.MustInvoke[drafts.DraftsService](i)
	auth := do.MustInvoke[user.UserService](i)

	return NewDraftController(draftService, auth), nil
}

func (rc *DraftController) Build(server *grpc.Server) {
	reportsgrpc.RegisterDraftsServer(server, rc)
}

func (rc *DraftController) HandlerError(err error) error {
	if err := rc.AuthErrorsHandler.Handle(err); err != nil {
		return err
	}

	logrus.WithFields(
		logrus.Fields{
			"controller": "drafts",
			"transport":  "grpc",
		},
	).Error(err)
	return rc.FormatError(codes.Internal, err)
}

// Return draft by id
// If draft not found return DRAFT_NOT_FOUND error
func (dc *DraftController) GetDraft(
	ctx context.Context,
	req *reportsgrpc.GetDraftReq,
) (*reportsgrpc.GetDraftResp, error) {
	user, err := dc.auth.AuthUser(
		ctx,
		TokenFromContext(ctx),
	)
	if err != nil {
		return nil, dc.AuthErrorsHandler.Handle(err)
	}

	draft, err := dc.draftService.GetDraft(
		ctx,
		drafts.GetDraftReq{
			ID: req.DraftId,
			By: user,
		},
	)
	if err == drafts.ErrDraftNotFound {
		return &reportsgrpc.GetDraftResp{
			Result: &reportsgrpc.GetDraftResp_Error{
				Error: reportsgrpc.DraftsServiceErrors_DRAFT_NOT_FOUND,
			},
		}, nil
	} else if errors.Is(err, drafts.ErrCantGetDraft) {
		return nil, dc.FormatError(codes.PermissionDenied, err)
	} else if err != nil {
		return nil, dc.HandlerError(err)
	}

	return &reportsgrpc.GetDraftResp{
		Result: &reportsgrpc.GetDraftResp_Draft{
			Draft: DraftFrom(draft),
		},
	}, nil

}

// Return list of drafts
func (dc *DraftController) GetDrafts(
	ctx context.Context,
	req *reportsgrpc.GetDraftsReq,
) (*reportsgrpc.GetDraftsResp, error) {
	user, err := dc.auth.AuthUser(
		ctx,
		TokenFromContext(ctx),
	)
	if err != nil {
		return nil, dc.AuthErrorsHandler.Handle(err)
	}

	drafts, err := dc.draftService.GetDrafts(
		ctx,
		dc.marshaller.MarshallGetDraftsReq(
			req,
			user,
		),
	)
	if err != nil {
		return nil, dc.HandlerError(err)
	}

	return &reportsgrpc.GetDraftsResp{
		Drafts: lo.Map(
			drafts,
			func(draft aggregate.Draft, _ int) *types.Draft {
				return DraftFrom(draft)
			},
		),
	}, nil
}

// Return paginated list of drafts
func (dc *DraftController) GetDraftsPaginated(
	ctx context.Context,
	req *reportsgrpc.GetDraftsPaginatedReq,
) (*reportsgrpc.GetDraftsPaginatedResp, error) {
	user, err := dc.auth.AuthUser(
		ctx,
		TokenFromContext(ctx),
	)
	if err != nil {
		return nil, dc.AuthErrorsHandler.Handle(err)
	}

	r := dc.marshaller.MarshallGetDraftsPaginatedReq(
		req,
		user,
	)

	drafts, err := dc.draftService.GetDrafts(
		ctx,
		r,
	)
	if err != nil {
		return nil, dc.HandlerError(err)
	}

	count, err := dc.draftService.CountDrafts(
		ctx,
		r,
	)
	if err != nil {
		return nil, dc.HandlerError(err)
	}

	return dc.marshaller.MarshallGetDraftsResp(req, drafts, count), nil

}
