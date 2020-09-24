package cli

import (
	"time"

	eosSearch "github.com/dfuse-io/dfuse-eosio/search"
	"github.com/dfuse-io/dlauncher/launcher"
	"github.com/spf13/cobra"
)

func init() {
	launcher.RegisterCommonFlags = func(cmd *cobra.Command) error {
		// Common stores configuration flags
		cmd.Flags().String("common-backup-store-url", PitreosURL, "[COMMON] Store URL (with prefix) where to read or write backups.")
		cmd.Flags().String("common-blocks-store-url", MergedBlocksStoreURL, "[COMMON] Store URL (with prefix) where to read/write. Used by: relayer, statedb, trxdb-loader, blockmeta, search-indexer, search-live, search-forkresolver, eosws, accounthist")
		cmd.Flags().String("common-oneblock-store-url", OneBlockStoreURL, "[COMMON] Store URL (with prefix) to read/write one-block files. Used by: mindreader, merger")
		cmd.Flags().String("common-blockstream-addr", RelayerServingAddr, "gRPC endpoint to get real-time blocks. Used by: statedb, trxdb-loader, blockmeta, search-indexer, search-live, eosws, accounthist. (relayer uses its own --relayer-blockstream-addr)")

		// Network config
		cmd.Flags().String("common-network-id", NetworkID, "Short network identifier, for billing purposes (usually maps namespaces on deployments). Used by: dgraphql")
		// TODO: eventually, pluck that from somewhere instead of asking for it here (!). You risk noticing its missing very late, and it'll require reprocessing if you want the pubkeys.
		cmd.Flags().String("common-chain-id", "", "Chain ID in hex. Used by: trxdb-loader (to reverse the signatures and extract public keys)")

		// Authentication, metering and rate limiter plugins
		cmd.Flags().String("common-auth-plugin", "null://", "Auth plugin URI, see dfuse-io/dauth repository")
		cmd.Flags().String("common-metering-plugin", "null://", "Metering plugin URI, see dfuse-io/dmetering repository")
		cmd.Flags().String("common-ratelimiter-plugin", "null://", "Rate Limiter plugin URI, see dfuse-io/dauth repository")

		// Database connection strings
		cmd.Flags().String("common-trxdb-dsn", TrxDBDSN, "kvdb connection string to trxdb database. Used by: trxdb-loader, abicodec, eosws, dgraphql")

		// Service addresses
		cmd.Flags().String("common-search-addr", RouterServingAddr, "gRPC endpoint to reach the Search Router. Used by: abicodec, eosws, dgraphql")
		cmd.Flags().String("common-blockmeta-addr", BlockmetaServingAddr, "gRPC endpoint to reach the Blockmeta. Used by: search-indexer, search-router, search-live, eosws, dgraphql")

		// Filtering
		cmd.Flags().String("common-include-filter-expr", "*", "[COMMON] CEL program to determine if a given action should be included for processing purposes, can be prefixed with lowblocknum `#123;` and multiple values separated by three semi-colons `;;;`, see https://docs.dfuse.io/eosio/admin-guide/filtering/ for more information.")
		cmd.Flags().String("common-exclude-filter-expr", "", "[COMMON] CEL program to determine if an included action should be excluded, can be prefixed with lowblocknum `#123;` and multiple values separated by three semi-colons `;;;`, see https://docs.dfuse.io/eosio/admin-guide/filtering/ for more information.")
		cmd.Flags().String("common-system-actions-include-filter-expr", "receiver == 'eosio' && action in ['updateauth', 'deleteauth', 'linkauth', 'unlinkauth', 'newaccount', 'setabi']", "[COMMON] CEL program to determine which actions to keep regardless of the include or exclude filter expressions, those are actions required by dfuse system(s) to function properly, can be prefixed with lowblocknum `#123;` and multiple values separated by three semi-colons `;;;`, change it only if you known what you are doing, see https://docs.dfuse.io/eosio/admin-guide/filtering/ for more information.")

		// Search flags
		cmd.Flags().String("search-common-mesh-store-addr", "", "[COMMON] Address of the backing etcd cluster for mesh service discovery.")
		cmd.Flags().String("search-common-mesh-dsn", DmeshDSN, "[COMMON] Dmesh DSN, supports local & etcd")
		cmd.Flags().String("search-common-mesh-service-version", DmeshServiceVersion, "[COMMON] Dmesh service version (v1)")
		cmd.Flags().Duration("search-common-mesh-publish-interval", 0*time.Second, "[COMMON] How often does search archive poll dmesh")
		cmd.Flags().String("search-common-dfuse-events-action-name", "", "[COMMON] The dfuse Events action name to intercept")
		cmd.Flags().Bool("search-common-dfuse-events-unrestricted", false, "[COMMON] Flag to disable all restrictions of dfuse Events specialize indexing, for example for a private deployment")
		cmd.Flags().String("search-common-indices-store-url", IndicesStoreURL, "[COMMON] Indices path to read or write index shards Used by: search-indexer, search-archiver.")
		cmd.Flags().String("search-common-indexed-terms", eosSearch.DefaultIndexedTerms, "[COMMON] Comma separated list of terms available for indexing. These include: receiver, account, action, auth, scheduled, status, notif, input, event, ram.consumed, ram.released, db.table, db.key, data.[freeform]. Ex: 'data.from', 'data.to', they are those fields dynamically specified by smart contracts as part of their action invocations.")
		cmd.Flags().Duration("common-system-shutdown-signal-delay", 0*time.Second, "[COMMON] Add a delay between receiving SIGTERM signal and shutting down apps. 'eosws' and 'dgraphql' will respond negatively to /healthz during this period")

		return nil
	}
}
