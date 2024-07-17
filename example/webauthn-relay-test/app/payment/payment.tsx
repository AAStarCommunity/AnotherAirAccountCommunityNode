"use client";

import api from "@/app/api";
import API from "../api/api";
import { startAuthentication } from "@simplewebauthn/browser";
import { PublicKeyCredentialRequestOptionsJSON } from "@simplewebauthn/types";

export const PasskeyPayment = async (formData: FormData) => {
  let amount = formData.get("amount") as string;
  await generateAuthPasskeyPublicKey("ab@de.com", amount);
};

const generateAuthPasskeyPublicKey = async (email: string, amount: string) => {
  const origin = window.location.origin;
  const nonce = Math.floor(Math.random() * 100001).toString();
  const resp = await api.post(
    API.PASSKEY_PAYMENT,
    {
      origin,
      amount,
      nonce,
      txdata: "this_is_a_fake_tx_data",
    },
    {
      headers: {
        Authorization: "Bearer " + localStorage.getItem("token"),
      },
    }
  );
  const json = resp.data.data as PublicKeyCredentialRequestOptionsJSON;
  if (json !== null) {
    const attest = await startAuthentication(json);
    const verifyResp = await api.post(
      API.PASSKEY_PAYMENT_VERIFY +
        "?origin=" +
        encodeURIComponent(origin) +
        "&email=" +
        email +
        "&nonce=" +
        nonce,
      attest,
      {
        headers: {
          Authorization: "Bearer " + localStorage.getItem("token"),
        },
      }
    );
    if (verifyResp.status === 200) {
      alert("Signature:\n" + verifyResp.data.data.sign);
    } else {
      alert("Signature FAILED!");
    }
  }
};
