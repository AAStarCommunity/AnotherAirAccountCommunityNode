"use client";

import api from "@/app/api";
import API from "../api/api";
import { startAuthentication } from "@simplewebauthn/browser";
import { PublicKeyCredentialRequestOptionsJSON } from "@simplewebauthn/types";
import { AuthData } from "../types/auth";
import { getAggr } from "./bls-helper/utility";


export const PasskeyPayment = async (formData: FormData) => {
  await isSecurePaymentConfirmationSupported();
  let txdata = formData.get("txdata") as string;
  let network = formData.get("network") as string;
  const bls: AuthData = await generateAuthPasskeyPublicKey(txdata, network);
  console.log(bls);
  const calldata = getAggr(bls.signatures, bls.pubkeys, bls.message);
  console.log(calldata);
  return calldata;
};

const generateAuthPasskeyPublicKey = async (txdata: string, network: string): Promise<AuthData | any> => {
  const origin = window.location.origin;
  const ticket = Math.floor(Math.random() * 100001).toString();
  const resp = await api.post(
    API.PASSKEY_PAYMENT,
    {
      origin,
      ticket,
      txdata: txdata,
      network: network,
    },
    {
      headers: {
        Authorization: "Bearer " + localStorage.getItem("token"),
      },
    }
  );
  const json = resp.data.data as PublicKeyCredentialRequestOptionsJSON;
  if (json !== null) {
    /*
    Practice the SPC flow
    // const req = await requestSPC(json);
    // await req.show();
    // return;
    */
    const attest = await startAuthentication(json);
    const verifyResp = await api.post(
      API.PASSKEY_PAYMENT_VERIFY +
      "?origin=" +
      encodeURIComponent(origin) +
      "&ticket=" +
      ticket +
      "&network=" +
      network,
      attest,
      {
        headers: {
          Authorization: "Bearer " + localStorage.getItem("token"),
        },
      }
    );
    if (verifyResp.status === 200) {
      alert("Signature:\n" + verifyResp.data.data.sign);
      return verifyResp.data.data.dvt;
    } else {
      alert("Signature FAILED!");
      return null;
    }
  }
};

const requestSPC = (
  json: PublicKeyCredentialRequestOptionsJSON
): PaymentRequest => {
  const { challenge, allowCredentials } = json;
  return new PaymentRequest(
    [
      {
        // Specify `secure-payment-confirmation` as payment method.
        supportedMethods: "secure-payment-confirmation",
        data: {
          // The RP ID
          rpId: "rp.example",

          // List of credential IDs obtained from the RP server.
          credentialIds: allowCredentials!.map(
            (credential) => new TextEncoder().encode(credential.id).buffer
          ),

          // The challenge is also obtained from the RP server.
          challenge: new TextEncoder().encode(challenge),

          // A display name and an icon that represent the payment instrument.
          instrument: {
            displayName: "Fancy Card ****1234",
            icon: "https://www.aastar.xyz/_next/image?url=%2Flogo.png&w=828&q=75",
            iconMustBeShown: false,
          },

          // The origin of the payee (merchant)
          payeeOrigin: "https://localhost:3000",

          // The number of milliseconds to timeout.
          timeout: 360000, // 6 minutes
        },
      },
    ],
    {
      // Payment details.
      total: {
        label: "Total",
        amount: {
          currency: "USD",
          value: "5.00",
        },
      },
    }
  );
};

const isSecurePaymentConfirmationSupported = async () => {
  if (!window.PaymentRequest) {
    return [false, "Payment Request API is not supported"];
  }

  try {
    // The data below is the minimum required to create the request and
    // check if a payment can be made.
    const supportedInstruments = [
      {
        supportedMethods: "secure-payment-confirmation",
        data: {
          // RP's hostname as its ID
          rpId: "rp.example",
          // A dummy credential ID
          credentialIds: [new Uint8Array(1)],
          // A dummy challenge
          challenge: new Uint8Array(1),
          instrument: {
            // Non-empty display name string
            displayName: " ",
            // Transparent-black pixel.
            icon: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg==",
          },
          // A dummy merchant origin
          payeeOrigin: "https://non-existent.example",
        },
      },
    ];

    const details = {
      // Dummy shopping details
      total: { label: "Total", amount: { currency: "USD", value: "0" } },
    };

    const request = new PaymentRequest(supportedInstruments, details);
    const canMakePayment = await request.canMakePayment();
    return [canMakePayment, canMakePayment ? "" : "SPC is not available"];
  } catch (error) {
    console.log(error);
    return [false, error];
  }
};

isSecurePaymentConfirmationSupported().then((result) => {
  const [isSecurePaymentConfirmationSupported, reason] = result;
  if (isSecurePaymentConfirmationSupported) {
    // Display the payment button that invokes SPC.
  } else {
    // Fallback to the legacy authentication method.
  }
});
