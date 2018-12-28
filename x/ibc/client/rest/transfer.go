package rest

import (
	"net/http"

	"github.com/ftlnetwork/ftlnetwork-sdk/client/context"
	"github.com/ftlnetwork/ftlnetwork-sdk/client/utils"
	"github.com/ftlnetwork/ftlnetwork-sdk/codec"
	"github.com/ftlnetwork/ftlnetwork-sdk/crypto/keys"
	sdk "github.com/ftlnetwork/ftlnetwork-sdk/types"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/ibc"

	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, kb keys.Keybase) {
	r.HandleFunc("/ibc/{destchain}/{address}/send", TransferRequestHandlerFn(cdc, kb, cliCtx)).Methods("POST")
}

type transferReq struct {
	BaseReq utils.BaseReq `json:"base_req"`
	Amount  sdk.Coins     `json:"amount"`
}

// TransferRequestHandler - http request handler to transfer coins to a address
// on a different chain via IBC.
func TransferRequestHandlerFn(cdc *codec.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		destChainID := vars["destchain"]
		bech32Addr := vars["address"]

		to, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req transferReq
		err = utils.ReadRESTReq(w, r, cdc, &req)
		if err != nil {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		info, err := kb.Get(req.BaseReq.Name)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		packet := ibc.NewIBCPacket(
			sdk.AccAddress(info.GetPubKey().Address()), to,
			req.Amount, req.BaseReq.ChainID, destChainID,
		)
		msg := ibc.IBCTransferMsg{IBCPacket: packet}

		utils.CompleteAndBroadcastTxREST(w, r, cliCtx, req.BaseReq, []sdk.Msg{msg}, cdc)
	}
}
