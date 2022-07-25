-- "lpACLib" is a library for Low-Power Approximate Computing Modules.
-- Copyright (C) 2016, Walaa El-Harouni, Muhammad Shafique, CES, KIT.
—- E-mail: walaa.elharouny@gmail.com, swahlah@yahoo.com

-- This program is free software: you can redistribute it and/or modify
-- it under the terms of the GNU General Public License as published by
-- the Free Software Foundation, either version 3 of the License, or
-- (at your option) any later version.

-- This program is distributed in the hope that it will be useful,
-- but WITHOUT ANY WARRANTY; without even the implied warranty of
-- MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
-- GNU General Public License for more details.

-- You should have received a copy of the GNU General Public License
-- along with this program. If not, see <http://www.gnu.org/licenses/>.

--------------------------------------------
-- AdderIMPACTZeroApproxOneBit
-- Author: Jorge Castro-Godinez 
-- Implementation for approximate 1-bit FA from: 
-- "Low-Power Digital Signal Processing Using Approximate Adders"
-------------------------------------------

--! Use standard library and logic elements
library ieee;
use ieee.std_logic_1164.all;

-- renamed from approx_full_adder_purdue2
entity AdderIMPACTZeroApproxOneBit is 
    port( A     : in  std_logic; --! First adder
          B     : in  std_logic; --! Second adder
          Cin  : in  std_logic; --! Carry in
          Sum   : out std_logic; --! Resulting addition
          Cout : out std_logic  --! Carry out
    );
end AdderIMPACTZeroApproxOneBit;


--! @brief Approximate 1-bit full adder
--! c_out = (a and b) or (c_in and (a or b))
--! sum   = not c_out
architecture AdderIMPACTZeroApproxOneBitArch of AdderIMPACTZeroApproxOneBit is
	signal coutSig: std_logic;
begin
	coutSig <= (A and B) or (Cin and (A or B));
	Cout <= coutSig;
    Sum   <= not coutSig;
end AdderIMPACTZeroApproxOneBitArch;