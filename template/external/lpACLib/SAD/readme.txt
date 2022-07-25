* pt_script.tcl is the same for all variants

* 8x1: variants for building 8x1 SAD using three different XA: Zero, First, and Third

* 8x8: variants for building 8x8 SAD using the 8x1 SAD variants of XAs Zero, First, and Third
	+ testbench_SAD8.vhd valid for all variants, should change: (i) the component name; (ii) uut instantiation; and (iii) bitWidth  and approxBits
	+ Three DC scripts are available, each created for a single variant (Zero, First, or Third) so that the changes are minimal
	+ For changing the number of approx LSBs, edit the files and re-compile  

* 32x32: variants for building 32x32 SAD using the 8x1 SAD variants of XAs Zero, First, and Third
	+ testbench_SAD32.vhd valid for all variants, should change: (i) the component name; (ii) uut instantiation; and (iii) bitWidth  and approxBits
	+ Three DC scripts are available, each created for a single variant (Zero, First, or Third) so that the changes are minimal
	+ For changing the number of approx LSBs, edit the files and re-compile 
		
* For questions please email Dr. Muhammad Shafique (swahlah@yahoo.com) and Walaa El-Harouni (walaa.elharouny@gmail.com)