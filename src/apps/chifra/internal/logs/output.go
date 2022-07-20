// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were generated with makeClass --run. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package logsPkg

// EXISTING_CODE
import (
	"net/http"

	"github.com/spf13/cobra"
)

// EXISTING_CODE

// RunLogs handles the logs command for the command line. Returns error only as per cobra.
func RunLogs(cmd *cobra.Command, args []string) (err error) {
	opts := logsFinishParse(args)
	// EXISTING_CODE
	// EXISTING_CODE
	err, _ = opts.LogsInternal()
	return
}

// ServeLogs handles the logs command for the API. Returns error and a bool if handled
func ServeLogs(w http.ResponseWriter, r *http.Request) (err error, handled bool) {
	opts := logsFinishParseApi(w, r)
	// EXISTING_CODE
	// EXISTING_CODE
	return opts.LogsInternal()
}

// LogsInternal handles the internal workings of the logs command.  Returns error and a bool if handled
func (opts *LogsOptions) LogsInternal() (err error, handled bool) {
	err = opts.validateLogs()
	if err != nil {
		return err, true
	}

	// EXISTING_CODE
	if opts.Globals.ApiMode {
		return nil, false
	}

	handled = true
	err = opts.Globals.PassItOn("getLogs", opts.Globals.Chain, opts.toCmdLine(), opts.getEnvStr())
	// EXISTING_CODE

	return
}

// EXISTING_CODE
// EXISTING_CODE
