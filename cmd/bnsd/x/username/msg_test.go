package username

import (
	"testing"

	"github.com/iov-one/weave"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/weavetest"
)

func TestRegisterNamespaceMsgValidate(t *testing.T) {
	cases := map[string]struct {
		Msg  weave.Msg
		Want *errors.Error
	}{
		"valid message": {
			Msg: &RegisterNamespaceMsg{
				Metadata: &weave.Metadata{Schema: 1},
				Label:    "evilcorp",
				Public:   true,
			},
			Want: nil,
		},
		"label must be provided": {
			Msg: &RegisterNamespaceMsg{
				Metadata: &weave.Metadata{Schema: 1},
				Label:    "",
				Public:   true,
			},
			Want: errors.ErrEmpty,
		},
	}

	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			if err := tc.Msg.Validate(); !tc.Want.Is(err) {
				t.Fatal(err)
			}
		})
	}
}

func TestRegisterTokenMsgValidate(t *testing.T) {
	cases := map[string]struct {
		Msg  weave.Msg
		Want *errors.Error
	}{
		"valid message": {
			Msg: &RegisterTokenMsg{
				Metadata: &weave.Metadata{Schema: 1},
				Username: "alice*iov",
				Targets: []BlockchainAddress{
					{BlockchainID: "blobchain", Address: "1234567890"},
				},
			},
			Want: nil,
		},
		"empty targets": {
			Msg: &RegisterTokenMsg{
				Metadata: &weave.Metadata{Schema: 1},
				Username: "alice*iov",
				Targets:  nil,
			},
			Want: nil,
		},
		"different address but the same blockchain ID is not allowed": {
			Msg: &RegisterTokenMsg{
				Metadata: &weave.Metadata{Schema: 1},
				Username: "alice*iov",
				Targets: []BlockchainAddress{
					{BlockchainID: "blobchain", Address: "a-blobchain-id-1"},
					{BlockchainID: "blobchain", Address: "a-blobchain-id-2"},
				},
			},
			Want: errors.ErrDuplicate,
		},
	}

	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			if err := tc.Msg.Validate(); !tc.Want.Is(err) {
				t.Fatal(err)
			}
		})
	}
}

func TestTransferTokenMsgValidate(t *testing.T) {
	cases := map[string]struct {
		Msg  weave.Msg
		Want *errors.Error
	}{
		"valid message": {
			Msg: &TransferTokenMsg{
				Metadata: &weave.Metadata{Schema: 1},
				Username: "alice*iov",
				NewOwner: weavetest.NewCondition().Address(),
			},
			Want: nil,
		},
		"invalid new owner address": {
			Msg: &TransferTokenMsg{
				Metadata: &weave.Metadata{Schema: 1},
				Username: "alice*iov",
				NewOwner: []byte("x"),
			},
			Want: errors.ErrInput,
		},
	}

	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			if err := tc.Msg.Validate(); !tc.Want.Is(err) {
				t.Fatal(err)
			}
		})
	}
}

func TestChangeTokenTargetsMsgValidate(t *testing.T) {
	cases := map[string]struct {
		Msg  weave.Msg
		Want *errors.Error
	}{
		"valid message": {
			Msg: &ChangeTokenTargetsMsg{
				Metadata: &weave.Metadata{Schema: 1},
				Username: "alice*iov",
				NewTargets: []BlockchainAddress{
					{BlockchainID: "blobchain", Address: "1234567890"},
				},
			},
			Want: nil,
		},
		"invalid new targets": {
			Msg: &ChangeTokenTargetsMsg{
				Metadata: &weave.Metadata{Schema: 1},
				Username: "alice*iov",
				NewTargets: []BlockchainAddress{
					{BlockchainID: "x", Address: "x"},
				},
			},
			Want: errors.ErrInput,
		},
		"missing new targets": {
			Msg: &ChangeTokenTargetsMsg{
				Metadata:   &weave.Metadata{Schema: 1},
				Username:   "alice*iov",
				NewTargets: []BlockchainAddress{},
			},
			Want: nil,
		},
		"different address but the same blockchain ID is not allowed": {
			Msg: &ChangeTokenTargetsMsg{
				Metadata: &weave.Metadata{Schema: 1},
				Username: "alice*iov",
				NewTargets: []BlockchainAddress{
					{BlockchainID: "blobchain", Address: "a-blobchain-id-1"},
					{BlockchainID: "blobchain", Address: "a-blobchain-id-2"},
				},
			},
			Want: errors.ErrDuplicate,
		},
	}

	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			if err := tc.Msg.Validate(); !tc.Want.Is(err) {
				t.Fatal(err)
			}
		})
	}
}
