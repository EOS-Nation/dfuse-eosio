package accounthist

import (
	"context"
	"fmt"
	"math"

	pbaccounthist "github.com/dfuse-io/dfuse-eosio/pb/dfuse/eosio/accounthist/v1"
	pbcodec "github.com/dfuse-io/dfuse-eosio/pb/dfuse/eosio/codec/v1"
	"github.com/dfuse-io/kvdb/store"
	"github.com/dfuse-io/logging"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

func (ws *Service) StreamActions(
	ctx context.Context,
	account uint64,
	limit uint64,
	cursor *pbaccounthist.Cursor,
	onAction func(cursor *pbaccounthist.Cursor, actionTrace *pbcodec.ActionTrace) error,
) error {
	logger := logging.Logger(ctx, zlog)

	queryShardNum := byte(0x00)
	querySeqNum := uint64(math.MaxUint64)
	if cursor != nil {
		queryShardNum = byte(cursor.ShardNum)
		querySeqNum = cursor.SequenceNumber - 1
	}

	startKey := encodeActionKey(account, queryShardNum, querySeqNum)
	endKey := store.Key(encodeActionPrefixKey(account)).PrefixNext()

	if limit == 0 || limit > ws.maxEntriesPerAccount {
		limit = ws.maxEntriesPerAccount
	}

	logger.Debug("scanning actions",
		zap.Stringer("account", EOSName(account)),
		zap.Stringer("start_key", Key(startKey)),
		zap.Stringer("end_key", Key(endKey)),
		zap.Uint64("limit", limit),
	)

	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	it := ws.kvStore.Scan(ctx, startKey, endKey, int(limit))
	for it.Next() {
		newact := &pbaccounthist.ActionRow{}
		err := proto.Unmarshal(it.Item().Value, newact)
		if err != nil {
			return fmt.Errorf("unmarshal action: %w", err)
		}

		if err := onAction(actionKeyToCursor(account, it.Item().Key), newact.ActionTrace); err != nil {
			return fmt.Errorf("on action: %w", err)
		}
	}

	if err := it.Err(); err != nil {
		return fmt.Errorf("fetching actions: %w", err)
	}

	return nil
}
