# sygma-shared-configuration

This repository uploads shared domains configuration to IPFS (CS Storage) for each Sygma environment.

## Shared Domains Configuration

Format of shared domains configuration for all Sygma services:

```json
{
  "domains": [
    {
      "id": 0,
      "name": "",
      "type": "",
      "bridgeContract": "",
      "handlers": [
        {
          "type": "", // erc20 / erc721 / permissionedGeneric / permissionlessGeneric / xc20
          "address": ""
        }
      ],
      "nativeTokenSymbol": "",
      "nativeTokenFullName": "",
      "nativeTokenDecimals": 0,
      "blockConfirmations": 0,
      "startBlock": "",
      "resources": [
        {
          "resourceId": "",
          "type": "",
          "address": "",
          "symbol": ""
        }
      ]
    }
  ]
}
```
- domains - array of domains that the bridge is supporting
  - id - domain id
  - name - domain name
  - type - domain type (supstrate, evm)
  - bridgeContract - address of the bridge contract
  - handlers - information about the type and address of the handlers contracts
    - type - handler type (erc20, erc721, permissionedGeneric, permissionlessGeneric)
    - address - handler contract address (hex)
  - nativeTokenSymbol - string that represents symbol of the native token on the domain
  - nativeTokenFullName - string that represents full name of the native token on the domain
  - nativeTokenDecimals - decimal count of the native token on the domain
  - resources - array of resources that the bridge is supporting.
    - resourceId - any hex string
    - type - resource type (erc20, erc721, permissionlessGeneric, permissionedGeneric)
    - address - address of the contract to which the resourceID is mapped (empty if type permissionlessGeneric)
    - symbol - symbol of the resource (eg. USDC) (empty if type permissionlessGeneric)
