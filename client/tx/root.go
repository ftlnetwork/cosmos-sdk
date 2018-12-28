package tx

import (
	"github.com/gorilla/mux"

	"github.com/ftlnetwork/ftlnetwork-sdk/client/context"
	"github.com/ftlnetwork/ftlnetwork-sdk/codec"
)

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/txs/{hash}", QueryTxRequestHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/txs", SearchTxRequestHandlerFn(cliCtx, cdc)).Methods("GET")
	r.HandleFunc("/txs", BroadcastTxRequest(cliCtx, cdc)).Methods("POST")
}
