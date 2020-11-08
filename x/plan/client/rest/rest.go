package rest

import (
	"encoding/base64"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

const (
	restName = "address"
)

// EncodeResp defines a tx encoding response.
type EncodeResp struct {
	Tx string `json:"tx" yaml:"tx"`
}

type DecodeReq struct {
	Tx string `json:"tx" yaml:"tx"`
}


// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/plan/profile/{%s}", restName), getHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/plan/encode", encodeTx(cliCtx)).Methods("POST")
	r.HandleFunc("/plan/decode", decodeTx(cliCtx)).Methods("POST")
}

//--------------------------------------------------------------------------------------
// Query Handlers

func getHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/plan/profile/%s", paramType), nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func encodeTx(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StdTx

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// re-encode it via the Amino wire protocol
		txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// base64 encode the encoded tx bytes
		txBytesBase64 := base64.StdEncoding.EncodeToString(txBytes)

		response := EncodeResp{Tx: txBytesBase64}
		rest.PostProcessResponse(w, cliCtx, response)
	}
}

func decodeTx(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DecodeReq
		var resp types.StdTx

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		decodedTx, err := base64.StdEncoding.DecodeString(req.Tx)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// re-encode it via the Amino wire protocol
		err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(decodedTx, &resp)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, resp)
	}
}
