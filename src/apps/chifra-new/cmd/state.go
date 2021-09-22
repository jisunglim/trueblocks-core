package cmd

/*-------------------------------------------------------------------------------------------
 * qblocks - fast, easily-accessible, fully-decentralized data from blockchains
 * copyright (c) 2016, 2021 TrueBlocks, LLC (http://trueblocks.io)
 *
 * This program is free software: you may redistribute it and/or modify it under the terms
 * of the GNU General Public License as published by the Free Software Foundation, either
 * version 3 of the License, or (at your option) any later version. This program is
 * distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even
 * the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details. You should have received a copy of the GNU General
 * Public License along with this program. If not, see http://www.gnu.org/licenses/.
 *-------------------------------------------------------------------------------------------*/

import (
	"fmt"

	"github.com/spf13/cobra"
)

// stateCmd represents the state command
var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("state called")
	},
}

func init() {
	rootCmd.AddCommand(stateCmd)
	stateCmd.SetHelpTemplate(getHelpTextState())
}

func getHelpTextState() string {
	return `chifra argc: 5 [1:state] [2:--help] [3:--verbose] [4:2] 
chifra state --help --verbose 2 
chifra state argc: 4 [1:--help] [2:--verbose] [3:2] 
chifra state --help --verbose 2 
PROG_NAME = [chifra state]

  Usage:    chifra state [-p|-c|-n|-v|-h] <address> [address...] [block...]  
  Purpose:  Retrieve account balance(s) for one or more addresses at given block(s).

  Where:
    addrs                 one or more addresses (0x...) from which to retrieve balances (required)
    blocks                an optional list of one or more blocks at which to report balances, defaults to 'latest'
    -p  (--parts <val>)   control which state to export, one or more of [none|some*|all|balance|nonce|code|storage|deployed|accttype]
    -c  (--changes)       only report a balance when it changes from one block to the next
    -n  (--no_zero)       suppress the display of zero balance accounts

    #### Hidden options
    -a  (--call <str>)    a bang-separated string consisting of address!4-byte!bytes
    #### Hidden options

    -x  (--fmt <val>)     export format, one of [none|json*|txt|csv|api]
    -v  (--verbose)       set verbose level (optional level defaults to 1)
    -h  (--help)          display this help screen

  Notes:
    - An address must start with '0x' and be forty-two characters long.
    - blocks may be a space-separated list of values, a start-end range, a special, or any combination.
    - If the queried node does not store historical state, the results are undefined.
    - special blocks are detailed under chifra when --list.
    - balance is the default mode. To select a single mode use none first, followed by that mode.
    - You may specify multiple modes on a single line.

  Powered by TrueBlocks
`
}
