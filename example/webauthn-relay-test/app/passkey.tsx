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
  const resp = await api.post(API.PASSKEY_REG, {
    email,
    origin,
    captcha: "111111",
  });
  const json = resp.data.data as PublicKeyCredentialCreationOptionsJSON;
  if (json !== null) {
    const attest = await startRegistration(json);
    const verifyResp = await api.post(
      API.PASSKEY_REG_VERIFY +
        "?origin=" +
        encodeURIComponent(origin) +
        "&email=" +
        email +
        "&network=optimism-sepolia",
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

export const PasskeyLogin = async (formData: FormData) => {
  let resp = await generateAuthPasskeyPublicKey();

  if (resp) {
    window.location.href = "/payment";
  } else {
    alert("signin failed");
  }
};

const generateAuthPasskeyPublicKey = async () => {
  const origin = window.location.origin;
  const resp = await api.post(API.PASSKEY_AUTH, { origin });
  const json = resp.data.data as PublicKeyCredentialRequestOptionsJSON;
  if (json !== null) {
    const attest = await startAuthentication(json);
    const verifyResp = await api.post(
      API.PASSKEY_AUTH_VERIFY + "?origin=" + encodeURIComponent(origin),
      attest
    );

    const signInRlt = verifyResp.status === 200 && verifyResp.data.code === 200;
    if (signInRlt) {
      if (verifyResp.data.token) {
        localStorage.setItem("token", verifyResp.data.token!);
      }
    }
    return signInRlt;
  }
};
