package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/consensys/gnark-crypto/accumulator/merkletree"
	"github.com/consensys/gnark-crypto/ecc"
	bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/joho/godotenv"
)

type MerkleCircuit struct {
	Path, Helper []frontend.Variable
	RootHash     frontend.Variable `gnark:",public"`
}

func (circuit *MerkleCircuit) Define(api frontend.API) error {
	hFunc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	merkle.VerifyProof(api, hFunc, circuit.RootHash, circuit.Path, circuit.Helper)
	return nil
}

func randomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func main() {
	godotenv.Load(fmt.Sprintf("..%c.env", os.PathSeparator))

	var buf bytes.Buffer
	for i := 0; i < 10; i++ {
		buf.Write([]byte(randomStr(10)))
	}

	fmt.Println("buf:", buf)

	// build & verify proof for an elmt in the file
	proofIndex := uint64(5)
	segmentSize := 10
	merkleRoot, merkleProof, numLeaves, err := merkletree.BuildReaderProof(&buf, bn254.NewMiMC(), segmentSize, proofIndex)
	if err != nil {
		return
	}
	fmt.Println("merkleRoot:", merkleRoot)
	fmt.Println("numLeaves:", numLeaves)
	fmt.Println("proof: ", merkleProof)

	proofHelper := merkle.GenerateProofHelper(merkleProof, proofIndex, numLeaves)

	fmt.Printf("ProofHelper: %v\n", proofHelper)

	verified := merkletree.VerifyProof(bn254.NewMiMC(), merkleRoot, merkleProof, proofIndex, numLeaves)
	if !verified {
		fmt.Printf("The merkle proof in plain go should pass")
	}

	// create cs
	circuit := MerkleCircuit{
		Path:   make([]frontend.Variable, len(merkleProof)),
		Helper: make([]frontend.Variable, len(merkleProof)-1),
	}
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

	assignment := &MerkleCircuit{
		Path:     make([]frontend.Variable, len(merkleProof)),
		Helper:   make([]frontend.Variable, len(merkleProof)-1),
		RootHash: merkleRoot,
	}
	for i := 0; i < len(merkleProof); i++ {
		assignment.Path[i] = merkleProof[i]
	}
	for i := 0; i < len(merkleProof)-1; i++ {
		assignment.Helper[i] = proofHelper[i]
	}
	witness, _ := frontend.NewWitness(assignment, ecc.BN254)

	proof, err := groth16.Prove(r1cs, pk, witness)
	if err != nil {
		fmt.Printf("Prove failed: %v\n", err)
		return
	}

	verifySolidityPath := fmt.Sprintf("..%chardhat%ccontracts%cmerkle_groth16.sol", os.PathSeparator, os.PathSeparator, os.PathSeparator)
	f, _ := os.OpenFile(verifySolidityPath, os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()
	vk.ExportSolidity(f)

	publicWitness, _ := witness.Public()
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		fmt.Printf("Verification failed: %v\n", err)
		return
	}
	fmt.Printf("Verification succeded\n")
}
