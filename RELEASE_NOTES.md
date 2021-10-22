# Cosmos SDK v0.44.3 Release Notes

Recently, the Cosmos-SDK team became aware of a high-severity security vulnerability that impacts Cosmos-SDK v0.43.x and v0.44.x and can result in a consensus halt. User funds are NOT at risk; however, the vulnerability can result in a chain halt. This vulnerability does not impact the current Cosmos Hub, though other Cosmos-SDK based blockchains using v0.43.x or v0.44.x may be affected and are advised to update to v0.44.2 immediately.

Nodes can update their software independently of each other (no coordinated chain restart necessary), but should do so as soon as they are able.

A full disclosure will be published a week after the release.
=======
The main performance improvement concerns gRPC queries, which are now able to run concurrently on the node ([\#10045](https://github.com/cosmos/cosmos-sdk/pull/10045)). To benefit from this performance boost, make sure to send your gRPC queries to the gRPC server directly (default port `9090`) instead of using the Tendermint RPC [`abci_query` endpoint](https://docs.tendermint.com/master/rpc/#/ABCI/abci_query) (default port `26657`).

This release notably also:

- bumps Tendermint to [v0.34.14](https://github.com/tendermint/tendermint/releases/tag/v0.34.14).
- bumps the `gin-gonic/gin` version to 1.7.0 to fix the upstream [security vulnerability](https://github.com/advisories/GHSA-h395-qcrw-5vmq).
- adds a null guard with a user-friendly error message for possible nil `Amount` in tx fee `Coins`.

See the [Cosmos SDK v0.44.3 milestone](https://github.com/cosmos/cosmos-sdk/blob/v0.44.3/CHANGELOG.md) on our issue tracker for the exhaustive list of all changes.
