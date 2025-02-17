// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were generated with makeClass --run. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package receiptsPkg

// EXISTING_CODE
import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/internal/globals"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	outputHelpers "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output/helpers"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpcClient"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/spf13/cobra"
)

// EXISTING_CODE

// RunReceipts handles the receipts command for the command line. Returns error only as per cobra.
func RunReceipts(cmd *cobra.Command, args []string) (err error) {
	opts := receiptsFinishParse(args)
	outputHelpers.SetEnabledForCmds("receipts", opts.IsPorted())
	outputHelpers.SetWriterForCommand("receipts", &opts.Globals)
	// EXISTING_CODE
	// EXISTING_CODE
	err, _ = opts.ReceiptsInternal()
	return
}

// ServeReceipts handles the receipts command for the API. Returns error and a bool if handled
func ServeReceipts(w http.ResponseWriter, r *http.Request) (err error, handled bool) {
	opts := receiptsFinishParseApi(w, r)
	outputHelpers.SetEnabledForCmds("receipts", opts.IsPorted())
	outputHelpers.InitJsonWriterApi("receipts", w, &opts.Globals)
	// EXISTING_CODE
	// EXISTING_CODE
	err, handled = opts.ReceiptsInternal()
	outputHelpers.CloseJsonWriterIfNeededApi("receipts", err, &opts.Globals)
	return
}

// ReceiptsInternal handles the internal workings of the receipts command.  Returns error and a bool if handled
func (opts *ReceiptsOptions) ReceiptsInternal() (err error, handled bool) {
	err = opts.validateReceipts()
	if err != nil {
		return err, true
	}

	// EXISTING_CODE
	if !opts.IsPorted() {
		if opts.Globals.IsApiMode() {
			return nil, false
		}
		return opts.Globals.PassItOn("getReceipts", opts.Globals.Chain, opts.toCmdLine(), opts.getEnvStr()), true
	}

	clientVersion, err := rpcClient.GetVersion(opts.Globals.Chain)
	if err != nil {
		return err, true
	}
	erigonUsed := utils.IsClientErigon(clientVersion)

	ctx, cancel := context.WithCancel(context.Background())

	// Note: Make sure to add an entry to enabledForCmd in src/apps/chifra/pkg/output/helpers.go
	fetchTransactionData := func(modelChan chan types.Modeler[types.RawReceipt], errorChan chan error) {
		// TODO: stream transaction identifiers
		for idIndex, rng := range opts.TransactionIds {
			txList, err := rng.ResolveTxs(opts.Globals.Chain)
			if err != nil {
				errorChan <- err
				if errors.Is(err, ethereum.NotFound) {
					continue
				}
				cancel()
				return
			}

			for _, tx := range txList {
				if tx.BlockNumber < uint32(byzantiumBlockNumber) && !erigonUsed {
					err = opts.Globals.PassItOn("getReceipts", opts.Globals.Chain, getReceiptsCmdLine(opts, []string{rng.Orig}), opts.getEnvStr())
					if err != nil {
						errorChan <- err
						cancel()
						return
					}
					continue
				}

				// TODO(cache): Can this be hidden behind the GetTransactionReceipt interface. No reason
				// TODO(cache): for this calling code to know the data is in the cache.
				// Try to load receipt from cache
				// TODO(cache): We should not be sending chain here. We have enough information to fully resolve the path at this level. Send only path.
				transaction, _ := cache.GetTransaction(
					opts.Globals.Chain,
					uint64(tx.BlockNumber),
					uint64(tx.TransactionIndex),
				)

				if transaction != nil && transaction.Receipt != nil {
					// Some values are not cached
					transaction.Receipt.BlockNumber = uint64(tx.BlockNumber)
					transaction.Receipt.TransactionHash = transaction.Hash
					transaction.Receipt.TransactionIndex = uint64(tx.TransactionIndex)
					modelChan <- transaction.Receipt
					continue
				}

				// TODO: Why does this interface always accept nil and zero at the end?
				receipt, err := rpcClient.GetTransactionReceipt(opts.Globals.Chain, uint64(tx.BlockNumber), uint64(tx.TransactionIndex), nil, 0)
				if err != nil {
					if errors.Is(err, ethereum.NotFound) {
						errorChan <- fmt.Errorf("transaction at %s returned an error: %s", opts.Transactions[idIndex], ethereum.NotFound)
						continue
					}
					errorChan <- err
					cancel()
					return
				}

				modelChan <- &receipt
			}
		}
	}

	err = output.StreamMany(ctx, fetchTransactionData, output.OutputOptions{
		Writer:     opts.Globals.Writer,
		Chain:      opts.Globals.Chain,
		TestMode:   opts.Globals.TestMode,
		NoHeader:   opts.Globals.NoHeader,
		ShowRaw:    opts.Globals.ShowRaw,
		Verbose:    opts.Globals.Verbose,
		LogLevel:   opts.Globals.LogLevel,
		Format:     opts.Globals.Format,
		OutputFn:   opts.Globals.OutputFn,
		Append:     opts.Globals.Append,
		JsonIndent: "  ",
	})

	handled = true
	// EXISTING_CODE

	return
}

// GetReceiptsOptions returns the options for this tool so other tools may use it.
func GetReceiptsOptions(args []string, g *globals.GlobalOptions) *ReceiptsOptions {
	ret := receiptsFinishParse(args)
	if g != nil {
		ret.Globals = *g
	}
	return ret
}

func (opts *ReceiptsOptions) IsPorted() (ported bool) {
	// EXISTING_CODE
	ported = !opts.Articulate
	// EXISTING_CODE
	return
}

// EXISTING_CODE
// TODO: create EXISTING CODE block at the beginning of this file to keep constants ther
var byzantiumBlockNumber = 4370000

// TODO: remove this function when rewrite to Go is completed. It is only used to send
// pre-Byzantium transactions to C++ version
func getReceiptsCmdLine(opts *ReceiptsOptions, txs []string) string {
	options := ""
	if opts.Articulate {
		options += " --articulate"
	}

	options += " " + strings.Join(txs, " ")
	return options
}

// EXISTING_CODE
