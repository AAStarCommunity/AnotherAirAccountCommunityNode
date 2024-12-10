"use client";

import api from "@/app/api";
import API from "./api/api";
import {
  startAuthentication,
  startRegistration,
} from "@simplewebauthn/browser";
import {
  PublicKeyCredentialCreationOptionsJSON,
  PublicKeyCredentialRequestOptionsJSON,
} from "@simplewebauthn/types";

export const PasskeyRegisterViaAccount = async (formData: FormData) => {
  let account = formData.get("zuzalu-city") as string;
  let user = await generateRegPasskeyPublicKeyViaAccount(account);

  if (user != null) {
    return "User already exists"; // TODO: Handle errors with useFormStatus
  } else {
    alert("signup success");
  }
};

const generateRegPasskeyPublicKeyViaAccount = async (account: string) => {
  const origin = window.location.origin;
  const resp = await api.post(API.PASSKEY_REG_BY_ACCOUNT, {
    account,
    origin,
    type: "ZuzaluCityID",
  });
  const json = resp.data.data as PublicKeyCredentialCreationOptionsJSON;
  if (json !== null) {
    const attest = await startRegistration(json);
    const verifyResp = await api.post(
      API.PASSKEY_REG_VERIFY_BY_ACCOUNT +
        "?origin=" +
        encodeURIComponent(origin) +
        "&account=" +
        account +
        "&type=ZuzaluCityID&network=optimism-sepolia",
      attest
    );
    const signInRlt = verifyResp.status === 200 && verifyResp.data.code === 200;
    if (signInRlt) {
      if (verifyResp.data.token) {
        localStorage.setItem("token", verifyResp.data.token!);
      }
    }
  }
};
