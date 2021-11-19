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
package validate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Usage(msg string, values ...string) error {
	ret := msg
	for index, val := range values {
		rep := "{" + strconv.FormatInt(int64(index), 10) + "}"
		ret = strings.Replace(ret, rep, val, -1)
	}
	return errors.New(ret)
}

func Deprecated(cmd string, rep string) error {
	msg := "The {0} flag has been deprecated."
	if len(rep) > 0 {
		msg += " Use {1} instead."
	}
	return Usage(msg, cmd, rep)
}

// TODO: check if ParseUint has better performance
// TODO: TJR - I did the test - ParseUint is slower and
// TODO: TJR - does not handle hashes that are too long
func IsHex(str string) bool {
	return len(strings.Trim(str[2:], "0123456789abcdefABCDEF")) == 0
}

func Is0xPrefixed(str string) bool {
	if len(str) < 3 {
		return false
	}
	return str[:2] == "0x"
}

func IsValidHex(typ string, val string, nBytes int) (bool, error) {
	if !Is0xPrefixed(val) {
		return false, Usage("The {0} option ({1}) must {2}.", typ, val, "start with '0x'")
	} else if len(val) != (2 + nBytes*2) {
		return false, Usage("The {0} option ({1}) must {2}.", typ, val, fmt.Sprintf("be %d bytes long", nBytes))
	} else if !IsHex(val) {
		return false, Usage("The {0} option ({1}) must {2}.", typ, val, "be hex")
	}
	return true, nil
}

func IsValidFourByte(val string) (bool, error) {
	return IsValidHex("fourbyte", val, 4)
}

func IsValidTopic(val string) (bool, error) {
	return IsValidHex("topic", val, 32)
}

func IsValidAddress(val string) bool {
	if strings.Contains(val, ".eth") {
		return true
	}
	ok, _ := IsValidHex("address", val, 20)
	return ok
}

func IsValidAddressE(val string) (bool, error) {
	if strings.Contains(val, ".eth") {
		return true, nil
	}
	return IsValidHex("address", val, 20)
}

func ValidateAtLeastOneAddr(args []string) error {
	hasOne := false
	for _, arg := range args {
		if hasOne {
			break
		}
		hasOne = IsValidAddress(arg)
	}
	if hasOne {
		return nil
	}
	return Usage("Please specify at least one {0}.", "valid Ethereum address")
}

func ValidateEnum(field, value, valid string) error {
	if len(value) == 0 {
		return nil
	}
	valid = strings.Replace(valid, "[", "|", 1)
	valid = strings.Replace(valid, "]", "|", 1)
	if strings.Contains(valid, "|"+value+"|") {
		return nil
	}
	parts := strings.Split(strings.Trim(valid, "|"), "|")
	list := ""
	for _, part := range parts {
		if len(list) > 0 {
			list += " | "
		}
		list += part
	}
	return Usage("The {0} option ({1}) must be one of [ {2} ]", field, value, list)
}

func ValidateEnumSlice(field string, values []string, valid string) error {
	for _, value := range values {
		err := ValidateEnum(field, value, valid)
		if err != nil {
			return err
		}
	}
	return nil
}
