"use client";

import api from "@/app/api";
import API from "../api/api";
import {
  startAuthentication,
  startRegistration,
} from "@simplewebauthn/browser";
import {
  PublicKeyCredentialRequestOptionsJSON,
} from "@simplewebauthn/types";

export const PasskeyPayment = async (formData: FormData) => {
  let amount = parseFloat(formData.get("amount") as string);
  let resp = await generateAuthPasskeyPublicKey("ab@de.com", amount);
  alert("payment started")
};

const generateAuthPasskeyPublicKey = async (email: string, amount: number) => {
  const origin = window.location.origin;
  const resp = await api.post(API.PASSKEY_PAYMENT, { email, origin, amount });
  const json = resp.data.data as PublicKeyCredentialRequestOptionsJSON;
  if (json !== null) {
    const attest = await startAuthentication(json);
    const verifyResp = await api.post(
      API.PASSKEY_PAYMENT_VERIFY +
        "?origin=" +
        encodeURIComponent(origin) +
        "&email=" +
        email +
        "&amount=" +
        amount,
      attest
    );
    return verifyResp.status === 200;
  }
};
