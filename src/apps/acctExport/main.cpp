/*-------------------------------------------------------------------------
 * This source code is confidential proprietary information which is
 * copyright (c) 2018, 2019 TrueBlocks, LLC (http://trueblocks.io)
 * All Rights Reserved
 *------------------------------------------------------------------------*/
#include "options.h"

//-----------------------------------------------------------------------
int main(int argc, const char* argv[]) {
    pinlib_init(quickQuitHandler);

    COptions options;
    if (!options.prepareArguments(argc, argv))
        return 0;

    bool once = true;
    for (auto command : options.commandLines) {
        if (!options.parseArguments(command))
            return 0;

        // clang-format off
        options.className =
            (options.count ? GETRUNTIME_CLASS(CMonitorCount)->m_ClassName
             : options.appearances
                 ? GETRUNTIME_CLASS(CAppearance)->m_ClassName
                 : (options.traces
                        ? GETRUNTIME_CLASS(CTrace)->m_ClassName
                        : (options.receipts
                               ? GETRUNTIME_CLASS(CReceipt)->m_ClassName
                               : (options.logs ? GETRUNTIME_CLASS(CLogEntry)->m_ClassName
                                               : GETRUNTIME_CLASS(CTransaction)->m_ClassName))));
        // clang-format on

        if (once)
            cout << exportPreamble(expContext().fmtMap["header"], options.className);

        options.firstOut = true;
        if (options.count) {
            options.handle_counts();

        } else if (options.appearances) {
            options.handle_appearances();

        } else if (options.receipts) {
            options.handle_receipts();

        } else if (options.traces) {
            options.handle_traces();

        } else if (options.logs) {
            options.handle_logs();

        } else {
            ASSERT(accounting);
            options.handle_accounting();
        }
        once = false;
    }

    ostringstream os;
    os << ", \"start\": " << (isTestMode() ? "\"0xdeadbeef\"" : uint_2_Str(options.scanRange.first)) << endl;
    os << ", \"end\": " << (isTestMode() ? "\"0xdeadbeef\"" : uint_2_Str(options.scanRange.second)) << endl;
    if (!options.count && options.allMonitors.size() == 1) {
        options.getNamedAccount(options.allMonitors[0], options.accountedFor);
        if (options.abi_spec.nInterfaces() == 0) {
            HIDE_FIELD(CMonitor, "abi_spec");
        }
        os << ", \"accountedFor\": " << options.allMonitors[0] << endl;
    }
    expContext().fmtMap["meta"] += os.str();

    cout << exportPostamble(options.errors, expContext().fmtMap["meta"]);

    if (!isTestMode() && !options.freshen && !options.count && !options.accounting && !options.logs) {
        ostringstream oss;
        oss << "Exported " << padNum6T(options.stats.nFileRecords) << " ";
        oss << (!options.className.empty() ? (plural(options.className) + " from ") : "of ");
        oss << padNum6T(options.stats.nFileRecords) << " transactions for address "
            << (options.allMonitors.size() ? options.accountedFor : "");
        LOG_INFO(oss.str());
    }

    if (options.stats.nFiles != options.stats.nSkipped) {
        ostringstream header;
        header << options.stats.Format(substitute(substitute(STR_DISPLAY_ACCTSCRAPESTATS, "\t", ","), "{", "{p:"));
        LOG4(header.str());
        header << endl;

        ostringstream data;
        data << options.stats.Format(substitute(STR_DISPLAY_ACCTSCRAPESTATS, "\t", ","));
        LOG4(data.str());
        data << endl;

        string_q statsFile = configPath("performance_scraper.csv");
        if (!fileExists(statsFile))
            stringToAsciiFile(statsFile, header.str());
        appendToAsciiFile(statsFile, data.str());
    }

    pinlib_cleanup();
    return 0;
}

//-----------------------------------------------------------------------
bool COptions::handle_counts(void) {
    for (auto cnt : counts) {
        if (shouldQuit())
            break;
        cout << ((isJson() && !firstOut) ? ", " : "");
        cout << cnt;
        firstOut = false;
    }
    return true;
}
