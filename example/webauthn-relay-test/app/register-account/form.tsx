import { Form } from "../form";
import { SubmitButton } from "../submit-button";
import Link from "next/link";
import { browserSupportsWebAuthn } from "@simplewebauthn/browser";
import { ConnectButton, RainbowKitProvider } from '@rainbow-me/rainbowkit';
import { useAccount, WagmiProvider } from "wagmi";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import ConnectWallet from "../connect-wallet";
import { config } from "../wagmi";
import "@rainbow-me/rainbowkit/styles.css";
import { useEffect } from "react";
import { FormAccount } from "../form-account";
import { PasskeyRegisterViaAccount } from "../passkey-via-account";
const queryClient = new QueryClient();

export default function RegisterForm() {
  return (
    <div>
      {browserSupportsWebAuthn() ? (
        <FormAccount action={PasskeyRegisterViaAccount} isDiscoverable={false}>
          <SubmitButton>Sign Up</SubmitButton>
          <div className="flex items-center my-4">
            <hr className="flex-grow border-t border-gray-300" />
            <span className="mx-2 text-gray-600">or</span>
            <hr className="flex-grow border-t border-gray-300" />
          </div>
          <WagmiProvider config={config}>
            <QueryClientProvider client={queryClient}>
              <RainbowKitProvider>
                <ConnectWallet />
              </RainbowKitProvider>
            </QueryClientProvider>
          </WagmiProvider>
          <p className="text-center text-sm text-gray-600">
            {"Already have an account? "}
            <Link href="/" className="font-semibold text-gray-800">
              Sign in
            </Link>
            {" instead."}
          </p>
        </FormAccount>
      ) : (
        <div>Your browser doesn&apos;t support Passkey yet</div>
      )}
    </div>
  );
}
