package VHDL

import (
	"math/big"
	"testing"
)

func TestOverflowCheckGenericPass(t *testing.T) {
	//Test Max Numbers
	for i := 0; i < 64; i++ {
		// Maxnumber := uint64(math.Exp2(float64(i)) - 1)
		var expo big.Int
		expo.Exp(big.NewInt(2), big.NewInt(int64(i)), big.NewInt(0))
		Maxnumber := expo.Uint64() - 1

		output, overflow := OverflowCheckGeneric(uint(Maxnumber), uint(i))
		if overflow {
			t.Errorf("Overflow detected for N=%d, input=%d", i, Maxnumber)
		}

		if output != uint(Maxnumber) {
			t.Errorf("Output not equal for non-overflow case N=%d, input=%d, output=%d", i, Maxnumber, output)
		}

	}
}

func TestOverflowCheckGenericFail(t *testing.T) {
	//Test 1 overflow
	for i := 0; i < 64; i++ {
		//Maxnumber := uint64(math.Exp2(float64(i)))
		var expo big.Int
		expo.Exp(big.NewInt(2), big.NewInt(int64(i)), big.NewInt(0))
		Maxnumber := expo.Uint64()

		output, overflow := OverflowCheckGeneric(uint(Maxnumber), uint(i))
		if !overflow {
			t.Errorf("No Overflow detected for N=%d, input=%d", i, uint(Maxnumber))
		}

		if output != 0 {
			t.Errorf("Output not 0 for 1 above Overflow case N=%d, input=%d, output=%d", i, Maxnumber, output)
		}

	}
}
