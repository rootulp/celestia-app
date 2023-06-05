# Shares

<!-- toc -->

## Overview

All available data in a Celestia block is split into fixed-size data chunks known as "shares". A share is associated with exactly one namespace. The shares in a Celestia block are eventually erasure-coded and committed to in [Namespace Merkle trees](./data_structures.md#namespace-merkle-tree).

A share sequence is a contiguous set of shares that contain semantically relevant data. A share sequence may contain one or more shares. In most cases, a share sequence should be parsed together because the original data may have been split across share boundaries. One share sequence exists per reserved namespace and per blob.

## Share Structure

For shares **with a namespace above [`MAX_RESERVED_NAMESPACE`](./consensus.md#constants) but below [`PARITY_SHARE_NAMESPACE`](./consensus.md#constants)**:

- The first [`NAMESPACE_SIZE`](./consensus.md#constants) of a share's raw data is the namespace of that share.
- The next [`SHARE_INFO_BYTES`](./consensus.md#constants) bytes are for share information with the following structure:
  - The first 7 bits represent the share version in big endian form (initially, this will be `0000000` for version `0`);
  - The last bit is a sequence start indicator. The indicator is `1` if the share is the first share of a sequence or `0` if the share is a continuation share of a sequence.
- If this is the first share of a sequence the next [`SEQUENCE_BYTES`](./consensus.md#constants) contain a big endian `uint32` that represents the length of the sequence that follows in bytes.
- The remaining [`SHARE_SIZE`](./consensus.md#constants)`-`[`NAMESPACE_SIZE`](./consensus.md#constants)`-`[`SHARE_INFO_BYTES`](./consensus.md#constants) `-` [`SEQUENCE_BYTES`](./consensus.md#constants) bytes (if first share) or [`SHARE_SIZE`](./consensus.md#constants)`-`[`NAMESPACE_SIZE`](./consensus.md#constants)`-`[`SHARE_INFO_BYTES`](./consensus.md#constants) bytes (if continuation share) are blob data. Blob data are opaque bytes of data that are included in the block but do not impact the state. In other words, the remaining bytes have no special meaning and are simply used to store data.
- If there is insufficient blob data to fill the share, the remaining bytes are filled with `0`.

First share in a sequence:

![fig: sparse start share.](./figures/sparse_start_share.svg)

Continuation share in a sequence:

![fig: sparse continuation share.](./figures/sparse_continuation_share.svg)

Since blob data that exceeds [`SHARE_SIZE`](./consensus.md#constants)`-`[`NAMESPACE_SIZE`](./consensus.md#constants)`-`[`SHARE_INFO_BYTES`](./consensus.md#constants) `-` [`SEQUENCE_BYTES`](./consensus.md#constants) bytes will span more than one share, developers may choose to encode additional metadata in their raw blob data prior to inclusion in a Celestia block.

## Celestia Transaction Shares

In order for clients to parse shares in the middle of a sequence without downloading antecedent shares, Celestia encodes additional metadata in the shares associated with reserved namespaces. At the time of writing this only applies to the [`TRANSACTION_NAMESPACE`](./consensus.md#reserved-namespaces) and [`PAY_FOR_BLOB_NAMESPACE`](./consensus.md#reserved-namespaces). This share structure is often reffered to as "compact shares" to differentiate from the share structure defined above (a.k.a "sparse shares").

For shares **with a reserved namespace through [`NAMESPACE_ID_MAX_RESERVED`](./consensus.md#constants)**:

- The first [`NAMESPACE_SIZE`](./consensus.md#constants) of a share's raw data is the namespace of that share.
- The next [`SHARE_INFO_BYTES`](./consensus.md#constants) bytes are for share information with the following structure:
  - The first 7 bits represent the share version in big endian form (initially, this will be `0000000` for version `0`);
  - The last bit is a sequence start indicator. The indicator is `1` if the share is the first share of a sequence or `0` if the share is a continuation share of a sequence.
- If this is the first share of a sequence the next [`SEQUENCE_BYTES`](./consensus.md#constants) contain a big endian `uint32` that represents the length of the sequence that follows in bytes.
- The next [`SHARE_RESERVED_BYTES`](./consensus.md#constants) bytes are the starting byte of the length of the [canonically serialized](./consensus.md#serialization) first request that starts in the share, or `0` if there is none, as an unsigned [varint](https://developers.google.com/protocol-buffers/docs/encoding).
- The remaining [`SHARE_SIZE`](./consensus.md#constants)`-`[`NAMESPACE_SIZE`](./consensus.md#constants)`-`[`SHARE_INFO_BYTES`](./consensus.md#constants) `-` [`SEQUENCE_BYTES`](./consensus.md#constants) `-` [`SHARE_RESERVED_BYTES`](./consensus.md#constants) bytes (if first share) or [`SHARE_SIZE`](./consensus.md#constants)`-`[`NAMESPACE_SIZE`](./consensus.md#constants)`-`[`SHARE_INFO_BYTES`](./consensus.md#constants) `-` [`SHARE_RESERVED_BYTES`](./consensus.md#constants) bytes (if continuation share) are transaction or PayForBlob transaction data. Each transaction or PayForBlob transaction is prefixed with a [varint](https://developers.google.com/protocol-buffers/docs/encoding) of the length of that unit.
- If there is insufficient transaction or PayForBlob transaction data to fill the share, the remaining bytes are filled with `0`.

First share in a sequence:

![fig: compact start share.](./figures/compact_start_share.svg)

where reserved bytes would be `42` as a binary big endian `uint32` (`[0b00000000, 0b00000000, 0b00000000, 0b00101010]`).

Continuation share in a sequence:

![fig: compact continuation share.](./figures/compact_continuation_share.svg)

where reserved bytes would be `80` as a binary big endian `uint32` (`[0b00000000, 0b00000000, 0b00000000, 0b01010000]`).

## Padding

### Namespace Padding Share

A namespace padding share acts as padding between blobs so that the subsequent blob may begin at an index that conforms to the [non-interactive default rules](../rationale/data_square_layout.md#non-interactive-default-rules). A namespace padding share contains the namespace of the blob that precedes it in the data square so that the data square can retain the property that all shares are ordered by namespace.

The first [`NAMESPACE_SIZE`](./consensus.md#constants) of a share's raw data is the namespace of the blob that precedes this padding share. The next [`SHARE_INFO_BYTES`](./consensus.md#constants) bytes are for share information. The sequence start indicator is always `1`. The version bits are filled with the share version. The sequence length is zeroed out. The remaining [`SHARE_SIZE`](./consensus.md#constants)`-`[`NAMESPACE_SIZE`](./consensus.md#constants)`-`[`SHARE_INFO_BYTES`](./consensus.md#constants) `-` [`SEQUENCE_BYTES`](./consensus.md#constants) bytes are filled with `0`.

### Reserved Padding Share

Reserved padding shares are placed after the last reserved namespace share in the data square so that the first blob can start at an index that conforms to non-interactive default rules. Clients can safely ignore the contents of these shares because they don't contain any significant data.

For shares **with a namespace equal to [`RESERVED_PADDING_NAMESPACE`](./consensus.md#constants)** (i.e. reserved padding shares):

The first [`NAMESPACE_SIZE`](./consensus.md#constants) of a share's raw data `rawData` is the namespace of that share. The next [`SHARE_INFO_BYTES`](./consensus.md#constants) bytes are for share information. The sequence start indicator is always `1`. The version bits are filled with the share version. The sequence length is zeroed out. The remaining [`SHARE_SIZE`](./consensus.md#constants)`-`[`NAMESPACE_SIZE`](./consensus.md#constants)`-`[`SHARE_INFO_BYTES`](./consensus.md#constants) `-` [`SEQUENCE_BYTES`](./consensus.md#constants) bytes are filled with `0`.

### Tail Padding Share

Tail padding shares are placed after the last blob in the data square so that the number of shares in the data square is a perfect square. Clients can safely ignore the contents of these shares because they don't contain any significant data.

For shares **with a namespace equal to [`TAIL_PADDING_NAMESPACE`](./consensus.md#constants)** (i.e. tail padding shares):

The first [`NAMESPACE_SIZE`](./consensus.md#constants) of a share's raw data `rawData` is the namespace of that share, `namespaceID`. The next [`SHARE_INFO_BYTES`](./consensus.md#constants) bytes are for share information. The sequence start indicator is always `1`. The version bits are filled with the share version. The sequence length is zeroed out. The remaining [`SHARE_SIZE`](./consensus.md#constants)`-`[`NAMESPACE_SIZE`](./consensus.md#constants)`-`[`SHARE_INFO_BYTES`](./consensus.md#constants) `-` [`SEQUENCE_BYTES`](./consensus.md#constants) bytes are filled with `0`.

## Parity Share

Parity shares are the output of the erasure coding step of the data square construction process. They occupy quadrants Q1, Q2, and Q3 of the extended data square and are used to reconstruct the original data square (a.k.a Q0) in the case of a data withholding attack.

For shares **with a namespace equal to [`PARITY_SHARE_NAMESPACE`](./consensus.md#constants)** (i.e. parity shares):

Bytes carry no special meaning.
