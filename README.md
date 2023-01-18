# sygma-shared-configuration

This repository uploads shared domains configuration to IPFS (CS Storage) and publishes new entry to IPNS for each Sygma environment,

## Shared Domains Configuration

Shared domains configuration for all Sygma services:

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
          "type": "", // erc20 / erc721 / plgh / pgh / xc20
          "address": ""
        }
      ],
      "nativeTokenSymbol": "",
      "nativeTokenFullName": "",
      "nativeTokenDecimals": "",
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