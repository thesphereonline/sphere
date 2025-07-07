import nacl from "tweetnacl";
import { encodeBase64, decodeBase64 } from "tweetnacl-util";

export type Wallet = {
  publicKey: string;
  privateKey: string;
};

export function generateWallet(): Wallet {
  const keyPair = nacl.sign.keyPair();
  return {
    publicKey: encodeBase64(keyPair.publicKey),
    privateKey: encodeBase64(keyPair.secretKey),
  };
}

export function signTx(privateKeyBase64: string, message: string): string {
  const secretKey = decodeBase64(privateKeyBase64);
  const signed = nacl.sign.detached(new TextEncoder().encode(message), secretKey);
  return encodeBase64(signed);
}
