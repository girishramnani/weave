package multisig

import (
	"github.com/iov-one/weave"
	"github.com/iov-one/weave/x"
)

// Decorator checks multisig contract if available
type Decorator struct {
	auth   x.Authenticator
	bucket ContractBucket
}

var _ weave.Decorator = Decorator{}

// NewDecorator returns a default multisig decorator
func NewDecorator() Decorator {
	return Decorator{}
}

// Check enforce multisig contract before calling down the stack
func (d Decorator) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx,
	next weave.Checker) (weave.CheckResult, error) {
	var res weave.CheckResult
	newCtx, err := d.withMultisig(ctx, store, tx)
	if err != nil {
		return res, err
	}

	return next.Check(newCtx, store, tx)
}

// Deliver enforces multisig contract before calling down the stack
func (d Decorator) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx,
	next weave.Deliverer) (weave.DeliverResult, error) {
	var res weave.DeliverResult
	newCtx, err := d.withMultisig(ctx, store, tx)
	if err != nil {
		return res, err
	}

	return next.Deliver(newCtx, store, tx)
}

func (d Decorator) withMultisig(ctx weave.Context, store weave.KVStore, tx weave.Tx) (weave.Context, error) {
	if multisigContract, ok := tx.(MultiSigTx); ok {
		// does tx have multisig ?
		addr := multisigContract.GetMultiSig()
		if addr == nil {
			return ctx, nil
		}

		if d.auth.HasAddress(ctx, addr) {
			return ctx, nil
		}

		// load contract
		obj, err := d.bucket.Get(store, multisigContract.GetMultiSig())
		if err != nil {
			return ctx, err
		}
		contract := obj.Value().(*Contract)

		// retrieve sigs
		var sigs []weave.Address
		for _, sig := range contract.Sigs {
			sigs = append(sigs, sig)
		}

		// check sigs
		authenticated := x.HasNAddresses(ctx, d.auth, sigs, int(contract.ActivationThreshold))
		if !authenticated {
			return ctx, ErrUnauthorizedMultiSig(addr)
		}

		return withMultisig(ctx, addr), nil
	}

	return ctx, nil
}
