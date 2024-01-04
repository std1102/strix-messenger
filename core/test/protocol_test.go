package test

import (
	"encoding/json"
	"fmt"
	"lidx-core-lib/common"
	"lidx-core-lib/keys"
	"lidx-core-lib/ratchet"
	"testing"
)

func TestProtocol(t *testing.T) {
	aKey := keys.NewInternalKeyBundle()
	bKey := keys.NewInternalKeyBundle()

	aExternalKeyBundle := aKey.GenerateExternalKey()
	bExternalKeyBundle := bKey.GenerateExternalKey()

	fmt.Print(aExternalKeyBundle)

	aKey.GenerateEphemeralKey()

	aRachet, _ := ratchet.NewRachetFromInternal(aKey, bExternalKeyBundle)

	bRachet, _ := ratchet.NewRachetFromExternal(bKey, aExternalKeyBundle, aKey.EphemeralKey.PublicKey(), aRachet.GetId())

	for i := 0; i < 10; i++ {
		msg := aRachet.PopulateMessage([]byte("THIS IS MESSAGE"))
		aRachet.OnSend(msg)

		sendedMsg := msg.ToDto()

		t.Log(sendedMsg)

		msgJson, _ := json.Marshal(sendedMsg)

		recvMsg := ratchet.CreateMessageFromJson(string(msgJson))

		bRachet.OnRecieved(recvMsg)
		t.Log(string(recvMsg.PlainMessage))
	}
}

