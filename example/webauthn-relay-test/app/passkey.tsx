"use client";

import { createUser, getUser } from "app/db";
import api from "@/app/api";
import API from "./api/api";
import { startRegistration } from "@simplewebauthn/browser";
import { PublicKeyCredentialCreationOptionsJSON } from "@simplewebauthn/types";

export const PasskeyRegister = async (formData: FormData) => {
  let email = formData.get("email") as string;
  let user = await generatePasskeyPublicKey(email);

  if (user != null) {
    return "User already exists"; // TODO: Handle errors with useFormStatus
  } else {
    await createUser(email, "");
    // register user with passkey successfully
    // TODO: redirect
  }
};

export const PasskeyLogin = async (formData: FormData) => {};

const generatePasskeyPublicKey = async (email: string) => {
  const origin = window.location.origin;
  const resp = await api.post(API.PASSKEY_REG, { email, origin });
  const json = resp.data.data as PublicKeyCredentialCreationOptionsJSON;
  console.log(JSON.stringify(json, null, 2));
  if (json !== null) {
    const attest = await startRegistration(json);
    const verifyResp = await api.post(
      API.PASSKEY_REG_VERIFY + "?origin=" + encodeURIComponent(origin) + "&email=" + email,
      attest
    );
    console.log(verifyResp.data);
  }
};
