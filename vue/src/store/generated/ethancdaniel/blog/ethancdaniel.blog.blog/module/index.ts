// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { StdFee } from "@cosmjs/launchpad";
import { SigningStargateClient } from "@cosmjs/stargate";
import { Registry, OfflineSigner, EncodeObject, DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import { Api } from "./rest";
import { MsgUpdateComment } from "./types/blog/tx";
import { MsgCreateComment } from "./types/blog/tx";
import { MsgCreatePost } from "./types/blog/post";
import { MsgDeleteComment } from "./types/blog/tx";


const types = [
  ["/ethancdaniel.blog.blog.MsgUpdateComment", MsgUpdateComment],
  ["/ethancdaniel.blog.blog.MsgCreateComment", MsgCreateComment],
  ["/ethancdaniel.blog.blog.MsgCreatePost", MsgCreatePost],
  ["/ethancdaniel.blog.blog.MsgDeleteComment", MsgDeleteComment],
  
];

const registry = new Registry(<any>types);

const defaultFee = {
  amount: [],
  gas: "200000",
};

interface TxClientOptions {
  addr: string
}

interface SignAndBroadcastOptions {
  fee: StdFee,
  memo?: string
}

const txClient = async (wallet: OfflineSigner, { addr: addr }: TxClientOptions = { addr: "http://localhost:26657" }) => {
  if (!wallet) throw new Error("wallet is required");

  const client = await SigningStargateClient.connectWithSigner(addr, wallet, { registry });
  const { address } = (await wallet.getAccounts())[0];

  return {
    signAndBroadcast: (msgs: EncodeObject[], { fee=defaultFee, memo=null }: SignAndBroadcastOptions) => memo?client.signAndBroadcast(address, msgs, fee,memo):client.signAndBroadcast(address, msgs, fee),
    msgUpdateComment: (data: MsgUpdateComment): EncodeObject => ({ typeUrl: "/ethancdaniel.blog.blog.MsgUpdateComment", value: data }),
    msgCreateComment: (data: MsgCreateComment): EncodeObject => ({ typeUrl: "/ethancdaniel.blog.blog.MsgCreateComment", value: data }),
    msgCreatePost: (data: MsgCreatePost): EncodeObject => ({ typeUrl: "/ethancdaniel.blog.blog.MsgCreatePost", value: data }),
    msgDeleteComment: (data: MsgDeleteComment): EncodeObject => ({ typeUrl: "/ethancdaniel.blog.blog.MsgDeleteComment", value: data }),
    
  };
};

interface QueryClientOptions {
  addr: string
}

const queryClient = async ({ addr: addr }: QueryClientOptions = { addr: "http://localhost:1317" }) => {
  return new Api({ baseUrl: addr });
};

export {
  txClient,
  queryClient,
};
