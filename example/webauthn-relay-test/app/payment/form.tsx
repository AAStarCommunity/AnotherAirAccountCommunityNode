import { browserSupportsWebAuthn } from "@simplewebauthn/browser";
import { PayForm } from "./payform";
import { PaymentButton } from "./payment-button";
import { PasskeyPayment } from "./payment";

export default function PaymentForm() {
  return (
    <div>
      {browserSupportsWebAuthn() ? (
        <PayForm action={PasskeyPayment}>
          <PaymentButton>Signature It</PaymentButton>
          <p className="text-center text-sm text-gray-600">
            3 nodes signature tx data at least
          </p>
        </PayForm>
      ) : (
        <div>Your browser doesn&apos;t support Passkey yet</div>
      )}
    </div>
  );
}
