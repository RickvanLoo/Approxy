-- "lpACLib" is a library for Low-Power Approximate Computing Modules.
-- Copyright (C) 2016, Walaa El-Harouni, Muhammad Shafique, CES, KIT.
-- E-mail: walaa.elharouny@gmail.com, swahlah@yahoo.com

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
-- AdderIMPACTThirdApproxMultiBit
-- Author: Walaa El-Harouni  
-- uses: AdderAccurateOneBit.vhd and AdderIMPACTThirdApproxOneBit.vhd
--------------------------------------------

library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;

entity AdderIMPACTThirdApproxMultiBit is
	generic (bitWidth : integer := 32; approxBits : integer := 6);
	port (A   	 : in std_logic_vector(bitWidth-1 downto 0);
		  B 	 : in std_logic_vector(bitWidth-1 downto 0);
		  Cin 	 : in std_logic;
		  Sub 	 : in std_logic; -- '0' to add and '1' to subtract
		  Sum 	 : out std_logic_vector(bitWidth-1 downto 0);
		  Cout   : out std_logic );
end AdderIMPACTThirdApproxMultiBit;

architecture AdderIMPACTThirdApproxMultiBitArch of AdderIMPACTThirdApproxMultiBit is
	component AdderAccurateOneBit is
		port (A   	 : in std_logic;
			  B 	 : in std_logic;
			  Cin 	 : in std_logic;
			  Sum 	 : out std_logic;
			  Cout   : out std_logic );
	end component;
	
	component AdderIMPACTThirdApproxOneBit is
	port (A   	 : in std_logic;
		  B 	 : in std_logic;
		  Cin 	 : in std_logic;
		  Sum 	 : out std_logic;
		  Cout   : out std_logic );
	end component;

	 signal carry_internal: std_logic_vector(bitWidth downto 0);
	 signal bIn : std_logic_vector(B'range);

begin

	init : process(B, sub)
	begin
		for i in B'range loop
			bIn(i) <= B(i) xor Sub;
		end loop;
	end process init;
	
  approx: for i in 0 to approxBits-1 generate
    approxInstances: AdderIMPACTThirdApproxOneBit
      port map (
        A  => A(i),
        B  => bIn(i),
		Sum => Sum(i),
        Cin => carry_internal(i),
        Cout => carry_internal(i+1)
      );
  end generate;
  
  accurate: for i in approxBits to bitWidth-1 generate
     accurateInstances: AdderAccurateOneBit
      port map (
        A  => A(i),
        B  => bIn(i),
		Sum => Sum(i),
        Cin => carry_internal(i),
        Cout => carry_internal(i+1)
      );

  end generate;

  carry_internal(0) <= Sub or Cin; -- either Cin is '1' for an add operation or Sub is '1' for two's complement

  Cout <= carry_internal(bitWidth);
	
end AdderIMPACTThirdApproxMultiBitArch;




