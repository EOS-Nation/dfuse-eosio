package injector

import (
	"context"
	"fmt"

	"github.com/dfuse-io/bstream"
	pbcodec "github.com/dfuse-io/dfuse-eosio/pb/dfuse/eosio/codec/v1"
	"github.com/eoscanada/eos-go"
	"go.uber.org/zap"
)

func (i *Injector) processAction(ctx context.Context, blk *bstream.Block, act *pbcodec.ActionTrace, rawTraceMap map[uint64][]byte) error {
	accts := map[string]bool{
		act.Receiver: true,
	}
	for _, v := range act.Action.Authorization {
		accts[v.Actor] = true
	}

	for acct := range accts {
		acctUint := eos.MustStringToName(acct)

		actionKey := ActionKeyGenerator(blk, act, acctUint)

		acctSeqData, err := i.getSequenceData(ctx, actionKey)
		if err != nil {
			return fmt.Errorf("error while getting sequence data for account %v: %w", acct, err)
		}

		if acctSeqData.MaxEntries == 0 {
			return nil
		}

		// when shard 1 starts it will based the first seen action on values in shard 0. the last action for an account
		// will always have a greater last global seq
		if act.Receipt.GlobalSequence <= acctSeqData.LastGlobalSeq {
			zlog.Debug("this block has already been processed for this account",
				zap.Stringer("block", blk),
				zap.Stringer("key", actionKey),
			)
			return nil
		}

		lastDeletedSeq, err := i.deleteStaleRows(ctx, actionKey, acctSeqData)
		if err != nil {
			return fmt.Errorf("unable to delete stale rows: %w", err)
		}

		acctSeqData.LastDeletedOrdinal = lastDeletedSeq
		rawTrace := rawTraceMap[act.Receipt.GlobalSequence]

		// since the current ordinal is the last assgined order number we need to
		// increment it before we write a new action
		acctSeqData.CurrentOrdinal++
		if err = i.WriteAction(ctx, actionKey, acctSeqData, rawTrace); err != nil {
			return fmt.Errorf("error while writing action to store: %w", err)
		}

		acctSeqData.LastGlobalSeq = act.Receipt.GlobalSequence

		i.UpdateSeqData(actionKey, acctSeqData)
	}
	return nil
}
