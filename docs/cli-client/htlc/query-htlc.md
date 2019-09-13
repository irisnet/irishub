# iriscli htlc query-htlc

## Description

Query details of a htlc

## Usage

```bash
iriscli htlc query-htlc <hash-lock>
```

## Examples

### Query details of a htlc

```bash
iriscli htlc query-htlc f054e34abd9ccc3cab12a5b797b8e9c053507f279e7e53fb3f9f44d178c94b20
```

After that, you will get the detail info for the account.
```
HTLC:
	Sender:               iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f646vaym
	Receiver:             iaa1zx6n0jussc3lx0dk0rax6zsk80vgzyy7kyfud5
	ReceiverOnOtherChain: 72656365697665724f6e4f74686572436861696e
	Amount:               10000000000000000000iris-atto
	Secret:               5f5f5f6162636465666768696a6b6c6d6e6f707172737475767778797a5f5f5f
	Timestamp:            0
	ExpireHeight:         107
	State:                completed
```