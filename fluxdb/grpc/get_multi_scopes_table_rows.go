package grpc

import (
	"context"
	"fmt"
	"sort"

	"github.com/dfuse-io/derr"
	"github.com/dfuse-io/dfuse-eosio/fluxdb"
	pbfluxdb "github.com/dfuse-io/dfuse-eosio/pb/dfuse/eosio/fluxdb/v1"
	"github.com/dfuse-io/dhammer"
	"github.com/dfuse-io/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

func (s *Server) GetMultiScopesTableRows(request *pbfluxdb.GetMultiScopesTableRowsRequest, stream pbfluxdb.State_GetMultiScopesTableRowsServer) error {
	ctx := stream.Context()
	zlogger := logging.Logger(ctx, zlog)
	zlogger.Debug("get multi scope tables rows",
		zap.Reflect("request", request),
	)

	blockNum := uint32(request.BlockNum)
	actualBlockNum, lastWrittenBlockID, upToBlockID, speculativeWrites, err := s.prepareRead(ctx, blockNum, request.IrreversibleOnly)
	if err != nil {
		return derr.Statusf(codes.Internal, "unable to prepare read: %s", err)
	}

	// Sort by scope so at least, a constant order is kept across calls
	sort.Slice(request.Scopes, func(leftIndex, rightIndex int) bool {
		return request.Scopes[leftIndex] < request.Scopes[rightIndex]
	})

	scopes := make([]interface{}, len(request.Scopes))
	for i, s := range request.Scopes {
		scopes[i] = string(s)
	}

	nailer := dhammer.NewNailer(64, func(ctx context.Context, i interface{}) (interface{}, error) {
		scope := i.(string)

		tablet := fluxdb.NewContractStateTablet(request.Contract, scope, request.Table)
		responseRows, err := s.readContractStateTable(
			ctx,
			tablet,
			actualBlockNum,
			request.KeyType,
			request.ToJson,
			request.WithBlockNum,
			speculativeWrites,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to read contract state tablet %q: %w", tablet, err)
		}

		resp := &pbfluxdb.TableRowsScopeResponse{
			Scope: scope,
			Row:   make([]*pbfluxdb.TableRowResponse, len(responseRows.Rows)),
		}

		for itr, row := range responseRows.Rows {
			resp.Row[itr] = processTableRow(&readTableRowResponse{
				Row: row,
			})
		}
		return resp, nil
	})

	nailer.PushAll(ctx, scopes)

	stream.SetHeader(getMetadata(upToBlockID, lastWrittenBlockID))

	for {
		select {
		case <-ctx.Done():
			zlog.Debug("stream terminated prior completion")
			return nil
		case next, ok := <-nailer.Out:
			if !ok {
				zlog.Debug("nailer completed")
				return nil
			}
			stream.Send(next.(*pbfluxdb.TableRowsScopeResponse))
		}
	}
}