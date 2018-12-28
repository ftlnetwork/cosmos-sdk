package rest

import (
	"bytes"
	"net/http"

	"github.com/ftlnetwork/ftlnetwork-sdk/client/context"
	"github.com/ftlnetwork/ftlnetwork-sdk/client/utils"
	"github.com/ftlnetwork/ftlnetwork-sdk/codec"
	"github.com/ftlnetwork/ftlnetwork-sdk/crypto/keys"
	sdk "github.com/ftlnetwork/ftlnetwork-sdk/types"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/slashing"

	"github.com/gorilla/mux"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, kb keys.Keybase) {
	r.HandleFunc(
		"/slashing/validators/{validatorAddr}/unjail",
		unjailRequestHandlerFn(cdc, kb, cliCtx),
	).Methods("POST")
}

// Unjail TX body
type UnjailReq struct {
	BaseReq utils.BaseReq `json:"base_req"`
}

func unjailRequestHandlerFn(cdc *codec.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		bech32validator := vars["validatorAddr"]

		var req UnjailReq
		err := utils.ReadRESTReq(w, r, cdc, &req)
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

		valAddr, err := sdk.ValAddressFromBech32(bech32validator)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if !bytes.Equal(info.GetPubKey().Address(), valAddr) {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "must use own validator address")
			return
		}

		msg := slashing.NewMsgUnjail(valAddr)
		utils.CompleteAndBroadcastTxREST(w, r, cliCtx, req.BaseReq, []sdk.Msg{msg}, cdc)
	}
}
