import { ethers } from "hardhat";
import { Signer, BigNumber, Contract, ContractFactory } from "ethers";
// import "ethers";
import { expect } from "chai";

describe("mimc", function () {
  let accounts: Signer[];
  let contractFactory: ContractFactory;
  let verifierContract: Contract;

  before(async function () {
    accounts = await ethers.getSigners();
    contractFactory = await ethers.getContractFactory(
      "../contracts/mimc_groth16.sol"
    );
  });

  it("deploy", async function () {
    verifierContract = await contractFactory.deploy();
  });

  it("verifyProof", async function () {});
});
