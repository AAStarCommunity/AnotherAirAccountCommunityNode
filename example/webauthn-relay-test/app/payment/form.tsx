import { browserSupportsWebAuthn } from "@simplewebauthn/browser";
import { PayForm } from "./payform";
import { PaymentButton } from "./payment-button";
import { PasskeyPayment } from "./payment";

export default function PaymentForm() {
  return (
    <div>
      {browserSupportsWebAuthn() ? (
        <PayForm action={PasskeyPayment}>
          <PaymentButton>Payment</PaymentButton>
          <p className="text-center text-sm text-gray-600">
            All your transfer will into my pocket
          </p>
        </PayForm>
      ) : (
        <div>Your browser doesn&apos;t support Passkey yet</div>
      )}
    </div>
  );
}
