import { Form } from "../form";
import { SubmitButton } from "../submit-button";
import { PasskeyRegister } from "../passkey";
import Link from "next/link";
import { browserSupportsWebAuthn } from "@simplewebauthn/browser";

export default function RegisterForm() {
  return (
    <div>
      {browserSupportsWebAuthn() ? (
        <Form action={PasskeyRegister}>
          <SubmitButton>Sign Up</SubmitButton>
          <p className="text-center text-sm text-gray-600">
            {"Already have an account? "}
            <Link href="/" className="font-semibold text-gray-800">
              Sign in
            </Link>
            {" instead."}
          </p>
        </Form>
      ) : (
        <div>Your browser doesn&apos;t support Passkey yet</div>
      )}
    </div>
  );
}
