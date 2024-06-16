import { Form } from '../form';
import { SubmitButton } from '../submit-button';
import { PasskeyRegister } from '../passkey';
import Link from 'next/link';

export default function RegisterForm() {
  return (
    <div>
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
    </div>
  );
}