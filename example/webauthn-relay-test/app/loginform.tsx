import { Form } from "./form";
import { SubmitButton } from "./submit-button";
import { browserSupportsWebAuthn } from "@simplewebauthn/browser";
import { PasskeyLogin } from "./passkey";
import Link from "next/link";

export default function LoginForm() {
  return (
    <div>
      {browserSupportsWebAuthn() ? (
        <Form
          isDiscoverable={true}
          action={
            browserSupportsWebAuthn() ? PasskeyLogin : <div>not support</div>
          }
        >
          <SubmitButton>Sign in</SubmitButton>
          <p className="text-center text-sm text-gray-600">
            {"Don't have an account? "}
            <Link href="/register-account" className="font-semibold text-gray-800">
              Sign up
            </Link>
            {" for free."}
          </p>
        </Form>
      ) : (
        <div>Your browser doesn&apos;t support Passkey yet</div>
      )}
    </div>
  );
}
