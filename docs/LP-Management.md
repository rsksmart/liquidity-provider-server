# Liquidity Provider Server - Liquidity Provider specification
The intent of this document is to explain all the points that the liquidity provider (LP) should know in order to operate his instance of the liquidity
provider server (LPS).
This document contains both technical and non-technical information, so it is recommended to be reviewed by the LP itself and the person who is in change of 
setting up the environment where the LPS is going to be deployed.

## Table of contents
- [Context](#Context)
- [LPS Configuration](#Justification)
- Minimum security requirements
- Private key management
- Other
    - Headers

## Context
In the Flyover Protocol, there are two main actors, the regular user (user), who is party interested in executing Peg-In/Peg-Out operations and the Liquidity 
Provider (LP), who puts liquidity to speed up the operation for the user in exchange for a fee as a reward. In order to do this, the user and the LP need to 
agree on the terms of the service (a Peg-In/Peg-Out *Quote*). This implies that the different LPs may offer different quotes, so the user needs to be able to
communicate with each one of the LPs to decide which one is going to use for the operation.

The user interacts with the Flyover Protocol through the [Flyover SDK](https://github.com/rsksmart/unified-bridges-sdk/tree/main/packages/flyover-sdk). This 
SDK fetches the list of the available LP from the liquidity bridge contract (LBC), this contract returns a list where each element has some information about
the LP, among this information will be the URL of the liquidity provider server (LPS) instance of that LP so the user can communicate with it. This means 
that **the LPS has an API that every user interacts with to do the quote agreement**.

The LP also needs to interact with the protocol to perform some management operations related to topics such as collateral, funds, fees management, configuration, 
etc. The LP does this operation through the LPS, that's the reason why the LPS also has an API that the LP interacts with to perform various management operations.

To summarize, the LPS has two main APIs:
- **User/Public API**: This API is used by the user to interact with the LP to agree on a quote.
- **LP/Management API**: This API is used by the LP to interact with the LPS to perform management operations.

<div align="center">
    <img src="./lp-management/img.png" alt="User fetching LP list">
</div>

If we zoom in on one LPS:

<div align="center">
    <img src="./lp-management/img_1.png" alt="Internal view of LPS">
</div>

The fact that LPS' API is divided in a public one and a private one implies that the Management API has some security requirements that need to be addressed in order
to ensure that it will be only used by the LP. Some of these measures are provided out of the box by the LPS but some others require additional configuration for the
environment where the LPS will run.

## LPS Configuration
By default, the Management API is disabled, and it can be enabled only by setting the `ENABLE_MANAGEMENT_API` environment variable to `true`. This is a security measure
to ensure that the API will only be accessible if it is explicitly enabled by the LP (or the person setting up the environment).

Once this variable is set to true **it is responsibility of the LP and the person setting up the environment to ensure that the API is properly secured**. 

TODO: complete this section with the explanation of the authentication mechanism once the LP tool epic is implemented.

## Minimum security requirements
The full detail of the endpoints and how to call them can be found in the [OpenAPI file](../OpenApi.yml) of the LPS, the following list contains a short description of each endpoint and
weather it should be treated as public or secured as a private endpoint

- PUBLIC: accessible by anyone
- PRIVATE: only accessible by LP
- ANY: is up to the administrator to set it as private or public

|        **Endpoint**        | **Method** | **Visibility** |                   **Description**                   |
|:--------------------------:|:----------:|:--------------:|:---------------------------------------------------:|
|          /health           |    GET     |      ANY       |                     Healthcheck                     |
|       /getProviders        |    GET     |     PUBLIC     |             Get list of registered LPs              |
|     /providers/details     |    GET     |     PUBLIC     |      Get details of the LP that owns this LPS       |
|      /pegin/getQuote       |    POST    |     PUBLIC     |                Get pegin quote terms                |
|     /pegin/acceptQuote     |    POST    |     PUBLIC     |              Accept pegin quote terms               |
|     /pegout/getQuotes      |    POST    |     PUBLIC     |               Get pegout quote terms                |
|    /pegout/acceptQuote     |    POST    |     PUBLIC     |              Accept pegout quote terms              |
|     /pegin/collateral      |    GET     |    PRIVATE     |        Get collateral locked by LP for pegin        |
|     /pegout/collateral     |    GET     |    PRIVATE     |       Get collateral locked by LP for pegout        |
|    /pegin/addCollateral    |    POST    |    PRIVATE     |              Lock collateral for pegin              |
|   /pegout/addCollateral    |    POST    |    PRIVATE     |             Lock collateral for pegout              |
| /pegin/withdrawCollateral  |    POST    |    PRIVATE     |        Withdraw collateral locked for pegin         |
| /pegout/withdrawCollateral |    POST    |    PRIVATE     |        Withdraw collateral locked for pegout        |
|   /provider/changeStatus   |    POST    |    PRIVATE     |     Change status of the LP that owns this LPS      |
|   /provider/resignation    |    POST    |    PRIVATE     |        Resign as flyover liquidity provider         |
|        /userQuotes         |    GET     |     PUBLIC     |     Get list of pegout deposits made by a user      |
|       /configuration       |    GET     |    PRIVATE     |          Get the configuration of this LPS          |
|       /configuration       |    POST    |    PRIVATE     |    Modify the general configuration of this LPS     |
|    /pegin/configuration    |    POST    |    PRIVATE     | Modify the pegin related configuration of this LPS  |
|   /pegout/configuration    |    POST    |    PRIVATE     | Modify the pegout related configuration of this LPS |

## Private key management
TODO: complete this section after wallet management epic is implemented to explain properly all the mechanisms of wallet management offered by LPS

The LPS performs operations in behalf of the LP during the process of the protocol, it means that it requires access to both LP's Bitcoin and
Rootstock wallets. There are three options supported in order to provide this access to the LPS:
- **Provide the private key (PK) in a file (NOT RECOMMENDED)**: in this option the LP can provide the PKs through text files, the LPS will read 
them from there and use them to sign the transactions.
- **Provide the PK through a secret management service (NOT RECOMMENDED)**: in this option the LPS will get the PK from a secret management service
and use it to sign the transactions, this is the list of supported services at the moment:
  - Amazon secrets manager
- **Use a third party wallet management service (RECOMMENDED)**: in this option the LPS won't have knowledge of any PK of the LP at any point, instead,
it will use a third party service to delegate the signing of the transactions. The services supported by the LPS right now are the following:
  - Fireblocks


