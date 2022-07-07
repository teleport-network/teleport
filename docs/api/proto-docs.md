<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [bitnetwork/rvesting/v1/genesis.proto](#bitnetwork/rvesting/v1/genesis.proto)
    - [GenesisState](#bitnetwork.rvesting.v1.GenesisState)
    - [Params](#bitnetwork.rvesting.v1.Params)
  
- [bitnetwork/rvesting/v1/query.proto](#bitnetwork/rvesting/v1/query.proto)
    - [QueryParamsRequest](#bitnetwork.rvesting.v1.QueryParamsRequest)
    - [QueryParamsResponse](#bitnetwork.rvesting.v1.QueryParamsResponse)
    - [QueryRemainingRequest](#bitnetwork.rvesting.v1.QueryRemainingRequest)
    - [QueryRemainingResponse](#bitnetwork.rvesting.v1.QueryRemainingResponse)
  
    - [Query](#bitnetwork.rvesting.v1.Query)
  
- [Scalar Value Types](#scalar-value-types)



<a name="bitnetwork/rvesting/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## bitnetwork/rvesting/v1/genesis.proto



<a name="bitnetwork.rvesting.v1.GenesisState"></a>

### GenesisState
GenesisState defines the module's genesis state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#bitnetwork.rvesting.v1.Params) |  | module parameters invariant |
| `from` | [string](#string) |  |  |
| `init_reward` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="bitnetwork.rvesting.v1.Params"></a>

### Params
Params defines the rvesting module params


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enable_vesting` | [bool](#bool) |  | parameter to enable the vesting of staking reward |
| `per_block_reward` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="bitnetwork/rvesting/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## bitnetwork/rvesting/v1/query.proto



<a name="bitnetwork.rvesting.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="bitnetwork.rvesting.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#bitnetwork.rvesting.v1.Params) |  | params defines the parameters of the module. |






<a name="bitnetwork.rvesting.v1.QueryRemainingRequest"></a>

### QueryRemainingRequest







<a name="bitnetwork.rvesting.v1.QueryRemainingResponse"></a>

### QueryRemainingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `remaining` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="bitnetwork.rvesting.v1.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#bitnetwork.rvesting.v1.QueryParamsRequest) | [QueryParamsResponse](#bitnetwork.rvesting.v1.QueryParamsResponse) | Params returns the total set of parameters. | GET|/bitnetwork/rvesting/v1/params|
| `Remaining` | [QueryRemainingRequest](#bitnetwork.rvesting.v1.QueryRemainingRequest) | [QueryRemainingResponse](#bitnetwork.rvesting.v1.QueryRemainingResponse) |  | GET|/bitnetwork/rvesting/v1/remaining|

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

