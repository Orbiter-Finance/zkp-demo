package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/joho/godotenv"
)

// CubicCircuit defines a simple circuit
type CubicCircuit struct {
	PreImage frontend.Variable
	Hash     frontend.Variable `gnark:",public"`
}

// Define declares the circuit constraints
func (circuit *CubicCircuit) Define(api frontend.API) error {
	mimc, _ := mimc.NewMiMC(api)
	mimc.Write(circuit.PreImage)
	api.AssertIsEqual(circuit.Hash, mimc.Sum())
	return nil
}

func mimcHash(data []byte) string {
	f := bn254.NewMiMC()
	f.Write(data)
	hash := f.Sum(nil)
	hashInt := big.NewInt(0).SetBytes(hash)
	return hashInt.String()
}

func main() {
	godotenv.Load(fmt.Sprintf("..%c.env", os.PathSeparator))

	preImage := []byte{0x01, 0x02, 0x03}
	hash := mimcHash(preImage)

	fmt.Printf("Hash: %s\n", hash)

	var circuit CubicCircuit
	r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	if err != nil {
		fmt.Printf("Compile failed : %v\n", err)
		return
	}

	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		fmt.Printf("Setup failed\n")
		return
	}

	assignment := &CubicCircuit{PreImage: preImage, Hash: hash}
	witness, _ := frontend.NewWitness(assignment, ecc.BN254)
	proof, err := groth16.Prove(r1cs, pk, witness)
	if err != nil {
		fmt.Printf("Prove failed: %v\n", err)
		return
	}

	publicWitness, _ := witness.Public()

	f, _ := os.OpenFile("mimc_groth16.sol", os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()
	vk.ExportSolidity(f)

	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		fmt.Printf("Verification failed: %v\n", err)
		return
	}
	fmt.Printf("Verification succeded\n")
}
