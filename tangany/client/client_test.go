package client

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	client := NewClient(Config{
		TanganyDefaultURL:   "https://api.tangany.com/v1/",
		TanganyClientID:     "9f78a896-80e9-4e86-b51e-9c7d1e9ebc0f",
		TanganyClientSecret: "PQK7Q~n.SFxL.KXipI0h8Y3ruk5USTgFzHUDS",
		TanganyVaultURL:     "https://cw-keyv-demo-finn.vault.azure.net",
		TanganySubscription: "845cb70213614a9abb95fe3c09d406df"})

	www, err := client.CallNamedSmartContractMethod(context.Background(), CallNamedSmartContractMethodRequest{
		TanganyEthereumNetwork: "ropsten",
		Contract:               "0xC32AE45504Ee9482db99CfA21066A59E877Bc0e6",
		Method:                 "symbol",
		Types:                  []string{},
	})
	if err != nil {
		print(err)
	}
	print(www.List[0].Type)

	qqq, err := client.CallSmartContractMethod(context.Background(), CallSmartContractMethodRequest{
		TanganyEthereumNetwork: "ropsten",
		Contract:               "0xC32AE45504Ee9482db99CfA21066A59E877Bc0e6",
		Function:               "name()",
		Outputs:                []string{"string"},
	})
	if err != nil {
		print(err)
	}
	print(qqq.List[0].Type)

	eqweq, err := client.EstimateEthereumSmartContract(context.Background(), EstimateEthereumSmartContractRequest{
		TanganyEthereumNetwork: "ropsten",
		Wallet:                 "qwe",
		Contract:               "0xe672c6ada6fBb5cB7A5c2770aE84a42E8C2106E7",
		Function:               "approve(address,uint256)",
		Inputs:                 []string{"0xab174eAb6761d6525A8A3a2E065CA042e74D0025", "0x0"},
	})
	if err != nil {
		print(err)
	}
	print(eqweq.GasPrice)

	qwewq, err := client.ExecuteEthereumSmartContract(context.Background(), ExecuteEthereumSmartContractRequest{
		TanganyEthereumNetwork: "ropsten",
		Wallet:                 "qwe",
		Contract:               "0x61B6a7b2b031Ca7053c3fD28F255AC4B17ecd5a4",
		Function:               "payme()",
		Amount:                 "0.000124",
	})
	if err != nil {
		print(err)
	}
	print(qwewq.StatusURI)

	qweq, err := client.CallWalletBasedNamedSmartContractMethod(context.Background(), CallWalletBasedNamedSmartContractMethodRequest{
		TanganyEthereumNetwork: "ropsten",
		Wallet:                 "qwe",
		Contract:               "0xC32AE45504Ee9482db99CfA21066A59E877Bc0e6",
		Method:                 "balanceOf",
		Types:                  []string{},
	})
	if err != nil {
		print(err)
	}
	print(qweq.List[0].Type)

	qwe, err := client.CallWalletBasedSmartContractMethod(context.Background(), CallWalletBasedSmartContractMethodRequest{
		TanganyEthereumNetwork: "ropsten",
		Wallet:                 "qwe",
		Contract:               "0x61B6a7b2b031Ca7053c3fD28F255AC4B17ecd5a4",
		Function:               "getVaultsWithTime()",
		Inputs:                 []string{},
		Outputs:                []string{"address[]", "uint"},
	})
	if err != nil {
		print(err)
	}
	print(qwe.List[0].Type)

	rrrr, err := client.GetERC20TokenBalance(context.Background(), GetERC20TokenBalanceRequest{
		TanganyEthereumNetwork: "ropsten",
		Wallet:                 "qwe",
		Token:                  "0xC32AE45504Ee9482db99CfA21066A59E877Bc0e6",
	})
	if err != nil {
		print(err)
	}
	print(rrrr.Balance)

	wallet := "qwe"

	resp, err := client.CreateWallet(context.Background(), CreateWalletRequest{
		// can't overwrite so add new number after each test or another wallet name
		Wallet: wallet,
		UseHsm: false,
		Tags:   []Tag{{"hasKYC": true}},
	})
	if err != nil {
		print(err)
	}
	print(resp.Version)

	//r, err := client.GetEthWalletBalance(context.Background(), GetEthWalletBalanceRequest{
	//	Wallet: wallet,
	//})
	//if err != nil {
	//	print(err)
	//}
	//print(r.Address)

	erere, err := client.EstimateTransactionFee(context.Background(), EstimateTransactionFeeRequest{
		TanganyEthereumNetwork: "ropsten",
		Wallet:                 wallet,
		Amount:                 "0.000124",
		To:                     "0xab174eAb6761d6525A8A3a2E065CA042e74D0025",
		Data:                   "0xf00ba7",
	})
	if err != nil {
		print(err)
	}
	print(erere.Fee)

	rere, err := client.GetTransactionStatus(context.Background(), GetTransactionStatusRequest{
		TanganyEthereumNetwork: "ropsten",
		Hash:                   "0xa529e851f99e7c93d40623b48a33fe6fb19a7d94c456c1cc6122a0d068fa4ca0",
	})
	if err != nil {
		print(err)
	}
	print(rere.To)

	ressp, err := client.SignPayload(context.Background(), SignPayloadRequest{
		Payload: "testPayload",
		Wallet:  wallet,
	})
	if err != nil {
		print(err)
	}
	print(ressp.Encoding)

	rresp, err := client.VerifySignature(context.Background(), VerifySignatureRequest{
		Wallet:    wallet,
		Payload:   "testPayload",
		Signature: ressp.Signature,
	})
	if err != nil {
		print(err)
	}
	print(rresp.IsValid)

	resp2, err := client.GetWallet(context.Background(), GetWalletRequest{Wallet: wallet})
	if err != nil {
		print(err)
	}

	tag := make(map[string]interface{})
	tag["op"] = "add"
	tag["path"] = "/tags/"
	tag["value"] = Tag{"isCustomer": true}

	resp4, err := client.UpdateWallet(context.Background(), UpdateWalletRequest{
		Wallet: wallet,
		Tags:   []Tag{tag},
	})
	if err != nil {
		print(err)
	}
	print(resp4)

	resp2, err = client.GetWallet(context.Background(), GetWalletRequest{Wallet: wallet})
	if err != nil {
		print(err)
	}

	resp5, err := client.ListAllWallets(context.Background(), ListAllWalletsRequest{
		Index: 0,
		Limit: 10,
		Sort:  "wallet",
	})
	if err != nil {
		print(err)
	}
	print(resp5.Hits.Hsm)

	resp3, err := client.DeleteWallet(context.Background(), DeleteWalletRequest{Wallet: wallet})
	if err != nil {
		print(err)
	}
	print(resp3.RecoveryID)

	resp2, err = client.GetWallet(context.Background(), GetWalletRequest{Wallet: wallet})
	if err != nil {
		print(err)
	}

	print(resp2.Version)
}
