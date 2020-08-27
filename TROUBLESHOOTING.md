# Troubleshooting

## Installing / compiling error

If you're getting one or more error message like:
* `warning Error running install script for optional dependency:`
* `No receipt for 'com.apple.pkg.CLTools_Executables' found at '/'.`
* `gyp: No Xcode or CLT version detected!`

It probably means that your local Command Line Tools for Xcode was somehow corrupted. You should do a clean uninstall of Command Line Tools for Xcode, and then reinstall Command Line Tools for Xcode. This article should point you in the right direction: [How to resolve, No Xcode or CLT version detected!](https://medium.com/@mrjohnkilonzi/how-to-resolve-no-xcode-or-clt-version-detected-d0cf2b10a750).

## $GOPATH/bin folder missing from `PATH` env variable

On macOS, open `.zshrc` with your editor of choice from the terminal (like `sudo nano ~/.zshrc)` and add this line:
```
export PATH="$GOPATH/bin:$PATH"
```

Save your changes and then either enter `source ~/.zshrc` to make the changes effective immediatly for that specific terminal window or quit Terminal and re-open it.

## Failed continuity check (mindreader)

* **Symptom**:
  * Mindreader instance refuses to start
* **Log messages**:
  * `{"error": "continuityChecker failed: block 1911 would creates a hole after highest seen block: 1909"}`
  * `{"error": "continuityChecker already locked"}`
* **Cause**: The mindreader process missed a few blocks (probably because of an unclean shutdown) and the nodeos instance is passed that "hole". A manual restore operation is needed.
* **Solution**: Call the 'snapshot_restore' endpoint on the Mindreader manager to initiate a restore from latest snapshot, while dfuse-eos is running:

```
curl -sS -XPOST localhost:13009/v1/snapshot_restore
```

## Mindreader Head Info

You can obtain the actual head block information seen by mindreader instance using:

```bash
grpcurl -v -plaintext -d '{}' localhost:13010 dfuse.headinfo.v1.HeadInfo.GetHeadInfo
```

## Relayer Stream

You can check the blocks that goes out of the relayer component with:

```bash
grpcurl -plaintext -d '{}' localhost:13011 dfuse.bstream.v1.BlockStream.Blocks | jq .number
```

## Blockmeta Last Irreversible Block ID

You can obtain the last irreversible Block ID as seen by the blockmeta component with:

```bash
grpcurl -plaintext    -d '{}' localhost:13015 dfuse.blockmeta.v1.BlockID.LIBID
```

## Blockmeta Block Information

You can obtain information about a given block number (like if it's known to the system) by querying
blockmeta component with:

```bash
grpcurl -plaintext    -d '{"blockNum":"545"}' localhost:13015 dfuse.blockmeta.v1.BlockID.NumToID
```

## Can't find a solution?

If your issue isn't listed here, search the [issues](https://github.com/dfuse-io/dfuse-eosio/issues) section for a similar issue. If you can't find anything, open a new issue and someone from the community or the dfuse team will get to it.