func TestProtocol2(t *testing.T) {
	//aKey := keys.LoadInternalKey("{\n  \"identity_key\": {\n    \"public_hash\": \"fMWLlqA45tstM_vmfCnr6s1ZxItCa66n6VAEV_FmSc0\",\n    \"public_key\": \"MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEPGpNsPAahwbMdGJKAQqz5f5fOj2TQ-JoxzSPN7-54czLOK2Jy3Dq5Az_8RT6KDVYGdl2ueeof3Y3jNn421l6-dFohUCgbyFpLVa61cR6MOXUNx5u5ks2vk9ewviYM5kp\",\n    \"private_key\": \"VKOLxkY7OxshBqGrXJafEunKbzIr38JhcQSK5463_D75muWpoKgHjnvsCyCo5rLaZ6BGY8gSaHA4c6Aj1Jh9lb36hxe1vIFKsj8FblgiL-JgLn7iTaGPH5b2zps_sNlePiNIkbsZOt69yyqJcLCvIpwGy62YFgqlYY3LlrUj1rkyhVMsM9Vt7lQ_Kzo5hLaol7IcrkclWRg60AAdLnNlENpQ1I_L6asa_1vBsqNelFaq4M6WFn6gPfvWEOGM3z3-WycWyF1CbITDRItDb6xu7P-Sl5NeUcqYtWH5aQWt0nKEFWI\",\n    \"private_hash\": \"rSrgc3lKJmTksbtHjBVghGQM0GVkoKhMmrEVUWFlQMg\"\n  },\n  \"pre_key\": {\n    \"private_key\": \"TNwQjzx03hC4Vh_wBvOEhvETvO0lS-R9kcAIAlQG-It9gpXafu_SyXtXurz7nnCAAQd6hX4fdAgHGGZYlUln9oHg_lyA-9TfL4ZDF5QFCvN1FWmN8upDGZ76vLjxMkqVn3gAV8Z65sRSWRFjgVHMk8tJkCrsZn8hCfrzUAahp-W4EfjDw_TVnbT4wnFf-ESnTVjyeAq0-r2w4QZNTC-F8Vqt_wR1C4uTd6UCNBQEMqCPpV_geyXUvZexwaqXk5UtPcmzZNST7e25bFaVBoZ5NjjPgpxtnoPAABeb3vk_L-msC_A\",\n    \"private_hash\": \"2NdnQFywKpGCEVw71gb-SOiYUV2ltGDjFwzbzZXydzI\",\n    \"public_hash\": \"jDoIO6jPjOgA-tOqOMgYkYmkBakwjlzixalyYOTZbyo\",\n    \"public_key\": \"MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEqIbdkL-T8UCqnHNbobJsy5ExJYUIH3nN2afpDf4VfwSVuJXFfvDvhO4zFrrbZnabrSA_aBwa82nIUu1o3PTSD_z0ecrVWqCVRy2_yEWwluWQkrJxMRa2Pf4SxwSjEspz\"\n  },\n  \"one_time_key\": {\n    \"60449810-90e6-11ee-8f43-9af30e14a1a4\": {\n      \"public_key\": \"MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE0kevJK2oG8ygTx_AcJcEjLzQ_ttG1Xe5EiXZXNw-FmuT77Fl8hPMAjfvLjODma4gLAPa9RRCiBdR1ChCGVRw2wMmdZ3oozK9Y_jmg9Y2UKObKsWjGVD52Hm9ONzppv6o\",\n      \"private_key\": \"ZY6hedVsthcls8F1GiAYXgc_nCL-Zf4KpvHIO6b3qUZWwXNEILIhfDzgdRrbr5DyA8bO1xI20i4EuGuzhmEdZQWNDMnRsJR0wgungdvqHBnBCHGdpSrPam6cnd2RPPnGN-pAynTCAdCTCrmeJHw1vxvrFDjvuEeT59e7prpyRcE6y6AbNsius0wbLS7wODVAz1AUtL2nZ3_6GFxOwyvBU8N0EmReyNMuAWqiR-v0mOBgpjkkXZ5mZc8fADAEpKQxlnUv4b-zml3CzJUj_yhnVkRk6Tdx7Y-LDodx7zwWTcbprLY\",\n      \"private_hash\": \"nm3Z5hFTbkbABo_ToaL5qm_51ycLs0jLSNgUI7va4V0\",\n      \"public_hash\": \"3CJbLfCuV8sXXnrw8nicdbL48qLXhH_5kK7-zcJXl-E\"\n    }\n  }\n}", common.StringToByte("1234"))

	//aKey.GenerateEphemeralKey()

	aRachet := ratchet.LoadRachet("{\n      \"ephemeralKey\": \"MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE3yRHeEgU8jCj8ymEefj5kJF75-9iSuRAMr4CY64ShJzejTSQ4NXZuU9EcO217pbupeCXLwFXfENRZyYtkBy7Y-0qfbstcUq6Vht0Ofxl0sr7ccj14usVSETPOTZyFM8R\",\n      \"keyBundle\": {\n        \"identityKey\": \"MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEPGpNsPAahwbMdGJKAQqz5f5fOj2TQ-JoxzSPN7-54czLOK2Jy3Dq5Az_8RT6KDVYGdl2ueeof3Y3jNn421l6-dFohUCgbyFpLVa61cR6MOXUNx5u5ks2vk9ewviYM5kp\",\n        \"preKey\": \"MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEqIbdkL-T8UCqnHNbobJsy5ExJYUIH3nN2afpDf4VfwSVuJXFfvDvhO4zFrrbZnabrSA_aBwa82nIUu1o3PTSD_z0ecrVWqCVRy2_yEWwluWQkrJxMRa2Pf4SxwSjEspz\",\n        \"preKeySig\": \"MGUCMFLdULjTsEt0cld3MrvgrriL0HlU6tkva3w23Wf1T15svt7yAjXY-Y_juMwhXjDikQIxAP0XH1MdV7GquzM3wiceS2nKOwDbaaGGbRABbIb2NEF_0c3CxO3_xwRY0fLdogsW7g\",\n        \"oneTimeKeyId\": \"00000000-0000-0000-0000-000000000000\"\n      },\n      \"ratchetId\": \"5dadca20-90f2-11ee-97d2-40f50bb6d5eb\"\n    }", common.StringToByte("1234"))

	for i := 0; i < 10; i++ {
		//msg := aRachet.PopulateMessage([]byte("THIS IS MESSAGE"))
		//aRachet.OnSend(msg)
		//
		//sendedMsg := msg.ToDto()
		//
		//t.Log(sendedMsg)
		//
		//msgJson, _ := json.Marshal(sendedMsg)

		//recvMsg := ratchet.CreateMessageFromJson(string(msgJson))

		var msg ratchet.MessageDto

		json.Unmarshal(common.StringToByte("{\"chatSessionId\":\"5dadca20-90f2-11ee-97d2-40f50bb6d5eb\",\"index\":1,\"cipherMessage\":\"zbPTRSOKzbArJ-A6x_VgHwASGrtRhkanEDpiqxxrDVnGUZwsjY3KG-QMwv4opqUrAmvlY6Sj4yRHjtHkFjubmc0YJBuf7FQVSGmg\",\"isBinary\":false}"), &msg)

		revMsg := ratchet.CreateMessageFromDto(&msg)
		aRachet.OnRecieved(revMsg)
		t.Log(revMsg.PlainMessage)
	}
}
