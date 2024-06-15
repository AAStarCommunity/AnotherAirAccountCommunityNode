import { eq } from "drizzle-orm";
import { genSaltSync, hashSync } from "bcrypt-ts";

type user = {
  id: number;
  email: string;
  password: string;
};

export async function getUser(email: string) {
  const response = await fetch("/api/user", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email }),
  });
}

export async function createUser(email: string, passkey: string) {
  return;
}
