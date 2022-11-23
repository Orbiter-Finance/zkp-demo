import { ethers, run } from "hardhat";
import { getGoerliFastPerGas, getChainId } from "../test/utils";
async function main() {
  await run("compile");
  const accounts = await ethers.getSigners();
  let chainId = await accounts[0].getChainId();
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
