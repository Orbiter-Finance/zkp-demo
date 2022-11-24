import { ethers } from "hardhat";
import { Signer, BigNumber, Contract, ContractFactory } from "ethers";
// import "ethers";
import { expect } from "chai";
import fs from "fs-extra";
import path from "path";

describe("mimc", function () {
  let accounts: Signer[];
  let contractFactory: ContractFactory;
  let verifierContract: Contract;

  before(async function () {
    accounts = await ethers.getSigners();
    contractFactory = await ethers.getContractFactory("Verifier");
  });

  it("deploy", async function () {
    verifierContract = await contractFactory.deploy();
  });

  it("verifyProof", async function () {
    const fpSize = 4 * 8;

    const inputBuffer = await fs.readFile(
      path.resolve(__dirname, "mimc_public_witness.input")
    );
    const input = JSON.parse(inputBuffer.toString()).Hash;
    console.warn(input);

    const proofBuffer = await fs.readFile(
      path.resolve(__dirname, "mimc.proof")
    );

    const a: Buffer[] = [];
    a[0] = proofBuffer.slice(fpSize * 0, fpSize * 1);
    a[1] = proofBuffer.slice(fpSize * 1, fpSize * 2);

    const b: Buffer[][] = [[], []];
    b[0][0] = proofBuffer.slice(fpSize * 2, fpSize * 3);
    b[0][1] = proofBuffer.slice(fpSize * 3, fpSize * 4);
    b[1][0] = proofBuffer.slice(fpSize * 4, fpSize * 5);
    b[1][1] = proofBuffer.slice(fpSize * 5, fpSize * 6);

    const c: Buffer[] = [];
    c[0] = proofBuffer.slice(fpSize * 6, fpSize * 7);
    c[1] = proofBuffer.slice(fpSize * 7, fpSize * 8);

    const resp = await verifierContract.verifyProof(a, b, c, [input]);
    expect(resp).to.be.true;
  });
});
