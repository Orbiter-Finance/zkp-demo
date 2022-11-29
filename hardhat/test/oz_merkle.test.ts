// import "ethers";
import { StandardMerkleTree } from "@openzeppelin/merkle-tree";
import { defaultAbiCoder, keccak256 } from "ethers/lib/utils";
import { concatBytes, hexToBytes } from "ethereum-cryptography/utils";
import fs from "fs-extra";
import path from "path";

describe("oz_merkle", function () {
  before(async function () {});

  it("test", async function () {
    const abiTypes = ["address", "uint256"];
    const values = [
      ["0x1111111111111111111111111111111111111111", "5000000000000000000"],
      ["0x2222222222222222222222222222222222222222", "2500000000000000000"],
      ["0x3333333333333333333333333333333333333333", "1500000000000000000"],
    ];

    // const ed0 = defaultAbiCoder.encode(abiTypes, values[0]);
    // const kk0 = keccak256(keccak256(ed0));
    // console.warn("kk0:", kk0);

    // const ed1 = defaultAbiCoder.encode(abiTypes, values[1]);
    // const kk1 = keccak256(keccak256(ed1));
    // console.warn("kk1:", kk1);

    // const ed2 = defaultAbiCoder.encode(abiTypes, values[2]);
    // const kk2 = keccak256(keccak256(ed2));
    // console.warn("kk2:", kk2);

    // const kk0AndKk1 = keccak256(concatBytes(hexToBytes(kk2), hexToBytes(kk1)));
    // console.warn("kk0AndKk1:", kk0AndKk1);

    // const kkRoot = keccak256(
    //   concatBytes(hexToBytes(kk0AndKk1), hexToBytes(kk0))
    // );
    // console.warn("kkRoot:", kkRoot);

    const tree = StandardMerkleTree.of(values, abiTypes);

    const proof = tree.getProof(2);
    console.log("Merkle proof:", proof);

    console.log("Merkle Root:", tree.root);

    console.log(tree.render());

    await fs.writeFile(
      path.resolve(__dirname, "oz_merkle-tree.json"),
      JSON.stringify(tree.dump())
    );
  });
});
