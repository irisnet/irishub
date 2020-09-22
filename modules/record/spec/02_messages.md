<!--
order: 2
-->

# Messages

In this section we describe the processing of the record messages and the corresponding updates to the state.

## MsgCreateRecord

A record is created using the `MsgCreateRecord` message.

```go
type MsgCreateRecord struct {
	Contents []Content
	Creator  sdk.AccAddress // the creator of the record
}
```

This message is expected to fail if:
- the length of contents is 0
- the creator is empty
- the parameters of each content are faulty, namely:
    - the `Digest` is empty
    - the `DigestAlgo` is empty