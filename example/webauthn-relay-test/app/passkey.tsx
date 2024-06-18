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

export const PasskeyRegister = async (formData: FormData) => {
  let email = formData.get("email") as string;
  let user = await generateRegPasskeyPublicKey(email);

  if (user != null) {
    return "User already exists"; // TODO: Handle errors with useFormStatus
  } else {
    alert("signup success");
  }
};

const generateRegPasskeyPublicKey = async (email: string) => {
  const origin = window.location.origin;
  const resp = await api.post(API.PASSKEY_REG, { email, origin });
  const json = resp.data.data as PublicKeyCredentialCreationOptionsJSON;
  console.log(JSON.stringify(json, null, 2));
  if (json !== null) {
    const attest = await startRegistration(json);
    const verifyResp = await api.post(
      API.PASSKEY_REG_VERIFY +
        "?origin=" +
        encodeURIComponent(origin) +
        "&email=" +
        email,
      attest
    );
    console.log(verifyResp.data);
  }
};

export const PasskeyLogin = async (formData: FormData) => {
  let email = formData.get("email") as string;
  let resp = await generateAuthPasskeyPublicKey(email);

  if (resp) {
    alert("signin success");
  } else {
    alert("signin failed")
  }
};

const generateAuthPasskeyPublicKey = async (email: string) => {
  const origin = window.location.origin;
  const resp = await api.post(API.PASSKEY_AUTH, { email, origin });
  const json = resp.data.data as PublicKeyCredentialRequestOptionsJSON;
  if (json !== null) {
    const attest = await startAuthentication(json);
    const verifyResp = await api.post(
      API.PASSKEY_AUTH_VERIFY +
        "?origin=" +
        encodeURIComponent(origin) +
        "&email=" +
        email,
      attest
    );
    return verifyResp.status === 200;
  }
};
