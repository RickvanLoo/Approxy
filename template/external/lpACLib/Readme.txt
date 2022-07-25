Notes:
======
* Each file lists all the necessary files (sub-circuits) needed for its compilation
* Implementation for state-of-the-art circuits => paper name is mentioned. Those are:
  i. IMPACT XA (approXimate Adders)
  ii. Approx2x2MultLit and Config2x2MultLit 
  
* The following are our approximate designs:
	Approx2x2MultV1, Config2x2MultV1, Approx2x2MultV2, Approx2x2MultV3, Approx2x2MultV4

* Testbenches are provided for most circuits "testbench.vhd". May require editing for a successful run.

* pt_script.tcl is for running primetime power analysis tool. 
  Locations that require editing before a run are marked. 

* script.tcl is for running Synopsys Design Compiler for area and power analysis.
  Locations that require editing before a run are marked. 

* For questions please email Dr. Muhammad Shafique (swahlah@yahoo.com) and Walaa El-Harouni (walaa.elharouny@gmail.com)
=============================================================================================	
	
1. XA:
	|
	---- 1-bit adders: AdderAccurateOneBit, AdderIMPACTSimplifiedOneBit, 
	|				   AdderIMPACTZeroApproxOneBit, AdderIMPACTFirstApproxOneBit, 
	|				   AdderIMPACTSecondApproxOneBit, AdderIMPACTThirdApproxOneBit
	|
	---- multi-bit adders: AdderIMPACTFirstApproxMultiBit, AdderIMPACTThirdApproxMultiBit, 
						   AdderIMPACTZeroApproxMultiBit
			
						   
2. XMAA (approXimate Multiplier with Accurate Adders)
	|
	---- 2x2: Accurate2x2Mult, Approx2x2MultLit, Config2x2MultLit
	|		 Approx2x2MultV1, Config2x2MultV1
	|		 Approx2x2MultV2, Approx2x2MultV3, Approx2x2MultV4
	|		 
	---- 4x4: Accurate4x4Mult, Approx4x4MultLit, Config4x4MultLit
	|		   Approx4x4MultV1, Config4x4MultV1
	|
	---- 8x8: Accurate8x8Mult, Approx8x8MultLit, Config8x8MultLit
	|			Approx8x8MultV1, Config8x8MultV1
	|		
	---- 16x16: Accurate16x16Mult, Approx16x16MultLit, Config16x16MultLit
				Approx16x16MultV1, Config16x16MultV1	
				
				
3. SAD (Sum-of-Absolute-Differences):
	|
	---- 8x1: SAD8x1Zero, SAD8x1First, SAD8x1Third
	|
	---- 8x8: SAD8x8Zero, SAD8x8First, SAD8x8Third
	|
	---- 32x32: SAD32x32Zero, SAD32x32First, SAD32x32Third
	


4. ConfigMultLit_XA:
   Building a mult using "Config2x2MultLit" for PPs generation and using XA.
   	|
   	---- Config4x4MultLitThird: 4x4 multiplier using Config2x2MultLit and AdderIMPACTThirdApproxMultiBit
   	|
   	---- Config8x8MultLitThird: 8x8 multiplier using Config2x2MultLit and AdderIMPACTThirdApproxMultiBit
	
	

5. ConfigMultV1_XA:	
   Building a mult using "Config2x2MultV1" for PPs generation and using XA.
   These were used for testing "AMXA" (Accurate Multiplication with approXimate Adder)
   by setting the enable
   |
   ---- Config4x4MultV1First: 4x4 multiplier using Config2x2MultV1 and AdderIMPACTFirstApproxMultiBit
   |
   ---- Config4x4MultV1Third: 4x4 multiplier using Config2x2MultV1 and AdderIMPACTThirdApproxMultiBit
   |
   ---- Config4x4MultV1Zero:  4x4 multiplier using Config2x2MultV1 and AdderIMPACTZeroApproxMultiBit
   |
   ---- Config8x8MultV1First: 8x8 multiplier using Config2x2MultV1 and AdderIMPACTFirstApproxMultiBit
   |
   ---- Config8x8MultV1Third: 8x8 multiplier using Config2x2MultV1 and AdderIMPACTThirdApproxMultiBit
   |
   ---- Config8x8MultV1Zero:  8x8 multiplier using Config2x2MultV1 and AdderIMPACTZeroApproxMultiBit
   |
   ---- Config16x16MultV1Zero: 16x16 multiplier using Config2x2MultV1 and AdderIMPACTZeroApproxMultiBit
   
  
